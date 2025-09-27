import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '@/views/LandingPage.vue'
import LoginPage from '@/views/LoginPage.vue'
import ChatPage from '@/views/ChatPage.vue'
import DebugPage from '@/views/DebugPage.vue'
import chatService from '@/services/chatService.js'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: LandingPage
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage
  },
  {
    path: '/home',
    name: 'Home',
    component: ChatPage
  },
  {
    path: '/debug',
    name: 'Debug',
    component: DebugPage
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 检查登录状态
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const user = localStorage.getItem('user')
  
  // 如果访问需要登录的页面
  if (to.path === '/home') {
    if (!token || !user) {
      // 未登录，跳转到登录页
      next('/login')
    } else {
      // 已登录，允许访问
      next()
    }
  } else if (to.path === '/login') {
    // 允许访问登录页，不强制跳转
    next()
  } else {
    // 其他页面，直接通过
    next()
  }
})

export default router
