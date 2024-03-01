package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/cache"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/models"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PublicService struct {
	db  *gorm.DB
	rds *cache.RedisClient
}
type DBClient gorm.DB
type RedisClient cache.RedisClient

func NewPublicService() *PublicService {
	return &PublicService{db: database.DB, rds: cache.Rds}
}

// []models.Advertisement, error
func (srv *PublicService) GetAdvertisements(req *schemas.PublicAdRequest) (any, error) {
	var ads []models.Advertisement
	// get from redis
	err := srv.getAdFromRedis(req, &ads)
	if err == nil {
		return ads, nil
	}

	err = srv.getAdFromDB(req, &ads)
	if err != nil {
		return nil, err
	}

	// TTL: 5 minutes
	// store full result to redis in goroutine
	go srv.setAdToRedis(req, &ads)

	// set offset and limit for sql result
	return ads[req.Offset : req.Offset+req.Limit], nil
}

func (srv *PublicService) getAdFromRedis(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	key := cache.PublicAdKey(req)
	rds := srv.rds.Client
	rds_ctx := srv.rds.Context

	// redis check key exist
	if rds.Exists(rds_ctx, key).Val() == 0 {
		fmt.Println("key does not exist")
		return redis.Nil
	}

	val, err := rds.ZRangeByScore(rds_ctx, key, &redis.ZRangeBy{
		Min: strconv.Itoa(req.Offset + 1),
		Max: strconv.Itoa(req.Offset + req.Limit),
	}).Result()
	if err != nil {
		return err
	}

	utils.PrintJson(val)

	// bind to ads
	*ads = make([]models.Advertisement, len(val))
	for i, v := range val {
		if err := json.Unmarshal([]byte(v), &(*ads)[i]); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (srv *PublicService) getAdFromDB(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
	// StartAt <= now <= EndAt
	now := time.Now()
	fmt.Println(now)

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
	if req.Age != 0 {
		tx = tx.Where("age_start <= ? AND age_end >= ?", req.Age, req.Age)
	}

	// limit and offset
	// tx = tx.Limit(req.Limit).Offset(req.Offset)

	if err := tx.Find(ads).Error; err != nil {
		return err
	}

	return nil
}

func (srv *PublicService) setAdToRedis(req *schemas.PublicAdRequest, ads *[]models.Advertisement) error {
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
			fmt.Println("setAdToRedis", err)
			return err
		}
		return nil
	}

	// use `redis.Do` is more efficient than `redis.ZAdd`
	cmd := make([]interface{}, 0, len(*ads)*2+2)
	cmd = append(cmd, "ZADD", key)
	for idx, ad := range *ads {
		cmd = append(cmd, float64(idx+1))
		cmd = append(cmd, fmt.Sprintf("{\"title\":\"%s\",\"endAt\":\"%s\"}", ad.Title, ad.EndAt.Format(time.RFC3339)))
	}
	_, err := rds.Do(rds_ctx, cmd...).Result()
	if err != nil {
		fmt.Println("setAdToRedis", err)
		return err
	}
	// set expire
	_, err = rds.Expire(rds_ctx, key, time.Minute*5).Result()
	if err != nil {
		fmt.Println("setAdToRedis", err)
		return err
	}

	return nil
}
