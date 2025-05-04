package repo_test

import (
	"fmt"
	"testing"
	"time"
	"vngom/models/account"
	"vngom/models/bases"
	"vngom/models/department"
	"vngom/models/employee"
	"vngom/repo"

	"github.com/google/uuid"
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

// test insert data
func TestInsertData(t *testing.T) {
	repoFactory := repo.NewRepoFactory()
	repoFactory.ConfigDb(
		"mysql",
		"localhost",
		3306,
		"root",
		"123456",
	)
	err := repoFactory.PingDb()
	if err != nil {
		t.Error(err)
	}
	repoDb, err := repoFactory.Get("TestInsertData")
	if err != nil {
		t.Error(err)
	}

	d := department.Department{
		Code:       "test",
		Name:       "test",
		CreatedOn:  time.Now().UTC(),
		ModifiedOn: time.Now().UTC(),
	}

	e := repoDb.Insert(&d)
	for i := 0; i < 10; i++ {
		e = repoDb.Insert(&employee.Employee{
			BaseModel: bases.BaseModel{
				ID:         uuid.New(),
				CreatedOn:  time.Now().UTC(),
				ModifiedOn: time.Now().UTC(),
			},
			Code:         fmt.Sprintf("code-%d", i), //code-0, code-1, code-2, code-3, code-4, code-5, code-6, code-7, code-8, code-9
			FirstName:    "test",
			LastName:     "test",
			Gender:       "test",
			JoinDate:     time.Now().UTC(),
			DepartmentID: &d.ID,
			Personal:     nil,
		})
		if e != nil {
			t.Error(e)
			return
		}
	}
	if e != nil {
		t.Error(e)
	}

}
