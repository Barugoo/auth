package redis

import (
	"github.com/barugoo/oscillo-auth/config"
	"github.com/go-redis/redis/v7"
)

func NewRedisClient(config *config.ServiceConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
