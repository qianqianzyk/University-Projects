<template>
  <div class="teachings-page">
    <el-card>
      <template #header>
        <div class="page-header">
          <span>授课管理</span>
          <el-button type="primary" @click="resetForm(); showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加授课
          </el-button>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="paginatedTeachings"
        style="width: 100%"
      >
        <el-table-column prop="teacher_name" label="教师姓名" />
        <el-table-column prop="teacher_id" label="教师工号" />
        <el-table-column prop="course_name" label="课程名称" />
        <el-table-column prop="course_info" label="课程信息" />
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
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑授课对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="isEdit ? '编辑授课' : '添加授课'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="教师" prop="teacher_id">
          <el-select v-model="form.teacher_id" placeholder="请选择教师" style="width: 100%">
            <el-option
              v-for="teacher in teachers"
              :key="teacher.teacher_id"
              :label="`${teacher.name} (${teacher.teacher_id})`"
              :value="teacher.teacher_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="课程" prop="course_id">
          <el-select v-model="form.course_id" placeholder="请选择课程" style="width: 100%">
            <el-option
              v-for="course in courses"
              :key="course.id"
              :label="`${course.name}（${course.class_id ? classMap[course.class_id] : '未知班级'}）${course.school_year}年${course.semester === '1' ? '上' : '下'}学期`"
              :value="course.id"
            />
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
import type { Teaching, Teacher, Course } from '../../types'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

const loading = ref(false)
const teachings = ref<Teaching[]>([])
const teachers = ref<Teacher[]>([])
const courses = ref<Course[]>([])
const showAddDialog = ref(false)
const isEdit = ref(false)
const formRef = ref()
const classes = ref<any[]>([])

// 分页相关
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 计算当前页显示的授课数据，并关联教师和课程信息
const paginatedTeachings = computed(() => {
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  const currentPageData = teachings.value.slice(start, end)
  
  // 为每条授课记录添加教师和课程信息
  return currentPageData.map(teaching => {
    const teacher = teachers.value.find(t => t.teacher_id === String(teaching.teacher_id))
    const course = courses.value.find(c => c.id === teaching.course_id)
    
    return {
      ...teaching,
      teacher_name: teacher?.name || '未知教师',
      course_name: course?.name || '未知课程',
      course_info: course ? `${(course as any).school_year}年${(course as any).semester === '1' ? '上' : '下'}学期` : '未知学期'
    }
  })
})

const form = reactive({
  id: 0,
  teacher_id: '',
  course_id: undefined as number | undefined
})

const rules = {
  teacher_id: [{ required: true, message: '请选择教师', trigger: 'change' }],
  course_id: [{ required: true, message: '请选择课程', trigger: 'change' }]
}

// 加载授课列表
const loadTeachings = async () => {
  loading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeachings(adminId)
    if (response.code >= 200500 && response.code < 200600) {
      ElMessage.error(response.msg || '系统异常，请稍后重试!')
      loading.value = false
      return
    }
    if (response.code === 200 || response.code === 0) {
      teachings.value = response.data.list || []
      pagination.total = teachings.value.length
    }
  } catch (error) {
    ElMessage.error('加载授课列表失败')
  } finally {
    loading.value = false
  }
}

// 加载教师列表
const loadTeachers = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeachers(adminId)
    if (response.code === 200 || response.code === 0) {
      teachers.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('加载教师列表失败')
  }
}

// 加载课程列表
const loadCourses = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCourses(adminId)
    if (response.code === 200 || response.code === 0) {
      courses.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('加载课程列表失败')
  }
}

// 加载班级列表
const loadClasses = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getClasses(adminId)
    if (response.code === 200 || response.code === 0) {
      classes.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('加载班级列表失败')
  }
}

// 课程id->班级名映射
const classMap = computed(() => {
  const map: Record<number, string> = {}
  classes.value.forEach(cls => {
    map[cls.id] = cls.name
  })
  return map
})

// 编辑授课
const handleEdit = (row: any) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    teacher_id: row.teacher_id,
    course_id: row.course_id
  })
  showAddDialog.value = true
}

// 删除授课
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这个授课安排吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const adminId = authStore.user?.id || 1
    await adminApi.deleteTeaching({ admin_id: adminId, id: row.id })
    ElMessage.success('删除成功')
    loadTeachings()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 检查课程是否已选择
    if (!form.course_id) {
      ElMessage.error('请选择课程')
      return
    }
    
    const adminId = authStore.user?.id || 1
    const submitData = {
      ...form,
      admin_id: adminId,
      course_id: form.course_id as number
    }
    
    if (isEdit.value) {
      // 编辑授课
      await adminApi.updateTeaching(submitData)
      ElMessage.success('更新成功')
    } else {
      // 添加授课
      await adminApi.createTeaching(submitData)
      ElMessage.success('添加成功')
    }
    
    showAddDialog.value = false
    loadTeachings()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    teacher_id: '',
    course_id: undefined
  })
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
  loadTeachings()
  loadTeachers()
  loadCourses()
  loadClasses()
})
</script>

<style scoped>
.teachings-page {
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