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

// cache thao tác tạo database| cache create database

var CacheCreateDatabase = make(map[string]bool)

// lock cache thao tác tạo database | lock create database
var lockCreateDatabase = sync.RWMutex{}

var CacheColumnInfo = make(map[string][]info.Column)

// cacheMutex bảo vệ CacheColumnInfo khỏi truy cập đồng thời
var cacheMutex sync.RWMutex

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
