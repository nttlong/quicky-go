package dbconfig_postgres_test

import (
	"fmt"
	"testing"
	"vngom/gormex/dbconfig/dbconfig_postgres"
	"vngom/gormex/dberrors"

	assert "github.com/stretchr/testify/assert"
)

var yamlFile = "E:/Docker/go/quicky-go/be/gormex/postgres.yaml"
var cnnNoDb = "postgres://postgres:123456@localhost:5432/"

type User struct {
	ID       string `gorm:"type:varchar(36);primary_key"`
	Username string `gorm:"type:varchar(50);uniqueIndex:idx_name_username"`
	Password string `gorm:"type:varchar(256);"`
}
type Working struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	StartDate string `gorm:"type:date;"`
	EndDate   string `gorm:"type:date;"`
	User      *User  `gorm:"foreignKey:ID"`
}
type Emp struct {
	ID        string     `gorm:"type:varchar(36);primary_key"`
	FirstName string     `gorm:"varchar(50);index:idx_name_firstname"`
	LastName  string     `gorm:"varchar(50);index:idx_name_lastname"`
	User      *User      `gorm:"foreignKey:ID"`
	Works     []*Working `gorm:"foreignKey:ID"`
}
type Dept struct {
	ID   string `gorm:"type:varchar(36);primary_key"`
	Name string `gorm:"type:varchar(50);uniqueIndex:idx_name_name"`
	Emps []*Emp `gorm:"foreignKey:ID"`
}

func TestNew(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)

}
func TestCreateDatabaseIfNotEx(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	err = cfg.CreateDbIfNotExist("test")
	assert.NoError(t, err)

}
func TestEntities(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	e := cfg.GetAllModelsInEntity(&Emp{})
	assert.Equal(t, 3, len(e))
	e = cfg.GetAllModelsInEntity(&Dept{})
	assert.Equal(t, 2, len(e))

}
func TestGetStorage(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	s, err := cfg.GetStorage("test")
	assert.NoError(t, err)
	fmt.Println(s)
}

func TestGetStorageAutoMigrate(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	s, err := cfg.GetStorage("test")
	assert.NoError(t, err)
	err = s.AutoMigrate(&Emp{})
	if err != nil {
		panic(err)
	}

	// check if table emp is created in db including user table
	rd := s.GetDb().Exec("SELECT * FROM emps")

	assert.NoError(t, rd.Error)
	db := s.GetDb()
	var user Emp
	err = db.Raw("SELECT * FROM emps WHERE FirstName = 'admin'").
		Scan(&user).Error
	assert.Error(t, err)
	rd = s.GetDb().Exec("SELECT * FROM users where username = ?", "admin")
	assert.NoError(t, rd.Error)
	getU := &Emp{}
	rd2 := s.GetDb().Model(&Emp{}).
		Where(`first_name = ?`, "username").
		First(&getU).Error

	assert.NoError(t, rd2)
	fmt.Println(s)
}
func TestDeleteData(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	s, err := cfg.GetStorage("test")
	assert.NoError(t, err)
	err = s.Delete(&User{}, "id = ?", "123456")
	assert.NoError(t, err)
	err = s.Delete(&User{ID: "123456"})
	assert.NoError(t, err)

}
func TestSaveData(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	s, err := cfg.GetStorage("test")
	assert.NoError(t, err)
	err = s.Delete(&User{ID: "123456"})
	assert.NoError(t, err)

	err = s.Save(&User{
		ID:       "123456",
		Username: "admin",
		Password: "123456",
	})
	assert.NoError(t, err)
	err = s.Save(&User{
		ID:       "123456",
		Username: "admin",
		Password: "123456",
	})
	assert.Error(t, err)

}
func TestTranslateError(t *testing.T) {
	cfg := dbconfig_postgres.New()
	cfg.LoadFromYamlFile(yamlFile)
	fmt.Println(cfg.GetConectionStringNoDatabase())
	assert.Equal(t, cnnNoDb, cfg.GetConectionStringNoDatabase())
	err := cfg.PingDb()
	assert.NoError(t, err)
	s, err := cfg.GetStorage("test")
	assert.NoError(t, err)
	err = s.Delete(&User{ID: "123456"})
	assert.NoError(t, err)

	err = s.Save(&User{
		ID:       "123456",
		Username: "admin",
		Password: "123456",
	})
	assert.NoError(t, err)
	err = s.Save(&User{
		ID:       "123456",
		Username: "admin",
		Password: "123456",
	})
	ft := s.GetDbConfig().TranslateError(err, &User{}, "save")
	assert.Equal(t, dberrors.Duplicate, ft.Code)
	assert.Equal(t, "save", ft.Action)
	assert.Equal(t, 1, len(ft.RefColumns))
	assert.Equal(t, "id", ft.RefColumns[0])

}
