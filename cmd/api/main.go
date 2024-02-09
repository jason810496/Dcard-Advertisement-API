package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/jason810496/Dcard-Advertisement-API/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
)
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
   g.JSON(http.StatusOK,"helloworld")
}

func main()  {
   router := gin.Default()
   docs.SwaggerInfo.BasePath = "/api/v1"
   v1_router := router.Group("/api/v1")
   {
      hello_example := v1_router.Group("/example")
      {
         hello_example.GET("/helloworld",Helloworld)
      }
   }
   router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
   
   router.Run(":8080")


   
}