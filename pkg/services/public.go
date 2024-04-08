package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"

	// "github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PublicService struct {
	db  *gorm.DB
	lc  *fastcache.Cache
	rds *cache.RedisClient
}

var PublicServiceInstance *PublicService

func GetPublicService() *PublicService {
	if PublicServiceInstance == nil {
		PublicServiceInstance = NewPublicService()
	}
	return PublicServiceInstance
}

func NewPublicService() *PublicService {
	return &PublicService{db: database.DB, rds: cache.RedisClientInstance, lc: cache.LocalCacheInstance}
}

// []models.Advertisement, error
func (srv *PublicService) GetAdvertisements(req *schemas.PublicAdRequest) (any, error) {
	var ads []models.Advertisement
	// get from local
	err, found := srv.GetAdFromLocal(req, &ads)
	if err == nil && found {
		return paginateAds(req, &ads), nil
	}

	// get from redis
	err, found = srv.GetAdFromRedis(req, &ads)
	if err == nil && found {
		// alredy use redis ZSET for pagination
		return ads, nil
	}
	go srv.SetAdToRedis(req, &ads)
	go srv.SetAdToLocal(req, &ads)

	err = srv.GetAdFromDB(req, &ads)
	if err != nil {
		return nil, err
	}

	// TTL: 5 minutes
	// store full result to redis in goroutine
	go srv.SetAdToRedis(req, &ads)
	go srv.SetAdToLocal(req, &ads)

	// if ads is empty, return [] instead of nil
	if len(ads) == 0 {
		ads = make([]models.Advertisement, 0)
		return ads, nil
	}

	return paginateAds(req, &ads), nil
}

func paginateAds(req *schemas.PublicAdRequest, ads *[]models.Advertisement) any {
	result_len := len(*ads)
	// case: offset >= result_len
	if req.Offset != nil && *req.Offset >= result_len {
		return nil
	}
	// case: offset+limit > result_len
	if *req.Offset+*req.Limit > result_len {
		return (*ads)[*req.Offset:]
	}

	// set offset and limit for sql result
	return (*ads)[*req.Offset : *req.Offset+*req.Limit]
}

func (srv *PublicService) GetAdFromLocal(req *schemas.PublicAdRequest, ads *[]models.Advertisement) (error, bool) {
	key := cache.PublicAdKey(req)
	val := srv.lc.Get(nil, []byte(key))
	if val == nil {
		return nil, false
	}
	if err := json.Unmarshal(val, ads); err != nil {
		return err, true
	}
	return nil, true
}

func (srv *PublicService) GetAdFromRedis(req *schemas.PublicAdRequest, ads *[]models.Advertisement) (error, bool) {
	key := cache.PublicAdKey(req)
	rds := srv.rds.Client
	rds_ctx := srv.rds.Context

	// redis check key exist
	if rds.Exists(rds_ctx, key).Val() == 0 {
		fmt.Println("key does not exist")
		return redis.Nil, false
	}

	val, err := rds.ZRangeByScore(rds_ctx, key, &redis.ZRangeBy{
		Min: strconv.Itoa(*req.Offset + 1),
		Max: strconv.Itoa(*req.Offset + *req.Limit),
	}).Result()
	if err != nil {
		return err, true
	}

	// bind to ads
	*ads = make([]models.Advertisement, len(val))
	for i, v := range val {
		if err := json.Unmarshal([]byte(v), &(*ads)[i]); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil, true
}

func (srv *PublicService) GetAdFromDB(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	// StartAt <= now <= EndAt
	now := time.Now()
	// fmt.Println(now)

	tx := srv.db
	tx = tx.Select("title", "end_at")
	tx = tx.Where("start_at <= ? AND end_at >= ?", now, now)
	// has Gender condition
	if req.Gender != "" {
		tx = tx.Where(datatypes.JSONQuery("gender").HasKey(req.Gender))
	}

	// has Country condition
	if req.Country != "" {
		tx = tx.Where(datatypes.JSONQuery("country").HasKey(req.Country))
	}

	// has Platform condition
	if req.Platform != "" {
		tx = tx.Where(datatypes.JSONQuery("platform").HasKey(req.Platform))
	}

	// has Age condition
	if req.Age != nil {
		tx = tx.Where("age_start <= ? AND age_end >= ?", *req.Age, *req.Age)
	}

	// limit and offset
	// tx = tx.Limit(req.Limit).Offset(req.Offset)

	if err := tx.Find(ads).Error; err != nil {
		return err
	}

	return nil
}

func (srv *PublicService) SetAdToLocal(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	key := cache.PublicAdKey(req)
	value, err := json.Marshal(ads)
	if err != nil {
		return err
	}

	srv.lc.Set([]byte(key), value)
	return nil
}

func (srv *PublicService) SetAdToRedis(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	key := cache.PublicAdKey(req)
	rds := srv.rds.Client
	rds_ctx := srv.rds.Context

	//no ads
	if len(*ads) == 0 {
		_, err := rds.ZAddArgs(rds_ctx, key, redis.ZAddArgs{
			XX: false,
			Members: []redis.Z{
				{Score: -999, Member: "empty"},
			},
		}).Result()
		if err != nil {
			fmt.Println("SetAdToRedis", err)
			return err
		}
		// set expire
		_, err = rds.Expire(rds_ctx, key, time.Minute*5).Result()
		if err != nil {
			return err
		}
		return nil
	}
	// use `redis.Do` is more efficient than `redis.ZAdd`
	cmd := make([]interface{}, 0, len(*ads)*2+2)
	cmd = append(cmd, "ZADD", key)
	for idx, ad := range *ads {
		cmd = append(cmd, float64(idx+1))
		cmd = append(cmd, fmt.Sprintf("{\"title\":\"%s\",\"endAt\":\"%s\"}", ad.Title, ad.EndAt.Format(utils.RFC3339Custom)))
	}
	_, err := rds.Do(rds_ctx, cmd...).Result()
	if err != nil {
		fmt.Println("SetAdToRedis", err)
		return err
	}
	// set expire
	_, err = rds.Expire(rds_ctx, key, time.Minute*5).Result()
	if err != nil {
		fmt.Println("SetAdToRedis", err)
		return err
	}

	return nil
}

func (srv *PublicService) SetHotSpotAdToRedis(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	key := cache.PublicAdKey(req)
	rds := srv.rds.Client
	rds_ctx := srv.rds.Context

	cmd := make([]interface{}, 0, len(*ads)*2)
	for idx, ad := range *ads {
		cmd = append(cmd, float64(idx+1))
		cmd = append(cmd, fmt.Sprintf("{\"title\":\"%s\",\"endAt\":\"%s\"}", ad.Title, ad.EndAt.Format(time.RFC3339)))
	}

	_, err := cache.UpdateCacheScript.Run(rds_ctx, rds, []string{key}, cmd).Result()
	if err != nil {
		fmt.Println("SetHotSpotAdToRedis", err)
		return err
	}

	return nil
}
