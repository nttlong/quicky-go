package fiber_wrapper

import (
	"strings"
	"vngom/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberAppWrapper struct {
	App    *fiber.App
	Db     *gorm.DB
	Tenant string
	Cfg    *config.IConfig
}

type IAppContext interface {
	GetApp() *fiber.Ctx
	GetTenant() string
	SetTenant(tenant string)

	GetConfig() config.IConfig
}
type AppContext struct {
	App    *fiber.Ctx
	Tenant string
	// Repo   repo.IRepo
	Cfg config.IConfig
	Rf  repo.IRepoFactory
}

func (c *AppContext) GetApp() *fiber.Ctx {
	return c.App
}
func (c *AppContext) GetTenant() string {
	return c.Tenant
}
func (c *AppContext) GetRepo() (repo.IRepo, error) {
	return c.Rf.Get(c.Tenant)
}
func (c *AppContext) GetConfig() config.IConfig {
	return c.Cfg
}
func (c *AppContext) SetTenant(tenant string) {
	c.Tenant = tenant
}
func (c *AppContext) GetRepoFactory() repo.IRepoFactory {
	return c.Rf
}
func NewAppContext(app *fiber.Ctx,
	tenant string,
	// rp repo.IRepo,
	cfg config.IConfig,
	rf repo.IRepoFactory) IAppContext {

	return &AppContext{
		App:    app,
		Tenant: tenant,
		// Repo:   rp,
		Cfg: cfg,
		Rf:  rf,
	}
}

type Handler func(c IAppContext) error
type Router struct {
	Method string

	Handler Handler
}

func InstallRouters(
	routers map[string]Router,
	app *fiber.App,
	startEnpont string,
	cfg config.IConfig,
	rf repo.IRepoFactory) {
	for route, val := range routers {
		method := strings.ToLower(val.Method)

		if method == "get" {
			app.Get(startEnpont+route, func(c *fiber.Ctx) error {

				tenant := c.Params("tenant")
				//r, err := rf.Get(tenant)
				// if err != nil {
				// 	return err
				// }
				appCxt := NewAppContext(c, tenant, cfg, rf)

				return val.Handler(appCxt)
			})
		}
		if method == "post" {
			app.Post(startEnpont+route, func(c *fiber.Ctx) error {
				tenant := c.Params("tenant")

				appCxt := NewAppContext(c, tenant, cfg, rf)

				return val.Handler(appCxt)
			})
		}
		if method == "put" {
			app.Put(startEnpont+route, func(c *fiber.Ctx) error {
				tenant := c.Params("tenant")

				appCxt := NewAppContext(c, tenant, cfg, rf)

				return val.Handler(appCxt)
			})
		}
		if method == "delete" {
			app.Delete(startEnpont+route, func(c *fiber.Ctx) error {
				tenant := c.Params("tenant")

				appCxt := NewAppContext(c, tenant, cfg, rf)

				return val.Handler(appCxt)
			})
		}
		if method == "patch" {
			app.Patch(startEnpont+route, func(c *fiber.Ctx) error {
				tenant := c.Params("tenant")

				appCxt := NewAppContext(c, tenant, cfg, rf)

				return val.Handler(appCxt)
			})
		}

	}
}
