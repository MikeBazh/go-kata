package cache

import (
	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host, port string) *redis.Client {
	// реализуйте создание клиента для Redis
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		//Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
