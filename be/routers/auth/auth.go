package auth

import (
	"vngom/fiber_wrapper"

	"github.com/gofiber/fiber/v2"
)

func Login(c fiber_wrapper.IAppContext) error {
	c.GetApp().Status(200).JSON(fiber.Map{
		"message": "login success",
	})

}
func GetTenant(c fiber_wrapper.IAppContext) error {
	return c.GetApp().SendString(c.GetTenant())
}
