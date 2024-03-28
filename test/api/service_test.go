package api_test

import (
	"encoding/json"
	"fmt"
	"testing"

	// "time"

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
			wantLen: 2,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age: &[]int{18}[0],
			},
			wantLen: 30,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Country: "TW",
			},
			wantLen: 6,
			wantErr: nil,
		},
		{
			req:     &schemas.PublicAdRequest{},
			wantLen: 30,
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
	req *schemas.PublicAdRequest
	got []models.Advertisement
	want []models.Advertisement
	wantErr error
}

func TestServiceSetAdFromRedis(t *testing.T){
	SetupFunction(t, ClearDB, GenerateAds)
	srv := services.NewPublicService()

	tests := []ServiceSetAdFromRedisCase{
		{
			req: &schemas.PublicAdRequest{
				Age:      &[]int{18}[0],
				Country:  "TW",
				Platform: "ios",
				Offset:  0,
				Limit:  100,
			},
			got: make([]models.Advertisement, 0),
			want: make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age: &[]int{18}[0],
				Offset:  0,
				Limit:  100,
			},
			got: make([]models.Advertisement, 0),
			want: make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Country: "TW",
				Offset:  0,
				Limit:  100,
			},
			got: make([]models.Advertisement, 0),
			want: make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req:     &schemas.PublicAdRequest{
				Offset:  0,
				Limit:  100,
			},
			got: make([]models.Advertisement, 0),
			want: make([]models.Advertisement, 0),
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age: &[]int{90}[0],
				Offset:  0,
				Limit:  100,
			},
			got: make([]models.Advertisement, 0),
			want: make([]models.Advertisement, 0),
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
			err = srv.SetAdToRedis(tt.req, &tt.want)
			if err != nil {
				fmt.Println("SetAdToRedis got err")
				utils.PrintJson(err)
			}
			defer TeardownFunction(t)

			err = srv.GetAdFromRedis(tt.req, &tt.want)
			if err != nil {
				fmt.Println("got err")
				utils.PrintJson(err)
				tt.wantErr = err
			}
			fmt.Println("want printJson")
			utils.PrintJson(tt.want)

			fmt.Println("got printJson")
			utils.PrintJson(tt.got)
			
			wantJson, _ := json.Marshal(tt.want)
			gotJson, _ := json.Marshal(tt.got)

			fmt.Println(("got:"))
			fmt.Println(string(gotJson))
			fmt.Println(("want:"))
			fmt.Println(string(wantJson))

			assert.Equal(t, tt.wantErr, err)
			// assert.Equal(t, string(wantJson), string(gotJson))
		})
	}

}


func TestServiceGetAdFromRedis(t *testing.T){

}