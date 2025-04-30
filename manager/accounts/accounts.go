package accounts

import (
	"crypto/rand"
	"fmt"
	"quicky-go/models/account"
	"quicky-go/models/bases"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPasswordWithSalt(password, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func ComparePasswordWithSalt(password, hashedPassword, salt string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err
}
func generateSalt() (string, error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes for salt: %w", err)
	}
	return fmt.Sprintf("%x", saltBytes), nil // Convert bytes to hexadecimal string
}

// CreateAccount creates a new account in the database.
func CreateAccount(db *gorm.DB, username, email, plainPassword string) (*account.Account, error) {
	// 1. Generate a unique salt
	salt, err := generateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// 2. Hash the password with the salt
	hashedPassword, err := HashPasswordWithSalt(plainPassword, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create a new Account instance
	newAccount := account.Account{
		BaseModel: bases.BaseModel{
			ID:         uuid.New(),
			CreatedBy:  "admin",
			ModifiedBy: "admin",
			CreatedOn:  time.Now(),
			ModifiedOn: time.Now(),
		},
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Salt:     salt,
	}

	// 4. Create the account in the database using GORM
	result := db.Create(&newAccount)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create account in database: %w", result.Error)
	}

	return &newAccount, nil
}
