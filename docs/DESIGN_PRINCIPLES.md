# 架构设计原则

本文档描述 PaperTok 后端服务采用的架构设计哲学和核心原则。

---

## 1. 设计目标

| 目标 | 说明 |
|------|------|
| **可维护** | 代码组织清晰，易于理解和修改 |
| **可扩展** | 新增功能只需新建模块，无需修改已有代码 |
| **松耦合** | 模块间通过接口隔离，无直接依赖 |
| **可测试** | 每个模块可独立单元测试 |

---

## 2. 架构模式：Vertical Slice Architecture

### 2.1 核心概念

| 术语 | 说明 |
|------|------|
| **Vertical Slice** | 按业务功能垂直划分代码，每个 Slice 包含完整处理逻辑 |
| **Repository Pattern** | 数据访问抽象层，隔离业务逻辑与数据存储 |
| **Dependency Injection** | 所有依赖通过构造函数注入，便于测试 |
| **Interface Segregation** | 每个模块只声明需要的接口方法（ISP 原则） |

### 2.2 架构优势

| 优势 | 说明 |
|------|------|
| **High Cohesion** | 每个功能的代码集中在一个目录 |
| **Loose Coupling** | 模块通过接口依赖，无直接耦合 |
| **Independent Testability** | 模块只依赖接口，Mock 简单 |
| **Progressive Development** | 新功能 = 新模块目录 |
| **Easy Deletion** | 废弃功能删除目录即可 |

### 2.3 与传统分层架构对比

```
传统 Layered Architecture（水平分层）:
┌─────────────────────────────────────┐
│           Controllers               │
├─────────────────────────────────────┤
│  GodService (所有逻辑混在一起)       │
├─────────────────────────────────────┤
│           Repositories              │
└─────────────────────────────────────┘

问题：
- 测试一个功能要 Mock 很多依赖
- 新功能要修改已有大文件
- 代码耦合严重

Vertical Slice Architecture（垂直切片）:
┌─────────────┐ ┌─────────────┐
│  Feature A  │ │  Feature B  │
│   Service   │ │   Service   │
│      ↓      │ │      ↓      │
│   Repo      │ │   Repo      │
└─────────────┘ └─────────────┘
        ↓               ↓
┌─────────────────────────────────────┐
│       Shared Kernel (Core)          │
└─────────────────────────────────────┘

优势：
- 高内聚，低耦合
- 独立测试，独立部署
- 新功能不影响已有代码
```

---

## 3. 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                      API Layer                               │
│                   (HTTP Handlers)                            │
├─────────────────────────────────────────────────────────────┤
│                     Facade Layer                             │
│               (Unified Entry Point)                          │
├─────────────────────────────────────────────────────────────┤
│    Feature Slices              │      Core Services          │
│  (Vertical Slices)             │    (Shared Services)        │
├─────────────────────────────────────────────────────────────┤
│                   Repository Layer                           │
│                      (Ports)                                 │
├─────────────────────────────────────────────────────────────┤
│                Infrastructure Layer                          │
│                  (Shared Kernel)                             │
│              Cache │ HTTPClient │ ...                        │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. 依赖方向原则

### 4.1 依赖图

```
API Layer
    ↓
Facade Layer
    ↓
Features ────→ Core Services
    ↓              ↓
Repository ←─────────┘
    ↓
Infrastructure
```

### 4.2 规则

| 规则 | 说明 |
|------|------|
| ✅ Features 可依赖 | Core, Repository, Infra |
| ✅ Core 可依赖 | Infra |
| ❌ Features 之间 | 不能相互依赖 |
| ❌ Core | 不能依赖 Features |
| ❌ 循环依赖 | 禁止 |

---

## 5. 模块标准规范

### 5.1 目录结构

```
internal/features/featurename/
├── README.md          # 模块文档（推荐）
├── interface.go       # 对外接口定义（必须）
├── deps.go            # 依赖接口定义（必须）
├── types.go           # 领域类型（可选）
├── errors.go          # 错误定义（推荐）
├── service.go         # 实现（必须）
└── service_test.go    # 单元测试（必须）
```

