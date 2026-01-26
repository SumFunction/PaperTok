/**
 * 本地存储 Hook
 */

import { useState, useEffect, useCallback } from 'react';
import { getUserPreferences, saveUserPreferences } from '../utils/storage';

/**
 * 管理点赞列表
 */
export function useLikedPapers() {
  const [likedPapers, setLikedPapers] = useState<Set<string>>(new Set());

  // 初始化加载
  useEffect(() => {
    const preferences = getUserPreferences();
    setLikedPapers(new Set(preferences.likedPapers));
  }, []);

  // 切换点赞状态
  const toggleLike = useCallback((paperId: string) => {
    setLikedPapers((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(paperId)) {
        newSet.delete(paperId);
      } else {
        newSet.add(paperId);
      }

      // 保存到本地存储
      saveUserPreferences({
        likedPapers: Array.from(newSet),
      });

      return newSet;
    });
  }, []);

  // 检查是否已点赞
  const isLiked = useCallback(
    (paperId: string) => {
      return likedPapers.has(paperId);
    },
    [likedPapers]
  );

  return {
    likedPapers,
    toggleLike,
    isLiked,
  };
}

/**
 * 管理收藏列表
 */
export function useFavoritePapers() {
  const [favoritePapers, setFavoritePapers] = useState<Set<string>>(new Set());

  // 初始化加载
  useEffect(() => {
    const preferences = getUserPreferences();
    setFavoritePapers(new Set(preferences.favoritePapers));
  }, []);

  // 切换收藏状态
  const toggleFavorite = useCallback((paperId: string) => {
    setFavoritePapers((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(paperId)) {
        newSet.delete(paperId);
      } else {
        newSet.add(paperId);
      }

      // 保存到本地存储
      saveUserPreferences({
        favoritePapers: Array.from(newSet),
      });

      return newSet;
    });
  }, []);

  // 检查是否已收藏
  const isFavorited = useCallback(
    (paperId: string) => {
      return favoritePapers.has(paperId);
    },
    [favoritePapers]
  );

  return {
    favoritePapers,
    toggleFavorite,
    isFavorited,
  };
}

/**
 * 管理浏览历史
 */
export function useViewedPapers() {
  const [viewedPapers, setViewedPapers] = useState<Set<string>>(new Set());

  // 初始化加载
  useEffect(() => {
    const preferences = getUserPreferences();
    setViewedPapers(new Set(preferences.viewedPapers));
  }, []);

  // 添加到浏览历史
  const addToViewed = useCallback((paperId: string) => {
    setViewedPapers((prev) => {
      const newSet = new Set(prev);
      newSet.add(paperId);

      // 保存到本地存储
      saveUserPreferences({
        viewedPapers: Array.from(newSet),
      });

      return newSet;
    });
  }, []);

  // 检查是否已浏览
  const isViewed = useCallback(
    (paperId: string) => {
      return viewedPapers.has(paperId);
    },
    [viewedPapers]
  );

  return {
    viewedPapers,
    addToViewed,
    isViewed,
  };
}
