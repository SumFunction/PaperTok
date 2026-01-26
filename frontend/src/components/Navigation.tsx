/**
 * 导航栏组件
 * 提供页面之间的快速切换
 */

import { Link, useLocation } from 'react-router-dom';
import './Navigation.css';

export function Navigation() {
  const location = useLocation();

  const navItems = [
    { path: '/', label: '推荐', icon: '🏠' },
    { path: '/search', label: '搜索', icon: '🔍' },
    { path: '/favorites', label: '收藏', icon: '⭐' },
  ];

  return (
    <nav className="navigation">
      <div className="navigation__container">
        <div className="navigation__logo">
          <Link to="/">
            <span className="navigation__logo-icon">📄</span>
            <span className="navigation__logo-text">PaperTok</span>
          </Link>
        </div>

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
        </ul>
      </div>
    </nav>
  );
}
