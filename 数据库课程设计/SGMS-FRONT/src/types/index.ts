// 基础响应类型
export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}

// 用户相关类型
export interface User {
  id: number;
  username: string;
  user_type: number; // 1:学生 2:教师 3:管理员
}

export interface LoginRequest {
  username: string;
  password: string;
  user_type: number;
}

export interface LoginResponse {
  id: number;
  user_type: number;
}

// 学生相关类型
export interface Student {
  id: number;
  name: string;
  student_id: string;
  class_id: number;
  department_id: number;
  gender: string;
  phone: string;
  email: string;
  address: string;
  birthday: string;
  enrollment_date: string;
  gpa: number;
}

export interface StudentListResponse {
  list: Student[];
  total: number;
}

// 教师相关类型
export interface Teacher {
  id: number;
  name: string;
  teacher_id: string;
  department_id: number;
  gender: string;
  phone: string;
  email: string;
  title: string;
  hire_date: string;
}

export interface TeacherListResponse {
  list: Teacher[];
  total: number;
}

// 班级相关类型
export interface ClassInfo {
  id: number;
  name: string;
  department_id: number;
  grade: string;
  student_count: number;
}

export interface ClassListResponse {
  list: ClassInfo[];
  total: number;
}

// 课程相关类型
export interface Course {
  id: number;
  name: string;
  code: string;
  credits: number;
  department_id: number;
  description: string;
  school_year?: number;
  class_id?: number;
  semester?: string;
}

export interface CourseListResponse {
  list: Course[];
  total: number;
}

// 部门相关类型
export interface Department {
  id: number;
  name: string;
  description: string;
}

export interface DepartmentListResponse {
  list: Department[];
  total: number;
}

// 授课相关类型
export interface Teaching {
  id: number;
  teacher_id: number;
  course_id: number;
  class_id: number;
  semester: string;
  academic_year: string;
}

export interface TeachingListResponse {
  list: Teaching[];
  total: number;
}

// 成绩相关类型
export interface Score {
  id: number;
  student_id: number;
  course_id: number;
  score: number;
  semester: string;
  academic_year: string;
}

export interface ScoreDistribution {
  score_range: string;
  count: number;
}

export interface AvgScore {
  course_id: number;
  course_name: string;
  avg_score: number;
}

// 统计分析类型
export interface CityStudentCount {
  city_id: number;
  city_name: string;
  student_count: number;
}

export interface Province {
  id: number;
  name: string;
}

export interface City {
  id: number;
  name: string;
  province_id: number;
}

export interface Score {
  id: number;
  student_id: number;
  course_id: number;
  score: number;
}


export interface GpaRank {
  student_id: number;
  student_name: string;
  gpa: number;
  rank: number;
}

export interface CityGpaTopTenCount {
  province_id: number;
  province_name: string;
  department_id: number;
  department_name: string;
  top10_count: number;
}

// 选课相关类型
export interface CourseSelection {
  id: number;
  student_id: number;
  course_id: number;
  selection_date: string;
}

// 教师统计类型
export interface TeacherStatistics {
  teacher_id: number;
  teacher_name: string;
  course_name: string;
  school_year: number;
  semester: string;
  course_count: number;
  student_count: number;
  avg_score: number;
} 