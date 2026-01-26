/**
 * PaperTok 主应用组件
 * 配置路由和全局状态管理
 */

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AppProvider } from './contexts/AppContext';
import { Navigation } from './components/Navigation';
import { Feed } from './pages/Feed';
import { Search } from './pages/Search';
import { Favorites } from './pages/Favorites';
import './index.css';

function App() {
  return (
    <AppProvider>
      <Router>
        <div className="app">
          <Navigation />
          <main className="app__main">
            <Routes>
              <Route path="/" element={<Feed />} />
              <Route path="/search" element={<Search />} />
              <Route path="/favorites" element={<Favorites />} />
            </Routes>
          </main>
        </div>
      </Router>
    </AppProvider>
  );
}

export default App;
