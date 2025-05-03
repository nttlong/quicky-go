package department

import (
	"time"
	"vngom/models/employee"
)

type Department struct {
	Employees []employee.Employee `json:"employees" gorm:"foreignKey:DepartmentID;references:ID"`
	// ID  is auto incremented int
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement";column:ID`
	Code       string    `json:"code" gorm:"type:varchar(50);index;column:Code"`
	Name       string    `json:"firstName" gorm:"type:varchar(191);index";column:Name`
	CreatedOn  time.Time `gorm:"index;column:CreatedOn"`
	ModifiedOn time.Time `gorm:"index;column:ModifiedOn"`
	ModifiedBy string    `gorm:"index,length:191;column:ModifiedBy"`
	CreatedBy  string    `gorm:"index,length:191;column:CreatedBy"`
	Level      uint      `json:"level" gorm:"type:int;index;column:Level"`
	ParentID   uint      `json:"parentID" gorm:"type:int;index;column:ParentID"`
	//example department A has id 1 and department B has id 2,
	// department A's parentID is B, levelCode in department A is "1.2",
	// why? when get all children of department A, we just query levelCode like 1.2.*,

	LevelCode string `json:"levelCode" gorm:"type:varchar(191);index;column:LevelCode"`
	//list of employees in this department

}

func (d *Department) TableName() string {
	return "Department"
}
