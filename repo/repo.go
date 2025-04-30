package repo

import (
	"database/sql"
	"fmt"
	"log"

	"quicky-go/configs"
	"quicky-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// global variable for gorm db connection used in all package
var Repo *gorm.DB

// hash dict strore db connection and tenant db connection
var RepoHash map[string]*gorm.DB = make(map[string]*gorm.DB)

// hash map connection string for tenant db
var CnnStrHash map[string]*string = make(map[string]*string)

func GetConnectionString(tenanDb string) (*string, error) {
	//if tenant db connection string exist in hash map return connection string
	if CnnStrHash[tenanDb] != nil {
		return CnnStrHash[tenanDb], nil
	}
	//create new connection string for tenant db
	dsn, err := createDbIfNotExistDB(tenanDb)
	if err != nil {
		return nil, err
	}
	CnnStrHash[tenanDb] = dsn
	return dsn, nil

}

// this function return connection string for tenant db without database name if tenant db not exist in database server it will create new database and return connection string for new database
func createDbIfNotExistDB(tenanDb string) (*string, error) {
	cfg := configs.Info.DB
	var dsnNoDB string
	var dsnWithDB string
	var err error
	var db *sql.DB

	switch cfg.DBType {
	case configs.DBTypeMySQL:
		dsnNoDB = fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
		dsnWithDB = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, tenanDb)
	case configs.DBTypePostgres:
		dsnNoDB = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)
		dsnWithDB = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, tenanDb, cfg.DBSchema)
	case configs.DBTypeSQLServer:
		dsnNoDB = fmt.Sprintf("sqlserver://%s:%s@%s:%d", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
		dsnWithDB = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, tenanDb)
		if cfg.DBSchema != "" {
			dsnWithDB += fmt.Sprintf("&schema=%s", cfg.DBSchema)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	// Connect to the database server without specifying the database initially
	db, err = sql.Open(string(cfg.DBType), dsnNoDB)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database server: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database server: %w", err)
	}
	log.Println("Successfully connected to the database server.")

	// Check if the database exists and create it if it doesn't
	switch cfg.DBType {
	case configs.DBTypeMySQL:
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", tenanDb))
		if err != nil {
			return nil, fmt.Errorf("failed to create MySQL database '%s': %w", cfg.DBName, err)
		}
		log.Printf("MySQL database '%s' checked/created successfully.", cfg.DBName)
	case configs.DBTypePostgres:
		var exists bool
		err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", tenanDb)).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check if PostgreSQL database '%s' exists: %w", cfg.DBName, err)
		}
		if !exists {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", tenanDb))
			if err != nil {
				return nil, fmt.Errorf("failed to create PostgreSQL database '%s': %w", tenanDb, err)
			}
			log.Printf("PostgreSQL database '%s' created successfully.", tenanDb)
		} else {
			log.Printf("PostgreSQL database '%s' already exists.", tenanDb)
		}
	case configs.DBTypeSQLServer:
		var exists bool
		err = db.QueryRow(fmt.Sprintf("SELECT CASE WHEN db_id('%s') IS NOT NULL THEN 1 ELSE 0 END", tenanDb)).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check if SQL Server database '%s' exists: %w", tenanDb, err)
		}
		if !exists {
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", tenanDb))
			if err != nil {
				return nil, fmt.Errorf("failed to create SQL Server database '%s': %w", tenanDb, err)
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

	return &dsnWithDB, nil
}

// / Create new repo connection to tenant db if not exist else return existing connection
func GetRepo(tenanDb string) (*gorm.DB, error) {
	//if tenant db connection exist in hash map return connection
	if RepoHash[tenanDb] != nil {
		return RepoHash[tenanDb], nil
	}
	//create new connection string for tenant db
	dsn, err := createDbIfNotExistDB(tenanDb)
	if err != nil {
		return nil, err
	}
	//create new repo connection to tenant db
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database '%s': %w", tenanDb, err)
	}
	RepoHash[tenanDb] = db
	if tenanDb != configs.Info.DB.DBName {
		err := models.AutoMigrate(db)
		if err != nil {
			return nil, fmt.Errorf("failed to auto migrate database '%s': %w", tenanDb, err)
		}
	} else {
		err := models.AutoMigrateSystemDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to auto migrate database '%s': %w", tenanDb, err)
		}
	}

	return db, nil

}
func GetManagerRepo() *gorm.DB {

	ret, err := GetRepo(configs.Info.DB.DBName)
	if err != nil {
		panic(err)
	}
	return ret

}

func init() {

	cfg := configs.Info.DB
	_, err := createDbIfNotExistDB(cfg.DBName)
	if err != nil {
		panic(err)
	}

}
