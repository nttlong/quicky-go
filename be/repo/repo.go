package repo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"vngom/repo/repo_types"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IRepo interface {
	Insert(data interface{}) repo_types.DataActionError
	Update(data interface{}) repo_types.DataActionError
	Delete(data interface{}) repo_types.DataActionError
	AutoMigrate(data interface{}) error
}

type IRepoFactory interface {
	Get(dbName string) (*IRepo, error)
	/// SetAppDbDriverName sets the driver name for the app database.
	SetAppDbDriverName(driverName string)
	ConfigDb(driverName string, host string, port int, username string, password string)
	GetConectionStringNoDatabase() string
	PingDb()
	GetFullEntityName(entity interface{}) string
	GetColumOfEntity(entity interface{}) ([]repo_types.Column, error)
}

// ==============================

type RepoFactory struct {
	appDbDriverName string
	appDbHost       string
	appDbPort       int
	appDbUsername   string
	appDbPassword   string
}

func (rf *RepoFactory) Get(dbName string) (*IRepo, error) {
	// TODO: implement
	return nil, nil
}

func (rf *RepoFactory) SetAppDbDriverName(driverName string) {
	rf.appDbDriverName = driverName

}
func (rf *RepoFactory) GetConectionStringNoDatabase() string {
	// check all required fields are set
	if rf.appDbDriverName == "" || rf.appDbHost == "" || rf.appDbPort == 0 || rf.appDbUsername == "" || rf.appDbPassword == "" {
		panic("app db configuration is not set")
	}
	switch rf.appDbDriverName {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/", rf.appDbUsername, rf.appDbPassword, rf.appDbHost, rf.appDbPort)
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/", rf.appDbUsername, rf.appDbPassword, rf.appDbHost, rf.appDbPort)
	// TODO: implement
	default:
		panic(fmt.Sprintf("unsupported driver: %s", rf.appDbDriverName))
	}
}
func (rf *RepoFactory) ConfigDb(driverName string, host string, port int, username string, password string) {
	rf.appDbDriverName = driverName
	rf.appDbHost = host
	rf.appDbPort = port
	rf.appDbUsername = username
	rf.appDbPassword = password

}

func (rf *RepoFactory) PingDb() {
	switch rf.appDbDriverName {
	case "mysql":
		cnn := rf.GetConectionStringNoDatabase()
		_, err := gorm.Open(mysql.Open(cnn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		// TODO: implement
	case "postgres":
		panic("not implemented for postgres")
	// TODO: implement
	default:
		panic(fmt.Sprintf("unsupported driver: %s", rf.appDbDriverName))
	}

}

var CacheColumnInfo = make(map[string][]repo_types.Column)
var cacheMutex = &sync.RWMutex{}

func (rf *RepoFactory) GetFullEntityName(entity interface{}) string {
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

func (rf *RepoFactory) GetColumOfEntity(enty interface{}) ([]repo_types.Column, error) {
	// Kiểm tra enty và lấy type
	typ := reflect.TypeOf(enty)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("enty must be a struct, got %v", typ.Kind())
	}

	// Kiểm tra cache
	typeName := rf.GetFullEntityName(enty)
	cacheMutex.RLock()
	if cached, exists := CacheColumnInfo[typeName]; exists {
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	// Tính toán danh sách cột
	columns, err := ComputeColumns(typ)
	if err != nil {
		return nil, err
	}

	// Lưu vào cache
	cacheMutex.Lock()
	CacheColumnInfo[typeName] = columns
	cacheMutex.Unlock()

	return columns, nil
}

func NewRepoFactory() IRepoFactory {
	return &RepoFactory{}
}
func GoTypeToSQLType(typ reflect.Type) string {
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
func ComputeColumns(typ reflect.Type) ([]repo_types.Column, error) {
	var columns []repo_types.Column
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Bỏ qua field không xuất khẩu (private)
		if !field.IsExported() {
			continue
		}

		// Lấy tag gorm
		gormTag := field.Tag.Get("gorm")
		if gormTag == "" || strings.Contains(gormTag, "embedded") {
			//checj if field is embedded struct
			if field.Type.Kind() == reflect.Struct {
				embeddedColumns, err := ComputeColumns(field.Type)
				if err != nil {
					return nil, err
				}
				columns = append(columns, embeddedColumns...)
			}
			continue
		}

		col := repo_types.Column{
			Name:      field.Name,
			Type:      GoTypeToSQLType(field.Type),
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
				if length, ok := ExtractLength(col.Type); ok {
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
func ExtractLength(sqlType string) (int, bool) {
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
