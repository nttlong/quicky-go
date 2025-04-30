package errors

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type DbErrorType int

const (
	NoError         DbErrorType = 0
	DuplicateError  DbErrorType = 1
	RequiredColumns DbErrorType = 2
	OtherError      DbErrorType = 99
)

type ErrorAnalysisResult struct {
	TableName    string
	Columns      []string
	DbErrorType  DbErrorType
	ErrorMessage string // Thêm thông báo lỗi gốc để debug
}

// Implement interface error cho ErrorAnalysisResult
func (e ErrorAnalysisResult) Error() string {
	return fmt.Sprintf("DB Error Type: %v, Table: %s, Columns: %v, Original Error: %s",
		e.DbErrorType, e.TableName, e.Columns, e.ErrorMessage)
}

// Cache entity model and name
var entityModelCache = make(map[interface{}]string)

// function get tablename form entity model
func GetTableName(db *gorm.DB, entity interface{}) (*string, error) {
	//check entity model in cache
	// get type of entity model
	t := reflect.TypeOf(entity)
	// check cache

	if tableName, ok := entityModelCache[t]; ok {
		return &tableName, nil
	}
	stmt := &gorm.Statement{DB: db}
	stmt.Model = entity // Thiết lập model cho statement bằng cách gán field
	err := stmt.Parse(stmt.Model)
	if err == nil && stmt.Schema != nil {
		// update cache
		entityModelCache[t] = stmt.Schema.Table
		return &stmt.Schema.Table, nil
	}
	return nil, fmt.Errorf("Cannot get table name from entity model: %v", err)
	// update cache

}

// decalglobal var has dict store entyty model and map[string][]string
var PrimaryColumnsCache = make(map[string][]string)

// getPrimaryColumns nhận vào một entity (struct) và trả về một slice chứa tên
// các cột (trong database) được đánh dấu là primary key thông qua tag `gorm`.
// getPrimaryColumns takes an entity (struct) and returns a slice containing the
// database column names marked as primary key using the `gorm` tag.
func GetPrimaryColumns(db *gorm.DB, entity interface{}) []string {

	//table is
	tableName, terr := GetTableName(db, entity)
	if terr != nil {
		panic(terr)
	}
	if tableName == nil {
		panic("Cannot get table name from entity model")
	}

	if _, ok := PrimaryColumnsCache[*tableName]; ok {

		return PrimaryColumnsCache[*tableName]
	}
	primaryColumns := make([]string, 0)
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	// If it's a pointer, get the actual value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Ensure the input is a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Error: Input is not a struct or a pointer to a struct.")
		return primaryColumns
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		gormTag := field.Tag.Get("gorm")
		tags := strings.Split(gormTag, ";")

		isPrimaryKey := false
		columnName := field.Name // Default column name is the field name

		// Parse the gorm tags
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag == "primaryKey" || tag == "primary_key" { // Both spellings are common
				isPrimaryKey = true
			} else if strings.HasPrefix(tag, "column:") {
				columnName = strings.TrimPrefix(tag, "column:")
			}
		}

		if isPrimaryKey {
			primaryColumns = append(primaryColumns, columnName)
		}
	}
	// Cập nhật PrimaryColumnsCache
	PrimaryColumnsCache[*tableName] = primaryColumns
	return primaryColumns
}

// decalre global var has dict store entyty model lool likes entityModelDictRequireCols[tableName] = []string{"column1", "column2"}
var RequiredColumnsCache = make(map[string][]string)

// getRequiredColumns nhận vào một entity (struct) và trả về một slice chứa tên
// các cột (trong database) được đánh dấu là bắt buộc (not null) thông qua tag `gorm`.
func GetRequiredColumns(db *gorm.DB, entity interface{}) []string {
	//get table name from entity model
	tableName, terr := GetTableName(db, entity)
	if terr != nil {
		panic(terr)
	}
	if tableName == nil {
		panic("Cannot get table name from entity model")
	}
	if _, ok := RequiredColumnsCache[*tableName]; ok {

		return RequiredColumnsCache[*tableName]
	}
	requiredColumns := make([]string, 0)
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	// Nếu là pointer, lấy giá trị thực tế
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Đảm bảo đầu vào là struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Error: Input is not a struct or a pointer to a struct.")
		return requiredColumns
	}

	// Lặp qua các field của struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		gormTag := field.Tag.Get("gorm")
		tags := strings.Split(gormTag, ";")

		isNotNull := false
		columnName := field.Name // Tên cột mặc định là tên field

		// Phân tích các tag gorm
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag == "not null" {
				isNotNull = true
			} else if strings.HasPrefix(tag, "column:") {
				columnName = strings.TrimPrefix(tag, "column:")
			}
		}

		// Các kiểu dữ liệu pointer thường ngụ ý nullable, nên bỏ qua chúng
		if isNotNull && field.Type.Kind() != reflect.Ptr {
			requiredColumns = append(requiredColumns, columnName)
		}
	}
	// Cập nhật entityModelDictRequireCols
	RequiredColumnsCache[*tableName] = requiredColumns
	return requiredColumns
}

