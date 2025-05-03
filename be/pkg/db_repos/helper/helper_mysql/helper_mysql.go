package helper_mysql

import (
	"errors"
	"fmt"
	"quicky-go/pkg/db_repos/helper/info"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type HelperMysql struct {
	connectionString string

	Host     string
	Port     string
	User     string
	Password string
}

// cache db connection string
var DbConnectionString map[string]string = make(map[string]string)

// lock db
var lockDbConnectionString = sync.RWMutex{}

// cache thao tác tạo database| cache create database

var CacheCreateDatabase = make(map[string]bool)

// lock cache thao tác tạo database | lock create database
var lockCreateDatabase = sync.RWMutex{}

var CacheColumnInfo = make(map[string][]info.Column)

// cacheMutex bảo vệ CacheColumnInfo khỏi truy cập đồng thời
var cacheMutex sync.RWMutex

// cache gorm.DB
var CacheGormDB = make(map[string]*gorm.DB)

// lock cache gorm.DB
var lockCacheGormDB = sync.RWMutex{}

func (m *HelperMysql) GetTypeNameOfEntity(entity interface{}) string {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("entity must be a struct, got %v", typ.Kind()))
	}

	// Kiểm tra cache
	fullName := typ.String()
	return fullName
}
func (m *HelperMysql) GetConnectionString() (string, error) {
	if m.connectionString == "" {
		// create connection with username password and no database
		if m.User == "" || m.Password == "" || m.Host == "" || m.Port == "" {
			return "", errors.New("User, Password, Host, Port must be set")

		}
		m.connectionString = m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/"
	}
	return m.connectionString, nil
}
func (m *HelperMysql) GetDbConnectionString(dbName string) (string, error) {
	lockDbConnectionString.RLock()
	if cached, exists := DbConnectionString[m.connectionString]; exists {
		lockDbConnectionString.RUnlock()
		return cached, nil
	}
	lockDbConnectionString.RUnlock()
	if m.User == "" || m.Password == "" || m.Host == "" || m.Port == "" {
		return "", errors.New("User, Password, Host, Port must be set")

	}
	//try create database if not exists
	err := m.CreateDatabase(dbName)
	if err != nil {
		return "", err
	}

	//create  mysql db connection to db
	// use m.Host, m.Port, m.User, m.Password and dbName to create connection string
	//"root:123456@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"

	cnn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, dbName)
	// try connect to db
	db, err := gorm.Open(mysql.Open(cnn), &gorm.Config{})
	if err != nil {
		return "", err
	}
	defer func() {
		sqlDB, _ := db.DB() // Get the underlying sql.DB
		if sqlDB != nil {
			sqlDB.Close() // Close the database connection
		}
	}()
	// check if connection is successful or not
	err = db.Exec(fmt.Sprintf("USE %s", dbName)).Error
	if err != nil {
		return "", err
	}
	fmt.Println("Connected to database", dbName)
	// add to cache
	lockDbConnectionString.Lock()
	DbConnectionString[m.connectionString] = cnn
	lockDbConnectionString.Unlock()
	return cnn, nil

}
func (m *HelperMysql) Connect() error {
	//create  mysql db connection
	dsn, err := m.GetConnectionString()
	if err != nil {
		return err
	}

	_, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	fmt.Println("Connected to database")
	// check if connection is successful or not
	return nil
}
func (m *HelperMysql) CreateDatabase(dbName string) error {
	lockCreateDatabase.RLock()
	if _, exists := CacheCreateDatabase[dbName]; exists {
		lockCreateDatabase.RUnlock()
		return nil
	}
	lockCreateDatabase.RUnlock() // Move RUnlock here

	dsn, err := m.GetConnectionString()
	if err != nil {
		return err
	}
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
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)).Error
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
func (m *HelperMysql) GetColumns(enty interface{}) ([]info.Column, error) {
	// Kiểm tra enty và lấy type
	typ := reflect.TypeOf(enty)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("enty must be a struct, got %v", typ.Kind())
	}

	// Kiểm tra cache
	typeName := typ.String()
	cacheMutex.RLock()
	if cached, exists := CacheColumnInfo[typeName]; exists {
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	// Tính toán danh sách cột
	columns, err := computeColumns(typ)
	if err != nil {
		return nil, err
	}

	// Lưu vào cache
	cacheMutex.Lock()
	CacheColumnInfo[typeName] = columns
	cacheMutex.Unlock()

	return columns, nil
}

