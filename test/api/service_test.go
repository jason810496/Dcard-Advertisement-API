package api_test

import (
	"testing"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestGetAdFromDB(t *testing.T) {
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
				Age: 	&[]int{18}[0],
			},
			wantLen: 30,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{},
			wantLen: 30,
			wantErr: nil,
		},
		{
			req: &schemas.PublicAdRequest{
				Age:      &[]int{90}[0],
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