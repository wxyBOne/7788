<template>
  <div class="sidebar">
      <div class="logo-section">
        <div class="logo">seven</div>
      </div>
    
    <nav class="nav-menu">
      <div class="nav-item" :class="{ active: activeSection === 'home' }" @click="$emit('setActiveSection', 'home')">
        <div class="nav-icon home-icon"></div>
      </div>
      <div class="nav-item" :class="{ active: activeSection === 'messages' }" @click="$emit('setActiveSection', 'messages')">
        <div class="nav-icon messages-icon"></div>
      </div>
      <div class="nav-item" :class="{ active: activeSection === 'settings' }" @click="$emit('setActiveSection', 'settings')">
        <div class="nav-icon settings-icon"></div>
      </div>
    </nav>
    
    <div class="user-section" @click="toggleUserCard">
      <div class="user-avatar">
        <img :src="userAvatar" :alt="userName" />
      </div>
    </div>

    <!-- AI伙伴初始化弹窗 -->
    <CompanionInitModal 
      :show="showCompanionInit" 
      @close="showCompanionInit = false"
      @created="onCompanionCreated"
    />

    <!-- 用户信息卡片弹窗 -->
    <div v-if="showUserCard" class="user-card-overlay" @click="closeUserCard">
      <div class="user-card" @click.stop>
        <div class="user-card-header">
          <div class="user-card-avatar">
            <img :src="userAvatar" :alt="userName" />
          </div>
          <div class="user-card-info">
            <h3 class="user-name">{{ userName }}</h3>
            <p class="user-email">{{ userEmail }}</p>
          </div>
        </div>
        <div class="user-card-actions">
          <button class="switch-account-btn" @click="switchAccount">
            <svg class="switch-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="8.5" cy="7" r="4"/>
              <line x1="20" y1="8" x2="20" y2="14"/>
              <line x1="23" y1="11" x2="17" y2="11"/>
            </svg>
            切换账号
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import chatService from '@/services/chatService.js'
import CompanionInitModal from './CompanionInitModal.vue'

defineProps({
  activeSection: {
    type: String,
    default: 'messages'
  },
  companionEmotion: {
    type: Object,
    default: () => ({
      emotion: '平静',
      intensity: 0.5,
      color: '#52b4b4',
      brightness: 0.7,
      particle_speed: 0.5
    })
  }
})

defineEmits(['setActiveSection'])

// 弹窗控制
const showUserCard = ref(false)
const showCompanionInit = ref(false)

// 用户信息
const userAvatar = computed(() => {
  return chatService.currentUser?.avatar_url || '/src/img/DefaultUserAvatar.jpg'
})

const userName = computed(() => {
  return chatService.currentUser?.name || chatService.currentUser?.username || '用户'
})

const userEmail = computed(() => {
  return chatService.currentUser?.email || 'user@example.com'
})

// 切换用户卡片显示
const toggleUserCard = () => {
  showUserCard.value = !showUserCard.value
}

// 关闭用户卡片
const closeUserCard = () => {
  showUserCard.value = false
}

// 处理AI伙伴创建
const onCompanionCreated = async (companion) => {
  console.log('AI伙伴创建成功:', companion)
  // 刷新好友列表
  await chatService.loadUserFriends()
  // 可以选择直接进入AI伙伴的聊天
  // emit('setActiveSection', 'messages')
}

// 切换账号
const switchAccount = () => {
  // 清除本地存储的登录信息
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  
  // 清除chatService中的用户信息
  chatService.currentUser = null
  chatService.friends = []
  chatService.currentChat = null
  chatService.messages = []
  
  // 关闭弹窗
  closeUserCard()
  
  // 跳转到登录页面
  window.location.href = '/login'
}

// 组件挂载时检查用户信息
onMounted(() => {
  if (!chatService.currentUser) {
    chatService.initializeFromStorage()
  }
})
</script>

<style lang="scss" scoped>
.sidebar {
  width: 60px;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
  border-right: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40px;
}

.logo {
  color: white;
  background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
  font-weight: bold;
  font-size: 16px;
  padding: 10px 14px 14px;
  border-radius: 12px;
  margin-bottom: 8px;
  width: fit-content;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 88;
  box-shadow: 0 4px 12px rgba(49, 49, 49, 0.14);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
  }
}

.logo-text {
  text-align: center;
  font-size: 10px;
  color: #64748b;
}

.logo-title {
  font-weight: 600;
  margin-bottom: 2px;
}

.logo-subtitle {
  font-size: 8px;
}

.nav-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #64748b;
  margin-bottom: 2px;
  width: 40px;
  
  &:hover {
    background: #e2e8f0;
  }
  
  &.active {
    color: transparent;
    background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    
    .nav-icon {
      background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
      -webkit-mask: var(--icon-mask) no-repeat center;
      mask: var(--icon-mask) no-repeat center;
      -webkit-mask-size: contain;
      mask-size: contain;
    }
  }
}

.nav-icon {
  width: 20px;
  height: 20px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.home-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z'/%3E%3Cpolyline points='9,22 9,12 15,12 15,22'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z'/%3E%3Cpolyline points='9,22 9,12 15,12 15,22'/%3E%3C/svg%3E");
}

.messages-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z'/%3E%3C/svg%3E");
}

.settings-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='3'/%3E%3Cpath d='M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='3'/%3E%3Cpath d='M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z'/%3E%3C/svg%3E");
}


.user-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: auto;
  width: 40px;
  
  &:hover {
    background: #e2e8f0;
  }
}

.user-avatar {
  width: 37px;
  height: 37px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 6px;
  border: 2px solid #2632449c;
  opacity: 0.9;
  cursor: pointer;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  cursor: pointer;
    
  }
}

// 用户信息卡片弹窗样式
.user-card-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 1000;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  padding: 20px;
}

.user-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  padding: 20px;
  min-width: 280px;
  max-width: 320px;
  margin-top: 60px;
  margin-left: 20px;
  animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.user-card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.user-card-avatar {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #e2e8f0;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.user-card-info {
  flex: 1;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 4px 0;
  line-height: 1.2;
}

.user-email {
  font-size: 14px;
  color: #64748b;
  margin: 0;
  line-height: 1.2;
}

.user-card-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.switch-account-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  color: #475569;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #f1f5f9;
    border-color: #cbd5e1;
    color: #334155;
  }
  
  &:active {
    transform: translateY(1px);
  }
}

.switch-icon {
  width: 16px;
  height: 16px;
  color: #64748b;
}

</style>
