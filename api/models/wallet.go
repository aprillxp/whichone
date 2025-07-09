package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID        uint  `gorm:"primaryKey"`
	PlayerID  uint  `gorm:"unique;not null"` // unique = 1 user 1 wallet
	Balance   int64 `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
