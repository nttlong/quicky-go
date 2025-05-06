package dbconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"gopkg.in/yaml.v2"
)

type IDbConfigBase interface {
	GetUser() string
	SetUser(User string)
	GetPassword() string
	SetPassword(password string)
	GetHost() string
	SetHost(host string)
	GetPort() string
	SetPort(port string)
	GetOptions() map[string]string
	SetOptions(options map[string]string)
	LoadFromYamlFile(yamlFile string) error
	CheckIsLoaded() bool
	//to json strin with pretty format password show as ****
	ToJSON() string
	GetAllColumnsInfoFromEntity(entity interface{}) []ColumInfo
	GetColumInfoOfField(reflect.StructField) *ColumInfo
	GetAllModelsInEntity(entity interface{}) []interface{}
}
type IStorage interface {
	AutoMigrate(entity interface{}) error
}
type IDbConfig interface {
	IDbConfigBase
	GetConectionString(dbname string) string
	GetConectionStringNoDatabase() string
	PingDb() error
	CreateDbIfNotExist(dbname string) error
	GetStorage(dbName string) (IStorage, error)
}
type DbConfigBase struct {
	User string `yaml:"user"`

	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`

	Options  map[string]string `yaml:"options"`
	IsLoaded bool
}

func (c *DbConfigBase) GetUser() string {
	if !c.IsLoaded {
		panic("DbConfigBase is not loaded,pleas call LoadFromYamlFile first")
	}
	return c.User
}

func (c *DbConfigBase) SetUser(User string) {

	c.User = User
}

func (c *DbConfigBase) GetPassword() string {
	if !c.IsLoaded {
		panic("DbConfigBase is not loaded,pleas call LoadFromYamlFile first")
	}
	return c.Password
}

func (c *DbConfigBase) SetPassword(password string) {

	c.Password = password
}

func (c *DbConfigBase) GetHost() string {
	if !c.IsLoaded {
		panic("DbConfigBase is not loaded,pleas call LoadFromYamlFile first")
	}
	return c.Host
}

func (c *DbConfigBase) SetHost(host string) {
	c.Host = host
}

func (c *DbConfigBase) GetPort() string {
	if !c.IsLoaded {
		panic("DbConfigBase is not loaded,pleas call LoadFromYamlFile first")
	}
	return c.Port
}
func (c *DbConfigBase) SetPort(port string) {
	c.Port = port
}
func (c *DbConfigBase) GetOptions() map[string]string {
	if !c.IsLoaded {
		panic("DbConfigBase is not loaded,pleas call LoadFromYamlFile first")
	}
	return c.Options
}
func (c *DbConfigBase) SetOptions(options map[string]string) {
	c.Options = options
}
func (c *DbConfigBase) CheckIsLoaded() bool {
	return c.IsLoaded
}

func (c *DbConfigBase) LoadFromYamlFile(yamlFile string) error {
	content, err := os.ReadFile(yamlFile)
	if err != nil {
		return err
	}
	var config map[string]map[string]interface{}
	err = yaml.Unmarshal(content, &config)

	if err != nil {
		return err
	}
	//dbConfig, ok := config["db"].(map[string]interface{})
	bffContent, err := yaml.Marshal(config["db"]) // Use yaml.Marshal
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bffContent, &c)

	if c.Password == "" {
		return errors.New(fmt.Sprintf("Password is empty in the config file %s", yamlFile))
	}
	if c.Host == "" {
		return errors.New(fmt.Sprintf("Host is empty in the config file %s", yamlFile))
	}
	if c.Port == "" {
		return errors.New(fmt.Sprintf("Port is empty in the config file %s", yamlFile))
	}
	if c.User == "" {
		return errors.New(fmt.Sprintf("User is empty in the config file %s", yamlFile))
	}
	if c.Options == nil {
		c.Options = make(map[string]string)
	}
	c.IsLoaded = true
	return nil
}
func (c *DbConfigBase) ToJSON() string {
	// use json.MarshalIndent to format the json string with indent
	jsonStr, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonStr)

}

var (
	cacheColumnInfo = make(map[reflect.Type][]ColumInfo)
)

type ColumInfo struct {
	Content      string
	Name         string
	DbType       string
	Length       int
	IsPk         bool
	IndexName    string
	IsUnique     bool
	Typ          reflect.StructField
	DefaultValue string
	HasDefault   bool
}

func (c *DbConfigBase) GetColumInfoOfField(field reflect.StructField) *ColumInfo {
	return getColumInfoOfField(field)
}
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	runes := []rune(s) // Xử lý chuỗi dưới dạng rune để hỗ trợ Unicode

	// Duyệt qua từng ký tự
	for i, r := range runes {
		// Nếu ký tự hiện tại là chữ hoa
		if unicode.IsUpper(r) {
			// Thêm dấu gạch dưới trước ký tự chữ hoa, trừ khi:
			// - Đây là ký tự đầu tiên (i == 0)
			// - Ký tự trước đó đã là dấu gạch dưới
			if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '_' {
				result.WriteRune('_')
			}
			// Chuyển ký tự thành chữ thường và thêm vào kết quả
			result.WriteRune(unicode.ToLower(r))
		} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			// Nếu không phải chữ cái hoặc số, thay thế bằng dấu gạch dưới
			// Tránh thêm nhiều dấu gạch dưới liên tiếp
			if result.Len() == 0 || result.String()[result.Len()-1] != '_' {
				result.WriteRune('_')
			}
		} else {
			// Nếu là chữ thường hoặc số, thêm trực tiếp
			result.WriteRune(r)
		}
	}

	// Loại bỏ dấu gạch dưới ở đầu và cuối, đồng thời thay thế nhiều dấu gạch dưới liên tiếp bằng một dấu
	snake := strings.Trim(result.String(), "_")
	snake = strings.ReplaceAll(snake, "__", "_")

	return snake
}
func getColumInfoOfField(field reflect.StructField) *ColumInfo {

	tag := field.Tag.Get("gorm")
	if tag == "" {
		return nil
	}
	if strings.HasPrefix(tag, "foreignKey:") {
		return nil
	}

	tags := strings.Split(tag, ";")
	ret := ColumInfo{
		Content: tag,
		Typ:     field,
		Name:    toSnakeCase(field.Name),
		Length:  -1,
	}
	for _, t := range tags {
		if strings.Contains(t, ":") {
			key := strings.Split(t, ":")[0]
			value := strings.Split(t, ":")[1]
			var strLen *string = nil
			if key == "column" {
				ret.Name = value

			}
			if key == "type" {
				if strings.Contains(value, "(") && strings.Contains(value, ")") {
					ret.DbType = strings.Split(value, "(")[0]

					strLen = &strings.Split(strings.Split(value, "(")[1], ")")[0]
					length, err := strconv.Atoi(*strLen)
					if err == nil {
						ret.Length = length
					}
				}

			}

			if key == "index" {
				ret.IndexName = value
			}
			if key == "uniqueIndex" {
				ret.IsUnique = true
				ret.IndexName = value
			}
			if key == "unique" {
				ret.IsUnique = true
			}
			if key == "default" {
				ret.DefaultValue = value
				ret.HasDefault = true
			}
		}
		if t == "primaryKey" {
			ret.IsPk = true
		}
	}

	return &ret
}

var (
	cacheCoummsInfo = make(map[reflect.Type][]ColumInfo)
	lockCoummsInfo  = new(sync.RWMutex)
)

func (c *DbConfigBase) GetAllColumnsInfoFromEntity(entity interface{}) []ColumInfo {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	lockCoummsInfo.RLock()
	ret, ok := cacheCoummsInfo[typ]
	lockCoummsInfo.RUnlock()
	if ok {
		return ret
	}
	lockCoummsInfo.Lock()
	defer lockCoummsInfo.Unlock()
	ret = getAllColumnsInfoFromEntity(typ)
	cacheCoummsInfo[typ] = ret
	return ret

}
func (c *DbConfigBase) GetAllModelsInEntity(entity interface{}) []interface{} {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	var models []interface{}
	models = append(models, entity)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct {
			if strings.HasPrefix(field.Tag.Get("gorm"), "foreignKey:") {
				//create new entity
				newEntity := reflect.New(field.Type).Interface()
				models = append(models, newEntity)
			}

		}
	}
	return models
}

// =
func NewDbConfigBase() IDbConfigBase {
	return &DbConfigBase{}
}

// ===================================================
func getAllColumnsInfoFromEntity(typ reflect.Type) []ColumInfo {
	//scan all fields of the entity
	var columInfos []ColumInfo = make([]ColumInfo, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct {
			if strings.HasPrefix(field.Tag.Get("gorm"), "foreignKey:") {
				continue
			}
			subColumInfos := getAllColumnsInfoFromEntity(field.Type)

			columInfos = append(columInfos, subColumInfos...)
		} else {
			columInfo := getColumInfoOfField(field)
			if columInfo == nil {
				continue
			}
			//add columInfo to columInfos
			columInfos = append(columInfos, *columInfo)
			//columInfos.append(columInfo)
			//columInfos = append(columInfos, {columInfo})
		}
	}
	return columInfos
}
