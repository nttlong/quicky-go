package repo_types

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/*
*
  - Struct Column để lưu thông tin của một cột trong một bảng.
    storegae information of a column in a table.
  - Lưu các thông tin như tên cột, kiểu dữ liệu, kích thước, có thể null hay không, tên index, có unique hay không.
  - Store information like column name, data type, size, can be null or not, index name, have unique or not.
*/
type Column struct {
	/* Tên cột, column name (property name of entity) */
	Name string
	/* Colum type, data type is in "number,text,date,boolean") */
	Type string
	/* Allow null or not, can be null or not */
	AllowNull bool
	/* Length of column if null unlimited (depends on database driver), length of column if not null */
	Length *int
	/* Index name, index name */
	IndexName string
	/* Is unique or not, unique or not */
	IsUnique bool
}
type ErrorCode int
type DbAction int

// Định nghĩa các hằng số cho các giá trị enum.
// Define enum values for ErrorCode and DbAction.
const (
	Unknown ErrorCode = iota // Lỗi không xác định, Unknown error
	// Lỗi trùng lặp, Duplicate error
	// Duplicate error
	Duplicate // iota tự động tăng giá trị của mỗi hằng số
	// Lỗi tham chiếu, Reference error, lỗ này thường gây ra do relation ship giữa các bảng.
	// Reference error, usually caused by relation ship between tables.

	Reference
	// Lỗi bắt buộc, Required error, lỗi này thường gây ra khi một thông tin cần thiết bị không được cung cấp.
	// Required error, usually caused by missing required data.
	Require
	// Lỗi kích thước, Invalid length error, lỗi này thường gây ra khi một thông tin có kích thước không hợp lệ.
	// Invalid length error, usually caused by invalid data length.
	InvalidLen
)

// Định nghĩa các hằng số cho các thao tác trên dabase gây ra lỗi.
const (
	// nguyên nhân do insert
	Insert DbAction = iota
	// nguyên nhân do update
	Update
	// nguyên nhân do delete
	Delete
)

// Hàm String để làm cho enum dễ đọc hơn khi in. make ErrorCode enum more readable
// String() method for ErrorCode enum
func (e ErrorCode) String() string {
	switch e {
	case Duplicate:
		return "Duplicate"
	case Reference:
		return "Reference"
	case Require:
		return "Require"
	case InvalidLen:
		return "InvalidLen"
	default:
		return "Unknown"
	}
}

// Hàm String để làm cho enum dễ đọc hơn khi in. make DbAction enum more readable
func (e DbAction) String() string {
	switch e {
	case Insert:
		return "Insert"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// Cấu trúc ghi nhận lại toàn bộ lỗi gây ra khi thao tác trên database.
// Structure to catch all errors occurred when operate on database.
type DataActionError struct {
	Err error
	// Thao tác gây ra lỗi, action cause error
	Action string
	// Mã lỗi, error code
	Code ErrorCode
	// các cột liên quan đến lỗi nếu có, reference columns cause error
	RefColumns []string // Các cột liên quan đến lỗi nếu có , Reference columns cause error
	// các bảng liên quan đến lỗi nếu có, reference table name cause error
	RefTableName string // Tên bảng liên quan đến lỗi nếu có, Reference table name cause error

}

// hàm diễn dịch lại lỗi gây ra khi thao tác trên database.
// translate DataActionError to readable message
func (e *DataActionError) Error() string {
	msg := "Error " + e.Action + " " + e.Code.String() + " " + e.Err.Error()
	if len(e.RefColumns) > 0 {
		msg += " RefColumns: " + strings.Join(e.RefColumns, ",")
	}
	if e.RefTableName != "" {
		msg += " RefTableName: " + e.RefTableName
	}
	return msg
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
func ComputeColumns(typ reflect.Type) ([]Column, error) {

	var columns []Column
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

		col := Column{
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
			if tag == "primary_key" || strings.Contains(gormTag, ";primaryKey;") || strings.Contains(gormTag, ";primary_key;") {
				col.IsUnique = true
			}
			//tag=="uniqueIndex:idx_name,length:191"
			if strings.Contains(tag, "uniqueIndex:") {

				indexName := strings.Split(strings.Split(tag, ":")[1], ",")[0]
				col.IndexName = indexName
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
func GetTypeNameOfEntity(entity interface{}) string {
	typ := reflect.TypeOf(entity)
	tableName := typ.String()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		tableName = typ.String()
	}
	return tableName
}
func GetTableNameOfEntity(entity interface{}) string {
	tableName := GetTypeNameOfEntity(entity)
	if strings.Contains(tableName, ".") {
		tableName = strings.Split(tableName, ".")[1]
	}
	return tableName
}
func GetReflectType(entity interface{}) (reflect.Type, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		return typ, nil

	}
	return nil, errors.New("entity must be a pointer to struct")
}
