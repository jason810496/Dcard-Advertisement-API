package handlers

import (
	"github.com/gin-gonic/gin"
	schemas "github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
)


// CreateAdAPI godoc
// @Summary admin API
// @tags ad
// @Accept json
// @Produce json
// @Description create advertisement with `startAt`, `endAt` and `condition`
// @Param advertisement body schemas.CreateAdRequest true "advertisement request schema"
// @Produce json
// @Success 200 {object} schemas.CreateAdResponse
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/ad [post]
func CreateAd(c *gin.Context) {

	var json schemas.CreateAdRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	c.JSON(200, gin.H{
		"message": "create ad",
	})
}
