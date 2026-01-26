/**
 * 加载骨架屏组件
 * 在数据加载时展示占位内容
 */

import './LoadingCard.css';

export function LoadingCard() {
  return (
    <div className="loading-card">
      <div className="loading-card__cover loading-card__shimmer"></div>
      <div className="loading-card__content">
        <div className="loading-card__title loading-card__shimmer"></div>
        <div className="loading-card__meta loading-card__shimmer"></div>
        <div className="loading-card__summary loading-card__shimmer"></div>
        <div className="loading-card__summary loading-card__shimmer"></div>
        <div className="loading-card__summary loading-card__shimmer"></div>
        <div className="loading-card__actions">
          <div className="loading-card__action loading-card__shimmer"></div>
          <div className="loading-card__action loading-card__shimmer"></div>
          <div className="loading-card__action loading-card__shimmer"></div>
          <div className="loading-card__action loading-card__shimmer"></div>
        </div>
      </div>
    </div>
  );
}

export function LoadingList({ count = 3 }: { count?: number }) {
  return (
    <div className="loading-list">
      {Array.from({ length: count }).map((_, index) => (
        <LoadingCard key={index} />
      ))}
    </div>
  );
}