### 5.2 interface.go

定义对外暴露的接口：

```go
// Service defines the interface for paper feed operations.
type Service interface {
    // GetFeed fetches papers for the feed.
    // @Params:
    //   - ctx: context for cancellation
    //   - req: fetch request parameters
    // @Returns:
    //   - []*Paper: list of papers
    //   - error: if fetch fails
    GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error)
}
```

### 5.3 deps.go

声明本模块需要的依赖接口：

```go
// arxivService defines the arXiv service capability required.
type arxivService interface {
    FetchByCategory(ctx context.Context, req *arxiv.FetchRequest) ([]*arxiv.Paper, error)
}

// paperRepository defines the paper repository capability required.
type paperRepository interface {
    GetByCategory(ctx context.Context, category string) ([]*Paper, bool)
    SaveByCategory(ctx context.Context, category string, papers []*Paper, ttl time.Duration)
}
```

### 5.4 errors.go

定义领域错误：

```go
var (
    ErrNotFound = errors.New("paper not found")
    ErrInvalidRequest = errors.New("invalid request")
)

func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }
```

### 5.5 service.go

实现核心逻辑：

```go
type Impl struct {
    arxivSvc  arxivService
    paperRepo paperRepository
    cacheTTL  time.Duration
}

var _ Service = (*Impl)(nil) // compile-time check

func New(arxivSvc arxivService, repo paperRepository, ttl time.Duration) *Impl {
    return &Impl{arxivSvc: arxivSvc, paperRepo: repo, cacheTTL: ttl}
}

func (s *Impl) GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error) {
    // 1. Check cache
    // 2. Fetch from arXiv
    // 3. Save to cache
    // 4. Return
}
```

---

## 6. 单元测试原则

### 6.1 测试原则

| 原则 | 说明 |
|------|------|
| **ISP** | 只声明需要的接口方法，减少 Mock 复杂度 |
| **Mock Objects** | 只 Mock 依赖的接口 |
| **Test Isolation** | 每个模块独立测试，测试间无状态共享 |
| **AAA Pattern** | Arrange-Act-Assert 结构化测试 |

### 6.2 测试示例

```go
func TestImpl_GetFeed(t *testing.T) {
    // Arrange
    mockArxiv := &mockArxivService{papers: [...]}
    mockRepo := newMockPaperRepository()
    svc := New(mockArxiv, mockRepo, 5*time.Minute)

    // Act
    papers, err := svc.GetFeed(ctx, &FetchRequest{Category: "cs.AI"})

    // Assert
    assert.NoError(t, err)
    assert.Len(t, papers, 1)
}
```

---

## 7. 添加新功能指南

### 7.1 添加新 Feature

1. 创建目录 `internal/features/newfeature/`
2. 创建标准文件：interface.go, deps.go, service.go, service_test.go
3. 实现 Service 接口
4. 编写单元测试
5. 在 Facade 中注册

### 7.2 添加新 Core Service

1. 创建目录 `internal/core/newservice/`
2. 创建标准文件
3. 实现 Service 接口
4. 编写单元测试
5. 在 Facade 中初始化

---

## 8. 代码风格

### 8.1 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写，单词 | `paperfeed` |
| 接口 | 首字母大写 | `Service` |
| 私有接口 | 首字母小写 | `arxivService` |
| 实现 | `Impl` | `Impl` |
| 构造函数 | `New` | `New()` |

### 8.2 注释规范

```go
// Service defines the interface for paper feed operations.
// This interface is used by the Facade layer.
type Service interface {
    // GetFeed fetches papers for the feed.
    // @Params:
    //   - ctx: context for cancellation and tracing
    //   - req: fetch request parameters
    // @Returns:
    //   - []*Paper: list of papers
    //   - error: if fetch fails
    GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error)
}
```

---

## 9. 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) - 架构设计文档
- [local-guides/](../backend/local-guides/) - 本地开发指南

---

**文档版本**：v1.0
**最后更新**：2026-01-26
