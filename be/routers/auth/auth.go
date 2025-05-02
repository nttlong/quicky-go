package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	tenant := c.Params("tenant") // Lấy giá trị của tham số tenant
	return c.SendString(fmt.Sprintf("Login for tenant: %s", tenant))

}
