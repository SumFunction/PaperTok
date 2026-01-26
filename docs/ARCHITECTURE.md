# PaperTok 技术架构设计文档

## 1. 系统架构总览

### 1.1 架构图
```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                            │
│  ┌───────────────────────────────────────────────────────┐  │
│  │           React 前端应用 (TypeScript)                  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌───────────────┐  │  │
│  │  │ 推荐流组件   │  │ 交互组件    │  │ 收藏列表      │  │  │
│  │  └─────────────┘  └─────────────┘  └───────────────┘  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌───────────────┐  │  │
│  │  │ 状态管理    │  │ 本地存储    │  │ 路由管理      │  │  │
│  │  └─────────────┘  └─────────────┘  └───────────────┘  │  │
│  └───────────────────────┬───────────────────────────────┘  │
└──────────────────────────┼─────────────────────────────────┘
                           │ HTTP/REST API
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go 后端服务                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                    HTTP Server                         │  │
│  │           (Gin Framework + middleware)                 │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │                                    │
│  ┌───────────────────────┴───────────────────────────────┐  │
│  │                    Handler 层                          │  │
│  │  ┌──────────┐  ┌──────────┐  ┌────────────────────┐  │  │
│  │  │ Paper    │  │ Search   │  │  Health Check      │  │  │
│  │  │ Handler  │  │ Handler  │  │                    │  │  │
│  │  └──────────┘  └──────────┘  └────────────────────┘  │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │                                    │
│  ┌───────────────────────┴───────────────────────────────┐  │
│  │                    Service 层                          │  │
│  │  ┌────────────────┐  ┌────────────────┐               │  │
│  │  │ ArxivService   │  │ CacheService   │               │  │
│  │  │ - 查询论文     │  │ - 内存缓存     │               │  │
│  │  │ - 解析响应     │  │ - TTL 管理     │               │  │
│  │  └────────────────┘  └────────────────┘               │  │
│  └───────────────────────┬───────────────────────────────┘  │
└──────────────────────────┼─────────────────────────────────┘
                           │ HTTP 请求
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    arXiv API                                 │
│            http://export.arxiv.org/api/query                 │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 技术栈清单

#### 前端技术栈
- **框架**：React 18.2+ (TypeScript)
- **构建工具**：Vite 5.0+
- **状态管理**：React Hooks (useState, useEffect, useContext)
- **路由**：React Router v6
- **UI 组件**：原生 CSS + TailwindCSS（可选）
- **虚拟滚动**：react-window 或 react-virtuoso
- **HTTP 客户端**：axios 或 fetch API
- **本地存储**：localStorage / sessionStorage

#### 后端技术栈
- **语言**：Go 1.21+
- **Web 框架**：Gin
- **HTTP 客户端**：http.Client
- **配置管理**：viper
- **日志**：logrus 或 zap
- **CORS**：gin-contrib/cors

---

## 2. 前端架构设计

### 2.1 目录结构
```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── components/          # 可复用组件
│   │   ├── PaperCard.tsx    # 论文卡片组件
│   │   ├── LoadingCard.tsx  # 加载骨架屏
│   │   ├── CategoryFilter.tsx # 分类筛选器
│   │   ├── ActionButtons.tsx  # 操作按钮组（点赞/收藏/分享）
│   │   └── Toast.tsx        # 提示消息
│   ├── pages/               # 页面组件
│   │   ├── Feed.tsx         # 推荐流页面
│   │   ├── Favorites.tsx    # 收藏列表页面
│   │   └── Search.tsx       # 搜索页面（可选）
│   ├── hooks/               # 自定义 Hooks
│   │   ├── usePapers.ts     # 获取论文数据
│   │   ├── useLocalStorage.ts # 本地存储 Hook
│   │   └── useKeyboard.ts   # 键盘快捷键 Hook
│   ├── services/            # API 服务
│   │   └── api.ts           # API 客户端
│   ├── types/               # TypeScript 类型定义
│   │   └── paper.ts         # 论文相关类型
│   ├── utils/               # 工具函数
│   │   ├── date.ts          # 日期格式化
│   │   ├── text.ts          # 文本处理
│   │   └── keyboard.ts      # 键盘事件
│   ├── context/             # React Context
│   │   └── AppContext.tsx   # 全局状态管理
│   ├── App.tsx              # 应用根组件
│   ├── main.tsx             # 应用入口
│   └── index.css            # 全局样式
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js       # 可选
```

### 2.2 核心组件设计

#### 2.2.1 PaperCard 组件
```typescript
interface PaperCardProps {
  paper: Paper;
  onLike: (id: string) => void;
  onFavorite: (id: string) => void;
  onShare: (paper: Paper) => void;
  isLiked: boolean;
  isFavorite: boolean;
}

