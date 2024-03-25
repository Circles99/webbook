package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: "webook-redis:16379",
	})
	return rdb
}
