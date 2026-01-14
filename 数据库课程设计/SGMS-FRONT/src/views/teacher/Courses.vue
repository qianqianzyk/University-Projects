<template>
  <div class="courses-content-multi">
    <el-card class="courses-card-multi styled-card-multi">
      <template #header>
        <div class="courses-header-multi">
          <span class="courses-title-multi">我的课程</span>
          <div class="header-actions-multi">
            <el-select v-model="selectedYear" placeholder="学年" size="small" class="select-multi" @change="loadCourses">
              <el-option label="全部" value="" />
              <el-option v-for="year in years" :key="year" :label="year" :value="year" />
            </el-select>
            <el-select v-model="selectedSemester" placeholder="学期" size="small" class="select-multi" @change="loadCourses">
              <el-option label="全部" value="" />
              <el-option label="上学期" value="1" />
              <el-option label="下学期" value="2" />
            </el-select>
          </div>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="courses"
        class="courses-table-multi"
        border stripe
        style="width: 100%"
      >
        <el-table-column prop="course_name" label="课程名称" min-width="140" />
        <el-table-column prop="class_name" label="班级" min-width="120" />
        <el-table-column prop="credits" label="学分" min-width="80" />
        <el-table-column prop="school_year" label="学年" min-width="100" />
        <el-table-column prop="semester" label="学期" min-width="100">
          <template #default="{ row }">
            {{ row.semester === 1 ? '上学期' : '下学期' }}
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="课时" min-width="80" />
        <el-table-column prop="exam_type" label="考试类型" min-width="100" />
        <el-table-column label="任课教师" min-width="120">
          <template #default="{ row }">
            {{ teacherNames(row.teachers) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" class="multi-btn-multi" @click="viewStudents(row)">查看学生</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { teacherApi } from '../../api/teacher'
import type { Course } from '../../api/teacher'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const courses = ref<Course[]>([])
const currentYear = new Date().getFullYear()
const years = ref<number[]>([])
const selectedYear = ref(String(currentYear))
const selectedSemester = ref('')
const teacherProfile = ref<{ teacher_id: string }>({ teacher_id: '' })

// 初始化学年列表 2020~今年
for (let y = currentYear; y >= 2020; y--) {
  years.value.push(y)
}

const fetchTeacherProfile = async () => {
  const userId = authStore.user?.id
  if (!userId) return
  const res = await teacherApi.getTeacher(userId)
  if ((res.code === 200 || res.code === 0) && res.data && res.data.teacher) {
    teacherProfile.value.teacher_id = res.data.teacher.teacher_id
  }
}

const loadCourses = async () => {
  loading.value = true
  try {
    const teacherId = authStore.user?.id
    if (!teacherId) return

    const year = selectedYear.value ? Number(selectedYear.value) : currentYear
    let allCourses: Course[] = []
    if (!selectedSemester.value) {
      // 学期为全部，分别请求上学期和下学期
      const [res1, res2] = await Promise.all([
        teacherApi.getCourse(teacherId, year, 1),
        teacherApi.getCourse(teacherId, year, 2)
      ])
      if ((res1.code === 200 || res1.code === 0) && Array.isArray(res1.data.courses)) {
        allCourses = allCourses.concat(res1.data.courses)
      }
      if ((res2.code === 200 || res2.code === 0) && Array.isArray(res2.data.courses)) {
        allCourses = allCourses.concat(res2.data.courses)
      }
    } else {
      // 只请求选中的学期
      const res = await teacherApi.getCourse(teacherId, year, Number(selectedSemester.value))
      if ((res.code === 200 || res.code === 0) && Array.isArray(res.data.courses)) {
        allCourses = res.data.courses
      }
    }
    courses.value = allCourses
  } catch (error) {
    ElMessage.error('加载课程列表失败')
  } finally {
    loading.value = false
  }
}

const viewStudents = (course: Course) => {
  router.push({
    path: '/teacher/students',
    query: {
      courseId: course.course_id.toString(),
      teacherId: String(teacherProfile.value.teacher_id)
    }
  })
}

function teacherNames(teachers: Array<{ teacher_name: string }>) {
  return teachers && teachers.length ? teachers.map(t => t.teacher_name).join('，') : '-';
}

onMounted(async () => {
  await fetchTeacherProfile()
  loadCourses()
})
</script>

<style scoped>
.courses-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.courses-card-multi {
  background: #f8faff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  padding: 40px 36px 36px 36px;
  min-width: 340px;
  max-width: 1200px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  border: 1.5px solid #e0e7ff;
}
.courses-header-multi {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 40px;
}
.courses-title-multi {
  font-size: 22px;
  font-weight: 700;
  color: #43a047;
  margin-bottom: 0;
}
.header-actions-multi {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 0;
}
.select-multi {
  border-radius: 8px;
  min-width: 90px;
  max-width: 110px;
  margin-right: 0;
}
.courses-table-multi {
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}
.multi-btn-multi {
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
  min-width: 90px;
}
.multi-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
@media (max-width: 900px) {
  .courses-card-multi {
    max-width: 98vw;
    padding: 24px 8px 16px 8px;
  }
}
</style> 