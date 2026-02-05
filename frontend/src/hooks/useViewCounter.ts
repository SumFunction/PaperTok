/**
 * 浏览计数器 Hook
 * 跟踪用户浏览的论文数量，达到阈值时触发回调
 * 支持登录/未登录状态区分，数据持久化到 localStorage
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { useAuth } from '../contexts/AuthContext';

const STORAGE_KEY = 'papertok_view_count';
const VIEWED_PAPERS_KEY = 'papertok_viewed_papers';

interface UseViewCounterOptions {
  /** 触发阈值，默认 10 */
  threshold?: number;
  /** 达到阈值时的回调 */
  onThresholdReached?: () => void;
}

interface UseViewCounterReturn {
  /** 当前浏览计数 */
  viewCount: number;
  /** 增加计数（记录一次新的浏览） */
  increment: () => void;
  /** 重置计数 */
  reset: () => void;
  /** 已浏览的论文 ID 集合 */
  viewedPapers: Set<string>;
  /** 记录论文浏览 */
  recordView: (paperId: string) => void;
}

/**
 * 从 localStorage 获取已浏览的论文列表
 */
function getViewedPapersFromStorage(): Set<string> {
  try {
    const stored = localStorage.getItem(VIEWED_PAPERS_KEY);
    if (stored) {
      return new Set(JSON.parse(stored));
    }
  } catch (error) {
    console.warn('Failed to read viewed papers from storage:', error);
  }
  return new Set();
}

/**
 * 保存已浏览的论文列表到 localStorage
 */
function saveViewedPapersToStorage(papers: Set<string>): void {
  try {
    localStorage.setItem(VIEWED_PAPERS_KEY, JSON.stringify([...papers]));
  } catch (error) {
    console.warn('Failed to save viewed papers to storage:', error);
  }
}

/**
 * 浏览计数器 Hook
 */
export function useViewCounter(options: UseViewCounterOptions = {}): UseViewCounterReturn {
  const { threshold = 10, onThresholdReached } = options;
  const { isAuthenticated } = useAuth();

  // 使用 Set 避免重复计数同一论文
  const [viewedPapers, setViewedPapers] = useState<Set<string>>(() => getViewedPapersFromStorage());
  const viewCount = viewedPapers.size;

  // 记录是否已经触发过回调（避免重复触发）
  const hasTriggeredRef = useRef(false);

  /**
   * 记录论文浏览
   */
  const recordView = useCallback((paperId: string) => {
    setViewedPapers((prev) => {
      // 如果已经浏览过，不重复计数
      if (prev.has(paperId)) {
        console.log('[ViewCounter] 论文已浏览过，跳过:', paperId);
        return prev;
      }

      const newSet = new Set(prev);
      newSet.add(paperId);

      console.log('[ViewCounter] 记录新浏览:', paperId, '总计:', newSet.size);

      // 持久化到 localStorage
      saveViewedPapersToStorage(newSet);

      return newSet;
    });
  }, []);

  /**
   * 增加计数（通用方法）
   */
  const increment = useCallback(() => {
    setViewedPapers((prev) => {
      const newSet = new Set(prev);
      // 使用时间戳作为临时 ID，适用于不需要去重的场景
      const tempId = `temp_${Date.now()}`;
      newSet.add(tempId);
      saveViewedPapersToStorage(newSet);
      return newSet;
    });
  }, []);

  /**
   * 重置计数
   */
  const reset = useCallback(() => {
    setViewedPapers(new Set());
    saveViewedPapersToStorage(new Set());
    hasTriggeredRef.current = false;
  }, []);

  // 监听计数变化，达到阈值时触发回调
  useEffect(() => {
    console.log('[ViewCounter] 状态变化:', {
      viewCount,
      threshold,
      isAuthenticated,
      hasTriggered: hasTriggeredRef.current
    });

    // 只对未登录用户触发弹窗
    if (!isAuthenticated && viewCount >= threshold && !hasTriggeredRef.current) {
      console.log('[ViewCounter] 触发阈值回调！');
      hasTriggeredRef.current = true;
      onThresholdReached?.();
    }

    // 如果用户登录了，重置触发标记
    if (isAuthenticated) {
      hasTriggeredRef.current = false;
    }
  }, [viewCount, threshold, isAuthenticated, onThresholdReached]);

  // 用户登录状态变化时，重置计数
  useEffect(() => {
    if (isAuthenticated) {
      reset();
    }
  }, [isAuthenticated, reset]);

  return {
    viewCount,
    increment,
    reset,
    viewedPapers,
    recordView,
  };
}
