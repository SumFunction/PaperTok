/**
 * 论文卡片组件
 * 抖音风格：全屏沉浸 + 右侧操作栏
 */

import { useState, useCallback, useRef } from 'react';
import type { Paper } from '../types';
import { formatRelativeTime } from '../utils/date';
import { formatAuthors, getCategoryName } from '../utils/text';
import './PaperCard.css';

interface PaperCardProps {
  paper: Paper;
  isLiked: boolean;
  isFavorited: boolean;
  onLike: (id: string) => void;
  onFavorite: (id: string) => void;
  onShare: (paper: Paper) => void;
  onView?: (id: string) => void;
}

// 分类对应的渐变色
const categoryGradients: Record<string, string> = {
  'cs.AI': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'cs.CL': 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
  'cs.LG': 'linear-gradient(135deg, #fc466b 0%, #3f5efb 100%)',
  'cs.CV': 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'cs.RO': 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  'cs.NE': 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  'cs.CR': 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
  'cs.DC': 'linear-gradient(135deg, #d299c2 0%, #fef9d7 100%)',
  'cs.SE': 'linear-gradient(135deg, #89f7fe 0%, #66a6ff 100%)',
  'stat.ML': 'linear-gradient(135deg, #f6d365 0%, #fda085 100%)',
  'math.OC': 'linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)',
  'default': 'linear-gradient(135deg, #2c3e50 0%, #4a5568 100%)',
};

