package dbconfig_mysql

import (
	"strings"
	"vngom/gormex/dbconfig"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlDbConfig struct {
	dbconfig.DbConfigBase
}

func (c *MySqlDbConfig) GetConectionString(dbname string) string {
	strOps := ""
	for k, v := range c.Options {
		strOps += k + "=" + v + "&"
	}
	if len(strOps) > 0 {
		strOps = strOps[:len(strOps)-1]
	}
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + dbname + "?" + strOps

}

func (c *MySqlDbConfig) GetConectionStringNoDatabase() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/"
}
func (c *MySqlDbConfig) PingDb() error {
	d := mysql.New(mysql.Config{
		DSN: c.GetConectionStringNoDatabase(),
	})
	_, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
func (c *MySqlDbConfig) CreateDbIfNotExist(dbname string) error {
	//create mysql connection string without database name
	dns := c.GetConectionStringNoDatabase()
	//create new connection
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	ret := db.Exec("CREATE DATABASE `" + dbname + "`")
	if ret.Error != nil && !strings.Contains(ret.Error.Error(), "Error 1007") {
		return ret.Error
	}
	return nil
}
func New() *MySqlDbConfig {
	return &MySqlDbConfig{}
}
