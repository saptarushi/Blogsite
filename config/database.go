package config

import (
	"Blogsite/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=Postgresql@1234 dbname=blogsite_db port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automigrate models
	err = DB.AutoMigrate(&models.User{}, &models.Blog{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	fmt.Println("Database connection established and models migrated!")
}
