<template>
  <div class="dashboard-content-multi">
    <el-row :gutter="32" class="dashboard-row-multi">
      <el-col :span="12">
        <el-card class="stat-card-multi styled-card-multi stat-card-purple">
          <div class="stat-content-multi stat-center-multi">
            <div class="stat-icon-multi course-icon-multi">
              <el-icon><Reading /></el-icon>
            </div>
            <div class="stat-info-multi">
              <div class="stat-number-multi big-number-multi purple-number">{{ stats.courseCount }}</div>
              <div class="stat-label-multi">授课课程</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="stat-card-multi styled-card-multi stat-card-green">
          <div class="stat-content-multi stat-center-multi">
            <div class="stat-icon-multi gpa-icon-multi">
              <el-icon><School /></el-icon>
            </div>
            <div class="stat-info-multi">
              <div class="stat-number-multi big-number-multi green-number">{{ stats.classCount }}</div>
              <div class="stat-label-multi">授课班级</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="32" class="dashboard-row-multi" style="margin-top: 32px;">
      <el-col :span="12">
        <el-card class="styled-card-multi dashboard-table-card-multi">
          <template #header>
            <div class="card-header-multi">
              <span>我的课程</span>
              <el-button size="small" type="primary" class="multi-btn-multi" @click="viewAllCourses">查看详情</el-button>
            </div>
          </template>
          <div class="my-courses-table-multi">
            <el-table :data="pagedCourses" class="dashboard-table-multi" border stripe>
              <el-table-column prop="course_name" label="课程名称" min-width="120" />
              <el-table-column prop="class_name" label="班级" min-width="100" />
              <el-table-column prop="credits" label="学分" min-width="80" />
              <el-table-column prop="hours" label="课时" min-width="80" />
              <el-table-column prop="exam_type" label="考试类型" min-width="100">
                <template #default="{ row }">
                  {{ row.exam_type === '1' ? '考试' : '考查' }}
                </template>
              </el-table-column>
            </el-table>
          </div>
          <el-pagination
            v-if="myCourses.length > pageSize"
            v-model:current-page="coursePage"
            :page-size="pageSize"
            :total="myCourses.length"
            layout="prev, pager, next"
            style="margin-top: 10px; text-align: right;"
          />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="styled-card-multi dashboard-table-card-multi">
          <template #header>
            <div class="card-header-multi">
              <span>课程成绩统计</span>
            </div>
          </template>
          <div class="course-stats-table-multi">
            <el-table :data="pagedCourseStats" class="dashboard-table-multi" border stripe>
              <el-table-column prop="course_name" label="课程" min-width="120" />
              <el-table-column prop="avg_score" label="平均分" min-width="80" />
              <el-table-column prop="class_name" label="班级" min-width="100" />
              <el-table-column prop="school_year" label="学年" min-width="100" />
            </el-table>
          </div>
          <el-pagination
            v-if="courseStats.length > pageSize"
            v-model:current-page="courseStatsPage"
            :page-size="pageSize"
            :total="courseStats.length"
            layout="prev, pager, next"
            style="margin-top: 10px; text-align: right;"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Reading, User, Document, School, Avatar } from '@element-plus/icons-vue'
import { teacherApi } from '../../api/teacher'
import type { Course, CourseAvgScore } from '../../api/teacher'
import { useAuthStore } from '../../stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const stats = ref({
  courseCount: 0,
  classCount: 0
})

const myCourses = ref<Course[]>([])
const courseStats = ref<CourseAvgScore[]>([])
const coursePage = ref(1)
const courseStatsPage = ref(1)
const pageSize = 10

const pagedCourses = computed(() => {
  const start = (coursePage.value - 1) * pageSize
  return myCourses.value.slice(start, start + pageSize)
})
const pagedCourseStats = computed(() => {
  const start = (courseStatsPage.value - 1) * pageSize
  return courseStats.value.slice(start, start + pageSize)
})

