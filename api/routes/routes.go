package routes

import (
	"api/controllers"
	"api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Connected!")
	})

	// authentication (unprotected route/endpoint)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// middleware
	protected := router.Group("/", middleware.AuthMiddleware())

	// protected route/endpoint
	protected.POST("/logout", controllers.Logout)

	protected.POST("/bank", controllers.RegisterBank)
	protected.POST("/wallet/topup", controllers.TopUpBalance)

	protected.POST("/bet", controllers.PlaceBet)
	protected.POST("/withdraw", controllers.Withdraw)
	protected.GET("/transaction", controllers.GetPlayerTransaction)

	protected.GET("/players", controllers.GetAllPlayers)
	protected.GET("/players/:id", controllers.GetPlayerByID)
	protected.GET("/players/me", controllers.GetMyProfile)
}
