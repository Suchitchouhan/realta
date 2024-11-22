package models

type Client struct {
	ID      uint    `gorm:"primaryKey"`
	Name    string  `gorm:"uniqueIndex"`
	Balance float64 `gorm:"not null"`
}
