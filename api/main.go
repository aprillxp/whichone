// main.go
package main

import (
	"api/configs" // Import configs untuk Redis
	"api/database"
	"api/models" // Untuk AutoMigrate
	"api/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from system.")
	} else {
		log.Println(".env file loaded successfully.")
	}
}

func main() {
	if os.Getenv("RESET_DB") == "true" {
		log.Println("⚠️ RESET_DB is true. Resetting database...")
		database.Reset()
		log.Println("✅ Database reset complete. Exiting.")
		return
	}

	database.ConnectDB()

	database.DB.AutoMigrate(&models.Player{}, &models.Bank{}, &models.Wallet{})

	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	// 2. Koneksi Redis
	configs.ConnectRedis()
	defer configs.CloseRedis()

	router := gin.Default()
	routes.Routes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(router.Run(":" + port))
}
