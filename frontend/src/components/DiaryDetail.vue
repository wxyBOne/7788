<template>
  <div class="diary-overlay" @click="closeDiary">
    <div class="diary-container" @click.stop>
      <!-- 日记头部 -->
      <div class="diary-header">
        <div class="diary-title">
          <h3>{{ diary.title }}</h3>
          <span class="diary-date">{{ formatDate(diary.date) }}</span>
        </div>
        <button class="close-btn" @click="closeDiary">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
      
      <!-- 日记内容 -->
      <div class="diary-content">
        <div class="diary-meta">
          <div class="mood-indicator">
            <span class="mood-label">心情:</span>
            <div class="mood-stars">
              <span 
                v-for="i in 5" 
                :key="i" 
                class="star"
                :class="{ filled: i <= diary.mood_score }"
              >★</span>
            </div>
          </div>
          <div class="writing-style">
            <span class="style-label">文笔风格:</span>
            <span class="style-value">{{ getWritingStyle(diary.content) }}</span>
          </div>
        </div>
        
        <div class="diary-text">
          {{ diary.content }}
        </div>
      </div>
      
      <!-- 日记底部 -->
      <div class="diary-footer">
        <div class="diary-stats">
          <span class="word-count">{{ diary.content.length }} 字</span>
          <span class="mention-indicator" v-if="diary.is_user_mentioned">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            提到了你
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  diary: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

// 关闭日记
const closeDiary = () => {
  emit('close')
}

// 格式化日期
const formatDate = (date) => {
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

// 根据内容判断文笔风格
const getWritingStyle = (content) => {
  const length = content.length
  const hasComplexSentences = content.includes('，') && content.includes('。')
  const hasEmotionalWords = /[开心|快乐|难过|伤心|兴奋|平静|期待|想念]/.test(content)
  
  if (length < 50) {
    return '稚嫩'
  } else if (length < 100 && !hasComplexSentences) {
    return '简单'
  } else if (length < 200 && hasComplexSentences) {
    return '流畅'
  } else if (hasEmotionalWords && hasComplexSentences) {
    return '成熟'
  } else {
    return '细腻'
  }
}
</script>

<style lang="scss" scoped>
.diary-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.diary-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
  max-width: 600px;
  width: 100%;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.diary-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px 24px 16px;
  border-bottom: 1px solid #e2e8f0;
}

.diary-title h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.3;
}

.diary-date {
  font-size: 14px;
  color: #64748b;
  font-weight: 400;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #f8fafc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #64748b;
  
  &:hover {
    background: #e2e8f0;
    color: #475569;
  }
  
  svg {
    width: 16px;
    height: 16px;
  }
}

.diary-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.diary-meta {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.mood-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mood-label {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.mood-stars {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 16px;
  color: #e2e8f0;
  transition: color 0.2s ease;
  
  &.filled {
    color: #fbbf24;
  }
}

.writing-style {
  display: flex;
  align-items: center;
  gap: 8px;
}

.style-label {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.style-value {
  font-size: 14px;
  color: #1e293b;
  font-weight: 600;
  padding: 4px 8px;
  background: #f1f5f9;
  border-radius: 6px;
}

.diary-text {
  font-size: 16px;
  line-height: 1.7;
  color: #374151;
  white-space: pre-wrap;
  word-break: break-word;
}

.diary-footer {
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.diary-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #64748b;
}

.word-count {
  font-weight: 500;
}

.mention-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #059669;
  font-weight: 500;
  
  svg {
    width: 14px;
    height: 14px;
  }
}

// 响应式设计
@media (max-width: 640px) {
  .diary-overlay {
    padding: 10px;
  }
  
  .diary-container {
    max-height: 90vh;
  }
  
  .diary-header {
    padding: 20px 20px 12px;
  }
  
  .diary-content {
    padding: 20px;
  }
  
  .diary-footer {
    padding: 12px 20px;
  }
  
  .diary-meta {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
