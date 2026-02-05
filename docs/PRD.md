# PaperTok 产品需求文档（PRD）

> **文档版本**: v2.1  
> **最后更新**: 2026-01-27  
> **维护者**: 产品团队  
> **状态**: 持续迭代中

---

## 📋 目录

1. [产品概述](#1-产品概述)
2. [功能清单](#2-功能清单)
3. [迭代计划](#3-迭代计划)
4. [技术架构映射](#4-技术架构映射)
5. [更新日志](#5-更新日志)

---

## 1. 产品概述

### 1.1 产品定位

PaperTok 是一个网页版"抖音式刷论文"应用，通过沉浸式的卡片流体验，让用户能够快速浏览和发现 arXiv 上的最新学术论文。

### 1.2 核心价值

- **高效浏览**：通过短视频式的滑动体验，快速筛选感兴趣的论文
- **降低认知负担**：精美的视觉呈现，让阅读论文变得轻松愉悦
- **个性化发现**：基于分类和兴趣筛选，发现相关领域最新研究

### 1.3 目标用户

- 科研人员和研究生
- AI/CS 领域从业者
- 学术爱好者

---

## 2. 功能清单

### 2.1 已实现功能 ✅

#### 后端功能

| 功能ID | 功能名称 | 状态 | 实现位置 | 完成时间 |
|--------|---------|------|---------|---------|
| BE-001 | 论文推荐流 API | ✅ 已完成 | `features/paperfeed` | 2026-01-26 |
| BE-002 | 论文搜索 API | ✅ 已完成 | `features/papersearch` | 2026-01-26 |
| BE-003 | 论文详情 API | ✅ 已完成 | `features/papersearch` | 2026-01-26 |
| BE-004 | arXiv API 集成 | ✅ 已完成 | `core/arxiv` | 2026-01-26 |
| BE-005 | 内存缓存 | ✅ 已完成 | `infra/cache` | 2026-01-26 |
| BE-006 | 健康检查接口 | ✅ 已完成 | `api/handlers/health` | 2026-01-26 |
| BE-007 | 用户注册 API | ✅ 已完成 | `features/userauth` | 2026-01-27 |
| BE-008 | 用户登录 API | ✅ 已完成 | `features/userauth` | 2026-01-27 |
| BE-009 | JWT Token 认证 | ✅ 已完成 | `core/auth` | 2026-01-27 |
| BE-010 | 用户信息 API | ✅ 已完成 | `features/userauth` | 2026-01-27 |
| BE-011 | 认证中间件 | ✅ 已完成 | `api/middleware/auth` | 2026-01-27 |
| BE-012 | 用户仓储层 | ✅ 已完成 | `repository/user` | 2026-01-27 |

#### 前端功能

| 功能ID | 功能名称 | 状态 | 实现位置 | 完成时间 |
|--------|---------|------|---------|---------|
| FE-001 | 论文推荐流页面 | ✅ 已完成 | `pages/Feed.tsx` | 2025-01-22 |
| FE-002 | 收藏列表页面 | ✅ 已完成 | `pages/Favorites.tsx` | 2025-01-22 |
| FE-003 | 搜索页面 | ✅ 已完成 | `pages/Search.tsx` | 2025-01-22 |
| FE-004 | 论文卡片组件 | ✅ 已完成 | `components/PaperCard.tsx` | 2025-01-22 |
| FE-005 | 分类筛选器 | ✅ 已完成 | `components/CategoryFilter.tsx` | 2025-01-22 |
| FE-006 | 点赞功能 | ✅ 已完成 | `contexts/AppContext.tsx` | 2025-01-22 |
| FE-007 | 收藏功能 | ✅ 已完成 | `contexts/AppContext.tsx` | 2025-01-22 |
| FE-008 | 键盘快捷键 | ✅ 已完成 | `hooks/useKeyboard.ts` | 2025-01-22 |
| FE-009 | 无限滚动加载 | ✅ 已完成 | `pages/Feed.tsx` | 2025-01-22 |
| FE-010 | 本地存储 | ✅ 已完成 | `utils/storage.ts` | 2025-01-22 |
| FE-011 | 登录页面 | ✅ 已完成 | `pages/Login.tsx` | 2026-01-27 |
| FE-012 | 注册页面 | ✅ 已完成 | `pages/Register.tsx` | 2026-01-27 |
| FE-013 | 认证模态框 | ✅ 已完成 | `components/AuthModal.tsx` | 2026-01-27 |
| FE-014 | 登录表单组件 | ✅ 已完成 | `components/LoginForm.tsx` | 2026-01-27 |
| FE-015 | 注册表单组件 | ✅ 已完成 | `components/RegisterForm.tsx` | 2026-01-27 |
| FE-016 | 认证状态管理 | ✅ 已完成 | `contexts/AuthContext.tsx` | 2026-01-27 |
| FE-017 | 路由保护组件 | ✅ 已完成 | `components/ProtectedRoute.tsx` | 2026-01-27 |

### 2.2 待实现功能 🔄

#### 高优先级（P0）

| 功能ID | 功能名称 | 优先级 | 预计工作量 | 依赖 |
|--------|---------|--------|-----------|------|
| FE-020 | 收藏列表完整显示 | P0 | 2天 | BE-003 |
| FE-021 | 分享功能完善 | P0 | 1天 | - |
| FE-022 | 错误重试机制 | P0 | 1天 | - |
| BE-013 | 批量获取论文接口 | P0 | 2天 | - |

#### 中优先级（P1）

| 功能ID | 功能名称 | 优先级 | 预计工作量 | 依赖 |
|--------|---------|--------|-----------|------|
| FE-023 | 论文详情页 | P1 | 3天 | BE-003 |
| FE-024 | 阅读历史记录 | P1 | 2天 | FE-010 |
| FE-025 | 搜索历史 | P1 | 1天 | FE-010 |
| FE-026 | 收藏云同步 | P1 | 3天 | FE-016 |
| BE-014 | Redis 缓存 | P1 | 3天 | - |
| BE-015 | 请求限流 | P1 | 2天 | - |

#### 低优先级（P2）

| 功能ID | 功能名称 | 优先级 | 预计工作量 | 依赖 |
|--------|---------|--------|-----------|------|
| FE-027 | 暗色/亮色模式切换 | P2 | 2天 | - |
| FE-028 | 多语言支持 | P2 | 5天 | - |
| FE-029 | PDF 预览 | P2 | 5天 | - |
| BE-016 | 推荐算法 | P2 | 10天 | BE-014 |

### 2.3 未来规划 💡

| 功能ID | 功能名称 | 优先级 | 说明 |
|--------|---------|--------|------|
| FE-030 | 论文笔记功能 | P3 | 支持高亮和注释 |
| FE-031 | 社交分享 | P3 | 分享到 Twitter/LinkedIn |
| FE-032 | 第三方登录 | P3 | 微信/Google OAuth |
| BE-017 | 推荐算法优化 | P3 | 基于浏览历史 |
| BE-018 | 引文关系图 | P3 | 可视化论文引用关系 |

---

## 3. 迭代计划

### 3.1 当前迭代：v2.1（2026-01-26 ~ 2026-02-09）

**目标**: 完善核心功能，提升用户体验

#### Sprint 1（Week 1）

- [ ] **FE-011**: 收藏列表完整显示
  - 后端：实现批量获取论文接口（BE-007）
  - 前端：完善收藏列表页面，显示完整论文信息
  - 测试：单元测试 + 集成测试

- [ ] **FE-012**: 分享功能完善
  - 实现复制到剪贴板功能
  - 添加分享成功提示
  - 支持分享到社交媒体（可选）

- [ ] **FE-013**: 错误重试机制
  - 网络错误自动重试（最多3次）
  - 友好的错误提示
  - 离线状态检测

#### Sprint 2（Week 2）

- [ ] **FE-014**: 论文详情页
  - 创建详情页路由和组件
  - 显示完整论文信息
  - 支持返回推荐流

- [ ] **BE-008**: Redis 缓存（可选）
  - 评估是否需要分布式缓存
  - 如需要，实现 Redis 缓存层

### 3.2 下一迭代：v2.2（2026-02-10 ~ 2026-02-23）

**目标**: 性能优化和体验增强

- [ ] **FE-015**: 阅读历史记录
- [ ] **FE-016**: 搜索历史
- [ ] **BE-009**: 请求限流
- [ ] **FE-017**: 暗色模式

### 3.3 迭代节奏

- **迭代周期**: 2周一个迭代
- **发布频率**: 每2周发布一次
- **紧急修复**: 随时发布 hotfix

---

## 4. 技术架构映射

### 4.1 功能到架构的映射

| 功能 | 前端位置 | 后端位置 | API 端点 |
|------|---------|---------|---------|
| 论文推荐流 | `pages/Feed.tsx` | `features/paperfeed` | `GET /api/v1/papers` |
| 论文搜索 | `pages/Search.tsx` | `features/papersearch` | `GET /api/v1/papers/search` |
| 论文详情 | - | `features/papersearch` | `GET /api/v1/papers/:id` |
| 收藏列表 | `pages/Favorites.tsx` | - | - |
| 点赞/收藏 | `contexts/AppContext.tsx` | - | 本地存储 |
| 用户注册 | `components/RegisterForm.tsx` | `features/userauth` | `POST /api/v1/auth/register` |
| 用户登录 | `components/LoginForm.tsx` | `features/userauth` | `POST /api/v1/auth/login` |
| 用户信息 | `contexts/AuthContext.tsx` | `features/userauth` | `GET /api/v1/auth/profile` |
| Token 刷新 | `services/api.ts` | `features/userauth` | `POST /api/v1/auth/refresh` |

### 4.2 新增功能的开发流程

1. **需求分析**: 在 PRD 中添加功能需求
2. **架构设计**: 确定功能归属的 Feature/Core Service
3. **后端开发**: 
   - 创建/扩展 Feature 或 Core Service
   - 实现接口和测试
   - 在 Facade 中注册
4. **前端开发**:
   - 创建/扩展页面或组件
   - 集成 API
   - 添加交互逻辑
5. **测试**: 单元测试 + 集成测试
6. **文档更新**: 更新 PRD 状态和架构文档

### 4.3 功能开发模板

#### 后端 Feature 开发

```bash
# 1. 创建 Feature 目录
mkdir -p internal/features/newfeature

# 2. 创建标准文件
touch internal/features/newfeature/{interface,deps,service,service_test}.go

# 3. 实现功能

# 4. 在 Facade 中注册
# 编辑 internal/facade/service.go

# 5. 添加 API Handler
# 编辑 internal/api/handlers/xxx.go
```

#### 前端功能开发

```bash
# 1. 创建页面/组件
touch src/pages/NewPage.tsx
touch src/components/NewComponent.tsx

# 2. 添加路由（如需要）
# 编辑 src/App.tsx

# 3. 集成 API
# 编辑 src/services/api.ts
```

---

## 5. 更新日志

### v2.1 (2026-01-27)

**用户认证系统**

- ✅ 实现完整的用户注册/登录功能
- ✅ JWT Token 认证机制
- ✅ 密码强度验证（8位+，2种字符类型）
- ✅ 认证中间件保护路由
- ✅ 用户信息 Repository 层（支持内存/MySQL）

**前端认证UI（抖音风格）**

- ✅ 登录页面（抖音风格动态背景）
- ✅ 注册页面（渐变动画效果）
- ✅ 认证模态框（Tab切换、滑动动画）
- ✅ 登录表单（密码可见性切换）
- ✅ 注册表单（详细密码强度检查器）
- ✅ 认证状态管理（AuthContext）
- ✅ 路由保护组件

**优化改进**

- ✅ 后端错误信息中文化
- ✅ 修复 AuthModal useState 导入问题
- ✅ 完善密码策略提示

### v2.0 (2026-01-26)

**架构重构**

- ✅ 重构后端为垂直切片架构（Vertical Slice Architecture）
- ✅ 创建 Facade 层统一业务入口
- ✅ 实现 paperfeed 和 papersearch Features
- ✅ 重构 arXiv 客户端为 Core Service
- ✅ 创建 Repository 层和 Infrastructure 层
- ✅ 完善文档体系（local-guides、README）

**功能状态**

- ✅ 后端 API 功能完整
- ✅ 前端核心功能完整
- ⚠️ 收藏列表需要完善（FE-020）

### v1.0 (2025-01-22)

**MVP 发布**

- ✅ 论文推荐流
- ✅ 分类筛选
- ✅ 点赞/收藏
- ✅ 搜索功能
- ✅ 收藏列表（基础版）

---

## 6. 功能详细说明

### 6.1 论文推荐流（FE-001）

**功能描述**: 抖音式垂直滑动浏览论文

**交互行为**:
- 垂直滚动浏览
- 自动加载下一批（无限滚动）
- 键盘快捷键支持（j/k 切换）
- 分类筛选

**技术要求**:
- 虚拟滚动优化（如需要）
- 图片懒加载
- 预加载策略

**API**: `GET /api/v1/papers?category={category}&limit={limit}&offset={offset}`

### 6.2 论文搜索（FE-003）

**功能描述**: 关键词搜索论文

**交互行为**:
- 输入关键词搜索
- 实时搜索建议（可选）
- 搜索结果列表展示

**API**: `GET /api/v1/papers/search?query={query}&limit={limit}`

### 6.3 收藏功能（FE-007）

**功能描述**: 收藏感兴趣的论文

**当前实现**:
- ✅ 本地存储收藏状态
- ⚠️ 收藏列表仅显示 ID（待完善）

**待完善**:
- 显示完整论文信息
- 支持批量操作
- 支持分类管理

**依赖**: BE-007（批量获取论文接口）

### 6.4 点赞功能（FE-006）

**功能描述**: 快速标记喜欢的论文

**实现**:
- 本地存储点赞状态
- 视觉反馈（心形图标）
- 键盘快捷键（l）

---

## 7. 非功能需求

### 7.1 性能要求

| 指标 | 目标 | 当前状态 |
|------|------|---------|
| 首屏加载 | < 2秒 | ✅ 达标 |
| API 响应 | < 5秒 | ✅ 达标 |
| 滑动流畅度 | 60fps | ✅ 达标 |
| 图片加载 | 懒加载 | ✅ 已实现 |

### 7.2 可用性要求

- ✅ 响应式设计（桌面 + 移动端）
- ✅ 错误处理和重试
- ⚠️ 离线支持（部分实现，待完善）

### 7.3 兼容性

- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+

---

## 8. 数据模型

### 8.1 Paper 对象

```typescript
interface Paper {
  id: string;                    // arXiv ID
  title: string;                 // 标题
  authors: string[];             // 作者列表
  summary: string;               // 摘要
  published: string;             // 发表时间 (ISO 8601)
  updated: string;               // 更新时间 (ISO 8601)
  categories: string[];          // 分类标签
  primaryCategory: string;      // 主分类
  arxivUrl: string;              // arXiv 链接
  pdfUrl: string;                // PDF 链接
  imageUrl: string;              // 封面图
}
```

### 8.2 用户偏好（本地存储）

```typescript
interface UserPreference {
  likedPapers: string[];         // 点赞的论文 ID
  favoritePapers: string[];      // 收藏的论文 ID
  viewedPapers: string[];         // 已浏览的论文 ID
  selectedCategory: string;      // 当前分类
  searchHistory: string[];       // 搜索历史（待实现）
}
```

---

## 9. API 规范

### 9.1 获取论文列表

```
GET /api/v1/papers

Query Parameters:
  - category: string (可选) - arXiv 分类
  - limit: number (可选) - 返回数量，默认 20
  - offset: number (可选) - 分页偏移量
  - sort_by: string (可选) - 排序方式

Response:
{
  "success": true,
  "data": {
    "papers": Paper[],
    "total": number,
    "page": number,
    "pageSize": number
  }
}
```

### 9.2 搜索论文

```
GET /api/v1/papers/search

Query Parameters:
  - query: string (必需) - 搜索关键词
  - limit: number (可选) - 返回数量

Response:
{
  "success": true,
  "data": {
    "papers": Paper[],
    "total": number
  }
}
```

### 9.3 获取论文详情

```
GET /api/v1/papers/:id

Response:
{
  "success": true,
  "data": Paper
}
```

### 9.4 批量获取论文（待实现）

```
POST /api/v1/papers/batch

Request Body:
{
  "ids": string[]
}

Response:
{
  "success": true,
  "data": {
    "papers": Paper[]
  }
}
```

### 9.5 用户注册

```
POST /api/v1/auth/register

Request Body:
{
  "username": string,    // 3-50字符，字母数字下划线
  "email": string,       // 有效邮箱格式
  "password": string     // 至少8位，需包含2种字符类型
}

Response:
{
  "success": true,
  "data": {
    "user": {
      "id": number,
      "username": string,
      "email": string,
      "createdAt": string,
      "updatedAt": string
    },
    "token": string    // JWT Token
  }
}
```

### 9.6 用户登录

```
POST /api/v1/auth/login

Request Body:
{
  "identifier": string,  // 邮箱或用户名
  "password": string
}

Response:
{
  "success": true,
  "data": {
    "user": User,
    "token": string
  }
}
```

### 9.7 获取用户信息

```
GET /api/v1/auth/profile

Headers:
  Authorization: Bearer {token}

Response:
{
  "success": true,
  "data": {
    "id": number,
    "username": string,
    "email": string,
    "createdAt": string,
    "updatedAt": string
  }
}
```

### 9.8 刷新 Token

```
POST /api/v1/auth/refresh

Request Body:
{
  "token": string  // 当前 token
}

Response:
{
  "success": true,
  "data": {
    "user": User,
    "token": string  // 新 token
  }
}
```

---

## 10. 开发规范

### 10.1 PRD 更新规范

**何时更新 PRD**:
- 新增功能时
- 功能状态变更时（开发中 → 已完成）
- 迭代计划调整时
- 架构变更影响功能时

**更新内容**:
1. 在"功能清单"中添加/更新功能
2. 在"更新日志"中记录变更
3. 更新"迭代计划"（如需要）
4. 更新"技术架构映射"（如需要）

### 10.2 功能开发流程

1. **需求确认**: 在 PRD 中明确功能需求
2. **技术设计**: 确定实现方案和架构位置
3. **开发**: 按照架构规范实现
4. **测试**: 编写并运行测试
5. **更新 PRD**: 将功能状态更新为"已完成"
6. **文档**: 更新相关技术文档

---

## 11. 成功指标

### 11.1 MVP 验收标准

- [x] 能够从 arXiv 成功获取论文数据
- [x] 推荐流滑动流畅，无卡顿
- [x] 点赞、收藏功能正常工作
- [x] 分类切换和搜索功能可用
- [ ] 收藏列表完整显示（进行中）
- [x] 在主流浏览器上正常运行

### 11.2 后续迭代指标

- 用户留存率（日活、周活）
- 平均浏览时长
- 收藏转化率
- 分享次数
- API 响应时间
- 错误率

---

## 12. 风险与依赖

### 12.1 外部依赖

- **arXiv API**: 公开可用，但有请求频率限制
- **浏览器兼容性**: 依赖现代浏览器 API

### 12.2 技术风险

- **API 稳定性**: arXiv API 可能波动，需要缓存和重试机制
- **性能瓶颈**: 大量图片加载可能影响性能
- **跨域问题**: 已通过后端代理解决

---

## 附录

### A. 功能ID命名规范

- **BE-XXX**: 后端功能
- **FE-XXX**: 前端功能
- **API-XXX**: API 相关功能

### B. 优先级说明

- **P0**: 高优先级，当前迭代必须完成
- **P1**: 中优先级，下一迭代计划
- **P2**: 低优先级，未来规划
- **P3**: 长期规划

### C. 状态说明

- **✅ 已完成**: 功能已实现并通过测试
- **🔄 进行中**: 功能正在开发
- **⏳ 待开始**: 功能已规划但未开始
- **⚠️ 待完善**: 功能已实现但需要改进

---

**文档维护**: 每次功能迭代后更新此文档  
**版本控制**: 使用 Git 跟踪 PRD 变更  
**评审周期**: 每个迭代开始前评审和更新
