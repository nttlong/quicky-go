// / declare  structs for configuration load ymal file
package configs

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
	DBType     DBType `yaml:"dbType"`
	DBName     string `yaml:"dbName"`
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
	DBHost     string `yaml:"dbHost"`
	DBPort     int    `yaml:"dbPort"`
	DBSchema   string `yaml:"dbSchema"`
}

// Config is the main configuration struct.
type Config struct {
	DB DBConfig `yaml:"db"`
}

// declare a global variabale to store the configuration
var Info Config

// / declare a global variabale to store the current app path
var CurrentAppPath string
var ConfigFilePath string

// LoadConfig loads the configuration from the given file path. make sure once call this function
func LoadConfig(filePath string) error {
	// read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// unmarshal the content into the Config struct
	err = yaml.Unmarshal(content, &Info)
	if err != nil {
		return err
	}
	return nil
}

// get config as json format with indentation
func GetConfigAsJSON() string {
	// marshal the Config struct into json format
	jsonBytes, err := yaml.Marshal(Info)
	if err != nil {
		panic(err)
	}
	// print the json format with indentation
	return (string(jsonBytes))
}

// init function is called when the package is imported and create a new instance of Config struct make sure just once
func init() {
	// set the current app path
	CurrentAppPath, _ = os.Getwd()
	// set the config file path
	ConfigFilePath = CurrentAppPath + "/config.yaml"
	// create a new instance of Config struct
	Info = Config{}
	// load the configuration from the config file
	LoadConfig(ConfigFilePath)
}
