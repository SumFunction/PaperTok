/**
 * PaperTok 主应用组件
 * 配置路由和全局状态管理
 */

import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AppProvider } from './contexts/AppContext';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { Navigation } from './components/Navigation';
import { ProtectedRoute } from './components/ProtectedRoute';
import { AuthModal } from './components/AuthModal';
import { Feed } from './pages/Feed';
import { Search } from './pages/Search';
import { Favorites } from './pages/Favorites';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import './index.css';

/**
 * 强制登录提示 Hook
 * 未登录用户进入时立即弹出登录框
 */
function useLoginPrompt() {
  const { isAuthenticated, openAuthModal, loading } = useAuth();

  useEffect(() => {
    console.log('[LoginPrompt] 状态:', { loading, isAuthenticated });

    // 等待认证状态加载完成
    if (loading) {
      console.log('[LoginPrompt] 等待认证加载...');
      return;
    }

    // 如果已登录，不弹窗
    if (isAuthenticated) {
      console.log('[LoginPrompt] 用户已登录，不弹窗');
      return;
    }

    // 未登录，延迟后弹出登录提示（使用 welcome 场景）
    console.log('[LoginPrompt] 用户未登录，准备弹出欢迎登录框');
    const timer = setTimeout(() => {
      console.log('[LoginPrompt] 弹出登录框');
      openAuthModal('login', 'welcome');
    }, 500); // 增加延迟，确保页面已渲染

    return () => clearTimeout(timer);
  }, [isAuthenticated, loading, openAuthModal]);
}

/**
 * 应用内部组件（使用 Auth 后的数据）
 */
function AppContent() {
  useLoginPrompt();

  return (
    <Router>
      <div className="app">
        <Navigation />
        <main className="app__main">
          <Routes>
            {/* 公共路由 */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

            {/* 受保护路由 */}
            <Route
              path="/favorites"
              element={
                <ProtectedRoute>
                  <Favorites />
                </ProtectedRoute>
              }
            />

            {/* 公开路由（也可以在登录后访问） */}
            <Route path="/" element={<Feed />} />
            <Route path="/search" element={<Search />} />
          </Routes>
        </main>
        {/* 认证模态框（全局） */}
        <AuthModal />
      </div>
    </Router>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppProvider>
        <AppContent />
      </AppProvider>
    </AuthProvider>
  );
}

export default App;
