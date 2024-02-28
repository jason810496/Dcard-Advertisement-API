package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	"github.com/stretchr/testify/assert"
)

func setup() {
	// Do something here.
	os.Args = []string{"", "-config", "test"}
	gin.SetMode(gin.TestMode)

	config.Init()
	database.Init()
	database.CheckConnection()

	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	// Do something here.
	database.CloseConnection()
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestPublicInvalidAgeLower(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Age","message":"Should be greater than 1"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"age": "0",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)

}

func TestPublicInvalidAgeUpper(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Age","message":"Should be less than 100"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"age": "101",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)

}

func TestPublicInvalidCountry(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Country","message":"Should be one of TW HK JP US KR"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"country": "UK",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)

}

func TestPublicInvalidPlatform(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Platform","message":"Should be one of ios android web"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"platform": "windows",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)

}

func TestPublicInvalidGender(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Gender","message":"Should be one of F M"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"gender": "X",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)
}

func TestPublicInvalidLimit(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"limit": "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)
}

func TestPublicInvalidOffset(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Offset","message":"Should be greater than 0"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"offset": "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)
}

func TestPublicInvalidCountryAndPlatform(t *testing.T) {
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Country","message":"Should be one of TW HK JP US KR"},{"field":"Platform","message":"Should be one of ios android web"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"country":  "UK",
			"platform": "windows",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)
}

func TestPublicValidInvalidOffsetAndLimit(t *testing.T){
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"},{"field":"Offset","message":"Should be greater than 0"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"limit":  "-1",
			"offset": "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		},
		)
}

func TestPublicValidInvalidMultipleCondition(t *testing.T){
	r := gofight.New()
	errorJson := `{"code":400,"errors":[{"field":"Age","message":"Should be less than 100"},{"field":"Country","message":"Should be one of TW HK JP US KR"},{"field":"Platform","message":"Should be one of ios android web"},{"field":"Offset","message":"Should be greater than 0"}]}`
	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"age":     "101",
			"country": "UK",
			"platform": "windows",
			"offset":  "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, errorJson, r.Body.String())
		})
}
