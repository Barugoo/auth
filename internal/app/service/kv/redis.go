package kv

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisKV struct {
	redisClient *redis.Client
}

func (kv *RedisKV) Set(ctx context.Context, key, value string) (bool, error) {
	span := kv.StartSpan(ctx, "RedisSet")
	defer span.Finish()

	ok, err := kv.set(key, value)
	if err != nil {
		err = kv.wrapError(err)
	}
	return ok, err
}

func (kv *RedisKV) set(key, value string) (bool, error) {
	err := kv.redisClient.Set(key, value, 0).Err()
	if err != nil {
		return false, err
	}
	val, err := kv.redisClient.Get(key).Result()
	if err != nil {
		return false, err
	}
	if val != value {
		return false, err
	}
	return true, nil
}

func (kv *RedisKV) Get(ctx context.Context, key string) (bool, error) {
	span := kv.StartSpan(ctx, "RedisGet")
	defer span.Finish()

	ok, err := kv.set(key, value)
	if err != nil {
		err = kv.wrapError(err)
	}
	return ok, err
}

func (kv *RedisKV) Get(ctx context.Context, key string) (string, error) {
	val, err := kv.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (kv *RedisKV) wrapError(err error) error {
	return fmt.Errorf("redis: %v", err)
}
