# PRD 快速参考

> 开发时的快速查询指南

---

## 📋 功能状态速查

### 已实现 ✅

**后端**:
- BE-001: 论文推荐流 API
- BE-002: 论文搜索 API
- BE-003: 论文详情 API
- BE-004: arXiv API 集成
- BE-005: 内存缓存
- BE-006: 健康检查接口

**前端**:
- FE-001: 论文推荐流页面
- FE-002: 收藏列表页面（基础版）
- FE-003: 搜索页面
- FE-004: 论文卡片组件
- FE-005: 分类筛选器
- FE-006: 点赞功能
- FE-007: 收藏功能
- FE-008: 键盘快捷键
- FE-009: 无限滚动加载
- FE-010: 本地存储

### 待实现 🔄

**高优先级（P0）**:
- FE-011: 收藏列表完整显示
- FE-012: 分享功能完善
- FE-013: 错误重试机制
- BE-007: 批量获取论文接口

**中优先级（P1）**:
- FE-014: 论文详情页
- FE-015: 阅读历史记录
- FE-016: 搜索历史
- BE-008: Redis 缓存
- BE-009: 请求限流

---

## 🔗 API 端点速查

| 端点 | 方法 | 说明 | 状态 |
|------|------|------|------|
| `/api/v1/papers` | GET | 获取论文列表 | ✅ |
| `/api/v1/papers/search` | GET | 搜索论文 | ✅ |
| `/api/v1/papers/:id` | GET | 获取论文详情 | ✅ |
| `/api/v1/papers/batch` | POST | 批量获取论文 | ⏳ |
| `/health` | GET | 健康检查 | ✅ |

---

## 🗂️ 架构映射速查

| 功能 | 前端 | 后端 | API |
|------|------|------|-----|
| 推荐流 | `pages/Feed.tsx` | `features/paperfeed` | `GET /api/v1/papers` |
| 搜索 | `pages/Search.tsx` | `features/papersearch` | `GET /api/v1/papers/search` |
| 详情 | - | `features/papersearch` | `GET /api/v1/papers/:id` |
| 收藏 | `pages/Favorites.tsx` | - | - |

---

## 📝 开发流程速查

### 新增功能

1. **更新 PRD**: 添加功能到"待实现功能"
2. **分配功能ID**: 按规则分配（FE-XXX 或 BE-XXX）
3. **开发**: 按照架构规范实现
4. **测试**: 编写并运行测试
5. **更新 PRD**: 移动到"已实现功能"，更新日志

### 功能完成

1. 更新功能状态为 ✅
2. 填写完成时间
3. 在更新日志中记录
4. Git 提交: `docs(prd): 更新功能 FE-XXX 状态为已完成`

---

## 🎯 当前迭代（v2.1）

**时间**: 2026-01-26 ~ 2026-02-09

**目标功能**:
- [ ] FE-011: 收藏列表完整显示
- [ ] FE-012: 分享功能完善
- [ ] FE-013: 错误重试机制
- [ ] BE-007: 批量获取论文接口

---

## 📊 功能ID分配

| 类型 | 下一个可用ID |
|------|-------------|
| BE | BE-007 |
| FE | FE-011 |

---

## 🔍 快速查找

- **完整 PRD**: [PRD.md](./PRD.md)
- **维护指南**: [PRD_MAINTENANCE.md](./PRD_MAINTENANCE.md)
- **架构文档**: [ARCHITECTURE.md](./ARCHITECTURE.md)
- **设计原则**: [DESIGN_PRINCIPLES.md](./DESIGN_PRINCIPLES.md)

---

**最后更新**: 2026-01-26  
**维护者**: 开发团队
