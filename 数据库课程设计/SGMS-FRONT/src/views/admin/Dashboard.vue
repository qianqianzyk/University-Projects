<template>
  <div class="page-root">
    <div class="dashboard">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon student-icon">
                <el-icon><UserFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.studentCount }}</div>
                <div class="stat-label">学生总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon teacher-icon">
                <el-icon><Avatar /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.teacherCount }}</div>
                <div class="stat-label">教师总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon course-icon">
                <el-icon><Collection /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.courseCount }}</div>
                <div class="stat-label">课程总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon class-icon">
                <el-icon><School /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.classCount }}</div>
                <div class="stat-label">班级总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon department-icon">
                <el-icon><OfficeBuilding /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.departmentCount }}</div>
                <div class="stat-label">部门总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon teaching-icon">
                <el-icon><Edit /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.teachingCount }}</div>
                <div class="stat-label">教师授课总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon province-icon">
                <el-icon><MapLocation /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.provinceCount }}</div>
                <div class="stat-label">省份总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon city-icon">
                <el-icon><House /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.cityCount }}</div>
                <div class="stat-label">城市总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon score-icon">
                <el-icon><Medal /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.scoreCount }}</div>
                <div class="stat-label">成绩总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <!-- 课程平均成绩统计 -->
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>课程平均成绩统计</span>
              </div>
            </template>
            <el-table
              v-loading="avgScoreLoading"
              :data="avgScoreWithClass"
              style="width: 100%"
              height="300"
            >
              <el-table-column prop="course_name" label="课程名称" />
              <el-table-column prop="class_name" label="班级" />
              <el-table-column prop="school_year" label="学年" />
              <el-table-column prop="avg_score" label="平均成绩">
                <template #default="{ row }">
                  {{ row.avg_score.toFixed(1) }}
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <!-- 生源地学生数量统计 -->
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>生源地学生数量统计</span>
              </div>
            </template>
            <el-table
              v-loading="cityCountLoading"
              :data="cityStudentCounts"
              style="width: 100%"
              height="300"
            >
              <el-table-column prop="province_name" label="省份" />
              <el-table-column prop="city_name" label="城市" />
              <el-table-column prop="student_count" label="学生数量" />
            </el-table>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <!-- GPA排名分析 -->
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>GPA排名分析</span>
                <div class="filter-section">
                  <el-select v-model="gpaRankType" placeholder="排名类型" size="small" @change="onRankTypeChange" style="width: 120px;">
                    <el-option label="专业排名" :value="1" />
                    <el-option label="班级排名" :value="2" />
                  </el-select>
                  <el-select v-model="selectedDepartment" placeholder="选择部门" size="small" @change="onDepartmentChange" style="width: 150px;">
                    <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
                  </el-select>
                </div>
              </div>
            </template>
            <!-- 班级切换按钮，仅在班级排名时显示 -->
            <div v-if="gpaRankType === 2 && classGroups.length > 0" class="class-navigation-container">
              <div class="class-navigation">
                <el-button size="small" @click="prevClass" :disabled="currentClassIndex === 0">
                  <el-icon><ArrowLeft /></el-icon>
                </el-button>
                <span class="class-info">{{ currentClassInfo }}</span>
                <el-button size="small" @click="nextClass" :disabled="currentClassIndex === classGroups.length - 1">
                  <el-icon><ArrowRight /></el-icon>
                </el-button>
              </div>
            </div>
            <el-table
              v-loading="gpaRankLoading"
              :data="currentGpaRanks"
              style="width: 100%"
              height="300"
            >
              <el-table-column prop="rank" label="排名" width="80" />
              <el-table-column prop="student_name" label="学生姓名" />
              <el-table-column prop="class_name" label="班级" />
              <el-table-column prop="gpa" label="GPA">
                <template #default="{ row }">
                  {{ row.gpa.toFixed(2) }}
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <!-- 生源地绩点前十统计 -->
        <el-col :span="12">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>生源地绩点前十统计</span>
              </div>
            </template>
            <el-table
              v-loading="topTenLoading"
              :data="cityGpaTopTenCounts"
              style="width: 100%"
              height="300"
            >
              <el-table-column prop="province_name" label="省份" />
              <el-table-column prop="department_name" label="专业" />
              <el-table-column prop="top10_count" label="前十名数量">
                <template #default="{ row }">
                  <span>{{ row.top10_count || '无数据' }}</span>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <!-- 教师授课统计 -->
        <el-col :span="24">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>教师授课统计</span>
              </div>
            </template>
            <el-table
              v-loading="teacherStatsLoading"
              :data="paginatedTeacherStatistics"
              style="width: 100%"
            >
              <el-table-column prop="teacher_name" label="教师姓名" />
              <el-table-column prop="course_name" label="课程名称" />
              <el-table-column prop="school_year" label="学年" />
              <el-table-column prop="semester" label="学期" />
              <el-table-column prop="student_count" label="学生数" />
              <el-table-column prop="avg_score" label="平均成绩">
                <template #default="{ row }">
                  {{ row.avg_score.toFixed(1) }}
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-container">
              <el-pagination
                v-model:current-page="teacherStatsCurrentPage"
                v-model:page-size="teacherStatsPageSize"
                :page-sizes="[10, 20, 50, 100]"
                :total="teacherStatistics.length"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleTeacherStatsSizeChange"
                @current-change="handleTeacherStatsCurrentChange"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <!-- 成绩分布统计 -->
        <el-col :span="24">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>成绩分布统计</span>
                <div class="filter-section">
                  <el-select v-model="selectedCourseForDistribution" placeholder="选择课程" size="small" @change="onCourseChange" style="width: 300px;">
                    <el-option
                      v-for="course in courseList"
                      :key="course.id"
                      :label="`${course.name}（${course.class_id ? classMap[course.class_id] : '未知班级'}）${course.school_year}年${course.semester === '1' ? '上' : '下'}学期`"
                      :value="course.id"
                    />
                  </el-select>
                  <el-select v-model="selectedSchoolYearForDistribution" placeholder="选择学年" size="small" @change="onSchoolYearChange" style="width: 120px;">
                    <el-option v-for="year in schoolYearList" :key="year" :label="year.toString()" :value="year" />
                  </el-select>
                </div>
              </div>
            </template>
            <div v-loading="scoreDistributionLoading" style="height: 400px;">
              <v-chart v-if="scoreDistributionData.length > 0" :option="scoreDistributionOption" style="height: 100%;" />
              <div v-else style="display: flex; justify-content: center; align-items: center; height: 100%; color: #909399;">
                请选择课程和学年查看成绩分布
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <!-- 学生成绩统计 -->
        <el-col :span="24">
          <el-card>
            <template #header>
              <div class="page-header">
                <span>学生成绩统计</span>
                <div class="filter-section">
                  <el-select v-model="selectedCourseForStudentScore" placeholder="选择课程" size="small" @change="onStudentScoreCourseChange" style="width: 300px;">
                    <el-option v-for="courseName in uniqueCourseNames" :key="courseName" :label="courseName" :value="courseName" />
                  </el-select>
                  <el-select v-model="selectedSchoolYearForStudentScore" placeholder="选择学年" size="small" @change="onStudentScoreSchoolYearChange" style="width: 120px;">
                    <el-option v-for="year in schoolYearList" :key="year" :label="year.toString()" :value="year" />
                  </el-select>
                </div>
              </div>
            </template>
            <el-table
              v-loading="studentScoreLoading"
              :data="studentScores"
              :empty-text="'请选择课程和学年查看学生成绩'"
              style="width: 100%"
              height="400"
            >
              <el-table-column prop="student_name" label="学生姓名" />
              <el-table-column prop="student_id" label="学生ID" />
              <el-table-column prop="course_name" label="课程名称" />
              <el-table-column prop="class_name" label="班级" />
              <el-table-column prop="school_year" label="学年" />
              <el-table-column prop="semester" label="学期" />
              <el-table-column prop="teacher_name" label="教师姓名">
                <template #default="{ row }">
                  <span>{{ Array.isArray(row.teacher_name) ? row.teacher_name.join(', ') : row.teacher_name }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="score" label="成绩">
                <template #default="{ row }">
                  <span :style="{ color: row.score >= 60 ? '#67C23A' : '#F56C6C' }">
                    {{ row.score }}
                  </span>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Avatar, School, ArrowLeft, ArrowRight, OfficeBuilding, House, UserFilled, Edit, MapLocation, Medal, Collection } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { adminApi } from '../../api/admin'
