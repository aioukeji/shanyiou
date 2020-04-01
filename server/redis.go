package server

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
)

func RedisSet(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(key, p, 0).Err()
}

func RedisZCount(key string) (int64, error) {
	return redisClient.ZCount(key, "-inf", "inf").Result()
}

func RedisZAdd(key string, score int, member interface{}) error {
	return redisClient.ZAdd(key, &redis.Z{Score: float64(score), Member: member}).Err()
}

func RedisGet(c *redis.Client, key string, dest interface{}) error {
	bytes, err := c.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

func RedisScan(key string, val interface{}) error {
	return redisClient.Get(key).Scan(val)
}

func RedisLPush(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.LPush(key, p).Err()
}
