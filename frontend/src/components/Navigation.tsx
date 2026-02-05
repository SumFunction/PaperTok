/**
 * 导航栏组件
 * 提供页面之间的快速切换
 * 根据登录状态显示不同的导航选项
 */

import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import './Navigation.css';

export function Navigation() {
  const location = useLocation();
  const { isAuthenticated, user, logout, loading, openAuthModal } = useAuth();

  // 公共导航项（所有用户可见）
  const publicNavItems = [
    { path: '/', label: '推荐', icon: '🏠' },
    { path: '/search', label: '搜索', icon: '🔍' },
  ];

  // 已登录用户额外的导航项
  const authenticatedNavItems = [
    ...publicNavItems,
    { path: '/favorites', label: '收藏', icon: '⭐' },
  ];

  // 根据登录状态选择导航项
  const navItems = isAuthenticated ? authenticatedNavItems : publicNavItems;

  /**
   * 处理登录按钮点击 - 打开 AuthModal 而不是跳转页面
   */
  const handleLoginClick = () => {
    openAuthModal('login', 'welcome');
  };

  /**
   * 处理登出
   */
  const handleLogout = async () => {
    await logout();
    // 登出后跳转到首页
    window.location.href = '/';
  };

  return (
    <nav className="navigation">
      <div className="navigation__container">
        {/* Logo */}
        <div className="navigation__logo">
          <Link to="/">
            <span className="navigation__logo-icon">📄</span>
            <span className="navigation__logo-text">PaperTok</span>
          </Link>
        </div>

        {/* 导航菜单 */}
        <ul className="navigation__menu">
          {navItems.map((item) => (
            <li key={item.path} className="navigation__item">
              <Link
                to={item.path}
                className={`navigation__link ${
                  location.pathname === item.path ? 'navigation__link--active' : ''
                }`}
              >
                <span className="navigation__icon">{item.icon}</span>
                <span className="navigation__label">{item.label}</span>
              </Link>
            </li>
          ))}

          {/* 未登录时显示登录按钮 */}
          {!isAuthenticated && (
            <li className="navigation__item">
              <button
                className="navigation__link navigation__login-btn"
                onClick={handleLoginClick}
                type="button"
              >
                <span className="navigation__icon">👤</span>
                <span className="navigation__label">登录</span>
              </button>
            </li>
          )}

          {/* 已登录用户信息 */}
          {isAuthenticated && (
            <li className="navigation__item navigation__item--user">
              <div className="navigation__user">
                <span className="navigation__user-name">{user?.username || '用户'}</span>
                <button
                  className="navigation__logout"
                  onClick={handleLogout}
                  disabled={loading}
                  aria-label="登出"
                >
                  登出
                </button>
              </div>
            </li>
          )}
        </ul>
      </div>
    </nav>
  );
}
