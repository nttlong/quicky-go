package config_test

import (
	"os"
	"testing"
	"vngom/config"
)

func TestConfig(t *testing.T) {
	c := config.NewConfig()
	t.Log(c)
}
func TestConfig_LoadConfig(t *testing.T) {
	currentDir, _ := os.Getwd()
	filePath := currentDir + "/../config.yaml"
	c := config.NewConfig()
	c.LoadConfig(filePath)
	t.Log(c.GetDBConfig())
	t.Log(c.GetServerConfig())
	t.Log(c)
}
