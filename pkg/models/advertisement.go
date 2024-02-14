package models

import (
	"gorm.io/gorm"
	"time"
)

// Advertisement model
type Advertisement struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Title   string
	StartAt time.Time
	EndAt   time.Time
	// conditions
	AgeStart uint8
	AgeEnd   uint8
	Country  string `gorm:"type:varchar(2)"`  // ISO 3166-1 alpha-2
	Platform string `gorm:"type:varchar(10)"` // android, ios, web
}
