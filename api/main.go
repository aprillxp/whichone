package main

import (
	"api/database"
	"api/models"
	"api/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if os.Getenv("RESET") == "true" {
		database.Reset()
		return
	}

	database.ConnectDB()
	database.DB.AutoMigrate(&models.Player{}, &models.Bank{}, &models.Wallet{})

	router := gin.Default()
	routes.Routes(router)

	log.Println("I Love You 2210")
	log.Fatal(router.Run(":8080"))

}
