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

// CustomClaims biar lebih aman daripada pake MapClaims mentah
type CustomClaims struct {
	PlayerID uint `json:"player_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Cek Blacklist di Redis
		redisClient := configs.GetRedisClient()
		_, err := redisClient.Get(configs.RedisCtx, "jwt_blacklist:"+tokenString).Result()
		if err == nil {
			// token ada di blacklist
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}
		if err != redis.Nil { // error selain key not found
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error during token check", "details": err.Error()})
			c.Abort()
			return
		}

		//Parse Token
		token, err := utils.ParseKeyWithClaims(tokenString, &CustomClaims{})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		//Ambil Claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		//Simpan ke Context
		c.Set("player_id", claims.PlayerID)

		c.Next()
	}
}
