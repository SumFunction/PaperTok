# PaperTok 项目改进方案（最终版）

> **文档版本**: v1.0
> **创建日期**: 2026-02-02
> **基于文档**: 项目调研报告、架构设计文档、产品需求文档、改进建议矩阵
> **状态**: 已批准，待执行

---

## 目录

1. [执行摘要](#1-执行摘要)
2. [当前架构评估](#2-当前架构评估)
3. [存在问题清单](#3-存在问题清单)
4. [优先级分类的改进建议](#4-优先级分类的改进建议)
5. [技术债务处理方案](#5-技术债务处理方案)
6. [功能演进路线图](#6-功能演进路线图)
7. [实施计划](#7-实施计划)
8. [风险管控](#8-风险管控)
9. [成功指标](#9-成功指标)

---

## 1. 执行摘要

### 1.1 项目现状概述

PaperTok 是一个网页版"抖音式刷论文"应用，采用现代化的全栈技术架构。经过全面调研，项目整体技术架构**优秀**，代码质量高，文档完善。

**核心优势**:
- ✅ 严格遵循垂直切片架构（Vertical Slice Architecture），模块化程度高
- ✅ 技术栈现代化：Go 1.24 + React 19 + TypeScript
- ✅ 后端测试覆盖率良好（5/5测试包通过）
- ✅ 完善的中文文档体系，与代码实现保持同步
- ✅ JWT认证系统已实现且配置安全（强制环境变量）

**核心问题**:
- ⚠️ 收藏功能不完整（仅显示ID列表）
- ⚠️ 前端测试缺失
- ⚠️ 性能优化空间大（内存缓存、无虚拟列表）
- ⚠️ 用户数据未持久化（仅localStorage）

### 1.2 改进目标

| 阶段 | 目标 | 时间周期 |
|------|------|---------|
| 第一阶段 | 安全与稳定性 | 2周 |
| 第二阶段 | 性能优化 | 3周 |
| 第三阶段 | 功能完善 | 4周 |
| 第四阶段 | 架构演进 | 持续进行 |

### 1.3 关键决策

1. **优先级排序**: 用户体验 > 性能 > 新功能
2. **渐进式改进**: 避免大规模重构���采用增量改进
3. **技术债务管理**: 每个迭代预留20%时间处理技术债务
4. **监控驱动**: 建立监控体系，数据驱动优化

---

## 2. 当前架构评估

### 2.1 技术栈评估

#### 后端技术栈（评分: ⭐⭐⭐⭐⭐ 5/5）

| 技术 | 版本 | 评估 | 说明 |
|------|------|------|------|
| Go | 1.24.0 | ✅ 优秀 | 最新稳定版，性能出色 |
| Gin | 1.11.0 | ✅ 优秀 | 轻量级Web框架，生态成熟 |
| JWT | golang-jwt/jwt/v5 | ✅ 优秀 | 标准实现，配置安全 |
| MySQL驱动 | go-sql-driver/mysql | ✅ 良好 | 成熟稳定 |
| Viper | 1.21.0 | ✅ 优秀 | 配置管理规范 |

#### 前端技术栈（评分: ⭐⭐⭐⭐⭐ 5/5）

| 技术 | 版本 | 评估 | 说明 |
|------|------|------|------|
| React | 19.2.0 | ✅ 优秀 | 最新稳定版，性能优秀 |
| TypeScript | 5.9.3 | ✅ 优秀 | 类型安全，开发体验好 |
| Vite | 7.2.4 | ✅ 优秀 | 快速构建，HMR高效 |
| Axios | 1.13.2 | ✅ 良好 | 成熟HTTP客户端 |
| React Router | 7.12.0 | ✅ 优秀 | 路由管理标准 |

### 2.2 架构设计评估

#### 分层架构（评分: ⭐⭐⭐⭐⭐ 5/5）

```
┌─────────────────────────────────────────────────────────────┐
│                      API Layer                               │
│                   (HTTP Handlers)                            │
├─────────────────────────────────────────────────────────────┤
│                     Facade Layer                             │
│               (Unified Entry Point)                          │
├─────────────────────────────────────────────────────────────┤
│    Feature Slices              │      Core Services          │
│  - paperfeed                   │    - arxiv                  │
│  - papersearch                 │    - auth                   │
│  - userauth                    │                             │
├─────────────────────────────────────────────────────────────┤
│                   Repository Layer                           │
│              (paper, user, memory/sql)                       │
├─────────────────────────────────────────────────────────────┤
│                Infrastructure Layer                          │
│                  (cache, httpclient)                         │
└─────────────────────────────────────────────────────────────┘
```

**优势**:
- 依赖方向清晰，无循环依赖
- 模块职责明确，高内聚低耦合
- 接口抽象良好，便于测试和扩展
- 配置驱动的存储选择（内存/MySQL）

### 2.3 代码质量评估

#### 后端代码（评分: ⭐⭐⭐⭐ 4/5）

| 指标 | 状态 | 说明 |
|------|------|------|
| 测试覆盖率 | ✅ 良好 | 5/5测试包通过，覆盖核心功能 |
| 错误处理 | ✅ 完善 | 自定义错误类型，错误传播规范 |
| 注释文档 | ✅ 详细 | 公共接口文档完整 |
| 代码规范 | ✅ 规范 | 遵循Go语言规范 |

**改进空间**:
- API Handler层测试缺失
- Repository层测试不足
- Infrastructure层测试缺失

#### 前端代码（评分: ⭐⭐⭐⭐ 4/5）

| 指标 | 状态 | 说明 |
|------|------|------|
| TypeScript覆盖 | ✅ 完整 | 严格类型检查 |
| 组件设计 | ✅ 良好 | 函数组件+Hooks |
| 状态管理 | ✅ 良好 | Context + Hooks |
| 代码规范 | ✅ 完善 | ESLint规则配置 |

**改进空间**:
- 缺乏自动化测试
- 部分组件复杂度较高
- 无统一错误处理层

### 2.4 文档完整性评估（评分: ⭐⭐⭐⭐⭐ 5/5）

| 文档 | 状态 | 说明 |
|------|------|------|
| PRD.md | ✅ 完整 | 产品需求、功能清单、API规范 |
| ARCHITECTURE.md | ✅ 完整 | 系统架构、技术栈、目录结构 |
| DESIGN_PRINCIPLES.md | ✅ 完整 | 设计原则、模块规范 |
| local-guides/ | ✅ 完整 | 4个本地开发指南 |
| API文档 | ✅ 完整 | 响应格式规范 |

---

## 3. 存在问题清单

### 3.1 高优先级问题（P0）

#### 问题P0-001: 收藏功能不完整

**问题描述**:
- 收藏列表页面仅显示论文ID列表
- 无法展示论文标题、作者、摘要等信息
- 用户体验极差，功能基本不可用

**影响范围**:
- 用户体验严重受损
- 违反产品核心价值（快速浏览和收藏）
- 用户留存率可能下降

**技术原因**:
- 前端: `Favorites.tsx` 仅读取localStorage中的ID列表
- 后端: 缺少批量获取论文接口 `POST /api/v1/papers/batch`
- 无后端数据同步机制

**解决方案**:
1. 后端实现批量获取论文接口
2. 前端调用批量接口获取完整信息
3. 实现收藏数据后端同步

**预计工作量**: 3天

---

#### 问题P0-002: 前端测试缺失

**问题描述**:
- 前端无任何自动化测试
- 组件测试、Hook测试、集成测试全部缺失
- 代码重构和修改风险高

**影响范围**:
- 代码质量无法保障
- Bug发现滞后
- 重构风险高

**解决方案**:
1. 搭建Jest + React Testing Library环境
2. 为核心组件编写单元测试
3. 为关键Hook编写测试
4. 建立测试覆盖率目标（>70%）

**预计工作量**: 4天

---

#### 问题P0-003: 用户数据未持久化

**问题描述**:
- 用户点赞、收藏数据仅存储在localStorage
- 清除浏览器数据会导致数据丢失
- 多设备无法同步

**影响范围**:
- 用户数据易丢失
- 无法实现多端同步
- 用户体验受限

**解决方案**:
1. 设计用户数据表结构
2. 实现后端用户数据API
3. 前端实现数据同步逻辑

**预计工作量**: 5天

---

#### 问题P0-004: 错误重试机制缺失

**问题描述**:
- 网络错误无自动重试
- arXiv API调用失败无降级处理
- 离线状态无检测

**影响范围**:
- 网络不稳定时用户体验差
- API暂时故障导致功能不可用

**解决方案**:
1. 前端实现指数退避重试
2. 添加离线状态检测
3. 实现友好的错误提示

**预计工作量**: 2天

---

### 3.2 中优先级问题（P1）

| 问题ID | 问题描述 | 影响 | 工作量 |
|--------|---------|------|-------|
| P1-001 | 内存缓存重启丢失 | 性能下降 | 3天 |
| P1-002 | 无请求限流保护 | 安全风险 | 1天 |
| P1-003 | 日志输出不规范 | 调试困难 | 1天 |
| P1-004 | 无性能监控 | 问题定位困难 | 3天 |
| P1-005 | 无虚拟列表 | 大列表性能问题 | 4天 |
| P1-006 | arXiv API调用未优化 | 响应慢 | 2天 |

---

### 3.3 低优先级问题（P2）

| 问题ID | 问题描述 | 影响 | 工作量 |
|--------|---------|------|-------|
| P2-001 | TypeScript严格模式未全开启 | 潜在类型错误 | 2天 |
| P2-002 | 无CI/CD流水线 | 部署效率低 | 3天 |
| P2-003 | 配置文件含测试域名 | 部署风险 | 0.5天 |
| P2-004 | 无统一错误处理层 | 体验不一致 | 1天 |

---

## 4. 优先级分类的改进建议

### 4.1 快速胜利（高影响 + 低难度）

| ID | 项目 | 影响 | 工时 | 关联PRD |
|----|------|------|------|---------|
| QW-001 | 批量获取论文接口 | 高 | 1天 | BE-013 |
| QW-002 | 收藏列表完整显示 | 高 | 2天 | FE-020 |
| QW-003 | 错误重试机制 | 高 | 1天 | FE-022 |
| QW-004 | 请求限流中间件 | 高 | 1天 | BE-015 |
| QW-005 | 日志标准化输出 | 中 | 1天 | 新增 |
| QW-006 | 分享功能完善 | 中 | 0.5天 | FE-021 |
| QW-007 | API响应时间监控 | 中 | 0.5天 | 新增 |
| QW-008 | 前端加载态优化 | 中 | 1天 | 新增 |

**总计工时**: 8天
**预计完成时间**: 2周

---

### 4.2 战略投资（高影响 + 高难度）

| ID | 项目 | 影响 | 工时 | 关联PRD |
|----|------|------|------|---------|
| SI-001 | Redis分布式缓存 | 高 | 3天 | BE-014 |
| SI-002 | 用户数据持久化(MySQL) | 高 | 5天 | 新增 |
| SI-003 | 收藏云同步功能 | 高 | 3天 | FE-026 |
| SI-004 | 虚拟列表性能优化 | 高 | 4天 | 新增 |
| SI-005 | 前端测试环境搭建 | 高 | 4天 | 新增 |
| SI-006 | 推荐算法基础版 | 高 | 10天 | BE-016 |
| SI-007 | CI/CD流水线搭建 | 高 | 3天 | 新增 |
| SI-008 | 监控告警系统 | 高 | 5天 | 新增 |

**总计工时**: 37天
**预计完成时间**: 8周

---

### 4.3 渐进改进（低影响 + 低难度）

| ID | 项目 | 影响 | 工时 | 关联PRD |
|----|------|------|------|---------|
| PI-001 | 阅读历史记录 | 低 | 2天 | FE-024 |
| PI-002 | 搜索历史 | 低 | 1天 | FE-025 |
| PI-003 | 论文详情页 | 中 | 3天 | FE-023 |
| PI-004 | 暗色/亮色模式切换 | 低 | 2天 | FE-027 |
| PI-005 | 键盘快捷键帮助面板 | 低 | 1天 | 新增 |
| PI-006 | 论文分类标签优化 | 低 | 1天 | 新增 |
| PI-007 | 图片懒加载优化 | 低 | 1天 | 新增 |

**总计工时**: 11天
**预计完成时间**: 3周

---

### 4.4 延后考虑（低影响 + 高难度）

| ID | 项目 | 影响 | 工时 | 关联PRD |
|----|------|------|------|---------|
| LC-001 | 多语言支持(i18n) | 低 | 5天 | FE-028 |
| LC-002 | PDF预览功能 | 中 | 5天 | FE-029 |
| LC-003 | 第三方OAuth登录 | 中 | 5天 | FE-032 |
| LC-004 | 论文笔记功能 | 低 | 8天 | FE-030 |
| LC-005 | 社交分享集成 | 低 | 3天 | FE-031 |
| LC-006 | 引文关系图可视化 | 低 | 15天 | BE-018 |

**总计工时**: 41天
**预计完成时间**: 长期规划

---

## 5. 技术债务处理方案

### 5.1 技术债务清单

#### 高优先级技术债务

| ID | 债务描述 | 影响 | 建议处理方案 | 计划时间 |
|----|---------|------|------------|---------|
| TD-001 | 用户数据仅localStorage存储 | 数据丢失风险 | 实现MySQL持久化 | Week 5-6 |
| TD-002 | 收藏列表仅显示ID | 功能不完整 | 实现批量接口+完整显示 | Week 1 |
| TD-003 | 无请求限流保护 | DDoS风险 | 实现限流中间件 | Week 1 |
| TD-004 | 日志输出不规范 | 调试困难 | 统一日志格式和级别 | Week 1 |
| TD-005 | 无错误重试机制 | 用户体验差 | 实现自动重试 | Week 1 |

#### 中优先级技术债务

| ID | 债务描述 | 影响 | 建议处理方案 | 计划时间 |
|----|---------|------|------------|---------|
| TD-006 | 内存缓存重启丢失 | 性能下降 | 引入Redis持久化缓存 | Week 3 |
| TD-007 | 无性能监控 | 问题定位困难 | 集成APM工具 | Week 2 |
| TD-008 | 前端无虚拟列表 | 大列表性能问题 | 引入react-window | Week 4 |
| TD-009 | arXiv API调用未优化 | 响应慢 | 批量查询+缓存 | Week 3 |
| TD-010 | 前端测试缺失 | 质量保障不足 | 搭建测试环境 | Week 2 |

#### 低优先级技术债务

| ID | 债务描述 | 影响 | 建议处理方案 | 计划时间 |
|----|---------|------|------------|---------|
| TD-011 | TypeScript严格模式未全开 | 潜在类型错误 | 逐步开启strict模式 | Week 6 |
| TD-012 | 部分组件无单元测试 | 维护风险 | 补充测试用例 | Week 2 |
| TD-013 | 无CI/CD流水线 | 部署效率低 | 搭建GitHub Actions | Week 7 |
| TD-014 | 前端无统一错误处理 | 体验不一致 | 统一错误处理层 | Week 1 |

### 5.2 技术债务处理原则

1. **优先级排序**: 按影响程度和风险等级排序
2. **增量偿还**: 每个迭代预留20%时间处理技术债务
3. **预防为主**: 新功能开发时避免引入新债务
4. **定期评审**: 每月评审技术债务清单

### 5.3 技术债务预防措施

#### 代码审查清单

- [ ] 是否引入了新的硬编码配置？
- [ ] 是否有充分的错误处理？
- [ ] 是否添加了必要的测试？
- [ ] 是否更新了相关文档？
- [ ] 是否考虑了性能影响？
- [ ] 是否有潜在的安全风险？

---

## 6. 功能演进路线图

### 6.1 短期路线图（1-3个月）

#### 目标：完善核心功能，提升用户体验

```
Month 1: 安全与稳定性
├── Week 1-2: 紧急修复
│   ├── JWT安全验证（已完成）
│   ├── 批量获取论文接口
│   ├── 收藏列表完整显示
│   ├── 错误重试机制
│   └── 请求限流中间件
└── Week 3-4: 体验改进
    ├── 日志标准化
    ├── API响应监控
    ├── 前端加载态优化
    └── 分享功能完善

Month 2: 性能优化
├── Week 5-6: 缓存升级
│   ├── Redis集成
│   ├── arXiv API调用优化
│   └── 缓存策略优化
└── Week 7-8: 前端性能
    ├── 虚拟列表实现
    ├── 图片懒加载优化
    └── 首屏加载优化

Month 3: 数据持久化
├── Week 9-10: 用户数据同步
│   ├── MySQL表结构设计
│   ├── 用户数据API实现
│   └── 前端数据同步
└── Week 11-12: 测试补充
    ├── 前端测试环境
    ├── 后端测试补充
    └── 集成测试实现
```

### 6.2 中期路线图（3-6个月）

#### 目标：个性化推荐，增强互动

```
Month 4-5: 推荐系统
├── 用户行为追踪
├── 推荐算法基础版
└── A/B测试框架

Month 5-6: 体验增强
├── 阅读历史记录
├── 搜索历史
├── 论文详情页
└── 暗色模式支持
```

### 6.3 长期路线图（6-12个月）

#### 目标：架构升级，生态扩展

```
Month 7-9: 架构演进
├── CI/CD流水线
├── 监控告警系统
└── 服务拆分准备

Month 10-12: 功能扩展
├── PWA支持
├── 社交分享
├── 第三方登录
└── 多语言支持
```

### 6.4 架构演进方向

#### 当前架构（单体）

```
┌─────────────────────────────────────────┐
│         PaperTok 单体应用                │
│  ┌───────────┐  ┌───────────┐          │
│  │ React前端 │  │  Go后端   │          │
│  └───────────┘  └───────────┘          │
└─────────────────────────────────────────┘
```

#### 演进架构（微服务）

```
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  User Svc   │  │  Paper Svc  │  │ Recommend   │         │
│  │  (用户服务)  │  │  (论文服务)  │  │  (推荐服务)  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   MySQL     │  │   Redis     │  │  RocketMQ   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

**触发条件**:
- QPS > 1000
- 团队规模 > 5人
- 单个服务部署周期 > 30分钟

---

## 7. 实施计划

### 7.1 第一阶段：安全与稳定性（Week 1-2）

#### Sprint 1.1 (Week 1)

| 任务 | 优先级 | 工时 | 负责模块 | 验收标准 |
|------|-------|------|---------|---------|
| 批量获取论文接口 | P0 | 1天 | Backend | POST /api/v1/papers/batch可用 |
| 收藏列表完整显示 | P0 | 2天 | Frontend | 显示完整论文信息 |
| 错误重试机制 | P0 | 1天 | Frontend | 网络错误自动重试3次 |
| 请求限流中间件 | P0 | 1天 | Backend | 每IP每分钟60次 |
| 日志标准化输出 | P1 | 1天 | Backend | JSON格式日志 |

#### Sprint 1.2 (Week 2)

| 任务 | 优先级 | 工时 | 负责模块 | 验收标准 |
|------|-------|------|---------|---------|
| API响应时间监控 | P1 | 0.5天 | Backend | 响应时间记录到日志 |
| 分享功能完善 | P1 | 0.5天 | Frontend | 分享成功有Toast提示 |
| 前端加载态优化 | P1 | 1天 | Frontend | 骨架屏流畅显示 |
| 单元测试补充 | P1 | 2天 | Both | 核心模块测试覆盖70% |
| 安全漏洞扫描 | P1 | 1天 | Both | 无高危漏洞 |

### 7.2 第二阶段：性能优化（Week 3-4）

#### Sprint 2.1 (Week 3)

| 任务 | 优先级 | 工时 | 技术方案 | 预期收益 |
|------|-------|------|---------|---------|
| Redis分布式缓存 | P0 | 3天 | go-redis/v9 | 缓存命中率>80% |
| arXiv API调用优化 | P1 | 2天 | 批量查询 | API调用减少70% |
| 前端资源加载优化 | P1 | 2天 | 代码分割 | 首屏加载<1.5s |

#### Sprint 2.2 (Week 4)

| 任务 | 优先级 | 工时 | 技术方案 | 预期收益 |
|------|-------|------|---------|---------|
| 虚拟列表性能优化 | P0 | 4天 | react-window | 滚动FPS稳定60 |
| 图片懒加载优化 | P1 | 1天 | IntersectionObserver | 内存占用减少60% |
| 首屏加载优化 | P1 | 1天 | 预加载 | 首屏<1.5s |

### 7.3 第三阶段：功能完善（Week 5-8）

#### Sprint 3.1 (Week 5-6)

| 任务 | 优先级 | 工时 | 负责模块 | 验收标准 |
|------|-------|------|---------|---------|
| 用户数据持久化 | P0 | 5天 | Backend | MySQL存储用户数据 |
| 收藏云同步功能 | P0 | 3天 | Frontend | 数据与后端同步 |
| 阅读历史记录 | P1 | 2天 | Frontend | 记录浏览历史 |

#### Sprint 3.2 (Week 7-8)

| 任务 | 优先级 | 工时 | 负责模块 | 验收标准 |
|------|-------|------|---------|---------|
| 前端测试环境搭建 | P0 | 3天 | Frontend | Jest+RTL可用 |
| 核心组件测试 | P0 | 3天 | Frontend | 覆盖率>70% |
| 论文详情页 | P1 | 3天 | Frontend | 详情页可用 |
| 搜索历史 | P1 | 1天 | Frontend | 记录搜索历史 |

### 7.4 第四阶段：架构演进（Week 9+）

#### 持续改进项

| 任务 | 优先级 | 预计工时 | 说明 |
|------|-------|---------|------|
| CI/CD流水线 | P1 | 3天 | GitHub Actions |
| 监控告警系统 | P1 | 5天 | Prometheus+Grafana |
| 推荐算法基础版 | P1 | 10天 | 基于协同过滤 |
| 服务拆分准备 | P2 | - | 评估必要性 |

---

## 8. 风险管控

### 8.1 实施风险

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| arXiv API限流 | 中 | 高 | 增加缓存、延长TTL、实现降级 |
| 数据库迁移失败 | 低 | 高 | 充分测试、灰度发布、回滚预案 |
| 性能优化效果不佳 | 中 | 中 | 基准测试、性能监控、渐进优化 |
| 第三方依赖不稳定 | 低 | 中 | 版本锁定、定期更新、备用方案 |

### 8.2 技术风险

| 风险 | 说明 | 缓解措施 |
|------|------|---------|
| Redis单点故障 | 缓存数据丢失 | 使用Redis Sentinel/Cluster |
| MySQL连接池耗尽 | 高并发场景 | 合理配置连接池参数 |
| 前端内存泄漏 | 虚拟列表实现不当 | 内存监控、正确cleanup |
| JWT密钥泄露 | 安全风险 | 定期轮换、环境变量隔离 |

### 8.3 产品风险

| 风险 | 说明 | 缓解措施 |
|------|------|---------|
| 用户获取困难 | 学术工具推广挑战 | SEO优化、社区运营 |
| 竞品压力 | 类似产品竞争 | 差异化功能、用户体验 |
| 数据源变更 | arXiv API变更 | 版本监控、适配层隔离 |

---

## 9. 成功指标

### 9.1 性能指标

| 指标 | 当前 | 目标 | 测量方式 |
|------|------|------|---------|
| API平均响应时间 | ~2s | <500ms | Prometheus |
| 首屏加载时间 | ~2s | <1.5s | Lighthouse |
| 滚动FPS | ~45 | 稳定60 | Chrome DevTools |
| 缓存命中率 | 0% | >80% | Redis监控 |
| API错误率 | <1% | <0.1% | 日志分析 |

### 9.2 业务指标

| 指标 | 当前 | 目标 |
|------|------|------|
| 日活用户(DAU) | - | 稳定增长 |
| 平均浏览时长 | - | >10分钟 |
| 收藏转化率 | - | >20% |
| 用户留存率(次日) | - | >40% |
| 用户留存率(7日) | - | >20% |

### 9.3 质量指标

| 指标 | 当前 | 目标 |
|------|------|------|
| 后端测试覆盖率 | ~60% | >80% |
| 前端测试覆盖率 | 0% | >70% |
| 代码重复率 | <5% | <3% |
| TypeScript覆盖率 | ~90% | 100% |
| 安全漏洞数 | 0 | 0 |

---

## 附录A：技术实现示例

### A.1 批量获取论文接口

```go
// internal/features/paperbatch/service.go
package paperbatch

type Service interface {
    GetBatch(ctx context.Context, ids []string) ([]*Paper, error)
}

type Impl struct {
    arxivSvc  arxivService
    paperRepo paperRepository
}

func (s *Impl) GetBatch(ctx context.Context, ids []string) ([]*Paper, error) {
    // 1. 尝试从缓存获取
    cached, missing := s.paperRepo.GetBatch(ctx, ids)

    // 2. 如果全部命中缓存，直接返回
    if len(missing) == 0 {
        return cached, nil
    }

    // 3. 从arXiv获取缺失的论文
    fetched, err := s.arxivSvc.FetchByID(ctx, missing)
    if err != nil {
        return nil, err
    }

    // 4. 合并结果并缓存
    result := append(cached, fetched...)
    s.paperRepo.SaveBatch(ctx, fetched)

    return result, nil
}
```

### A.2 Redis缓存实现

```go
// internal/infra/cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return &RedisCache{client: client, ttl: ttl}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(val, dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, c.ttl).Err()
}
```

### A.3 前端错误重试Hook

```typescript
// frontend/src/hooks/useRetry.ts
import { useState, useCallback } from 'react';

interface RetryOptions {
  maxAttempts?: number;
  delay?: number;
  backoff?: boolean;
}

export function useRetry<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  options: RetryOptions = {}
) {
  const { maxAttempts = 3, delay = 1000, backoff = true } = options;
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(
    async (...args: Parameters<T>) => {
      setIsLoading(true);
      setError(null);

      let lastError: Error | null = null;
      let currentDelay = delay;

      for (let attempt = 1; attempt <= maxAttempts; attempt++) {
        try {
          const result = await fn(...args);
          setIsLoading(false);
          return result;
        } catch (err) {
          lastError = err as Error;
          if (attempt < maxAttempts) {
            await new Promise(resolve => setTimeout(resolve, currentDelay));
            if (backoff) currentDelay *= 2;
          }
        }
      }

      setIsLoading(false);
      setError(lastError);
      throw lastError;
    },
    [fn, maxAttempts, delay, backoff]
  );

  return { execute, isLoading, error };
}
```

### A.4 虚拟列表实现

```typescript
// frontend/src/components/VirtualFeed.tsx
import { FixedSizeList as List } from 'react-window';
import { Paper } from '../types';
import PaperCard from './PaperCard';

interface VirtualFeedProps {
  papers: Paper[];
  height: number;
}

export function VirtualFeed({ papers, height }: VirtualFeedProps) {
  const Row = ({ index, style }: { index: number; style: React.CSSProperties }) => (
    <div style={style}>
      <PaperCard paper={papers[index]} />
    </div>
  );

  return (
    <List
      height={height}
      itemCount={papers.length}
      itemSize={height} // 全屏卡片
      width="100%"
    >
      {Row}
    </List>
  );
}
```

---

## 附录B：数据库Schema设计

### B.1 用户收藏表

```sql
CREATE TABLE user_favorites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    paper_id VARCHAR(255) NOT NULL COMMENT 'arXiv论文ID',
    paper_data JSON NOT NULL COMMENT '论文数据缓存',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_paper (user_id, paper_id),
    KEY idx_user_id (user_id),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';
```

### B.2 用户浏览历史

```sql
CREATE TABLE user_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    paper_id VARCHAR(255) NOT NULL COMMENT 'arXiv论文ID',
    paper_data JSON NOT NULL COMMENT '论文数据缓存',
    view_duration INT DEFAULT 0 COMMENT '浏览时长(秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_user_id (user_id),
    KEY idx_created_at (created_at),
    KEY idx_user_created (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户浏览历史';
```

### B.3 用户点赞表

```sql
CREATE TABLE user_likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    paper_id VARCHAR(255) NOT NULL COMMENT 'arXiv论文ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_paper (user_id, paper_id),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户点赞表';
```

---

## 附录C：相关文档索引

| 文档 | 路径 | 说明 |
|------|------|------|
| 产品需求文档 | docs/PRD.md | 功能清单、API规范 |
| 架构设计文档 | docs/ARCHITECTURE.md | 系统架构、技术栈 |
| 设计原则文档 | docs/DESIGN_PRINCIPLES.md | 架构原则、模块规范 |
| 项目调研报告 | docs/PROJECT_RESEARCH_REPORT.md | 代码质量分析 |
| 改进建议提案 | docs/IMPROVEMENT_PROPOSAL.md | 改进优先级矩阵 |
| 改进建议矩阵 | docs/IMPROVEMENT_MATRIX.md | 可视化矩阵图 |
| 本地开发指南 | backend/local-guides/ | 后端开发指南 |

---

**文档维护**: 每月评审更新
**版本控制**: Git跟踪变更
**评审周期**: 每个迭代开始前
**责任人**: 技术架构师

---

**变更历史**

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| v1.0 | 2026-02-02 | 初始版本，基于调研报告生成 | 技术架构师 |
