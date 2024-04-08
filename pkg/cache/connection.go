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

type RedisClusterClient struct {
	Client  *redis.ClusterClient
	Context context.Context
}

var RedisClientInstance *RedisClient
var RedisFailoverClusterClientReadInstance *RedisClusterClient
var RedisFailoverClusterClientWriteInstance *RedisClient

// https://github.com/redis/go-redis/issues/1169

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
}

func CloseConnection() {
	RedisClientInstance.Client.Close()
}

func InitClusterReadClient() {
	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:       config.Settings.Redis.Sentinel.MasterName,
		SentinelAddrs:    config.Settings.Redis.Sentinel.Addrs,
		Password:         config.Settings.Redis.Password,
		SentinelPassword: config.Settings.Redis.Password,
		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		RouteRandomly: true,
	})

	RedisFailoverClusterClientReadInstance = &RedisClusterClient{
		Client:  rdb,
		Context: context.Background(),
	}
}

func InitClusterWriteClient() {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       config.Settings.Redis.Sentinel.MasterName,
		SentinelAddrs:    config.Settings.Redis.Sentinel.Addrs,
		Password:         config.Settings.Redis.Password,
		SentinelPassword: config.Settings.Redis.Password,
	})

	RedisFailoverClusterClientWriteInstance = &RedisClient{
		Client:  rdb,
		Context: context.Background(),
	}
}

func (r *RedisClusterClient) CheckConnection() {
	val, err := r.Client.Ping(r.Context).Result()
	if err != nil {
		fmt.Println("Redis connection failed")
		fmt.Println(err)
	}
	fmt.Println("Redis connection success")
	fmt.Print(val)
}

func GetClusterWriteClient() *redis.Client {
	return RedisFailoverClusterClientWriteInstance.Client
}

func GetClusterReadClient() *redis.ClusterClient {
	return RedisFailoverClusterClientReadInstance.Client
}

func CloseClusterReadClient() {
	RedisFailoverClusterClientReadInstance.Client.Close()
}

func CloseClusterWriteClient() {
	RedisFailoverClusterClientWriteInstance.Client.Close()
}
