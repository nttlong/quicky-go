package main_test

import (
	"testing"

	"fmt"
	"quicky-go/pkg/config_reader"
	_ "quicky-go/pkg/config_reader"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// TODO: implement test cases
	fmt.Println("TestLoadConfig")
	fmt.Println(config_reader.ConfigFilePath)
	fmt.Println(config_reader.CurrentAppPath)
	assert.NotNil(t, config_reader.Info, "Config should not be nil")
	assert.NotNil(t, config_reader.Info.DB, "info.DB should not be nil")
	assert.NotNil(t, config_reader.Info.DB.DBName, "DB Name should not be nil")
	assert.NotNil(t, config_reader.Info.DB.DBHost, "DB Host should not be nil")
	// check config_reader.Info.DB.DBPort>0
	assert.True(t, config_reader.Info.DB.DBPort > 0, "DB Port should be greater than 0")
	assert.NotNil(t, config_reader.Info.DB.DBUser, "DB User should not be nil")
	assert.NotNil(t, config_reader.Info.DB.DBPassword, "DB Password should not be nil")
	assert.NotNil(t, config_reader.Info.DB.DBType, "DB DBType should not be nil")
	assert.NotNil(t, config_reader.CurrentAppPath, "CurrentAppPathshould not be nil")
	assert.NotNil(t, config_reader.ConfigFilePath, "ConfigFilePath should not be nil")

}