// computeColumns tính toán danh sách cột từ type
func computeColumns(typ reflect.Type) ([]info.Column, error) {
	var columns []info.Column
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Bỏ qua field không xuất khẩu (private)
		if !field.IsExported() {
			continue
		}

		// Lấy tag gorm
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" || strings.Contains(gormTag, "embedded") {
			continue
		}

		col := info.Column{
			Name:      field.Name,
			Type:      goTypeToSQLType(field.Type),
			AllowNull: true,
		}

		// Phân tích tag gorm
		tags := strings.Split(gormTag, ";")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if strings.HasPrefix(tag, "column:") {
				col.Name = strings.TrimPrefix(tag, "column:")
			}
			if strings.HasPrefix(tag, "type:") {
				col.Type = strings.TrimPrefix(tag, "type:")
				if length, ok := extractLength(col.Type); ok {
					col.Length = &length
				}
			}
			if tag == "not null" {
				col.AllowNull = false
			}
			if tag == "unique" {
				col.IsUnique = true
			}
			if strings.HasPrefix(tag, "index:") {
				col.IndexName = strings.TrimPrefix(tag, "index:")
			} else if tag == "index" {
				col.IndexName = "idx_" + strings.ToLower(col.Name)
			}
			if tag == "primary_key" {
				col.IsUnique = true
			}
		}

		columns = append(columns, col)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no valid columns found in %v", typ.Name())
	}

	return columns, nil
}

// goTypeToSQLType suy ra kiểu SQL từ kiểu Go
func goTypeToSQLType(typ reflect.Type) string {
	switch typ.Kind() {
	case reflect.String:
		return "text"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "number"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "number"
	case reflect.Struct:
		if typ.String() == "time.Time" {
			return "datetime"
		}
	}
	return "text" // Mặc định
}

// extractLength trích xuất độ dài từ type (ví dụ: varchar(100) -> 100)
func extractLength(sqlType string) (int, bool) {
	start := strings.Index(sqlType, "(")
	end := strings.Index(sqlType, ")")
	if start != -1 && end != -1 && start < end {
		lengthStr := sqlType[start+1 : end]
		length, err := strconv.Atoi(lengthStr)
		if err == nil {
			return length, true
		}
	}
	return 0, false
}
func (m *HelperMysql) GetDb(dbName string) (*gorm.DB, error) {
	//check if db connection is already cached
	lockCacheGormDB.RLock()
	if cached, exists := CacheGormDB[dbName]; exists {
		lockCacheGormDB.RUnlock()
		return cached, nil
	}
	lockCacheGormDB.RUnlock()
	//create new db connection
	dsn, err := m.GetDbConnectionString(dbName)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//add to cache
	lockCacheGormDB.Lock()
	CacheGormDB[dbName] = db
	lockCacheGormDB.Unlock()
	return db, nil

}

type RepositoryMySql struct {
	db     *gorm.DB
	DbName string
}

func (r *RepositoryMySql) SetDb(db *gorm.DB, dbName string) error {
	r.db = db
	r.DbName = dbName
	return nil
}
func (r *RepositoryMySql) Insert(entity interface{}) *info.DataActionError {
	// use gorm to insert entity
	//r.db.Model(entity).Save()
	ret := r.db.Create(entity)
	// check if error

	if ret.Error != nil {
		return &info.DataActionError{
			Err:          ret.Error,
			Action:       "insert",
			Code:         info.Duplicate,
			RefColumns:   nil,
			RefTableName: "",
		}
	}
	return nil

}
func (r *RepositoryMySql) Update(entity interface{}) *info.DataActionError {
	panic("unimplemented")
}
func (r *RepositoryMySql) Delete(entity interface{}) *info.DataActionError {
	panic("unimplemented")
}
func (r *RepositoryMySql) AutoMigrate(entity interface{}) error {
	err := r.db.AutoMigrate(entity)
	if err != nil {
		return err
	}
	return nil

}

// use gorm to find all entity
