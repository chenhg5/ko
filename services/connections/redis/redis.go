package redis

import (
	"github.com/go-redis/redis"
	"ninepic/config"
	"time"
)

type ClientType struct {
	RedisCon *redis.Client
}

var Client *ClientType

func init() {
	Client = &ClientType{
		RedisCon: redis.NewClient(&redis.Options{
			Addr:     config.GetEnv().REDIS_IP + ":" + config.GetEnv().REDIS_PORT,
			Password: config.GetEnv().REDIS_PASSWORD, // no password set
			DB:       config.GetEnv().REDIS_DB,       // use default DB
		}),
	}
}

func (Client *ClientType) Set(key string, value interface{}, expiration time.Duration) *redis.Client {
	err := (*Client).RedisCon.Set(key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
	return (*Client).RedisCon
}

func (Client *ClientType) Get(key string) (string, *redis.Client) {
	val, err := (*Client).RedisCon.Get(key).Result()

	if err == redis.Nil {
		return "", (*Client).RedisCon
	}

	if err != nil {
		panic(err)
	}

	return val, (*Client).RedisCon
}

func (Client *ClientType) Del(key string) *redis.Client {
	_, err := (*Client).RedisCon.Del(key).Result()
	if err != nil {
		panic(err)
	}
	return (*Client).RedisCon
}
