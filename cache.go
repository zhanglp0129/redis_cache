package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

const (
	CacheIncrByLuaScript = `
    if redis.call("EXISTS", KEYS[1]) == 1 then
        return redis.call("INCRBY", KEYS[1], ARGV[1])
    else
        return 0
    end
    `
)

type QueryFunc[T any] func() (T, error)

// QueryWithCache 带着缓存查询
// model，传出参数，数据模型
// query，查询函数，缓存未命中时调用
// return，是否命中缓存
func QueryWithCache[T any](rdb redis.UniversalClient, key string, model *T, query QueryFunc[*T], options ...Option) (bool, error) {
	// 获取缓存配置
	config := CacheConfig{
		cacheTime: getDefaultCacheTime(),
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
			return true, rdb.Expire(config.ctx, key, config.cacheTime).Err()
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

	return rdb.SetEx(config.ctx, key, string(data), config.cacheTime).Err()
}

// DeleteCache 删除缓存
func DeleteCache(rdb redis.UniversalClient, key string) error {
	return rdb.Del(context.Background(), key).Err()
}

// CacheIncrBy 如果只缓存了一个整数，可以修改缓存，让这个整数增加value
func CacheIncrBy(rdb redis.UniversalClient, key string, value int64) error {
	return rdb.Eval(context.Background(), CacheIncrByLuaScript, []string{key}, value).Err()
}

// DeleteCacheToPipe 将删除缓存的操作加入pipe，不会执行
func DeleteCacheToPipe(pipe redis.Pipeliner, keys ...string) {
	pipe.Del(context.Background(), keys...)
}

// CacheIncrByToPipe 将缓存的整数增加value。会把操作加入pipe，不会执行。一次只能添加一个
func CacheIncrByToPipe(pipe redis.Pipeliner, key string, value int64) {
	pipe.Eval(context.Background(), CacheIncrByLuaScript, []string{key}, value)
}
