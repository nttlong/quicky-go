package info

type Column struct {
	Name      string
	Type      string
	AllowNull bool
	Length    *int
	IndexName string
	IsUnique  bool
}
type ErrorCode int
type DbAction int

// Định nghĩa các hằng số cho các giá trị enum.
const (
	Duplicate ErrorCode = iota // iota tự động tăng giá trị của mỗi hằng số
	Reference
	Require
	InvalidLen
)
const (
	Insert DbAction = iota
	Update
	Delete
)

// Hàm String để làm cho enum dễ đọc hơn khi in. make ErrorCode enum more readable
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

type DataActionError struct {
	Err          error
	Action       string
	Code         ErrorCode
	RefColumns   []string // Các cột liên quan đến lỗi nếu có , Reference columns cause error
	RefTableName string   // Tên bảng liên quan đến lỗi nếu có, Reference table name cause error

}
