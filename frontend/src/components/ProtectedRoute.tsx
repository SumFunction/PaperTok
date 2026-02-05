/**
 * 路由保护组件
 * 保护需要��录才能访问的路由
 */

import { useEffect, type ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

/**
 * ProtectedRoute 组件属性
 */
interface ProtectedRouteProps {
  children: ReactNode;
  /**
   * 未登录时重定向的��径，默认为 /login
   */
  redirectTo?: string;
}

/**
 * ProtectedRoute 组件
 * 如果用户未登录，重定向到登录页
 * 否则渲染子组件
 */
export function ProtectedRoute({ children, redirectTo = '/login' }: ProtectedRouteProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, loading } = useAuth();

  useEffect(() => {
    // 等待认证状态初始化完成
    if (loading) {
      return;
    }

    // 如果用户未登录，重定向到登录页
    if (!isAuthenticated) {
      // 保存当前路径，登录后可以返回
      navigate(redirectTo, {
        replace: true,
        state: { from: location },
      });
    }
  }, [isAuthenticated, loading, navigate, location, redirectTo]);

  // 加载状态
  if (loading) {
    return (
      <div className="protected-route-loading">
        <div className="protected-route-loading__spinner" />
        <p>加载中...</p>
      </div>
    );
  }

  // 未认证时不渲染任何内容（即将重定向）
  if (!isAuthenticated) {
    return null;
  }

  // 已认证，渲染子组件
  return <>{children}</>;
}

/**
 * 高阶组件版本
 * 用于包装需要认证的组件
 */
export function withAuth<P extends object>(
  WrappedComponent: React.ComponentType<P>,
  redirectTo?: string
): React.ComponentType<P> {
  return function AuthWrappedComponent(props: P) {
    return (
      <ProtectedRoute redirectTo={redirectTo}>
        <WrappedComponent {...props} />
      </ProtectedRoute>
    );
  };
}
