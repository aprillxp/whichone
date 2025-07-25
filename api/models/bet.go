package models

import (
	"time"

	"gorm.io/gorm"
)

type Bet struct {
	gorm.Model
	PlayerID    uint   `json:"player_id"`
	Player      Player `gorm:"foreignKey:PlayerID"`
	BetAmount   int64  `json:"bet_amount"`
	Result      string `json:"result"`
	Payout      int64  `json:"payout"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
