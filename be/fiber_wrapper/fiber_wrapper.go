package fiber_wrapper

import (
	"vngom/repo"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberAppWrapper struct {
	App    *fiber.App
	Db     *gorm.DB
	Tenant string
}

type IAppContext interface {
	GetApp() *fiber.Ctx
	GetTenant() string
	GetRepo() *repo.IRepo
}
type AppContext struct {
	App    *fiber.Ctx
	Tenant string
	Repo   *repo.IRepo
}

func (c *AppContext) GetApp() *fiber.Ctx {
	return c.App
}
func (c *AppContext) GetTenant() string {
	return c.Tenant
}
func (c *AppContext) GetRepo() *repo.IRepo {
	return c.Repo
}
func NewAppContext(app *fiber.Ctx, tenant string, rp *repo.IRepo) IAppContext {
	return &AppContext{
		App:    app,
		Tenant: tenant,
		Repo:   rp,
	}
}

type Handler func(c IAppContext) error
type Router struct {
	Method string

	Handler Handler
}