import { useAuthStore } from '../../stores/auth'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'

use([
  CanvasRenderer,
  BarChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const authStore = useAuthStore()

const stats = ref({
  studentCount: 0,
  teacherCount: 0,
  courseCount: 0,
  classCount: 0,
  departmentCount: 0,
  teachingCount: 0,
  provinceCount: 0,
  cityCount: 0,
  scoreCount: 0
})

// 加载状态
const avgScoreLoading = ref(false)
const cityCountLoading = ref(false)
const gpaRankLoading = ref(false)
const topTenLoading = ref(false)
const teacherStatsLoading = ref(false)
const scoreDistributionLoading = ref(false)

// 数据
const avgScores = ref<any[]>([])
const cityStudentCounts = ref<any[]>([])
const gpaRanks = ref<any[]>([])
const cityGpaTopTenCounts = ref<any[]>([])
const teacherStatistics = ref<any[]>([])
const departments = ref<any[]>([])
const classes = ref<any[]>([])

// 成绩分布统计相关数据
const courseList = ref<any[]>([])
const schoolYearList = ref<number[]>([])
const selectedCourseForDistribution = ref<number | null>(null)
const selectedSchoolYearForDistribution = ref<number | null>(null)
const scoreDistributionData = ref<any[]>([])

// 学生成绩统计相关数据
const selectedCourseForStudentScore = ref<string | null>(null)
const selectedSchoolYearForStudentScore = ref<number | null>(null)
const studentScores = ref<any[]>([])
const studentScoreLoading = ref(false)

// 班级分组相关
const classGroups = ref<any[]>([])
const currentClassIndex = ref(0)
const currentGpaRanks = ref<any[]>([])

// 教师授课统计分页
const teacherStatsCurrentPage = ref(1)
const teacherStatsPageSize = ref(10)

// 筛选条件
const gpaRankType = ref(1)
const selectedDepartment = ref<number | null>(null)
const selectedClass = ref<string | null>(null)

// 计算当前班级信息
const currentClassInfo = computed(() => {
  console.log('计算当前班级信息:', {
    classGroupsLength: classGroups.value.length,
    currentClassIndex: currentClassIndex.value,
    gpaRankType: gpaRankType.value,
    selectedClass: selectedClass.value
  })
  
  if (classGroups.value.length === 0) return ''
  const currentClass = classGroups.value[currentClassIndex.value]
  const info = `${currentClass.class_name} (${currentClassIndex.value + 1}/${classGroups.value.length})`
  console.log('当前班级信息:', info)
  return info
})

// 计算分页后的教师授课统计数据
const paginatedTeacherStatistics = computed(() => {
  const start = (teacherStatsCurrentPage.value - 1) * teacherStatsPageSize.value
  const end = start + teacherStatsPageSize.value
  return teacherStatistics.value.slice(start, end)
})

// 计算去重后的课程名称列表
const uniqueCourseNames = computed(() => {
  const courseNames = new Set<string>()
  courseList.value.forEach(course => {
    if (course.name) {
      courseNames.add(course.name)
    }
  })
  return Array.from(courseNames).sort()
})

// 课程id->班级名映射
const classMap = computed(() => {
  const map: Record<number, string> = {}
  classes.value.forEach(cls => {
    map[cls.id] = cls.name
  })
  return map
})

// 计算课程id到class_id的映射
const courseIdToClassId = computed(() => {
  const map: Record<number, number> = {}
  courseList.value.forEach(course => {
    if (course.id && course.class_id) {
      map[course.id] = course.class_id
    }
  })
  return map
})

// 课程平均成绩表格用的班级名
const avgScoreWithClass = computed(() => {
  return avgScores.value.map(item => {
    const classId = courseIdToClassId.value[item.course_id]
    const className = classId ? classMap.value[classId] : '未知班级'
    return {
      ...item,
      class_name: className
    }
  })
})

// 成绩分布图表配置
const scoreDistributionOption = computed(() => {
  if (scoreDistributionData.value.length === 0) {
    return {}
  }

  const xAxisData = scoreDistributionData.value.map(item => item.score_range)
  const seriesData = scoreDistributionData.value.map(item => item.count)

  return {
    title: {
      text: '成绩分布统计',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: xAxisData,
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      name: '学生人数'
    },
    series: [
      {
        name: '学生人数',
        type: 'bar',
        data: seriesData,
        itemStyle: {
          color: '#409EFF'
        }
      }
    ]
  }
})

const loadDashboardData = async () => {
  try {
    console.log('开始加载仪表盘数据...')
    console.log('authStore:', authStore)
    console.log('authStore.user:', authStore.user)
    const adminId = authStore.user?.id || 1
    console.log('管理员ID:', adminId)

    // 获取统计数据
    console.log('开始调用API...')
    const results = await Promise.allSettled([
      adminApi.getStudents(adminId),
      adminApi.getTeachers(adminId),
      adminApi.getCourses(adminId),
      adminApi.getClasses(adminId)
    ])
    console.log('API调用完成')
    
    // 处理结果
    const [studentsResult, teachersResult, coursesResult, classesResult] = results
    
    console.log('API调用结果:', results)
    
    // 更新统计数据
    if (studentsResult.status === 'fulfilled' && (studentsResult.value.code === 200 || studentsResult.value.code === 0)) {
      console.log('学生响应:', studentsResult.value)
      console.log('学生数据详情:', studentsResult.value.data)
      console.log('学生列表:', studentsResult.value.data.list)
      console.log('学生列表长度:', studentsResult.value.data.list?.length)
      stats.value.studentCount = studentsResult.value.data.list?.length || 0
    } else {
      console.log('学生API失败:', studentsResult)
    }
    
    if (teachersResult.status === 'fulfilled' && (teachersResult.value.code === 200 || teachersResult.value.code === 0)) {
      console.log('教师数据详情:', teachersResult.value.data)
      stats.value.teacherCount = teachersResult.value.data?.list?.length || 0
    } else {
      console.log('教师API失败:', teachersResult)
    }
    
    if (coursesResult.status === 'fulfilled' && (coursesResult.value.code === 200 || coursesResult.value.code === 0)) {
      console.log('课程数据详情:', coursesResult.value.data)
      stats.value.courseCount = coursesResult.value.data?.list?.length || 0
    } else {
      console.log('课程API失败:', coursesResult)
    }
    
    if (classesResult.status === 'fulfilled' && (classesResult.value.code === 200 || classesResult.value.code === 0)) {
      console.log('班级数据详情:', classesResult.value.data)
      stats.value.classCount = classesResult.value.data?.list?.length || 0
    } else {
      console.log('班级API失败:', classesResult)
    }

  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    // 如果API调用失败，使用默认数据
    stats.value = {
      studentCount: 1250,
      teacherCount: 89,
      courseCount: 156,
      classCount: 45,
      departmentCount: 0,
      teachingCount: 0,
      provinceCount: 0,
      cityCount: 0,
      scoreCount: 0
    }
  }
}

// 加载课程平均成绩
const loadAvgScores = async () => {
  avgScoreLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getAvgScore(adminId)
    if (response.code === 200 || response.code === 0) {
      avgScores.value = response.data.avg_scores || []
    }
  } catch (error) {
    ElMessage.error('加载课程平均成绩失败')
  } finally {
    avgScoreLoading.value = false
  }
}

// 加载生源地学生数量
const loadCityStudentCount = async () => {
  cityCountLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCityStudentCount(adminId)
    if (response.code === 200 || response.code === 0) {
      cityStudentCounts.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('加载生源地学生数量失败')
  } finally {
    cityCountLoading.value = false
  }
}

// 处理排名类型选择变化
const onRankTypeChange = () => {
  console.log('排名类型改变，重新请求API确保数据一致性')
  
  // 重置班级相关状态
  currentClassIndex.value = 0
  classGroups.value = []
  selectedClass.value = null
  
  // 重新请求API以确保数据一致性
  loadGpaRank()
}

// 处理部门选择变化
const onDepartmentChange = () => {
  console.log('部门选择改变，强制重新请求API')
  
  // 强制重新请求API，完全重置状态
  currentClassIndex.value = 0
  classGroups.value = []
  selectedClass.value = null
  currentGpaRanks.value = []
  
  // 加载GPA排名
  loadGpaRank()
}

// 加载GPA排名
const loadGpaRank = async () => {
  if (!selectedDepartment.value) {
    // 如果没有选择部门，使用默认值
    if (departments.value.length > 0) {
      selectedDepartment.value = departments.value[0].id
    } else {
      return
    }
  }
  
  console.log('加载GPA排名，部门ID:', selectedDepartment.value, '排名类型:', gpaRankType.value)
  
  gpaRankLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    // 确保selectedDepartment.value不为null
    if (selectedDepartment.value) {
      // 无论选择什么排名类型，都请求专业排名数据（type=1）
      const response = await adminApi.getDepartmentGpaRank(adminId, selectedDepartment.value, 1)
      if (response.code === 200 || response.code === 0) {
        gpaRanks.value = response.data.rank || []
        console.log('获取到GPA排名数据:', gpaRanks.value)
        
        // 如果是班级排名，进行分组处理
        if (gpaRankType.value === 2) {
          processClassGroups()
        } else {
          // 专业排名直接显示
          currentGpaRanks.value = gpaRanks.value
          classGroups.value = []
          currentClassIndex.value = 0
          selectedClass.value = null
          console.log('专业排名模式，直接显示所有数据')
        }
      } else if (response.code >= 200500 && response.code < 200600) {
        ElMessage.error(response.msg || '系统异常，请稍后重试!')
        gpaRankLoading.value = false
        return
      }
    }
  } catch (error) {
    ElMessage.error('加载GPA排名失败')
  } finally {
    gpaRankLoading.value = false
  }
}

