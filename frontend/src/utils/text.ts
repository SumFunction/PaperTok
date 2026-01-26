/**
 * 文本处理工具函数
 */

/**
 * 截断文本，添加省略号
 */
export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) {
    return text;
  }
  return text.slice(0, maxLength) + '...';
}

/**
 * 格式化作者列表（最多显示前3个）
 */
export function formatAuthors(authors: string[]): string {
  if (authors.length === 0) {
    return '未知作者';
  }

  if (authors.length <= 3) {
    return authors.join(', ');
  }

  return authors.slice(0, 3).join(', ') + ' et al.';
}

/**
 * 获取分类显示名称
 */
export function getCategoryName(category: string): string {
  const categoryMap: Record<string, string> = {
    'cs.AI': '人工智能',
    'cs.CL': '计算语言学',
    'cs.LG': '机器学习',
    'cs.CV': '计算机视觉',
    'cs.RO': '机器人',
    'cs.NE': '神经网络',
    'cs.CR': '密码学',
    'cs.DB': '数据库',
    'cs.DS': '数据结构与算法',
  };

  return categoryMap[category] || category;
}

/**
 * 清理 HTML 标签
 */
export function stripHtmlTags(html: string): string {
  const tmp = document.createElement('div');
  tmp.innerHTML = html;
  return tmp.textContent || tmp.innerText || '';
}

/**
 * 计算阅读时间（基于字数）
 */
export function calculateReadingTime(text: string): string {
  const wordsPerMinute = 200; // 平均阅读速度
  const wordCount = text.split(/\s+/).length;
  const minutes = Math.ceil(wordCount / wordsPerMinute);

  if (minutes < 1) {
    return '1 分钟';
  }

  return `${minutes} 分钟`;
}
