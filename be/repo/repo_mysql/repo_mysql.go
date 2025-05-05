package repo_mysql

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"vngom/repo/repo_types"

	"gorm.io/gorm"
)

type RepoMysql struct {
	Db     *gorm.DB
	DbName string
}

func (r *RepoMysql) GetDbName() string {
	return r.DbName
}
func (r *RepoMysql) GetError(err error, typ reflect.Type, tableName string, action string) *repo_types.DataActionError {
	//"Duplicate entry '00000000-0000-0000-0000-000000000000' for key 'tenants.PRIMARY'"
	errStr := err.Error()
	if strings.Contains(errStr, "Duplicate entry") {
		ret := &repo_types.DataActionError{
			Err:          err,
			Code:         repo_types.Duplicate,
			RefTableName: tableName,
			Action:       action,
		}

		if strings.Contains(errStr, ".PRIMARY") {
			cols, ex := repo_types.ComputeColumns(typ)
			if ex != nil {
				return &repo_types.DataActionError{
					Err:  err,
					Code: repo_types.Unknown,
				}
			}
			ret.RefColumns = make([]string, 0)
			//select the primary key column in cols
			for _, col := range cols {
				if col.IsUnique {
					ret.RefColumns = append(ret.RefColumns, col.Name)
				}
			}

		} else {
			// extract the index_name column name from the error message
			// message "Error 1062 (23000): Duplicate entry 'test' for key 'tenants.idx_name'"
			indexName := strings.Split(errStr, "'")[3]
			indexName = strings.Split(indexName, ".")[1]
			cols, ex := repo_types.ComputeColumns(typ)
			if ex != nil {
				return &repo_types.DataActionError{
					Err:  err,
					Code: repo_types.Unknown,
				}
			}
			ret.RefColumns = make([]string, 0)
			//select the primary key column in cols
			for _, col := range cols {
				if col.IndexName == indexName {
					ret.RefColumns = append(ret.RefColumns, col.Name)
				}
			}

		}

		return ret
	}
	return &repo_types.DataActionError{
		Err:  err,
		Code: repo_types.Unknown,
	}
}
func (r *RepoMysql) Insert(data interface{}) *repo_types.DataActionError {

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
func (r *RepoMysql) Update(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	result := r.Db.Model(&data).Where(cond, args).Updates(data)
	if result.Error != nil {
		return r.GetError(result.Error, reflect.TypeOf(data), repo_types.GetTableNameOfEntity(data), "update")
	}

	return nil
}
func (r *RepoMysql) Select(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	panic("implement")
}
func (r *RepoMysql) Get(data interface{}, cond string, args ...interface{}) *repo_types.DataActionError {
	result := r.Db.Model(&data).Where(cond, args).First(data)
	if result.Error != nil {
		return r.GetError(result.Error, reflect.TypeOf(data), repo_types.GetTableNameOfEntity(data), "get")
	}

	return nil
}

func (r *RepoMysql) Delete(data interface{}) *repo_types.DataActionError {
	panic("implement")
}

var (
	autoMigrateCache     = make(map[string]bool)
	autoMigrateCacheLock = new(sync.RWMutex)
	autoMigrateWaitGroup sync.WaitGroup
)

// AutoMigrate performs database migration for the given data structure.
// It checks cache to avoid redundant migrations and handles nested pointers with gorm tags.
func (r *RepoMysql) AutoMigrate(data interface{}) error {
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
