/**
 * API 服务层
 * 封装所有后端 API 调用
 */

import axios, { type AxiosInstance, type AxiosError } from 'axios';
import type {
  Paper,
  ApiResponse,
  PapersResponse,
  User,
  AuthResponse,
  ApiError,
  STORAGE_KEYS,
} from '../types';

// API 基础 URL
// 生产环境使用相对路径（通过 Nginx 代理），开发环境使用 localhost:8080
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

/**
 * Token 管理类
 */
class TokenManager {
  private static readonly TOKEN_KEY = 'papertok_token';
  private static readonly USER_KEY = 'papertok_user';

  /**
   * 获取存储的 token
   */
  static getToken(): string | null {
    try {
      return localStorage.getItem(this.TOKEN_KEY);
    } catch {
      return null;
    }
  }

  /**
   * 保存 token
   */
  static setToken(token: string): void {
    try {
      localStorage.setItem(this.TOKEN_KEY, token);
    } catch (error) {
      console.error('Failed to save token:', error);
    }
  }

  /**
   * 清除 token
   */
  static clearToken(): void {
    try {
      localStorage.removeItem(this.TOKEN_KEY);
      localStorage.removeItem(this.USER_KEY);
    } catch (error) {
      console.error('Failed to clear token:', error);
    }
  }

  /**
   * 保存用户信息
   */
  static setUser(user: User): void {
    try {
      localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    } catch (error) {
      console.error('Failed to save user:', error);
    }
  }

  /**
   * 获取用户信息
   */
  static getUser(): User | null {
    try {
      const userStr = localStorage.getItem(this.USER_KEY);
      if (userStr) {
        return JSON.parse(userStr);
      }
      return null;
    } catch {
      return null;
    }
  }

  /**
   * 清除所有认证信息
   */
  static clearAuth(): void {
    this.clearToken();
  }
}

/**
 * 创建 axios 实例
 */
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 秒超时
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * 请求拦截器
 * 自动添加认证 token
 */
apiClient.interceptors.request.use(
  (config) => {
    const token = TokenManager.getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * 响应拦截器
 * 统一错误处理
 */
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError<ApiError>) => {
    // 处理 401 未授权错误
    if (error.response?.status === 401) {
      TokenManager.clearAuth();
      // 不再强制跳转到登录页，让 AuthContext 统一管理认证状态
      // 组件会根据 isAuthenticated 状态决定是否显示 AuthModal
    }

    // 统一错误处理
    if (error.response) {
      console.error('API Error:', error.response.status, error.response.data);
    } else if (error.request) {
      console.error('Network Error:', error.message);
    } else {
      console.error('Request Error:', error.message);
    }
    return Promise.reject(error);
  }
);

/**
 * PaperTok API 类
 */
export class PaperTokAPI {
  // ==================== 论文相关 API ====================

  // 请求去重：存储正在进行的请求
  private static pendingRequests: Map<string, Promise<PapersResponse>> = new Map();
  // 最后一次请求的时间戳
  private static lastRequestTime: number = 0;
  // 最小请求间隔（毫秒）
  private static MIN_REQUEST_INTERVAL = 300;

  /**
   * 获取论文列表（带去重和节流）
   */
  static async getPapers(params: {
    category?: string;
    limit?: number;
    offset?: number;
  }): Promise<PapersResponse> {
    // 生成请求的唯一标识
    const requestKey = `papers_${params.category || ''}_${params.limit || 20}_${params.offset || 0}`;
    
    // 检查是否有相同的请求正在进行
    const pendingRequest = this.pendingRequests.get(requestKey);
    if (pendingRequest) {
      console.log('[API] 复用已有请求:', requestKey);
      return pendingRequest;
    }

    // 节流：检查距离上次请求的时间（仅对相同 offset 的请求进行节流）
    const now = Date.now();
    const timeSinceLastRequest = now - this.lastRequestTime;
    // 只对 offset=0 的初始请求进行节流，加载更多请求（offset>0）不节流
    if (params.offset === 0 && timeSinceLastRequest < this.MIN_REQUEST_INTERVAL) {
      console.log('[API] 请求过于频繁，跳过:', requestKey);
      // 返回空结果而不是报错，避免影响UI
      return { papers: [], total: 0, page: 1, pageSize: 20 };
    }

    this.lastRequestTime = now;
    console.log('[API] 发起新请求:', requestKey);

    // 创建请求 Promise
    const requestPromise = (async (): Promise<PapersResponse> => {
      try {
        const response = await apiClient.get<ApiResponse<PapersResponse>>('/api/v1/papers', {
          params: {
            category: params.category || '',
            limit: params.limit || 20,
            offset: params.offset || 0,
          },
        });

        if (response.data.success && response.data.data) {
          return {
            papers: response.data.data.papers || [],
            total: response.data.data.total || 0,
            page: response.data.data.page || 1,
            pageSize: response.data.data.pageSize || 20,
          };
        }

        throw new Error(response.data.error || '获取论文列表失败');
      } catch (error) {
        console.error('Failed to fetch papers:', error);
        throw error;
      } finally {
        // 请求完成后移除
        this.pendingRequests.delete(requestKey);
      }
    })();

    // 存储正在进行的请求
    this.pendingRequests.set(requestKey, requestPromise);

    return requestPromise;
  }

