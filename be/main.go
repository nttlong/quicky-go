package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"vngom/config"

	"vngom/fiber_wrapper"
	"vngom/routers"

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
		di.Provide(func(appInfo AppInfo) config.IConfig {
			c := config.NewConfig()
			c.LoadConfig(appInfo.CurrentYamlFile)
			return c
		}), // provide config
		di.Provide(func() map[string]fiber_wrapper.Router {
			return routers.Routes

		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	// invoke function
	if err := Container.Invoke(func(
		app *fiber.App,
		tx context.Context,
		cfg config.IConfig,
		routers map[string]fiber_wrapper.Router) {

		//decalre routes hash dict string and function

		//scan routes and add to hash map
		//add routes to app
		startEnpont := "/api/:tenant"

		for route, val := range routers {
			method := strings.ToLower(val.Method)

			if method == "get" {
				app.Get(startEnpont+route, func(c *fiber.Ctx) error {
					tenant := c.Params("tenant")
					a := fiber_wrapper.NewAppContext(c, tenant, nil)
					return val.Handler(a)
				})
			}

		}

		app.Get("/", func(c *fiber.Ctx) error {

			return c.SendString("Hello, Fiber!")
		})
		app.Listen(cfg.GetServerConfig().Host + ":" + cfg.GetServerConfig().Port)

	}); err != nil {
		log.Fatal(err)
	}

}
