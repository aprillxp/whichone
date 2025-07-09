package database

import (
	"api/models"
	"log"
)

func Reset() {
	// Connect ke DB
	ConnectDB()

	// Drop semua tabel
	err := DB.Migrator().DropTable(&models.Player{}, &models.Wallet{}, &models.Bank{})
	if err != nil {
		log.Fatal("❌ Failed to drop table:", err)
	}
	log.Println("✅ All table dropped successfully")

	// Migrasi ulang
	err = DB.AutoMigrate(&models.Player{}, &models.Wallet{}, &models.Bank{})
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}
	log.Println("✅ Migrating is successfully")
}
