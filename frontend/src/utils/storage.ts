/**
 * 本地存储工具函数
 */

import { STORAGE_KEYS } from '../types';
import type { UserPreference } from '../types';

/**
 * 获取本地存储的数据
 */
export function getStorageData<T>(key: string): T | null {
  try {
    const item = window.localStorage.getItem(key);
    if (item === null) {
      return null;
    }
    return JSON.parse(item);
  } catch (error) {
    console.error(`Error reading from localStorage key "${key}":`, error);
    return null;
  }
}

/**
 * 设置本地存储的数据
 */
export function setStorageData<T>(key: string, value: T): boolean {
  try {
    window.localStorage.setItem(key, JSON.stringify(value));
    return true;
  } catch (error) {
    console.error(`Error writing to localStorage key "${key}":`, error);
    return false;
  }
}

/**
 * 删除本地存储的数据
 */
export function removeStorageData(key: string): boolean {
  try {
    window.localStorage.removeItem(key);
    return true;
  } catch (error) {
    console.error(`Error removing localStorage key "${key}":`, error);
    return false;
  }
}

/**
 * 获取用户偏好设置
 */
export function getUserPreferences(): UserPreference {
  const liked = getStorageData<string[]>(STORAGE_KEYS.LIKED_PAPERS) || [];
  const favorites = getStorageData<string[]>(STORAGE_KEYS.FAVORITE_PAPERS) || [];
  const viewed = getStorageData<string[]>(STORAGE_KEYS.VIEWED_PAPERS) || [];
  const preferences = getStorageData<{ selectedCategory: string }>(STORAGE_KEYS.PREFERENCES);

  return {
    likedPapers: liked,
    favoritePapers: favorites,
    viewedPapers: viewed,
    selectedCategory: preferences?.selectedCategory || '',
  };
}

/**
 * 保存用户偏好设置
 */
export function saveUserPreferences(preferences: Partial<UserPreference>): void {
  if (preferences.likedPapers !== undefined) {
    setStorageData(STORAGE_KEYS.LIKED_PAPERS, preferences.likedPapers);
  }
  if (preferences.favoritePapers !== undefined) {
    setStorageData(STORAGE_KEYS.FAVORITE_PAPERS, preferences.favoritePapers);
  }
  if (preferences.viewedPapers !== undefined) {
    setStorageData(STORAGE_KEYS.VIEWED_PAPERS, preferences.viewedPapers);
  }
  if (preferences.selectedCategory !== undefined) {
    setStorageData(STORAGE_KEYS.PREFERENCES, { selectedCategory: preferences.selectedCategory });
  }
}
