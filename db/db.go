package db

import (
	"fmt"
	"quicky-go/configs"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func ConnectGORM() (*gorm.DB, error) {
	var dsn string
	cfg := configs.Info.DB
	var err error

	switch cfg.DBType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBSchema,
		)
		GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlserver":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		if cfg.DBSchema != "" {
			dsn += fmt.Sprintf("&schema=%s", cfg.DBSchema)
		}
		GormDB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := GormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Check if the connection is successful
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return GormDB, nil
}
