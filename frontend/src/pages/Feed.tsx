/**
 * 推荐流页面
 * 抖音风格：全屏垂直滑动浏览论文
 */

import { useState, useCallback, useEffect, useRef } from 'react';
import { PaperCard } from '../components/PaperCard';
import { CategoryFilter } from '../components/CategoryFilter';
import { useAppContext } from '../contexts/AppContext';
import { useAuth } from '../contexts/AuthContext';
import { usePapers } from '../hooks/usePapers';
import { usePaperKeyboard } from '../hooks/useKeyboard';
import { useViewCounter } from '../hooks/useViewCounter';
import './Feed.css';

export function Feed() {
  const { selectedCategory, toggleLike, toggleFavorite, isLiked, isFavorited } =
    useAppContext();
  const { isAuthenticated, openAuthModal } = useAuth();

  const [toasts, setToasts] = useState<
    Array<{ id: string; message: string; type: 'success' | 'error' | 'info' }>
  >([]);
  const [showSwipeHint, setShowSwipeHint] = useState(true);

  // 使用论文数据 Hook
  const { papers, loading, error, hasMore, fetchPapers } = usePapers(selectedCategory);
  
  // 调试：输出当前状态
  console.log('[Feed] 渲染状态:', { papersCount: papers.length, hasMore, loading, selectedCategory });

  // 使用 ref 存储最新的状态和函数，避免闭包问题
  const stateRef = useRef({ hasMore, loading, selectedCategory, papersCount: papers.length });
  stateRef.current = { hasMore, loading, selectedCategory, papersCount: papers.length };
  
  const fetchPapersRef = useRef(fetchPapers);
  fetchPapersRef.current = fetchPapers;

  // 存储最新的 papers 数组引用，避免滚动事件的闭包问题
  const papersRef = useRef(papers);
  papersRef.current = papers;

  // 当前查看的论文索引
  const [currentIndex, setCurrentIndex] = useState(0);
  const feedRef = useRef<HTMLDivElement>(null);
  const loadMoreRef = useRef<HTMLDivElement>(null);
  
  // 防抖标志，防止短时间内重复请求
  const isLoadingMoreRef = useRef(false);

  // 浏览计数器 - 未登录用户浏览 10 篇论文后强制登录
  const { recordView, viewCount } = useViewCounter({
    threshold: 10,
    onThresholdReached: () => {
      console.log('[Feed] 达到浏览阈值，弹出登录框');
      if (!isAuthenticated) {
        // 使用 'limit_reached' 场景，不可关闭弹窗
        openAuthModal('login', 'limit_reached');
      }
    },
  });

  // 隐藏滑动提示
  useEffect(() => {
    const timer = setTimeout(() => setShowSwipeHint(false), 5000);
    return () => clearTimeout(timer);
  }, []);

  // 记录第一篇论文的浏览（页面加载时）
  const hasRecordedFirstRef = useRef(false);
  useEffect(() => {
    if (papers.length > 0 && !isAuthenticated && !hasRecordedFirstRef.current) {
      hasRecordedFirstRef.current = true;
      const firstPaper = papers[0];
      console.log('[Feed] 记录第一篇论文浏览:', firstPaper.id);
      recordView(firstPaper.id);
    }
  }, [papers, isAuthenticated, recordView]);

  // 存储最新的 recordView 引用，避免闭包问题
  const recordViewRef = useRef(recordView);
  recordViewRef.current = recordView;

  // 存储当前索引的 ref，避免闭包问题
  const currentIndexRef = useRef(currentIndex);
  currentIndexRef.current = currentIndex;

  // 监听滚动更新当前索引 - 使用 ref 避免闭包陷阱
  useEffect(() => {
    const feed = feedRef.current;
    if (!feed) return;

    const handleScroll = () => {
      const scrollTop = feed.scrollTop;
      const itemHeight = feed.clientHeight;
      const newIndex = Math.round(scrollTop / itemHeight);
      const currentPapers = papersRef.current;
      const prevIndex = currentIndexRef.current;
      const { hasMore, loading, selectedCategory } = stateRef.current;

      console.log('[Feed] 滚动事件:', { newIndex, prevIndex, papersCount: currentPapers.length, hasMore, loading });

      // 更新当前索引（限制在有效论文范围内）
      if (newIndex >= 0) {
        const validIndex = Math.min(newIndex, currentPapers.length - 1);
        
        if (validIndex !== prevIndex && validIndex >= 0) {
          setCurrentIndex(validIndex);
          currentIndexRef.current = validIndex;
          setShowSwipeHint(false);

          // 记录论文浏览
          const currentPaper = currentPapers[validIndex];
          if (currentPaper) {
            console.log('[Feed] 记录论文浏览:', currentPaper.id, '索引:', validIndex);
            recordViewRef.current(currentPaper.id);
          }
        }

        // 当滚动到接近末尾时加载（提前 3 篇开始预加载）
        const distanceToEnd = currentPapers.length - 1 - validIndex;
        const shouldLoadMore = distanceToEnd <= 3 &&
                               hasMore && 
                               !loading && 
                               !isLoadingMoreRef.current &&
                               currentPapers.length > 0;

        if (shouldLoadMore) {
          console.log('[Feed] ✅ 触发加载更多！当前索引:', validIndex, '距末尾:', distanceToEnd, '论文数:', currentPapers.length);
          isLoadingMoreRef.current = true;
          fetchPapersRef.current(selectedCategory, true).finally(() => {
            setTimeout(() => {
              isLoadingMoreRef.current = false;
              console.log('[Feed] 加载完成');
            }, 500);
          });
        }
      }
    };

    feed.addEventListener('scroll', handleScroll, { passive: true });
    return () => feed.removeEventListener('scroll', handleScroll);
  }, []); // 空依赖，使用 ref 获取最新值

  // 使用 IntersectionObserver 监听"加载更多"元素，确保触发加载
  useEffect(() => {
    const loadMoreEl = loadMoreRef.current;
    if (!loadMoreEl) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting) {
          const { hasMore, loading, selectedCategory } = stateRef.current;
          console.log('[Feed] IntersectionObserver 触发:', { hasMore, loading, isLoadingMore: isLoadingMoreRef.current });
          
          if (hasMore && !loading && !isLoadingMoreRef.current) {
            console.log('[Feed] ✅ IntersectionObserver 触发加载更多！');
            isLoadingMoreRef.current = true;
            fetchPapersRef.current(selectedCategory, true).finally(() => {
              setTimeout(() => {
                isLoadingMoreRef.current = false;
                console.log('[Feed] IntersectionObserver 加载完成');
              }, 500);
            });
          }
        }
      },
      {
        root: feedRef.current,
        threshold: 0.1, // 10% 可见时触发
      }
    );

    observer.observe(loadMoreEl);
    console.log('[Feed] IntersectionObserver 已设置，监听加载更多元素');
    return () => observer.disconnect();
  }, [papers.length, hasMore, loading]); // 当论文数量、hasMore 或 loading 变化时重新设置 observer

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
    const currentPapers = papersRef.current;
    if (!feed || index < 0 || index >= currentPapers.length) return;
    
    feed.scrollTo({
      top: index * feed.clientHeight,
      behavior: 'smooth',
    });
  }, []); // 使用 ref 获取最新值，无需依赖

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

        {/* 加载更多触发区域 - 自动加载，滚动到这里时会触发加载 */}
        {hasMore && !loading && (
          <div ref={loadMoreRef} className="feed__item feed__load-more">
            <div className="feed__load-more-content">
              <div className="feed__spinner" />
              <p>正在加载更多...</p>
            </div>
          </div>
        )}

        {/* 加载中 */}
        {loading && (
          <div className="feed__item feed__loading-item">
            <div className="feed__loading">
              <div className="feed__spinner" />
              <p>加载更多论文...</p>
            </div>
          </div>
        )}

        {/* 没有更多内容 */}
        {!hasMore && papers.length > 0 && (
          <div className="feed__item feed__end">
            <div className="feed__end-content">
              <div className="feed__end-icon">🎉</div>
              <p>你已经看完了所有论文</p>
              <p style={{ fontSize: '0.85rem', opacity: 0.7 }}>换个分类继续探索吧</p>
            </div>
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