// 处理班级分组
const processClassGroups = () => {
  console.log('开始处理班级分组，原始数据:', gpaRanks.value)
  
  // 确保从干净的状态开始
  classGroups.value = []
  currentClassIndex.value = 0
  selectedClass.value = null
  
  const groups: { [key: string]: any[] } = {}
  
  // 按班级分组
  gpaRanks.value.forEach((student) => {
    const className = student.class_name
    console.log(`学生 ${student.student_name} 属于班级: ${className}`)
    if (!groups[className]) {
      groups[className] = []
    }
    // 为每个班级内的学生重新计算排名
    groups[className].push({
      ...student,
      rank: groups[className].length + 1
    })
  })
  
  // 转换为数组格式
  classGroups.value = Object.keys(groups).map(className => ({
    class_name: className,
    students: groups[className]
  }))
  
  console.log('班级分组结果:', classGroups.value)
  console.log('班级列表:', classGroups.value.map(c => c.class_name))
  
  // 设置默认班级选择
  if (classGroups.value.length > 0) {
    selectedClass.value = classGroups.value[0].class_name
    console.log('设置默认班级:', selectedClass.value)
  }
  
  // 更新当前显示的排名
  updateCurrentGpaRanks()
}

// 更新当前显示的GPA排名
const updateCurrentGpaRanks = () => {
  console.log('更新当前显示的GPA排名:', {
    gpaRankType: gpaRankType.value,
    classGroupsLength: classGroups.value.length,
    currentClassIndex: currentClassIndex.value,
    selectedClass: selectedClass.value,
    classGroups: classGroups.value.map(c => c.class_name)
  })
  
  if (gpaRankType.value === 2 && classGroups.value.length > 0 && selectedClass.value) {
    const currentClass = classGroups.value.find(c => c.class_name === selectedClass.value)
    if (currentClass) {
      currentGpaRanks.value = currentClass.students || []
      console.log('班级排名模式，当前班级:', selectedClass.value, '学生数量:', currentGpaRanks.value.length)
    } else {
      console.error('找不到选中的班级:', selectedClass.value)
      currentGpaRanks.value = []
    }
  } else {
    currentGpaRanks.value = gpaRanks.value
    console.log('专业排名模式，学生数量:', currentGpaRanks.value.length)
  }
}

