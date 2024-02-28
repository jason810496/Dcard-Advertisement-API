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

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"age": "101",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)

}

func TestPublicInvalidCountry(t *testing.T) {
	r := gofight.New()

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"country": "UK",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)

}

func TestPublicInvalidPlatform(t *testing.T) {
	r := gofight.New()

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"platform": "windows",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)

}

func TestPublicInvalidGender(t *testing.T) {
	r := gofight.New()

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"gender": "X",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)
}

func TestPublicInvalidLimit(t *testing.T) {
	r := gofight.New()

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"limit": "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)
}

func TestPublicInvalidOffset(t *testing.T) {
	r := gofight.New()

	r.GET("/api/v1/ad").
		SetQuery(gofight.H{
			"offset": "-1",
		}).
		Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		},
		)
}
