// / declare  structs for configuration load ymal file
package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import driver MySQL
	_ "github.com/lib/pq"              // Import driver PostgreSQL (nếu bạn dùng)

	//_ "github.com/microsoft/go-sqlcmd/mssql" // Import driver SQL Server (nếu bạn dùng)
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
func CheckDBConnect() error {
	cfg := Info.DB
	var dsnNoDB string
	var dsnWithDB string
	var err error
	var db *sql.DB

	switch cfg.DBType {
	case DBTypeMySQL:
		dsnNoDB = fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
		dsnWithDB = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	case DBTypePostgres:
		dsnNoDB = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)
		dsnWithDB = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSchema)
	case DBTypeSQLServer:
		dsnNoDB = fmt.Sprintf("sqlserver://%s:%s@%s:%d", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
		dsnWithDB = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		if cfg.DBSchema != "" {
			dsnWithDB += fmt.Sprintf("&schema=%s", cfg.DBSchema)
		}
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	// Connect to the database server without specifying the database initially
	db, err = sql.Open(string(cfg.DBType), dsnNoDB)
	if err != nil {
		return fmt.Errorf("failed to open connection to database server: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database server: %w", err)
	}
	log.Println("Successfully connected to the database server.")

	// Check if the database exists and create it if it doesn't
	switch cfg.DBType {
	case DBTypeMySQL:
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName))
		if err != nil {
			return fmt.Errorf("failed to create MySQL database '%s': %w", cfg.DBName, err)
		}
		log.Printf("MySQL database '%s' checked/created successfully.", cfg.DBName)
	case DBTypePostgres:
		var exists bool
		err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", cfg.DBName)).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if PostgreSQL database '%s' exists: %w", cfg.DBName, err)
		}
		if !exists {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				return fmt.Errorf("failed to create PostgreSQL database '%s': %w", cfg.DBName, err)
			}
			log.Printf("PostgreSQL database '%s' created successfully.", cfg.DBName)
		} else {
			log.Printf("PostgreSQL database '%s' already exists.", cfg.DBName)
		}
	case DBTypeSQLServer:
		var exists bool
		err = db.QueryRow(fmt.Sprintf("SELECT CASE WHEN db_id('%s') IS NOT NULL THEN 1 ELSE 0 END", cfg.DBName)).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if SQL Server database '%s' exists: %w", cfg.DBName, err)
		}
		if !exists {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				return fmt.Errorf("failed to create SQL Server database '%s': %w", cfg.DBName, err)
			}
			log.Printf("SQL Server database '%s' created successfully.", cfg.DBName)
		} else {
			log.Printf("SQL Server database '%s' already exists.", cfg.DBName)
		}
	default:
		// Should not reach here as the type is checked earlier
	}

	// Optionally, you can now try to connect to the specific database
	// dbWithDB, err := sql.Open(string(cfg.DBType), dsnWithDB)
	// if err != nil {
	// 	return fmt.Errorf("failed to open connection to database '%s': %w", cfg.DBName, err)
	// }
	// defer dbWithDB.Close()
	// if err := dbWithDB.Ping(); err != nil {
	// 	return fmt.Errorf("failed to ping database '%s': %w", cfg.DBName, err)
	// }
	// log.Printf("Successfully connected and pinged database '%s'.", cfg.DBName)

	return nil
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
	err := CheckDBConnect()
	if err != nil {
		log.Fatal(err)
	}
}
