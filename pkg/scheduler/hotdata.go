package scheduler

import (
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

type callbackType func(*string, *string, *string, *int)

func GenerateHotData(interval *time.Duration, callback callbackType) {
	for _, g := range utils.GenderList {
		// country
		for _, c := range utils.CountryList {
			// platform
			for _, p := range utils.PlatformList {
				// age
				for a := config.Settings.HotData.AgeStart; a <= config.Settings.HotData.AgeEnd; a++ {
					callback(&g, &c, &p, &a)
					time.Sleep(*interval)
				}
			}
		}
	}
}
