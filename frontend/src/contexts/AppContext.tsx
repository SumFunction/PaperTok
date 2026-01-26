/**
 * 应用全局状态管理 Context
 */

import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import type { Paper, AppState } from '../types';
import { useLikedPapers, useFavoritePapers } from '../hooks/useLocalStorage';

interface AppContextType extends AppState {
  // Actions
  setPapers: (papers: Paper[]) => void;
  toggleLike: (id: string) => void;
  toggleFavorite: (id: string) => void;
  isLiked: (id: string) => boolean;
  isFavorited: (id: string) => boolean;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setSelectedCategory: (category: string) => void;
  setSearchQuery: (query: string) => void;
  setSearchResults: (results: Paper[]) => void;
  clearSearchResults: () => void;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

interface AppProviderProps {
  children: ReactNode;
}

export function AppProvider({ children }: AppProviderProps) {
  // 自定义 Hooks
  const { likedPapers, toggleLike: toggleLikePaper, isLiked } = useLikedPapers();
  const { favoritePapers, toggleFavorite: toggleFavoritePaper, isFavorited } = useFavoritePapers();

  // 应用状态
  const [papers, setPapers] = useState<Paper[]>([]);
  const [viewedPapers, setViewedPapers] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [searchResults, setSearchResults] = useState<Paper[]>([]);

  /**
   * 切换点赞状态
   */
  const toggleLike = useCallback(
    (id: string) => {
      toggleLikePaper(id);
    },
    [toggleLikePaper]
  );

  /**
   * 切换收藏状态
   */
  const toggleFavorite = useCallback(
    (id: string) => {
      toggleFavoritePaper(id);
    },
    [toggleFavoritePaper]
  );

  /**
   * 清空搜索结果
   */
  const clearSearchResults = useCallback(() => {
    setSearchResults([]);
    setSearchQuery('');
  }, []);

  // Context 值
  const value: AppContextType = {
    // State
    papers,
    likedPapers,
    favoritePapers,
    viewedPapers,
    selectedCategory,
    loading,
    error,
    searchQuery,
    searchResults,

    // Actions
    setPapers,
    toggleLike,
    toggleFavorite,
    isLiked,
    isFavorited,
    setLoading,
    setError,
    setSelectedCategory,
    setSearchQuery,
    setSearchResults,
    clearSearchResults,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}

/**
 * 使用 App Context 的 Hook
 */
export function useAppContext(): AppContextType {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useAppContext must be used within an AppProvider');
  }
  return context;
}
