package repo_test

import (
	"fmt"
	"testing"
	"vngom/models/account"
	"vngom/models/department"
	"vngom/repo"
)

func TestRepoFactory(t *testing.T) {
	repoFactory := repo.NewRepoFactory()
	repoFactory.ConfigDb(
		"mysql",
		"localhost",
		3306,
		"root",
		"123456",
	)
	repoFactory.PingDb()

}
func TestGetFullEntityNameOfRepoFactory(t *testing.T) {
	repoFactory := repo.NewRepoFactory()

	f := repoFactory.GetFullEntityName(&account.Account{})
	fmt.Println(f)
	t.Log(f)
}
func TestGetColumOfRepoFactory(t *testing.T) {
	repoFactory := repo.NewRepoFactory()
	type BaseEntity struct {
		Code   string `gorm:"type:varchar(191);uniqueIndex:idx_code,length:191;column:Code"`
		Name   string `gorm:"type:varchar(191);column:Name"`
		Status int    `gorm:"column:Status"`
	}
	type Account struct {
		BaseEntity
		Username string `gorm:"type:varchar(191);uniqueIndex:idx_username,length:191;column:Username"`
		Email    string `gorm:"type:varchar(191);uniqueIndex:idx_email,length:191;column:Email"`
		Password string `gorm:"column:Password"`
		Salt     string `json:"-" gorm:"not null;column:Salt"` // Lưu salt, không hiển thị trong JSON
	}
	f, e := repoFactory.GetColumOfEntity(&Account{})
	if e != nil {
		t.Error(e)
	}
	fmt.Println(f)
	t.Log(f)
}
func TestGetRepoFromRepoFactory(t *testing.T) {
	repoFactory := repo.NewRepoFactory()
	repoFactory.ConfigDb(
		"mysql",
		"localhost",
		3306,
		"root",
		"123456",
	)
	repoFactory.PingDb()
	repo, e := repoFactory.Get("test")
	if e != nil {
		t.Error(e)
	}
	t.Log(repo)
}
func TestAutomigrateEntry(t *testing.T) {
	repoFactory := repo.NewRepoFactory()
	repoFactory.ConfigDb(
		"mysql",
		"localhost",
		3306,
		"root",
		"123456",
	)
	repoFactory.PingDb()
	repoDb, err := repoFactory.Get("TestAutomigrateEntry")
	if err != nil {
		t.Error(err)
	}

	e := repoDb.AutoMigrate(&department.Department{})
	if e != nil {
		t.Error(e)
	}
	t.Log(e)
}
