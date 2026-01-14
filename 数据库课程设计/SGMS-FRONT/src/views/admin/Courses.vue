<template>
  <div class="courses-page">
    <el-card>
      <template #header>
        <div class="page-header">
          <span>课程管理</span>
          <el-button type="primary" @click="resetForm(); showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加课程
          </el-button>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="paginatedCourses"
        style="width: 100%"
      >
        <el-table-column prop="id" label="课程ID" width="80" />
        <el-table-column prop="name" label="课程名称" />
        <el-table-column prop="school_year" label="学年" />
        <el-table-column prop="semester" label="学期">
          <template #default="{ row }">
            {{ row.semester === '1' ? '上' : '下' }}
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="学时" />
        <el-table-column prop="credit" label="学分" />
        <el-table-column prop="class_id" label="所属班级">
          <template #default="{ row }">
            {{ getClassName(row.class_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="exam_type" label="考核方式" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页组件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[15, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑课程对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="isEdit ? '编辑课程' : '添加课程'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="学年" prop="school_year">
          <el-input-number v-model="form.school_year" :min="2020" :max="2030" />
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-select v-model="form.semester" placeholder="请选择学期">
            <el-option label="上" value="1" />
            <el-option label="下" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="学时" prop="hours">
          <el-input-number v-model="form.hours" :min="16" :max="200" />
        </el-form-item>
        <el-form-item label="学分" prop="credit">
          <el-input-number v-model="form.credit" :min="0.5" :max="10" :precision="1" />
        </el-form-item>
        <el-form-item label="所属班级" prop="class_id">
          <el-select v-model="form.class_id" placeholder="请选择班级">
            <el-option
              v-for="classItem in classes"
              :key="classItem.id"
              :label="classItem.name"
              :value="classItem.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="考核方式" prop="exam_type">
          <el-select v-model="form.exam_type" placeholder="请选择考核方式">
            <el-option label="考试" value="考试" />
            <el-option label="考查" value="考查" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { adminApi } from '../../api/admin'
import type { Course, Teacher, ClassInfo } from '../../types'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

const loading = ref(false)
const courses = ref<Course[]>([])
const teachers = ref<Teacher[]>([])
const classes = ref<ClassInfo[]>([])
const showAddDialog = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 分页相关
const pagination = reactive({
  currentPage: 1,
  pageSize: 15,
  total: 0
})

// 计算当前页显示的课程数据
const paginatedCourses = computed(() => {
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return courses.value.slice(start, end)
})

const form = reactive({
  id: 0,
  name: '',
  school_year: 2024,
  semester: '',
  hours: 48,
  credit: 3,
  class_id: undefined as number | undefined,
  exam_type: ''
})

const rules = {
  name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
  school_year: [{ required: true, message: '请输入学年', trigger: 'blur' }],
  semester: [{ required: true, message: '请选择学期', trigger: 'change' }],
  hours: [{ required: true, message: '请输入学时', trigger: 'blur' }],
  credit: [{ required: true, message: '请输入学分', trigger: 'blur' }],
  class_id: [{ required: true, message: '请选择班级', trigger: 'change' }],
  exam_type: [{ required: true, message: '请选择考核方式', trigger: 'change' }]
}

const loadCourses = async () => {
  loading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCourses(adminId)
    if ((response.code === 200 || response.code === 0) && Array.isArray(response.data.list)) {
      courses.value = response.data.list
      pagination.total = response.data.total || courses.value.length
    } else {
      courses.value = []
    }
  } catch (error) {
    ElMessage.error('加载课程列表失败')
  } finally {
    loading.value = false
  }
}

const loadTeachers = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeachers(adminId)
    if (response.code === 200 || response.code === 0) {
      teachers.value = response.data.list || []
    }
  } catch (error) {
    console.error('加载教师列表失败:', error)
  }
}

const loadClasses = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getClasses(adminId)
    if (response.code === 200 || response.code === 0) {
      classes.value = response.data.list || []
    }
  } catch (error) {
    console.error('加载班级列表失败:', error)
  }
}

