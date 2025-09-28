<template>
  <div class="companion-profile-panel">
    <!-- 粒子小球头像 -->
    <div class="companion-avatar-section">
      <ParticleAvatar
        :emotion="emotionState.emotion"
        :intensity="emotionState.intensity"
        :color="emotionState.color"
        :brightness="emotionState.brightness"
        :particle-speed="emotionState.particle_speed"
        :growth-percentage="companion.growth_percentage"
        size="large"
      />
    </div>
    
    <!-- 基本信息 -->
    <div class="companion-info">
      <h2 class="companion-name">{{ companion.name }}</h2>
      <p class="companion-signature">{{ companion.personality_signature }}</p>
    </div>
    
    <!-- 成长进度 -->
    <div class="growth-section">
      <div class="growth-header">
        <h3>成长进度</h3>
        <span class="growth-percentage">{{ Math.round(companion.growth_percentage) }}%</span>
      </div>
      <div class="growth-progress">
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            :style="{ width: companion.growth_percentage + '%', backgroundColor: emotionState.color }"
          ></div>
        </div>
        <div class="growth-details">
          <div class="growth-item">
            <span class="label">等级:</span>
            <span class="value">{{ companion.current_level }}</span>
          </div>
          <div class="growth-item">
            <span class="label">经验:</span>
            <span class="value">{{ companion.total_experience }}</span>
          </div>
          <div class="growth-item">
            <span class="label">模式:</span>
            <span class="value">{{ companion.growth_mode === 'short' ? '短周期' : '长周期' }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 能力值 -->
    <div class="abilities-section">
      <h3>能力值</h3>
      <div class="abilities-grid">
        <div class="ability-item">
          <span class="ability-name">语言流畅度</span>
          <div class="ability-bar">
            <div 
              class="ability-fill" 
              :style="{ width: (displayAbilities.conversation_fluency / 10) * 100 + '%' }"
            ></div>
          </div>
          <span class="ability-value">{{ displayAbilities.conversation_fluency }}/10</span>
        </div>
        <div class="ability-item">
          <span class="ability-name">知识广度</span>
          <div class="ability-bar">
            <div 
              class="ability-fill" 
              :style="{ width: (displayAbilities.knowledge_breadth / 10) * 100 + '%' }"
            ></div>
          </div>
          <span class="ability-value">{{ displayAbilities.knowledge_breadth }}/10</span>
        </div>
        <div class="ability-item">
          <span class="ability-name">共情深度</span>
          <div class="ability-bar">
            <div 
              class="ability-fill" 
              :style="{ width: (displayAbilities.empathy_depth / 10) * 100 + '%' }"
            ></div>
          </div>
          <span class="ability-value">{{ displayAbilities.empathy_depth }}/10</span>
        </div>
        <div class="ability-item">
          <span class="ability-name">创造力</span>
          <div class="ability-bar">
            <div 
              class="ability-fill" 
              :style="{ width: (displayAbilities.creativity_level / 10) * 100 + '%' }"
            ></div>
          </div>
          <span class="ability-value">{{ displayAbilities.creativity_level }}/10</span>
        </div>
        <div class="ability-item">
          <span class="ability-name">幽默感</span>
          <div class="ability-bar">
            <div 
              class="ability-fill" 
              :style="{ width: (displayAbilities.humor_sense / 10) * 100 + '%' }"
            ></div>
          </div>
          <span class="ability-value">{{ displayAbilities.humor_sense }}/10</span>
        </div>
      </div>
    </div>
    
    <!-- 日记列表 -->
    <div class="diary-section">
      <div class="diary-header">
        <h3>成长日记</h3>
        <span class="diary-count">{{ diaries.length }} 篇</span>
      </div>
      <div class="diary-list">
        <div 
          v-for="diary in diaries" 
          :key="diary.id"
          class="diary-item"
          @click="openDiary(diary)"
        >
          <div class="diary-item-header">
            <span class="diary-title">{{ diary.title || '无标题' }}</span>
            <span class="diary-date">{{ formatDate(diary.date) }}</span>
          </div>
          <div class="diary-preview">
            {{ diary.content.substring(0, 50) }}{{ diary.content.length > 50 ? '...' : '' }}
          </div>
          <div class="diary-meta">
            <div class="mood-stars">
              <span 
                v-for="i in 5" 
                :key="i" 
                class="star"
                :class="{ filled: i <= diary.mood_score }"
              >★</span>
            </div>
            <span class="word-count">{{ diary.content.length }}字</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 日记详情弹窗 -->
    <DiaryDetail
      v-if="selectedDiary"
      :diary="selectedDiary"
      @close="closeDiary"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ParticleAvatar from './ParticleAvatar.vue'
import DiaryDetail from './DiaryDetail.vue'
import api from '@/services/api.js'

const props = defineProps({
  companion: {
    type: Object,
    required: true
  }
})

// 状态
const diaries = ref([])
const selectedDiary = ref(null)
const emotionState = ref({
  emotion: '平静',
  intensity: 0.5,
  color: '#52b4b4',
  brightness: 0.7,
  particle_speed: 0.5
})

// 计算属性：处理AI伙伴的能力值显示
const displayAbilities = computed(() => {
  // 直接显示数据库中的实际能力值
  return {
    conversation_fluency: props.companion.conversation_fluency || 1,
    knowledge_breadth: props.companion.knowledge_breadth || 1,
    empathy_depth: props.companion.empathy_depth || 1,
    creativity_level: props.companion.creativity_level || 1,
    humor_sense: props.companion.humor_sense || 1
  }
})

// 格式化日期
const formatDate = (date) => {
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  })
}

