package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

func main() {
	Init()
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(
			renewCache,
		),
	)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Second*6):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}


func renewCache() {
	fmt.Println("start renew cache")
	srv := services.NewPublicService()
	schema := schemas.NewPublicAdRequest()
	ads := []models.Advertisement{}
	// loop through gender, country, platform, age
	// gender
	for _,g := range utils.GenderList {
		// country
		for _,c := range utils.CountryList {
			// platform
			for _,p := range utils.PlatformList {
				// age
				for a:=18; a<=40; a++ {
					schema.Gender = g
					schema.Country = c
					schema.Platform = p
					schema.Age = a

					fmt.Println( fmt.Sprintln("ad:%s:%s:%s:%d",g,c,p,a))

					err := srv.GetAdFromDB(&schema, &ads)
					if err != nil {
						fmt.Println("GetAdFromDB error")
						fmt.Println(err)
					}

					err_rds := srv.SetHotSpotAdToRedis(&schema,&ads)
					if err_rds != nil {
						fmt.Println("SetHotSpotAdToRedis error")
						fmt.Println(err_rds)
					}

					time.Sleep( time.Microsecond * 100 )
				}
			}
		}
	}

	// test for once
	// schema.Gender = "F"
	// schema.Age = 18
	// schema.Country = "TW"
	// schema.Platform = "ios"

	// err := srv.GetAdFromDB(&schema, &ads)
	// if err != nil {
	// 	fmt.Println("GetAdFromDB error")
	// 	fmt.Println(err)
	// }

	// err_rds := srv.SetHotSpotAdToRedis(&schema,&ads)
	// if err_rds != nil {
	// 	fmt.Println("SetHotSpotAdToRedis error")
	// 	fmt.Println(err_rds)
	// }
	fmt.Println("finish renew cache")
}

func Init(){
	config.Init()
	database.Init()
	cache.Init()
	cache.Rds.CheckConnection()
}