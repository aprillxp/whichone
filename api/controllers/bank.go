package controllers

import (
	"api/database"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterBankInput struct {
	NamaRekening  string `json:"nama_rekening" binding:"required"`
	NomorRekening string `json:"nomor_rekening" binding:"required"`
	NamaBank      string `json:"nama_bank" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

func RegisterBank(c *gin.Context) {
	var input RegisterBankInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	playerID, exist := c.Get("player_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Player ID not found"})
		return
	}
	pID := playerID.(uint) // convert player id to interger

	var existingBank models.Bank
	result := database.DB.Where("player_id = ? AND nomor_rekening = ?", pID, input.NomorRekening).First(&existingBank)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			bank := models.Bank{
				PlayerID:      pID,
				NamaRekening:  input.NamaRekening,
				NomorRekening: input.NomorRekening,
				NamaBank:      input.NamaBank,
				IsDefault:     input.IsDefault,
			}
			if err := database.DB.Create(&bank).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register this account number", "details": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "Bank account registered successfully", "bank": bank})

			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Bank account is already registered for this account"})
		return
	}
}
