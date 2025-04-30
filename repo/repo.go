package repo

import (
	"fmt"
	"log"

	"quicky-go/configs"

	"gorm.io/driver/mysql" // Hoặc driver của bạn (postgres, sqlserver)
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// global variable for gorm db connection used in all package
var Repo *gorm.DB

func init() {
	var err error
	cfg := configs.Info.DB
	dsn := ""

	switch cfg.DBType {
	case configs.DBTypeMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		Repo, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case configs.DBTypePostgres:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBSchema,
		)
		Repo, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case configs.DBTypeSQLServer:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		if cfg.DBSchema != "" {
			dsn += fmt.Sprintf("&schema=%s", cfg.DBSchema)
		}
		Repo, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database type: %s", cfg.DBType)
		return
	}

	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", cfg.DBType, err)
	}

	log.Printf("Successfully connected to %s database.", cfg.DBType)

	// Optional: AutoMigrate your models here if needed
	// if err := Repo.AutoMigrate(&YourModel1{}, &YourModel2{}); err != nil {
	// 	log.Fatalf("Failed to auto migrate models: %v", err)
	// }
}
