package models

import (
	"gorm.io/datatypes"
	// "gorm.io/gorm"
	"time"
)

// Advertisement model
// `datatypes.JSON` is a custom type for gorm to store JSON data
// https://github.com/go-gorm/datatypes
type Advertisement struct {
	// gorm.Model
	ID        uint       `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Title   string     `json:"title,omitempty"`
	StartAt *time.Time `json:"startAt,omitempty"`
	EndAt   *time.Time `json:"endAt,omitempty" time_format:"RFC3339"`
	// conditions
	AgeStart uint8          `json:"ageStart,omitempty"`
	AgeEnd   uint8          `json:"ageEnd,omitempty"`
	Gender   datatypes.JSON `json:"gender,omitempty"`   // ["F", "M"]
	Country  datatypes.JSON `json:"country,omitempty"`  // ["TW", "HK", "JP", "US", "KR"]
	Platform datatypes.JSON `json:"platform,omitempty"` // ["ios", "android", "web"]
	// have loaded to redis
	Active bool `json:"active,omitempty"`
}

// func (ad *Advertisement) AfterFind(tx *gorm.DB) (err error) {
// 	if ad.EndAt != nil {
// 		*ad.EndAt = ad.EndAt.In(time.Local)
// 	}
// 	return
//   }
