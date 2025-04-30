package tenants

// quản lý các tenant trong hệ thống
import (
	"fmt"
	"quicky-go/models/tenants"
	"quicky-go/repo"
	"regexp"
	"strings"
	"time"
	"unicode"
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
		Name:        name,
		Description: "Mô tả của tenant",
		Status:      1,
		DbTenant:    name.lower(),
		CreatedBy:   "admin",
		ModifiedBy:  "admin",
		CreatedOn:   time.Now(),
		ModifiedOn:  time.Now(),
	}
	// lưu vào cơ sở dữ liệu
	var repo = repo.GetManagerRepo()
	err := repo.Create(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil

}
