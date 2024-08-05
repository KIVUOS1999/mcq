package rediscache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bkend-redis/constants"
	"github.com/go-redis/redis/v8"
)

type CacheStruct struct {
	redisClient *redis.Client
	expireTTL   int
}

func New() *CacheStruct {
	redisHost := os.Getenv(constants.RedisHostKey)
	expireTTL := os.Getenv(constants.RedisExpiryKey)

	expireTTLInt, err := strconv.Atoi(expireTTL)
	if err != nil {
		expireTTLInt = 1 // Passed value is in days
	}

	if redisHost == "" {
		panic("No redis host is provided")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0,
	})

	return &CacheStruct{
		redisClient: redisClient,
		expireTTL:   expireTTLInt,
	}
}

func (c *CacheStruct) SetValue(key string, value string) error {
	cmd := c.redisClient.Set(context.Background(), key, value, time.Duration(c.expireTTL)*time.Hour*constants.Day)
	if cmd.Err() != nil {
		log.Println("Err in setKey value", cmd.Err())
		return cmd.Err()
	}

	return nil
}

func (c *CacheStruct) GetValue(key string) (string, error) {
	value, err := c.redisClient.Get(context.Background(), key).Result()
	return value, err
}

func (c *CacheStruct) DeleteKey(key string) error {
	cmd := c.redisClient.Del(context.Background(), key)
	if cmd.Err() != nil {
		log.Println("Err in deleteKey value", cmd.Err())
		return cmd.Err()
	}

	return nil
}
