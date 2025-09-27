<template>
  <div class="debug-page">
    <h1>调试页面</h1>
    
    <div class="debug-section">
      <h2>登录状态</h2>
      <p>Token: {{ token || '无' }}</p>
      <p>User: {{ user ? JSON.stringify(user, null, 2) : '无' }}</p>
      <p>ChatService User: {{ chatService.currentUser ? JSON.stringify(chatService.currentUser, null, 2) : '无' }}</p>
    </div>

    <div class="debug-section">
      <h2>好友列表</h2>
      <p>好友数量: {{ chatService.friends.length }}</p>
      <div v-for="friend in chatService.friends" :key="friend.id">
        <p>{{ friend.name }} (ID: {{ friend.id }}, Character ID: {{ friend.character_id }})</p>
      </div>
    </div>

    <div class="debug-section">
      <h2>当前聊天</h2>
      <p>选中聊天: {{ chatService.currentChat ? JSON.stringify(chatService.currentChat, null, 2) : '无' }}</p>
      <p>消息数量: {{ chatService.messages.length }}</p>
    </div>

    <div class="debug-section">
      <h2>操作</h2>
      <button @click="clearLogin">清除登录状态</button>
      <button @click="loadFriends">重新加载好友</button>
      <button @click="goToLogin">跳转登录页</button>
      <button @click="goToHome">跳转首页</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import chatService from '@/services/chatService.js'

const router = useRouter()
const token = ref('')
const user = ref(null)

onMounted(() => {
  token.value = localStorage.getItem('token')
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
})

const clearLogin = () => {
  chatService.forceLogout()
  token.value = ''
  user.value = null
  console.log('登录状态已清除')
}

const loadFriends = async () => {
  try {
    await chatService.loadUserFriends()
    console.log('好友列表已重新加载')
  } catch (error) {
    console.error('加载好友失败:', error)
  }
}

const goToLogin = () => {
  router.push('/login')
}

const goToHome = () => {
  router.push('/home')
}
</script>

<style scoped>
.debug-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.debug-section {
  margin: 20px 0;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.debug-section h2 {
  margin-top: 0;
  color: #333;
}

button {
  margin: 5px;
  padding: 10px 15px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

button:hover {
  background: #0056b3;
}
</style>

