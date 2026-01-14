/**
 * 学生API接口
 * Base URL: http://121.43.236.83:8888
 */

import api from './config'

// 类型定义
export interface Student {
  id: number // 学生ID
  student_id: string // 学号
  name: string // 姓名
  gender: string // 性别
  age: number // 年龄
  city_id: number // 生源地城市编号
  class_id: number // 所属班级编号
  gpa: number // 学生平均绩点
  total_credits: number // 已修总学分
}

export interface UpdateStudentReq {
  id: number // 学生ID 用来指定要更新的学生
  name: string // 姓名
  gender: string // 性别
  age: number // 年龄
  city_id: number // 城市ID
}

export interface Course {
  course_id: number // 课程ID
  course_name: string // 课程名称
  credits: number // 学分
  year: number // 学年
  semester: number // 学期
  hours: number // 课时
  exam_type: string // 考试类型 1:考试 2:考查
  teachers: Teacher[] // 教师列表
}

export interface Teacher {
  teacher_id: number // 教师ID
  teacher_name: string // 教师姓名
}

export interface Score {
  course_id: number // 课程ID
  course_name: string // 课程名称
  school_year: number // 学年
  semester: string // 学期
  hours: number // 课时
  exam_type: string // 考试类型
  teachers: Teacher[] // 教师列表
  score: number // 成绩
  credits: number // 学分
  retake_required: boolean // 是否需要重修
}

export interface SelectCourseReq {
  student_id: number // 学生ID
  course_id: number // 课程ID
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  msg: string // 更新为后端实际返回的字段
  data: T
}

export interface StudentResponse {
  id: number
  student_id: string
  name: string
  gender: string
  age: number
  city_id: number
  city_name: string
  class_id: number
  class_name: string
  gpa: number
  total_credits: number
}

export interface CourseResponse {
  courses: Course[]
}

export interface ScoreResponse {
  scores: Score[]
}

/**
 * 学生API接口
 */
export const studentApi = {
  /**
   * 获取学生信息
   * @param id 学生ID
   * @returns 学生信息
   */
  getStudent(id: number): Promise<ApiResponse<StudentResponse>> {
    return api.get('/api/student/info', { params: { id } })
  },

  /**
   * 修改学生信息
   * @param data 更新学生参数
   * @returns 更新结果
   */
  updateStudent(data: UpdateStudentReq): Promise<ApiResponse> {
    return api.put('/api/student/info', data)
  },

  /**
   * 获取学生课程成绩
   * @param id 学生ID
   * @param year 学年，可选
   * @returns 成绩列表
   */
  getScore(id: number, year?: number): Promise<ApiResponse<ScoreResponse>> {
    const params: any = { id }
    if (year) params.year = year
    return api.get('/api/student/score', { params })
  },

  /**
   * 获取班级课程列表
   * @param class_id 班级ID
   * @param year 学年
   * @param semester 学期
   * @param student_id 学生ID 如果传了代表查询该学生的可选课程
   * @returns 课程列表
   */
  getClassCourse(class_id: number, year: number, semester: number, student_id?: number): Promise<ApiResponse<CourseResponse>> {
    const params: any = { class_id, year, semester }
    if (student_id) params.student_id = student_id
    return api.get('/api/student/class/course', { params })
  },

  /**
   * 获取教师课程列表
   * @param teacher_id 教师ID
   * @param year 学年，可选
   * @param semester 学期，可选
   * @returns 课程列表
   */
  getTeacherCourse(teacher_id: number, year?: number, semester?: number): Promise<ApiResponse<CourseResponse>> {
    const params: any = { teacher_id }
    if (year) params.year = year
    if (semester) params.semester = semester
    return api.get('/api/student/teacher/course', { params })
  },

  /**
   * 学生选课
   * @param data 选课参数
   * @returns 选课结果
   */
  selectCourse(data: SelectCourseReq): Promise<ApiResponse> {
    return api.post('/api/student/course/select', data)
  }
} 