<template>
  <div class="message received">
    <div class="message-avatar">
      <img :src="character?.avatar_url || message.avatar_url || message.avatar" :alt="character?.name || message.name" />
    </div>
    <div class="message-content">
      <!-- 文本消息 -->
      <div v-if="message.message_type === 'text'" class="message-text">
        {{ message.ai_response || message.content }}
      </div>
      
      <!-- 表情消息 -->
      <div v-else-if="message.message_type === 'emoji'" class="message-text emoji-message">
        {{ message.ai_response || message.content }}
      </div>
      
      <!-- 图片消息 -->
      <div v-else-if="message.message_type === 'image'" class="message-image">
        <img :src="message.image_url || message.imageUrl" :alt="message.content" />
        <div class="message-text">{{ message.ai_response || message.content }}</div>
      </div>
      
      <!-- 语音消息 - 显示文字内容而不是音频播放按钮 -->
      <div v-else-if="message.message_type === 'voice'" class="message-text">
        {{ message.ai_response || message.content }}
      </div>
      
      <div class="message-time">{{ formatTime(message.created_at) }}</div>
    </div>
  </div>
</template>
  
<script setup>
import chatService from '@/services/chatService.js'

defineProps({
  message: {
    type: Object,
    required: true
  },
  character: {
    type: Object,
    default: null
  }
})

// 格式化时间
const formatTime = (timestamp) => {
  return chatService.formatTime(timestamp)
}
  </script>
  
  <style lang="scss" scoped>
  .message {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 16px;
    width: 100%;
    max-width: 100%;
  }
  
.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  // 关键：让头像顶部与气泡顶部对齐，通过调整 margin 实现
  margin-top: 2px; 
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
  
  .message-content {
    max-width: 70%; 
    background: #f2f5f7;
    border-radius: 16px;
    padding: 12px 16px;
    word-break: break-word;
    word-wrap: break-word;
    overflow-wrap: break-word;
    position: relative; // 为时间定位做准备
  }
  
  .message-text {
    font-size: 14px;
    line-height: 1.4;
    margin-bottom: 4px;
    color: #1e293b;
    white-space: pre-wrap;
  }
  
  .message-image {
    margin-bottom: 8px;
    
    img {
      max-width: 200px;
      max-height: 120px;
      border-radius: 8px;
      margin-bottom: 8px;
    }
  }
  
  .message-audio {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
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
    gap: 8px;
    margin-bottom: 4px;
  }
  
  .file-icon {
    font-size: 16px;
  }
  
  .file-info {
    display: flex;
    flex-direction: column;
  }
  
  .file-name {
    font-size: 14px;
    color: #1e293b;
    font-weight: 500;
    word-break: break-word;
  }
  
  .file-size {
    font-size: 12px;
    color: #64748b;
  }
  
  .message-time {
    font-size: 11px;
    color: #94a3b8;
    margin-top: 4px;
    text-align: right;
  }
  
  @media (max-width: 768px) {
    .message-content {
      max-width: 85%;
    }
  }
  
  @media (max-width: 480px) {
    .message-content {
      max-width: 90%;
    }
  }
  </style>