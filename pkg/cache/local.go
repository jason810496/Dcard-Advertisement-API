package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
)

var LocalCacheInstance *fastcache.Cache

func LocalCacheInit() {
	LocalCacheInstance = fastcache.New(config.Settings.LocalCache.MaxSize)
}
