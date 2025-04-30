package tenants

import (
	"time"

	"github.com/google/uuid"
)

type Tenants struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey;column:ID"`
	Name        string     `gorm:"uniqueIndex:idx_tenants_name,length:191;column:Name"` // Tên của tenant
	Description string     `gorm:"type:text;column:Description"`                        // Mô tả của tenant
	Status      int        `gorm:"default:1;column:Status"`                             // Trạng thái của tenant: 1: Active, 0: Inactive
	DeletedAt   *time.Time `gorm:"index;column:DeletedAt"`                              // Thời gian xóa mềm (soft delete)
	DeletedBy   *string    `gorm:"index;column:DeletedBy"`                              // Người xóa mềm (soft delete)
	DbTenant    string     `gorm:"index;column:DbTenant"`                               // Tên của database của tenant

	CreatedOn  time.Time `gorm:"index;column:CreatedOn"`
	ModifiedOn time.Time `gorm:"index;column:ModifiedOn"`
	ModifiedBy string    `gorm:"index;column:ModifiedBy"`
	CreatedBy  string    `gorm:"index;column:CreatedBy"`
}

// TableName sets the desired table name
func (Tenants) TableName() string {
	return "Tenants"
}
