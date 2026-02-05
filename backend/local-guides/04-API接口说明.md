# API 接口说明

> **目标**: PaperTok 后端 REST API 完整文档

---

## 1. 基础信息

| 项目 | 值 |
|------|---|
| Base URL | `http://localhost:8080` |
| 协议 | HTTP/HTTPS |
| 数据格式 | JSON |

---

## 2. 通用响应格式

### 成功响应

```json
{
  "success": true,
  "data": { ... },
  "timestamp": 1706123456
}
```

### 错误响应

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述",
    "details": "详细信息（可选）"
  },
  "timestamp": 1706123456
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| `INVALID_PARAMS` | 参数无效 |
| `NOT_FOUND` | 资源不存在 |
| `INTERNAL_ERROR` | 服务器内部错误 |

---

## 3. 接口列表

### 3.1 健康检查

**GET /health**

检查服务运行状态。

**请求示例**：
```bash
curl http://localhost:8080/health
```

**响应示例**：
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

---

### 3.2 获取论文列表

**GET /api/v1/papers**

获取论文推荐流。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `category` | string | 否 | `cs.AI` | arXiv 分类 |
| `limit` | int | 否 | `20` | 返回数量（1-100） |
| `offset` | int | 否 | `0` | 分页偏移量 |
| `sort_by` | string | 否 | `lastUpdatedDate` | 排序方式 |

**排序方式**：
- `lastUpdatedDate` - 按更新时间
- `submittedDate` - 按提交时间

**常用分类**：
- `cs.AI` - 人工智能
- `cs.LG` - 机器学习
- `cs.CV` - 计算机视觉
- `cs.CL` - 计算语言学
- `cs.NE` - 神经网络
- `stat.ML` - 统计学习

**请求示例**：
```bash
curl "http://localhost:8080/api/v1/papers?category=cs.AI&limit=10"
```

**响应示例**：
```json
{
  "success": true,
  "data": {
    "papers": [
      {
        "id": "2401.12345",
        "title": "A Novel Approach to AI",
        "authors": ["John Doe", "Jane Smith"],
        "summary": "This paper presents...",
        "published": "2024-01-15T10:00:00Z",
        "updated": "2024-01-16T10:00:00Z",
        "categories": ["cs.AI", "cs.LG"],
        "primaryCategory": "cs.AI",
        "arxivUrl": "https://arxiv.org/abs/2401.12345",
        "pdfUrl": "https://arxiv.org/pdf/2401.12345.pdf",
        "imageUrl": "https://arxiv.org/html/2401.12345/x1.png"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 10
  },
  "timestamp": 1706123456
}
```

---

### 3.3 搜索论文

**GET /api/v1/papers/search**

按关键词搜索论文。

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `query` | string | 是 | - | 搜索关键词 |
| `limit` | int | 否 | `20` | 返回数量（1-100） |

**请求示例**：
```bash
curl "http://localhost:8080/api/v1/papers/search?query=transformer&limit=5"
```

**响应示例**：
```json
{
  "success": true,
  "data": {
    "papers": [...],
    "total": 5,
    "page": 1,
    "pageSize": 5
  },
  "timestamp": 1706123456
}
```

**错误响应**（缺少 query 参数）：
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMS",
    "message": "Query parameter 'query' is required"
  },
  "timestamp": 1706123456
}
```

---

### 3.4 获取论文详情

**GET /api/v1/papers/:id**

获取单篇论文详情。

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string | 论文 ID（如 `2401.12345`） |

**请求示例**：
```bash
curl http://localhost:8080/api/v1/papers/2401.12345
```

**响应示例**：
```json
{
  "success": true,
  "data": {
    "id": "2401.12345",
    "title": "A Novel Approach to AI",
    "authors": ["John Doe", "Jane Smith"],
    "summary": "This paper presents...",
    "published": "2024-01-15T10:00:00Z",
    "updated": "2024-01-16T10:00:00Z",
    "categories": ["cs.AI", "cs.LG"],
    "primaryCategory": "cs.AI",
    "arxivUrl": "https://arxiv.org/abs/2401.12345",
    "pdfUrl": "https://arxiv.org/pdf/2401.12345.pdf",
    "imageUrl": "https://arxiv.org/html/2401.12345/x1.png"
  },
  "timestamp": 1706123456
}
```

**错误响应**（论文不存在）：
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Paper not found"
  },
  "timestamp": 1706123456
}
```

---

## 4. Paper 对象

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 论文 ID |
| `title` | string | 标题 |
| `authors` | string[] | 作者列表 |
| `summary` | string | 摘要 |
| `published` | string | 发布时间（ISO 8601） |
| `updated` | string | 更新时间（ISO 8601） |
| `categories` | string[] | 分类列表 |
| `primaryCategory` | string | 主分类 |
| `arxivUrl` | string | arXiv 页面链接 |
| `pdfUrl` | string | PDF 下载链接 |
| `imageUrl` | string | 封面图链接 |

---

## 5. 测试脚本

### 5.1 获取 AI 论文

```bash
curl -s "http://localhost:8080/api/v1/papers?category=cs.AI&limit=5" | jq
```

### 5.2 搜索 GPT 论文

```bash
curl -s "http://localhost:8080/api/v1/papers/search?query=GPT&limit=5" | jq
```

### 5.3 获取论文详情

```bash
curl -s "http://localhost:8080/api/v1/papers/2401.12345" | jq
```

---

## 6. 相关文档

- [01-项目架构总览.md](./01-项目架构总览.md) - 架构概览
- [02-快速开始.md](./02-快速开始.md) - 快速开始
