# PaperTok 后端变更记录

## 2025-01-22

### 新增功能
- **图片 URL 提取**: 为论文数据添加 `imageUrl` 字段，从 arXiv HTML 版本提取封面图
  - 格式：`https://arxiv.org/html/{paper_id}/x1.png`
  - 在 `convertFeedToPapers` 中自动生成

### 优化改进
- **API 响应格式统一**
  - 使用标准的 `ErrorInfo` 结构体返回错误信息
  - 添加 Unix 时间戳字段
  - 完善错误处理的 HTTP 状态码

- **分页参数优化**
  - 将 `max_results` 改为更语义化的 `limit`
  - 添加 `offset` 参数支持分页
  - 添加默认值和边界验证

### 新增接口
- `GET /api/v1/papers/:id` - 获取单篇论文详情
  - 支持通过 arXiv ID 查询论文
  - 返回完整的论文信息

### 修改文件
- `internal/model/paper.go` - 已包含 ImageUrl 字段
- `internal/service/arxiv.go` - 添加图片 URL 生成逻辑和 GetPaperByID 实现
- `internal/handler/paper.go` - 完善响应格式，添加详情接口
- `cmd/server/main.go` - 注册新的路由

### 验证方式
```bash
# 启动服务
cd backend && go run cmd/server/main.go

# 测试健康检查
curl http://localhost:8080/health

# 测试论文列表
curl "http://localhost:8080/api/v1/papers?limit=2"

# 测试论文详情
curl "http://localhost:8080/api/v1/papers/2401.12345"

# 测试搜索
curl "http://localhost:8080/api/v1/papers/search?query=transformer&limit=5"
```
