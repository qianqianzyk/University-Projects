<template>
  <div class="dashboard-content-multi">
    <div class="dashboard-row-multi">
      <el-card class="student-info-card-multi styled-card-multi">
        <div class="student-info-row-multi">
          <el-icon class="info-icon-multi"><User /></el-icon>
          <span>姓名：<b>{{ studentInfo.name }}</b></span>
          <span>学号：<b>{{ studentInfo.student_id }}</b></span>
          <span>性别：<b>{{ studentInfo.gender === 'M' ? '男' : studentInfo.gender === 'F' ? '女' : studentInfo.gender }}</b></span>
          <span>班级：<b>{{ studentInfo.class_name }}</b></span>
          <span>生源地：<b>{{ studentInfo.city_name }}</b></span>
          <span>年龄：<b>{{ studentInfo.age }}</b></span>
        </div>
      </el-card>
    </div>
    <el-row :gutter="32" class="dashboard-row-multi">
      <el-col :span="12">
        <el-card class="stat-card-multi styled-card-multi stat-card-green">
          <div class="stat-content-multi stat-center-multi">
            <div class="stat-icon-multi gpa-icon-multi">
              <el-icon><Reading /></el-icon>
            </div>
            <div class="stat-info-multi">
              <div class="stat-number-multi big-number-multi green-number">{{ myCourses.length }}</div>
              <div class="stat-label-multi">在学课程</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="stat-card-multi styled-card-multi stat-card-green">
          <div class="stat-content-multi stat-center-multi">
            <div class="stat-icon-multi gpa-icon-multi">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info-multi">
              <div class="stat-number-multi big-number-multi green-number">{{ stats.gpa }}</div>
              <div class="stat-label-multi">GPA</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <div class="dashboard-refresh-row">
      <el-button size="small" icon="el-icon-Refresh" @click="loadDashboardData" :loading="loading" class="refresh-btn-multi">刷新数据</el-button>
    </div>
    <el-row :gutter="32" class="dashboard-row-multi">
      <el-col :span="12">
        <el-card class="styled-card-multi">
          <template #header>
            <div class="card-header-multi"><el-icon style="vertical-align: middle;"><Reading /></el-icon> <span style="margin-left: 6px; font-weight: bold;">我的课程</span></div>
          </template>
          <el-table :data="myCourses" style="width: 100%" v-loading="loading" class="styled-table-multi">
            <el-table-column prop="course_name" label="课程名称" header-align="center" align="center" />
            <el-table-column prop="exam_type" label="考试类型" header-align="center" align="center" />
            <el-table-column label="授课教师" header-align="center" align="center">
              <template #default="{ row }">
                <span>{{ row.teachers && row.teachers.length ? row.teachers.map((t: any) => t.teacher_name).join('，') : '-' }}</span>
              </template>
            </el-table-column>
            <template #empty>
              <div style="text-align:center;color:#999;">暂无在学课程</div>
            </template>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="styled-card-multi">
          <template #header>
            <div class="card-header-multi"><el-icon style="vertical-align: middle;"><TrendCharts /></el-icon> <span style="margin-left: 6px; font-weight: bold;">最近成绩</span></div>
          </template>
          <el-table :data="recentScores" style="width: 100%" v-loading="loading" class="styled-table-multi">
            <el-table-column prop="course_name" label="课程名称" header-align="center" align="center" />
            <el-table-column prop="credits" label="学分" header-align="center" align="center" />
            <el-table-column prop="score" label="成绩" header-align="center" align="center" />
            <template #empty>
              <div style="text-align:center;color:#999;">暂无成绩数据</div>
            </template>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="courseDetailVisible" title="课程详情" width="400px">
      <div v-if="courseDetail">
        <div>课程名称：{{ courseDetail.name }}</div>
        <div>授课教师：{{ courseDetail.teacher }}</div>
        <div>学分：{{ courseDetail.credits }}</div>
        <div>学期：{{ courseDetail.semester }}</div>
      </div>
    </el-dialog>
    <el-dialog v-model="scoreDetailVisible" title="成绩详情" width="400px">
      <div v-if="scoreDetail">
        <div>课程名称：{{ scoreDetail.course_name }}</div>
        <div>成绩：{{ scoreDetail.score }}</div>
        <div>学期：{{ scoreDetail.semester }}</div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Reading, TrendCharts, User } from '@element-plus/icons-vue'
