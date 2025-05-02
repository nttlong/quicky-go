package helper

import (
	"errors"
	"fmt"
	"quicky-go/pkg/db_repos/helper/helper_mysql"
	"quicky-go/pkg/db_repos/helper/info"
)

// IHelper is an interface for database helper

type IHelper interface {
	Connect() error

	GetConnectionString() (string, error)
	CreateDatabase(dbName string) error
	GetColumns(enty interface{}) ([]info.Column, error)
	GetTypeNameOfEntity(enty interface{}) string
}

// cache for the helper instance
var helperCache map[string]IHelper = make(map[string]IHelper)

func CreateHelper(driverName string, host string, port string, user string, password string) {
	switch driverName {
	case "mysql":
		helperCache[driverName] = &helper_mysql.HelperMysql{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
		}
	case "postgres":
		panic(errors.New("not implemented yet"))

	case "mssql":
		panic(errors.New("not implemented yet"))
	default:
		panic(fmt.Sprintf("Invalid driver name: %s", driverName))

	}
}
func getCurrentPackageName() string {

	return "quicky-go/pkg/db_repos/helper"
}
func GetHelper(driverName string) IHelper {
	//check if the helper is already created
	if helperCache[driverName] == nil {
		// get current packege name
		currentPackageName := getCurrentPackageName()
		erroMsg := fmt.Sprintf("Helper not created yet,please call CreateHelper() in %s package", currentPackageName)

		panic(errors.New(erroMsg))
	}
	return helperCache[driverName]
}
