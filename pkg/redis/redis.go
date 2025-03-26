package redis_utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	Client *redis.Client
	Ctx    context.Context
}

type CacheData struct {
	Exp          int    `json:"exp"`
	LoginSession string `json:"login_session"`
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
}

func NewRedisUtil(rdb *redis.Client) *RedisUtil {
	return &RedisUtil{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisUtil) SetCacheKey(key string, value interface{}, ctx context.Context) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	duration := time.Duration(0)
	return r.Client.Set(ctx, key, data, duration).Err()
}

func (r *RedisUtil) SetEx(key string, data interface{}, expiredTime time.Duration, ctx context.Context) error {
	jsonData, errData := json.Marshal(data)
	if errData != nil {
		return errData
	}

	err := r.Client.Set(ctx, key, jsonData, expiredTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisUtil) GetCache(key string, result interface{}) error {
	value, err := r.Client.Get(r.Ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), result)
}

func Get[T any](r *RedisUtil, key string, ctx context.Context) (*T, error) {
	// Get data from Redis
	result, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Unmarshal into generic type
	var data T
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return &data, nil
}

func (r *RedisUtil) GetCacheKey(key string, ctx context.Context) (*CacheData, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return nil, err
	}
	var data CacheData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &data, nil
}

func (r *RedisUtil) DeleteCache(key string) error {
	deleted, err := r.Client.Del(r.Ctx, key).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return fmt.Errorf("key %s does not exist", key)
	}
	return nil
}

func (r *RedisUtil) AddToBlockList(token string, expiration time.Duration) error {
	return r.Client.Set(r.Ctx, "blokclist:"+token, "revoked", expiration).Err()
}

func (r *RedisUtil) IsTokenRevoked(token string) (bool, error) {
	val, err := r.Client.Get(r.Ctx, "blocklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, nil
	}
	return val == "revoked", nil
}

func (r *RedisUtil) RateLimit(key string, limit int, window time.Duration) (bool, error) {
	count, err := r.Client.Incr(r.Ctx, key).Result()
	if err != nil {
		return false, nil
	}
	if count == 1 {
		_ = r.Client.Expire(r.Ctx, key, window).Err()
	}
	return count <= int64(limit), nil
}

func (r *RedisUtil) CloseConnection() error {
	return r.Client.Close()
}
