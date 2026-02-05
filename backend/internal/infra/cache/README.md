# Cache Infrastructure

> 通用缓存基础设施，提供内存缓存能力

---

## 职责

- 提供通用缓存接口
- 实现内存缓存
- 管理 TTL 和过期清理

---

## 接口

```go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
}
```

---

## 文件结构

| 文件 | 说明 |
|------|------|
| `interface.go` | 缓存接口定义 |
| `memory.go` | 内存缓存实现 |

---

## 实现

### MemoryCache

基于 Go map 和 RWMutex 的线程安全内存缓存。

特性：
- 支持 TTL 过期
- 自动清理过期项（每分钟）
- 线程安全

```go
cache := cache.NewMemoryCache()

// 设置（5分钟过期）
cache.Set("key", value, 5*time.Minute)

// 获取
value, found := cache.Get("key")

// 删除
cache.Delete("key")

// 清空
cache.Clear()
```

---

## 扩展

可以添加其他缓存实现：
- Redis 缓存
- 分布式缓存
