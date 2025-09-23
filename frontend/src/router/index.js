import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '@/views/LandingPage.vue'
import ChatPage from '@/views/ChatPage.vue'
import CharacterChatPage from '@/views/CharacterChatPage.vue'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: LandingPage
  },
  {
    path: '/home',
    name: 'Home',
    component: ChatPage
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
