<template>
  <div class="chat-area">
    <div v-if="selectedChat" class="chat-conversation">
      <div class="chat-header">
        <div class="chat-user-info" @click="$emit('toggleProfile')">
          <div class="chat-user-avatar">
            <img :src="selectedChat.avatar" :alt="selectedChat.name" />
          </div>
          <div class="chat-user-details">
            <div class="chat-user-name">{{ selectedChat.name }}</div>
            <div class="chat-user-status">在线</div>
          </div>
        </div>
        <div class="chat-actions">
          <button class="action-btn phone-btn"></button>
          <button class="action-btn more-btn"></button>
        </div>
      </div>
      
      <div class="messages-container">
        <template v-for="message in messages" :key="message.id">
          <ReceivedMessage 
            v-if="message.name !== 'You'" 
            :message="message" 
          />
          <SentMessage 
            v-else 
            :message="message" 
          />
        </template>
      </div>
      
      <div class="message-input">
         <input type="text" placeholder="输入消息..." />
        <div class="input-actions">
          <button class="input-btn attach-btn"></button>
          <button class="input-btn emoji-btn"></button>
          <button class="send-btn"></button>
        </div>
      </div>
    </div>
    
    <div v-else class="no-chat">
      <div class="no-chat-content">
        <div class="no-chat-icon">💬</div>
         <div class="no-chat-text">选择一个聊天开始对话</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ReceivedMessage from './ReceivedMessage.vue'
import SentMessage from './SentMessage.vue'

defineProps({
  selectedChat: {
    type: Object,
    default: null
  }
})

defineEmits(['toggleProfile'])

// 消息数据
const messages = ref([
  {
    id: 1,
    type: 'image',
    avatar: '/src/img/Hermione.webp',
    name: '赫敏',
    content: '我觉得这张图片很适合我们的设计。',
    time: '今天 14:45',
    imageUrl: 'https://via.placeholder.com/200x120/87CEEB/FFFFFF?text=Mountain'
  },
  {
    id: 2,
    type: 'audio',
    avatar: '/src/img/Hermione.webp',
    name: '赫敏',
    content: '',
    time: '今天 14:45',
    duration: '02:23'
  },
  {
    id: 3,
    type: 'text',
    avatar: '/src/img/Hermione.webp',
    name: '赫敏',
    content: '你好！很想看看一些设计。😊',
    time: '今天 14:45'
  },
  {
    id: 4,
    type: 'text',
    avatar: '/src/img/Hermione.webp',
    name: '赫敏',
    content: '这里有一些图片素材，你有时间的时候可以看看。',
    time: '15:00'
  },
  {
    id: 5,
    type: 'file',
    avatar: '/src/img/Hermione.webp',
    name: '赫敏',
    content: '',
    time: '15:00',
    fileName: '图片素材.zip',
    fileSize: '16 MB'
  },
  {
    id: 6,
    type: 'text',
    avatar: '/src/img/Hermione.webp',
    name: 'You',
    content: '谢谢！！！',
    time: '15:05'
  },
  {
    id: 7,
    type: 'audio',
    avatar: '/src/img/Hermione.webp',
    name: 'You',
    content: '',
    time: '14:50',
    duration: '02:23'
  },
  {
    id: 8,
    type: 'text',
    avatar: '/src/img/Hermione.webp',
    name: 'You',
    content: '太棒了 🎉 这是个很好的设计想法。🤩',
    time: '15:05'
  }
])
</script>

<style lang="scss" scoped>
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);

  min-width: 0;
}

.chat-conversation {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.chat-user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 10px 8px 8px;
  border-radius: 15px;
  transition: all 0.2s ease;
  
  &:hover {
    background: #f8fafc;
  }
  
  .chat-user-avatar {
    cursor: pointer;
  }
  
  .chat-user-details {
    cursor: pointer;
  }
}

.chat-user-avatar {
  width: 46px;
  height: 46px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 12px;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    cursor: pointer;

  }
}

.chat-user-details {
  .chat-user-name {
    font-size: 16px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 2px;
    cursor: pointer;
  }
  
  .chat-user-status {
    font-size: 12px;
    color: #10b981;
    cursor: pointer;
  }
}