import { studentApi } from '../../api/student'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const stats = ref({
  gpa: 0
})
const studentInfo = ref({
  name: '',
  student_id: '',
  gender: '',
  class_name: '',
  city_name: '',
  age: ''
})
const loading = ref(false)
const myCourses = ref<any[]>([])
const recentScores = ref<any[]>([])
const courseDetailVisible = ref(false)
const courseDetail = ref<any>(null)
const scoreDetailVisible = ref(false)
const scoreDetail = ref<any>(null)

const currentYear = new Date().getFullYear()

const loadDashboardData = async () => {
  loading.value = true
  try {
    const studentId = authStore.user?.id
    if (!studentId) throw new Error('未获取到学生ID')
    // 获取学生信息
    const studentInfoResp = await studentApi.getStudent(studentId)
    const studentData = studentInfoResp.data as any
    if (studentInfoResp.code >= 200500 && studentInfoResp.code < 200600) {
      ElMessage.error(studentInfoResp.msg || '系统异常，请稍后重试!')
      loading.value = false
      return
    }
    if (studentInfoResp.code === 200 && studentData) {
      Object.assign(studentInfo.value, studentData)
      stats.value.gpa = studentData.gpa || 0
    }
    // 获取成绩
    const scoresResp = await studentApi.getScore(studentId, currentYear)
    if (scoresResp.code >= 200500 && scoresResp.code < 200600) {
      ElMessage.error(scoresResp.msg || '系统异常，请稍后重试!')
      myCourses.value = []
      recentScores.value = []
      loading.value = false
      return
    }
    if ((scoresResp.code === 0 || scoresResp.code === 200) && Array.isArray(scoresResp.data.scores)) {
      const allScores = scoresResp.data.scores.map((item: any) => ({
        ...item,
        credits: item.credit, // 字段映射
        teachers: item.teacher_list // 字段映射
      }))
      myCourses.value = allScores.filter(s => s.score === 0)
      recentScores.value = allScores.filter(s => typeof s.score === 'number' && s.score !== 0)
    } else {
      myCourses.value = []
      recentScores.value = []
    }
  } catch (error) {
    ElMessage.error('加载仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

const viewCourseDetail = (row: any) => {
  courseDetail.value = row
  courseDetailVisible.value = true
}
const viewScoreDetail = (row: any) => {
  scoreDetail.value = row
  scoreDetailVisible.value = true
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.dashboard-row-multi {
  width: 100%;
  margin-bottom: 24px;
}
.dashboard-refresh-row {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  margin: 16px 0 32px 0;
}
.styled-card-multi {
  background: #f8faff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  border: 1.5px solid #e0e7ff;
}
.student-info-card-multi {
  margin-bottom: 0;
}
.student-info-row-multi {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  font-size: 16px;
  color: #6c63ff;
  align-items: center;
  padding: 8px 0 2px 0;
}
.info-icon-multi {
  font-size: 22px;
  margin-right: 10px;
  color: #6c63ff;
}
.stat-card-multi {
  margin-bottom: 0;
  background: #f3e7fa;
  border: none;
}
.stat-card-purple {
  background: linear-gradient(135deg, #ede7f6 0%, #f3e7fa 100%);
  border: 1.5px solid #b39ddb;
}
.stat-card-green {
  background: linear-gradient(135deg, #e8f5e9 0%, #f3e7fa 100%);
  border: 1.5px solid #a5d6a7;
}
.stat-content-multi {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 80px;
}
.stat-center-multi {
  justify-content: center;
}
.stat-icon-multi {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 18px;
  font-size: 28px;
  color: white;
  box-shadow: 0 2px 8px 0 rgba(120,80,200,0.13);
}
.course-icon-multi {
  background: linear-gradient(135deg, #6c63ff 0%, #42a5f5 100%);
}
.gpa-icon-multi {
  background: linear-gradient(135deg, #43e695 0%, #6c63ff 100%);
}
.stat-info-multi {
  flex: 1;
  text-align: left;
}
.stat-number-multi {
  font-size: 28px;
  font-weight: bold;
  line-height: 1;
}
.big-number-multi {
  font-size: 38px;
  font-weight: 700;
}
.purple-number {
  color: #6c63ff;
}
.green-number {
  color: #43a047;
}
.stat-label-multi {
  font-size: 15px;
  color: #6c63ff;
}
.card-header-multi {
  font-size: 16px;
  color: #6c63ff;
  font-weight: 600;
  display: flex;
  align-items: center;
}
.styled-table-multi {
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}
.refresh-btn-multi {
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
  min-width: 110px;
}
.refresh-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
@media (max-width: 900px) {
  .student-info-row-multi {
    gap: 12px;
    font-size: 14px;
  }
}
</style> 