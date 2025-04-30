// this is declare personal info of user
// The personal info are all attributes of person since they birth. such as name, age, address, phone number, e
package personal

import (
	"time"

	"quicky-go/models/bases"

	_ "github.com/jinzhu/gorm"
)

// PersonalInfo represents the personal information of a user.
type PersonalInfo struct {
	bases.BaseModel
	FirstName string `json:"first_name" gorm:"type:varchar(100);index"`
	LastName  string `json:"last_name" gorm:"type:varchar(100);index"`

	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender" gorm:"type:varchar(10);index"`
	Address     string     `json:"address" gorm:"type:varchar(255)"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(20);index"`
	Nationality string     `json:"nationality" gorm:"type:varchar(50)"`
	BirthPlace  string     `json:"birth_place" gorm:"type:varchar(100)"`
	// Add other relevant personal information fields here
}
