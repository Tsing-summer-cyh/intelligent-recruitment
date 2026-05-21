import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/',
    redirect: '/jobs' // 🌟 修改点：默认进入职位大厅，而不是登录页
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '候选人登录' }
  },
  {
    path: '/jobs',
    name: 'JobBoard',
    component: () => import('../views/JobBoard.vue'),
    meta: { requiresAuth: false, title: '职位大厅' } // 🌟 确保这里是 false
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true, role: 'candidate', title: '我的档案' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  document.title = to.meta.title || '智能招聘系统'
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('role')

  if (to.path === '/login' && token && userRole === 'candidate') {
    return next('/jobs')
  }

  if (to.meta.requiresAuth) {
    if (!token) {
      ElMessage.warning('请先登录！')
      return next('/login')
    }
    if (userRole !== 'candidate') {
      ElMessage.error('该页面仅限求职者访问！')
      localStorage.clear()
      return next('/login')
    }
  }
  next()
})

export default router