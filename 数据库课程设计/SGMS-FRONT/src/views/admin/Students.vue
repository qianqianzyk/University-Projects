<template>
  <div class="students-page">
    <el-card>
      <template #header>
        <div class="page-header">
          <span>学生管理</span>
          <el-button type="primary" @click="resetForm(); showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加学生
          </el-button>
        </div>
      </template>
      
      <!-- 搜索筛选区域 -->
      <div class="search-container">
        <el-input
          v-model="searchName"
          placeholder="请输入学生姓名搜索"
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
        :data="paginatedStudents"
        style="width: 100%"
      >
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="student_id" label="学号" />
        <el-table-column prop="gender" label="性别">
          <template #default="{ row }">
            {{ row.gender === 'M' ? '男' : '女' }}
          </template>
        </el-table-column>
        <el-table-column prop="age" label="年龄" />
        <el-table-column prop="province_id" label="省份">
          <template #default="{ row }">
            {{ getProvinceNameByCity(row.city_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="city_id" label="城市">
          <template #default="{ row }">
            {{ getCityName(row.city_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="department_id" label="部门">
          <template #default="{ row }">
            {{ getDepartmentNameByClass(row.class_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="class_id" label="班级">
          <template #default="{ row }">
            {{ getClassName(row.class_id) }}
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
    
    <!-- 添加/编辑学生对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="isEdit ? '编辑学生' : '添加学生'"
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
        <el-form-item label="学号" prop="student_id">
          <el-input v-model="form.student_id" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select v-model="form.gender" placeholder="请选择性别">
            <el-option label="男" value="M" />
            <el-option label="女" value="F" />
          </el-select>
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input v-model="form.age" type="number" />
        </el-form-item>
        <el-form-item label="省份" prop="province_id">
          <el-select v-model="form.province_id" placeholder="请选择省份" @change="onProvinceChange">
            <el-option 
              v-for="province in provinces" 
              :key="province.id" 
              :label="province.name" 
              :value="province.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="城市" prop="city_id" v-if="form.province_id">
          <el-select v-model="form.city_id" placeholder="请选择城市">
            <el-option 
              v-for="city in filteredCities" 
              :key="city.id" 
              :label="city.name" 
              :value="city.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="部门" prop="department_id">
          <el-select v-model="form.department_id" placeholder="请选择部门" @change="onDepartmentChange">
            <el-option 
              v-for="department in departments" 
              :key="department.id" 
              :label="department.name" 
              :value="department.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="班级" prop="class_id" v-if="form.department_id">
          <el-select v-model="form.class_id" placeholder="请选择班级">
            <el-option 
              v-for="classItem in filteredClasses" 
              :key="classItem.id" 
              :label="classItem.name" 
              :value="classItem.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="可选，留空则后端自动生成" />
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
import type { Student } from '../../types'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

const loading = ref(false)
const students = ref<Student[]>([])
const cities = ref<any[]>([])
const classes = ref<any[]>([])
const departments = ref<any[]>([])
const provinces = ref<any[]>([])
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

// 计算当前页显示的学生数据，支持搜索筛选
const paginatedStudents = computed(() => {
  let filteredStudents = students.value
  
  // 根据姓名搜索筛选
  if (searchName.value.trim()) {
    filteredStudents = students.value.filter(student => 
      student.name.toLowerCase().includes(searchName.value.toLowerCase())
    )
  }
  
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return filteredStudents.slice(start, end)
})

const form = reactive({
  id: 0,
  name: '',
  student_id: '',
  password: '',
  gender: '',
  age: 0,
  city_id: '',
  class_id: '',
  province_id: '',
  department_id: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  student_id: [{ required: true, message: '请输入学号', trigger: 'blur' }],
  password: [{ required: false, message: '请输入密码', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [
    { required: true, message: '请输入年龄', trigger: 'blur' },
    { validator: (rule: any, value: any, callback: any) => {
      const age = Number(value)
      if (isNaN(age) || age < 1 || age > 120) {
        callback(new Error('年龄必须在1-120之间'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ],
  city_id: [{ required: true, message: '请选择城市', trigger: 'change' }],
  class_id: [{ required: true, message: '请选择班级', trigger: 'change' }],
  province_id: [{ required: true, message: '请选择省份', trigger: 'change' }],
  department_id: [{ required: true, message: '请选择部门', trigger: 'change' }]
}

const loadStudents = async () => {
  loading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getStudents(adminId)
    if ((response.code === 200 || response.code === 0) && Array.isArray(response.data.list)) {
      students.value = response.data.list
      pagination.total = students.value.length
    } else {
      students.value = []
    }
  } catch (error) {
    ElMessage.error('加载学生列表失败')
  } finally {
    loading.value = false
  }
}

const loadCities = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCities(adminId)
    if (response.code === 200 || response.code === 0) {
      cities.value = response.data.cities || []
    }
  } catch (error) {
    console.error('加载城市列表失败:', error)
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

const loadDepartments = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getDepartments(adminId)
    if (response.code === 200 || response.code === 0) {
      departments.value = response.data.list || []
    }
  } catch (error) {
    console.error('加载部门列表失败:', error)
  }
}

const loadProvinces = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getProvinces(adminId)
    if (response.code === 200 || response.code === 0) {
      provinces.value = response.data.provinces || []
    }
  } catch (error) {
    console.error('加载省份列表失败:', error)
  }
}

const handleEdit = (row: Student) => {
  isEdit.value = true
  
  // 使用类型断言来处理可能缺失的字段
  const studentRow = row as any
  
  // 根据城市ID找到对应的省份ID
  const city = cities.value.find(c => c.id === studentRow.city_id)
  const provinceId = city ? city.province_id : ''
  
  // 根据班级ID找到对应的部门ID
  const classItem = classes.value.find(c => c.id === studentRow.class_id)
  const departmentId = classItem ? classItem.department_id : ''
  
  Object.assign(form, {
    ...studentRow,
    province_id: provinceId,
    department_id: departmentId
  })
  
  showAddDialog.value = true
}

const resetForm = () => {
  isEdit.value = false
  Object.assign(form, {
    id: 0,
    name: '',
    student_id: '',
    password: '',
    gender: '',
    age: 0,
    city_id: '',
    class_id: '',
    province_id: '',
    department_id: ''
  })
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const handleDelete = async (row: Student) => {
  try {
    await ElMessageBox.confirm('确定要删除这个学生吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const adminId = authStore.user?.id || 1
    await adminApi.deleteStudent({ admin_id: adminId, id: row.id })
    ElMessage.success('删除成功')
    loadStudents()
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
    
    // 准备提交数据，将字符串ID转换为数字
    const submitData = {
      ...form,
      admin_id: adminId,
      age: Number(form.age),
      city_id: Number(form.city_id),
      class_id: Number(form.class_id),
      province_id: Number(form.province_id),
      department_id: Number(form.department_id)
    }
    
    if (isEdit.value) {
      // 编辑学生
      await adminApi.updateStudent(submitData)
      ElMessage.success('更新成功')
    } else {
      // 添加学生
      await adminApi.createStudent(submitData)
      ElMessage.success('添加成功')
    }
    
    showAddDialog.value = false
    resetForm()
    loadStudents()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 计算属性：根据城市ID获取城市名称
const getCityName = (cityId: number) => {
  const city = cities.value.find(c => c.id === cityId)
  return city ? city.name : `城市${cityId}`
}

// 计算属性：根据班级ID获取班级名称
const getClassName = (classId: number) => {
  const classItem = classes.value.find(c => c.id === classId)
  return classItem ? classItem.name : `班级${classId}`
}

// 计算属性：根据部门ID获取部门名称
const getDepartmentName = (departmentId: number) => {
  const department = departments.value.find(d => d.id === departmentId)
  return department ? department.name : `部门${departmentId}`
}

// 计算属性：根据省份ID获取省份名称
const getProvinceName = (provinceId: number) => {
  const province = provinces.value.find(p => p.id === provinceId)
  return province ? province.name : `省份${provinceId}`
}

// 计算属性：根据城市ID获取省份名称
const getProvinceNameByCity = (cityId: number) => {
  const city = cities.value.find(c => c.id === cityId)
  if (city) {
    const province = provinces.value.find(p => p.id === city.province_id)
    return province ? province.name : `省份${city.province_id}`
  }
  return `省份-`
}

// 计算属性：根据班级ID获取部门名称
const getDepartmentNameByClass = (classId: number) => {
  const classItem = classes.value.find(c => c.id === classId)
  if (classItem) {
    const department = departments.value.find(d => d.id === classItem.department_id)
    return department ? department.name : `部门${classItem.department_id}`
  }
  return `部门-`
}

const filteredCities = computed(() => {
  if (!form.province_id) return [];
  return cities.value.filter(city => city.province_id === form.province_id);
})

const filteredClasses = computed(() => {
  if (!form.department_id) return [];
  return classes.value.filter(classItem => classItem.department_id === form.department_id);
})

const onProvinceChange = () => {
  form.city_id = '';
}

const onDepartmentChange = () => {
  form.class_id = '';
}

// 分页处理函数
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
}

const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
}

const handleSearch = () => {
  pagination.currentPage = 1 // 重置到第一页
}

onMounted(() => {
  loadStudents()
  loadCities()
  loadClasses()
  loadDepartments()
  loadProvinces()
})
</script>

<style scoped>
.students-page {
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