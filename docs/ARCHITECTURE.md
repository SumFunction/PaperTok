# PaperTok 技术架构设计文档

## 1. 系统架构总览

### 1.1 架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户浏览器                               │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │              React 前端应用 (TypeScript)                   │ │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────────────────┐ │ │
│  │  │ 推荐流    │  │ 搜索      │  │ 收藏列表              │ │ │
│  │  └───────────┘  └───────────┘  └───────────────────────┘ │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST API
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Go 后端服务                                   │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                    API Layer (Handlers)                    │ │
│  │           /api/v1/papers, /api/v1/papers/search            │ │
│  ├───────────────────────────────────────────────────────────┤ │
│  │                      Facade Layer                          │ │
│  │                    统一业务入口                             │ │
│  ├───────────────────────────────────────────────────────────┤ │
│  │     Features (垂直切片)       │      Core Services         │ │
│  │  ├─ paperfeed (论文流)        │   └─ arxiv (arXiv客户端)   │ │
│  │  └─ papersearch (搜索)        │                            │ │
│  ├───────────────────────────────────────────────────────────┤ │
│  │                    Repository Layer                        │ │
│  │                  paper (论文数据访问)                       │ │
│  ├───────────────────────────────────────────────────────────┤ │
│  │                  Infrastructure Layer                      │ │
│  │               cache (缓存) │ httpclient (HTTP)             │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP 请求
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        arXiv API                                 │
│               http://export.arxiv.org/api/query                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 技术栈清单

#### 前端技术栈
- **框架**：React 18.2+ (TypeScript)
- **构建工具**：Vite 5.0+
- **状态管理**：React Hooks (useState, useEffect, useContext)
- **路由**：React Router v6
- **HTTP 客户端**：fetch API
- **本地存储**：localStorage

#### 后端技术栈
- **语言**：Go 1.21+
- **Web 框架**：Gin
- **HTTP 客户端**：http.Client
- **配置管理**：viper
- **缓存**：内存缓存

---

## 2. 后端架构设计

### 2.1 架构模式：Vertical Slice Architecture

采用**垂直切片架构**，按业务功能垂直划分代码：

```
┌─────────────┐ ┌─────────────┐
│  PaperFeed  │ │ PaperSearch │
│   Feature   │ │   Feature   │
│      ↓      │ │      ↓      │
│   Service   │ │   Service   │
└──────┬──────┘ └──────┬──────┘
       │               │
       ▼               ▼
┌─────────────────────────────────┐
│    Shared Kernel (Core + Infra) │
└─────────────────────────────────┘
```

### 2.2 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── api/                     # API 层
│   │   ├── handlers/            # HTTP 处理器
│   │   │   ├── paper.go
│   │   │   └── health.go
│   │   └── middleware/          # 中间件
│   │       ├── cors.go
│   │       └── logger.go
│   │
│   ├── facade/                  # Facade 模式入口
│   │   └── service.go
│   │
│   ├── features/                # Feature Slices（垂直切片）
│   │   ├── paperfeed/           # 论文推荐流功能
│   │   │   ├── interface.go
│   │   │   ├── deps.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── papersearch/         # 论文搜索功能
│   │       ├── interface.go
│   │       ├── deps.go
│   │       ├── service.go
│   │       └── service_test.go
│   │
│   ├── core/                    # Core Services（共享服务）
│   │   └── arxiv/               # arXiv API 客户端
│   │       ├── interface.go
│   │       ├── deps.go
│   │       ├── types.go
│   │       ├── errors.go
│   │       ├── client.go
│   │       └── client_test.go
│   │
│   ├── repository/              # Repository 层
│   │   └── paper/
│   │       ├── interface.go
│   │       └── memory.go
│   │
│   ├── infra/                   # 基础设施
│   │   ├── cache/
│   │   │   ├── interface.go
│   │   │   └── memory.go
│   │   └── httpclient/
│   │       └── client.go
│   │
│   └── config/
│       └── config.go
│
├── local-guides/                # 本地开发指南
│   ├── 00-本地指南索引.md
│   ├── 01-项目架构总览.md
│   ├── 02-快速开始.md
│   ├── 03-后端架构详解.md
│   └── 04-API接口说明.md
│
└── config.yaml
```

### 2.3 分层职责

| 层 | 职责 | 依赖 |
|---|------|------|
| **API Layer** | HTTP 请求处理、参数校验 | Facade |
| **Facade Layer** | 统一业务入口、依赖组装 | Features, Core |
| **Features** | 垂直切片、业务功能 | Core, Repository |
| **Core Services** | 共享服务、外部 API | Infra |
| **Repository** | 数据访问抽象 | Infra |
| **Infrastructure** | 基础设施（缓存、HTTP） | - |

### 2.4 依赖方向

```
API Layer → Facade → Features → Core Services
                 ↘          ↘
               Repository   Infrastructure
