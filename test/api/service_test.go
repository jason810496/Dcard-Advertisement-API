package api_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetAdFromDB(t *testing.T) {
	SetupFunction(t, ClearDB, GenerateAds)
	defer TeardownFunction(t)

	tests := []struct {
		req     *schemas.PublicAdRequest
		wantLen int
		wantErr error
	}{
		{
			req: &schemas.PublicAdRequest{
				Age:      &[]int{18}[0],
				Country:  "TW",
				Platform: "ios",
			},
			wantLen: 3,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age: &[]int{18}[0],
			},
			wantLen: 72,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Country: "TW",
			},
			wantLen: 12,
			wantErr: nil,
		},
		{
			req:     &schemas.PublicAdRequest{},
			wantLen: 72,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age: &[]int{90}[0],
			},
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(cache.PublicAdKey(tt.req), func(t *testing.T) {
			SetupFunction(t)
			defer TeardownFunction(t)

			srv := services.NewPublicService()
			got := make([]models.Advertisement, 0)
			err := srv.GetAdFromDB(tt.req, &got)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantLen, len(got))
		})
	}
}

type ServiceSetAdFromRedisCase struct {
	req     *schemas.PublicAdRequest
	got     []models.Advertisement
	want    []models.Advertisement
	wantErr error
}

func TestServiceSetAdAndGetAdFromRedis(t *testing.T) {
	// ClearRedis in every test case setup
	SetupFunction(t, ClearDB, GenerateAds)
	srv := services.NewPublicService()

	tests := []ServiceSetAdFromRedisCase{
		{
			req: &schemas.PublicAdRequest{
				Age:      &[]int{18}[0],
				Country:  "TW",
				Platform: "ios",
				Offset:   &[]int{0}[0],
				Limit:    &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age:    &[]int{18}[0],
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Country: "TW",
				Offset:  &[]int{0}[0],
				Limit:   &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age:    &[]int{90}[0],
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(cache.PublicAdKey(tt.req), func(t *testing.T) {
			SetupFunction(t, ClearRedis)
			err := srv.GetAdFromDB(tt.req, &tt.want)
			if err != nil {
				fmt.Println("GetAdFromDB got err")
				utils.PrintJson(err)
			}
			err = srv.SetAdToRedis(tt.req, &tt.want)
			if err != nil {
				fmt.Println("SetAdToRedis got err")
				utils.PrintJson(err)
			}
			defer TeardownFunction(t)

			time.Sleep(100 * time.Millisecond)
			err, found := srv.GetAdFromRedis(tt.req, &tt.got)
			if err != nil {
				fmt.Println("got err")
				utils.PrintJson(err)
				tt.wantErr = err
			}
			if !found {
				fmt.Println("not found")
			}

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, len(tt.want), len(tt.got))

			// compare in raw json string
			wantJson, _ := json.Marshal(tt.want)
			gotJson, _ := json.Marshal(tt.got)
			assert.Equal(t, string(wantJson), string(gotJson))
		})
	}
}

func TestServiceSetAdAndGetAdFromLocalCache(t *testing.T) {
	SetupFunction(t, ClearDB, GenerateAds)
	srv := services.NewPublicService()

	tests := []ServiceSetAdFromRedisCase{
		{
			req: &schemas.PublicAdRequest{
				Age:      &[]int{18}[0],
				Country:  "TW",
				Platform: "ios",
				Offset:   &[]int{0}[0],
				Limit:    &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age:    &[]int{18}[0],
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Country: "TW",
				Offset:  &[]int{0}[0],
				Limit:   &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age:    &[]int{90}[0],
				Offset: &[]int{0}[0],
				Limit:  &[]int{100}[0],
			},
			got:     make([]models.Advertisement, 0),
			want:    make([]models.Advertisement, 0),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(cache.PublicAdKey(tt.req), func(t *testing.T) {
			SetupFunction(t)
			err := srv.GetAdFromDB(tt.req, &tt.want)
			if err != nil {
				fmt.Println("GetAdFromDB got err")
				utils.PrintJson(err)
			}
			err = srv.SetAdToLocal(tt.req, &tt.want)
			if err != nil {
				fmt.Println("SetAdToLocal got err")
				utils.PrintJson(err)
			}
			defer TeardownFunction(t)

			time.Sleep(10 * time.Millisecond)
			err, found := srv.GetAdFromLocal(tt.req, &tt.got)
			if err != nil {
				fmt.Println("got err")
				utils.PrintJson(err)
				tt.wantErr = err
			}
			if !found {
				fmt.Println("not found")
			}

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, len(tt.want), len(tt.got))

			// compare in raw json string
			wantJson, _ := json.Marshal(tt.want)
			gotJson, _ := json.Marshal(tt.got)
			assert.Equal(t, string(wantJson), string(gotJson))
		})
	}
}
