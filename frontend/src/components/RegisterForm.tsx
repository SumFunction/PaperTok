/**
 * 注册表单组件
 * 抖音风格：简洁现代、详细密码强度提示
 */

import { useState, useEffect, useCallback, useMemo, type FormEvent } from 'react';
import { useAuth } from '../contexts/AuthContext';
import './RegisterForm.css';

/**
 * 表单数据接口
 */
interface RegisterFormData {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

/**
 * 表单验证错误接口
 */
interface FormErrors {
  username?: string;
  email?: string;
  password?: string;
  confirmPassword?: string;
}

/**
 * 密码要求项
 */
interface PasswordRequirement {
  id: string;
  label: string;
  check: (password: string) => boolean;
}

interface RegisterFormProps {
  /**
   * 注册成功后的回调函数
   */
  onSuccess?: () => void;
  /**
   * 显示登录链接
   */
  showLoginLink?: boolean;
}

/**
 * 用户名验证规则
 */
const USERNAME_MIN_LENGTH = 3;
const USERNAME_MAX_LENGTH = 20;
const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/;

/**
 * 邮箱格式验证正则
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * 密码要求列表
 */
const PASSWORD_REQUIREMENTS: PasswordRequirement[] = [
  { id: 'length', label: '至少 8 个字符', check: (p) => p.length >= 8 },
  { id: 'upper', label: '包含大写字母', check: (p) => /[A-Z]/.test(p) },
  { id: 'lower', label: '包含小写字母', check: (p) => /[a-z]/.test(p) },
  { id: 'number', label: '包含数字', check: (p) => /[0-9]/.test(p) },
  { id: 'special', label: '包含特殊字符', check: (p) => /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(p) },
];

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
 * 勾选图标
 */
const CheckIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="20 6 9 17 4 12" />
  </svg>
);

/**
 * RegisterForm 组件
 */
