package redis_cache

import "embed"

var (
	// 缓存自增的lua脚本
	cacheIncrByLuaScript string
)

//go:embed lua
var luaScriptStatic embed.FS

func init() {
	// 加载lua脚本
	tmp, err := luaScriptStatic.ReadFile("lua/cache_incr_by.lua")
	if err != nil {
		panic(err)
	}
	cacheIncrByLuaScript = string(tmp)
}
