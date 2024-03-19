package api_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

func SetupPackage() {
	os.Args = []string{"", "-config", "test"}
	gin.SetMode(gin.TestMode)

	config.Init()
	database.Init()
	database.CheckConnection()
	cache.Init()
	cache.Rds.CheckConnection()

	fmt.Printf("\033[1;33m%s\033[0m", "> Package Setup completed\n")
}

func TeardownPackage() {
	database.CloseConnection()
	fmt.Printf("\033[1;33m%s\033[0m", "> Package Teardown completed")
	fmt.Printf("\n")
}

func SetupFunction(t *testing.T, args ...func()) {
	// execute the function
	for _, f := range args {
		f()
	}

	// fmt.Printf("\033[1;33m%s\033[0m \033[4m%s\033[0m \033[1;33m%s\033[0m",">",t.Name(), " Function Setup completed\n")
	fmt.Printf("\033[94m%s %s\033[0m\n", ">", t.Name())

}

func TeardownFunction(t *testing.T, args ...func()) {
	// execute the function
	for _, f := range args {
		f()
	}
	// fmt.Printf("\033[1;33m%s\033[0m \033[4m%s\033[0m \033[1;33m%s\033[0m",">",t.Name(), " Function Teardown completed\n")
}

func PrintTestTitle(title string) {
	fmt.Printf("\033[1;34m> %s\033[0m", title)
	fmt.Printf("\n")
}

func ClearDB() {
	db := database.DB
	db.Exec("DELETE FROM advertisements")

}

func ClearRedis() {
	cache := cache.Rds
	rds := cache.Client
	rds_ctx := cache.Context
	rds.FlushAll(rds_ctx)
}

func GenerateAds() {
	db := database.DB
	ageStart := 18
	ageEnd := 30

	for _, country := range utils.CountryList {
		for _, platform := range utils.PlatformList {
			for _, gender := range utils.GenderList {
				countryList := []string{
					country,
				}
				platformList := []string{
					platform,
				}
				genderList := []string{
					gender,
				}

				ad := models.Advertisement{
					Title:    fmt.Sprintf("AD-%s-%s-%s-%d-%d", country, platform, gender, ageStart, ageEnd),
					StartAt:  &[]time.Time{time.Now()}[0],
					EndAt:    &[]time.Time{time.Now().Add(time.Hour * 24 * 30)}[0],
					AgeStart: uint8(ageStart),
					AgeEnd:   uint8(ageEnd),
					Gender:   utils.ToJsonArray(genderList),
					Country:  utils.ToJsonArray(countryList),
					Platform: utils.ToJsonArray(platformList),
				}
				db.Create(&ad)
			}
		}
	}
}