export function RegisterForm({ onSuccess, showLoginLink = true }: RegisterFormProps) {
  const { register, loading, error, clearError } = useAuth();

  // 表单状态
  const [formData, setFormData] = useState<RegisterFormData>({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  });

  // 表单验证错误
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  // 表单提交状态
  const [isSubmitting, setIsSubmitting] = useState(false);
  // 密码可见性
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  // 输入框聚焦状态
  const [focusedField, setFocusedField] = useState<string | null>(null);
  // 是否显示密码要求
  const [showPasswordRequirements, setShowPasswordRequirements] = useState(false);

  /**
   * 计算密码要求满足情况
   */
  const passwordChecks = useMemo(() => {
    return PASSWORD_REQUIREMENTS.map((req) => ({
      ...req,
      passed: req.check(formData.password),
    }));
  }, [formData.password]);

  /**
   * 计算满足的要求数量
   */
  const passedCount = useMemo(() => {
    return passwordChecks.filter((c) => c.passed).length;
  }, [passwordChecks]);

  /**
   * 密码是否满足要求（至少长度 + 2种字符类型）
   */
  const isPasswordValid = useMemo(() => {
    const hasLength = passwordChecks[0]?.passed;
    const characterTypesCount = passwordChecks.slice(1).filter((c) => c.passed).length;
    return hasLength && characterTypesCount >= 2;
  }, [passwordChecks]);

  /**
   * 密码强度等级
   */
  const passwordStrength = useMemo((): 'weak' | 'medium' | 'strong' | null => {
    if (!formData.password) return null;
    if (passedCount <= 2) return 'weak';
    if (passedCount <= 3) return 'medium';
    return 'strong';
  }, [passedCount, formData.password]);

  /**
   * 验证用户名
   */
  const validateUsername = useCallback((username: string): boolean => {
    if (!username) {
      setFormErrors((prev) => ({ ...prev, username: '请输入用户名' }));
      return false;
    }
    if (username.length < USERNAME_MIN_LENGTH) {
      setFormErrors((prev) => ({
        ...prev,
        username: `用户名至少需要${USERNAME_MIN_LENGTH}个字符`,
      }));
      return false;
    }
    if (username.length > USERNAME_MAX_LENGTH) {
      setFormErrors((prev) => ({
        ...prev,
        username: `用户名不能超过${USERNAME_MAX_LENGTH}个字符`,
      }));
      return false;
    }
    if (!USERNAME_REGEX.test(username)) {
      setFormErrors((prev) => ({
        ...prev,
        username: '用户名只能包含字母、数字和下划线',
      }));
      return false;
    }
    setFormErrors((prev) => ({ ...prev, username: undefined }));
    return true;
  }, []);

  /**
   * 验证邮箱
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
    // 检查是否满足至少2种字符类型
    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password);
    const typeCount = [hasUpper, hasLower, hasNumber, hasSpecial].filter(Boolean).length;
    
    if (typeCount < 2) {
      setFormErrors((prev) => ({
        ...prev,
        password: '密码需要包含至少2种字符类型（大写、小写、数字、特殊字符）',
      }));
      return false;
    }
    setFormErrors((prev) => ({ ...prev, password: undefined }));
    return true;
  }, []);

  /**
   * 验证确认密码
   */
  const validateConfirmPassword = useCallback((password: string, confirmPassword: string): boolean => {
    if (!confirmPassword) {
      setFormErrors((prev) => ({ ...prev, confirmPassword: '请确认密码' }));
      return false;
    }
    if (password !== confirmPassword) {
      setFormErrors((prev) => ({ ...prev, confirmPassword: '两次输入的密码不一致' }));
      return false;
    }
    setFormErrors((prev) => ({ ...prev, confirmPassword: undefined }));
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
    if (fieldName === 'password') {
      setShowPasswordRequirements(true);
    }
  };

  /**
   * 处理用户名失焦验证
   */
  const handleUsernameBlur = () => {
    setFocusedField(null);
    if (formData.username) {
      validateUsername(formData.username);
    }
  };

  /**
   * 处理邮箱失焦验证
   */
  const handleEmailBlur = () => {
    setFocusedField(null);
    if (formData.email) {
      validateEmail(formData.email);
    }
  };

  /**
   * 处理密码失焦验证
   */
  const handlePasswordBlur = () => {
    setFocusedField(null);
    if (formData.password) {
      validatePassword(formData.password);
    }
    // 延迟隐藏密码要求，让用户看到最终状态
    setTimeout(() => {
      if (document.activeElement?.getAttribute('name') !== 'password') {
        setShowPasswordRequirements(false);
      }
    }, 200);
  };

  /**
   * 处理确认密码失焦验证
   */
  const handleConfirmPasswordBlur = () => {
    setFocusedField(null);
    if (formData.confirmPassword) {
      validateConfirmPassword(formData.password, formData.confirmPassword);
    }
  };

  /**
   * 处理表单提交
   */
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // 清除之前的错误
    clearError();

    // 验证表单
    const isUsernameValid = validateUsername(formData.username);
    const isEmailValid = validateEmail(formData.email);
    const isPwdValid = validatePassword(formData.password);
    const isConfirmPasswordValid = validateConfirmPassword(formData.password, formData.confirmPassword);

    if (!isUsernameValid || !isEmailValid || !isPwdValid || !isConfirmPasswordValid) {
      return;
    }

    setIsSubmitting(true);

    try {
      await register(formData.username, formData.email, formData.password);

      // 注册成功，调用 onSuccess 回调
      // 注意：不在这里直接 navigate，让页面组件或 AuthModal 处理导航逻辑
      // 这样可以避免 React 状态更新异步导致的竞态条件
      if (onSuccess) {
        onSuccess();
      }
      // 如果没有 onSuccess（在 Register 页面使用时），由 Register.tsx 的 useEffect 处理导航
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
    <div className="register-form">
      <div className="register-form__header">
        <h1 className="register-form__title">创建账号</h1>
        <p className="register-form__subtitle">加入 PaperTok，探索精彩论文</p>
      </div>

      {/* API 错误提示 */}
      {error && (
        <div className="register-form__error" role="alert">
          <svg
            className="register-form__error-icon"
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

      <form className="register-form__form" onSubmit={handleSubmit} noValidate>
        {/* 用户名输入 */}
        <div className={`register-form__field ${focusedField === 'username' ? 'register-form__field--focused' : ''}`}>
          <label htmlFor="reg-username" className="register-form__label">
            用户名
          </label>
          <div className="register-form__input-wrapper">
            <input
              id="reg-username"
              name="username"
              type="text"
              autoComplete="username"
              className={`register-form__input ${
                formErrors.username ? 'register-form__input--error' : ''
              } ${formData.username && !formErrors.username ? 'register-form__input--valid' : ''}`}
              placeholder="字母、数字、下划线"
              value={formData.username}
              onChange={handleInputChange}
              onFocus={() => handleFocus('username')}
              onBlur={handleUsernameBlur}
              disabled={isDisabled}
              maxLength={USERNAME_MAX_LENGTH}
              aria-invalid={Boolean(formErrors.username)}
              aria-describedby={formErrors.username ? 'reg-username-error' : undefined}
            />
            {formData.username && !formErrors.username && (
              <span className="register-form__char-count">
                {formData.username.length}/{USERNAME_MAX_LENGTH}
              </span>
            )}
          </div>
          {formErrors.username && (
            <span id="reg-username-error" className="register-form__field-error" role="alert">
              {formErrors.username}
            </span>
          )}
        </div>

        {/* 邮箱输入 */}
        <div className={`register-form__field ${focusedField === 'email' ? 'register-form__field--focused' : ''}`}>
          <label htmlFor="reg-email" className="register-form__label">
            邮箱地址
          </label>
          <div className="register-form__input-wrapper">
            <input
              id="reg-email"
              name="email"
              type="email"
              autoComplete="email"
              className={`register-form__input ${
                formErrors.email ? 'register-form__input--error' : ''
              } ${formData.email && !formErrors.email ? 'register-form__input--valid' : ''}`}
              placeholder="请输入邮箱"
              value={formData.email}
              onChange={handleInputChange}
              onFocus={() => handleFocus('email')}
              onBlur={handleEmailBlur}
              disabled={isDisabled}
              aria-invalid={Boolean(formErrors.email)}
              aria-describedby={formErrors.email ? 'reg-email-error' : undefined}
            />
          </div>
          {formErrors.email && (
            <span id="reg-email-error" className="register-form__field-error" role="alert">
              {formErrors.email}
            </span>
          )}
        </div>

        {/* 密码输入 */}
        <div className={`register-form__field ${focusedField === 'password' ? 'register-form__field--focused' : ''}`}>
          <label htmlFor="reg-password" className="register-form__label">
            密码
          </label>
          <div className="register-form__input-wrapper">
            <input
              id="reg-password"
              name="password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="new-password"
              className={`register-form__input register-form__input--password ${
                formErrors.password ? 'register-form__input--error' : ''
              } ${formData.password && isPasswordValid ? 'register-form__input--valid' : ''}`}
              placeholder="设置密码"
              value={formData.password}
              onChange={handleInputChange}
              onFocus={() => handleFocus('password')}
              onBlur={handlePasswordBlur}
              disabled={isDisabled}
              aria-invalid={Boolean(formErrors.password)}
              aria-describedby="password-requirements"
            />
            <button
              type="button"
              className="register-form__toggle-password"
              onClick={() => setShowPassword(!showPassword)}
              tabIndex={-1}
              aria-label={showPassword ? '隐藏密码' : '显示密码'}
            >
              {showPassword ? <EyeOffIcon /> : <EyeIcon />}
            </button>
          </div>
          
          {/* 密码强度指示条和简化提示 */}
          {formData.password && (
            <>
              <div className="register-form__strength-bar">
                <div 
                  className={`register-form__strength-fill register-form__strength-fill--${passwordStrength}`}
                  style={{ width: `${(passedCount / PASSWORD_REQUIREMENTS.length) * 100}%` }}
                />
              </div>
              <p className="register-form__strength-hint">
                {passwordStrength === 'weak' && '密码强度：弱'}
                {passwordStrength === 'medium' && '密码强度：中'}
                {passwordStrength === 'strong' && '密码强度：强'}
                {!passwordStrength && '请输入密码'}
              </p>
            </>
          )}
          
          {formErrors.password && (
            <span className="register-form__field-error" role="alert">
              {formErrors.password}
            </span>
          )}
        </div>

        {/* 确认密码输入 */}
        <div className={`register-form__field ${focusedField === 'confirmPassword' ? 'register-form__field--focused' : ''}`}>
          <label htmlFor="reg-confirmPassword" className="register-form__label">
            确认密码
          </label>
          <div className="register-form__input-wrapper">
            <input
              id="reg-confirmPassword"
              name="confirmPassword"
              type={showConfirmPassword ? 'text' : 'password'}
              autoComplete="new-password"
              className={`register-form__input register-form__input--password ${
                formErrors.confirmPassword ? 'register-form__input--error' : ''
              } ${formData.confirmPassword && formData.password === formData.confirmPassword ? 'register-form__input--valid' : ''}`}
              placeholder="再次输入密码"
              value={formData.confirmPassword}
              onChange={handleInputChange}
              onFocus={() => handleFocus('confirmPassword')}
              onBlur={handleConfirmPasswordBlur}
              disabled={isDisabled}
              aria-invalid={Boolean(formErrors.confirmPassword)}
              aria-describedby={formErrors.confirmPassword ? 'reg-confirmPassword-error' : undefined}
            />
            <button
              type="button"
              className="register-form__toggle-password"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              tabIndex={-1}
              aria-label={showConfirmPassword ? '隐藏密码' : '显示密码'}
            >
              {showConfirmPassword ? <EyeOffIcon /> : <EyeIcon />}
            </button>
          </div>
          {formErrors.confirmPassword && (
            <span id="reg-confirmPassword-error" className="register-form__field-error" role="alert">
              {formErrors.confirmPassword}
            </span>
          )}
        </div>

        {/* 提交按钮 */}
        <button
          type="submit"
          className={`register-form__submit ${isSubmitting ? 'register-form__submit--loading' : ''}`}
          disabled={isDisabled}
          aria-busy={isDisabled}
        >
          {isSubmitting ? (
            <>
              <span className="register-form__spinner" />
              注册中...
            </>
          ) : (
            '立即注册'
          )}
        </button>

        {/* 温馨提示 */}
        {!showLoginLink && (
          <p className="register-form__switch-hint">
            已有账号？点击上方<strong>「登录」</strong>标签
          </p>
        )}
      </form>

      {/* 登录链接 */}
      {showLoginLink && (
        <p className="register-form__footer">
          已有账号？{' '}
          <a href="/login" className="register-form__link">
            立即登录
          </a>
        </p>
      )}
    </div>
  );
}
