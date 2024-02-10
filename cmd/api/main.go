package main

import (
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	// "net/http"

	// "github.com/gin-gonic/gin"
	docs "github.com/jason810496/Dcard-Advertisement-API/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	docs.SwaggerInfo.Title = "Dcard Advertisement API"
	docs.SwaggerInfo.Description = "This is a swagger for Dcard Backend Intern 2024"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "0.0.0.0:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := handlers.SetupRouter()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(":8080")

}
