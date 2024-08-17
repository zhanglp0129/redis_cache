# RedisCache
Golang使用redis缓存，将数据转换为json，并存入redis

## 使用
1. 安装依赖
```shell
go get -u github.com/zhanglp0129/redis_cache
```

2. 带着缓存查询
```go
// 创建redis客户端连接
rdb := redis.NewClient(&redis.Options{
Addr: "127.0.0.1:6379",
})
// redis key
key := "test_key"
// model可以是任意类型
model := ""

// 带着缓存查询，需要传入查询函数，在缓存未命中时调用
hit, err := QueryWithCache(rdb, key, &model, func() (*string, error) {
    data := "test_data"
    return &data, nil
})

// 打印结果
fmt.Println("是否命中缓存：", hit)
fmt.Println("结果：", model)
```

3. 删除缓存
```go
rdb := redis.NewClient(&redis.Options{
    Addr: "127.0.0.1:6379",
})
key := "test_key"

err := DeleteCache(rdb, key)
```