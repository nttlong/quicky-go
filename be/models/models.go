package models

import (
	"vngom/models/account"
	"vngom/models/department"
	"vngom/models/personal"
	"vngom/models/tenants"
)

type Account account.Account
type PersonalInfo personal.PersonalInfo
type Tenants tenants.TenantInfo
type Department department.Department
