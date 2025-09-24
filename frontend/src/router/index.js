import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '@/views/LandingPage.vue'
import LoginPage from '@/views/LoginPage.vue'
import Home from '@/views/Home.vue'
import CharacterChatPage from '@/views/CharacterChatPage.vue'

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
    component: Home
  },
  {
    path: '/chat',
    name: 'CharacterChat',
    component: CharacterChatPage
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
