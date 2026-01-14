<template>
  <div class="courses-content-multi">
    <el-card class="courses-card-multi">
      <template #header>
        <div class="page-header-multi">
          <span>我的课程</span>
          <div class="header-actions-multi">
            <el-select v-model="selectedYear" placeholder="选择学年" size="small" class="select-multi" @change="loadCourses">
              <el-option v-for="year in years" :key="year" :label="year" :value="year" />
            </el-select>
            <el-select v-model="selectedSemester" :placeholder="selectedSemester === '' ? '全部' : '选择学期'" size="small" class="select-multi" @change="loadCourses">
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
        class="courses-table-multi"
        style="width: 100%"
      >
        <el-table-column prop="course_name" label="课程名称" />
        <el-table-column prop="credits" label="学分" />
        <el-table-column prop="school_year" label="学年" />
        <el-table-column prop="semester" label="学期" />
        <el-table-column prop="hours" label="课时" />
        <el-table-column prop="exam_type" label="考试类型" />
        <el-table-column label="授课教师">
          <template #default="{ row }">
            <span>{{ row.teachers && row.teachers.length ? row.teachers.map((t: any) => t.teacher_name).join('，') : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="retake_required" label="需重修">
          <template #default="{ row }">
            <span :class="row.retake_required ? 'red-text' : 'green-text'">{{ row.retake_required ? '是' : '否' }}</span>
          </template>
        </el-table-column>
        <template #empty>
          <div style="text-align:center;color:#999;">暂无课程数据</div>
        </template>
      </el-table>
      <el-dialog v-model="courseDetailVisible" title="课程详情" width="400px" class="dialog-multi">
        <div v-if="courseDetail">
          <div>课程名称：{{ courseDetail.name }}</div>
          <div>课程代码：{{ courseDetail.code }}</div>
          <div>授课教师：{{ courseDetail.teacher }}</div>
          <div>学分：{{ courseDetail.credits }}</div>
          <div>学期：{{ courseDetail.semester }}</div>
        </div>
      </el-dialog>
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
const courseDetailVisible = ref(false)
const courseDetail = ref<any>(null)
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
    list = list.filter((s: any) => s.semester === selectedSemester.value)
  }
  return list
})

const loadCourses = async () => {
  loading.value = true
  try {
    const studentId = authStore.user?.id
    if (!studentId) throw new Error('未获取到学生ID')
    const scoresResp = await studentApi.getScore(studentId, selectedYear.value)
    if (scoresResp.code >= 200500 && scoresResp.code < 200600) {
      ElMessage.error(scoresResp.msg || '系统异常，请稍后重试!')
      courses.value = []
      loading.value = false
      return
    }
    if ((scoresResp.code === 0 || scoresResp.code === 200) && Array.isArray(scoresResp.data.scores)) {
      courses.value = scoresResp.data.scores.map((item: any) => ({
        ...item,
        credits: item.credit, // 字段映射
        teachers: item.teacher_list // 字段映射
      }))
    } else {
      courses.value = []
    }
  } catch (error) {
    ElMessage.error('加载课程列表失败')
  } finally {
    loading.value = false
  }
}

const viewCourseDetail = (row: any) => {
  courseDetail.value = row
  courseDetailVisible.value = true
}

onMounted(() => {
  years.value = getYearOptions()
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
.select-btn-multi {
  background: linear-gradient(90deg, #6c63ff 0%, #42e695 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(76, 201, 240, 0.08);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 90px;
}
.select-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
.courses-table-multi {
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}
.dialog-multi {
  border-radius: 16px;
}
.dialog-header-multi {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}
.red-text {
  color: #e53935;
  font-weight: bold;
}
.green-text {
  color: #43a047;
  font-weight: bold;
}
@media (max-width: 900px) {
  .page-header-multi {
    font-size: 15px;
  }
  .header-actions-multi {
    gap: 4px;
  }
  .courses-card-multi {
    max-width: 98vw;
    padding: 24px 8px 16px 8px;
  }
}
</style> 