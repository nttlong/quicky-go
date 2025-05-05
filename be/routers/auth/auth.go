package auth

import (
	"fmt"
	"time"
	"vngom/fiber_wrapper"
	"vngom/models/tenants"

	"vngom/manager/accounts"

	_ "vngom/repo/repo_types"

	"github.com/google/uuid"
)

func Login(c fiber_wrapper.IAppContext) error {
	//check time
	starAt := time.Now().UTC()
	tenantRepo, err := c.GetRepoFactory().Get("tenants")
	cdb, err := c.GetRepo()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	elapseTime := time.Now().UTC().Sub(starAt)
	fmt.Println("Login time: ", elapseTime.Milliseconds())
	desc := fmt.Sprint("create tenant: %s", c.GetTenant())
	starAt = time.Now().UTC()
	_ = tenantRepo.Insert(&tenants.Tenants{
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
	elapseTime = time.Now().UTC().Sub(starAt)
	fmt.Println("Login time in ms: ", elapseTime.Milliseconds())
	// if dbErr != nil {
	// 	return c.GetApp().SendString(dbErr.Error())
	// }

	mgg := accounts.NewAccountManager(cdb)
	elapseTime = time.Now().UTC().Sub(starAt)

	acc, errv := mgg.CreateAccount("admin", "admin@localhost", "123456")
	elapseTime = time.Now().UTC().Sub(starAt)
	fmt.Println("CreateAccount time in ms: ", elapseTime.Milliseconds())

	//nacc, errv := mgg.ValidateAccount("admin", "123456")
	if errv != nil {
		return c.GetApp().SendString(errv.Error())
	} else {
		//return json.Marshal(nacc)
		return c.GetApp().JSON(acc)
	}

	// if errDb != nil {
	// 	return c.GetApp().SendString(errDb.Error())
	// }
	// return c.GetApp().SendString(c.GetTenant())

	// tenant := c.Params("tenant") // Lấy giá trị của tham số tenant
	// return c.SendString(fmt.Sprintf("Login for tenant: %s", tenant))

}
func GetTenant(c fiber_wrapper.IAppContext) error {
	return c.GetApp().SendString(c.GetTenant())
}
