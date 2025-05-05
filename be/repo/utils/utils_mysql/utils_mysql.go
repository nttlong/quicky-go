package utils_mysql

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UtilsMysql struct {
	user     string
	password string
	host     string
	port     int
}

var (
	lockCreateDatabase  = &sync.RWMutex{}
	CacheCreateDatabase = make(map[string]bool)
)

func (u *UtilsMysql) GetConectionStringNoDatabase() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/", u.user, u.password, u.host, u.port)
}
func (u *UtilsMysql) GetConectionString(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", u.user, u.password, u.host, u.port, dbName)
}
func (u *UtilsMysql) ConfigDb(host string, port int, username string, password string) {
	u.host = host
	u.port = port
	u.user = username
	u.password = password

}
func (u *UtilsMysql) GetDbType() string {
	return "mysql"
}

func (u *UtilsMysql) CreateDatabaseIfNotExists(dbName string) error {
	lockCreateDatabase.RLock()
	if _, exists := CacheCreateDatabase[dbName]; exists {
		lockCreateDatabase.RUnlock()
		return nil
	}
	lockCreateDatabase.RUnlock() // Move RUnlock here
	dsn := u.GetConectionStringNoDatabase()

	dsn += "mysql"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, _ := db.DB() // Get the underlying sql.DB
		if sqlDB != nil {
			sqlDB.Close() // Close the database connection
		}
	}()

	// create database
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)).Error
	if err != nil {
		return err
	}
	fmt.Println("Database created successfully")
	// add to cache
	lockCreateDatabase.Lock()
	CacheCreateDatabase[dbName] = true
	lockCreateDatabase.Unlock()
	return nil
}
func (u *UtilsMysql) PingDb() error {
	cnn := u.GetConectionStringNoDatabase()
	_, err := gorm.Open(mysql.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	} else {
		return nil
	}

}
func (u *UtilsMysql) GetUser() string {
	return u.user
}
func (u *UtilsMysql) GetPassword() string {
	return u.password
}
func (u *UtilsMysql) GetHost() string {
	return u.host
}
func (u *UtilsMysql) GetPort() int {
	return u.port
}
