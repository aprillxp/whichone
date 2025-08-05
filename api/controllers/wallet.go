package controllers

import (
	"api/database"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TopUpInput struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
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

type WithdrawRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

func Withdraw(c *gin.Context) {
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerID := c.MustGet("player_id").(uint)

	transaction := database.DB.Begin()
	if transaction.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	var player models.Player
	if err := transaction.Preload("Wallet").First(&player, playerID).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Player not found"})
		return
	}

	if player.Wallet.Balance < req.Amount {
		transaction.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance for withdrawal"})
		return
	}

	if player.Bank.ID == 0 {
		transaction.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bank is not found or not registered yet."})
		return
	}

	player.Wallet.Balance -= req.Amount
	if err := transaction.Save(&player.Wallet).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	withdrawalTx := models.Transaction{
		PlayerID:    playerID,
		Type:        "withdraw",
		Amount:      -req.Amount,
		Description: "Withdrawal request..",
		Status:      "pending",
	}
	if err := transaction.Create(&withdrawalTx).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch withdrawal transaction"})
		return
	}

	if err := transaction.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Withdrawal request submitted successfully",
		"amount":             req.Amount,
		"new_balance":        player.Wallet.Balance,
		"transaction_status": "pending",
	})
}

func GetPlayerTransaction(c *gin.Context) {
	playerID := c.MustGet("player_id").(uint)

	var transaction []models.Transaction
	if err := database.DB.Where("player_id = ?", playerID).Order("created_at desc").Find(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
	})
}
