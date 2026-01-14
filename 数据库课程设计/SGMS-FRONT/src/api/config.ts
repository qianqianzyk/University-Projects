import axios from 'axios';
import type { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';

// API基础配置 - 使用相对路径，通过Vite代理访问
export const BASE_URL = '';

// 备用API地址（如果代理不工作）
export const FALLBACK_BASE_URL = 'http://121.43.236.83:8888';

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加调试信息（移除环境检查）
    console.log('API请求:', config.method?.toUpperCase(), config.url, config.params);

    return config;
  },
  (error) => {
    console.error('请求拦截器错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log('API响应成功:', response.config.url, response.data);
    return response.data;
  },
  (error) => {
    console.error('API请求错误:', error.config?.url, error);

    // 如果是CORS错误，提供友好的错误信息
    if (error.code === 'ERR_NETWORK' || error.message.includes('CORS')) {
      console.error('CORS错误，请检查代理配置或API服务器设置');
      // 不要在这里处理401跳转，让具体的业务逻辑处理
    }

    // 对于401错误，不在这里自动跳转，让具体的业务逻辑决定如何处理
    // 这样可以避免在登录页面出现循环跳转的问题

    // 确保错误信息能够正确传递
    if (error.response?.status === 401) {
      // 只在非登录页面时清除token
      const currentPath = window.location.pathname;
      if (currentPath !== '/login' && currentPath !== '/') {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

export default api; 