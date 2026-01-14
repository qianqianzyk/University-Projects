<template>
  <div class="scores-content-multi">
    <el-card class="scores-card-multi">
      <template #header>
        <div class="page-header-multi">
          <span>我的成绩</span>
          <div>
            <el-select v-model="selectedYear" placeholder="选择学年" size="small" class="select-multi" @change="loadScores">
              <el-option v-for="year in years" :key="year" :label="year" :value="year" />
            </el-select>
            <el-select v-model="selectedSemester" :placeholder="selectedSemester === '' ? '全部' : '选择学期'" size="small" class="select-multi" @change="loadScores">
              <el-option label="全部" value="" />
              <el-option label="上学期" value="1" />
              <el-option label="下学期" value="2" />
            </el-select>
          </div>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="filteredScores"
        class="scores-table-multi"
        style="width: 100%"
      >
        <el-table-column prop="course_name" label="课程名称" />
        <el-table-column prop="credits" label="学分" />
        <el-table-column prop="score" label="成绩" />
        <el-table-column prop="school_year" label="学年" />
        <el-table-column prop="semester" label="学期" />
        <el-table-column prop="exam_type" label="考试类型" />
        <el-table-column label="授课教师">
          <template #default="{ row }">
            <span>{{ row.teachers && row.teachers.length ? row.teachers.map((t: any) => t.teacher_name).join('，') : '-' }}</span>
          </template>
        </el-table-column>
        <template #empty>
          <div style="text-align:center;color:#999;">暂无成绩数据</div>
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { studentApi } from '../../api/student'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const scores = ref<any[]>([])
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

const loadScores = async () => {
  loading.value = true
  try {
    const studentId = authStore.user?.id
    if (!studentId) throw new Error('未获取到学生ID')
    const response = await studentApi.getScore(studentId, selectedYear.value)
    if (response.code >= 200500 && response.code < 200600) {
      ElMessage.error(response.msg || '系统异常，请稍后重试!')
      scores.value = []
      loading.value = false
      return
    }
    if ((response.code === 0 || response.code === 200) && Array.isArray(response.data.scores)) {
      scores.value = response.data.scores.map((item: any) => ({
        ...item,
        credits: item.credit,
        teachers: item.teacher_list
      }))
      const apiYears = response.data.scores.map((s: any) => s.school_year)
      const defaultYears = getYearOptions()
      years.value = Array.from(new Set([...apiYears, ...defaultYears])).filter(Boolean).sort((a, b) => b - a)
    } else {
      scores.value = []
    }
  } catch (error) {
    ElMessage.error('加载成绩列表失败')
  } finally {
    loading.value = false
  }
}

const filteredScores = computed(() => {
  let list = scores.value.filter((s: any) => s.score > 0)
  if (selectedYear.value) {
    list = list.filter((s: any) => s.school_year === selectedYear.value)
  }
  if (selectedSemester.value) {
    list = list.filter((s: any) => s.semester === selectedSemester.value)
  }
  return list
})

onMounted(() => {
  years.value = getYearOptions()
  loadScores()
})
</script>

<style scoped>
.scores-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.scores-card-multi {
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
.select-multi {
  border-radius: 8px;
  min-width: 100px;
  margin-right: 8px;
}
.scores-table-multi {
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}
@media (max-width: 900px) {
  .page-header-multi {
    font-size: 15px;
  }
  .scores-card-multi {
    max-width: 98vw;
    padding: 24px 8px 16px 8px;
  }
}
</style> 