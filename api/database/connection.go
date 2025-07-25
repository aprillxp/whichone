package database

import (
	"api/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connect error:", err)
	}
	log.Println("✅ Connected to Database")

	// hanya migrasi tanpa drop table
	if err := DB.AutoMigrate(&models.Player{}, &models.Bank{}, &models.Wallet{}, &models.Bet{}, &models.Transaction{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
}
