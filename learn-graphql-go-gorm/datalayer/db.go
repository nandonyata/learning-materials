package database

import (
	"fmt"
	"learn-graphql-go-gorm/datalayer/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection using GORM
func InitDB() {
	// Connection string — adjust user, password, host, port, dbname as needed
	dsn := "host=localhost user=postgres password=postgres dbname=graphql-go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("failed to get generic DB from GORM: %v", err)
	}

	// Verify connection
	if err = sqlDB.Ping(); err != nil {
		log.Panicf("failed to ping database: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")
	DB = db
}

// CloseDB closes the underlying SQL database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs GORM auto migrations for your models
func Migrate() {
	fmt.Println("🚀 Running GORM migrations...")

	// Example: add your model structs here
	err := DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("✅ Database migrated successfully!")
}
