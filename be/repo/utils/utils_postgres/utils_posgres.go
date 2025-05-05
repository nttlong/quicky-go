package utils_postgres

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UtilsPostgres struct {
	user     string
	password string
	host     string
	port     int
}

func (u *UtilsPostgres) GetConectionStringNoDatabase() string {

	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		u.host,
		u.port,
		u.user,
		u.password,
	)
	return dns

}
func (u *UtilsPostgres) GetConectionString(dbName string) string {
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		u.host,
		u.port,
		u.user,
		u.password,
		dbName,
	)
	return dns

}
func (u *UtilsPostgres) ConfigDb(host string, port int, username string, password string) {
	u.host = host
	u.port = port
	u.user = username
	u.password = password

}

var (
	lockCreateDatabase  = &sync.RWMutex{}
	CacheCreateDatabase = make(map[string]bool)
)

func (u *UtilsPostgres) CreateDatabaseIfNotExists(dbName string) error {
	lockCreateDatabase.RLock()
	if _, exists := CacheCreateDatabase[dbName]; exists {
		lockCreateDatabase.RUnlock()
		return nil
	}
	lockCreateDatabase.RUnlock() // Move RUnlock here
	dns := u.GetConectionStringNoDatabase()

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := db.DB() // Get the underlying sql.DB
		if sqlDB != nil {
			sqlDB.Close() // Close the database connection
		}
	}()
	// execute the query to create the database
	postgresSQLCreateDatabaseIfNotExit := fmt.Sprintf("CREATE DATABASE \"%s\"", dbName)
	err = db.Exec(postgresSQLCreateDatabaseIfNotExit).Error
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}

	}
	fmt.Println("Database created successfully")
	// add to cache
	lockCreateDatabase.Lock()
	CacheCreateDatabase[dbName] = true
	lockCreateDatabase.Unlock()
	return nil
}
func (u *UtilsPostgres) GetDbType() string {
	return "postgres"
}
func (u *UtilsPostgres) PingDb() error {
	cnn := u.GetConectionStringNoDatabase()
	_, err := gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	} else {
		return nil
	}

}
func (u *UtilsPostgres) GetUser() string {
	return u.user
}
func (u *UtilsPostgres) GetPassword() string {
	return u.password
}
func (u *UtilsPostgres) GetHost() string {
	return u.host
}
func (u *UtilsPostgres) GetPort() int {
	return u.port
}
