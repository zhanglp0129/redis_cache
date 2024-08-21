package redis_cache

import (
	"context"
	"math/rand"
	"time"
)

type CacheConfig struct {
	cacheTime time.Duration
	ctx       context.Context
	flush     bool
	write     bool
}

// 获取默认缓存时间，默认为29-31分钟
func getDefaultCacheTime() time.Duration {
	return 29*time.Minute + time.Duration(rand.Int63()%(2*int64(time.Minute)))
}

type Option func(c *CacheConfig)

// WithCacheTime 指定缓存时间，默认为29-31分钟
func WithCacheTime(cacheTime time.Duration) Option {
	return func(c *CacheConfig) {
		c.cacheTime = cacheTime
	}
}

// WithContext 指定使用redis时的context
func WithContext(ctx context.Context) Option {
	return func(c *CacheConfig) {
		c.ctx = ctx
	}
}

// WithFlushCacheTime 在命中缓存时刷新缓存时间，默认为false
func WithFlushCacheTime(flush bool) Option {
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
