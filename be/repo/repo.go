package repo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"vngom/repo/repo_mysql"
	"vngom/repo/repo_postgres"
	"vngom/repo/repo_types"
	"vngom/repo/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IRepo interface {
	Insert(data interface{}) *repo_types.DataActionError
	Update(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError
	Get(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError
	Select(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError
	Delete(data interface{}) *repo_types.DataActionError
	AutoMigrate(data interface{}) error
	GetError(err error, typ reflect.Type, tableName string, action string) *repo_types.DataActionError
	GetDbName() string
}

type IRepoFactory interface {
	utils.IUtils
	Get(dbName string) (IRepo, error)

	GetFullEntityName(entity interface{}) string
	GetColumOfEntity(enty interface{}) ([]repo_types.Column, error)
}

// ==============================

type RepoFactory struct {
	utils.IUtils
}

// caceh for repo
var repoCache = make(map[string]IRepo)
var repoCacheMutex = &sync.RWMutex{}

func (rf *RepoFactory) Get(dbName string) (IRepo, error) {

	// check if the repo is already created
	repoCacheMutex.RLock()
	if repo, exists := repoCache[dbName]; exists {
		repoCacheMutex.RUnlock()
		return repo, nil
	}
	repoCacheMutex.RUnlock()
	// lock the cache for writing
	repoCacheMutex.Lock()

	err := rf.CreateDatabaseIfNotExists(dbName)
	if err != nil {
		repoCacheMutex.Unlock()
		return nil, err
	}
	// create a new repo
	repo, err := NewRepo(
		dbName,
		rf.IUtils.GetDbType(),
		rf.IUtils.GetHost(),
		rf.IUtils.GetPort(),
		rf.IUtils.GetUser(),
		rf.IUtils.GetPassword(),
	)
	if err != nil {
		repoCacheMutex.Unlock()
		return nil, err
	}
	// add the repo to the cache
	repoCache[dbName] = repo
	// unlock the cache
	repoCacheMutex.Unlock()
	return repo, nil
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

var (
	RepoFactoryInstance IRepoFactory
	once                sync.Once
)

func NewRepoFactory(dbType string) IRepoFactory {
	once.Do(func() {
		RepoFactoryInstance = &RepoFactory{
			IUtils: utils.NewUtils(dbType),
		}
	})

	return RepoFactoryInstance
}
func NewRepo(
	dbName string,
	driverName string,
	host string,
	port int,
	username string,
	password string,
) (IRepo, error) {
	fmt.Print("NewRepo: " + dbName)

	switch driverName {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username,
			password,

			host,
			port,
			dbName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		// Set reasonable pool sizes (tune nếu cần)
		sqlDB.SetMaxOpenConns(200)          // Max concurrent conns
		sqlDB.SetMaxIdleConns(50)           // Idle connections giữ lại
		sqlDB.SetConnMaxLifetime(time.Hour) // Refresh sau 1 giờ
		return &repo_mysql.RepoMysql{
			Db:     db,
			DbName: dbName,
		}, nil
	case "postgres":
		dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			username,
			password,
			dbName)
		db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		return &repo_postgres.RepoPostgres{
			Db:     db,
			DbName: dbName,
		}, nil

	default:
		panic(fmt.Sprintf("unsupported driver: %s", driverName))
	}
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
