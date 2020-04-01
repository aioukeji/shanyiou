package server

import (
	"github.com/go-redis/redis/v7"
	"sync"
)

var fabricFlushLock sync.Mutex
var redisLock sync.Mutex

var redisOption = &redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
}
var redisClient = redis.NewClient(redisOption)
