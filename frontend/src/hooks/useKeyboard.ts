/**
 * 键盘快捷键 Hook
 */

import { useEffect, useCallback } from 'react';

interface KeyboardShortcut {
  key: string;
  action: () => void;
  description?: string;
}

export function useKeyboard(shortcuts: KeyboardShortcut[], enabled: boolean = true) {
  useEffect(() => {
    if (!enabled) return;

    const handleKeyDown = (event: KeyboardEvent) => {
      // 检查是否在输入框中
      const target = event.target as HTMLElement;
      const isInputFocused =
        target.tagName === 'INPUT' ||
        target.tagName === 'TEXTAREA' ||
        target.isContentEditable;

      // 如果在输入框中，不触发快捷键
      if (isInputFocused) return;

      // 查找匹配的快捷键
      const shortcut = shortcuts.find((s) => {
        // 支持大小写不敏感的字符键
        if (event.key.toLowerCase() === s.key.toLowerCase()) {
          return true;
        }
        // 支持特殊键名（如 ArrowUp, ArrowDown 等）
        if (event.key === s.key) {
          return true;
        }
        return false;
      });

      if (shortcut) {
        event.preventDefault();
        shortcut.action();
      }
    };

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [shortcuts, enabled]);
}

/**
 * 预定义的快捷键 Hook
 */
export function usePaperKeyboard(handlers: {
  onNext?: () => void;
  onPrevious?: () => void;
  onLike?: () => void;
  onFavorite?: () => void;
  onExpand?: () => void;
}) {
  const shortcuts: KeyboardShortcut[] = [];

  if (handlers.onNext) {
    shortcuts.push({ key: 'j', action: handlers.onNext, description: '下一篇' });
    shortcuts.push({ key: 'ArrowDown', action: handlers.onNext, description: '下一篇' });
  }

  if (handlers.onPrevious) {
    shortcuts.push({ key: 'k', action: handlers.onPrevious, description: '上一篇' });
    shortcuts.push({ key: 'ArrowUp', action: handlers.onPrevious, description: '上一篇' });
  }

  if (handlers.onLike) {
    shortcuts.push({ key: 'l', action: handlers.onLike, description: '点赞' });
  }

  if (handlers.onFavorite) {
    shortcuts.push({ key: 'f', action: handlers.onFavorite, description: '收藏' });
  }

  if (handlers.onExpand) {
    shortcuts.push({ key: ' ', action: handlers.onExpand, description: '展开摘要' });
  }

  useKeyboard(shortcuts);
}
