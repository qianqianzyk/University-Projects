<template>
  <div class="profile-content-multi">
    <el-card class="profile-card-multi styled-card-multi">
      <div class="profile-header-multi">
        <span class="profile-logo-multi">个人信息</span>
      </div>
      <el-divider />
      <el-form
        v-loading="loading"
        :model="profile"
        label-width="100px"
        class="profile-form-multi"
      >
        <div class="form-grid-multi">
          <el-form-item label="工号">
            <el-input v-model="profile.teacher_id" disabled class="disabled-input-multi" />
          </el-form-item>
          <el-form-item label="姓名">
            <el-input v-model="profile.name" disabled class="disabled-input-multi" />
          </el-form-item>
          <el-form-item label="性别">
            <el-input :value="profile.gender === 'M' ? '男' : profile.gender === 'F' ? '女' : profile.gender" disabled class="disabled-input-multi" />
          </el-form-item>
          <el-form-item label="年龄">
            <el-input v-model="profile.age" disabled class="disabled-input-multi" />
          </el-form-item>
          <el-form-item label="职称">
            <el-input v-model="profile.title" disabled class="disabled-input-multi" />
          </el-form-item>
          <el-form-item label="电话">
            <el-input v-model="profile.phone" disabled class="disabled-input-multi" />
          </el-form-item>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Star } from '@element-plus/icons-vue'
import { teacherApi } from '../../api/teacher'
import type { Teacher } from '../../api/teacher'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)

const profile = reactive<Teacher>({
  id: 0,
  teacher_id: '',
  name: '',
  gender: '',
  age: 0,
  title: '',
  phone: '',
  is_admin: false
})

const loadProfile = async () => {
  loading.value = true
  try {
    const teacherId = authStore.user?.id
    if (!teacherId) return
    const response = await teacherApi.getTeacher(teacherId)
    if (response.code >= 200500 && response.code < 200600) {
      ElMessage.error(response.msg || '系统异常，请稍后重试!')
      loading.value = false
      return
    }
    if ((response.code === 200 || response.code === 0) && response.data) {
      Object.assign(profile, response.data)
    }
  } catch (error) {
    ElMessage.error('加载个人信息失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.profile-card-multi {
  background: #f8faff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  padding: 40px 36px 36px 36px;
  min-width: 340px;
  max-width: 700px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  border: 1.5px solid #e0e7ff;
}
.profile-header-multi {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.profile-logo-multi {
  font-size: 26px;
  font-weight: 700;
  color: #6c63ff;
  letter-spacing: 2px;
}
.profile-form-multi {
  width: 100%;
}
.form-grid-multi {
  display: grid;
  grid-template-columns: repeat(2, minmax(220px, 1fr));
  gap: 20px;
}
.disabled-input-multi {
  background: #f3f6fa;
  color: #b0b3b8;
  width: 100%;
  box-sizing: border-box;
}
@media (max-width: 600px) {
  .profile-card-multi {
    min-width: 90vw;
    padding: 24px 8px 16px 8px;
  }
  .form-grid-multi {
    grid-template-columns: 1fr;
    gap: 10px;
  }
}
</style> 