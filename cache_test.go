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

}
