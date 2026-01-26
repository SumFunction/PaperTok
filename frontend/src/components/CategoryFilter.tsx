/**
 * 分类筛选器组件
 * 允许用户快速切换不同的论文分类
 */

import { CATEGORIES } from '../types';
import './CategoryFilter.css';

interface CategoryFilterProps {
  selectedCategory: string;
  onCategoryChange: (category: string) => void;
}

export function CategoryFilter({ selectedCategory, onCategoryChange }: CategoryFilterProps) {
  return (
    <nav className="category-filter" role="navigation" aria-label="分类筛选">
      <div className="category-filter__container">
        <div className="category-filter__list">
          {CATEGORIES.map((category) => (
            <div key={category.id} className="category-filter__item">
              <button
                className={`category-filter__button ${
                  selectedCategory === category.id ? 'category-filter__button--active' : ''
                }`}
                onClick={() => onCategoryChange(category.id)}
                aria-pressed={selectedCategory === category.id}
                title={category.description}
              >
                {category.name}
              </button>
            </div>
          ))}
        </div>
      </div>
    </nav>
  );
}
