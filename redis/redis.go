package redis

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-redis/redis"
)

var (
	once     sync.Once
	redisCli *redis.Client
)

// TODO: Change to nacos
func InitRedis() {
	once.Do(func() {
		redisCli = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis server address
			Password: "",               // No password
			DB:       0,                // Default database
		})
	})

	// Ping the Redis server to check if the connection is successful
	_, err := redisCli.Ping().Result()
	if err != nil {
		logs.Error("[InitRedis] Error connecting to Redis:", err)
	} else {
		logs.Info("[InitRedis] Init Redis Success")

	}
}

func Set(key string, val any, expiration ...time.Duration) error {
	if val == nil {
		val = ""
	}
	return redisCli.Set(key, val, getDuration(expiration[0])).Err()
}

func Get(key string) (string, error) {
	return redisCli.Get(key).Result()
}

func Del(key ...string) (int64, error) {
	return redisCli.Del(key...).Result()
}

func Exists(key ...string) (bool, error) {
	r, err := redisCli.Exists(key...).Result()
	return r == int64(len(key)), err
}

func getDuration(t time.Duration) time.Duration {
	if t < time.Second {
		return t * time.Second
	}
	return t
}
