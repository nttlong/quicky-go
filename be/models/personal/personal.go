// this is declare personal info of user
// The personal info are all attributes of person since they birth. such as name, age, address, phone number, e
package personal

import (
	"time"

	"vngom/models/bases"

	_ "github.com/jinzhu/gorm"
)

// PersonalInfo represents the personal information of a user.
type PersonalInfo struct {
	bases.BaseModel
	FirstName string `json:"first_name" gorm:"type:varchar(100);index";column:FirstName`
	LastName  string `json:"last_name" gorm:"type:varchar(100);index";column:LastName`

	DateOfBirth *time.Time `json:"date_of_birth;";column:DateOfBirth;index`
	Gender      string     `json:"gender" gorm:"type:varchar(10);index";column:Gender`
	Address     string     `json:"address" gorm:"type:varchar(255)";column:Address"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(20);index";column:PhoneNumber`
	Nationality string     `json:"nationality" gorm:"type:varchar(50)";column:Nationality`
	BirthPlace  string     `json:"birth_place" gorm:"type:varchar(100)";column:BirthPlace`
	// Add other relevant personal information fields here
}

func (d *PersonalInfo) TableName() string {
	return "PersonalInfo"
}
