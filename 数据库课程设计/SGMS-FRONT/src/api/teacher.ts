/**
 * 教师API接口
 * Base URL: http://121.43.236.83:8888
 */

import api from './config'

// 类型定义
export interface Teacher {
  id: number // 主键ID
  teacher_id: string // 工号
  name: string // 姓名
  gender: string // 性别
  age: number // 年龄
  title: string // 职称
  phone: string // 电话
  is_admin: boolean // 是否管理员
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

export interface TeacherInCourse {
  teacher_id: number // 教师ID
  teacher_name: string // 教师姓名
}

export interface CourseRank {
  course_id: number // 课程ID
  course_name: string // 课程名称
  score: number // 成绩
  student_id: string // 学生ID
  student_name: string // 学生姓名
  class_id: number // 班级ID
  class_name: string // 班级名称
  school_year: number // 学年
  semester: number // 学期
  rank: number // 排名
}

export interface CourseAvgScore {
  course_name: string // 课程名称
  class_name: string // 班级名称
  school_year: number // 学年
  avg_score: number // 平均成绩
}

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

export interface SetStudentScoreReq {
  teacher_id: number // 教师ID
  student_id: string // 学生ID
  course_id: number // 课程ID
  score: number // 成绩
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  msg: string // 更新为后端实际返回的字段
  data: T
}

export interface TeacherResponse {
  teacher: Teacher
}

export interface CourseResponse {
  courses: Course[]
}

export interface CourseRankResponse {
  rank: CourseRank[]
}

export interface CourseAvgScoreResponse {
  avg_scores: CourseAvgScore[]
}

export interface StudentResponse {
  student: Student
}

export interface StudentListResponse {
  students: Student[]
}

/**
 * 教师API接口
 */
export const teacherApi = {
  /**
   * 获取教师信息
   * @param id 教师ID
   * @returns 教师信息
   */
  getTeacher(id: number): Promise<ApiResponse<TeacherResponse>> {
    return api.get('/api/teacher', { params: { id } })
  },

  /**
   * 获取教师所授课程
   * @param teacher_id 教师ID
   * @param year 学年，可选
   * @param semester 学期，可选
   * @returns 课程列表
   */
  getCourse(teacher_id: number, year?: number, semester?: number): Promise<ApiResponse<CourseResponse>> {
    const params: any = { teacher_id }
    if (year) params.year = year
    if (semester) params.semester = semester
    return api.get('/api/teacher/course', { params })
  },

  /**
   * 获取教师所授课程平均成绩
   * @param teacher_id 教师ID
   * @returns 课程平均成绩
   */
  getCourseAvgScore(teacher_id: number): Promise<ApiResponse<CourseAvgScoreResponse>> {
    return api.get('/api/teacher/course/avgscore', { params: { teacher_id } })
  },

  /**
   * 获取教师所授课程成绩排名
   * @param teacher_id 教师ID
   * @param course_id 课程ID
   * @param year 学年
   * @param semester 学期
   * @returns 课程成绩排名
   */
  getCourseRank(teacher_id: number, course_id: number, year: number, semester: number): Promise<ApiResponse<CourseRankResponse>> {
    return api.get('/api/teacher/course/rank', { params: { teacher_id, course_id, year, semester } })
  },

  /**
   * 获取教师所授课程下的学生列表
   * @param teacher_id 教师ID
   * @param course_id 课程ID
   * @param year 学年，可选
   * @param semester 学期，可选
   * @returns 学生列表
   */
  getCourseStudentList(teacher_id: number, course_id: number, year?: number, semester?: number): Promise<ApiResponse<StudentListResponse>> {
    const params: any = { teacher_id, course_id }
    if (year) params.year = year
    if (semester) params.semester = semester
    return api.get('/api/teacher/course/student', { params })
  },

  /**
   * 获取教师所授课程下的学生信息
   * @param teacher_id 教师ID
   * @param student_id 学生ID
   * @returns 学生信息
   */
  getStudent(teacher_id: number, student_id: string): Promise<ApiResponse<StudentResponse>> {
    return api.get('/api/teacher/student', { params: { teacher_id, student_id } })
  },

  /**
   * 设置学生成绩
   * @param data 设置成绩参数
   * @returns 设置结果
   */
  setStudentScore(data: SetStudentScoreReq): Promise<ApiResponse> {
    return api.post('/api/teacher/course/student/score', data)
  }
} 