const loadDashboardData = async () => {
  try {
    const teacherId = authStore.user?.id
    if (!teacherId) return

    // 获取当前年份
    const year = new Date().getFullYear()

    // 分别请求学期1和2的课程
    const [coursesRes1, coursesRes2] = await Promise.all([
      teacherApi.getCourse(teacherId, year, 1),
      teacherApi.getCourse(teacherId, year, 2)
    ])

    if (coursesRes1.code >= 200500 && coursesRes1.code < 200600) {
      ElMessage.error(coursesRes1.msg || '系统异常，请稍后重试!')
      return
    }

    let allCourses: Course[] = []
    if (coursesRes1.code === 200 || coursesRes1.code === 0 && Array.isArray(coursesRes1.data.courses)) {
      allCourses = allCourses.concat(coursesRes1.data.courses)
    }
    if (coursesRes2.code === 200 || coursesRes2.code === 0 && Array.isArray(coursesRes2.data.courses)) {
      allCourses = allCourses.concat(coursesRes2.data.courses)
    }
    myCourses.value = allCourses // 展示全部课程
    stats.value.courseCount = myCourses.value.length

    // 加载课程平均成绩
    const avgScoreResponse = await teacherApi.getCourseAvgScore(teacherId)
    if (avgScoreResponse.code === 200 || avgScoreResponse.code === 0 && Array.isArray(avgScoreResponse.data.avg_scores)) {
      let avgScores = avgScoreResponse.data.avg_scores
      if (avgScores.length > 0) {
        // 按学年倒序排序
        avgScores = avgScores.slice().sort((a, b) => (b.school_year || 0) - (a.school_year || 0))
        courseStats.value = avgScores // 展示全部课程统计
        // 授课班级数量为所有班级名去重后的数量
        const classSet = new Set(avgScores.map(item => item.class_name))
        stats.value.classCount = classSet.size
      }
    }
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
  }
}

const viewAllCourses = () => {
  router.push('/teacher/courses')
}

const goToCourses = () => {
  router.push('/teacher/courses')
}

const goToStudents = () => {
  router.push('/teacher/students')
}

const goToProfile = () => {
  router.push('/teacher/profile')
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
  margin: 0;
}
.styled-card-multi {
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
.stat-card-multi {
  background: #f8faff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  border: 1.5px solid #e0e7ff;
  padding: 0;
  min-width: 0;
  max-width: 100%;
  margin-bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
.stat-content-multi {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 0;
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
  color: #fff;
}
.course-icon-multi {
  background: linear-gradient(135deg, #6c63ff 0%, #42a5f5 100%);
}
.gpa-icon-multi {
  background: linear-gradient(135deg, #42e695 0%, #3bb2b8 100%);
}
.stat-info-multi {
  flex: 1;
}
.stat-number-multi {
  font-size: 32px;
  font-weight: bold;
  line-height: 1;
}
.big-number-multi {
  font-size: 38px;
}
.purple-number {
  color: #6c63ff;
}
.green-number {
  color: #43a047;
}
.stat-label-multi {
  font-size: 15px;
  color: #666;
  margin-top: 5px;
}
.card-header-multi {
  font-weight: bold;
  color: #1976d2;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
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
  min-width: 120px;
}
.multi-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
.dashboard-table-card-multi {
  min-width: 340px;
  max-width: 1200px;
  width: 100%;
  margin-bottom: 0;
  margin-top: 0;
  padding: 32px 24px 24px 24px;
  border-radius: 18px;
  box-shadow: 0 2px 12px rgba(120, 80, 200, 0.08);
  border: 1.5px solid #e0e7ff;
  background: #fff;
}
.dashboard-table-multi {
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  background: #fff;
}
@media (max-width: 900px) {
  .dashboard-content-multi {
    padding: 16px 4px 0 4px;
  }
  .styled-card-multi, .dashboard-table-card-multi {
    max-width: 98vw;
    padding: 24px 8px 16px 8px;
  }
}
</style> 