// 加载生源地绩点前十统计
const loadCityGpaTopTenCount = async () => {
  topTenLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCityGpaTopTenCount(adminId)
    console.log('生源地绩点前十统计响应:', response)
    if (response.code === 200 || response.code === 0) {
      cityGpaTopTenCounts.value = response.data.list || []
      console.log('生源地绩点前十统计数据:', cityGpaTopTenCounts.value)
    }
  } catch (error) {
    console.error('加载生源地绩点前十统计失败:', error)
    ElMessage.error('加载生源地绩点前十统计失败')
  } finally {
    topTenLoading.value = false
  }
}

// 加载教师授课统计
const loadTeacherStatistics = async () => {
  teacherStatsLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeacherStatistics(adminId)
    if (response.code === 200 || response.code === 0) {
      teacherStatistics.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('加载教师授课统计失败')
  } finally {
    teacherStatsLoading.value = false
  }
}

// 加载部门列表
const loadDepartments = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getDepartments(adminId)
    if (response.code === 200 || response.code === 0) {
      departments.value = response.data.list || []
      if (departments.value.length > 0) {
        selectedDepartment.value = departments.value[0].id
        // 加载默认部门的GPA排名
        await loadGpaRank()
      }
    }
  } catch (error) {
    ElMessage.error('加载部门列表失败')
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

// 切换到上一个班级
const prevClass = () => {
  console.log('切换到上一个班级，当前索引:', currentClassIndex.value)
  if (currentClassIndex.value > 0) {
    currentClassIndex.value--
    selectedClass.value = classGroups.value[currentClassIndex.value].class_name
    updateCurrentGpaRanks()
    console.log('切换到班级:', selectedClass.value)
  }
}

// 切换到下一个班级
const nextClass = () => {
  console.log('切换到下一个班级，当前索引:', currentClassIndex.value)
  if (currentClassIndex.value < classGroups.value.length - 1) {
    currentClassIndex.value++
    selectedClass.value = classGroups.value[currentClassIndex.value].class_name
    updateCurrentGpaRanks()
    console.log('切换到班级:', selectedClass.value)
  }
}

// 处理教师授课统计分页大小变化
const handleTeacherStatsSizeChange = (size: number) => {
  teacherStatsPageSize.value = size
  teacherStatsCurrentPage.value = 1
}

// 处理教师授课统计当前页变化
const handleTeacherStatsCurrentChange = (page: number) => {
  teacherStatsCurrentPage.value = page
}

// 加载部门总数
const loadDepartmentCount = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getDepartments(adminId)
    if (response.code === 200 || response.code === 0) {
      stats.value.departmentCount = response.data.list?.length || 0
    }
  } catch (error) {
    console.error('加载部门总数失败:', error)
  }
}

