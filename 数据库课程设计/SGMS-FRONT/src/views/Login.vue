<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <span class="login-logo">高校成绩管理系统</span>
        <span class="login-title">用户登录</span>
      </div>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="rules"
        class="login-form-new"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="user_type">
          <el-radio-group v-model="loginForm.user_type" class="user-type-radio">
            <el-radio-button :label="1">学生</el-radio-button>
            <el-radio-button :label="2">教师</el-radio-button>
            <el-radio-button :label="3">管理员</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名（6位及以上，仅包含数字）"
            size="large"
            clearable
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码（6位及以上，含字母和数字）"
            size="large"
            clearable
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="success"
            size="large"
            :loading="loading"
            class="login-btn-new"
            @click="handleLogin"
            auto-insert-space
          >
            {{ loading ? '登录中...' : '立即登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  user_type: 1
})

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 6, message: '用户名长度不能少于6个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d).+$/, message: '密码必须包含字母和数字', trigger: 'blur' }
  ],
  user_type: [
    { required: true, message: '请选择用户类型', trigger: 'change' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  try {
    await loginFormRef.value.validate()
    loading.value = true
    const result = await authStore.login(loginForm.username, loginForm.password, loginForm.user_type)
    if (result.success) {
      ElMessage.success('登录成功')
      let targetRoute = ''
      if (authStore.hasRole('admin')) {
        targetRoute = '/admin/dashboard'
      } else if (authStore.hasRole('student')) {
        targetRoute = '/student/dashboard'
      } else if (authStore.hasRole('teacher')) {
        targetRoute = '/teacher/dashboard'
      }
      if (targetRoute) {
        try {
          await router.push(targetRoute)
        } catch (routerError) {
          ElMessage.error('页面跳转失败，请刷新页面重试')
        }
      } else {
        ElMessage.error('用户角色识别失败')
      }
    } else {
      const errorMessage = result.message || '登录失败，请检查用户名和密码'
      ElMessage.error(errorMessage)
    }
  } catch (error: any) {
    let errorMessage = '登录失败，请检查输入信息'
    if (error.message) {
      if (error.message.includes('Network Error') || error.message.includes('ERR_NETWORK')) {
        errorMessage = '网络连接失败，请检查网络设置或联系管理员'
      } else if (error.message.includes('timeout')) {
        errorMessage = '请求超时，请稍后重试'
      } else if (error.message.includes('CORS')) {
        errorMessage = '服务器连接异常，请联系管理员'
      } else {
        errorMessage = error.message
      }
    }
    ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  position: fixed;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #e0eafc 0%, #cfdef3 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}
.login-card {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(80, 120, 180, 0.13);
  padding: 48px 36px 36px 36px;
  min-width: 340px;
  max-width: 400px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  animation: fadeInUp 0.7s;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
}
.login-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 32px;
}
.login-logo {
  font-size: 38px;
  font-weight: 900;
  color: #4f8cff;
  letter-spacing: 2px;
  margin-bottom: 6px;
}
.login-title {
  font-size: 22px;
  font-weight: 600;
  color: #333;
  letter-spacing: 1px;
}
.login-form-new {
  width: 100%;
}
.user-type-radio {
  width: 100%;
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}
.login-btn-new {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(90deg, #4f8cff 0%, #38e7ff 100%);
  border: none;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(80, 120, 180, 0.08);
}
.login-btn-new:hover {
  box-shadow: 0 6px 18px rgba(80, 120, 180, 0.18);
  transform: translateY(-2px) scale(1.03);
}
@media (max-width: 480px) {
  .login-card {
    min-width: 90vw;
    padding: 24px 8px 16px 8px;
  }
  .login-header {
    margin-bottom: 18px;
  }
  .login-title {
    font-size: 18px;
  }
}
</style>
<style>
body, html {
  margin: 0;
  padding: 0;
}
</style> 