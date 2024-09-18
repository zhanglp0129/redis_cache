package redis_cache

import (
	"context"
	"math/rand"
	"time"
)

var (
	// DefaultCacheTime 默认缓存时间，30min
	DefaultCacheTime = 30 * time.Minute
	// DefaultCacheTimeDiff 默认缓存时间的上下偏差，5min
	DefaultCacheTimeDiff = 5 * time.Minute
)

type CacheConfig struct {
	cacheTime time.Duration
	ctx       context.Context
	flush     bool
	write     bool
}

// 获取默认缓存时间，默认为 [25, 35) min
func getDefaultCacheTime() time.Duration {
	diff := time.Duration(rand.Int63()%int64(2*DefaultCacheTimeDiff)) - DefaultCacheTimeDiff
	return max(0, DefaultCacheTime+diff)
}

type Option func(c *CacheConfig)

// WithCacheTime 指定缓存时间，默认为 [25, 35) min
func WithCacheTime(cacheTime time.Duration) Option {
	return func(c *CacheConfig) {
		c.cacheTime = cacheTime
	}
}

// WithContext 指定使用redis时的context，默认为background
func WithContext(ctx context.Context) Option {
	return func(c *CacheConfig) {
		c.ctx = ctx
	}
}

// FlushCacheTime 在命中缓存时刷新缓存时间，默认为false
func FlushCacheTime(flush bool) Option {
	return func(c *CacheConfig) {
		c.flush = flush
	}
}

// WriteCache 是否在缓存未命中时，写入缓存，默认为true
func WriteCache(write bool) Option {
	return func(c *CacheConfig) {
		c.write = write
	}
}
