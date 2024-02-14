package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

// CreateAdAPI godoc
//
//	@Summary		admin API
//	@tags			ad
//	@Accept			json
//	@Produce		json
//	@Description	create advertisement with `startAt`, `endAt` and `condition`
//	@Param			advertisement	body	schemas.CreateAdRequest	true	"advertisement request schema"
//	@Produce		json
//	@Success		200	{object}	schemas.CreateAdResponse
//	@Failure		400	{object}	utils.HTTPError
//	@Router			/api/v1/ad [post]
func CreateAd(ctx *gin.Context) {
	var json schemas.CreateAdRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		utils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v\n", json)

	ctx.JSON(200, gin.H{
		"message": "create ad",
	})
}
