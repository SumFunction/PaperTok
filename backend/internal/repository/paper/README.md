# Paper Repository

> 论文数据访问层，提供缓存存储能力

---

## 职责

- 按分类缓存论文列表
- 按 ID 缓存单篇论文
- 管理缓存失效

---

## 接口

```go
type Repository interface {
    GetByCategory(ctx context.Context, category string) ([]*Paper, bool)
    SaveByCategory(ctx context.Context, category string, papers []*Paper, ttl time.Duration)
    GetByID(ctx context.Context, id string) (*Paper, bool)
    Save(ctx context.Context, paper *Paper, ttl time.Duration)
    InvalidateCategory(ctx context.Context, category string)
    Clear(ctx context.Context)
}
```

---

## 文件结构

| 文件 | 说明 |
|------|------|
| `interface.go` | 接口和数据类型定义 |
| `memory.go` | 内存缓存实现 |

---

## 实现

### MemoryRepository

基于内存缓存的实现，使用 `infra/cache` 模块。

```go
repo := paper.NewMemoryRepository(cache)

// 保存
repo.SaveByCategory(ctx, "cs.AI", papers, 5*time.Minute)

// 获取
papers, found := repo.GetByCategory(ctx, "cs.AI")

// 失效
repo.InvalidateCategory(ctx, "cs.AI")
```

---

## 缓存键设计

| 键模式 | 说明 |
|--------|------|
| `papers:category:{category}` | 分类论文列表 |
| `papers:id:{id}` | 单篇论文 |

---

## 扩展

可以添加其他存储实现：
- `redis.go` - Redis 缓存
- `postgres.go` - PostgreSQL 持久化
