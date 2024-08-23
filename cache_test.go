package redis_cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestQueryWithCacheHit(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	key := "test_key"

	var result string
	_, err := QueryWithCache(rdb, key, &result, func() (*string, error) {
		data := "test_data"
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "test_data" {
		t.Fatal("数据不一致")
	}

	result = ""
	hit, err := QueryWithCache(rdb, key, &result, func() (*string, error) {
		data := "test_data"
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "test_data" {
		t.Fatal("数据不一致")
	}
	if !hit {
		t.Fatal("未命中缓存")
	}
}

func TestQueryWithCacheMiss(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	key := "test_key"

	err := DeleteCache(rdb, key)
	if err != nil {
		t.Fatal(err)
	}

	var result string
	hit, err := QueryWithCache(rdb, key, &result, func() (*string, error) {
		data := "test_data"
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "test_data" {
		t.Fatal("数据不一致")
	}
	if hit {
		t.Fatal("命中缓存")
	}
}

func TestDeleteCache(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	key := "test_key"

	err := DeleteCache(rdb, key)
	if err != nil {
		t.Fatal(err)
	}
	if n, err := rdb.Exists(context.Background(), key).Result(); err != nil || n > 0 {
		t.Fatal("删除缓存失败")
	}

	// 删除不存在的缓存
	err = DeleteCache(rdb, key)
	if err != nil {
		t.Fatal(err)
	}
	if n, err := rdb.Exists(context.Background(), key).Result(); err != nil || n > 0 {
		t.Fatal("删除缓存失败")
	}
}

func TestCacheIncrBy(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	key := "test_key"

	// 删除缓存
	err := DeleteCache(rdb, key)
	if err != nil {
		t.Fatal(err)
	}

	// 写入缓存，10
	var result int64
	_, err = QueryWithCache(rdb, key, &result, func() (*int64, error) {
		var data int64 = 10
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// 缓存+5
	err = CacheIncrBy(rdb, key, 5)
	if err != nil {
		t.Fatal(err)
	}

	// 再读取缓存，并判断是否为15
	_, err = QueryWithCache(rdb, key, &result, func() (*int64, error) {
		var data int64 = 10
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != 15 {
		t.Fatal("缓存错误")
	}

	// 删除缓存
	err = DeleteCache(rdb, key)
	if err != nil {
		t.Fatal(err)
	}

	// 缓存+5
	err = CacheIncrBy(rdb, key, 5)
	if err != nil {
		t.Fatal(err)
	}

	// 再读取缓存，并判断是否命中
	hit, err := QueryWithCache(rdb, key, &result, func() (*int64, error) {
		var data int64 = 10
		return &data, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if hit || result != 10 {
		t.Fatal("缓存错误")
	}
}
