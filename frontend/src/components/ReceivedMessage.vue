<template>
    <div class="message received">
      <div class="message-avatar">
        <img :src="message.avatar" :alt="message.name" />
      </div>
      <div class="message-content">
        <!-- 文本消息 -->
        <div v-if="message.type === 'text'" class="message-text">
          {{ message.content }}
        </div>
        
        <!-- 图片消息 -->
        <div v-else-if="message.type === 'image'" class="message-image">
          <img :src="message.imageUrl" :alt="message.content" />
          <div class="message-text">{{ message.content }}</div>
        </div>
        
        <!-- 音频消息 -->
        <div v-else-if="message.type === 'audio'" class="message-audio">
          <div class="audio-waveform">🎵</div>
          <div class="audio-duration">{{ message.duration }}</div>
        </div>
        
        <!-- 文件消息 -->
        <div v-else-if="message.type === 'file'" class="message-file">
          <div class="file-icon">📄</div>
          <div class="file-info">
            <div class="file-name">{{ message.fileName }}</div>
            <div class="file-size">{{ message.fileSize }}</div>
          </div>
        </div>
        
        <div class="message-time">{{ message.time }}</div>
      </div>
    </div>
  </template>
  
  <script setup>
  defineProps({
    message: {
      type: Object,
      required: true
    }
  })
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