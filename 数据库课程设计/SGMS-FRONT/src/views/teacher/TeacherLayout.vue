<template>
  <div class="layout-root-blue">
    <aside class="sidebar-blue">
      <div class="logo-blue">
        <h2>高校成绩管理系统</h2>
        <p>教师端</p>
      </div>
      <el-menu
        :default-active="route.path"
        class="sidebar-menu-blue"
        router
        background-color="transparent"
        text-color="#fff"
        active-text-color="#1976d2"
      >
        <el-menu-item index="/teacher/dashboard">
          <el-icon><DataBoard /></el-icon>
          <span>概览</span>
        </el-menu-item>
        <el-menu-item index="/teacher/courses">
          <el-icon><Document /></el-icon>
          <span>我的课程</span>
        </el-menu-item>
        <el-menu-item v-if="authStore.user?.id === 3" index="/teacher/students">
          <el-icon><User /></el-icon>
          <span>学生管理</span>
        </el-menu-item>
        <el-menu-item index="/teacher/profile">
          <el-icon><Avatar /></el-icon>
          <span>个人信息</span>
        </el-menu-item>
      </el-menu>
    </aside>
    <div class="main-area-blue">
      <header class="header-blue">
        <div class="header-left-blue">
          <el-breadcrumb separator="/" class="breadcrumb-blue styled-breadcrumb-multi">
            <el-breadcrumb-item>
              <el-icon style="vertical-align: middle; margin-right: 2px;"><DataBoard /></el-icon>
              首页
            </el-breadcrumb-item>
            <el-breadcrumb-item>
              <el-icon v-if="getCurrentPageTitle() === '概览'" style="vertical-align: middle; margin-right: 2px;"><DataBoard /></el-icon>
              <el-icon v-else-if="getCurrentPageTitle() === '我的课程'" style="vertical-align: middle; margin-right: 2px;"><Document /></el-icon>
              <el-icon v-else-if="getCurrentPageTitle() === '学生管理'" style="vertical-align: middle; margin-right: 2px;"><User /></el-icon>
              <el-icon v-else-if="getCurrentPageTitle() === '个人信息'" style="vertical-align: middle; margin-right: 2px;"><Avatar /></el-icon>
              <span class="breadcrumb-current-multi">{{ getCurrentPageTitle() }}</span>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right-blue">
          <el-dropdown @command="handleCommand">
            <span class="user-info-blue">
              <el-icon class="user-avatar-blue"><User /></el-icon>
              <span class="username-blue">{{ username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      <main class="main-content-blue">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  DataBoard,
  Reading,
  User,
  Document,
  Avatar,
  ArrowDown
} from '@element-plus/icons-vue'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const username = computed(() => authStore.user?.username || '')

const getCurrentPageTitle = () => {
  const routeMap: Record<string, string> = {
    '/teacher/dashboard': '概览',
    '/teacher/courses': '我的课程',
    '/teacher/students': '学生管理',
    '/teacher/profile': '个人信息'
  }
  return routeMap[route.path] || '首页'
}

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      authStore.logout()
      router.push('/login')
    } catch {
      // 用户取消
    }
  } else if (command === 'profile') {
    router.push('/teacher/profile')
  }
}
</script>

<style scoped>
.layout-root-blue {
  position: fixed;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  overflow: hidden;
  background: linear-gradient(135deg, #e3f0ff 0%, #b3e0ff 100%);
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
.sidebar-blue {
  width: 220px !important;
  min-width: 220px !important;
  max-width: 220px !important;
  flex-shrink: 0 !important;
  height: 100vh;
  background: linear-gradient(180deg, #2196f3 0%, #90caf9 100%);
  color: #fff;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
  box-shadow: 2px 0 12px 0 rgba(25,118,210,0.07);
}
.logo-blue {
  padding: 28px 20px 16px 20px;
  text-align: center;
  border-bottom: 1px solid #e3f0ff;
}
.logo-blue h2 {
  margin: 0;
  font-size: 20px;
  color: #fff;
  font-weight: 900;
  letter-spacing: 2px;
}
.logo-blue p {
  margin: 5px 0 0 0;
  font-size: 13px;
  color: #e3f0ff;
  letter-spacing: 1px;
}
.sidebar-menu-blue {
  border: none;
  background: transparent;
}
.sidebar-menu-blue .el-menu-item {
  border-radius: 8px;
  margin: 8px 12px;
  font-size: 16px;
  transition: background 0.2s, color 0.2s;
}
.sidebar-menu-blue .el-menu-item.is-active {
  background: #fff !important;
  color: #1976d2 !important;
  font-weight: bold;
}
.sidebar-menu-blue .el-menu-item:hover {
  background: rgba(255,255,255,0.18) !important;
  color: #fff !important;
}
.main-area-blue {
  flex: 1;
  height: 100vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
.header-blue {
  height: 64px;
  background: #e3f0ff;
  border-bottom: 1px solid #b3e0ff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  box-sizing: border-box;
  margin: 0;
  box-shadow: 0 2px 8px rgba(25,118,210,0.04);
}
.header-left-blue {
  flex: none;
  display: flex;
  align-items: center;
  height: 100%;
}
.breadcrumb-blue.styled-breadcrumb-multi {
  background: linear-gradient(90deg, #e3f0ff 0%, #b3e0ff 100%);
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(120, 80, 200, 0.10);
  padding: 8px 24px;
  display: inline-flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  margin-left: 24px;
  margin-top: 8px;
  margin-bottom: 8px;
}
.styled-breadcrumb-multi .el-breadcrumb__separator {
  color: #42e695;
  font-weight: bold;
  font-size: 18px;
  margin: 0 8px;
}
.styled-breadcrumb-multi .el-breadcrumb__item {
  color: #1976d2;
  transition: color 0.2s;
}
.styled-breadcrumb-multi .el-breadcrumb__item:last-child .breadcrumb-current-multi {
  color: #6c63ff;
  font-weight: 900;
  text-shadow: 0 2px 8px rgba(120, 80, 200, 0.10);
}
.styled-breadcrumb-multi .el-breadcrumb__item .el-icon {
  font-size: 18px;
  margin-right: 2px;
  vertical-align: middle;
}
.header-right-blue {
  display: flex;
  align-items: center;
}
.user-info-blue {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px 14px;
  border-radius: 6px;
  transition: background-color 0.3s;
  font-size: 16px;
}
.user-info-blue:hover {
  background-color: #bbdefb;
}
.username-blue {
  margin: 0 8px;
  color: #1976d2;
  font-weight: 600;
}
.user-avatar-blue {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1976d2 0%, #90caf9 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-right: 2px;
  box-shadow: 0 2px 8px rgba(25,118,210,0.08);
}
.main-content-blue {
  flex: 1;
  overflow: auto;
  background: transparent;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  margin: 0;
  padding: 32px 32px 0 32px;
}
@media (max-width: 900px) {
  .main-content-blue {
    padding: 16px 4px 0 4px;
  }
  .header-blue {
    padding: 0 8px;
  }
}
</style> 