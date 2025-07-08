package models

import "time"

type Wallet struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"` // unique = 1 user 1 wallet
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
