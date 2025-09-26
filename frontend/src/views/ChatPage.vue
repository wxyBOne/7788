<template>
    <div class="chat-page">
      <ChatSidebar 
        :activeSection="activeSection" 
        @setActiveSection="setActiveSection" 
      />
  
      <!-- 使用 v-show 替代 v-if 保持 DOM 存在 -->
      <Transition name="slide-fade">
        <ChatDetailsPanel 
          v-show="showDetails" 
          @selectChat="selectChat" 
          class="details-panel"
        />
      </Transition>
  
      <ChatArea 
        :selectedChat="selectedChat" 
        @toggleProfile="toggleProfile" 
        class="chat-area"
        :class="{ 
          'with-details': showDetails, 
          'with-profile': showProfile 
        }"
      />
  
      <Transition name="slide-fade">
        <ChatProfilePanel 
          v-show="showProfile" 
          :selectedChat="selectedChat" 
          class="profile-panel"
        />
      </Transition>
    </div>
  </template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import ChatSidebar from '@/components/ChatSidebar.vue'
import ChatDetailsPanel from '@/components/ChatDetailsPanel.vue'
import ChatArea from '@/components/ChatArea.vue'
import ChatProfilePanel from '@/components/ChatProfilePanel.vue'
import chatService from '@/services/chatService.js'

// 响应式数据
const activeSection = ref('messages')
const showDetails = ref(false)
const showProfile = ref(false)
const selectedChat = ref(null)
const isLoading = ref(true)

// 初始化数据
onMounted(async () => {
  try {
    // 检查用户是否已登录
    const token = localStorage.getItem('token')
    if (!token) {
      // 未登录，跳转到登录页
      window.location.href = '/login'
      return
    }

    // 只有当好友列表为空时才加载，避免重复加载
    if (chatService.friends.length === 0) {
      await chatService.loadUserFriends()
    }
    
    // 如果有好友，并且当前没有选中聊天，则默认选择第一个
    if (chatService.friends.length > 0 && !selectedChat.value) {
      selectedChat.value = chatService.friends[0]
      await chatService.switchChat(selectedChat.value)
    }
    
    isLoading.value = false
  } catch (error) {
    console.error('初始化失败:', error)
    isLoading.value = false
  }
})

// 方法
const setActiveSection = (section) => {
  if (activeSection.value === section && showDetails.value) {
    // 如果点击的是当前激活的功能，则收起详情栏
    showDetails.value = false
    showProfile.value = false
  } else {
    activeSection.value = section
    if (section === 'messages') {
      showDetails.value = true
    } else {
      showDetails.value = false
      showProfile.value = false
    }
  }
}

const selectChat = async (chat) => {
  selectedChat.value = chat
  await chatService.switchChat(chat)
  showDetails.value = false
}

const toggleProfile = () => {
  if (selectedChat.value) {
    showProfile.value = !showProfile.value
  }
}

</script>

<style lang="scss" scoped>
.chat-page {
  display: flex;
  height: 100vh;
  background: #f8fafc;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  overflow: hidden;
  position: relative;
}

// 面板基础样式
.details-panel,
.profile-panel {
  position: absolute;
  top: 0;
  height: 100%;
  z-index: 10;
}

.details-panel {
  left: 60px; // 侧边栏宽度
  border-right: 1px solid #e2e8f0; // 添加1px间距
}

.profile-panel {
  right: 0;
  border-left: 1px solid #e2e8f0; // 添加1px间距
}

.chat-area {
  flex: 1;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); // 延长过渡时间
  min-width: 0;
  z-index: 5;
  
  &.with-details {
    margin-left: 300px; // 侧边栏160px + 面板宽度300px + 1px间距
  }
  
  &.with-profile {
    margin-right: 281px; // 面板宽度280px + 1px间距
  }
}

// 优化过渡动画
.slide-fade-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); // 离开时稍快
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
</style>