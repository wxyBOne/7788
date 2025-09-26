<template>
  <div class="profile-panel">
    <div class="profile-header">
      <div class="profile-avatar">
        <img :src="selectedChat?.avatar_url || selectedChat?.avatar" :alt="selectedChat?.name" />
      </div>
      <div class="profile-name">{{ selectedChat?.name }}</div>
      <div class="profile-description">{{ selectedChat?.personality_signature || '暂无个性签名' }}</div>
    </div>
    
    
    <div class="profile-actions">
      <button class="block-btn">屏蔽 {{ selectedChat?.name }}</button>
    </div>
  </div>
</template>

<script setup>
import chatService from '@/services/chatService.js'

defineProps({
  selectedChat: {
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
.profile-panel {
  width: 280px; // 加宽到280px
  background: #f8fafc;
  border-left: 1px solid #e2e8f0;
  padding: 24px;
  overflow-y: auto;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  
  // 隐藏滚动条但保持滚动功能
  scrollbar-width: none; // Firefox
  -ms-overflow-style: none; // IE/Edge
  
  &::-webkit-scrollbar {
    display: none; // Chrome/Safari
  }
}

.profile-header {
  text-align: center;
  margin-top: 20px;
  margin-bottom: 24px;
}

.profile-avatar {
  width: 75px;
  height: 75px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 16px;
  border: 4px solid #2632449c;

  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.profile-name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
}

.profile-description {
  font-size: 14px;
  color: #64748b;
  font-style: italic;
  margin-bottom: 8px;
}

.profile-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  
  span {
    font-size: 14px;
    font-weight: 600;
    color: #1e293b;
  }
}


.section-arrow {
  color: #94a3b8;
  cursor: pointer;
}



.profile-actions {
  margin-top: auto;
  padding-top: 20px;
}

.block-btn {
  width: 100%;
  padding: 12px 16px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  
  &:hover {
    background: #dc2626;
  }
  
  &::before {
    content: "×";
    font-size: 16px;
  }
}
</style>
