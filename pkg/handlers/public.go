package handlers

import (
	"github.com/gin-gonic/gin"
	schemas "github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
)

// PublicAdAPI godoc
// @Summary public API
// @tags ad
// @Produce json
// @Description query ad by query parameters age, country
// @Param age query int false "age" minimum(1) maximum(100)
// @Param country query string false "country" Enums(TW, HK, JP, US, KR)
// @Param platform query string false "platform" Enums(ios, android,web)
// @Param gender query string false "gender" Enums(F,M)
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} schemas.PublicAdResponse
// @Failure 400 {object} utils.HTTPError
// @Router /api/v1/ad [get]
func PublicAd(c *gin.Context) {

	items := schemas.PublicAdResponse{
		Items: []schemas.PublicAdItem{
			{
				Title: "title",
				EndAt: "endAt",
			},
		},
	}

	c.JSON(200, items)
}
