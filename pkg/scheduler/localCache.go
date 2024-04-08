package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

func RefreshLocalCache(lc *fastcache.Cache, intervalArg ...time.Duration) {
	fmt.Println("Start Refresh local cache")
	// Preheat cache
	srv := services.NewPublicService()	
	req := schemas.PublicAdRequest{}
	req.Offset = new(int)
	*req.Offset = 0
	req.Limit = new(int)
	*req.Limit = 10000

	interval := time.Millisecond
	// use custom interval if pass args
	if len(intervalArg) > 0 {
		interval = intervalArg[0]
	}
	GenerateHotData(&interval, func(g *string, c *string, p *string, a *int) {
		req.Country = *c
		req.Gender = *g
		req.Platform = *p
		if *a == config.Settings.HotData.AgeEnd+1 {
			req.Age = nil
		} else {
			req.Age = a
		}

		key := cache.PublicAdKey(&req)
		var ads []models.Advertisement
		err, found := srv.GetAdFromRedis(&req, &ads)
		if err != nil {
			utils.PrintJson(err)
		}
		if !found {
			err = srv.GetAdFromDB(&req, &ads)
		}
		err = srv.GetAdFromDB(&req, &ads)
		
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
