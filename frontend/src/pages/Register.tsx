/**
 * 注册页面
 * 提供用户注册功能，注册成功后跳转到首页
 */

import { useEffect, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { RegisterForm } from '../components/RegisterForm';
import './Register.css';

/**
 * Register 页面组件
 */
export function Register() {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, loading } = useAuth();

  // 获取注册前想要访问的页面
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/';

  /**
   * 如果用户已登录，重定向到原访问页面或首页
   */
  useEffect(() => {
    if (isAuthenticated && !loading) {
      console.log('[Register] 用户已登录，跳转到:', from);
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, loading, navigate, from]);

  /**
   * 注册成功回调 - 导航到目标页面
   */
  const handleRegisterSuccess = useCallback(() => {
    console.log('[Register] 注册成功，跳转到:', from);
    navigate(from, { replace: true });
  }, [navigate, from]);

  // 加载状态
  if (loading) {
    return (
      <div className="register-page">
        <div className="register-page__loading">
          <div className="register-page__spinner" />
          <p>加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="register-page">
      <div className="register-page__container">
        {/* Logo */}
        <div className="register-page__logo">
          <span className="register-page__logo-icon">📄</span>
          <span className="register-page__logo-text">PaperTok</span>
        </div>

        {/* 注册表单 */}
        <RegisterForm onSuccess={handleRegisterSuccess} />
      </div>
    </div>
  );
}
