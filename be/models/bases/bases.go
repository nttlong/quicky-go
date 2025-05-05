package bases

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`

	CreatedOn  time.Time `gorm:"index"`
	ModifiedOn time.Time `gorm:"index"`
	ModifiedBy string    `gorm:"index;type:varchar(50)"`
	CreatedBy  string    `gorm:"index;varchar(50)"`
}
