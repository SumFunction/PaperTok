# 用户认证前端模块实现文档

## 变更摘要

### 新增文件

#### 类型定义更新
- `src/types/index.ts` - 添加了用户认证相关类型

#### Context 层
- `src/contexts/AuthContext.tsx` - 认证状态管理 Context
- `src/contexts/AuthContext.tsx` - 导出 `useAuth` Hook 和 `requireAuth` HOC

#### 组件层
- `src/components/LoginForm.tsx` - 登录表单组件
- `src/components/LoginForm.css` - 登录表单样式
- `src/components/RegisterForm.tsx` - 注册表单组件
- `src/components/RegisterForm.css` - 注册表单样式
- `src/components/ProtectedRoute.tsx` - 路由保护组件
- `src/components/ProtectedRoute.css` - 路由保护样式

#### 页面层
- `src/pages/Login.tsx` - 登录页面
- `src/pages/Login.css` - 登录页面样式
- `src/pages/Register.tsx` - 注册页面
- `src/pages/Register.css` - 注册页面样式

### 修改文件

- `src/services/api.ts` - 扩展 API 服务，添加认证相关接口（login, register, getProfile, logout）
- `src/App.tsx` - 添加 AuthProvider、登录注册路由、路由保护
- `src/components/Navigation.tsx` - 根据登录状态显示不同导航选项
- `src/components/Navigation.css` - 添加用户信息和登出按钮样式

## 关键设计决策

### 1. 认证状态管理
- 使用独立的 `AuthContext` 管理认证状态，与现有的 `AppContext` 分离
- 认证状态包括：user, token, isAuthenticated, loading, error
- 使用 localStorage 存储 token 和用户信息，实现持久化登录

### 2. Token 管理
- 在 API 服务层实现 `TokenManager` 类，统一管理 token 的存储、获取和清除
- 通过 axios 请求拦截器自动添加 Bearer token
- 通过 axios 响应拦截器处理 401 未授权错误，自动清除认证信息并重定向到登录页

### 3. 路由保护
- `ProtectedRoute` 组件用于保护需要登录才能访问的路由
- 未登录时自动重定向到登录页，并保存原始访问路径
- 登录成功后自动跳转回原访问页面

### 4. 表单验证
- 实时验证：输入时清除错误，失焦时进行验证
- 邮箱格式验证使用正则表达式
- 密码强度指示器（弱/中/强）
- 用户名验证：长度限制、字符限制（字母、数字、下划线）

### 5. UI/UX 设计
- 抖音风格：深色主题、渐变色强调、简洁现代
- 响应式设计：适配移动端、桌面端、横屏手机
- 加载状态：统一的 spinner 加载动画
- 错误提示：清晰的错误消息展示

## 类型定义

```typescript
// 用户信息
interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

// 认证响应
interface AuthResponse {
  user: User;
  token: string;
}

// 登录请求
interface LoginRequest {
  email: string;
  password: string;
}

// 注册请求
interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}
```

## API 接口

### 登录
- **接口**: `POST /api/v1/auth/login`
- **请求**: `{ email, password }`
- **响应**: `{ success: boolean, data: { user, token } }`

### 注册
- **接口**: `POST /api/v1/auth/register`
- **请求**: `{ username, email, password }`
- **响应**: `{ success: boolean, data: { user, token } }`

### 获取用户信息
- **接口**: `GET /api/v1/auth/profile`
- **请求头**: `Authorization: Bearer {token}`
- **响应**: `{ success: boolean, data: User }`

### 登出
- **接口**: `POST /api/v1/auth/logout`
- **请求头**: `Authorization: Bearer {token}`

## 路由配置

```typescript
// 公共路由
/login  - 登录页面
/register  - 注册页面
/  - 首页（推荐）
/search  - 搜索

// 受保护路由
/favorites  - 收藏（需要登录）
```

## 导航状态

### 未登录状态
- 推荐、搜索、登录

### 已登录状态
- 推荐、搜索、收藏、用户信息（显示用户名和登出按钮）

## 验证方式

### 本地开发
```bash
cd frontend
npm install
npm run dev
```

### 构建验证
```bash
cd frontend
npm run build
```

### 功能测试
1. 访问 `/login` 页面，查看登录表单是否正常显示
2. 访问 `/register` 页面，查看注册表单是否正常显示
3. 测试表单验证（输入无效邮箱、短密码等）
4. 测试登录/注册流程（需要后端 API 支持）
5. 测试路由保护（未登录访问 `/favorites` 应重定向到登录页）
6. 测试登出功能

## 后续 TODO

1. 添加"记住我"功能
2. 添加密码重置功能
3. 添加第三方登录（Google、GitHub 等）
4. 添加邮箱验证流程
5. 添加用户个人资料编辑页面
6. 优化表单提交性能（防抖、节流）
7. 添加单元测试
8. 添加国际化支持

## 注意事项

1. **安全性**：密码在前端不以明文形式存储，token 存储在 localStorage（后续可考虑使用 httpOnly cookie）
2. **兼容性**：使用现代浏览器 API（localStorage、fetch），需要考虑旧浏览器降级
3. **性能**：认证状态初始化可能需要时间，已在 App 层面处理加载状态
4. **错误处理**：API 错误会统一显示在表单顶部，401 错误会自动清除认证信息
