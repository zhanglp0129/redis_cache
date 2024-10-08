package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

type QueryFunc[T any] func() (T, error)

// QueryWithCache 带着缓存查询
// model，传出参数，数据模型
// query，查询函数，缓存未命中时调用
// return，是否命中缓存
func QueryWithCache[T any](rdb redis.UniversalClient, key string, model *T, query QueryFunc[*T], options ...Option) (bool, error) {
	// 获取缓存配置
	config := CacheConfig{
		cacheTime: DefaultCacheTime,
		diff:      DefaultCacheTimeDiff,
		ctx:       context.Background(),
		flush:     false,
		write:     true,
	}
	for _, opt := range options {
		opt(&config)
	}

	data, err := rdb.Get(config.ctx, key).Result()
	if err == nil {
		// 命中缓存
		err = json.Unmarshal([]byte(data), model)
		if err != nil {
			// json解析出错，以未命中缓存处理
			return false, cacheMiss(rdb, key, model, query, &config)
		}

		// 刷新缓存时间
		if config.flush {
			return true, rdb.Expire(config.ctx, key, computeCacheTime(config.cacheTime, config.diff)).Err()
		}
		return true, nil
	} else if errors.Is(err, redis.Nil) {
		// 未命中缓存
		return false, cacheMiss(rdb, key, model, query, &config)
	}

	// 出现错误
	return false, err
}

// 缓存未命中
func cacheMiss[T any](rdb redis.UniversalClient, key string, model *T, query QueryFunc[*T], config *CacheConfig) error {
	res, err := query()
	if err != nil {
		return err
	}
	*model = *res

	// 不写入缓存
	if !config.write {
		return nil
	}

	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return rdb.SetEx(config.ctx, key, string(data), computeCacheTime(config.cacheTime, config.diff)).Err()
}
