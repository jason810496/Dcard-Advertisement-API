package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/handlers"
	"github.com/stretchr/testify/assert"
)

func TestPublicAgeLowerBound(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		age       string
		code      int
		errorJson string
	}{
		{"0", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Age","message":"Should be greater than 1"}]}`},
		{"-1", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Age","message":"Should be greater than 1"}]}`},
		{"-55", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Age","message":"Should be greater than 1"}]}`},
		{"1", http.StatusOK, "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.age, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"age": tt.age,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicAgeUpperBound(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		age       string
		code      int
		errorJson string
	}{
		{"101", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Age","message":"Should be less than 100"}]}`},
		{"404", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Age","message":"Should be less than 100"}]}`},
		{"100", http.StatusOK, "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.age, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"age": tt.age,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicCountry(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		country   string
		code      int
		errorJson string
	}{
		{"TW", http.StatusOK, "[]"},
		{"HK", http.StatusOK, "[]"},
		{"JP", http.StatusOK, "[]"},
		{"US", http.StatusOK, "[]"},
		{"KR", http.StatusOK, "[]"},
		{"UK", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Country","message":"Should be one of TW HK JP US KR"}]}`},
		{"CN", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Country","message":"Should be one of TW HK JP US KR"}]}`},
		{"CA", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Country","message":"Should be one of TW HK JP US KR"}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"country": tt.country,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}

}

func TestPublicPlatform(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		platform  string
		code      int
		errorJson string
	}{
		{"ios", http.StatusOK, "[]"},
		{"android", http.StatusOK, "[]"},
		{"web", http.StatusOK, "[]"},
		{"windows", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Platform","message":"Should be one of ios android web"}]}`},
		{"mac", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Platform","message":"Should be one of ios android web"}]}`},
		{"linux", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Platform","message":"Should be one of ios android web"}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"platform": tt.platform,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicGender(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		gender    string
		code      int
		errorJson string
	}{
		{"F", http.StatusOK, "[]"},
		{"M", http.StatusOK, "[]"},
		{"X", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Gender","message":"Should be one of F M"}]}`},
		{"Female", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Gender","message":"Should be one of F M"}]}`},
		{"Male", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Gender","message":"Should be one of F M"}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.gender, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"gender": tt.gender,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicLimit(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		limit     string
		code      int
		errorJson string
	}{
		{"0", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"}]}`},
		{"-1", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"}]}`},
		{"-55", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"}]}`},
		{"1", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"}]}`},
		{"2", http.StatusOK, "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.limit, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"limit": tt.limit,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicOffset(t *testing.T) {

	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		offset    string
		code      int
		errorJson string
	}{
		{"0", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Offset","message":"Should be greater than 0"}]}`},
		{"-1", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Offset","message":"Should be greater than 0"}]}`},
		{"-55", http.StatusBadRequest, `{"code":400,"errors":[{"field":"Offset","message":"Should be greater than 0"}]}`},
		{"1", http.StatusOK, "[]"},
		{"2", http.StatusOK, "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.offset, func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(gofight.H{
					"offset": tt.offset,
				}).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}

func TestPublicMultipleCondition(t *testing.T) {
	SetupFunction(t, ClearDB, ClearRedis)
	defer TeardownFunction(t)

	tests := []struct {
		condition map[string]string
		code      int
		errorJson string
	}{
		{
			map[string]string{
				"age":      "101",
				"country":  "UK",
				"platform": "windows",
				"offset":   "-1"},
			http.StatusBadRequest,
			`{"code":400,"errors":[{"field":"Age","message":"Should be less than 100"},{"field":"Country","message":"Should be one of TW HK JP US KR"},{"field":"Platform","message":"Should be one of ios android web"},{"field":"Offset","message":"Should be greater than 0"}]}`,
		},
		{
			map[string]string{
				"age":      "50",
				"country":  "TW",
				"platform": "linux",
			},
			http.StatusBadRequest,
			`{"code":400,"errors":[{"field":"Platform","message":"Should be one of ios android web"}]}`,
		},
		{
			map[string]string{},
			http.StatusOK,
			"[]",
		},
		{
			map[string]string{
				"age":      "0",
				"offset":   "10",
				"limit":    "1",
			},
			http.StatusBadRequest,
			`{"code":400,"errors":[{"field":"Age","message":"Should be greater than 1"}]}`,
		},
		{
			map[string]string{
				"age":      "50",
				"offset":   "0",
				"limit":    "0",
			},
			http.StatusBadRequest,
			`{"code":400,"errors":[{"field":"Limit","message":"Should be greater than 1"},{"field":"Offset","message":"Should be greater than 0"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.condition), func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			r := gofight.New()
			r.GET("/api/v1/ad").
				SetQuery(tt.condition).
				Run(handlers.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					fmt.Println(r.Body.String())
					assert.Equal(t, tt.code, r.Code)
					assert.Equal(t, tt.errorJson, r.Body.String())
				},
				)
		})
	}
}
