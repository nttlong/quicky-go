package accounts

import (
	"crypto/rand"
	"fmt"
	"quicky-go/models/account"
	"quicky-go/models/bases"
	"quicky-go/repo"
	"time"

	qErr "quicky-go/repo/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
func CreateAccount(dbName string, username, email, plainPassword string) (*account.Account, *qErr.ErrorAnalysisResult) {
	// 1. Generate a unique salt
	salt, err := generateSalt()
	if err != nil {
		return nil, qErr.AnalizeError(nil, nil, err)
	}

	// 2. Hash the password with the salt
	hashedPassword, err := HashPasswordWithSalt(plainPassword, salt)
	if err != nil {
		return nil, qErr.AnalizeError(nil, nil, err)
	}

	// 3. Create a new Account instance
	newAccount := account.Account{
		BaseModel: bases.BaseModel{
			ID:         uuid.New(),
			CreatedBy:  "admin",
			ModifiedBy: "admin",
			CreatedOn:  time.Now().UTC(),
			ModifiedOn: time.Now().UTC(),
		},
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Salt:     salt,
	}

	// 4. Create the account in the database using GORM

	qErr := repo.Insert(dbName, &newAccount)
	if qErr != nil {
		return nil, qErr
	}

	return &newAccount, nil
}
func CheckIsAccountExist(dbName string, username, email string) (bool, error) {
	db, err := repo.GetRepo(dbName)
	if err != nil {
		return false, fmt.Errorf("failed to get repository: %w", err)
	}
	var count int64
	result := db.Model(&account.Account{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to count account in database: %w", result.Error)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
