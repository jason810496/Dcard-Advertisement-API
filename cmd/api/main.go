package main

import (
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	// "net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
)


// @title Dcard Advertisement API
// @version 1.0
// @description This is a swagger for Dcard Backend Intern 2024
// @host 0.0.0.0:8080
// @BasePath
// @schemes http, https
func main() {

	router := handlers.SetupRouter()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(config.Settings.App.Address)

}
