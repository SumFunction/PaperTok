/**
 * 论文数据管理 Hook
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import type { Paper } from '../types';
import { PaperTokAPI } from '../services/api';

interface UsePapersResult {
  papers: Paper[];
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  fetchPapers: (category?: string, append?: boolean) => Promise<void>;
  refreshPapers: () => Promise<void>;
}

const LIMIT = 20;

export function usePapers(initialCategory: string = ''): UsePapersResult {
  const [papers, setPapers] = useState<Paper[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState<boolean>(true);

  // 使用 ref 存储可变状态，避免 useCallback 依赖变化
  const stateRef = useRef({
    currentCategory: initialCategory,
    offset: 0,
    isLoading: false,
    initialized: false,
  });

  /**
   * 获取论文列表 - 稳定的引用，不会因状态变化而改变
   */
  const fetchPapers = useCallback(async (category?: string, append: boolean = false): Promise<void> => {
    const state = stateRef.current;
    
    // 防止并发请求
    if (state.isLoading) {
      return;
    }

    const targetCategory = category !== undefined ? category : state.currentCategory;
    const currentOffset = append ? state.offset : 0;

    state.isLoading = true;
    setLoading(true);
    setError(null);

    try {
      const response = await PaperTokAPI.getPapers({
        category: targetCategory,
        limit: LIMIT,
        offset: currentOffset,
      });

      if (append) {
        // 追加模式：去重后添加
        setPapers((prevPapers) => {
          const existingIds = new Set(prevPapers.map((p) => p.id));
          const newPapers = response.papers.filter((p) => !existingIds.has(p.id));
          console.log('[usePapers] 追加模式: 返回', response.papers.length, '篇, 去重后新增', newPapers.length, '篇');
          // 无论是否有新论文，offset 都应该增加（按请求的数量）
          // 这样下次请求会跳过这批论文
          state.offset = currentOffset + LIMIT;
          console.log('[usePapers] 更新 offset:', state.offset);
          if (newPapers.length === 0) {
            console.log('[usePapers] 警告: 没有新论文被添加（全部重复）');
          }
          return [...prevPapers, ...newPapers];
        });
      } else {
        // 替换模式：去重后替换
        const uniquePapers = response.papers.filter(
          (paper, index, self) => self.findIndex((p) => p.id === paper.id) === index
        );
        setPapers(uniquePapers);
        state.offset = uniquePapers.length;
        state.currentCategory = targetCategory;
      }

      // 判断是否还有更多：如果返回的论文数量大于 0 且接近 LIMIT，则认为还有更多
      // 这样可以容忍后端返回略少于 LIMIT 的情况
      const newHasMore = response.papers.length > 0 && response.papers.length >= LIMIT - 5;
      console.log('[usePapers] 设置 hasMore:', newHasMore, '返回论文数:', response.papers.length, 'LIMIT:', LIMIT);
      setHasMore(newHasMore);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取论文失败';
      setError(errorMessage);
      console.error('Failed to fetch papers:', err);
    } finally {
      state.isLoading = false;
      setLoading(false);
    }
  }, []); // 空依赖，函数引用永远不变

  /**
   * 刷新论文列表
   */
  const refreshPapers = useCallback(async (): Promise<void> => {
    stateRef.current.offset = 0;
    setHasMore(true);
    await fetchPapers(stateRef.current.currentCategory, false);
  }, [fetchPapers]);

  /**
   * 初始加载 - 只在组件挂载时执行一次
   */
  useEffect(() => {
    if (stateRef.current.initialized) return;
    stateRef.current.initialized = true;
    
    fetchPapers(initialCategory, false);
  }, [initialCategory, fetchPapers]);

  return {
    papers,
    loading,
    error,
    hasMore,
    fetchPapers,
    refreshPapers,
  };
}
