package config

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/websoket-chat/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *Config) *gorm.DB {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "prefer"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.SSLMode,
	)

	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	sqlDB, _ := db.DB()
	maxIdleConns := 30
	if cfg.SetMaxIdleConns != "" {
		maxIdleConns, _ = strconv.Atoi(cfg.SetMaxIdleConns)
	}
	MaxOpenConns := 100
	if cfg.SetMaxOpenConns != "" {
		MaxOpenConns, _ = strconv.Atoi(cfg.SetMaxOpenConns)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(MaxOpenConns)

	log.Info("Database connected")

	DB = db

	return db
}

func MigrateTable(db *gorm.DB) (err error) {
	// AutoMigrate will create tables, missing columns, and missing indexes
	err = db.AutoMigrate(
		&model.Message{},
		&model.User{},
		&model.Role{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}

	return nil
}
