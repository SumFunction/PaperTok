# PaperSearch Feature

> 论文搜索功能，支持关键词搜索和按 ID 获取

---

## 职责

- 按关键词搜索论文
- 按 ID 获取单篇论文

---

## 接口

```go
type Service interface {
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
| `service.go` | 实现 |
| `service_test.go` | 单元测试 |

---

## 依赖

- `arxiv.Service` - arXiv API 客户端

---

## 使用示例

```go
svc := papersearch.New(arxivSvc)

// 搜索论文
papers, err := svc.Search(ctx, "machine learning", 20)

// 获取单篇论文
paper, err := svc.GetByID(ctx, "2401.12345")
```

---

## 数据流

```
Search:
1. Search() 被调用
   ↓
2. 调用 arxiv.Search()
   ↓
3. 返回论文列表

GetByID:
1. GetByID() 被调用
   ↓
2. 调用 arxiv.GetByID()
   ↓
3. 返回论文（或 nil）
```
