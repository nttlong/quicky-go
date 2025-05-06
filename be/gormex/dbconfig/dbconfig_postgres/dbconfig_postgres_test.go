package dbconfig_postgres_test

import (
	"fmt"
	"testing"
	"vngom/gormex/dbconfig/dbconfig_postgres"

	assert "github.com/stretchr/testify/assert"
)

var yamlFile = "E:/Docker/go/quicky-go/be/gormex/postgres.yaml"
var cnnNoDb = "postgres://postgres:123456@localhost:5432/"

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
