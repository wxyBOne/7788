<template>
  <div class="details-panel">
    <div class="details-header">
      <h3>{{ isAddFriendMode ? '添加好友' : '消息' }}</h3>
      <div class="header-actions">
        <button v-if="isAddFriendMode" class="back-btn" @click="toggleAddFriendMode">←</button>
        <button v-if="isAddFriendMode" class="delete-btn" @click="toggleAddFriendMode">×</button>
        <button v-else class="add-btn" @click="toggleAddFriendMode">+</button>
      </div>
    </div>
    
    <div class="search-bar">
      <div class="search-icon"></div>
      <input 
        v-model="searchKeyword"
        type="text" 
        :placeholder="isAddFriendMode ? '搜索新好友...' : '搜索好友...'"
        @keyup.enter="handleSearch"
      />
      <button v-if="searchKeyword" class="cancel-btn" @click="clearSearch">取消</button>
    </div>
    
    <div class="chat-tabs">
      <div class="tab active">{{ isAddFriendMode ? '新好友' : '我的好友' }}</div>
    </div>
    
    <div class="chat-list">
      <!-- 好友列表 -->
      <div v-if="!isAddFriendMode">
        <ChatItem 
          v-for="friend in filteredFriends"
          :key="friend.id"
          :chatData="friend" 
          :isActive="selectedChatId === friend.id"
          @selectChat="$emit('selectChat', friend)"
        />
      </div>
      
      <!-- 可添加的角色列表 -->
      <div v-else>
        <div 
          v-for="character in searchResults"
          :key="character.id"
          class="character-item"
          :class="{ 'added': character.is_added }"
        >
          <div class="character-info">
            <img :src="character.avatar_url" :alt="character.name" class="character-avatar" />
            <div class="character-details">
              <div class="character-name">{{ character.name }}</div>
              <div class="character-description">{{ character.description }}</div>
            </div>
          </div>
          <button 
            v-if="!character.is_added"
            class="add-character-btn"
            @click="addFriend(character.id)"
            :disabled="isAdding"
          >
            {{ isAdding ? '添加中...' : '添加' }}
          </button>
          <span v-else class="added-text">已添加</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import ChatItem from './ChatItem.vue'
import chatService from '@/services/chatService.js'

const props = defineProps({
  selectedChatId: {
    type: Number,
    default: null
  }
})

defineEmits(['selectChat'])

// 响应式数据
const isAddFriendMode = ref(false)
const searchKeyword = ref('')
const searchResults = ref([])
const isAdding = ref(false)

// 计算属性
const filteredFriends = computed(() => {
  if (!searchKeyword.value) {
    return chatService.friends
  }
  return chatService.friends.filter(friend => 
    friend.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

// 监听搜索关键词变化
watch(searchKeyword, (newValue) => {
  if (isAddFriendMode.value && newValue) {
    handleSearch()
  }
})

// 方法
const toggleAddFriendMode = () => {
  isAddFriendMode.value = !isAddFriendMode.value
  if (!isAddFriendMode.value) {
    searchResults.value = []
    searchKeyword.value = ''
  }
}

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) return
  
  try {
    await chatService.searchFriends(searchKeyword.value)
    searchResults.value = chatService.searchResults
  } catch (error) {
    console.error('搜索失败:', error)
  }
}

const clearSearch = () => {
  searchKeyword.value = ''
  searchResults.value = []
}

const addFriend = async (characterId) => {
  isAdding.value = true
  try {
    await chatService.addFriend(characterId)
    // 更新搜索结果中的状态
    const character = searchResults.value.find(c => c.id === characterId)
    if (character) {
      character.is_added = true
    }
  } catch (error) {
    console.error('添加好友失败:', error)
    alert('添加好友失败：' + error.message)
  } finally {
    isAdding.value = false
  }
}
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

.header-actions {
  display: flex;
  gap: 8px;
}

.add-btn, .back-btn, .delete-btn {
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
  transition: all 0.2s;

  &:hover {
    background: #4da6a6;
  }
}

.back-btn {
  background: #6b7280;
  
  &:hover {
    background: #4b5563;
  }
}

.delete-btn {
  background: #ef4444;
  
  &:hover {
    background: #dc2626;
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
  
  .cancel-btn {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: #6b7280;
    font-size: 12px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: all 0.2s;
    
    &:hover {
      background: #f3f4f6;
      color: #374151;
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

// 角色列表样式
.character-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid #f1f5f9;
  transition: all 0.2s;
  
  &:hover {
    background: #f8fafc;
  }
  
  &.added {
    opacity: 0.6;
  }
}

.character-info {
  display: flex;
  align-items: center;
  flex: 1;
}

.character-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
}

.character-details {
  flex: 1;
}

.character-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 2px;
}

.character-description {
  font-size: 12px;
  color: #64748b;
  line-height: 1.4;
}

.add-character-btn {
  padding: 6px 12px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover:not(:disabled) {
    background: #2563eb;
  }
  
  &:disabled {
    background: #9ca3af;
    cursor: not-allowed;
  }
}

.added-text {
  font-size: 12px;
  color: #10b981;
  font-weight: 500;
}

</style>
