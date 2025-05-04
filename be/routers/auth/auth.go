package auth

import (
	"fmt"
	"time"
	"vngom/fiber_wrapper"
	"vngom/models/tenants"

	"github.com/google/uuid"
)

func Login(c fiber_wrapper.IAppContext) error {
	tenantRepo, err := c.GetRepoFactory().Get("tenants")
	if err != nil {
		return err
	}
	desc := fmt.Sprint("create tenant: %s", c.GetTenant())
	dbErr := tenantRepo.Insert(&tenants.Tenants{
		ID:          uuid.New(),
		Name:        c.GetTenant(),
		Description: desc,
		Status:      1,
		DbTenant:    c.GetTenant(),
		CreatedBy:   "admin",
		ModifiedBy:  "admin",
		CreatedOn:   time.Now().UTC(),
		ModifiedOn:  time.Now().UTC(),
	})
	if dbErr != nil {
		return c.GetApp().SendString(dbErr.Error())
	}

	return c.GetApp().SendString(c.GetTenant())
	// tenant := c.Params("tenant") // Lấy giá trị của tham số tenant
	// return c.SendString(fmt.Sprintf("Login for tenant: %s", tenant))

}
