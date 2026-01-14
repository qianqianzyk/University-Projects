import api from './config';
import type { LoginRequest, LoginResponse, ApiResponse } from '../types';

export const authApi = {
  // 用户登录
  login: (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    return api.post('/api/login', data);
  },
}; 