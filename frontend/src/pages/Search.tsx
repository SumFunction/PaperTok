/**
 * 搜索页面
 * 允许用户搜索论文
 */

import { useState, useCallback, useEffect, type ChangeEvent } from 'react';
import { PaperCard } from '../components/PaperCard';
import { LoadingList } from '../components/LoadingCard';
import { useSearch } from '../hooks/useSearch';
import { useAppContext } from '../contexts/AppContext';
import './Search.css';

export function Search() {
  const { toggleLike, toggleFavorite, isLiked, isFavorited } = useAppContext();

  const [searchQuery, setSearchQuery] = useState<string>('');
  const [debouncedQuery, setDebouncedQuery] = useState<string>('');
  const { searchResults, searching, searchError, searchPapers, clearSearch } = useSearch();

  // 防抖处理搜索输入
  useEffect(() => {
    const timer = setTimeout(() => {
      if (debouncedQuery.trim()) {
        searchPapers(debouncedQuery);
      } else {
        clearSearch();
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [debouncedQuery, searchPapers, clearSearch]);

  // 处理搜索输入
  const handleSearchChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchQuery(value);
    setDebouncedQuery(value);
  };

  // 处理分享
  const handleShare = useCallback(() => {
    // 可以在这里显示 Toast 提示
    console.log('Shared');
  }, []);

  // 处理标记为已浏览
  const handleView = useCallback((id: string) => {
    console.log('Viewed paper:', id);
  }, []);

  return (
    <div className="search">
      <div className="search__container">
        {/* 搜索框 */}
        <div className="search__header">
          <div className="search__input-wrapper">
            <svg
              className="search__icon"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
            <input
              type="text"
              className="search__input"
              placeholder="搜索论文标题、作者或摘要..."
              value={searchQuery}
              onChange={handleSearchChange}
              autoFocus
            />
            {searchQuery && (
              <button
                className="search__clear"
                onClick={() => {
                  setSearchQuery('');
                  setDebouncedQuery('');
                  clearSearch();
                }}
              >
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            )}
          </div>
        </div>

        {/* 搜索结果 */}
        <div className="search__results">
          {searching && <LoadingList count={2} />}

          {!searching && searchError && (
            <div className="search__error">
              <p>搜索失败: {searchError}</p>
            </div>
          )}

          {!searching && !searchError && debouncedQuery && searchResults.length === 0 && (
            <div className="search__empty">
              <p>未找到相关论文</p>
              <p className="search__empty-hint">尝试使用不同的关键词搜索</p>
            </div>
          )}

          {!searching && searchResults.length > 0 && (
            <>
              <div className="search__result-count">
                找到 {searchResults.length} 篇论文
              </div>
              {searchResults.map((paper) => (
                <div key={paper.id} className="search__item">
                  <PaperCard
                    paper={paper}
                    isLiked={isLiked(paper.id)}
                    isFavorited={isFavorited(paper.id)}
                    onLike={toggleLike}
                    onFavorite={toggleFavorite}
                    onShare={handleShare}
                    onView={handleView}
                  />
                </div>
              ))}
            </>
          )}

          {!searching && !debouncedQuery && (
            <div className="search__placeholder">
              <div className="search__placeholder-icon">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </div>
              <h3>搜索论文</h3>
              <p>输入关键词搜索论文标题、作者或摘要</p>
              <div className="search__examples">
                <p>示例:</p>
                <ul>
                  <li>transformer architecture</li>
                  <li>computer vision</li>
                  <li>reinforcement learning</li>
                </ul>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
