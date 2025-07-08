package database

import "log"

func Reset() {
	// Connect ke DB
	ConnectDB()

	// Drop semua tabel
	// err := DB.Migrator().DropTable(&models.User{}, &models.Wallet{}, &models.Bet{})
	// if err != nil {
	// 	log.Fatal("❌ Failed to drop table:", err)
	// }
	// log.Println("✅ All table dropped successfully")

	// // Migrasi ulang
	// // err = DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Bet{})
	// if err != nil {
	// 	log.Fatal("❌ Failed to migrate:", err)
	// }
	log.Println("✅ Migrating is successfully")
}
