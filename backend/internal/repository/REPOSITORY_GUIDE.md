# Repository 开发指南

本指南定义 Repository 模块的开发规范。

---

## 1. Repository 定位

Repository 是业务逻辑与数据存储之间的抽象层。

| 角色 | 说明 |
|------|------|
| **数据访问抽象** | 隐藏底层存储细节 |
| **领域边界** | 每个领域实体一个 Repository |
| **可测试性** | 便于 Mock 和单元测试 |

---

## 2. 目录结构

```
repository/entityname/
├── README.md          # 模块文档（推荐）
├── interface.go       # 接口和数据类型（必须）
└── memory.go          # 内存实现（或其他实现）
```

---

## 3. 文件规范

### interface.go

```go
// Entity represents the domain entity.
type Entity struct {
    ID   string
    Name string
}

// Repository defines the data access interface.
type Repository interface {
    GetByID(ctx context.Context, id string) (*Entity, bool)
    Save(ctx context.Context, entity *Entity, ttl time.Duration)
    Delete(ctx context.Context, id string)
}
```

### memory.go

```go
type MemoryRepository struct {
    cache cache.Cache
}

var _ Repository = (*MemoryRepository)(nil)

func NewMemoryRepository(c cache.Cache) *MemoryRepository {
    return &MemoryRepository{cache: c}
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*Entity, bool) {
    // implementation
}
```

---

## 4. 存储选择原则

| 数据特征 | 存储选择 |
|---------|---------|
| 高频读写、临时数据 | 内存缓存 |
| 持久化、复杂查询 | PostgreSQL |
| 分布式缓存 | Redis |

---

## 5. 现有 Repositories

| Repository | 职责 | 存储 |
|------------|------|------|
| `paper` | 论文数据缓存 | 内存 |
