package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vngom/config"
	"vngom/repo"

	"vngom/fiber_wrapper"
	"vngom/routers"

	"github.com/defval/di"
	"github.com/gofiber/fiber/v2"
)

type AppInfo struct {
	CurrentDir      string
	CurrentYamlFile string
}
type GetTenantFunc func(ctx *fiber.App) string

func main() {
	//runtime.GOMAXPROCS(4)
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
			return fiber.New(fiber.Config{
				// Prefork:      true, // Sử dụng multi-core
				// ServerHeader: "Fiber",
				// ReadTimeout:  5 * time.Second,
				// WriteTimeout: 10 * time.Second,
				// IdleTimeout:  30 * time.Second,
				// Concurrency:  100000, // Tăng số connection đồng thời

			})
		}),
		di.Provide(func(appInfo AppInfo) config.IConfig {
			c := config.NewConfig()
			c.LoadConfig(appInfo.CurrentYamlFile)
			return c
		}), // provide config
		di.Provide(func() map[string]fiber_wrapper.Router {
			return routers.Routes

		}),
		di.Provide(func(cfg config.IConfig) repo.IRepoFactory {

			dbCfg := cfg.GetDBConfig()
			repoFactory := repo.NewRepoFactory(string(dbCfg.Type))
			repoFactory.ConfigDb(

				dbCfg.Host,
				dbCfg.Port,
				dbCfg.User,
				dbCfg.Password,
			)

			return repoFactory
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
		routers map[string]fiber_wrapper.Router,
		repoFactory repo.IRepoFactory,
	) {

		//decalre routes hash dict string and function

		//scan routes and add to hash map
		//add routes to app
		startEnpont := "/api/:tenant"

		app.Use(func(c *fiber.Ctx) error {
			start := time.Now()

			// Xử lý các middleware và handler tiếp theo
			err := c.Next()

			duration := time.Since(start)

			// Thêm Server-Timing header
			c.Append("Server-Timing", fmt.Sprintf("total;dur=%d", duration.Milliseconds()))

			return err
		})
		fiber_wrapper.InstallRouters(routers, app, startEnpont, cfg, repoFactory)
		app.Get("/health", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		app.Listen(cfg.GetServerConfig().Host + ":" + cfg.GetServerConfig().Port)

	}); err != nil {
		log.Fatal(err)
	}

}