// 加载教师授课总数
const loadTeachingCount = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getTeachings(adminId)
    if (response.code === 200 || response.code === 0) {
      stats.value.teachingCount = response.data.list?.length || 0
    }
  } catch (error) {
    console.error('加载教师授课总数失败:', error)
  }
}

// 加载省份总数
const loadProvinceCount = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getProvinces(adminId)
    if (response.code === 200 || response.code === 0) {
      stats.value.provinceCount = response.data.provinces?.length || 0
    }
  } catch (error) {
    console.error('加载省份总数失败:', error)
  }
}

// 加载城市总数
const loadCityCount = async () => {
  try {
    const adminId = authStore.user?.id || 1
    // 直接获取所有城市
    const citiesResponse = await adminApi.getCities(adminId)
    if (citiesResponse.code === 200 || citiesResponse.code === 0) {
      stats.value.cityCount = citiesResponse.data.cities?.length || 0
    }
  } catch (error) {
    console.error('加载城市总数失败:', error)
  }
}

// 加载成绩总数
const loadScoreCount = async () => {
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getScores(adminId)
    
    // 检查响应结构
    if (response && (response.code === 200 || response.code === 0)) {
      // 如果data是undefined，可能数据直接在response中
      const scoreList = response.data?.list || []
      stats.value.scoreCount = scoreList.length
    } else {
      stats.value.scoreCount = 0
    }
  } catch (error) {
    stats.value.scoreCount = 0
  }
}