export function PaperCard({
  paper,
  isLiked,
  isFavorited,
  onLike,
  onFavorite,
  onShare,
  onView,
}: PaperCardProps) {
  const [summaryExpanded, setSummaryExpanded] = useState(false);
  const [imageError, setImageError] = useState(false);
  const [showDoubleTapHeart, setShowDoubleTapHeart] = useState(false);
  const lastTapRef = useRef<number>(0);
  const tapTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // 获取作者首字母作为头像
  const getAuthorInitial = () => {
    if (paper.authors.length > 0) {
      return paper.authors[0].charAt(0).toUpperCase();
    }
    return 'P';
  };

  // 获取背景渐变色
  const getBackgroundGradient = () => {
    return categoryGradients[paper.primaryCategory] || categoryGradients['default'];
  };

  // 处理双击点赞
  const handleDoubleTap = useCallback(() => {
    const now = Date.now();
    const DOUBLE_TAP_DELAY = 300;
    
    if (now - lastTapRef.current < DOUBLE_TAP_DELAY) {
      // 双击触发点赞
      if (!isLiked) {
        onLike(paper.id);
      }
      setShowDoubleTapHeart(true);
      setTimeout(() => setShowDoubleTapHeart(false), 800);
      
      if (tapTimeoutRef.current) {
        clearTimeout(tapTimeoutRef.current);
      }
    } else {
      // 单击标记已浏览
      tapTimeoutRef.current = setTimeout(() => {
        if (onView) {
          onView(paper.id);
        }
      }, DOUBLE_TAP_DELAY);
    }
    
    lastTapRef.current = now;
  }, [isLiked, onLike, onView, paper.id]);

  // 处理分享
  const handleShare = async () => {
    try {
      if (navigator.share) {
        await navigator.share({
          title: paper.title,
          text: paper.summary.slice(0, 100) + '...',
          url: paper.arxivUrl,
        });
      } else {
        await navigator.clipboard.writeText(`${paper.title}\n${paper.arxivUrl}`);
      }
      onShare(paper);
    } catch (err) {
      // 用户取消分享或复制失败
      console.log('Share cancelled or failed:', err);
    }
  };

  return (
    <article className="paper-card" onClick={handleDoubleTap}>
      {/* 背景层 */}
      <div className="paper-card__background">
        {paper.imageUrl && !imageError ? (
          <img
            src={paper.imageUrl}
            alt=""
            className="paper-card__bg-image"
            onError={() => setImageError(true)}
            loading="lazy"
          />
        ) : (
          <div 
            className="paper-card__bg-gradient"
            style={{ background: getBackgroundGradient() }}
          >
            <div className="paper-card__bg-pattern" />
            <span className="paper-card__bg-category">
              {paper.primaryCategory.split('.')[1] || 'AI'}
            </span>
          </div>
        )}
      </div>

      {/* 渐变遮罩 */}
      <div className="paper-card__overlay" />

      {/* 双击点赞动画 */}
      {showDoubleTapHeart && (
        <div className="paper-card__double-tap-heart">
          <svg viewBox="0 0 24 24">
            <path d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        </div>
      )}

      {/* 内容区域 */}
      <div className="paper-card__content">
        {/* 分类标签 */}
        <div className="paper-card__tags">
          <span className="paper-card__tag">
            {getCategoryName(paper.primaryCategory)}
          </span>
          {paper.categories.slice(1, 3).map((cat) => (
            <span key={cat} className="paper-card__tag paper-card__tag--secondary">
              {getCategoryName(cat)}
            </span>
          ))}
        </div>

        {/* 标题 */}
        <h2 className="paper-card__title">{paper.title}</h2>

        {/* 作者信息 */}
        <div className="paper-card__meta">
          <div className="paper-card__author-avatar">
            {getAuthorInitial()}
          </div>
          <span className="paper-card__authors">
            {formatAuthors(paper.authors)}
          </span>
          <span className="paper-card__time">
            {formatRelativeTime(paper.published)}
          </span>
        </div>

        {/* 摘要 */}
        <div className="paper-card__summary">
          <p
            className={`paper-card__summary-text ${
              summaryExpanded ? 'paper-card__summary-text--expanded' : ''
            }`}
          >
            {paper.summary}
          </p>
          <button
            className="paper-card__expand-btn"
            onClick={(e) => {
              e.stopPropagation();
              setSummaryExpanded((prev) => !prev);
            }}
          >
            {summaryExpanded ? '收起 ↑' : '展开全文 ↓'}
          </button>
        </div>

        {/* 查看原文链接 */}
        <a
          href={paper.arxivUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="paper-card__link"
          onClick={(e) => e.stopPropagation()}
        >
          <svg className="paper-card__link-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6M15 3h6v6M10 14L21 3" />
          </svg>
          查看 arXiv 原文
        </a>
      </div>

      {/* 右侧操作栏 */}
      <div className="paper-card__actions">
        {/* 点赞 */}
        <button
          className={`paper-card__action ${isLiked ? 'paper-card__action--liked' : ''}`}
          onClick={(e) => {
            e.stopPropagation();
            onLike(paper.id);
          }}
          aria-label={isLiked ? '取消点赞' : '点赞'}
        >
          <div className="paper-card__action-icon">
            <svg
              viewBox="0 0 24 24"
              fill={isLiked ? 'currentColor' : 'none'}
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </svg>
          </div>
          <span className="paper-card__action-count">喜欢</span>
        </button>

        {/* 收藏 */}
        <button
          className={`paper-card__action ${isFavorited ? 'paper-card__action--favorited' : ''}`}
          onClick={(e) => {
            e.stopPropagation();
            onFavorite(paper.id);
          }}
          aria-label={isFavorited ? '取消收藏' : '收藏'}
        >
          <div className="paper-card__action-icon">
            <svg
              viewBox="0 0 24 24"
              fill={isFavorited ? 'currentColor' : 'none'}
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
          </div>
          <span className="paper-card__action-count">收藏</span>
        </button>

        {/* 分享 */}
        <button
          className="paper-card__action paper-card__action--share"
          onClick={(e) => {
            e.stopPropagation();
            handleShare();
          }}
          aria-label="分享"
        >
          <div className="paper-card__action-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
            </svg>
          </div>
          <span className="paper-card__action-count">分享</span>
        </button>

        {/* 查看原文 */}
        <a
          href={paper.pdfUrl || paper.arxivUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="paper-card__action paper-card__action--link"
          onClick={(e) => e.stopPropagation()}
          aria-label="查看 PDF"
        >
          <div className="paper-card__action-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              <path d="M12 3v6h6" />
              <path d="M9 13h6M9 17h4" />
            </svg>
          </div>
          <span className="paper-card__action-count">PDF</span>
        </a>
      </div>
    </article>
  );
}
