package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		di.Provide(func(cfg config.IConfig) repo.IRepoFactory {
			repoFactory := repo.NewRepoFactory()
			dbCfg := cfg.GetDBConfig()
			repoFactory.ConfigDb(
				string(dbCfg.Type),
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

		fiber_wrapper.InstallRouters(routers, app, startEnpont, cfg, repoFactory)

		app.Listen(cfg.GetServerConfig().Host + ":" + cfg.GetServerConfig().Port)

	}); err != nil {
		log.Fatal(err)
	}

}
