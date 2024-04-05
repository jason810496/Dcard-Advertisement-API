package main

import (
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/scheduler"

	"github.com/jason810496/Dcard-Advertisement-API/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.Init()
	database.Init()
	cache.Init()
	cache.RedisClientInstance.CheckConnection()
	cache.LocalCacheInit()
	scheduler.PreHeatLocalCache(cache.LocalCacheInstance)

	router := handlers.SetupRouter()
	docs.SwaggerInfo.Title = "Dcard Advertisement API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "This is a swagger for Dcard Backend Intern 2024"
	docs.SwaggerInfo.Host = config.Settings.App.Address
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	database.CheckConnection()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(config.Settings.App.Address)

}
