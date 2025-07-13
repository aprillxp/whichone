package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
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
	protected.GET("/players", controllers.GetAllPlayers)
	protected.GET("/players/:id", controllers.GetPlayerByID)
}
