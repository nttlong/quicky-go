package dbconfig_postgres

import (
	"fmt"
	"strings"
	"sync"
	"vngom/gormex/dbconfig"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDbConfig struct {
	dbconfig.DbConfigBase
}
type PostgresStorage struct {
	db *gorm.DB
}

func (c *PostgresDbConfig) GetConectionString(dbname string) string {
	strOps := ""
	for k, v := range c.Options {
		if k == "collation" {
			continue
		}
		strOps += k + "=" + v + "&"
	}
	if len(strOps) > 0 {
		strOps = strOps[:len(strOps)-1]
	}
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + dbname + "?" + strOps

}

func (c *PostgresDbConfig) GetConectionStringNoDatabase() string {
	//create postgres connection string without database name
	strOps := ""
	for k, v := range c.Options {
		strOps += k + "=" + v + "&"
	}
	if len(strOps) > 0 {
		strOps = strOps[:len(strOps)-1]
	}
	//"host=%s port=%d user=%s password=%s sslmode=disable"
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/"

}
func (c *PostgresDbConfig) PingDb() error {
	d := postgres.New(postgres.Config{
		DSN: c.GetConectionStringNoDatabase(),
	})
	_, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) AutoMigrate(entity interface{}) error {
	return s.db.AutoMigrate(entity)
}

func (c *PostgresDbConfig) GetStorage(dbName string) (dbconfig.IStorage, error) {
	err := c.PingDb()
	if err != nil {
		return nil, err
	}
	if err = c.createDbIfNotExist(dbName); err != nil {
		return nil, err
	}
	dns := c.GetConectionString(dbName)
	d, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: d}, nil

}
func New() dbconfig.IDbConfig {
	return &PostgresDbConfig{}
}

var (
	cacheCreateDbIfNotExist = make(map[string]bool)
	lockCreateDbIfNotExist  = sync.RWMutex{}
)

func (c *PostgresDbConfig) CreateDbIfNotExist(dbname string) error {
	lockCreateDbIfNotExist.RLock()
	isCreated := cacheCreateDbIfNotExist[dbname]
	lockCreateDbIfNotExist.RUnlock()
	if isCreated {
		return nil
	}
	lockCreateDbIfNotExist.Lock()
	defer lockCreateDbIfNotExist.Unlock()
	if cacheCreateDbIfNotExist[dbname] {
		return nil
	}
	//create database if not exist
	err := c.createDbIfNotExist(dbname)
	if err != nil {
		return err
	}
	cacheCreateDbIfNotExist[dbname] = true
	return nil
}

// ======================================================================
func (c *PostgresDbConfig) createDbIfNotExist(dbname string) error {
	//create postgres connection string without database name
	dns := c.GetConectionStringNoDatabase()
	//create new connection
	d, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	//create database if not exist
	/**
		CREATE DATABASE mydb
	WITH ENCODING 'UTF8'
	LC_COLLATE 'vi_VN.UTF-8'
	LC_CTYPE 'vi_VN.UTF-8';
	*/
	collate := c.DbConfigBase.Options["collate"]
	sql := fmt.Sprintf("CREATE DATABASE \"%s\" WITH ENCODING 'UTF8' LC_COLLATE '%s' LC_CTYPE '%s'", dbname, collate, collate)
	err = d.Exec(sql).Error
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}
	return nil
}
