package cache

import (
	"context"
	"fmt"
	"github.com/LiteMove/light-stack/internal/shared/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

// Init 初始化Redis连接
func Init() error {
	cfg := config.Get()

	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试连接
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

// GetRDB 获取Redis客户端
func GetRDB() *redis.Client {
	return RDB
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除缓存
func Del(keys ...string) error {
	return RDB.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func Exists(key string) (bool, error) {
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}

// SetNX 设置缓存（仅当key不存在时）
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return RDB.SetNX(ctx, key, value, expiration).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}

// HSet 设置哈希字段
func HSet(key string, field string, value interface{}) error {
	return RDB.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func HGet(key string, field string) (string, error) {
	return RDB.HGet(ctx, key, field).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return RDB.HDel(ctx, key, fields...).Err()
}

// Close 关闭Redis连接
func Close() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}
