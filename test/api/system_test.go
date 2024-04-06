package api_test

/*
import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	"github.com/stretchr/testify/assert"
)

func TestGetColdNonExistAd(t *testing.T) {
	PrintTestTitle(t.Name())
	ClearDB()
	ClearRedis()

	r := gofight.New()
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"age":      "20",
			"country":  "TW",
			"platform": "ios",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "[]", r.Body.String())
		})
}

func TestGetColdExistAd(t *testing.T) {
	PrintTestTitle(t.Name())
	ClearDB()
	ClearRedis()

	// add ad
}

*/
