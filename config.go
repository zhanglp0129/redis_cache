package redis_cache

import (
	"context"
	"time"
)

type CacheConfig struct {
	cacheTime time.Duration
	ctx       context.Context
}
