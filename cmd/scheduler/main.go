package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/scheduler"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
)

func main() {
	Init()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// renewCache(&config.Settings.Schedule.FirstInterval)
	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			config.Settings.Schedule.Period,
		),
		gocron.NewTask(
			renewCache,
			&config.Settings.Schedule.Interval,
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
	case <-ctx.Done():
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	fmt.Println("scheduler is shut down")
}

func Init() {
	config.Init()
	database.Init()
	// cache.Init()
	cache.InitClusterReadClient()
	cache.InitClusterWriteClient()
	cache.RedisFailoverClusterClientReadInstance.CheckConnection()
	cache.RedisFailoverClusterClientWriteInstance.CheckConnection()
}

func renewCache(interval *time.Duration) {
	fmt.Println("start renew cache")
	srv := services.NewPublicService()
	schema := schemas.PublicAdRequest{}
	ads := []models.Advertisement{}

	// loop through gender, country, platform, age
	scheduler.GenerateHotData(interval, func(g, c, p *string, a *int) {
		schema.Country = *c
		schema.Gender = *g
		schema.Platform = *p
		if *a == config.Settings.HotData.AgeEnd+1 {
			schema.Age = nil
		} else {
			schema.Age = a
		}

		key := cache.PublicAdKey(&schema)
		fmt.Println("key: ", key)

		err := srv.GetAdFromDB(&schema, &ads)
		if err != nil {
			fmt.Println("GetAdFromDB error")
			fmt.Println(err)
		}

		err_rds := srv.SetHotSpotAdToRedis(&schema, &ads)
		if err_rds != nil {
			fmt.Println("SetHotSpotAdToRedis error")
			fmt.Println(err_rds)
		}
	})

	fmt.Println("finish renew cache")
}
