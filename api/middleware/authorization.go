package middleware

import (
	"api/configs"
	"api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token is required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		redisClient := configs.GetRedisClient()
		_, err := redisClient.Get(configs.RedisCtx, "jwt_blacklist"+tokenString).Result()
		if err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}
		if err != redis.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Redis error during check", "details": err.Error()})
			c.Abort()
			return
		}

		token, err := utils.ParseKey(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token"})
			c.Abort()
			return
		}

		playerIDconv, ok := claims["player_id"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Player ID not found or Invalid in token"})
			c.Abort()
			return
		}
		c.Set("player_id", uint(playerIDconv))

		c.Next()
	}
}
