package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	PlayerID    uint   `json:"player_id"`
	Player      Player `gorm:"foreignKey:PlayerID"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
