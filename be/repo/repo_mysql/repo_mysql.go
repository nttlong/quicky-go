package repo_mysql

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"vngom/repo/repo_types"

	"gorm.io/gorm"
)

type RepoMysql struct {
	Db *gorm.DB
}

func (r *RepoMysql) Insert(data interface{}) *repo_types.DataActionError {
	err := r.AutoMigrate(data)
	if err != nil {
		return &repo_types.DataActionError{Err: err}
	}
	err = r.Db.Create(data).Error

	if err != nil {
		return &repo_types.DataActionError{Err: err}
	}
	return nil

}
func (r *RepoMysql) Update(data interface{}) *repo_types.DataActionError {
	panic("implement")
}
func (r *RepoMysql) Delete(data interface{}) *repo_types.DataActionError {
	panic("implement")
}

var (
	autoMigrateCache     = make(map[reflect.Type]bool)
	autoMigrateCacheLock = new(sync.RWMutex)
	autoMigrateWaitGroup sync.WaitGroup
)

// AutoMigrate performs database migration for the given data structure.
// It checks cache to avoid redundant migrations and handles nested pointers with gorm tags.
func (r *RepoMysql) AutoMigrate(data interface{}) error {
	typ := reflect.TypeOf(data)

	// Check cache with read lock to avoid redundant migration
	autoMigrateCacheLock.RLock()
	if autoMigrateCache[typ] {
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
		fname := field.Name
		fmt.Println(fname)
		fk := field.Type.Kind().String()
		fmt.Println(fk)

		// Handle slice of interface (e.g., []interface{})
		if field.Type.Kind() == reflect.Slice {
			r.Db.AutoMigrate(reflect.New(field.Type.Elem()).Interface())
		}
		KN := field.Type.Kind().String()
		fmt.Println(KN)
		//check if the field is a array
		// if field.Type.Kind() == reflect.Array {
		// 	r.Db.AutoMigrate(reflect.New(field.Type.Elem()).Interface())
		// }

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
	if autoMigrateCache[typ] {
		return nil
	}

	// Perform the actual auto-migration
	err := r.Db.AutoMigrate(data)
	if err != nil {
		return err
	}

	// Mark this type as migrated
	autoMigrateCache[typ] = true
	fmt.Println("Migrated:", typ.Name())

	return nil
}
