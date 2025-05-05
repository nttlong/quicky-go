package bases

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey;column:ID"`

	CreatedOn  time.Time `gorm:"index;column:CreatedOn"`
	ModifiedOn time.Time `gorm:"index;column:ModifiedOn"`
	ModifiedBy string    `gorm:"index;type:varchar(50);column:ModifiedBy"`
	CreatedBy  string    `gorm:"index;varchar(50);;column:CreatedBy"`
}
