package services

import (
	"log"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
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
		StartAt:  &adv.StartAt,
		EndAt:    &adv.EndAt,
		AgeStart: uint8(adv.Conditions.AgeStart), // convert int to uint8 (age is between 0-127
		AgeEnd:   uint8(adv.Conditions.AgeEnd),
		Gender:   utils.ToJsonArray(adv.Conditions.Gender),
		Country:  utils.ToJsonArray(adv.Conditions.Country),
		Platform: utils.ToJsonArray(adv.Conditions.Platform),
	}

	if err := srv.db.Create(&advertisement).Error; err != nil {
		return err
	}

	log.Printf("%#v\n", advertisement)

	return nil
}