// 加载课程列表
const loadCourseList = async () => {
  try {
    const adminId = authStore.user?.id || 1
    
    // 获取课程列表
    const coursesResponse = await adminApi.getCourses(adminId)
    if (coursesResponse.code === 200 || coursesResponse.code === 0) {
      const courses = coursesResponse.data.list || []
      
      // 获取教师授课统计信息，用于关联课程和教师
      const teacherStatsResponse = await adminApi.getTeacherStatistics(adminId)
      if (teacherStatsResponse.code === 200 || teacherStatsResponse.code === 0) {
        const teacherStats = teacherStatsResponse.data.list || []
        
        // 为每个课程添加教师信息
        courseList.value = courses.map(course => {
          // 查找该课程的教师信息
          const courseTeacher = teacherStats.find(stat => 
            stat.course_name === course.name && 
            stat.school_year === course.school_year
          )
          
          return {
            ...course,
            teacher_name: courseTeacher?.teacher_name || '未分配教师'
          }
        })
      } else {
        courseList.value = courses
      }
      
      // 提取学年列表
      const years = new Set<number>()
      courseList.value.forEach(course => {
        if (course.school_year) {
          years.add(course.school_year)
        }
      })
      schoolYearList.value = Array.from(years).sort((a, b) => b - a) // 降序排列
    }
  } catch (error) {
    ElMessage.error('加载课程列表失败')
  }
}

