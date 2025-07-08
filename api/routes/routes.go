package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "88Wager Backen Up")
	})

	// authentication
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
}
