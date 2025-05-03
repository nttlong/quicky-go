package helper_test

import (
	"quicky-go/pkg/db_repos/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConnectionString(t *testing.T) {
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")

	cnn, err := h.GetConnectionString()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "root:123456@tcp(localhost:3306)/", cnn)
}
func Test_GetDbConnectionString(t *testing.T) {
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")
	expectedCnn := "root:123456@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	cnn, err := h.GetDbConnectionString("test_db")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedCnn, cnn)
}
func Test_GetHelper(t *testing.T) {
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")
	h2 := helper.GetHelper("mysql")
	assert.Equal(t, h, h2)
}
func Test_Helper_Connect(t *testing.T) {
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")
	err := h.Connect()
	if err != nil {
		t.Error(err)
	}
}
func Test_CreateDatabase(t *testing.T) {
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")
	err := h.CreateDatabase("test_db")
	if err != nil {
		t.Error(err)
	}
}
func Test_GetColumnInfo(t *testing.T) {
	//use grom create model
	type TestModel struct {
		ID   int    `gorm:"column:id;primary_key"`
		Name string `gorm:"column:name;index:idx_name"`
		Age  int    `gorm:"column:age"`
	}
	type Account struct {
		Username string `gorm:"type:varchar(191);uniqueIndex:idx_username,length:191;column:Username"`
		Email    string `gorm:"type:varchar(191);uniqueIndex:idx_email,length:191;column:Email"`
		Password string `gorm:"column:Password"`
		Salt     string `json:"-" gorm:"not null;column:Salt"` // Lưu salt, không hiển thị trong JSON
	}
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")
	columns, err := h.GetColumns(&Account{})

	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 3, len(columns))
} // func Test_GetConnectionString(t *testing.T) {
func Test_GetTypeNameOfEntity(t *testing.T) {
	//use grom create model
	type TestModel struct {
		ID   int    `gorm:"column:id;primary_key"`
		Name string `gorm:"column:name;index:idx_name"`
		Age  int    `gorm:"column:age"`
	}
	type Account struct {
		Username string `gorm:"type:varchar(191);uniqueIndex:idx_username,length:191;column:Username"`
		Email    string `gorm:"type:varchar(191);uniqueIndex:idx_email,length:191;column:Email"`
		Password string `gorm:"column:Password"`
		Salt     string `json:"-" gorm:"not null;column:Salt"` // Lưu salt, không hiển thị trong JSON
	}
	helper.CreateHelper("mysql", "localhost", "3306", "root", "123456")
	h := helper.GetHelper("mysql")

	typeName := h.GetTypeNameOfEntity(&Account{})
	assert.Equal(t, "Account", typeName)

}

// 	helper := helper_mysql.HelperMysql{
// 		Host:     "localhost",
// 		Port:     "3306",
// 		User:     "root",
// 		Password: "123456",
// 	}
// 	cnn, err := helper.GetConnectionString()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.Equal(t, "root:123456@tcp(localhost:3306)/", cnn)
// }
