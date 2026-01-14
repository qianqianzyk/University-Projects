import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { User } from '../types';
import { authApi } from '../api/auth';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(null);

  // 初始化状态
  const initAuth = () => {
    const savedToken = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');

    if (savedToken && savedUser) {
      token.value = savedToken;
      user.value = JSON.parse(savedUser);
    }
  };

  // 登录
  const login = async (username: string, password: string, userType: number) => {
    try {
      const response = await authApi.login({ username, password, user_type: userType });
      console.log('登录响应完整数据:', response); // 调试信息
      console.log('响应code:', response.code, '响应msg:', response.msg); // 调试信息

      // 后端返回 code: 200 表示成功，其他值表示失败
      if (response.code === 200) {
        // 根据API文档，登录响应包含id和user_type
        const userData = {
          id: response.data.id,
          user_type: response.data.user_type,
          username: username
        };

        token.value = `token_${response.data.id}`; // 临时token
        user.value = userData;

        localStorage.setItem('token', token.value);
        localStorage.setItem('user', JSON.stringify(userData));

        return { success: true };
      } else {
        // 后端返回错误信息，code不为0表示失败
        const errorMessage = response.msg || response.msg || '登录失败，用户名或密码错误'
        console.error('登录失败 - 后端错误码:', response.code, '错误信息:', errorMessage)
        console.log('返回错误结果:', { success: false, message: errorMessage })
        return { success: false, message: errorMessage };
      }
    } catch (error: any) {
      console.error('登录错误:', error); // 调试信息

      // 处理不同类型的错误
      let errorMessage = '登录失败，请稍后重试'

      if (error.response) {
        // 服务器返回了错误状态码
        const status = error.response.status
        const data = error.response.data

        if (status === 401) {
          errorMessage = '用户名或密码错误'
        } else if (status === 403) {
          errorMessage = '账户已被禁用，请联系管理员'
        } else if (status === 404) {
          errorMessage = '登录接口不存在，请联系管理员'
        } else if (status === 500) {
          errorMessage = '服务器内部错误，请稍后重试'
        } else if (data && data.msg) {
          errorMessage = data.msg
        } else {
          errorMessage = `服务器错误 (${status})`
        }
      } else if (error.request) {
        // 请求已发出但没有收到响应
        if (error.code === 'ERR_NETWORK') {
          errorMessage = '网络连接失败，请检查网络设置'
        } else if (error.message.includes('timeout')) {
          errorMessage = '请求超时，请稍后重试'
        } else {
          errorMessage = '无法连接到服务器，请检查网络连接'
        }
      } else if (error.message) {
        // 其他错误
        errorMessage = error.message
      }

      console.log('返回异常错误结果:', { success: false, message: errorMessage })
      return { success: false, message: errorMessage };
    }
  };

  // 登出
  const logout = () => {
    user.value = null;
    token.value = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  // 检查是否已登录
  const isLoggedIn = () => {
    return !!token.value && !!user.value;
  };

  // 检查用户角色
  const hasRole = (role: string) => {
    if (!user.value) {
      console.log('hasRole: 用户未登录')
      return false;
    }

    // 根据user_type判断角色
    const userType = user.value.user_type;
    console.log('hasRole: 检查角色', role, '用户类型:', userType)

    switch (role) {
      case 'student':
        const isStudent = userType === 1;
        console.log('hasRole: 学生角色检查结果:', isStudent)
        return isStudent;
      case 'teacher':
        const isTeacher = userType === 2;
        console.log('hasRole: 教师角色检查结果:', isTeacher)
        return isTeacher;
      case 'admin':
        const isAdmin = userType === 3;
        console.log('hasRole: 管理员角色检查结果:', isAdmin)
        return isAdmin;
      default:
        console.log('hasRole: 未知角色:', role)
        return false;
    }
  };

  // 获取用户类型
  const getUserType = () => {
    return user.value?.user_type;
  };

  // 获取用户ID
  const getUserId = () => {
    return user.value?.id;
  };

  return {
    user,
    token,
    initAuth,
    login,
    logout,
    isLoggedIn,
    hasRole,
    getUserType,
    getUserId,
  };
}); 