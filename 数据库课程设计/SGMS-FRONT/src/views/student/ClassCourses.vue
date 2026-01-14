<template>
  <div class="class-courses-content-multi">
    <el-card class="class-courses-card-multi">
      <template #header>
        <div class="page-header-multi">
          <span>班级课程</span>
          <div class="header-actions-multi">
            <el-select v-model="selectedYear" placeholder="选择学年" size="small" class="select-multi" @change="loadClassCourses">
              <el-option v-for="year in years" :key="year" :label="year" :value="year" />
            </el-select>
            <el-select v-model="selectedSemester" :placeholder="selectedSemester === '' ? '全部' : '选择学期'" size="small" class="select-multi" @change="loadClassCourses">
              <el-option label="全部" value="" />
              <el-option label="上学期" value="1" />
              <el-option label="下学期" value="2" />
            </el-select>
          </div>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="filteredCourses"
        class="class-courses-table-multi"
        style="width: 100%"
      >
        <el-table-column prop="course_name" label="课程名称" />
        <el-table-column prop="credits" label="学分" />
        <el-table-column prop="school_year" label="学年" />
        <el-table-column prop="semester" label="学期">
          <template #default="{ row }">
            <span>{{ row.semester === 1 ? '上学期' : '下学期' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="课时" />
        <el-table-column prop="exam_type" label="考试类型" />
        <el-table-column label="授课教师">
          <template #default="{ row }">
            <span>{{ row.teachers && row.teachers.length ? row.teachers.map((t: any) => t.teacher_name).join('，') : '-' }}</span>
          </template>
        </el-table-column>
        <template #empty>
          <div style="text-align:center;color:#999;">暂无课程数据</div>
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { studentApi } from '../../api/student'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const courses = ref<any[]>([])
const years = ref<number[]>([])
const selectedYear = ref<number>(new Date().getFullYear())
const selectedSemester = ref<string>('')

const getYearOptions = () => {
  const thisYear = new Date().getFullYear()
  const minYear = 2020
  const arr = []
  for (let y = thisYear; y >= minYear; y--) arr.push(y)
  return arr
}

const filteredCourses = computed(() => {
  let list = courses.value
  if (selectedYear.value) {
    list = list.filter((s: any) => s.school_year === selectedYear.value)
  }
  if (selectedSemester.value) {
    list = list.filter((s: any) => s.semester === parseInt(selectedSemester.value))
  }
  return list
})

const loadClassCourses = async () => {
  loading.value = true
  try {
    const studentId = authStore.user?.id
    if (!studentId) throw new Error('未获取到学生ID')
    
    // 获取学生信息以获取班级ID
    const studentInfo = await studentApi.getStudent(studentId)
    if (studentInfo.code >= 200500 && studentInfo.code < 200600) {
      ElMessage.error(studentInfo.msg || '系统异常，请稍后重试!')
      courses.value = []
      loading.value = false
      return
    }
    
    if (!studentInfo.data?.class_id) {
      ElMessage.error('未获取到班级信息')
      courses.value = []
      loading.value = false
      return
    }
    
    // 获取班级课程
    const classCoursesResp = await studentApi.getClassCourse(
      studentInfo.data.class_id,
      selectedYear.value,
      selectedSemester.value ? parseInt(selectedSemester.value) : 1
    )
    
    if (classCoursesResp.code >= 200500 && classCoursesResp.code < 200600) {
      ElMessage.error(classCoursesResp.msg || '系统异常，请稍后重试!')
      courses.value = []
      loading.value = false
      return
    }
    
    if ((classCoursesResp.code === 0 || classCoursesResp.code === 200) && Array.isArray(classCoursesResp.data.courses)) {
      courses.value = classCoursesResp.data.courses
    } else {
      courses.value = []
    }
  } catch (error) {
    ElMessage.error('加载班级课程列表失败')
    courses.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  years.value = getYearOptions()
  loadClassCourses()
})
</script>

<style scoped>
.class-courses-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.class-courses-card-multi {
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  border: 1.5px solid #e0e7ff;
  background: #f8faff;
  width: 100%;
  max-width: 1200px;
  padding: 40px 36px 36px 36px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
.page-header-multi {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  color: #6c63ff;
  font-weight: 600;
}
.header-actions-multi {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 8px;
}
.select-multi {
  border-radius: 8px;
  min-width: 100px;
  margin-right: 8px;
}
.class-courses-table-multi {
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
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

:deep(.el-select) {
  border-radius: 8px;
}
</style> 