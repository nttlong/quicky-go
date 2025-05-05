package accounts

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
	"time"
	"vngom/models/account"
	"vngom/models/bases"
	"vngom/repo"
	"vngom/repo/repo_types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountManager struct {
	repo.IRepo
}

func generateSalt() (string, error) {
	saltBytes := make([]byte, 16) // Use 16 bytes for a strong salt
	_, err := io.ReadFull(rand.Reader, saltBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(saltBytes), nil
}
func hashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword trả về byte slice của hash.
	// Tham số thứ hai là cost, càng cao thì càng tốn thời gian băm, nhưng càng an toàn.
	// Cost 10 là một giá trị hợp lý cho nhiều ứng dụng.
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Log lỗi để bạn có thể thấy nó trong console hoặc log file.
		log.Printf("Lỗi khi băm mật khẩu: %v", err)
		return "", err // Trả về lỗi để người gọi hàm có thể xử lý.
	}

	// Chuyển đổi byte slice thành string trước khi trả về.
	hashedString := string(hashedBytes)
	return hashedString, nil
}
func validatePassword(password, storedHash string) bool {
	// bcrypt.CompareHashAndPassword xử lý việc lấy salt từ hash đã lưu trữ
	// và so sánh hash của mật khẩu đã cho với hash đã lưu trữ.
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		// Nếu có lỗi, có thể là do hash không hợp lệ hoặc mật khẩu không khớp.
		log.Printf("Lỗi so sánh mật khẩu: %v", err)
		return false // Trả về false cho cả hai trường hợp để tránh tiết lộ thông tin
	}

	// Nếu không có lỗi, mật khẩu khớp.
	return true
}
func (m *AccountManager) HashPassword(acc *account.Account) error {
	//generate salt

	hashPass, err := hashPassword(acc.Password)
	if err != nil {
		return err
	}
	acc.Password = hashPass
	return nil
}

func (m *AccountManager) CreateAccount(name string, email string, password string) (*account.Account, *repo_types.DataActionError) {

	acc := account.Account{
		Username: name,
		Email:    email,

		BaseModel: bases.BaseModel{
			ID:         uuid.New(),
			CreatedOn:  time.Now(),
			ModifiedOn: time.Now(),
			ModifiedBy: "admin",
			CreatedBy:  "admin",
		},
	}
	m.HashPassword(&acc)
	startAt := time.Now()
	err := m.Insert(&acc)
	elapseTime := time.Since(startAt)
	fmt.Println("CreateAccount time in ms: ", elapseTime.Milliseconds())
	if err != nil {
		if err.Code == repo_types.Duplicate {
			err = m.Get(&acc, `username =?`, acc.Username)

			if err != nil {
				return nil, err
			}
			return &acc, nil

		} else {
			return nil, err
		}

	}
	return &acc, nil
}
func (m *AccountManager) ValidateAccount(name string, password string) (*account.Account, *repo_types.DataActionError) {
	acc := account.Account{}
	err := m.IRepo.Get(&acc, "Username = ?", name)
	if err != nil {
		return nil, err
	}
	hp, _ := hashPassword(password)
	fmt.Print(hp)

	if !validatePassword(password, acc.Password) {
		return nil, m.GetError(errors.New("Invalid password"), reflect.TypeOf(acc), "account", "validate")
	}
	return &acc, nil
}

var (
	cacheAccountManager    map[string]*AccountManager = make(map[string]*AccountManager)
	cacheAccountManageLock sync.RWMutex
)

func NewAccountManager(repo repo.IRepo) *AccountManager {
	// Lock đọc
	cacheAccountManageLock.RLock()
	if cache, ok := cacheAccountManager[repo.GetDbName()]; ok {
		cacheAccountManageLock.RUnlock() // Giải phóng lock đọc ngay lập tức
		return cache
	}
	cacheAccountManageLock.RUnlock()

	// Lock ghi chỉ khi cần thiết
	cacheAccountManageLock.Lock()
	defer cacheAccountManageLock.Unlock()

	// Kiểm tra lại lần nữa trong khi đã có khóa ghi
	if cache, ok := cacheAccountManager[repo.GetDbName()]; ok {
		return cache
	}

	// Tạo mới và lưu vào cache
	cacheAccountManager[repo.GetDbName()] = &AccountManager{
		IRepo: repo,
	}
	return cacheAccountManager[repo.GetDbName()]
}
