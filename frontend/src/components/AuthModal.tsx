/**
 * 认证模态框组件
 * 抖音风格：全屏半透明遮罩 + 中心卡片式登录/注册表单
 * 支持 Tab 切换登录/注册，根据场景显示不同提示
 */

import { useState, useEffect, useCallback, useRef, type MouseEvent, type KeyboardEvent } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { LoginForm } from './LoginForm';
import { RegisterForm } from './RegisterForm';
import type { AuthModalTab, AuthModalScene } from '../contexts/AuthContext';
import './AuthModal.css';

export interface AuthModalProps {
  /**
   * 模态框打开时的回调
   */
  onOpen?: () => void;
  /**
   * 模态框关闭时的回调
   */
  onClose?: () => void;
}

/**
 * 根据场景获取提示文案
 */
function getSceneConfig(scene: AuthModalScene) {
  switch (scene) {
    case 'welcome':
      return {
        title: '欢迎来到 PaperTok',
        subtitle: '登录后开启精彩论文之旅',
        closable: true, // 可关闭，关闭后可免费浏览 10 篇
      };
    case 'limit_reached':
      return {
        title: '免费浏览已达上限',
        subtitle: '登录后即可无限浏览所有论文',
        closable: false, // 不可关闭，必须登录
      };
    case 'protected_route':
      return {
        title: '需要登录',
        subtitle: '该功能需要登录后才能使用',
        closable: true,
      };
    default:
      return {
        title: '登录 PaperTok',
        subtitle: '探索精彩论文世界',
        closable: true,
      };
  }
}

/**
 * AuthModal 组件
 */
export function AuthModal({ onOpen, onClose }: AuthModalProps) {
  const { showAuthModal, authModalDefaultTab, authModalScene, closeAuthModal } = useAuth();
  const contentRef = useRef<HTMLDivElement>(null);
  
  // 获取场景配置
  const sceneConfig = getSceneConfig(authModalScene);

  // 当前激活的 Tab
  const [activeTab, setActiveTab] = useState<AuthModalTab>(authModalDefaultTab);
  // 关闭动画状态
  const [isClosing, setIsClosing] = useState(false);

  // 同步外部控制的默认 Tab
  useEffect(() => {
    if (showAuthModal) {
      setActiveTab(authModalDefaultTab);
      setIsClosing(false);
      onOpen?.();
    } else {
      onClose?.();
    }
  }, [showAuthModal, authModalDefaultTab, onOpen, onClose]);

  // 关闭动画处理
  const handleClose = useCallback(() => {
    setIsClosing(true);
    setTimeout(() => {
      closeAuthModal();
      setIsClosing(false);
    }, 200);
  }, [closeAuthModal]);

  // 登录成功后的处理
  const handleLoginSuccess = useCallback(() => {
    handleClose();
  }, [handleClose]);

  // 注册成功后的处理
  const handleRegisterSuccess = useCallback(() => {
    handleClose();
  }, [handleClose]);

  // 点击遮罩关闭（仅当场景允许关闭时）
  const handleBackdropClick = useCallback((e: MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget && sceneConfig.closable) {
      handleClose();
    }
  }, [handleClose, sceneConfig.closable]);

  // ESC 键关闭（仅当场景允许关闭时）
  const handleKeyDown = useCallback((e: KeyboardEvent<HTMLDivElement>) => {
    if (e.key === 'Escape' && sceneConfig.closable) {
      handleClose();
    }
  }, [handleClose, sceneConfig.closable]);

  // 切换 Tab
  const handleTabChange = useCallback((tab: AuthModalTab) => {
    setActiveTab(tab);
  }, []);

  // 防止背景滚动
  useEffect(() => {
    if (showAuthModal) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [showAuthModal]);

  if (!showAuthModal) {
    return null;
  }

  return (
    <div
      className={`auth-modal ${isClosing ? 'auth-modal--closing' : ''}`}
      onClick={handleBackdropClick}
      onKeyDown={handleKeyDown}
      role="dialog"
      aria-modal="true"
      aria-labelledby="auth-modal-title"
      tabIndex={-1}
    >
      <div 
        ref={contentRef}
        className={`auth-modal__content ${isClosing ? 'auth-modal__content--closing' : ''}`}
      >
        {/* 关闭按钮（仅当场景允许关闭时显示） */}
        {sceneConfig.closable && (
          <button
            className="auth-modal__close"
            onClick={handleClose}
            aria-label="关闭"
            type="button"
          >
            <svg
              width="22"
              height="22"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        )}

        {/* Logo 和场景提示 */}
        <div className="auth-modal__logo">
          <span className="auth-modal__logo-icon">📄</span>
          <span className="auth-modal__logo-text">PaperTok</span>
        </div>
        
        {/* 场景标题和副标题 */}
        <div className="auth-modal__scene-header">
          <h2 className="auth-modal__scene-title">{sceneConfig.title}</h2>
          <p className="auth-modal__scene-subtitle">{sceneConfig.subtitle}</p>
        </div>

        {/* Tab 切换 */}
        <div className="auth-modal__tabs">
          <button
            className={`auth-modal__tab ${activeTab === 'login' ? 'auth-modal__tab--active' : ''}`}
            onClick={() => handleTabChange('login')}
            type="button"
            aria-label="切换到登录"
          >
            登录
          </button>
          <button
            className={`auth-modal__tab ${activeTab === 'register' ? 'auth-modal__tab--active' : ''}`}
            onClick={() => handleTabChange('register')}
            type="button"
            aria-label="切换到注册"
          >
            注册
          </button>
          <div 
            className="auth-modal__tab-indicator" 
            style={{ transform: `translateX(${activeTab === 'login' ? '0' : '100%'})` }}
            aria-hidden="true"
          />
        </div>

        {/* 表单内容 */}
        <div className="auth-modal__body">
          <div 
            className="auth-modal__forms" 
            style={{ transform: `translateX(${activeTab === 'login' ? '0' : '-50%'})` }}
          >
            <div className="auth-modal__form-panel">
              <LoginForm
                onSuccess={handleLoginSuccess}
                showRegisterLink={false}
              />
            </div>
            <div className="auth-modal__form-panel">
              <RegisterForm
                onSuccess={handleRegisterSuccess}
                showLoginLink={false}
              />
            </div>
          </div>
        </div>

        {/* 底部提示 */}
        <div className="auth-modal__footer">
          <p className="auth-modal__terms">
            继续即表示您同意 PaperTok 的
            <a href="/terms" className="auth-modal__terms-link">服务条款</a>
            和
            <a href="/privacy" className="auth-modal__terms-link">隐私政策</a>
          </p>
        </div>
      </div>
    </div>
  );
}

/**
 * 导出一个简化的 Hook 用于快速打开模态框
 */
export function useAuthModal() {
  const { openAuthModal, closeAuthModal, showAuthModal } = useAuth();

  return {
    showAuthModal,
    openLogin: () => openAuthModal('login'),
    openRegister: () => openAuthModal('register'),
    closeAuthModal,
  };
}
