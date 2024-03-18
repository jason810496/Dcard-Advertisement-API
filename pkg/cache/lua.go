package cache

import (
	"github.com/redis/go-redis/v9"
)

// update ZSET atomically
// 1. DEL key
// 2. ZADD key score member [score member ...]
// EVAL "redis.call('DEL', KEYS[1]) redis.call('ZADD', KEYS[1], 'NX',unpack(ARGV))" 1 test 1 a 2 b 3 d
var UpdateCacheScript = redis.NewScript(`
	if redis.call("EXISTS", KEYS[1]) == 1 then
		redis.call("DEL", KEYS[1])
	end
	redis.call("ZADD", KEYS[1], unpack(ARGV))
`)
