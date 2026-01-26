/**
 * 搜索功能 Hook
 */

import { useState, useCallback } from 'react';
import type { Paper } from '../types';
import { PaperTokAPI } from '../services/api';

interface UseSearchResult {
  searchResults: Paper[];
  searching: boolean;
  searchError: string | null;
  searchPapers: (query: string) => Promise<void>;
  clearSearch: () => void;
}

export function useSearch(): UseSearchResult {
  const [searchResults, setSearchResults] = useState<Paper[]>([]);
  const [searching, setSearching] = useState<boolean>(false);
  const [searchError, setSearchError] = useState<string | null>(null);

  /**
   * 搜索论文
   */
  const searchPapers = useCallback(async (query: string) => {
    if (!query.trim()) {
      setSearchResults([]);
      return;
    }

    setSearching(true);
    setSearchError(null);

    try {
      const response = await PaperTokAPI.searchPapers({
        query,
        limit: 50,
      });

      setSearchResults(response.papers);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '搜索失败';
      setSearchError(errorMessage);
      console.error('Failed to search papers:', err);
    } finally {
      setSearching(false);
    }
  }, []);

  /**
   * 清空搜索结果
   */
  const clearSearch = useCallback(() => {
    setSearchResults([]);
    setSearchError(null);
  }, []);

  return {
    searchResults,
    searching,
    searchError,
    searchPapers,
    clearSearch,
  };
}
