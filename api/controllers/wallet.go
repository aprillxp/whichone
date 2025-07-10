package controllers

import (
	"api/database"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TopUpInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func TopUpBalance(c *gin.Context) {
	var input TopUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid top up amount"})
	}

	playerID, exist := c.Get("player_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, "Player ID not found in this context")
		return
	}
	pID := playerID.(uint)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Where("player_id = ?", pID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found for this player"})
				return err
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error retreiving wallet", "details": err.Error()})
			return err
		}

		wallet.Balance += int64(input.Amount)

		if err := tx.Save(&wallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up wallet", "details": err.Error()})
			return err
		}
		return nil
	})

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet top up successfully"})
}
