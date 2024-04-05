package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/middleware"
)

func SetupRouter() *gin.Engine {
	configModeToGinMode := map[string]string{
		"prod": gin.ReleaseMode,
		"dev":  gin.DebugMode,
		"test": gin.TestMode,
	}
	gin.SetMode(configModeToGinMode[config.Settings.App.Env])

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	v1_router := router.Group("/api/v1")
	{
		ad_router := v1_router.Group("/ad")
		{
			ad_router.GET("", PublicAd)
			ad_router.POST("", CreateAd)
		}
	}

	return router
}