  /**
   * 搜索论文
   */
  static async searchPapers(params: {
    query: string;
    limit?: number;
    offset?: number;
  }): Promise<PapersResponse> {
    try {
      const response = await apiClient.get<ApiResponse<PapersResponse>>('/api/v1/papers/search', {
        params: {
          query: params.query,
          limit: params.limit || 20,
          offset: params.offset || 0,
        },
      });

      if (response.data.success && response.data.data) {
        return {
          papers: response.data.data.papers || [],
          total: response.data.data.total || 0,
          page: response.data.data.page || 1,
          pageSize: response.data.data.pageSize || 20,
        };
      }

      throw new Error(response.data.error || '搜索失败');
    } catch (error) {
      console.error('Failed to search papers:', error);
      throw error;
    }
  }

  // ==================== 认证相关 API ====================

  /**
   * 错误信息映射表
   * 根据后端错误码显示更友好的提示
   */
  private static readonly ERROR_MESSAGES: Record<string, string> = {
    'INVALID_CREDENTIALS': '邮箱或密码错误',
    'USER_NOT_FOUND': '用户不存在',
    'WEAK_PASSWORD': '密码强度不足，至少需要8个字符',
    'USER_EXISTS': '��户已存在',
    'EMAIL_EXISTS': '邮箱已被注册',
    'INVALID_EMAIL': '邮箱格式不正确',
    'INVALID_USERNAME': '用户名格式不正确',
    'UNAUTHORIZED': '请先登录',
    'FORBIDDEN': '无权限访问',
    'NOT_FOUND': '请求的资源不存在',
    'INTERNAL_ERROR': '服务器错误，请稍后重试',
    'NETWORK_ERROR': '网络连接失败，请检查网络',
  };

  /**
   * 获取友好的错误信息
   */
  private static getErrorMessage(error: unknown): string {
    const apiError = error as AxiosError<ApiError>;

    // 网络错误
    if (!apiError.response) {
      if (apiError.code === 'ERR_NETWORK') {
        return this.ERROR_MESSAGES.NETWORK_ERROR;
      }
      return '网络连接失败，请检查网络设置';
    }

    // 从响应中获取错误信息
    const serverError = apiError.response.data?.error;

    // 尝试匹配预定义的错误码
    if (serverError && this.ERROR_MESSAGES[serverError]) {
      return this.ERROR_MESSAGES[serverError];
    }

    // 返回服务器原始错误或默认错误
    return serverError || '操作失败，请稍后重试';
  }

  /**
   * 用户登录
   */
  static async login(email: string, password: string): Promise<AuthResponse> {
    try {
      const response = await apiClient.post<{ success: boolean; data: AuthResponse }>(
        '/api/v1/auth/login',
        { identifier: email, password }
      );

      if (response.data.success && response.data.data) {
        const { user, token } = response.data.data;
        // 保存 token 和用户信息
        TokenManager.setToken(token);
        TokenManager.setUser(user);
        return { user, token };
      }

      throw new Error('登录失败');
    } catch (error) {
      throw new Error(this.getErrorMessage(error));
    }
  }

  /**
   * 用户注册
   */
  static async register(username: string, email: string, password: string): Promise<AuthResponse> {
    try {
      const response = await apiClient.post<{ success: boolean; data: AuthResponse }>(
        '/api/v1/auth/register',
        { username, email, password }
      );

      if (response.data.success && response.data.data) {
        const { user, token } = response.data.data;
        // 保存 token 和用户信息
        TokenManager.setToken(token);
        TokenManager.setUser(user);
        return { user, token };
      }

      throw new Error('注册失败');
    } catch (error) {
      throw new Error(this.getErrorMessage(error));
    }
  }

  /**
   * 获取当前用户信息
   */
  static async getProfile(): Promise<User> {
    try {
      const response = await apiClient.get<{ success: boolean; data: User }>('/api/v1/auth/profile');

      if (response.data.success && response.data.data) {
        const user = response.data.data;
        // 更新本地用户信息
        TokenManager.setUser(user);
        return user;
      }

      throw new Error('获取用户信息失败');
    } catch (error) {
      throw new Error(this.getErrorMessage(error));
    }
  }

  /**
   * 用户登出
   */
  static async logout(): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/logout');
    } catch (error) {
      console.error('Logout error:', error);
      // 即使接口调用失败，也要清除本地认证信息
    } finally {
      TokenManager.clearAuth();
    }
  }

  // ==================== 健康检查 ====================

  /**
   * 健康检查
   */
  static async healthCheck(): Promise<{ status: string; version: string }> {
    try {
      const response = await apiClient.get('/health');
      return response.data;
    } catch (error) {
      console.error('Health check failed:', error);
      throw error;
    }
  }
}

/**
 * 导出 TokenManager 供认证模块使用
 */
export { TokenManager };

/**
 * 导出 API 实例供其他地方使用
 */
export default apiClient;
