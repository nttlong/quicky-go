package helper

import (
	"errors"
	"fmt"
	"quicky-go/pkg/db_repos/helper/helper_mysql"
	"quicky-go/pkg/db_repos/helper/info"
)

// IHelper is an interface for database helper

type IHelper interface {
	/**
	* Connect to the database use when the application starts
	* use when the application starts to check the connection
	(ensure the database is up and running)
	*
	*/
	Connect() error

	/**
	 * Create connection string for the database server only
	 * that mean without the database name
	 */
	GetConnectionString() (string, error)
	/**
	* Get the database connection string with the given database name
	* Why we need this function?
	* The project is serve for multiple databases tenant, so we need to create a connection string for each database.
	* when application start the project do not know which database it will use
	* Database will be create when new tenant register
	@param dbName the name of the database
	*/
	GetDbConnectionString(dbName string) (string, error)
	//onetime call in the whole lifecycle of the application
	CreateDatabase(dbName string) error
	/**
	*  This function will return the list of all columns by given entity
	/*
	@param enty the entity (is entity or pointer to entity) make sure GORM can identify entity
	*/
	GetColumns(enty interface{}) ([]info.Column, error)
	/**
	* Get the full name of the entity (inclue package name)
	* @param enty the entity (is entity or pointer to entity) make sure GORM can identify entity
	 */
	GetTypeNameOfEntity(enty interface{}) string
}

// cache for the helper instance
var helperCache map[string]IHelper = make(map[string]IHelper)

/*
* CreateHelper creates a new helper instance and caches it for later use. \n
* </br>
* Heed: This function should be called only once in the whole lifecycle of the application.
 */
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

/*
*
* GetHelper returns the cached helper instance for the given driver name.
* trả về instance của helper được lưu trữ trong cache cho driverName đã cho.
*返回给定 driverName 的缓存辅助实例
*  before calling this function, make sure the helper is created using CreateHelper() function.
* trước khi gọi hàm này, hãy chắc chắn rằng helper đã được tạo bằng hàm CreateHelper().
* 在调用此函数之前，请确保已使用 CreateHelper() 函数创建了助手。
* this function will cause error if the helper is not created yet.
* nếu hàm này gặp lỗi, là do helper chưa được tạo đầy đủ.
* 如果此功能失败，那是因为助手尚未完全创建。
@param driverName the name of the database driver
@return the cached helper instance for the given driver name.
*/
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
