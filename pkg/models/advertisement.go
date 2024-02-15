package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// Advertisement model
// `datatypes.JSON` is a custom type for gorm to store JSON data
// https://github.com/go-gorm/datatypes
type Advertisement struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Title   string
	StartAt time.Time
	EndAt   time.Time
	// conditions
	AgeStart uint8
	AgeEnd   uint8
	Country  datatypes.JSON // ["TW", "JP"]
	Platform datatypes.JSON // ["ios", "android"]
}
