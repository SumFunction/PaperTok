/**
 * API 服务层
 * 封装所有后端 API 调用
 */

import axios, { type AxiosInstance, type AxiosError } from 'axios';
import type { Paper, ApiResponse, PapersResponse } from '../types';

// API 基础 URL
// 生产环境使用相对路径（通过 Nginx 代理），开发环境使用 localhost:8080
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

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
 */
apiClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证 token 等信息
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * 响应拦截器
 */
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // 统一错误处理
    if (error.response) {
      // 服务器返回错误状态码
      console.error('API Error:', error.response.status, error.response.data);
    } else if (error.request) {
      // 请求已发送但没有收到响应
      console.error('Network Error:', error.message);
    } else {
      // 请求配置出错
      console.error('Request Error:', error.message);
    }
    return Promise.reject(error);
  }
);

/**
 * PaperTok API 类
 */
export class PaperTokAPI {
  /**
   * 获取论文列表
   */
  static async getPapers(params: {
    category?: string;
    limit?: number;
    offset?: number;
  }): Promise<PapersResponse> {
    try {
      const response = await apiClient.get<ApiResponse<PapersResponse>>('/api/v1/papers', {
        params: {
          category: params.category || '',
          limit: params.limit || 20,
          offset: params.offset || 0,
        },
      });

      if (response.data.success && response.data.data) {
        // 后端返回的 data 已经是 { papers, total, page, pageSize } 结构
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
    }
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
        // 后端返回的 data 已经是 { papers, total, page, pageSize } 结构
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
 * 导出 API 实例供其他地方使用
 */
export default apiClient;
