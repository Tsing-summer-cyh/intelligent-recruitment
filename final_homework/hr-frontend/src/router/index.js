import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录 - 智能招聘系统' }
  },
  {
    // HR 管理端路由组
    path: '/hr',
    name: 'HRLayout',
    component: () => import('../views/hr/Layout.vue'),
    meta: { requiresAuth: true, role: 'hr' },
    children: [
      {
        path: 'dashboard',
        name: 'HRDashboard',
        component: () => import('../views/hr/Dashboard.vue'),
        meta: { title: '工作台 - HR管理端' }
      },
      {
        path: 'stats',
        name: 'HRStats',
        component: () => import('../views/hr/Stats.vue'),
        meta: { title: '统计分析 - HR管理端' }
      },
      {
        path: 'profiles',
        name: 'HRProfiles',
        component: () => import('../views/hr/Profiles.vue'),
        meta: { title: '候选人档案 - HR管理端' }
      },
      {
        path: 'applications',
        name: 'HRApplications',
        component: () => import('../views/hr/Applications.vue'),
        meta: { title: '投递记录 - HR管理端' }
      },
      {
        path: 'account',
        name: 'HRAccount',
        component: () => import('../views/hr/Account.vue'),
        meta: { title: '账户设置 - HR管理端' }
      },
      {
        path: 'ai-history',
        name: 'HRAIHistory',
        component: () => import('../views/hr/AIHistory.vue'),
        meta: { title: 'AI对话历史 - HR管理端' }
      }
    ]
  },
  {
    // 候选人端路由组
    path: '/candidate',
    name: 'CandidateLayout',
    component: () => import('../views/candidate/Layout.vue'),
    meta: { requiresAuth: true, role: 'candidate' },
    children: [
      {
        path: 'jobs',
        name: 'CandidateJobs',
        component: () => import('../views/candidate/Jobs.vue'),
        meta: { title: '职位大厅 - 智能招聘' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫：类似后端的 JWT 中间件
router.beforeEach((to, from, next) => {
  // 设置浏览器标签页标题
  document.title = to.meta.title || '智能招聘系统'

  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('role')

  // 如果去登录页，且已经有 token 了，直接踢回对应的首页
  if (to.path === '/login' && token) {
    if (userRole === 'hr') return next('/hr/dashboard')
    if (userRole === 'candidate') return next('/candidate/jobs')
  }

  // 如果访问的页面需要登录 (requiresAuth)
  if (to.meta.requiresAuth) {
    if (!token) {
      ElMessage.warning('请先登录！')
      return next('/login')
    }
    // 检查角色是否匹配，防止候选人强行访问 HR 页面
    if (to.meta.role && to.meta.role !== userRole) {
      ElMessage.error('权限不足，禁止跨角色访问！')
      localStorage.clear()
      return next('/login')
    }
  }

  next() // 放行
})

export default router