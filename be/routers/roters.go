package routers

import (
	"quicky-go/routers/auth"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	Method  string
	Path    string
	Handler func(*fiber.Ctx) error
}
type RouterMapper map[string]Route

var Routes RouterMapper = make(RouterMapper)

func GetRouter() RouterMapper {
	return Routes
}
func init() {
	Routes["/auth/login"] = Route{
		Method:  "GET",
		Path:    "/auth/login",
		Handler: auth.Login,
	}
}