.chat-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8fafc;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background-size: 18px;
  background-repeat: no-repeat;
  background-position: center;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: #e2e8f0;
    background-size: 16px; // hover时图标缩小
  }
}

.phone-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3C/svg%3E");
  }
}

.video-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolygon points='23 7 16 12 23 17 23 7'/%3E%3Crect x='1' y='5' width='15' height='14' rx='2' ry='2'/%3E%3C/svg%3E");
}

.more-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='1'/%3E%3Ccircle cx='19' cy='12' r='1'/%3E%3Ccircle cx='5' cy='12' r='1'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='1'/%3E%3Ccircle cx='19' cy='12' r='1'/%3E%3Ccircle cx='5' cy='12' r='1'/%3E%3C/svg%3E");
  }
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: flex-start;
  
  // 确保消息容器正确布局
  width: 100%;
  box-sizing: border-box;
  
  // 隐藏滚动条但保持滚动功能
  scrollbar-width: none;
  -ms-overflow-style: none;
  
  &::-webkit-scrollbar {
    display: none;
  }
}

// 为发送的消息添加右对齐
.message.sent {
  align-self: flex-end;
  align-items: flex-end;
  
  .message-content {
    align-self: flex-end;
  }
}

.message.received {
  align-self: flex-start;
  align-items: flex-start;
  
  .message-content {
    align-self: flex-start;
  }
}


.message-image {
  margin-bottom: 8px;
  
  img {
    width: 200px;
    height: 120px;
    object-fit: cover;
    border-radius: 8px;
  }
}

.message-audio {
  display: flex;
  align-items: center;
  gap: 8px;
}

.audio-waveform {
  font-size: 16px;
}

.audio-duration {
  font-size: 12px;
  color: #64748b;
}

.message-file {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-icon {
  font-size: 20px;
}

.file-info {
  .file-name {
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 2px;
  }
  
  .file-size {
    font-size: 12px;
    color: #64748b;
  }
}

.message-input {
  display: flex;
  align-items: center;
  padding: 20px;
  border-top: 1px solid #e2e8f0;
  gap: 12px;
  
  input {
    flex: 1;
    border: none;
    outline: none;
    padding: 12px 16px;
    background: #f8fafc;
    border-radius: 24px;
    font-size: 14px;
    
    &::placeholder {
      color: #94a3b8;
    }
  }
}

.input-actions {
  display: flex;
  gap: 8px;
}

.input-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8fafc;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background-size: 18px;
  background-repeat: no-repeat;
  background-position: center;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: #e2e8f0;
    background-size: 16px; // hover时图标缩小
  }
}

.attach-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66L9.64 16.2a2 2 0 0 1-2.83-2.83l8.49-8.49'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66L9.64 16.2a2 2 0 0 1-2.83-2.83l8.49-8.49'/%3E%3C/svg%3E");
  }
}

.emoji-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10'/%3E%3Cpath d='M8 14s1.5 2 4 2 4-2 4-2'/%3E%3Cline x1='9' y1='9' x2='9.01' y2='9'/%3E%3Cline x1='15' y1='9' x2='15.01' y2='9'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10'/%3E%3Cpath d='M8 14s1.5 2 4 2 4-2 4-2'/%3E%3Cline x1='9' y1='9' x2='9.01' y2='9'/%3E%3Cline x1='15' y1='9' x2='15.01' y2='9'/%3E%3C/svg%3E");
  }
}

.send-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #52b4b4da;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: all 0.2s ease;
  
  &::after {
    content: '';
    position: absolute;
    width: 18px;
    height: 18px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cline x1='22' y1='2' x2='11' y2='13'/%3E%3Cpolygon points='22,2 15,22 11,13 2,9 22,2'/%3E%3C/svg%3E");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
  }
  
  &:hover {
    background: #4da6a6;
    
    &::after {
      width: 16px;
      height: 16px;
    }
  }
}

.no-chat {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.no-chat-content {
  text-align: center;
  color: #64748b;
}

.no-chat-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-chat-text {
  font-size: 16px;
}
</style>
