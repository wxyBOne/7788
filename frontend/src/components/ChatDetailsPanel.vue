<template>
  <div class="details-panel">
    <div class="details-header">
      <h3>聊天</h3>
      <button class="add-btn">+</button>
    </div>
    
    <div class="search-bar">
      <div class="search-icon"></div>
      <input type="text" placeholder="搜索..." />
    </div>
    
    <div class="chat-tabs">
      <div class="tab active">全部聊天</div>
    </div>
    
    
    <div class="chat-list">
      <ChatItem 
        :chatData="currentChat" 
        :isActive="true"
        @selectChat="$emit('selectChat', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ChatItem from './ChatItem.vue'

defineEmits(['selectChat'])

// 聊天数据
const currentChat = ref({
  id: 'hermione',
  name: '赫敏',
  avatar: '/src/img/Hermione.webp',
  lastMessage: 'Okay😊 I know very well about...',
  time: 'Now',
  unreadCount: 3,
  isOnline: true
})
</script>

<style lang="scss" scoped>
.details-panel {
  width: 300px;
  background: white;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.details-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  
  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
  }
}

.add-btn {
  width: 32px;
  height: 32px;
  background: #52b4b4da;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background: #4da6a6;
  }
}

.search-bar {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  position: relative;
  
  .search-icon {
    position: absolute;
    left: 32px;
    top: 50%;
    transform: translateY(-50%);
    width: 16px;
    height: 16px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='11' cy='11' r='8'/%3E%3Cpath d='M21 21l-4.35-4.35'/%3E%3C/svg%3E");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    z-index: 1;
    pointer-events: none;
  }
  
  input {
    flex: 1;
    border: none;
    outline: none;
    padding: 8px 12px 8px 36px;
    background: #eff2f4;
    border-radius: 8px;
    font-size: 14px;
    color: #64748b;
    transition: all 0.3s ease;
    box-sizing: border-box;
    
    &:focus {
      outline: none;
      border-color: #cacaca;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
      background: rgba(255, 255, 255, 1);
    }
    
    &::placeholder {
      color: #94a3b8;
    }
  }
}


.chat-tabs {
  display: flex;
  padding: 8px 20px 0;
}

.tab {
  padding: 8px 16px;
  font-size: 14px;
  color: #64748b;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  
  &.active {
  color: #50A7B0;
  border-bottom-color: #50A7B0;
  }
}


.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}



</style>
