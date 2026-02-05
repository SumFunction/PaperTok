# Core Service 开发指南

本指南定义 Core Service 模块的开发规范。

---

## 1. Core Service 定位

Core Service 是**共享的服务模块**，为多个 Feature 提供通用能力。

| 维度 | Core Service | Feature |
|------|--------------|---------|
| 职责 | 通用服务能力 | 特定业务场景 |
| 依赖方向 | 被 Feature 依赖 | 依赖 Core/Infra |
| 稳定性 | 高（接口变更影响大） | 中（相对独立） |

**依赖约束**：
- ✅ Core 可依赖 Infra
- ✅ Feature 可依赖 Core
- ❌ Core 不能依赖 Feature

---

## 2. 目录结构

```
core/servicename/
├── README.md          # 模块文档（推荐）
├── interface.go       # 对外接口定义（必须）
├── deps.go            # 依赖接口定义（必须）
├── types.go           # 数据类型（可选）
├── errors.go          # 错误定义（推荐）
├── client.go          # 实现（必须）
└── client_test.go     # 单元测试（必须）
```

---

## 3. 文件规范

### interface.go

```go
// Service defines the interface for xxx operations.
type Service interface {
    // Method does something.
    Method(ctx context.Context, param Type) (Result, error)
}
```

### deps.go

```go
// httpClient defines the HTTP capability required.
type httpClient interface {
    Get(ctx context.Context, url string) (*http.Response, error)
}
```

### types.go

```go
// Paper represents a paper from arXiv.
type Paper struct {
    ID    string
    Title string
}

// FetchRequest contains fetch parameters.
type FetchRequest struct {
    Category string
    Limit    int
}
```

### errors.go

```go
var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid")
)

func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }
```

### client.go

```go
type Client struct {
    httpClient httpClient
}

var _ Service = (*Client)(nil)

func NewClient(cfg Config, client httpClient) *Client {
    return &Client{httpClient: client}
}

func (c *Client) Method(ctx context.Context, param Type) (Result, error) {
    // implementation
}
```

---

## 4. 接口设计原则

| 原则 | 说明 |
|------|------|
| **稳定性** | 接口一旦发布，保持向后兼容 |
| **通用性** | 避免为特定 Feature 定制 |
| **最小化** | 只暴露必要的方法 |

---

## 5. 现有 Core Services

| Service | 职责 |
|---------|------|
| `arxiv` | arXiv API 客户端 |
