package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// init all component
	Init()
	// setup Router
	router := InitRouter()
	apiSrv := &http.Server{
		Addr:    config.Settings.App.Address,
		Handler: router,
	}
	//  API goroutine
	go func() {
		if err := apiSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// renew localCache goroutine
	go func() {
		for {
			scheduler.RefreshLocalCache(cache.LocalCacheInstance, config.Settings.LocalCache.Interval)
			time.Sleep(config.Settings.LocalCache.Period)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func Init() {
	config.Init()
	database.Init()
	database.CheckConnection()
	cache.Init()
	cache.RedisClientInstance.CheckConnection()
	cache.LocalCacheInit()
	scheduler.RefreshLocalCache(cache.LocalCacheInstance)
}

func InitRouter() *gin.Engine {
	router := handlers.SetupRouter()
	docs.SwaggerInfo.Title = "Dcard Advertisement API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "This is a swagger for Dcard Backend Intern 2024"
	docs.SwaggerInfo.Host = config.Settings.App.Address
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
