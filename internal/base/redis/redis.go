package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"graph-med/internal/base/conf"
	"time"
)

var Client *redis.Client

func Setup(config *conf.AllConfig) error {
	redisConf := &config.Redis
	Client = redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       redisConf.Db,
		PoolSize: redisConf.PoolSize,
	})

	// 验证连接
	res, err := Client.Ping().Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return err
	} else {
		fmt.Printf("Connect Successful! Ping => %v\n", res)
		return nil
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	value, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(key, value, expiration).Err()
}

func Get[T any](key string) (T, error) {
	var value T
	err := json.Unmarshal([]byte(Client.Get(key).Val()), &value)
	return value, err
}

func Del(key string) error {
	return Client.Del(key).Err()
}
