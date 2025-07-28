package controllers

import (
	"api/database"
	"api/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BetRequest struct {
	BetAmount int64 `json:"bet_amount" binding:"required,min=1"`
}

func PlaceBet(c *gin.Context) {
	var req BetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerID := c.MustGet("player_id").(uint)

	var player models.Player
	if err := database.DB.Preload("Wallet").First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Player not found"})
		return
	}

	// check balance amount
	if player.Wallet.Balance < req.BetAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance."})
		return
	}

	// simulate win/lose
	rand.Seed(time.Now().UnixNano())
	isWin := rand.Intn(2) == 1

	var result string
	var payout int64
	var description string

	transaction := database.DB.Begin()
	if transaction.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	// deduct bet amount
	player.Wallet.Balance -= req.BetAmount
	if err := transaction.Save(&player.Wallet).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance."})
		return
	}

	// record bet tx
	betDeductionTx := models.Transaction{
		PlayerID:    playerID,
		Type:        "bet",
		Amount:      -req.BetAmount,
		Description: "Bet placed",
		Status:      "completed",
	}
	if err := transaction.Create(&betDeductionTx).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record bet deduction transaction."})
		return
	}

	if isWin {
		result = "win"
		payout = req.BetAmount * 2
		player.Wallet.Balance += payout
		description = "Won bet."

		winTransaction := models.Transaction{
			PlayerID:    playerID,
			Type:        "win",
			Amount:      payout,
			Description: "Bet win payout",
			Status:      "completed",
		}
		if err := transaction.Create(&winTransaction).Error; err != nil {
			transaction.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
			return
		}
	} else {
		result = "lose"
		payout = 0
		description = "Lost bet"
	}

	// update wallet after win/lose
	if err := transaction.Save(&player.Wallet).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	betHistory := models.Bet{
		PlayerID:    playerID,
		BetAmount:   req.BetAmount,
		Result:      result,
		Payout:      payout,
		Description: description,
	}
	if err := transaction.Create(&betHistory).Error; err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update history"})
		return
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Bet placed successfully",
		"bet_id":      betHistory.ID,
		"bet_amount":  req.BetAmount,
		"result":      result,
		"payout":      payout,
		"new_balance": player.Wallet.Balance,
	})
}
