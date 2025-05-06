package gormex

import (
	"fmt"
	"vngom/gormex/dbconfig"
	"vngom/gormex/dbconfig/dbconfig_mysql"
	"vngom/gormex/dbconfig/dbconfig_postgres"

	_ "gorm.io/gorm"
)

type IGormEx interface {
	GetDbConfig() dbconfig.IDbConfig
}

func NewDbConfig(dbType string) dbconfig.IDbConfig {
	if dbType == "mysql" {
		return &dbconfig_mysql.MySqlDbConfig{}
	}
	if dbType == "postgres" {
		//TODO: implement postgres config
		return &dbconfig_postgres.PostgresDbConfig{}
	}
	panic(fmt.Sprintf("Unsupported database type %s", dbType))
}
