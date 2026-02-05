# PaperFeed Feature

> 论文推荐流功能，支持按分类获取最新论文

---

## 职责

- 获取指定分类的论文列表
- 管理论文缓存
- 分页支持

---

## 接口

```go
type Service interface {
    GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error)
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
- `paper.Repository` - 论文数据存储

---

## 使用示例

```go
svc := paperfeed.New(arxivSvc, paperRepo, 5*time.Minute)

papers, err := svc.GetFeed(ctx, &paperfeed.FetchRequest{
    Category: "cs.AI",
    Limit:    20,
    Offset:   0,
    SortBy:   "lastUpdatedDate",
})
```

---

## 数据流

```
1. GetFeed() 被调用
   ↓
2. 检查 Repository 缓存
   ├─ 命中 → 返回缓存数据
   └─ 未命中 → 继续
   ↓
3. 调用 arxiv.FetchByCategory()
   ↓
4. 保存到 Repository 缓存
   ↓
5. 返回论文列表
```