```

**原则**：
- 依赖只能向下，禁止循环依赖
- Features 之间不能相互依赖
- Core Services 不能依赖 Features

---

## 3. API 接口设计

### 3.1 获取论文列表

```http
GET /api/v1/papers

Query Parameters:
  - category: string (可选) - arXiv 分类，如 "cs.AI"
  - limit: number (可选) - 返回数量，默认 20，最大 100
  - offset: number (可选) - 分页偏移量
  - sort_by: string (可选) - 排序方式

Response 200:
{
  "success": true,
  "data": {
    "papers": [...],
    "total": 20,
    "page": 1,
    "pageSize": 20
  },
  "timestamp": 1706123456
}
```

### 3.2 搜索论文

```http
GET /api/v1/papers/search

Query Parameters:
  - query: string (必需) - 搜索关键词
  - limit: number (可选) - 返回数量，默认 20

Response 200:
{
  "success": true,
  "data": {
    "papers": [...],
    "total": 10,
    "page": 1,
    "pageSize": 10
  },
  "timestamp": 1706123456
}
```

### 3.3 获取论文详情

```http
GET /api/v1/papers/:id

Response 200:
{
  "success": true,
  "data": {
    "id": "2401.12345",
    "title": "Paper Title",
    "authors": ["Author 1", "Author 2"],
    "summary": "Abstract...",
    ...
  },
  "timestamp": 1706123456
}
```

---

## 4. 数据模型

### 4.1 Paper

```go
type Paper struct {
    ID              string    `json:"id"`
    Title           string    `json:"title"`
    Authors         []string  `json:"authors"`
    Summary         string    `json:"summary"`
    Published       time.Time `json:"published"`
    Updated         time.Time `json:"updated"`
    Categories      []string  `json:"categories"`
    PrimaryCategory string    `json:"primaryCategory"`
    ArxivURL        string    `json:"arxivUrl"`
    PDFURL          string    `json:"pdfUrl"`
    ImageURL        string    `json:"imageUrl"`
}
```

---

## 5. 配置管理

```yaml
# config.yaml
server:
  port: 8080
  mode: debug

arxiv:
  base_url: "http://export.arxiv.org/api/query"
  timeout: 10s

cache:
  enabled: true
  ttl: 300s

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "http://localhost:3000"
```

---

## 6. 前端架构设计

### 6.1 目录结构

```
frontend/
├── src/
│   ├── components/          # 可复用组件
│   │   ├── PaperCard.tsx
│   │   ├── LoadingCard.tsx
│   │   ├── CategoryFilter.tsx
│   │   ├── Navigation.tsx
│   │   └── Toast.tsx
│   ├── pages/               # 页面组件
│   │   ├── Feed.tsx
│   │   ├── Favorites.tsx
│   │   └── Search.tsx
│   ├── hooks/               # 自定义 Hooks
│   │   ├── usePapers.ts
│   │   ├── useLocalStorage.ts
│   │   ├── useKeyboard.ts
│   │   └── useSearch.ts
│   ├── services/            # API 服务
│   │   └── api.ts
│   ├── types/               # TypeScript 类型
│   │   └── index.ts
│   ├── utils/               # 工具函数
│   │   ├── date.ts
│   │   ├── text.ts
│   │   └── storage.ts
│   ├── contexts/            # React Context
│   │   └── AppContext.tsx
│   ├── App.tsx
│   └── main.tsx
├── package.json
└── vite.config.ts
```

---

## 7. 开发规范

### 7.1 模块标准结构

每个模块目录应包含：

| 文件 | 用途 | 是否必须 |
|------|------|---------|
| `README.md` | 模块文档 | 推荐 |
| `interface.go` | 对外接口 | 必须 |
| `deps.go` | 依赖接口 | 必须 |
| `types.go` | 领域类型 | 可选 |
| `errors.go` | 错误定义 | 推荐 |
| `service.go` | 实现 | 必须 |
| `service_test.go` | 测试 | 必须 |

### 7.2 测试规范

- 使用 Mock 依赖进行单元测试
- AAA 模式：Arrange-Act-Assert
- 覆盖正常路径和错误路径

---

## 8. 部署架构

### 8.1 开发环境

```bash
# 前端
cd frontend && npm run dev  # http://localhost:5173

# 后端
cd backend && go run cmd/server/main.go  # http://localhost:8080
```

### 8.2 生产部署

```
┌─────────────────┐
│   Nginx/Caddy   │ ← 反向代理
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼─────┐
│ 前端   │ │ 后端   │
│ 静态   │ │ Go     │
│ 文件   │ │ :8080  │
└───────┘ └────────┘
```

---

## 9. 相关文档

- [DESIGN_PRINCIPLES.md](./DESIGN_PRINCIPLES.md) - 设计原则
- [local-guides/](../backend/local-guides/) - 本地开发指南
- [PRD.md](./PRD.md) - 产品需求文档

---

**文档版本**：v2.0
**最后更新**：2026-01-26
**状态**：架构重构完成
