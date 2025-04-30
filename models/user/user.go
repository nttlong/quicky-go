package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// models/user/user.go

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"-" gorm:"not null"` // Lưu bản hash, không hiển thị trong JSON
	Salt      string    `json:"-" gorm:"not null"` // Lưu salt, không hiển thị trong JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
