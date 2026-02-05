# 登录功能修复 + 抖音式登录弹窗实现

> 实施日期: 2026-01-27

## 变更摘要

### 第一阶段：修复登录功能

| 文件 | 变更内容 |
|------|----------|
| `frontend/src/services/api.ts` | 修复登录 API 字段不匹配，将 `email` 改为 `identifier`；增强错误提示映射 |
| `frontend/src/components/LoginForm.tsx` | 更新密码验证规则从 6 字符改为 8 字符 |
| `frontend/src/components/RegisterForm.tsx` | 更新密码验证规则从 6 字符改为 8 字符 |

### 第二阶段：抖音式登录弹窗

| 文件 | 变更内容 |
|------|----------|
| `frontend/src/components/AuthModal.tsx` | 新增：登录/注册模态框组件 |
| `frontend/src/components/AuthModal.css` | 新增：模态框样式（抖音风格） |
| `frontend/src/hooks/useViewCounter.ts` | 新增：浏览计数器 Hook |
| `frontend/src/contexts/AuthContext.tsx` | 扩展：添加模态框状态管理 |
| `frontend/src/pages/Feed.tsx` | 集成：浏览计数 + 弹窗触发 |
| `frontend/src/App.tsx` | 新增：首次访问弹窗逻辑 |

---

## 关键决策

### 1. API 字段命名
- **决策**：使用 `identifier` 而非 `email` 作为登录字段
- **理由**：后端支持邮箱或用户名登录，`identifier` 更具语义通用性

### 2. 密码验证规则
- **决策**：前后端统一为最少 8 字符
- **理由**：提高安全性，避免前端通过验证后端拒绝的情况

### 3. 浏览计数器实现
- **决策**：使用 Set 记录已浏览论文 ID，避免重复计数
- **理由**：用户来回滑动不应该重复计数，确保计数准确性

### 4. 模态框触发策略
- **决策**：
  - 首次访问：延迟 2 秒弹出（可关闭）
  - 浏览计数：滑动 10 个不同论文后弹出
  - 登录后：重置计数，不再自动弹出
- **理由**：平衡用户体验和推广需求，避免过度打扰

### 5. localStorage 持久化
- **决策**：浏览计数和首次访问标记都持久化到 localStorage
- **理由**：刷新页面后状态保留，符合用户预期

---

## 验证方式

### 启动服务

```bash
# 启动后端
cd backend && go run cmd/server/main.go

# 启动前端
cd frontend && npm run dev
```

### 测试场景

#### 1. 登录功能验证
- [ ] 使用 8+ 字符密码注册新用户
- [ ] 使用邮箱登录
- [ ] 使用用户名登录
- [ ] 错误密码显示正确提示
- [ ] 登录后刷新页面保持登录状态

#### 2. 首次访问弹窗
- [ ] 清除 localStorage 后首次访问
- [ ] 2 秒后自动弹出登录提示
- [ ] 关闭后不再自动弹出
- [ ] 点击遮罩可关闭

#### 3. 浏览计数弹窗
- [ ] 未登录状态下滑动浏览
- [ ] 浏览 10 个不同论文后弹出登录提示
- [ ] 登录成功后弹窗关闭且计数重置
- [ ] 登录后不再触发弹窗

#### 4. 模态框交互
- [ ] Tab 切换登录/注册
- [ ] ESC 键关闭
- [ ] 点击遮罩关闭
- [ ] 登录/注册成功后自动关闭
- [ ] 移动端样式正常（底部弹出效果）

---

## 技术细节

### API 请求格式（修复后）

```json
// 登录
POST /api/v1/auth/login
{
  "identifier": "user@example.com",  // 支持邮箱或用户名
  "password": "password123"           // 至少 8 字符
}

// 注册
POST /api/v1/auth/register
{
  "username": "testuser",
  "email": "user@example.com",
  "password": "password123"           // 至少 8 字符
}
```

### 错误码映射

| 错误码 | 提示信息 |
|--------|----------|
| `INVALID_CREDENTIALS` | 邮箱或密码错误 |
| `USER_NOT_FOUND` | 用户不存在 |
| `WEAK_PASSWORD` | 密码强度不足，至少需要8个字符 |
| `USER_EXISTS` | 用户已存在 |
| `EMAIL_EXISTS` | 邮箱已被注册 |
| `NETWORK_ERROR` | 网络连接失败，请检查网络 |

### 浏览计数器逻辑

```typescript
// 使用 Set 记录已浏览的论文 ID（避免重复计数）
const viewedPapers = new Set<string>();

// 记录浏览
viewedPapers.add(paperId);

// 达到阈值触发回调
if (viewedPapers.size >= threshold && !isAuthenticated) {
  openAuthModal('login');
}

// 登录后重置
if (isAuthenticated) {
  viewedPapers.clear();
}
```

### 模态框动画

```css
/* 进入动画：淡入 + 上移 + 缩放 */
@keyframes modalEnter {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* 移动端：底部滑入 */
@media (max-width: 480px) {
  @keyframes slideUp {
    from { opacity: 0; transform: translateY(100%); }
    to { opacity: 1; transform: translateY(0); }
  }
}
```

---

## 后续 TODO

1. **后端持久化**：当前用户数据存储在内存中，服务重启后丢失
   - 建议：实现 MySQL 持久化

2. **更多认证方式**：
   - OAuth 第三方登录（Google、GitHub）
   - 邮箱验证码登录

3. **记住我功能**：
   - 长期 token（30天）
   - 自动登录

4. **密码重置**：
   - 忘记密码流程
   - 邮件验证

5. **A/B 测试**：
   - 测试不同弹窗时机对转化率的影响
