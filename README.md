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

## TODO
- [ ] 文档和注释改为英文，同时优化文档
- [ ] 带着缓存查询时，能够加锁，防止缓存多次写入
- [ ] 添加布隆过滤器的实现
- [ ] 能针对不同的数据类型，使用不同的数据结构缓存
- [ ] 针对不同的redis数据结构，使用不同的方式修改缓存
- [ ] 能主动刷新缓存时间
