package services

import (
	"fmt"
	// "log"
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PublicService struct {
	db *gorm.DB
}

func NewPublicService() *PublicService {
	return &PublicService{db: database.DB}
}

func (srv *PublicService) GetAdvertisements(req schemas.PublicAdRequest) ([]models.Advertisement, error) {
	var ads []models.Advertisement
	// StartAt <= now <= EndAt
	now := time.Now()
	fmt.Println(now)

	tx := srv.db
	tx = tx.Select("title", "end_at")
	tx = tx.Where("start_at <= ? AND end_at >= ?", now, now)
	// has Gender condition
	if req.Gender != "" {
		tx = tx.Where(datatypes.JSONQuery("gender").HasKey(req.Gender))
	}

	// has Country condition
	if req.Country != "" {
		tx = tx.Where(datatypes.JSONQuery("country").HasKey(req.Country))
	}

	// has Platform condition
	if req.Platform != "" {
		tx = tx.Where(datatypes.JSONQuery("platform").HasKey(req.Platform))
	}

	// has Age condition
	if req.Age != 0 {
		tx = tx.Where("age_start <= ? AND age_end >= ?", req.Age, req.Age)
	}

	// limit and offset
	// tx = tx.Limit(req.Limit).Offset(req.Offset)

	if err := tx.Find(&ads).Error; err != nil {
		return nil, err
	}

	return ads, nil
}
