package repo_postgres

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"vngom/repo/repo_types"
	_ "vngom/repo/repo_types"

	"gorm.io/gorm"
)

type RepoPostgres struct {
	Db     *gorm.DB
	DbName string
}

var autoMigrateCache = make(map[string]bool)
var autoMigrateCacheLock = new(sync.RWMutex)

func (r *RepoPostgres) GetDbName() string {
	return r.DbName
}

// AutoMigrate performs database migration for the given data structure.
// It checks cache to avoid redundant migrations and handles nested pointers with gorm tags.
func (r *RepoPostgres) AutoMigrate(data interface{}) error {
	typ := reflect.TypeOf(data)
	cachekey := r.GetDbName() + "/" + typ.String()

	// Check cache with read lock to avoid redundant migration
	autoMigrateCacheLock.RLock()
	if autoMigrateCache[cachekey] {
		autoMigrateCacheLock.RUnlock()
		return nil
	}
	autoMigrateCacheLock.RUnlock()

	// Handle pointer type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Ensure data is a struct
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct, got %v", typ.Kind())
	}

	// Iterate over fields to handle nested structures and pointers
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// fname := field.Name
		// fmt.Println(fname)
		// fk := field.Type.Kind().String()
		// fmt.Println(fk)

		// Handle slice of interface (e.g., []interface{})
		if field.Type.Kind() == reflect.Slice {
			r.Db.AutoMigrate(reflect.New(field.Type.Elem()).Interface())
		}

		// Handle pointer with gorm tag
		if field.Type.Kind() == reflect.Ptr {
			elementType := field.Type.Elem()
			if elementType.Kind() == reflect.Struct {
				tag := field.Tag.Get("gorm")
				if tag != "" && strings.Contains(tag, "foreignKey:") {
					// Create a new instance for the pointer type
					instance := reflect.New(field.Type.Elem()).Interface()
					err := r.Db.AutoMigrate(instance)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Acquire write lock to perform migration and update cache
	autoMigrateCacheLock.Lock()
	defer autoMigrateCacheLock.Unlock()

	// Double-check cache to handle concurrent calls
	if autoMigrateCache[cachekey] {
		return nil
	}

	// Perform the actual auto-migration
	err := r.Db.AutoMigrate(data)
	if err != nil {
		return err
	}

	// Mark this type as migrated
	autoMigrateCache[cachekey] = true
	fmt.Println("Migrated:", cachekey)

	return nil
}
func (r *RepoPostgres) Insert(data interface{}) *repo_types.DataActionError {

	err := r.AutoMigrate(data)

	if err != nil {
		return &repo_types.DataActionError{Err: err}
	}
	err = r.Db.Create(data).Error
	if err != nil {
		tableName := repo_types.GetTableNameOfEntity(data)
		typ, errT := repo_types.GetReflectType(data)
		if errT != nil {
			return &repo_types.DataActionError{Err: err}
		}
		startAt := time.Now()

		rr := r.GetError(err, typ, tableName, "insert")
		elapseTime := time.Since(startAt)
		fmt.Println("GetError: ", elapseTime.Milliseconds())
		return rr
	}
	return nil
}
func (r *RepoPostgres) Update(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	panic("not implemented")
}
func (r *RepoPostgres) Get(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	panic("not implemented")
}
func (r *RepoPostgres) Select(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	panic("not implemented")
}
func (r *RepoPostgres) Delete(data interface{}) *repo_types.DataActionError {
	panic("not implemented")
}

func (r *RepoPostgres) GetError(err error, typ reflect.Type, tableName string, action string) *repo_types.DataActionError {

	//"duplicate key value violates unique constraint \"idx_db_tenants_name\""
	errStr := err.Error()
	if strings.Contains(errStr, "duplicate key value violates unique constraint") {

		index_name := strings.Split(errStr, "\"")[1]
		fmt.Println(index_name)
		cols, ex := repo_types.ComputeColumns(typ)
		if ex != nil {
			return &repo_types.DataActionError{
				Err:  err,
				Code: repo_types.Unknown,
			}
		}
		retErr := &repo_types.DataActionError{
			Err:    err,
			Action: action,
			Code:   repo_types.Duplicate,

			RefTableName: tableName,
		}
		retErr.RefColumns = make([]string, 0)
		for _, col := range cols {

			if col.IndexName == index_name {
				retErr.RefColumns = append(retErr.RefColumns, col.Name)
			}
		}
		return retErr
	}
	return &repo_types.DataActionError{
		Err:  err,
		Code: repo_types.Unknown,
	}

}
