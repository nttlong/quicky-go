package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"quicky-go/pkg/config"

	"quicky-go/routers"

	"github.com/defval/di"
	"github.com/gofiber/fiber/v2"
)

type AppInfo struct {
	CurrentDir      string
	CurrentYamlFile string
}

func main() {

	// get current directory

	// load config

	di.SetTracer(&di.StdTracer{})

	// create container
	Container, err := di.New(
		di.Provide(func() context.Context {
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				stop := make(chan os.Signal)
				signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
				<-stop
				cancel()
			}()
			return ctx
		}), // provide application context
		di.Provide(func() AppInfo {
			CurrentDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			CurrentYamlFile := CurrentDir + "/config.yaml"
			return AppInfo{
				CurrentDir:      CurrentDir,
				CurrentYamlFile: CurrentYamlFile,
			}
		}), // provide application info
		di.Provide(func() *fiber.App {
			return fiber.New()
		}),
		di.Provide(func(appInfo AppInfo) config.Config {
			return config.LoadConfig(appInfo.CurrentYamlFile)
		}), // provide config
		di.Provide(func() routers.RouterMapper {
			return routers.Routes

		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	// invoke function
	if err := Container.Invoke(func(app *fiber.App, tx context.Context, cfg config.Config, routers routers.RouterMapper) {

		//decalre routes hash dict string and function

		//scan routes and add to hash map
		//add routes to app
		startEnpont := "/api/:tenant"
		for route, val := range routers {
			method := strings.ToLower(val.Method)

			if method == "get" {
				app.Get(startEnpont+route, val.Handler)
			} else if method == "post" {
				app.Post(startEnpont+route, val.Handler)
			} else if val.Method == "put" {
				app.Put(startEnpont+route, val.Handler)
			} else if method == "delete" {
				app.Delete(startEnpont+route, val.Handler)
			} else if method == "patch" {
				app.Patch(startEnpont+route, val.Handler)
			} else if method == "options" {
				app.Options(startEnpont+route, val.Handler)
			} else if method == "head" {
				app.Head(startEnpont+route, val.Handler)
			} else if method == "all" {
				app.All(startEnpont+route, val.Handler)
			}

		}

		app.Post("/", func(c *fiber.Ctx) error {

			return c.SendString("Hello, Fiber!")
		})
		app.Listen(cfg.Server.Host + ":" + cfg.Server.Port)

	}); err != nil {
		log.Fatal(err)
	}

}
