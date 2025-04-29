package user

import (
	"fmt"
	"time"
)

// models/user/user.go

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate kiểm tra tính hợp lệ của User (ví dụ đơn giản)
func (u *User) Validate() error {
	if u.FirstName == "" || u.LastName == "" {
		return fmt.Errorf("first and last name cannot be empty")
	}
	return nil
}
