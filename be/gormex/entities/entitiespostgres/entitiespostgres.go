package entitiespostgres

import (
	"reflect"
	"vngom/gormex/dbconfig"
	_ "vngom/gormex/dbconfig"
	"vngom/gormex/dbconfig/dbconfig_postgres"
	_ "vngom/gormex/dbconfig/dbconfig_postgres"
	"vngom/gormex/entities"
	_ "vngom/gormex/entities"
)

type EntityPostgres[T any] struct {
	storage dbconfig.IStorage
}

func (e *EntityPostgres[T]) First(cond T) (*T, error) {

	err := e.storage.First(&cond)
	if err != nil {
		return nil, err
	}
	return &cond, nil
}
func (e *EntityPostgres[T]) Find(conds ...[]interface{}) ([]T, error) {
	var entities []T
	err := e.storage.Find(&entities, conds)
	return entities, err
}
func (e *EntityPostgres[T]) Create(entity T) (*T, error) {
	err := e.storage.Create(&entity)
	if err != nil {
		return nil, err
	} else {
		return &entity, nil
	}
}
func (e *EntityPostgres[T]) Update(entity T) error {
	return e.storage.Update(&entity)
}
func (e *EntityPostgres[T]) Delete(entity T) error {
	return e.storage.Delete(&entity)
}
func (e *EntityPostgres[T]) Count(conds ...[]interface{}) (int64, error) {
	var zero T
	t := reflect.TypeOf(zero)
	return e.storage.Count(t, conds)

}
func (e *EntityPostgres[T]) Save(entity T) error {
	return e.storage.Save(&entity)
}

func New[T any](Storage dbconfig.IStorage) entities.IEntity[T] {
	// get type of Storage
	typ := reflect.TypeOf(Storage)

	if typ != reflect.TypeOf(new(dbconfig_postgres.PostgresStorage)) {
		panic("Storage must be a pointer to dbconfig_postgres.PostgresStorage")
	}

	return &EntityPostgres[T]{
		storage: Storage,
	}
}
