# Facade Layer

> 统一业务入口，管理所有 Feature 和 Core Service 的实例化和依赖注入

---

## 职责

- **统一入口**：API Layer 只与 Facade 交互
- **依赖组装**：管理所有模块的实例化
- **依赖注入**：将依赖注入到各个 Feature

---

## 主要文件

| 文件 | 说明 |
|------|------|
| `service.go` | Facade 实现 |

---

## 使用方式

```go
// 初始化 Facade
f := facade.New(facade.Config{
    ArxivBaseURL: "http://export.arxiv.org/api/query",
    HTTPTimeout:  10 * time.Second,
    CacheTTL:     5 * time.Minute,
    CacheEnabled: true,
})

// 在 Handler 中使用
papers, err := f.GetPaperFeed(ctx, "cs.AI", 20, 0, "lastUpdatedDate")
```

---

## 公开方法

| 方法 | 说明 |
|------|------|
| `GetPaperFeed()` | 获取论文推荐流 |
| `SearchPapers()` | 搜索论文 |
| `GetPaperByID()` | 获取论文详情 |

---

## 依赖关系

```
Facade
├── paperfeed.Service
├── papersearch.Service
├── arxiv.Service
└── paper.Repository
```
