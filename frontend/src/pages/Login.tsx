/**
 * 登录页面
 * 提供用户登录功能，登录成功后跳转到首页
 */

import { useEffect, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { LoginForm } from '../components/LoginForm';
import './Login.css';

/**
 * Login 页面组件
 */
export function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, loading } = useAuth();

  // 获取登录前要访问的页面
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/';

  /**
   * 如果用户已登录，重定向到原访问页面或首页
   */
  useEffect(() => {
    if (isAuthenticated && !loading) {
      console.log('[Login] 用户已登录，跳转到:', from);
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, loading, navigate, from]);

  /**
   * 登录成功回调 - 导航到目标页面
   */
  const handleLoginSuccess = useCallback(() => {
    console.log('[Login] 登录成功，跳转到:', from);
    navigate(from, { replace: true });
  }, [navigate, from]);

  // 加载状态
  if (loading) {
    return (
      <div className="login-page">
        <div className="login-page__loading">
          <div className="login-page__spinner" />
          <p>加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="login-page">
      <div className="login-page__container">
        {/* Logo */}
        <div className="login-page__logo">
          <span className="login-page__logo-icon">📄</span>
          <span className="login-page__logo-text">PaperTok</span>
        </div>

        {/* 登录表单 */}
        <LoginForm onSuccess={handleLoginSuccess} />
      </div>
    </div>
  );
}
