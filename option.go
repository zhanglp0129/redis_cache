package redis_cache

import (
	"context"
	"time"
)

type Option func(c *CacheConfig)

// WithCacheTime 指定缓存时间，默认为55-65秒；当再次访问缓存时，会刷新缓存时间
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
