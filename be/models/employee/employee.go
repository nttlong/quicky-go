package employee

import (
	"quicky-go/models/bases"
	"quicky-go/models/user"
	"time"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
)

type Employee struct {
	bases.BaseModel
	Code      string `gorm:"uniqueIndex:idx_employee_code,length:191"`
	FirstName string
	LastName  string
	Gender    string

	JoinDate time.Time
	UserID   *uuid.UUID `gorm:"unique"` // Khóa ngoại duy nhất, có thể nil
	User     *user.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
