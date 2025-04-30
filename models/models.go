package models

import (
	"fmt"
	"quicky-go/models/account"
	"quicky-go/models/personal"
	"reflect"
	"strings"
	"unicode"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

var ModelList = []interface{}{
	&account.Account{},
	&personal.PersonalInfo{},
	// Thêm các model khác của bạn vào đây
	// &your_other_package.YourOtherModel{},
}

func extractColumnNameFromTag(fieldName, gormTag string) string {
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return toSnakeCase(fieldName)
}

// toSnakeCase converts a camelCase string to snake_case
func toSnakeCase(name string) string {
	var sb strings.Builder
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
func AutoMigrateWithCollateMySQL(db *gorm.DB, collation string) {
	for _, model := range ModelList {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(err)
		}

		// Detect and alter text columns for collation (MySQL specific)
		modelValue := reflect.ValueOf(model).Elem()
		modelType := modelValue.Type()

		for i := 0; i < modelValue.NumField(); i++ {
			field := modelType.Field(i)
			fieldValue := modelValue.Field(i)

			// Check if the field is a string and has a gorm tag
			if fieldValue.Kind() == reflect.String && field.Tag.Get("gorm") != "" {
				columnName := field.Tag.Get("gorm")
				actualColumnName := extractColumnNameFromTag(field.Name, columnName)

				// Alter the column to add collation (MySQL specific syntax)
				if actualColumnName != "" {
					sql := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s TEXT COLLATE %s;", toSnakeCase(modelType.Name()), actualColumnName, collation)
					err := db.Exec(sql).Error
					if err != nil {
						fmt.Printf("Error altering column %s on table %s (MySQL): %v\n", actualColumnName, toSnakeCase(modelType.Name()), err)
						// Decide if you want to panic here or just log the error
						// panic(err)
					} else {
						fmt.Printf("Successfully added COLLATE %s to column %s on table %s (MySQL)\n", collation, actualColumnName, toSnakeCase(modelType.Name()))
					}
				}
			}
			// Handle *string as well
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.String && field.Tag.Get("gorm") != "" {
				columnName := field.Tag.Get("gorm")
				actualColumnName := extractColumnNameFromTag(field.Name, columnName)

				if actualColumnName != "" {
					sql := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s TEXT COLLATE %s;", toSnakeCase(modelType.Name()), actualColumnName, collation)
					err := db.Exec(sql).Error
					if err != nil {
						fmt.Printf("Error altering column %s on table %s (MySQL): %v\n", actualColumnName, toSnakeCase(modelType.Name()), err)
						// Decide if you want to panic here or just log the error
						// panic(err)
					} else {
						fmt.Printf("Successfully added COLLATE %s to column %s on table %s (MySQL)\n", collation, actualColumnName, toSnakeCase(modelType.Name()))
					}
				}
			}
		}
	}
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
