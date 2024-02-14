package services

import (
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService() *AdminService {
	return &AdminService{db: database.DB}
}

func (srv *AdminService) CreateAdvertisement(adv *schemas.CreateAdRequest) error {
	advertisement := models.Advertisement{
		Title:    adv.Title,
		StartAt:  adv.StartAt,
		EndAt:    adv.EndAt,
		AgeStart: uint8(adv.Conditions.AgeStart), // convert int to uint8 (age is between 0-127
		AgeEnd:   uint8(adv.Conditions.AgeEnd),
		Country:  adv.Conditions.Country[0],  // only store the first country
		Platform: adv.Conditions.Platform[0], // only store the first platform
	}

	if err := srv.db.Create(&advertisement).Error; err != nil {
		return err
	}

	return nil
}
