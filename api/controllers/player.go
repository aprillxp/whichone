package controllers

import (
	"api/configs"
	"api/database"
	"api/models"
	"api/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request", "details": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		player := models.Player{
			Username: req.Username,
			Password: hashedPassword,
		}
		if err := tx.Create(&player).Error; err != nil {
			// cek unique contraint
			if err.Error() == `ERROR: Duplicate key value violates unique constraint` {
				c.JSON(http.StatusConflict, gin.H{"error": "Username already exist"})
				return err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return err
		}

		wallet := models.Wallet{
			PlayerID: player.ID,
			Balance:  0,
		}
		if err := tx.Create(&wallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet", "details": err.Error()})
			return err
		}
		return nil
	})

	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request", "details": err.Error()})
		return
	}

	var player models.Player
	if err := database.DB.Where("username = ?", req.Username).First(&player).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if !utils.CheckHashedPassword(req.Password, player.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateJWT(player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successfully",
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorize: Bearer token required"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := utils.ParseKey(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired token", "details": err.Error()})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired token", "details": "Token is invalid"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse key"})
		return
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse key"})
		return
	}
	expTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expTime)

	if duration <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Token already expired"})
		return
	}

	redisClient := configs.GetRedisClient()
	err = redisClient.SetEX(configs.RedisCtx, "jwt_blacklist"+tokenString, "ture", duration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token in redis", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout Successfully"})
}

func GetAllPlayers(c *gin.Context) {
	db := database.DB

	var players []models.Player
	query := db.Model(&models.Player{})

	// filter by query params
	if username := c.Query("username"); username != "" {
		query = query.Where("username ILIKE ?", "%"+username+"%")
	}
	if namaRekening := c.Query("nama_rekening"); namaRekening != "" {
		query = query.Joins("LEFT JOIN banks ON players.id = banks.player_id").Where("banks.nama_rekening ILIKE ?", "%"+namaRekening+"%")
	}
	if nomorRekening := c.Query("nomor_rekening"); nomorRekening != "" {
		query = query.Joins("LEFT JOIN banks ON players.id = banks.player_id").Where("banks.nomor_rekening ILIKE ?", "%"+nomorRekening+"%")
	}
	if namaBank := c.Query("nama_bank"); namaBank != "" {
		query = query.Joins("LEFT JOIN banks ON player.id = banks.player_id").Where("banks.nama_bank ILIKE ?", "%"+namaBank+"%")
	}
	if registerAt := c.Query("register_at"); registerAt != "" {
		query = query.Where("DATE(players.created_at) = ?", registerAt)
	}

	if minBalanceStr := c.Query("min_balance"); minBalanceStr != "" {
		if minBalance, err := strconv.ParseFloat(minBalanceStr, 64); err == nil {
			query = query.Joins("LEFT JOIN wallets ON players.id = wallets.player.id").Where("wallets.balance >= ?", minBalance)
		}
	}
	if maxBalanceStr := c.Query("min_balance"); maxBalanceStr != "" {
		if maxBalance, err := strconv.ParseFloat(maxBalanceStr, 64); err == nil {
			query = query.Joins("LEFT JOIN wallets ON players.id = wallets.player.id").Where("wallets.balance <= ?", maxBalance)
		}
	}

	query = query.Preload("Bank").Preload("Wallet")
	if err := query.Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreived player", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})
}

func GetPlayerByID(c *gin.Context) {
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Player ID format"})
		return
	}

	db := database.DB

	var player models.Player
	// eager load
	if err := db.Preload("Bank").Preload("Wallet").First(&player, playerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request or Failed to retrieved player", "details": err.Error()})
		}
		return

	}
	c.JSON(http.StatusOK, gin.H{"player": player})
}
