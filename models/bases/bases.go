package bases

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`

	CreatedOn  time.Time `gorm:"index;column:CreatedOn"`
	ModifiedOn time.Time `gorm:"index;column:ModifiedOn"`
	ModifiedBy string    `gorm:"index,length:191;column:ModifiedBy"`
	CreatedBy  string    `gorm:"index,length:191;column:CreatedBy"`
}
