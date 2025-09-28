<template>
  <div class="chat-item" :class="{ active: isActive }" @click="handleChatClick">
    <div class="chat-avatar">
      <!-- AI伙伴使用粒子小球头像 -->
      <ParticleAvatar
        v-if="isCompanion"
        :emotion="companionEmotion.emotion"
        :intensity="companionEmotion.intensity"
        :color="companionEmotion.color"
        :brightness="companionEmotion.brightness"
        :particle-speed="companionEmotion.particle_speed"
        :growth-percentage="chatData.growth_percentage || 0"
        size="small"
      />
      <!-- 普通角色使用普通头像 -->
      <img v-else :src="chatData.avatar_url || chatData.avatar" :alt="chatData.name" />
    </div>
    <div class="status-dot" v-if="chatData.is_online"></div>
    <div class="chat-content">
      <div class="chat-name">{{ chatData.name }}</div>
      <div class="chat-preview">{{ truncateMessage(chatData.last_message) }}</div>
    </div>
    <div class="chat-meta">
      <div class="chat-time">{{ formatTime(chatData.last_message_at) }}</div>
      <!-- AI伙伴成长进度指示器 -->
      <div v-if="isCompanion && chatData.growth_percentage < 100" class="growth-indicator">
        <div class="growth-bar">
          <div 
            class="growth-fill" 
            :style="{ width: chatData.growth_percentage + '%' }"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ParticleAvatar from './ParticleAvatar.vue'
import chatService from '@/services/chatService.js'
import api from '@/services/api.js'

const props = defineProps({
  chatData: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      name: '',
      avatar_url: '',
      last_message: '',
      last_message_at: null,
      is_online: false
    })
  },
  isActive: {
    type: Boolean,
    default: false
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

const emit = defineEmits(['selectChat', 'showCompanionInit'])

// 判断是否为AI伙伴
const isCompanion = computed(() => {
  return props.chatData.name === '空白AI' || props.chatData.type === 'companion'
})

// AI伙伴情绪状态
const companionEmotion = ref({
  emotion: '平静',
  intensity: 0.5,
  color: '#52b4b4',
  brightness: 0.7,
  particle_speed: 0.5
})

// 加载AI伙伴情绪状态
const loadCompanionEmotion = async () => {
  if (!isCompanion.value) return
  
  try {
    const token = localStorage.getItem('token')
    const response = await api.companion.getEmotionState(token, props.chatData.id)
    if (response.success) {
      companionEmotion.value = response.emotion
    }
  } catch (error) {
    console.error('加载AI伙伴情绪状态失败:', error)
  }
}

// 工具方法
const truncateMessage = (message) => {
  return chatService.truncateMessage(message, 30)
}

const formatTime = (timestamp) => {
  return chatService.formatTime(timestamp)
}

// 处理聊天点击
const handleChatClick = () => {
  // 如果是空白AI，检查用户是否已有AI伙伴
  if (props.chatData.name === '空白AI') {
    // 检查是否为AI伙伴（有growth_percentage字段说明是已创建的AI伙伴）
    if (props.chatData.type === 'companion' || props.chatData.growth_percentage !== undefined) {
      // 用户已有AI伙伴，正常选择聊天
      emit('selectChat', props.chatData)
    } else {
      // 用户还没有AI伙伴，显示初始化弹窗
      emit('showCompanionInit')
    }
  } else {
    // 普通角色，正常选择聊天
    emit('selectChat', props.chatData)
  }
}

onMounted(() => {
  if (isCompanion.value) {
    loadCompanionEmotion()
  }
})
</script>

<style lang="scss" scoped>
.chat-item {
  display: flex;
  align-items: center;
  padding: 14px 18px 14px 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 8px;
  margin: 0 12px;
  position: relative;
  
  &:hover {
    background: #f8fafc;
  }
  
  &.active {
    background: transparent;
    color: #1e293b;
    
    .chat-name, .chat-preview {
      color: #1e293b;
    }
    
    .chat-time {
      color: #64748b;
    }
  }
}

.chat-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 12px;
  cursor: pointer;
  position: relative;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.status-dot {
  position: absolute;
  top: 15px;
  left: 46px;
  width: 12px;
  height: 12px;
  background: #10b981;
  border-radius: 50%;
  border: 2px solid white;
  z-index: 10;
}

.chat-content {
  flex: 1;
  min-width: 0;
}

.chat-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  cursor: pointer;
  margin-bottom: 2px;
}

.chat-preview {
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  cursor: pointer;
  
  &.typing {
    color: #000000;
    font-style: italic;
  }
}

.chat-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: flex-start;
  margin-left: 8px;
  min-height: 32px;
  position: relative;
  cursor: pointer;
}

.chat-time {
  font-size: 11px;
  color: #94a3b8;
  margin-bottom: auto;
  cursor: pointer;
}

.growth-indicator {
  margin-top: 4px;
  width: 40px;
}

.growth-bar {
  width: 100%;
  height: 3px;
  background: #e2e8f0;
  border-radius: 2px;
  overflow: hidden;
}

.growth-fill {
  height: 100%;
  background: linear-gradient(90deg, #52b4b4, #059669);
  border-radius: 2px;
  transition: width 0.3s ease;
}
</style>
