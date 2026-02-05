# arXiv Core Service

> arXiv API 客户端，封装与 arXiv API 的交互

---

## 职责

- 按分类获取论文
- 关键词搜索论文
- 按 ID 获取单篇论文
- 解析 Atom XML 响应

---

## 接口

```go
type Service interface {
    FetchByCategory(ctx context.Context, req *FetchRequest) ([]*Paper, error)
    Search(ctx context.Context, query string, limit int) ([]*Paper, error)
    GetByID(ctx context.Context, id string) (*Paper, error)
}
```

---

## 文件结构

| 文件 | 说明 |
|------|------|
| `interface.go` | 对外接口定义 |
| `deps.go` | 依赖接口定义 |
| `types.go` | 数据类型定义 |
| `errors.go` | 错误定义 |
| `client.go` | 实现 |
| `client_test.go` | 单元测试 |

---

## 依赖

- `httpclient.HTTPClient` - HTTP 客户端

---

## 使用示例

```go
client := arxiv.NewClient(arxiv.Config{
    BaseURL: "http://export.arxiv.org/api/query",
    Timeout: 10 * time.Second,
}, httpClient)

// 按分类获取
papers, err := client.FetchByCategory(ctx, &arxiv.FetchRequest{
    Category:   "cs.AI",
    MaxResults: 20,
    SortBy:     "lastUpdatedDate",
})

// 搜索
papers, err := client.Search(ctx, "transformer", 10)

// 按 ID 获取
paper, err := client.GetByID(ctx, "2401.12345")
```

---

## 错误处理

```go
var (
    ErrFetchFailed     = errors.New("failed to fetch papers from arXiv")
    ErrSearchFailed    = errors.New("failed to search papers")
    ErrInvalidResponse = errors.New("invalid response from arXiv API")
    ErrNotFound        = errors.New("paper not found")
)

// 使用
if arxiv.IsFetchFailed(err) {
    // 处理获取失败
}
```

---

## arXiv API 参考

- **Base URL**: `http://export.arxiv.org/api/query`
- **响应格式**: Atom XML
- **常用分类**: cs.AI, cs.LG, cs.CV, cs.CL, stat.ML
