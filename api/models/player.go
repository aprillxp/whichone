package models

import "time"

type Player struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Wallet Wallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // one to one
}
