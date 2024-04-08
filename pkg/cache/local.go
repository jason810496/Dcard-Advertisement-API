package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	ginCachePersist "github.com/chenyahui/gin-cache/persist"
)

var LocalCacheInstance *fastcache.Cache
var GinCachePersistInstance *ginCachePersist.MemoryStore

func LocalCacheInit() {
	LocalCacheInstance = fastcache.New(config.Settings.LocalCache.MaxSize)
	GinCachePersistInstance = ginCachePersist.NewMemoryStore(config.Settings.LocalCache.TTL)
}
