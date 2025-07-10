package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "88Wager Backen Up")
	})

	// authentication (unprotected route/endpoint)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// middleware
	protected := router.Group("/", middleware.AuthMiddleware())

	// protected route/endpoint
	protected.POST("/logout", controllers.Logout)
}
