<template>
  <div class="chat-item" :class="{ active: isActive }" @click="$emit('selectChat', chatData)">
    <div class="chat-avatar">
      <img :src="chatData.avatar_url || chatData.avatar" :alt="chatData.name" />
    </div>
    <div class="status-dot" v-if="chatData.is_online"></div>
    <div class="chat-content">
      <div class="chat-name">{{ chatData.name }}</div>
      <div class="chat-preview">{{ truncateMessage(chatData.last_message) }}</div>
    </div>
    <div class="chat-meta">
      <div class="chat-time">{{ formatTime(chatData.last_message_at) }}</div>
    </div>
  </div>
</template>

<script setup>
import chatService from '@/services/chatService.js'

defineProps({
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
  }
})

defineEmits(['selectChat'])

// 工具方法
const truncateMessage = (message) => {
  return chatService.truncateMessage(message, 30)
}

const formatTime = (timestamp) => {
  return chatService.formatTime(timestamp)
}
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

</style>
