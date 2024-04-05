package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

var RedisClientInstance *RedisClient

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:        config.Settings.Redis.Host + ":" + config.Settings.Redis.Port,
		Password:    config.Settings.Redis.Password,
		PoolSize:    10,
		PoolTimeout: 5 * time.Second,
	})

	RedisClientInstance = &RedisClient{
		Client:  rdb,
		Context: context.Background(),
	}
}

func (r *RedisClient) CheckConnection() {
	val, err := r.Client.Ping(r.Context).Result()
	if err != nil {
		fmt.Println("Redis connection failed")
	}
	fmt.Println("Redis connection success")
	fmt.Print(val)
	// try set key
	err = r.Client.Set(r.Context, "key", "value", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	// get all keys
	keys, _, err := r.Client.Scan(r.Context, 0, "*", 100).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(keys)

}

func CloseConnection() {
	RedisClientInstance.Client.Close()
}
