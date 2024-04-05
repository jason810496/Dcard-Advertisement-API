package handlers

import (

	"github.com/gin-gonic/gin"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

// PublicAdAPI godoc
//
//	@Summary		public API
//	@tags			ad
//	@Produce		json
//	@Description	query ad by query parameters age, country
//	@Param			age			query		int		false	"age"		minimum(1)	maximum(100)
//	@Param			country		query		string	false	"country"	Enums(TW, HK, JP, US, KR)
//	@Param			platform	query		string	false	"platform"	Enums(ios, android,web)
//	@Param			gender		query		string	false	"gender"	Enums(F,M)
//	@Param			limit		query		int		false	"limit"
//	@Param			offset		query		int		false	"offset"
//	@Success		200			{object}	schemas.PublicAdResponse
//	@Failure		400			{object}	utils.HTTPError
//	@Router			/api/v1/ad [get]
func PublicAd(ctx *gin.Context) {
	json := schemas.NewPublicAdRequest()

	if err := ctx.ShouldBindQuery(&json); err != nil {
		utils.NewError(ctx, 400, err)
		return
	}

	srv := services.GetPublicService()

	res, err := srv.GetAdvertisements(&json)
	if err != nil {
		utils.NewError(ctx, 400, err)
		return
	}

	ctx.JSON(200, res)
}
