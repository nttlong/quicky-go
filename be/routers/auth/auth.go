package auth

import (
	"vngom/fiber_wrapper"
)

func Login(c fiber_wrapper.IAppContext) error {

	return c.GetApp().SendString("login")
	// tenant := c.Params("tenant") // Lấy giá trị của tham số tenant
	// return c.SendString(fmt.Sprintf("Login for tenant: %s", tenant))

}
