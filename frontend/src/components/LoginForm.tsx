/**
 * 登录表单组件
 * 抖音风格：简洁现代、流畅动画
 */

import { useState, useEffect, useCallback, type FormEvent } from 'react';
import { useAuth } from '../contexts/AuthContext';
import './LoginForm.css';

/**
 * 表单数据接口
 */
interface LoginFormData {
  email: string;
  password: string;
}

/**
 * 表单验证错误接口
 */
interface FormErrors {
  email?: string;
  password?: string;
}

interface LoginFormProps {
  /**
   * 登录成功后的回调函数
   */
  onSuccess?: () => void;
  /**
   * 显示注册链接
   */
  showRegisterLink?: boolean;
}

/**
 * 邮箱格式验证正则
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * 眼睛图标 - 显示密码
 */
const EyeIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
);

/**
 * 眼睛关闭图标 - 隐藏密码
 */
const EyeOffIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
    <line x1="1" y1="1" x2="23" y2="23" />
  </svg>
);

/**
 * LoginForm 组件
 */
export function LoginForm({ onSuccess, showRegisterLink = true }: LoginFormProps) {
  const { login, loading, error, clearError } = useAuth();

  // 表单状态
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  });

  // 表单验证错误
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  // 表单提交状态
  const [isSubmitting, setIsSubmitting] = useState(false);
  // 密码可见性
  const [showPassword, setShowPassword] = useState(false);
  // 输入框聚焦状态
  const [focusedField, setFocusedField] = useState<string | null>(null);

  /**
   * 验证邮箱格式
   */
  const validateEmail = useCallback((email: string): boolean => {
    if (!email) {
      setFormErrors((prev) => ({ ...prev, email: '请输入邮箱地址' }));
      return false;
    }
    if (!EMAIL_REGEX.test(email)) {
      setFormErrors((prev) => ({ ...prev, email: '请输入有效的邮箱地址' }));
      return false;
    }
    setFormErrors((prev) => ({ ...prev, email: undefined }));
    return true;
  }, []);

  /**
   * 验证密码
   */
  const validatePassword = useCallback((password: string): boolean => {
    if (!password) {
      setFormErrors((prev) => ({ ...prev, password: '请输入密码' }));
      return false;
    }
    if (password.length < 8) {
      setFormErrors((prev) => ({ ...prev, password: '密码至少需要8个字符' }));
      return false;
    }
    setFormErrors((prev) => ({ ...prev, password: undefined }));
    return true;
  }, []);

  /**
   * 处理输入变化
   */
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));

    // 清除对应字段的错误
    setFormErrors((prev) => ({ ...prev, [name]: undefined }));

    // 清除 API 错误
    if (error) {
      clearError();
    }
  };

  /**
   * 处理输入框聚焦
   */
  const handleFocus = (fieldName: string) => {
    setFocusedField(fieldName);
  };

  /**
   * 处理邮箱输入失焦验证
   */
  const handleEmailBlur = () => {
    setFocusedField(null);
    if (formData.email) {
      validateEmail(formData.email);
    }
  };

  /**
   * 处理密码输入失焦验证
   */
  const handlePasswordBlur = () => {
    setFocusedField(null);
    if (formData.password) {
      validatePassword(formData.password);
    }
  };

  /**
   * 切换密码可见性
   */
  const togglePasswordVisibility = () => {
    setShowPassword((prev) => !prev);
  };

  /**
   * 处理表单提交
   */
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // 清除之前的错误
    clearError();

    // 验证表单
    const isEmailValid = validateEmail(formData.email);
    const isPasswordValid = validatePassword(formData.password);

    if (!isEmailValid || !isPasswordValid) {
      return;
    }

    setIsSubmitting(true);

    try {
      await login(formData.email, formData.password);

      // 登录成功，调用 onSuccess 回调
      // 注意：不在这里直接 navigate，让页面组件或 AuthModal 处理导航逻辑
      // 这样可以避免 React 状态更新异步导致的竞态条件
      if (onSuccess) {
        onSuccess();
      }
      // 如果没有 onSuccess（在 Login 页面使用时），由 Login.tsx 的 useEffect 处理导航
    } catch {
      // 错误已在 AuthContext 中处理
    } finally {
      setIsSubmitting(false);
    }
  };

  // 组件卸载时清除错误
  useEffect(() => {
    return () => {
      clearError();
    };
  }, [clearError]);

  const isDisabled = loading || isSubmitting;

  return (
    <div className="login-form">
      <div className="login-form__header">
        <h1 className="login-form__title">欢迎回来</h1>
        <p className="login-form__subtitle">登录 PaperTok，探索精彩论文</p>
      </div>

      {/* API 错误提示 */}
      {error && (
        <div className="login-form__error" role="alert">
          <svg
            className="login-form__error-icon"
            width="18"
            height="18"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clipRule="evenodd"
            />
          </svg>
          <span>{error}</span>
        </div>
      )}

      <form className="login-form__form" onSubmit={handleSubmit} noValidate>
        {/* 邮箱输入 */}
        <div className={`login-form__field ${focusedField === 'email' ? 'login-form__field--focused' : ''}`}>
          <label htmlFor="login-email" className="login-form__label">
            邮箱地址
          </label>
          <div className="login-form__input-wrapper">
            <input
              id="login-email"
              name="email"
              type="email"
              autoComplete="email"
              className={`login-form__input ${
                formErrors.email ? 'login-form__input--error' : ''
              } ${formData.email && !formErrors.email ? 'login-form__input--valid' : ''}`}
              placeholder="请输入邮箱"
              value={formData.email}
              onChange={handleInputChange}
              onFocus={() => handleFocus('email')}
              onBlur={handleEmailBlur}
              disabled={isDisabled}
              aria-invalid={Boolean(formErrors.email)}
              aria-describedby={formErrors.email ? 'login-email-error' : undefined}
            />
          </div>
          {formErrors.email && (
            <span id="login-email-error" className="login-form__field-error" role="alert">
              {formErrors.email}
            </span>
          )}
        </div>

        {/* 密码输入 */}
        <div className={`login-form__field ${focusedField === 'password' ? 'login-form__field--focused' : ''}`}>
          <label htmlFor="login-password" className="login-form__label">
            密码
          </label>
          <div className="login-form__input-wrapper">
            <input
              id="login-password"
              name="password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="current-password"
              className={`login-form__input login-form__input--password ${
                formErrors.password ? 'login-form__input--error' : ''
              }`}
              placeholder="请输入密码"
              value={formData.password}
              onChange={handleInputChange}
              onFocus={() => handleFocus('password')}
              onBlur={handlePasswordBlur}
              disabled={isDisabled}
              aria-invalid={Boolean(formErrors.password)}
              aria-describedby={formErrors.password ? 'login-password-error' : undefined}
            />
            <button
              type="button"
              className="login-form__toggle-password"
              onClick={togglePasswordVisibility}
              tabIndex={-1}
              aria-label={showPassword ? '隐藏密码' : '显示密码'}
            >
              {showPassword ? <EyeOffIcon /> : <EyeIcon />}
            </button>
          </div>
          {formErrors.password && (
            <span id="login-password-error" className="login-form__field-error" role="alert">
              {formErrors.password}
            </span>
          )}
        </div>

        {/* 提交按钮 */}
        <button
          type="submit"
          className={`login-form__submit ${isSubmitting ? 'login-form__submit--loading' : ''}`}
          disabled={isDisabled}
          aria-busy={isDisabled}
        >
          {isSubmitting ? (
            <>
              <span className="login-form__spinner" />
              登录中...
            </>
          ) : (
            '立即登录'
          )}
        </button>

        {/* 温馨提示 */}
        {!showRegisterLink && (
          <p className="login-form__switch-hint">
            还没有账号？点击上方<strong>「注册」</strong>标签
          </p>
        )}
      </form>

      {/* 注册链接 */}
      {showRegisterLink && (
        <p className="login-form__footer">
          还没有账号？{' '}
          <a href="/register" className="login-form__link">
            立即注册
          </a>
        </p>
      )}
    </div>
  );
}
