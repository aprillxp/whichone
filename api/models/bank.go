package models

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	ID            uint `grom:"primaryKey"`
	PlayerID      uint
	NamaRekening  string `gorm:"size:255;not null"`
	NomorRekening string `gorm:"size:255;unique;not null"` // No rekening harus unik
	NamaBank      string `gorm:"size:255;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
