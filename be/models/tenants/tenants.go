package tenants

import (
	"time"

	"github.com/google/uuid"
)

type TenantInfo struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Name        string     `gorm:"type:char(191);uniqueIndex:idx_name"`        // Tên của tenant
	Description string     `gorm:"type:text"`                                  // Mô tả của tenant
	Status      int        `gorm:"default:1;column:Status"`                    // Trạng thái của tenant: 1: Active, 0: Inactive
	DeletedAt   *time.Time `gorm:"index;column:DeletedAt"`                     // Thời gian xóa mềm (soft delete)
	DeletedBy   *string    `gorm:"index;column:DeletedBy"`                     // Người xóa mềm (soft delete)
	DbTenant    string     `gorm:"uniqueIndex:idx_db_tenants_name,length:191"` // Tên của database của tenant

	CreatedOn  time.Time  `gorm:"index"`
	ModifiedOn *time.Time `gorm:"index"`
	ModifiedBy string     `gorm:"index"`
	CreatedBy  string     `gorm:"index"`
}
