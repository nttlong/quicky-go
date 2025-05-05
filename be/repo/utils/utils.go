package utils

import (
	"fmt"
	"sync"
	"vngom/repo/utils/utils_mysql"
	"vngom/repo/utils/utils_postgres"
)

type IUtils interface {
	GetDbType() string
	GetConectionStringNoDatabase() string
	GetConectionString(dbName string) string
	ConfigDb(host string, port int, username string, password string)
	CreateDatabaseIfNotExists(dbName string) error
	PingDb() error
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
}

var (
	cacheUtils map[string]IUtils
	lockUtils  sync.RWMutex
)

func NewUtils(dbType string) IUtils {
	// Bước 1: Kiểm tra nhanh với RLock
	lockUtils.RLock()
	if utils, ok := cacheUtils[dbType]; ok {
		lockUtils.RUnlock()
		return utils
	}
	lockUtils.RUnlock()

	// Bước 2: Tạo mới nếu cần, nhưng cần Lock để ghi
	lockUtils.Lock()
	defer lockUtils.Unlock()

	// Khởi tạo map nếu chưa có
	if cacheUtils == nil {
		cacheUtils = make(map[string]IUtils)
	}

	// Kiểm tra lại cache sau khi Lock (double-check)
	if utils, ok := cacheUtils[dbType]; ok {
		return utils
	}

	// Tạo instance theo dbType
	var u IUtils
	switch dbType {
	case "mysql":
		u = &utils_mysql.UtilsMysql{}
	case "postgres":
		u = &utils_postgres.UtilsPostgres{}
	default:
		panic(fmt.Sprint("Unsupported dbType: ", dbType))
	}

	// Lưu vào cache
	cacheUtils[dbType] = u
	return u
}
