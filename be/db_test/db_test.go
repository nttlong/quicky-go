package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nttlong/regorm"
	_ "github.com/nttlong/regorm"

	//"github.com/nttlong/regorm/dbconfig"
	"vngom/models"

	"github.com/nttlong/regorm/expr/compiler"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	assert.True(t, compiler.ToSnakeCase("AAAA_bbbbb") == "aaaa_bbbbb")
	fmt.Print(compiler.ToSnakeCase("AAAAbbbbb"))
}
func TestLoadConfig(t *testing.T) {
	cfg := regorm.New("postgres")
	err := cfg.LoadFromYamlFile("./../config.yaml")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "postgres", cfg.GetUser())
	assert.Equal(t, "123456", cfg.GetPassword())
	assert.Equal(t, "localhost", cfg.GetHost())
	assert.Equal(t, "5432", cfg.GetPort())

}
func TestCreateRepos(t *testing.T) {
	cfg := regorm.New("postgres")
	err := cfg.LoadFromYamlFile("./../config.yaml")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "postgres", cfg.GetUser())
	assert.Equal(t, "123456", cfg.GetPassword())
	assert.Equal(t, "localhost", cfg.GetHost())
	assert.Equal(t, "5432", cfg.GetPort())
	err = cfg.PingDb()
	assert.True(t, err == nil)
	storage, err := cfg.GetStorage("test")
	assert.True(t, err == nil)
	assert.True(t, storage != nil)
	tananStorage, err := cfg.GetStorage("tenant")
	if err != nil {
		panic(err)
	}
	err = tananStorage.Create(models.Tenants{
		ID:          uuid.New(),
		Name:        "tanan",
		Status:      1,
		Description: "test",
		DeletedBy:   nil,
		DeletedAt:   nil,
		CreatedOn:   time.Now().UTC(),
		ModifiedOn:  nil,
	})
	assert.True(t, err == nil)

}
