/**
 * 推荐流页面
 * 抖音风格：全屏垂直滑动浏览论文
 */

import { useState, useCallback, useEffect, useRef } from 'react';
import { PaperCard } from '../components/PaperCard';
import { CategoryFilter } from '../components/CategoryFilter';
import { useAppContext } from '../contexts/AppContext';
import { usePapers } from '../hooks/usePapers';
import { usePaperKeyboard } from '../hooks/useKeyboard';
import './Feed.css';

export function Feed() {
  const { selectedCategory, toggleLike, toggleFavorite, isLiked, isFavorited } =
    useAppContext();

  const [toasts, setToasts] = useState<
    Array<{ id: string; message: string; type: 'success' | 'error' | 'info' }>
  >([]);
  const [showSwipeHint, setShowSwipeHint] = useState(true);

  // 使用论文数据 Hook
  const { papers, loading, error, hasMore, fetchPapers } = usePapers(selectedCategory);

  // 当前查看的论文索引
  const [currentIndex, setCurrentIndex] = useState(0);
  const feedRef = useRef<HTMLDivElement>(null);
  const observerTarget = useRef<HTMLDivElement>(null);

  // 隐藏滑动提示
  useEffect(() => {
    const timer = setTimeout(() => setShowSwipeHint(false), 5000);
    return () => clearTimeout(timer);
  }, []);

  // 监听滚动更新当前索引
  useEffect(() => {
    const feed = feedRef.current;
    if (!feed) return;

    const handleScroll = () => {
      const scrollTop = feed.scrollTop;
      const itemHeight = feed.clientHeight;
      const newIndex = Math.round(scrollTop / itemHeight);
      
      if (newIndex !== currentIndex && newIndex >= 0 && newIndex < papers.length) {
        setCurrentIndex(newIndex);
        // 首次滑动后隐藏提示
        if (showSwipeHint) setShowSwipeHint(false);
      }
    };

    feed.addEventListener('scroll', handleScroll, { passive: true });
    return () => feed.removeEventListener('scroll', handleScroll);
  }, [currentIndex, papers.length, showSwipeHint]);

  // 显示 Toast 提示
  const showToast = useCallback(
    (message: string, type: 'success' | 'error' | 'info' = 'success') => {
      const id = Date.now().toString();
      setToasts((prev) => [...prev, { id, message, type }]);
      // 自动消失
      setTimeout(() => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
      }, 2000);
    },
    []
  );

  // 滚动到指定索引
  const scrollToIndex = useCallback((index: number) => {
    const feed = feedRef.current;
    if (!feed || index < 0 || index >= papers.length) return;
    
    feed.scrollTo({
      top: index * feed.clientHeight,
      behavior: 'smooth',
    });
  }, [papers.length]);

  // 处理分享
  const handleShare = useCallback(() => {
    showToast('已复制到剪贴板', 'success');
  }, [showToast]);

  // 处理分类切换
  const handleCategoryChange = useCallback(
    (category: string) => {
      setCurrentIndex(0);
      // 滚动回顶部
      if (feedRef.current) {
        feedRef.current.scrollTo({ top: 0 });
      }
      fetchPapers(category, false);
    },
    [fetchPapers]
  );

  // 处理点赞
  const handleLike = useCallback((id: string) => {
    const wasLiked = isLiked(id);
    toggleLike(id);
    showToast(wasLiked ? '已取消点赞' : '已点赞 ❤️', 'success');
  }, [isLiked, toggleLike, showToast]);

  // 处理收藏
  const handleFavorite = useCallback((id: string) => {
    const wasFavorited = isFavorited(id);
    toggleFavorite(id);
    showToast(wasFavorited ? '已取消收藏' : '已收藏 ⭐', 'success');
  }, [isFavorited, toggleFavorite, showToast]);

  // 键盘快捷键
  usePaperKeyboard({
    onNext: () => scrollToIndex(currentIndex + 1),
    onPrevious: () => scrollToIndex(currentIndex - 1),
    onLike: () => {
      if (papers[currentIndex]) {
        handleLike(papers[currentIndex].id);
      }
    },
    onFavorite: () => {
      if (papers[currentIndex]) {
        handleFavorite(papers[currentIndex].id);
      }
    },
    onExpand: () => {
      // 可以在这里实现展开摘要的功能
    },
  });

  // 无限滚动加载
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading) {
          fetchPapers(selectedCategory, true);
        }
      },
      { threshold: 0.1 }
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [hasMore, loading, fetchPapers, selectedCategory]);

  // 错误处理
  if (error && papers.length === 0) {
    return (
      <div className="feed feed--error">
        <div className="feed__error">
          <div className="feed__error-icon">😵</div>
          <h2>加载失败</h2>
          <p>{error}</p>
          <button 
            onClick={() => fetchPapers(selectedCategory, false)} 
            className="feed__retry-btn"
          >
            重新加载
          </button>
        </div>
      </div>
    );
  }

  // 空状态
  if (!loading && papers.length === 0) {
    return (
      <div className="feed feed--empty">
        <div className="feed__empty">
          <div className="feed__empty-icon">📭</div>
          <h2>暂无论文</h2>
          <p>该分类下暂时没有论文，换个分类试试？</p>
        </div>
      </div>
    );
  }

  return (
    <div className="feed" ref={feedRef}>
      {/* 分类筛选器 */}
      <CategoryFilter
        selectedCategory={selectedCategory}
        onCategoryChange={handleCategoryChange}
      />

      {/* 论文列表 */}
      <div className="feed__list">
        {papers.map((paper, index) => (
          <div key={paper.id} className="feed__item">
            <PaperCard
              paper={paper}
              isLiked={isLiked(paper.id)}
              isFavorited={isFavorited(paper.id)}
              onLike={handleLike}
              onFavorite={handleFavorite}
              onShare={handleShare}
            />
          </div>
        ))}

        {/* 加载中 */}
        {loading && (
          <div className="feed__loading">
            <div className="feed__spinner" />
          </div>
        )}

        {/* 无限滚动触发器 */}
        <div ref={observerTarget} className="feed__observer" />

        {/* 没有更多内容 */}
        {!hasMore && papers.length > 0 && (
          <div className="feed__end">
            <div className="feed__end-icon">🎉</div>
            <p>你已经看完了所有论文</p>
            <p style={{ fontSize: '0.85rem', opacity: 0.7 }}>换个分类继续探索吧</p>
          </div>
        )}
      </div>

      {/* 滑动提示 */}
      {showSwipeHint && papers.length > 1 && (
        <div className="feed__swipe-hint">
          <span className="feed__swipe-hint-icon">↓</span>
          <span>上滑查看下一篇</span>
        </div>
      )}

      {/* Toast 提示 */}
      <div className="feed__toasts">
        {toasts.map((toast) => (
          <div key={toast.id} className={`toast toast--${toast.type}`}>
            {toast.message}
          </div>
        ))}
      </div>
    </div>
  );
}
