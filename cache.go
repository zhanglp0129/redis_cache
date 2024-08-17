package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

// QueryWithCache 带着缓存查询
// model，传出参数，数据模型
// query，查询函数，缓存未命中时调用
// return，是否命中缓存
func QueryWithCache[T any](rdb redis.UniversalClient, key string, model *T, query func() (*T, error), options ...Option) (bool, error) {
	// 获取缓存配置
	config := CacheConfig{
		cacheTime: getDefaultCacheTime(),
		ctx:       context.Background(),
	}
	for _, opt := range options {
		opt(&config)
	}

	if n, err := rdb.Exists(config.ctx, key).Result(); err == nil && n > 0 {
		// 命中缓存
		return true, readFromRedis(rdb, key, model, &config)
	}

	// 未命中缓存
	data, err := query()
	if err != nil {
		return false, err
	}
	*model = *data
	return false, writeToRedis(rdb, key, model, &config)
}

// 获取默认缓存时间
func getDefaultCacheTime() time.Duration {
	return time.Duration(55+rand.Int()%11) * time.Second
}

// 从redis读取数据，并保存在model中
func readFromRedis[T any](rdb redis.UniversalClient, key string, model *T, cfg *CacheConfig) error {
	data, err := rdb.Get(cfg.ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), model)
	if err != nil {
		return err
	}

	// 刷新缓存时间
	rdb.Expire(cfg.ctx, key, cfg.cacheTime)

	return nil
}

// 将model的数据写入redis
func writeToRedis[T any](rdb redis.UniversalClient, key string, model *T, cfg *CacheConfig) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return rdb.SetEx(cfg.ctx, key, string(data), cfg.cacheTime).Err()
}

// DeleteCache 删除缓存
func DeleteCache(rdb redis.UniversalClient, key string) error {
	return rdb.Del(context.Background(), key).Err()
}