// 加载成绩分布数据
const loadScoreDistribution = async () => {
  if (!selectedCourseForDistribution.value || !selectedSchoolYearForDistribution.value) {
    return
  }

  scoreDistributionLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCourseScoreDistribution(
      adminId,
      selectedCourseForDistribution.value,
      selectedSchoolYearForDistribution.value
    )
    if (response.code === 200 || response.code === 0) {
      scoreDistributionData.value = response.data.distribution || []
    }
  } catch (error) {
    ElMessage.error('加载成绩分布数据失败')
    scoreDistributionData.value = []
  } finally {
    scoreDistributionLoading.value = false
  }
}

// 处理课程选择变化
const onCourseChange = () => {
  if (selectedCourseForDistribution.value && selectedSchoolYearForDistribution.value) {
    loadScoreDistribution()
  }
}

// 处理学年选择变化
const onSchoolYearChange = () => {
  if (selectedCourseForDistribution.value && selectedSchoolYearForDistribution.value) {
    loadScoreDistribution()
  }
}

// 加载学生成绩数据
const loadStudentScores = async () => {
  if (!selectedCourseForStudentScore.value || !selectedSchoolYearForStudentScore.value) {
    return
  }

  studentScoreLoading.value = true
  try {
    const adminId = authStore.user?.id || 1
    const response = await adminApi.getCourseAllStudentScore(
      adminId,
      selectedCourseForStudentScore.value,
      selectedSchoolYearForStudentScore.value
    )
    if (response.code === 200 || response.code === 0) {
      // 为每个学生成绩记录添加学年信息
      studentScores.value = (response.data.scores || []).map(score => ({
        ...score,
        school_year: selectedSchoolYearForStudentScore.value
      }))
    }
  } catch (error) {
    ElMessage.error('加载学生成绩数据失败')
    studentScores.value = []
  } finally {
    studentScoreLoading.value = false
  }
}

