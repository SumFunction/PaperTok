/**
 * 论文数据管理 Hook
 */

import { useState, useEffect, useCallback } from 'react';
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

export function usePapers(initialCategory: string = ''): UsePapersResult {
  const [papers, setPapers] = useState<Paper[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [currentCategory, setCurrentCategory] = useState<string>(initialCategory);
  const [offset, setOffset] = useState<number>(0);
  const LIMIT = 20;

  /**
   * 获取论文列表
   */
  const fetchPapers = useCallback(async (category?: string, append: boolean = false) => {
    const targetCategory = category !== undefined ? category : currentCategory;

    if (loading) return;

    setLoading(true);
    setError(null);

    try {
      const response = await PaperTokAPI.getPapers({
        category: targetCategory,
        limit: LIMIT,
        offset: append ? offset : 0,
      });

      if (append) {
        setPapers((prevPapers) => [...prevPapers, ...response.papers]);
        setOffset((prevOffset) => prevOffset + response.papers.length);
      } else {
        setPapers(response.papers);
        setOffset(response.papers.length);
        setCurrentCategory(targetCategory);
      }

      setHasMore(response.papers.length === LIMIT);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取论文失败';
      setError(errorMessage);
      console.error('Failed to fetch papers:', err);
    } finally {
      setLoading(false);
    }
  }, [currentCategory, offset, loading]);

  /**
   * 刷新论文列表
   */
  const refreshPapers = useCallback(async () => {
    setOffset(0);
    setHasMore(true);
    await fetchPapers(currentCategory, false);
  }, [currentCategory, fetchPapers]);

  /**
   * 初始加载
   */
  useEffect(() => {
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