// 功能：
// - 渲染论文完整信息
// - 处理点赞/收藏/分享交互
// - 展开/折叠摘要
// - 跳转原文
```

#### 2.2.2 Feed 页面
```typescript
// 功能：
// - 虚拟滚动展示论文列表
// - 自动加载更多
// - 分类筛选
// - 键盘快捷键支持
// - 空状态和错误处理
```

#### 2.2.3 Favorites 页面
```typescript
// 功能：
// - 展示收藏的论文列表
// - 支持取消收藏
// - 按收藏时间排序
```

### 2.3 状态管理

使用 React Context 进行全局状态管理：
```typescript
interface AppState {
  papers: Paper[];              // 当前论文列表
  likedPapers: Set<string>;     // 点赞的论文 ID
  favoritePapers: Set<string>;  // 收藏的论文 ID
  selectedCategory: string;     // 当前分类
  loading: boolean;             // 加载状态
  error: string | null;         // 错误信息
  // Actions
  fetchPapers: (category: string) => Promise<void>;
  toggleLike: (id: string) => void;
  toggleFavorite: (id: string) => void;
}
```

### 2.4 本地存储策略

存储键名规范：
```typescript
const STORAGE_KEYS = {
  LIKED_PAPERS: 'papertok_liked',
  FAVORITE_PAPERS: 'papertok_favorites',
  VIEWED_PAPERS: 'papertok_viewed',
  PREFERENCES: 'papertok_preferences',
};
```

数据结构：
```json
{
  "liked": ["2401.12345", "2401.67890"],
  "favorites": ["2401.12345"],
  "preferences": {
    "category": "cs.AI",
    "timeRange": "week"
  }
}
```

---

## 3. 后端架构设计

### 3.1 目录结构
```
backend/
├── cmd/
│   └── server/
│       └── main.go            # 应用入口
├── internal/
│   ├── handler/               # HTTP 处理器
│   │   ├── paper.go           # 论文相关接口
│   │   └── health.go          # 健康检查
│   ├── service/               # 业务逻辑层
│   │   ├── arxiv.go           # arXiv API 服务
│   │   └── cache.go           # 缓存服务
│   ├── model/                 # 数据模型
│   │   └── paper.go           # 论文数据结构
│   ├── config/                # 配置管理
│   │   └── config.go          # 配置结构
│   └── middleware/            # 中间件
│       ├── cors.go            # CORS 处理
│       └── logger.go          # 日志中间件
├── pkg/
│   └── arxiv/                 # arXiv API 客户端
│       ├── client.go          # API 客户端
│       └── parser.go          # Atom 响应解析
├── go.mod
├── go.sum
└── config.yaml                # 配置文件
```

### 3.2 API 接口设计

#### 3.2.1 获取论文列表
```http
GET /api/v1/papers

Query Parameters:
  - category: string (可选)  - arXiv 分类，如 "cs.AI"
  - max_results: number (可选) - 返回数量，默认 20，最大 100
  - sort_by: string (可选) - 排序方式，"lastUpdatedDate" 或 "submittedDate"

