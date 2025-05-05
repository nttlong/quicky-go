package account

import (
	"vngom/models/bases"

	"golang.org/x/crypto/bcrypt"
)

// models/user/user.go
type Account struct {
	bases.BaseModel
	Username string `gorm:"type:varchar(191);uniqueIndex:idx_username;"`
	Email    string `gorm:"type:varchar(191);uniqueIndex:idx_email;"`
	Password string `gorm:"type:varchar(191);"`
	Salt     string `json:"-" gorm:"not null;"` // Lưu salt, không hiển thị trong JSON
}

// TableName sets the desired table name

func HashPasswordWithSalt(password, salt string) (string, error) {
	// **LƯU Ý QUAN TRỌNG:** Trong ứng dụng thực tế, bạn NÊN sử dụng bcrypt thay vì cách này.
	// bcrypt tự động tạo salt an toàn và tích hợp nó vào hash.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// comparePasswordWithSalt so sánh mật khẩu đã cho với mật khẩu băm và salt đã lưu trữ
func ComparePasswordWithSalt(password, hashedPassword, salt string) error {
	// **LƯU Ý QUAN TRỌNG:** Tương tự, trong ứng dụng thực tế, bạn NÊN sử dụng bcrypt.CompareHashAndPassword.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err
}
