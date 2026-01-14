<template>
  <div class="students-content-multi">
    <el-card class="students-card-multi styled-card-multi">
      <template #header>
        <div class="students-title-multi">学生管理</div>
      </template>
      <el-table
        v-loading="loading"
        :data="students"
        class="students-table-multi"
        border stripe
        style="width: 100%"
      >
        <el-table-column prop="student_name" label="姓名" min-width="120" />
        <el-table-column prop="class_name" label="班级" min-width="120" />
        <el-table-column prop="score" label="成绩" min-width="100" />
        <el-table-column prop="rank" label="班级排名" min-width="100" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <div class="action-btns-vertical">
              <el-button size="small" class="multi-btn-multi" @click="viewStudent(row)">查看详情</el-button>
              <el-button size="small" type="primary" class="multi-btn-multi" @click="setScore(row)">设置成绩</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <!-- 其余弹窗等结构保持不变 -->
    <el-dialog
      v-model="studentDialogVisible"
      title="学生详细信息"
      width="500px"
    >
      <el-card v-if="selectedStudent" class="student-detail-card">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="学号">{{ selectedStudent.student_id }}</el-descriptions-item>
          <el-descriptions-item label="姓名">{{ selectedStudent.name }}</el-descriptions-item>
          <el-descriptions-item label="性别">{{ selectedStudent.gender === 'M' ? '男' : selectedStudent.gender === 'F' ? '女' : selectedStudent.gender }}</el-descriptions-item>
          <el-descriptions-item label="年龄">{{ selectedStudent.age }}</el-descriptions-item>
          <el-descriptions-item label="城市">{{ (selectedStudent as any).city_name }}</el-descriptions-item>
          <el-descriptions-item label="班级">{{ (selectedStudent as any).class_name }}</el-descriptions-item>
          <el-descriptions-item label="GPA">{{ selectedStudent.gpa }}</el-descriptions-item>
          <el-descriptions-item label="已修学分">{{ selectedStudent.total_credits }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-dialog>
    <el-dialog
      v-model="scoreDialogVisible"
      title="设置成绩"
      width="400px"
    >
      <el-form v-if="selectedStudent" :model="scoreForm" label-width="80px">
        <el-form-item label="学生姓名">
          <span>{{ selectedStudent.name }}</span>
        </el-form-item>
        <el-form-item label="成绩">
          <el-input v-model="scoreForm.score" type="number" min="0" max="100" placeholder="请输入成绩" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="multi-btn-multi" @click="submitScore">保存成绩</el-button>
          <el-button @click="scoreDialogVisible = false">取消</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { teacherApi } from '../../api/teacher'
import type { Student } from '../../api/teacher'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const students = ref<Student[]>([])
const selectedYear = ref('')
const selectedSemester = ref('')
const years = ref<number[]>([])
const studentDialogVisible = ref(false)
const selectedStudent = ref<Student | null>(null)
const scoreDialogVisible = ref(false)
const scoreForm = ref({ score: 0 as number })

const loadStudents = async () => {
  loading.value = true
  try {
    const teacherId = authStore.user?.id
    const courseId = route.query.courseId ? Number(route.query.courseId) : undefined
    
    if (!teacherId || !courseId) {
      ElMessage.error('缺少必要参数')
      return
    }

    const response = await teacherApi.getCourseStudentList(
      teacherId,
      courseId,
      selectedYear.value ? Number(selectedYear.value) : undefined,
      selectedSemester.value ? Number(selectedSemester.value) : undefined
    )
    
    if ((response.code === 200 || response.code === 0) && Array.isArray(response.data.students)) {
      students.value = response.data.students
      
      // 生成年份列表（这里可能需要从课程数据中获取）
      years.value = [2024, 2023, 2022]
    } else {
      students.value = []
    }
  } catch (error) {
    ElMessage.error('加载学生列表失败')
  } finally {
    loading.value = false
  }
}

const viewStudent = async (student: Student) => {
  try {
    const teacherId = authStore.user?.id
    if (!teacherId) return
    // 调用获取学生详细信息接口
    const response = await teacherApi.getStudent(teacherId, student.student_id)
    console.log('getStudent response:', response)
    if (response.code === 200 || response.code === 0) {
      selectedStudent.value = response.data as unknown as Student
      studentDialogVisible.value = true
      console.log('selectedStudent:', selectedStudent.value)
      console.log('studentDialogVisible:', studentDialogVisible.value)
    }
  } catch (error) {
    ElMessage.error('获取学生详情失败')
  }
}

const setScore = async (student: Student) => {
  try {
    const teacherId = authStore.user?.id
    if (!teacherId) return
    // 直接设置学生信息并弹窗
    selectedStudent.value = student
    scoreForm.value.score = (student as any).score || 0
    scoreDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取学生详情失败')
  }
}

const submitScore = async () => {
  try {
    const teacherId = authStore.user?.id
    const courseId = route.query.courseId ? Number(route.query.courseId) : undefined
    if (!teacherId || !courseId || !selectedStudent.value) return
    // 调用设置成绩接口
    const response = await teacherApi.setStudentScore({
      teacher_id: teacherId,
      student_id: selectedStudent.value.student_id,
      course_id: courseId,
      score: Number(scoreForm.value.score)
    })
    if (response.code === 200 || response.code === 0) {
      ElMessage.success('成绩设置成功')
      scoreDialogVisible.value = false
      loadStudents() // 刷新列表
    }
  } catch (error) {
    ElMessage.error('设置成绩失败')
  }
}

onMounted(() => {
  loadStudents()
})
</script>

<style scoped>
.students-content-multi {
  width: 100%;
  min-height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 0 0 0;
  box-sizing: border-box;
}
.students-card-multi {
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
.students-title-multi {
  font-size: 22px;
  font-weight: 700;
  color: #6c63ff;
  margin-bottom: 24px;
}
.students-table-multi {
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 90px;
  margin-right: 8px;
}
.multi-btn-multi:hover {
  background: linear-gradient(90deg, #42e695 0%, #6c63ff 100%);
  color: #fff;
}
@media (max-width: 900px) {
  .students-card-multi {
    max-width: 98vw;
    padding: 24px 8px 16px 8px;
  }
}
/* 操作按钮竖直对齐 */
.action-btns-vertical {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.action-btns-vertical .multi-btn-multi {
  min-width: 0;
  width: auto;
  padding-left: 12px;
  padding-right: 12px;
  margin-left: 0;
  box-shadow: none;
}
</style> 