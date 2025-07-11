package models

import (
	"gorm.io/gorm"
	"time"
)

type Player struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Banks     []Bank `gorm:"foreignKey:PlayerID"` // Punya banyak bank
	Wallet    Wallet `gorm:"foreignKey:PlayerID"` // Punya 1 dompet
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
