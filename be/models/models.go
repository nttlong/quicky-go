package models

import (
	"fmt"
	"quicky-go/models/account"
	"quicky-go/models/personal"
	"quicky-go/models/tenants"
	"reflect"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

var ModelList = []interface{}{
	&account.Account{},
	&personal.PersonalInfo{},
	// Thêm các model khác của bạn vào đây
	// &your_other_package.YourOtherModel{},
}
var ModelListManager = []interface{}{
	&tenants.Tenants{},
	// Thêm các model khác của bạn vào đây
	// &your_other_package.YourOtherModel{},
}

// declare error type for AutoMigrate function
type AutoMigrateError struct {
	Model   interface{}
	Err     error
	Message string
}

func getTableAndStringColumnsFromStruct(s interface{}) (tableName string, stringColumns []string) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	} else if val.Kind() != reflect.Struct {
		fmt.Println("Lỗi: Đầu vào không phải là struct hoặc con trỏ đến struct.")
		return "", nil
	}

	typ := val.Type()
	tableName = typ.Name() // Lấy tên bảng trực tiếp từ tên kiểu struct

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.String {
			mapdatabaseTag := field.Tag.Get("mapdatabase")
			if mapdatabaseTag != "" {
				stringColumns = append(stringColumns, mapdatabaseTag)
			} else {
				stringColumns = append(stringColumns, field.Name) // Nếu không có tag, dùng tên field
			}
		}
	}

	return tableName, stringColumns
}

func alterColumnIgnoreCaseSensitiveForMySQL(db *gorm.DB, tableName, columnName string) error {
	sql := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s TEXT COLLATE utf8mb4_general_ci;", tableName, columnName)
	err := db.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("failed to alter column '%s' on table '%s' to be case-insensitive (utf8mb4_general_ci): %w", columnName, tableName, err)
	}
	return nil
}

func AutoMigrate(db *gorm.DB) *AutoMigrateError {
	// get curent database name in db

	for _, model := range ModelList {
		err := db.AutoMigrate(model)
		if err != nil {
			modelName := reflect.TypeOf(model).Elem().Name()
			modelPackage := reflect.TypeOf(model).Elem().PkgPath()
			return &AutoMigrateError{
				Model:   model,
				Err:     err,
				Message: fmt.Sprintf("Error auto-migrating model %s: %v in package %s", modelName, err, modelPackage),
			}

		}
		//get all columns of the model
		tableName, stringColumns := getTableAndStringColumnsFromStruct(model)
		if len(stringColumns) > 0 {
			for _, column := range stringColumns {
				// check db if run on mysql
				if db.Config.Dialector.Name() == "mysql" {
					err := alterColumnIgnoreCaseSensitiveForMySQL(db, tableName, column)
					if err != nil {
						return &AutoMigrateError{
							Model:   model,
							Err:     err,
							Message: fmt.Sprintf("Error altering column %s on table %s to be case-insensitive (utf8mb4_general_ci): %v", column, tableName, err),
						}
					}
				} else {
					panic(fmt.Sprintf("Unsupported dialect: %s", db.Config.Dialector.Name()))

				}
			}
		}
	}
	return nil
}
func AutoMigrateSystemDB(db *gorm.DB) *AutoMigrateError {
	// get curent database name in db

	for _, model := range ModelListManager {
		err := db.AutoMigrate(model)
		if err != nil {
			modelName := reflect.TypeOf(model).Elem().Name()
			modelPackage := reflect.TypeOf(model).Elem().PkgPath()
			return &AutoMigrateError{
				Model:   model,
				Err:     err,
				Message: fmt.Sprintf("Error auto-migrating model %s: %v in package %s", modelName, err, modelPackage),
			}

		}
		//get all columns of the model
		tableName, stringColumns := getTableAndStringColumnsFromStruct(model)
		if len(stringColumns) > 0 {
			for _, column := range stringColumns {
				// check db if run on mysql
				if db.Config.Dialector.Name() == "mysql" {
					err := alterColumnIgnoreCaseSensitiveForMySQL(db, tableName, column)
					if err != nil {
						return &AutoMigrateError{
							Model:   model,
							Err:     err,
							Message: fmt.Sprintf("Error altering column %s on table %s to be case-insensitive (utf8mb4_general_ci): %v", column, tableName, err),
						}
					}
				} else {
					panic(fmt.Sprintf("Unsupported dialect: %s", db.Config.Dialector.Name()))

				}
			}
		}
	}
	return nil
}
