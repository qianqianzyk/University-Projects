<template>
  <div class="teachers-page">
    <el-card>
      <template #header>
        <div class="page-header">
          <span>教师管理</span>
          <el-button type="primary" @click="resetForm(); showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加教师
          </el-button>
        </div>
      </template>
      
      <!-- 搜索筛选区域 -->
      <div class="search-container">
        <el-input
          v-model="searchName"
          placeholder="请输入教师姓名搜索"
          style="width: 300px"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      
      <el-table
        v-loading="loading"
        :data="paginatedTeachers"
        style="width: 100%"
      >
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="teacher_id" label="工号" />
        <el-table-column prop="gender" label="性别">
          <template #default="{ row }">
            {{ row.gender === 'M' ? '男' : '女' }}
          </template>
        </el-table-column>
        <el-table-column prop="age" label="年龄" />
        <el-table-column prop="title" label="职称" />
        <el-table-column prop="phone" label="联系电话" />
        <el-table-column prop="is_admin" label="管理员权限">
          <template #default="{ row }">
            <el-tag :type="row.is_admin ? 'danger' : 'info'">
              {{ row.is_admin ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
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

    <!-- 添加/编辑教师对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="isEdit ? '编辑教师' : '添加教师'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="工号" prop="teacher_id">
          <el-input v-model="form.teacher_id" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select v-model="form.gender" placeholder="请选择性别">
            <el-option label="男" value="M" />
            <el-option label="女" value="F" />
          </el-select>
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input-number v-model="form.age" :min="18" :max="100" />
        </el-form-item>
        <el-form-item label="职称" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="form.password" type="password" />
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
import { Plus, Search } from '@element-plus/icons-vue'
import { adminApi } from '../../api/admin'
import type { Teacher } from '../../types'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

const loading = ref(false)
const teachers = ref<Teacher[]>([])
const showAddDialog = ref(false)
const isEdit = ref(false)
const formRef = ref()
const searchName = ref('')

// 分页相关
const pagination = reactive({
  currentPage: 1,
  pageSize: 15,
  total: 0
})

// 计算当前页显示的教师数据，支持搜索筛选
const paginatedTeachers = computed(() => {
  let filteredTeachers = teachers.value
  
  // 根据姓名搜索筛选
  if (searchName.value.trim()) {
    filteredTeachers = teachers.value.filter(teacher => 
      teacher.name.toLowerCase().includes(searchName.value.toLowerCase())
    )
  }
  
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return filteredTeachers.slice(start, end)
})

const form = reactive({
  id: 0,
  teacher_id: '',
  name: '',
  gender: '',
  age: 30,
  title: '',
  phone: '',
  is_admin: false,
  password: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  teacher_id: [{ required: true, message: '请输入工号', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [{ required: true, message: '请输入年龄', trigger: 'blur' }],
  title: [{ required: true, message: '请输入职称', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }]
}

const loadTeachers = async () => {
  loading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeachers(adminId)
    if ((response.code === 200 || response.code === 0) && Array.isArray(response.data.list)) {
      teachers.value = response.data.list
      pagination.total = teachers.value.length
    } else {
      teachers.value = []
    }
  } catch (error) {
    ElMessage.error('加载教师列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: Teacher) => {
  isEdit.value = true
  Object.assign(form, row)
  showAddDialog.value = true
}

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields(); // 先重置校验和表单
  }
  isEdit.value = false;
  Object.assign(form, {
    id: 0,
    teacher_id: '',
    name: '',
    gender: '',
    age: 30,
    title: '',
    phone: '',
    is_admin: false,
    password: ''
  });
}

const handleDelete = async (row: Teacher) => {
  try {
    await ElMessageBox.confirm('确定要删除这个教师吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const adminId = authStore.user?.id || 1
    await adminApi.deleteTeacher({ admin_id: adminId, id: row.id })
    ElMessage.success('删除成功')
    loadTeachers()
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
    
    const adminId = authStore.user?.id || 1
    
    if (isEdit.value) {
      // 编辑教师
      await adminApi.updateTeacher({ ...form, admin_id: adminId })
      ElMessage.success('更新成功')
    } else {
      // 添加教师
      await adminApi.createTeacher({ ...form, admin_id: adminId })
      ElMessage.success('添加成功')
    }
    
    showAddDialog.value = false
    loadTeachers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 分页处理函数
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
}

const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
}

// 处理搜索
const handleSearch = () => {
  pagination.currentPage = 1 // 重置到第一页
}

onMounted(() => {
  loadTeachers()
})
</script>

<style scoped>
.teachers-page {
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

.search-container {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
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

/* 标签样式 */
:deep(.el-tag) {
  border-radius: 6px;
  font-weight: 600;
}

:deep(.el-tag--danger) {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff4757 100%);
  border: none;
  color: #fff;
}

:deep(.el-tag--info) {
  background: linear-gradient(135deg, #6c757d 0%, #495057 100%);
  border: none;
  color: #fff;
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