// decalre global var has dict store entyty model and map[string][]string
// like entityModelDict["tableName"] ["indexName"] = []string{"column1", "column2"}
var IndexColoumnsCache = make(map[string]map[string][]string)

// getColumnsGroupedByIndex takes an entity (struct) and returns a map where keys are index names
// and values are slices of column names that belong to that index.  It considers "index",
// "uniqueIndex", and "primaryKey" tags.
func GetColumnsGroupedByIndex(db *gorm.DB, entity interface{}) map[string][]string {
	// get table name from entity model
	tableName, terr := GetTableName(db, entity)
	if terr != nil {
		panic(terr)
	}
	if tableName == nil {
		panic("Cannot get table name from entity model")
	}
	if _, ok := IndexColoumnsCache[*tableName]; ok {

		return IndexColoumnsCache[*tableName]
	}

	columnsByIndex := make(map[string][]string)
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	// If it's a pointer, get the actual value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Ensure the input is a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Error: Input is not a struct or a pointer to a struct.")
		return columnsByIndex
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		gormTag := field.Tag.Get("gorm")
		tags := strings.Split(gormTag, ";")

		columnName := field.Name // Default column name is the field name
		indexName := ""

		// Parse the gorm tags
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if strings.HasPrefix(tag, "column:") {
				columnName = strings.TrimPrefix(tag, "column:")
			} else if strings.HasPrefix(tag, "index:") {
				indexName = strings.Split(strings.TrimPrefix(tag, "index:"), ",")[0]
			} else if strings.HasPrefix(tag, "uniqueIndex:") {
				indexName = strings.Split(strings.TrimPrefix(tag, "uniqueIndex:"), ",")[0]
			} else if tag == "primaryKey" || tag == "primary_key" {
				indexName = "PRIMARY" // Consistent name for primary key
			}
		}

		if indexName != "" {
			if _, ok := columnsByIndex[indexName]; !ok {
				columnsByIndex[indexName] = make([]string, 0)
			}
			columnsByIndex[indexName] = append(columnsByIndex[indexName], columnName)
		}
	}
	// Update entityModelDict
	IndexColoumnsCache[*tableName] = columnsByIndex
	return columnsByIndex
}

// AnalizeError phân tích lỗi GORM để xác định loại lỗi, tên bảng và các cột liên quan
// AnalizeError phân tích lỗi GORM để xác định loại lỗi, tên bảng và các cột liên quan
func AnalizeError(db *gorm.DB, entityModel interface{}, err error) *ErrorAnalysisResult {
	if db == nil {
		return &ErrorAnalysisResult{
			DbErrorType:  OtherError,
			ErrorMessage: err.Error(),
		}
	}

	result := ErrorAnalysisResult{
		DbErrorType:  NoError,
		ErrorMessage: "",
	}

	if err == nil {
		return nil
	}

	result.ErrorMessage = err.Error()
	dialect := db.Dialector.Name()

	stmt := &gorm.Statement{DB: db}
	stmt.Model = entityModel // Thiết lập model cho statement bằng cách gán field
	if err := stmt.Parse(stmt.Model); err == nil && stmt.Schema != nil {
		result.TableName = stmt.Schema.Table
	}

	switch dialect {
	case "mysql":
		if strings.Contains(result.ErrorMessage, "Duplicate entry") {
			result.DbErrorType = DuplicateError
			// Cố gắng trích xuất cột gây ra lỗi (dựa trên primary key)
			primaryKeys := GetPrimaryColumns(db, entityModel)
			if len(primaryKeys) > 0 && strings.Contains(result.ErrorMessage, ".PRIMARY'") {
				result.Columns = primaryKeys
			} else {
				// Cố gắng trích xuất cột từ thông báo lỗi (cho unique index khác)
				indexCols := GetColumnsGroupedByIndex(db, entityModel)
				parts := strings.Split(result.ErrorMessage, "for key '")
				if len(parts) > 1 {
					keyPart := parts[1]
					//kkey looks like "'table_name.index_name'" extract index_name
					keyPart = strings.Split(keyPart, ".")[1]
					keyPart = strings.TrimSuffix(keyPart, "'")

					result.Columns = indexCols[keyPart]

				}
			}
		} else if strings.Contains(result.ErrorMessage, "cannot be null") {
			result.DbErrorType = RequiredColumns
			// Cố gắng trích xuất tên cột bị null
			parts := strings.Split(result.ErrorMessage, "Column '")
			for _, part := range parts[1:] {
				columnParts := strings.SplitN(part, "' ", 2)
				if len(columnParts) > 0 {
					result.Columns = append(result.Columns, columnParts[0])
				}
			}
			// Lấy danh sách các cột required từ tag (có thể chính xác hơn)
			requiredFromTag := GetRequiredColumns(db, entityModel)
			if len(result.Columns) == 0 && len(requiredFromTag) > 0 {
				result.Columns = requiredFromTag
			}
		} else {
			result.DbErrorType = OtherError
		}
	case "postgres":
		//return not implemented yet
		panic("Not implemented yet")
	case "sqlite":
		panic("Not implemented yet")
	case "sqlserver":
		panic("Not implemented yet")
	default:
		result.DbErrorType = OtherError
	}

	return &result
}
