/**
 * 认证状态管理 Context
 * 提供登录、注册、登出等认证功能
 */

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from 'react';
import { PaperTokAPI, TokenManager } from '../services/api';
import type { User, AuthResponse } from '../types';

/**
 * 认证状态接口
 */
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

/**
 * 认证模态框 Tab 类型
 */
export type AuthModalTab = 'login' | 'register';

/**
 * 认证模态框场景类型
 */
export type AuthModalScene = 'welcome' | 'limit_reached' | 'protected_route';

/**
 * AuthContext 接口
 */
interface AuthContextType extends AuthState {
  // Actions
  login: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
  clearError: () => void;
  // 模态框状态
  showAuthModal: boolean;
  authModalDefaultTab: AuthModalTab;
  authModalScene: AuthModalScene;
  openAuthModal: (defaultTab?: AuthModalTab, scene?: AuthModalScene) => void;
  closeAuthModal: () => void;
}

// 创建 Context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * Provider Props
 */
interface AuthProviderProps {
  children: ReactNode;
}

/**
 * AuthProvider 组件
 * 包裹应用以提供认证状态和方法
 */
export function AuthProvider({ children }: AuthProviderProps) {
  // 状态初始化
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // 模态框状态
  const [showAuthModal, setShowAuthModal] = useState<boolean>(false);
  const [authModalDefaultTab, setAuthModalDefaultTab] = useState<AuthModalTab>('login');
  const [authModalScene, setAuthModalScene] = useState<AuthModalScene>('welcome');

  /**
   * 初始化认证状态
   * 从 localStorage 恢复 token 和用户信息
   */
  useEffect(() => {
    const initAuth = async () => {
      try {
        const savedToken = TokenManager.getToken();
        const savedUser = TokenManager.getUser();

        console.log('[Auth] 初始化认证状态:', { hasToken: !!savedToken, hasUser: !!savedUser });

        if (savedToken && savedUser) {
          // 尝试从服务器验证 token
          try {
            const freshUser = await PaperTokAPI.getProfile();
            console.log('[Auth] 服务器验证成功，用户:', freshUser.username);
            setToken(savedToken);
            setUser(freshUser);
          } catch (error) {
            // 如果获取失败（包括 404、401 等），说明 token 无效，清除本地认证信息
            console.warn('[Auth] 服务器验证失败，清除本地认证信息:', error);
            TokenManager.clearAuth();
            // 不设置 user 和 token，保持未登录状态
          }
        }
      } catch (err) {
        console.error('[Auth] 认证初始化错误:', err);
        // 清除无效的认证信息
        TokenManager.clearAuth();
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  /**
   * 计算是否已认证
   */
  const isAuthenticated = Boolean(user && token);

  /**
   * 用户登录
   */
  const login = useCallback(async (email: string, password: string) => {
    setLoading(true);
    setError(null);

    try {
      const response: AuthResponse = await PaperTokAPI.login(email, password);
      setUser(response.user);
      setToken(response.token);
    } catch (err) {
      const message = err instanceof Error ? err.message : '登录失败，请稍后重试';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * 用户注册
   */
  const register = useCallback(async (username: string, email: string, password: string) => {
    setLoading(true);
    setError(null);

    try {
      const response: AuthResponse = await PaperTokAPI.register(username, email, password);
      setUser(response.user);
      setToken(response.token);
    } catch (err) {
      const message = err instanceof Error ? err.message : '注册失败，请稍后重试';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * 用户登出
   */
  const logout = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      await PaperTokAPI.logout();
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      // 无论接口调用是否成功，都清除本地状态
      setUser(null);
      setToken(null);
      setLoading(false);
    }
  }, []);

  /**
   * 刷新用户信息
   */
  const refreshUser = useCallback(async () => {
    if (!token) {
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const freshUser = await PaperTokAPI.getProfile();
      setUser(freshUser);
    } catch (err) {
      const message = err instanceof Error ? err.message : '获取用户信息失败';
      setError(message);
      // 如果获取用户信息失败（如 token 过期），清除认证状态
      if (message.includes('未授权') || message.includes('token')) {
        setUser(null);
        setToken(null);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [token]);

  /**
   * 清除错误信息
   */
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  /**
   * 打开认证模态框
   * @param defaultTab 默认显示的标签页
   * @param scene 弹窗场景（决定显示的提示文案）
   */
  const openAuthModal = useCallback((defaultTab: AuthModalTab = 'login', scene: AuthModalScene = 'welcome') => {
    setAuthModalDefaultTab(defaultTab);
    setAuthModalScene(scene);
    setShowAuthModal(true);
    // 打开模态框时清除之前的错误
    setError(null);
  }, []);

  /**
   * 关闭认证模态框
   */
  const closeAuthModal = useCallback(() => {
    setShowAuthModal(false);
    setError(null);
  }, []);

  // Context 值
  const value: AuthContextType = {
    // State
    user,
    token,
    isAuthenticated,
    loading,
    error,
    showAuthModal,
    authModalDefaultTab,
    authModalScene,

    // Actions
    login,
    register,
    logout,
    refreshUser,
    clearError,
    openAuthModal,
    closeAuthModal,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * 使用 AuthContext 的 Hook
 * @throws {Error} 如果在 AuthProvider 外部使用
 */
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

/**
 * HOC: 要求组件必须在认证状态下才能渲染
 * 用于保护需要登录的页面或组件
 */
export function requireAuth<P extends object>(
  WrappedComponent: React.ComponentType<P>
): React.ComponentType<P> {
  return function AuthRequiredComponent(props: P) {
    const { isAuthenticated, loading } = useAuth();

    if (loading) {
      return (
        <div className="auth-loading">
          <div className="auth-loading__spinner" />
          <p>加载中...</p>
        </div>
      );
    }

    if (!isAuthenticated) {
      // 这里应该由 ProtectedRoute 组件处理重定向
      return null;
    }

    return <WrappedComponent {...props} />;
  };
}
