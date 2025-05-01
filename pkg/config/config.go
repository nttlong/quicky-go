package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DBType string

const (
	DBTypeMySQL     DBType = "mysql"
	DBTypePostgres  DBType = "postgres"
	DBTypeSQLServer DBType = "sqlserver"
)

// DBConfig represents the database configuration.
type DBConfig struct {
	Type     DBType `yaml:"type"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type Config struct {
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
	// Add other configurations here if needed.
}

// load the configuration from the YAML file.
func LoadConfig(filePath string) Config {
	// read the file content
	ret := Config{}
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// unmarshal the content into the Config struct
	err = yaml.Unmarshal(content, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}