// 处理学生成绩统计课程选择变化
const onStudentScoreCourseChange = () => {
  if (selectedCourseForStudentScore.value && selectedSchoolYearForStudentScore.value) {
    loadStudentScores()
  }
}

// 处理学生成绩统计学年选择变化
const onStudentScoreSchoolYearChange = () => {
  if (selectedCourseForStudentScore.value && selectedSchoolYearForStudentScore.value) {
    loadStudentScores()
  }
}

// 初始化数据
const initData = async () => {
  await Promise.all([
    loadDashboardData(),
    loadAvgScores(),
    loadCityStudentCount(),
    loadCityGpaTopTenCount(),
    loadTeacherStatistics(),
    loadDepartments(),
    loadDepartmentCount(),
    loadTeachingCount(),
    loadProvinceCount(),
    loadCityCount(),
    loadScoreCount(),
    loadCourseList()
  ])
}

onMounted(async () => {
  try {
    await initData();
    await loadClasses();
  } catch (error) {
    console.error('初始化数据失败:', error);
  }
});
</script>

<style scoped>
.page-root {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.dashboard {
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}

.stat-card {
  margin-bottom: 20px;
  background: #f8faff;
  border-radius: 20px;
  box-shadow: 0 6px 32px rgba(120, 80, 200, 0.10), 0 1.5px 6px 0 rgba(0,0,0,0.03);
  border: 1.5px solid #e0e7ff;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(120, 80, 200, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 32px 0;
}

.stat-icon {
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

.student-icon {
  background: linear-gradient(135deg, #6c63ff 0%, #42a5f5 100%);
}

.teacher-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.course-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.class-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.department-icon {
  background: linear-gradient(135deg, #ffd700 0%, #ffa500 100%);
}

.teaching-icon {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff4757 100%);
}

.province-icon {
  background: linear-gradient(135deg, #00b894 0%, #00cec9 100%);
}

.city-icon {
  background: linear-gradient(135deg, #a29bfe 0%, #6c5ce7 100%);
}

.score-icon {
  background: linear-gradient(135deg, #fd79a8 0%, #e84393 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #6c63ff;
  line-height: 1;
}

.stat-label {
  font-size: 15px;
  color: #6c63ff;
  margin-top: 5px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #6c63ff;
  font-weight: 600;
  font-size: 16px;
}

.filter-section {
  display: flex;
  gap: 10px;
}

.class-navigation {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 10px;
}

.class-info {
  font-size: 12px;
  color: #6c63ff;
  background: linear-gradient(135deg, #f8faff 0%, #e0e7ff 100%);
  padding: 6px 12px;
  border-radius: 8px;
  min-width: 80px;
  text-align: center;
  border: 1px solid #e0e7ff;
}

.class-navigation-container {
  margin-top: 10px;
  margin-bottom: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
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

:deep(.el-button) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(120, 80, 200, 0.15);
}

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
</style> 