/**
 * 论文数据结构
 */
export interface Paper {
  id: string;                    // arXiv ID (e.g., "2401.12345")
  title: string;                 // 论文标题
  authors: string[];             // 作者列表
  summary: string;               // 摘要
  published: string;             // 发表时间 (ISO 8601)
  updated: string;               // 更新时间 (ISO 8601)
  categories: string[];          // 分类标签
  primaryCategory: string;       // 主分类
  imageUrl?: string;             // 封面图 URL
  arxivUrl: string;              // arXiv 原文链接
  pdfUrl: string;                // PDF 链接
}

/**
 * API 响应结构
 */
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  total?: number;
  page?: number;
  pageSize?: number;
}

/**
 * 论文列表响应
 */
export interface PapersResponse {
  papers: Paper[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 分类标签信息
 */
export interface Category {
  id: string;
  name: string;
  description: string;
}

/**
 * 用户偏好设置
 */
export interface UserPreference {
  likedPapers: string[];      // 点赞的论文 ID 列表
  favoritePapers: string[];   // 收藏的论文 ID 列表
  viewedPapers: string[];     // 已浏览的论文 ID 列表
  selectedCategory: string;   // 当前选择的分类
}

/**
 * 本地存储键名
 */
export const STORAGE_KEYS = {
  LIKED_PAPERS: 'papertok_liked',
  FAVORITE_PAPERS: 'papertok_favorites',
  VIEWED_PAPERS: 'papertok_viewed',
  PREFERENCES: 'papertok_preferences',
  AUTH_TOKEN: 'papertok_token',
  AUTH_USER: 'papertok_user',
} as const;

/**
 * 常用分类列表
 */
export const CATEGORIES: Category[] = [
  { id: '', name: '全部', description: '所有分类' },
  { id: 'cs.AI', name: '人工智能', description: 'Artificial Intelligence' },
  { id: 'cs.CL', name: '计算语言学', description: 'Computation and Language' },
  { id: 'cs.LG', name: '机器学习', description: 'Machine Learning' },
  { id: 'cs.CV', name: '计算机视觉', description: 'Computer Vision' },
  { id: 'cs.RO', name: '机器人', description: 'Robotics' },
  { id: 'cs.NE', name: '神经网络', description: 'Neural and Evolutionary Computing' },
  { id: 'cs.CR', name: '密码学', description: 'Cryptography and Security' },
  { id: 'cs.DB', name: '数据库', description: 'Databases' },
  { id: 'cs.DS', name: '数据结构与算法', description: 'Data Structures and Algorithms' },
];

/**
 * 用户信息
 */
export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

/**
 * 认证响应
 */
export interface AuthResponse {
  user: User;
  token: string;
}

/**
 * 登录请求
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * 注册请求
 */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

/**
 * API 错误响应
 */
export interface ApiError {
  success: false;
  error: string;
  timestamp: number;
}

/**
 * 应用全局状态
 */
export interface AppState {
  papers: Paper[];
  likedPapers: Set<string>;
  favoritePapers: Set<string>;
  viewedPapers: string[];
  selectedCategory: string;
  loading: boolean;
  error: string | null;
  searchQuery: string;
  searchResults: Paper[];
}