// 根据班级ID获取班级名称
const getClassName = (classId: number) => {
  const classItem = classes.value.find(c => c.id === classId)
  return classItem ? classItem.name : `班级${classId}`
}

const handleEdit = (row: Course) => {
  isEdit.value = true
  Object.assign(form, row)
  showAddDialog.value = true
}

const handleDelete = async (row: Course) => {
  try {
    await ElMessageBox.confirm('确定要删除这个课程吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const adminId = authStore.user?.id || 1
    await adminApi.deleteCourse({ admin_id: adminId, id: row.id })
    ElMessage.success('删除成功')
    loadCourses()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 检查班级是否已选择
    if (!form.class_id) {
      ElMessage.error('请选择班级')
      return
    }
    
    const adminId = authStore.user?.id || 1
    const submitData = {
      ...form,
      admin_id: adminId,
      class_id: form.class_id as number
    }
    
    if (isEdit.value) {
      // 编辑课程
      await adminApi.updateCourse(submitData)
      ElMessage.success('更新成功')
    } else {
      // 添加课程
      await adminApi.createCourse(submitData)
      ElMessage.success('添加成功')
    }
    
    showAddDialog.value = false
    loadCourses()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  isEdit.value = false;
  Object.assign(form, {
    id: 0,
    name: '',
    school_year: 2024,
    semester: '',
    hours: 48,
    credit: 3,
    class_id: undefined,
    exam_type: ''
  });
}

// 分页处理函数
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
}

const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
}

onMounted(() => {
  loadCourses()
  loadTeachers()
  loadClasses()
})
</script>

<style scoped>
.courses-page {
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #6c63ff;
  font-weight: 600;
  font-size: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

/* 卡片样式统一 */
:deep(.el-card) {
  background: #f8faff;
  border: 1.5px solid #e0e7ff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  box-shadow: 0 8px 25px rgba(120, 80, 200, 0.15);
}

:deep(.el-card__header) {
  background: linear-gradient(135deg, #f8faff 0%, #e0e7ff 100%);
  border-bottom: 1px solid #e0e7ff;
  border-radius: 20px 20px 0 0;
  color: #6c63ff;
  font-weight: 600;
}

/* 表格样式优化 */
:deep(.el-table) {
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}

:deep(.el-table th) {
  background: linear-gradient(135deg, #f8faff 0%, #e0e7ff 100%);
  color: #6c63ff;
  font-weight: 600;
}

:deep(.el-table td) {
  color: #666;
}

/* 按钮样式统一 */
:deep(.el-button) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(120, 80, 200, 0.15);
}

:deep(.el-button--primary) {
  background: linear-gradient(90deg, #6c63ff 0%, #42e695 100%);
  border: none;
  color: #fff;
  font-weight: 600;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}

:deep(.el-button--danger) {
  background: linear-gradient(90deg, #ff6b6b 0%, #ff4757 100%);
  border: none;
  color: #fff;
  font-weight: 600;
}

:deep(.el-button--danger:hover) {
  background: linear-gradient(90deg, #ff4757 0%, #ff6b6b 100%);
  color: #fff;
}

/* 输入框和选择器样式 */
:deep(.el-input) {
  border-radius: 8px;
}

:deep(.el-select) {
  border-radius: 8px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-input-number) {
  border-radius: 8px;
}

/* 分页器样式 */
:deep(.el-pagination) {
  margin-top: 20px;
}

:deep(.el-pagination .el-pager li) {
  border-radius: 6px;
  transition: all 0.3s ease;
}

:deep(.el-pagination .el-pager li:hover) {
  background: linear-gradient(135deg, #6c63ff 0%, #42a5f5 100%);
  color: white;
}

/* 对话框样式 */
:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #f8faff 0%, #e0e7ff 100%);
  color: #6c63ff;
  font-weight: 600;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  background: #f8faff;
  border-top: 1px solid #e0e7ff;
}

/* 表单样式 */
:deep(.el-form-item__label) {
  color: #6c63ff;
  font-weight: 600;
}

:deep(.el-form-item__content) {
  color: #666;
}
</style> 