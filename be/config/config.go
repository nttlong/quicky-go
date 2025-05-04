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
type IConfig interface {
	GetDBConfig() DBConfig
	GetServerConfig() ServerConfig
	LoadConfig(filePath string) error
}

func (c *Config) GetDBConfig() DBConfig {
	return c.DB
}

func (c *Config) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *Config) LoadConfig(filePath string) error {
	// read the file content

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// unmarshal the content into the Config struct
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return err
	}
	return nil
}

func NewConfig() IConfig {
	return &Config{}
}
