import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        name: 'AdminLayout',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'students',
        name: 'AdminStudents',
        component: () => import('../views/admin/Students.vue')
      },
      {
        path: 'teachers',
        name: 'AdminTeachers',
        component: () => import('../views/admin/Teachers.vue')
      },
      {
        path: 'classes',
        name: 'AdminClasses',
        component: () => import('../views/admin/Classes.vue')
      },
      {
        path: 'courses',
        name: 'AdminCourses',
        component: () => import('../views/admin/Courses.vue')
      },
      {
        path: 'departments',
        name: 'AdminDepartments',
        component: () => import('../views/admin/Departments.vue')
      },
      {
        path: 'teachings',
        name: 'AdminTeachings',
        component: () => import('../views/admin/Teachings.vue')
      }
    ]
  },
  {
    path: '/student',
    component: () => import('../views/student/StudentLayout.vue'),
    meta: { requiresAuth: true, role: 'student' },
    children: [
      {
        path: '',
        name: 'StudentLayout',
        redirect: '/student/dashboard'
      },
      {
        path: 'dashboard',
        name: 'StudentDashboard',
        component: () => import('../views/student/Dashboard.vue')
      },
      {
        path: 'courses',
        name: 'StudentCourses',
        component: () => import('../views/student/Courses.vue')
      },
      {
        path: 'class-courses',
        name: 'StudentClassCourses',
        component: () => import('../views/student/ClassCourses.vue')
      },
      {
        path: 'scores',
        name: 'StudentScores',
        component: () => import('../views/student/Scores.vue')
      },
      {
        path: 'profile',
        name: 'StudentProfile',
        component: () => import('../views/student/Profile.vue')
      }
    ]
  },
  {
    path: '/teacher',
    component: () => import('../views/teacher/TeacherLayout.vue'),
    meta: { requiresAuth: true, role: 'teacher' },
    children: [
      {
        path: '',
        name: 'TeacherLayout',
        redirect: '/teacher/dashboard'
      },
      {
        path: 'dashboard',
        name: 'TeacherDashboard',
        component: () => import('../views/teacher/Dashboard.vue')
      },
      {
        path: 'courses',
        name: 'TeacherCourses',
        component: () => import('../views/teacher/Courses.vue')
      },
      {
        path: 'students',
        name: 'TeacherStudents',
        component: () => import('../views/teacher/Students.vue')
      },
      {
        path: 'profile',
        name: 'TeacherProfile',
        component: () => import('../views/teacher/Profile.vue')
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 路由守卫
router.beforeEach((to, from, next) => {
  try {
    const authStore = useAuthStore();

    // 初始化认证状态
    if (!authStore.isLoggedIn()) {
      authStore.initAuth();
    }

    // 检查是否需要认证
    if (to.meta.requiresAuth && !authStore.isLoggedIn()) {
      next('/login');
      return;
    }

    // 检查角色权限
    if (to.meta.role && !authStore.hasRole(to.meta.role as string)) {
      // 根据用户角色重定向到对应的首页
      if (authStore.hasRole('admin')) {
        next('/admin/dashboard');
      } else if (authStore.hasRole('student')) {
        next('/student/dashboard');
      } else if (authStore.hasRole('teacher')) {
        next('/teacher/dashboard');
      } else {
        next('/login');
      }
      return;
    }

    // 如果已登录且访问登录页，重定向到对应的首页
    if (to.path === '/login' && authStore.isLoggedIn()) {
      if (authStore.hasRole('admin')) {
        next('/admin/dashboard');
      } else if (authStore.hasRole('student')) {
        next('/student/dashboard');
      } else if (authStore.hasRole('teacher')) {
        next('/teacher/dashboard');
      }
      return;
    }

    next();
  } catch (error) {
    console.error('路由守卫错误:', error);
    next('/login');
  }
});

export default router; 