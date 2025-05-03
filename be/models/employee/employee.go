package employee

import (
	"time"
	"vngom/models/account"
	"vngom/models/bases"
	"vngom/models/personal"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
)

type Employee struct {
	bases.BaseModel
	User      *account.Account `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Code      string           `gorm:"uniqueIndex:idx_employee_code,length:191"`
	FirstName string           `gorm:"index:idx_employee_firstname,length:191"`
	LastName  string           `gorm:"index:idx_employee_lastname,length:191"`
	Gender    string           `gorm:"index:idx_employee_gender,length:191"`

	JoinDate     time.Time
	UserID       *uuid.UUID             `gorm:"unique"` // Khóa ngoại duy nhất, có thể nil
	Personal     *personal.PersonalInfo `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DepartmentID *uint                  `gorm:"index";column:DepartmentID` // Khóa ngoại duy nhất, có thể nil
}

func (e *Employee) TableName() string {
	return "Employee"
}
