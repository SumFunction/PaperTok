/**
 * 收藏列表页面
 * 展示用户收藏的所有论文
 */

import { useMemo, useState, useEffect } from 'react';
import { PaperCard } from '../components/PaperCard';
import { useAppContext } from '../contexts/AppContext';
import { getUserPreferences } from '../utils/storage';
import { PaperTokAPI } from '../services/api';
import { LoadingList } from '../components/LoadingCard';
import './Favorites.css';

interface FavoritesState {
  papers: any[];
  loading: boolean;
  error: string | null;
}

export function Favorites() {
  const { favoritePapers, toggleLike, toggleFavorite, isLiked } = useAppContext();

  const [state, setState] = useState<FavoritesState>({
    papers: [],
    loading: true,
    error: null,
  });

  // 从 favoritePapers Set 转换为数组
  const favoriteIds = useMemo(() => {
    return Array.from(favoritePapers).reverse(); // 最新的在前
  }, [favoritePapers]);

  // 加载收藏的论文详情
  useEffect(() => {
    const loadFavorites = async () => {
      if (favoriteIds.length === 0) {
        setState({ papers: [], loading: false, error: null });
        return;
      }

      setState((prev) => ({ ...prev, loading: true, error: null }));

      try {
        // 由于后端没有提供根据 ID 批量获取论文的接口
        // 这里我们通过搜索 API 来获取论文（通过 arXiv ID）
        // 注意：这不是最优解，理想情况下应该有专门的批量获取接口

        // 暂时显示空状态，实际应用中可以：
        // 1. 在本地缓存完整的论文数据
        // 2. 后端提供批量查询接口
        // 3. 使用后端的搜索接口逐个查询

        setState({
          papers: [],
          loading: false,
          error: null,
        });
      } catch (err) {
        setState({
          papers: [],
          loading: false,
          error: err instanceof Error ? err.message : '加载失败',
        });
      }
    };

    loadFavorites();
  }, [favoriteIds]);

  // 处理分享
  const handleShare = () => {
    console.log('Shared');
  };

  // 处理标记为已浏览
  const handleView = (id: string) => {
    console.log('Viewed paper:', id);
  };

  // 空状态
  if (!state.loading && favoriteIds.length === 0) {
    return (
      <div className="favorites favorites--empty">
        <div className="favorites__empty">
          <div className="favorites__empty-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
              />
            </svg>
          </div>
          <h2>暂无收藏</h2>
          <p>点击论文卡片上的书签图标收藏论文</p>
        </div>
      </div>
    );
  }

  return (
    <div className="favorites">
      <div className="favorites__container">
        {/* 页面标题 */}
        <div className="favorites__header">
          <h1>我的收藏</h1>
          <p className="favorites__count">{favoriteIds.length} 篇论文</p>
        </div>

        {/* 论文列表 */}
        {state.loading && <LoadingList count={3} />}

        {!state.loading && state.error && (
          <div className="favorites__error">
            <p>加载失败: {state.error}</p>
          </div>
        )}

        {!state.loading && !state.error && (
          <div className="favorites__list">
            {favoriteIds.map((id) => (
              <div key={id} className="favorites__item">
                {/* 暂时显示 ID，实际应用中应该显示完整的论文信息 */}
                <div className="favorites__placeholder">
                  <p>Paper ID: {id}</p>
                  <p className="favorites__placeholder-hint">
                    收藏功能已保存，完整信息将在后续版本中显示
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
