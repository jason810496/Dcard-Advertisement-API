package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

func RefreshLocalCache(lc *fastcache.Cache, intervalArg ...time.Duration){
	fmt.Println("Start Refresh local cache")
	// Preheat cache
	srv := services.GetPublicService()
	req := schemas.NewPublicAdRequest()
	interval := time.Millisecond
	// use custom interval if pass args
	if len(intervalArg) > 0{
		interval = intervalArg[0]
	}
	GenerateHotData(&interval, func(g *string, c *string, p *string, a *int) {
		req.Age = a
		req.Country = *c
		req.Platform = *p
		req.Gender = *g

		key := cache.PublicAdKey(&req)
		var ads []models.Advertisement
		err := srv.GetAdFromDB(&req, &ads)

		if err != nil {
			utils.PrintJson(err)
		}

		value, err := json.Marshal(ads)
		if err != nil {
			utils.PrintJson(err)
		}

		lc.Set([]byte(key), []byte(value))
	})

	fmt.Println("Finish Refresh local cache")
}