<template>
  <div class="profile-content-multi">
    <div class="profile-card-multi">
      <div class="profile-header-multi">
        <span class="profile-logo-multi">个人信息</span>
        <el-button type="primary" @click="saveProfile" :loading="loading" class="profile-save-btn-multi">
          <el-icon><Check /></el-icon>
          保存修改
        </el-button>
      </div>
      <el-form
        v-loading="loading"
        :model="profile"
        label-width="110px"
        class="profile-form-multi"
        :rules="rules"
        ref="formRef"
      >
        <div class="form-section-multi">
          <div class="form-grid-multi">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="profile.name" placeholder="请输入姓名" class="input-multi" />
            </el-form-item>
            <el-form-item label="学号">
              <el-input v-model="profile.student_id" disabled class="disabled-input-multi" />
            </el-form-item>
            <el-form-item label="性别" prop="gender">
              <el-select v-model="profile.gender" placeholder="请选择性别" class="input-multi">
                <el-option label="男" value="M" />
                <el-option label="女" value="F" />
              </el-select>
            </el-form-item>
            <el-form-item label="年龄" prop="age">
              <el-input-number v-model="profile.age" :min="16" :max="30" class="input-multi" />
            </el-form-item>
            <el-form-item label="班级">
              <el-input v-model="profile.class_name" disabled class="disabled-input-multi" />
            </el-form-item>
            <el-form-item label="生源地">
              <el-input v-model="profile.city_name" disabled class="disabled-input-multi" />
            </el-form-item>
            <el-form-item label="GPA">
              <el-input v-model="profile.gpa" disabled class="gpa-input-multi" />
            </el-form-item>
            <el-form-item label="总学分">
              <el-input v-model="profile.total_credits" disabled class="disabled-input-multi" />
            </el-form-item>
          </div>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, UserFilled, School, Check } from '@element-plus/icons-vue'
import { studentApi } from '../../api/student'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const profile = reactive({
  id: 0,
  student_id: '',
  name: '',
  gender: '',
  age: 0,
  city_id: 0,
  city_name: '',
  class_id: 0,
  class_name: '',
  gpa: 0,
  total_credits: 0,
})

const formRef = ref()
const rules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 10, message: '姓名长度在 2 到 10 个字符', trigger: 'blur' }
  ],
  gender: [
    { required: true, message: '请选择性别', trigger: 'change' }
  ],
  age: [
    { required: true, message: '请输入年龄', trigger: 'blur' },
    { type: 'number', min: 16, max: 30, message: '年龄必须在 16 到 30 之间', trigger: 'blur' }
  ]
}

const loadProfile = async () => {
  loading.value = true
  try {
    const studentId = authStore.user?.id
    if (!studentId) throw new Error('未获取到学生ID')
    const response = await studentApi.getStudent(studentId)
    if (response.code === 200 && response.data) {
      Object.assign(profile, response.data)
    }
  } catch (error) {
    ElMessage.error('加载个人信息失败')
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    loading.value = true
    try {
      const updateData = {
        id: profile.id,
        name: profile.name,
        gender: profile.gender,
        age: profile.age,
        city_id: profile.city_id
      }
      await studentApi.updateStudent(updateData)
      ElMessage.success('保存成功')
      loadProfile()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
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
  animation: fadeInUp 0.7s;
  border: 1.5px solid #e0e7ff;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
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
.profile-save-btn-multi {
  border-radius: 10px;
  font-weight: 600;
  background: linear-gradient(90deg, #6c63ff 0%, #42e695 100%);
  border: none;
  color: #fff;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(76, 201, 240, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 120px;
}
.profile-save-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
.profile-form-multi {
  width: 100%;
}
.form-section-multi {
  margin-bottom: 0;
  background: #fff;
  border-radius: 14px;
  padding: 24px;
  border: 1.5px solid #e0e7ff;
}
.form-grid-multi {
  display: grid;
  grid-template-columns: repeat(2, minmax(220px, 1fr));
  gap: 20px;
}
.input-multi,
.disabled-input-multi,
.gpa-input-multi {
  width: 100%;
  box-sizing: border-box;
}
.input-multi {
  --el-input-border-color: #b39ddb;
  --el-input-hover-border-color: #6c63ff;
  --el-input-focus-border-color: #42e695;
}
.disabled-input-multi {
  background: #f3f6fa;
  color: #b0b3b8;
}
.gpa-input-multi {
  font-weight: bold;
  color: #43a047;
  background: #e8f5e9;
  border-radius: 6px;
  border: none;
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