// 打开日记详情
const openDiary = (diary) => {
  selectedDiary.value = diary
}

// 关闭日记详情
const closeDiary = () => {
  selectedDiary.value = null
}

// 加载日记列表
const loadDiaries = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await api.companion.getDiary(token, props.companion.id, 10)
    if (response.success) {
      diaries.value = response.diaries || []
    }
  } catch (error) {
    console.error('加载日记失败:', error)
  }
}

// 加载情绪状态
const loadEmotionState = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await api.companion.getEmotionState(token, props.companion.id)
    if (response.success) {
      emotionState.value = response.emotion
    }
  } catch (error) {
    console.error('加载情绪状态失败:', error)
  }
}

onMounted(() => {
  loadDiaries()
  loadEmotionState()
})
</script>

<style lang="scss" scoped>
.companion-profile-panel {
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 80vh;
  overflow-y: auto;
}

.companion-avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.companion-info {
  text-align: center;
  margin-bottom: 32px;
}

.companion-name {
  font-size: 24px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 8px 0;
}

.companion-signature {
  font-size: 16px;
  color: #64748b;
  font-style: italic;
  margin: 0;
  line-height: 1.5;
}

.growth-section {
  margin-bottom: 32px;
}

.growth-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.growth-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.growth-percentage {
  font-size: 20px;
  font-weight: 700;
  color: #059669;
}

.growth-progress {
  background: #f8fafc;
  border-radius: 12px;
  padding: 20px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 16px;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

.growth-details {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.growth-item {
  text-align: center;
}

.growth-item .label {
  display: block;
  font-size: 12px;
  color: #64748b;
  margin-bottom: 4px;
}

.growth-item .value {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.abilities-section {
  margin-bottom: 32px;
}

.abilities-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 16px 0;
}

.abilities-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ability-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ability-name {
  width: 80px;
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.ability-bar {
  flex: 1;
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  overflow: hidden;
}

.ability-fill {
  height: 100%;
  background: linear-gradient(90deg, #52b4b4, #059669);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.ability-value {
  width: 40px;
  font-size: 12px;
  color: #64748b;
  font-weight: 600;
  text-align: right;
}

.diary-section {
  margin-bottom: 16px;
}

.diary-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.diary-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.diary-count {
  font-size: 14px;
  color: #64748b;
  background: #f1f5f9;
  padding: 4px 8px;
  border-radius: 6px;
}

.diary-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.diary-item {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #f1f5f9;
    border-color: #cbd5e1;
    transform: translateY(-1px);
  }
}

.diary-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.diary-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.diary-date {
  font-size: 12px;
  color: #64748b;
}

.diary-preview {
  font-size: 14px;
  color: #64748b;
  line-height: 1.5;
  margin-bottom: 12px;
}

.diary-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mood-stars {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 12px;
  color: #e2e8f0;
  
  &.filled {
    color: #fbbf24;
  }
}

.word-count {
  font-size: 12px;
  color: #64748b;
}

// 响应式设计
@media (max-width: 640px) {
  .companion-profile-panel {
    padding: 16px;
  }
  
  .growth-details {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .ability-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .ability-name {
    width: auto;
  }
  
  .ability-bar {
    width: 100%;
  }
}
</style>
