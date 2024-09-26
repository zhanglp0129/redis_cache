package redis_cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// DeleteCacheCtx 删除缓存
func DeleteCacheCtx(ctx context.Context, rdb redis.UniversalClient, key string) error {
	return rdb.Del(ctx, key).Err()
}

// DeleteCache 删除缓存
func DeleteCache(rdb redis.UniversalClient, key string) error {
	return DeleteCacheCtx(context.Background(), rdb, key)
}

// DeleteCacheToPipeCtx 将删除缓存的操作加入pipe，不会执行
func DeleteCacheToPipeCtx(ctx context.Context, pipe redis.Pipeliner, keys ...string) {
	pipe.Del(ctx, keys...)
}

// DeleteCacheToPipe 将删除缓存的操作加入pipe，不会执行
func DeleteCacheToPipe(pipe redis.Pipeliner, keys ...string) {
	DeleteCacheToPipeCtx(context.Background(), pipe, keys...)
}

// CacheIncrByCtx 如果只缓存了一个整数，可以修改缓存，让这个整数增加value
func CacheIncrByCtx(ctx context.Context, rdb redis.UniversalClient, key string, value int64) error {
	return rdb.Eval(ctx, cacheIncrByLuaScript, []string{key}, value).Err()
}

// CacheIncrBy 如果只缓存了一个整数，可以修改缓存，让这个整数增加value
func CacheIncrBy(rdb redis.UniversalClient, key string, value int64) error {
	return CacheIncrByCtx(context.Background(), rdb, key, value)
}

// CacheIncrByToPipeCtx 将缓存的整数增加value。会把操作加入pipe，不会执行。一次只能添加一个
func CacheIncrByToPipeCtx(ctx context.Context, pipe redis.Pipeliner, key string, value int64) {
	pipe.Eval(ctx, cacheIncrByLuaScript, []string{key}, value)
}

// CacheIncrByToPipe 将缓存的整数增加value。会把操作加入pipe，不会执行。一次只能添加一个
func CacheIncrByToPipe(pipe redis.Pipeliner, key string, value int64) {
	CacheIncrByToPipeCtx(context.Background(), pipe, key, value)
}
