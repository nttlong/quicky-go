package tenants

// quản lý các tenant trong hệ thống
import (
	"fmt"
	"quicky-go/manager/accounts"
	"quicky-go/models/tenants"
	"quicky-go/repo"
	"regexp"
	"strings"
	"time"
	"unicode"

	repo_err "quicky-go/repo/errors"

	"github.com/google/uuid"
)

// Định nghĩa custom error struct
type TenantNameError struct {
	Input   string
	Length  int
	Time    time.Time
	Message string
}

// Implement interface error cho InputError
func (e *TenantNameError) Error() string {
	return fmt.Sprintf("Input error: %s, length: %d, at: %v - %s", e.Input, e.Length, e.Time, e.Message)
}
func ValidateTenantName(tenantName string) (bool, string) {
	if strings.TrimSpace(tenantName) == "" {
		return false, "Tenant name cannot be empty or just whitespace."
	}

	if len(tenantName) < 3 {
		return false, "Tenant name must be at least 3 characters long."
	}

	if len(tenantName) > 63 { // Common limit for subdomain/database names
		return false, "Tenant name cannot be longer than 63 characters."
	}

	// Check for invalid characters (only allow lowercase letters, numbers, and hyphens/underscores)
	for _, r := range tenantName {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false, "Tenant name can only contain lowercase letters, numbers, hyphens (-), and underscores (_)."
		}
	}

	// Ensure it doesn't start or end with a hyphen or underscore
	if strings.HasPrefix(tenantName, "-") || strings.HasPrefix(tenantName, "_") ||
		strings.HasSuffix(tenantName, "-") || strings.HasSuffix(tenantName, "_") {
		return false, "Tenant name cannot start or end with a hyphen (-) or underscore (_)."
	}

	// Ensure no consecutive hyphens or underscores
	if strings.Contains(tenantName, "--") || strings.Contains(tenantName, "__") || strings.Contains(tenantName, "-_") || strings.Contains(tenantName, "_-") {
		return false, "Tenant name cannot contain consecutive hyphens or underscores."
	}

	// Optional: Check if the name is a reserved keyword (e.g., "www", "admin", "api")
	reservedKeywords := []string{"www", "admin", "api", "root", "public"}
	lowerTenantName := strings.ToLower(tenantName)
	for _, keyword := range reservedKeywords {
		if lowerTenantName == keyword {
			return false, fmt.Sprintf("Tenant name '%s' is a reserved keyword.", tenantName)
		}
	}

	// Optional: Check if the name matches a specific pattern (e.g., using regex)
	pattern := `^[a-z0-9]+(?:[-_][a-z0-9]+)*$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(tenantName) {
		return false, "Tenant name does not match the required pattern (lowercase alphanumeric with single hyphens or underscores as separators)."
	}

	return true, "" // Tenant name is valid
}
func CreateTenant(name string, descrition string) (*tenants.Tenants, error) {
	// Kiểm tra tên tenant
	isValid, message := ValidateTenantName(name)
	if !isValid {
		return nil, &TenantNameError{Input: name, Length: len(name), Time: time.Now(), Message: message}
	}

	// tạo tenant mới
	t := tenants.Tenants{
		ID:          uuid.New(),
		Name:        name,
		Description: descrition,
		Status:      1,
		DbTenant:    strings.ToLower(name),
		CreatedBy:   "admin",
		ModifiedBy:  "admin",
		CreatedOn:   time.Now().UTC(),
		ModifiedOn:  time.Now().UTC(),
	}
	// lưu vào cơ sở dữ liệu
	var repoManager = repo.GetManagerRepo()
	err := repoManager.Create(&t).Error
	if err != nil {

		// check duplicate tenant name
		errType := repo_err.AnalizeError(repoManager, &t, err)
		if errType.DbErrorType == repo_err.DuplicateError {
			if errType.Columns[0] == "ID" {
				repoManager.Where("ID = ?", t.ID).First(&t)

			} else if errType.Columns[0] == "Name" {
				repoManager.Where("NameD = ?", t.Name).First(&t)
			} else {
				return nil, err
			}

		} else {
			return nil, err
		}

	}

	_, errGet := repo.GetRepo(t.DbTenant)
	if errGet != nil {
		return nil, errGet
	}
	// tạo database cho tenant mới

	_, qErr := accounts.CreateAccount(t.DbTenant, "admin", "admin", "admin")
	if qErr.DbErrorType == repo_err.DuplicateError {
		if qErr.Columns[0] == "Username" || qErr.Columns[0] == "Email" {
			return &t, nil
		}
	} else {
		return nil, qErr
	}

	return &t, nil

}
