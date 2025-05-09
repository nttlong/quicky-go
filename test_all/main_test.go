package main__test

import (
	"fmt"
	"testing"

	"github.com/nttlong/regorm"
	_ "github.com/nttlong/regorm"

	//"github.com/nttlong/regorm/dbconfig"
	"github.com/nttlong/regorm/expr/compiler"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	assert.True(t, compiler.ToSnakeCase("AAAA_bbbbb") == "aaaa_bbbbb")
	fmt.Print(compiler.ToSnakeCase("AAAAbbbbb"))
}
func TestLoadConfig(t *testing.T) {
	cfg := regorm.New("postgres")
	cfg.LoadFromYamlFile("./config.yaml")
	assert.True(t, cfg.GetUser() == "root")
	assert.True(t, cfg.GetPassword() == "123456")
	assert.True(t, cfg.GetHost() == "localhost")
	assert.True(t, cfg.GetPort() == "3306")

}
func TestCreateRepos(t *testing.T) {
	cfg := regorm.New("postgres")
	cfg.LoadFromYamlFile("./postgres.yaml")
	assert.Equal(t, "postgres", cfg.GetUser())
	assert.Equal(t, "123456", cfg.GetPassword())
	assert.Equal(t, "localhost", cfg.GetHost())
	assert.Equal(t, "5432", cfg.GetPort())
	err := cfg.PingDb()
	assert.True(t, err == nil)
	storage, err := cfg.GetStorage("test")
	assert.True(t, err == nil)
	assert.True(t, storage != nil)
	//list all tables
	type Table struct {
		TableName string `gorm:"type:vachar(50)"`
	}

}
