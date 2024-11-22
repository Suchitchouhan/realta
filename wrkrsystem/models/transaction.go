package models

import (
	"time"
)

type Transaction struct {
	ID         uint    `gorm:"primaryKey"`
	SenderID   uint    `gorm:"index"`
	ReceiverID uint    `gorm:"index"`
	Amount     float64 `gorm:"not null"`
	Status     string  `gorm:"default:'pending'"`
	CreatedAt  time.Time
}