Response 200:
{
  "success": true,
  "data": [
    {
      "id": "2401.12345",
      "title": "Paper Title",
      "authors": ["Author 1", "Author 2"],
      "summary": "Abstract...",
      "published": "2024-01-15T10:00:00Z",
      "updated": "2024-01-15T10:00:00Z",
      "categories": ["cs.AI", "cs.LG"],
      "primaryCategory": "cs.AI",
      "arxivUrl": "https://arxiv.org/abs/2401.12345",
      "pdfUrl": "https://arxiv.org/pdf/2401.12345.pdf"
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 20
}

Response 500:
{
  "success": false,
  "error": "Failed to fetch papers from arXiv"
}
```

#### 3.2.2 搜索论文
```http
GET /api/v1/papers/search

Query Parameters:
  - query: string (必需) - 搜索关键词
  - max_results: number (可选) - 返回数量，默认 20

Response 200: (同上)
```

#### 3.2.3 健康检查
```http
GET /health

Response 200:
{
  "status": "ok",
  "version": "1.0.0"
}
```

### 3.3 核心模块设计

#### 3.3.1 arXiv Service
```go
type ArxivService interface {
    FetchPapers(ctx context.Context, req *FetchRequest) ([]*Paper, error)
    SearchPapers(ctx context.Context, query string, limit int) ([]*Paper, error)
}

type FetchRequest struct {
    Category    string
    MaxResults  int
    SortBy      string
}
```

#### 3.3.2 Cache Service
```go
type CacheService interface {
    Get(key string) ([]*Paper, bool)
    Set(key string, papers []*Paper, ttl time.Duration)
    Invalidate(category string)
}
```

#### 3.3.3 Handler
```go
type PaperHandler struct {
    arxivService service.ArxivService
    cacheService service.CacheService
}

func (h *PaperHandler) GetPapers(c *gin.Context)
func (h *PaperHandler) SearchPapers(c *gin.Context)
```

### 3.4 数据模型

```go
type Paper struct {
    ID               string    `json:"id"`
    Title            string    `json:"title"`
    Authors          []string  `json:"authors"`
    Summary          string    `json:"summary"`
    Published        time.Time `json:"published"`
    Updated          time.Time `json:"updated"`
    Categories       []string  `json:"categories"`
    PrimaryCategory  string    `json:"primaryCategory"`
    ArxivURL         string    `json:"arxivUrl"`
    PDFURL           string    `json:"pdfUrl"`
}
```

### 3.5 配置管理

```yaml
# config.yaml
server:
  port: 8080
  mode: debug  # debug, release

arxiv:
  base_url: "http://export.arxiv.org/api/query"
  timeout: 10s
  max_retries: 3

cache:
  enabled: true
  ttl: 300s  # 5 分钟

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "http://localhost:3000"
```

---

## 4. 数据流设计

### 4.1 论文获取流程
```
1. 用户切换分类/首次加载
   ↓
2. 前端调用 GET /api/v1/papers?category=cs.AI
   ↓
3. 后端检查缓存
   ├─ 缓存命中 → 返回缓存数据
   └─ 缓存未命中 → 继续
   ↓
4. 后端调用 arXiv API
   ↓
5. 解析 Atom 格式响应
   ↓
6. 数据转换为 Paper 模型
   ↓
7. 写入缓存
   ↓
8. 返回 JSON 给前端
   ↓
9. 前端更新状态并渲染
```

### 4.2 用户交互流程
```
用户点击点赞/收藏
   ↓
前端更新本地状态（乐观更新）
   ↓
同步到 localStorage
   ↓
更新 UI 反馈
```

---

## 5. 性能优化策略

### 5.1 前端优化
1. **虚拟滚动**：只渲染可见区域的论文卡片
2. **图片懒加载**：使用 Intersection Observer API
3. **预加载**：提前加载下一批数据
4. **防抖/节流**：搜索输入、滚动事件使用节流
5. **代码分割**：使用 React.lazy 和 Suspense

### 5.2 后端优化
1. **内存缓存**：缓存 arXiv API 响应，减少外部请求
2. **连接池**：复用 HTTP 连接
3. **超时控制**：避免慢请求阻塞
4. **并发限制**：使用 goroutine 池限制并发

### 5.3 网络优化
1. **压缩**：启用 gzip 压缩
2. **CDN**：静态资源使用 CDN 加速
3. **HTTP/2**：后端启用 HTTP/2

---

## 6. 错误处理与容错

### 6.1 前端错误处理
```typescript
// 网络错误
- 重试机制（3 次）
- 离线提示
- 骨架屏降级

// 数据错误
- 空状态提示
- 默认占位图
- 错误边界（Error Boundary）
```

### 6.2 后端错误处理
```go
// 重试策略
- 指数退避重试
- 最大重试次数限制

// 超时处理
- 请求超时 10s
- 连接超时 5s

// 降级策略
- 缓存过期时返回旧数据
- arXiv API 失败时返回缓存
```

---

## 7. 安全考虑

### 7.1 前端安全
- **XSS 防护**：React 自动转义，避免 dangerouslySetInnerHTML
- **CSRF**：MVP 阶段无用户登录，暂不需要
- **内容安全策略**：配置 CSP 头

### 7.2 后端安全
- **CORS 配置**：仅允许信任的域名
- **速率限制**：防止 API 滥用（可选）
- **输入验证**：验证所有查询参数
- **日志脱敏**：不记录敏感信息

---

## 8. 部署架构

### 8.1 开发环境
```bash
# 前端
cd frontend && npm run dev  # http://localhost:5173

# 后端
cd backend && go run cmd/server/main.go  # http://localhost:8080
```

### 8.2 生产部署（MVP 简化方案）
```
┌─────────────────┐
│   Nginx / Caddy │ ← 反向代理 + 静态文件服务
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼─────┐
│ 前端   │ │ 后端   │
│ 静态   │ │ Go 进程│
│ 文件   │ │ :8080  │
└───────┘ └────────┘
```

### 8.3 Docker 部署（可选）
```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - ENV=production

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
```

---

## 9. 监控与日志

### 9.1 前端监控
- 错误收集：window.onerror + React Error Boundary
- 性能监控：Performance API
- 用户行为：关键操作埋点（可选）

### 9.2 后端日志
```go
// 结构化日志
logger.WithFields(logrus.Fields{
    "endpoint": "/api/v1/papers",
    "category": "cs.AI",
    "duration": "234ms",
}).Info("Fetched papers successfully")
```

---

## 10. 测试策略

### 10.1 前端测试
- **单元测试**：工具函数、自定义 Hooks
- **组件测试**：核心组件（PaperCard、Feed）
- **E2E 测试**：关键用户流程（可选）

### 10.2 后端测试
- **单元测试**：Service 层逻辑
- **集成测试**：API 端点测试
- **Mock 测试**：Mock arXiv API 响应

---

## 11. 开发规范

### 11.1 Git 提交规范
```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具链相关
```

### 11.2 代码规范
- **Go**：遵循 gofmt、golint、go vet
- **TypeScript**：使用 ESLint + Prettier
- **注释**：关键逻辑添加中文注释

---

## 12. MVP 后续扩展

### 12.1 数据库扩展
- 引入 MySQL 存储用户数据
- 引入 Redis 做分布式缓存
- 引入 RocketMQ 做异步任务

### 12.2 功能扩展
- 用户认证系统
- 推荐算法优化
- 社交功能
- 论文笔记功能

---

**文档版本**：v1.0
**最后更新**：2025-01-22
**状态**：架构设计完成
