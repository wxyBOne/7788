<template>
  <div class="profile-panel">
    <div class="profile-header">
      <div class="profile-avatar">
        <!-- AI伙伴使用粒子小球头像 -->
        <ParticleAvatar
          v-if="isCompanion"
          :emotion="companionEmotion.emotion"
          :intensity="companionEmotion.intensity"
          :color="companionEmotion.color"
          :brightness="companionEmotion.brightness"
          :particle-speed="companionEmotion.particle_speed"
          size="large"
        />
        <!-- 普通角色使用普通头像 -->
        <img v-else :src="selectedChat?.avatar_url || selectedChat?.avatar" :alt="selectedChat?.name" />
      </div>
      <div class="profile-name">{{ selectedChat?.name }}</div>
      
      <!-- 技能标签 -->
      <div class="skills-container" v-if="skills && skills.length > 0">
        <div class="skill-tag" v-for="skill in skills" :key="skill">
          {{ skill }}
        </div>
      </div>
      
      <div class="profile-description">{{ selectedChat?.personality_signature || '暂无个性签名' }}</div>
    </div>
    
    <!-- AI伙伴的成长进度 -->
    <div v-if="isCompanion" class="profile-section">
      <div class="growth-section">
        <div class="growth-header">
          <h3>成长进度</h3>
          <span class="growth-percentage">{{ Math.round(selectedChat.growth_percentage || 0) }}%</span>
        </div>
        <div class="growth-progress">
          <div class="progress-bar">
            <div 
              class="progress-fill" 
              :style="{ width: (selectedChat.growth_percentage || 0) + '%', backgroundColor: companionEmotion.color }"
            ></div>
          </div>
          <div class="growth-details">
            <div class="growth-item">
              <span class="label">等级:</span>
              <span class="value">{{ selectedChat.current_level || 1 }}</span>
            </div>
            <div class="growth-item">
              <span class="label">经验:</span>
              <span class="value">{{ selectedChat.total_experience || 0 }}</span>
            </div>
            <div class="growth-item">
              <span class="label">模式:</span>
              <span class="value">{{ selectedChat.growth_mode === 'short' ? '短周期' : '长周期' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- AI伙伴的能力值 -->
    <div v-if="isCompanion" class="profile-section">
      <div class="abilities-section">
        <h3>能力值</h3>
        <div class="abilities-grid">
          <div class="ability-item">
            <span class="ability-name">语言流畅度</span>
            <div class="ability-bar">
              <div 
                class="ability-fill" 
                :style="{ width: ((selectedChat.conversation_fluency || 1) / 10) * 100 + '%' }"
              ></div>
            </div>
            <span class="ability-value">{{ selectedChat.conversation_fluency || 1 }}/10</span>
          </div>
          <div class="ability-item">
            <span class="ability-name">知识广度</span>
            <div class="ability-bar">
              <div 
                class="ability-fill" 
                :style="{ width: ((selectedChat.knowledge_breadth || 1) / 10) * 100 + '%' }"
              ></div>
            </div>
            <span class="ability-value">{{ selectedChat.knowledge_breadth || 1 }}/10</span>
          </div>
          <div class="ability-item">
            <span class="ability-name">共情深度</span>
            <div class="ability-bar">
              <div 
                class="ability-fill" 
                :style="{ width: ((selectedChat.empathy_depth || 1) / 10) * 100 + '%' }"
              ></div>
            </div>
            <span class="ability-value">{{ selectedChat.empathy_depth || 1 }}/10</span>
          </div>
          <div class="ability-item">
            <span class="ability-name">创造力</span>
            <div class="ability-bar">
              <div 
                class="ability-fill" 
                :style="{ width: ((selectedChat.creativity_level || 1) / 10) * 100 + '%' }"
              ></div>
            </div>
            <span class="ability-value">{{ selectedChat.creativity_level || 1 }}/10</span>
          </div>
          <div class="ability-item">
            <span class="ability-name">幽默感</span>
            <div class="ability-bar">
              <div 
                class="ability-fill" 
                :style="{ width: ((selectedChat.humor_sense || 1) / 10) * 100 + '%' }"
              ></div>
            </div>
            <span class="ability-value">{{ selectedChat.humor_sense || 1 }}/10</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 角色简介 -->
    <div class="profile-section" v-if="selectedChat?.background_story">
      <div class="background-story">简介:
        {{ selectedChat.background_story }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import ParticleAvatar from './ParticleAvatar.vue'
import api from '@/services/api.js'

const props = defineProps({
  selectedChat: {
    type: Object,
    default: null
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

// 判断是否为AI伙伴
const isCompanion = computed(() => {
  return props.selectedChat?.name === '空白AI' || props.selectedChat?.type === 'companion'
})

// 解析技能数据
const skills = computed(() => {
  if (!props.selectedChat?.skills) {
    return []
  }
  
  try {
    let skillsData = props.selectedChat.skills
    
    // 如果skills是字符串，解析为JSON
    if (typeof skillsData === 'string') {
      // 处理可能的转义字符
      skillsData = skillsData.replace(/\\"/g, '"')
      const parsed = JSON.parse(skillsData)
      return parsed
    }
    // 如果已经是数组，直接返回
    if (Array.isArray(skillsData)) {
      return skillsData
    }
    
    return []
  } catch (error) {
    console.error('解析技能数据失败:', error)
    return []
  }
})

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
  border: 4px solid #e9ebed;
  display: flex;
  align-items: center;
  justify-content: center;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  // 粒子小球样式调整
  .particle-avatar {
    width: 100%;
    height: 100%;
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

.skills-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  margin-bottom: 12px;
}

.skill-tag {
  background: #c5cad0;
  color: white;
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
  transition: all 0.2s ease;
  
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(82, 180, 180, 0.4);
  }
}

.profile-section {
  margin-bottom: 24px;
}


.section-arrow {
  color: #94a3b8;
  cursor: pointer;
}

.background-story {
  font-size: 13px;
  color: #475569;
  line-height: 1.6;
  font-weight: 300;
  background: rgba(255, 255, 255, 0.675);
  padding: 13px 10px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  text-align: left;
}

// AI伙伴相关样式
.growth-section {
  margin-bottom: 24px;
}

.growth-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  
  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
  
  .growth-percentage {
    font-size: 18px;
    font-weight: 700;
    color: #10b981;
  }
}

.growth-progress {
  .progress-bar {
    width: 100%;
    height: 8px;
    background: #e5e7eb;
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 12px;
    
    .progress-fill {
      height: 100%;
      border-radius: 4px;
      transition: width 0.3s ease;
    }
  }
  
  .growth-details {
    display: flex;
    justify-content: space-between;
    
    .growth-item {
      text-align: center;
      
      .label {
        display: block;
        font-size: 12px;
        color: #6b7280;
        margin-bottom: 4px;
      }
      
      .value {
        font-size: 14px;
        font-weight: 600;
        color: #1f2937;
      }
    }
  }
}

.abilities-section {
  margin-bottom: 24px;
  
  h3 {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
  
  .abilities-grid {
    display: flex;
    flex-direction: column;
    gap: 12px;
    
    .ability-item {
      display: flex;
      align-items: center;
      gap: 12px;
      
      .ability-name {
        font-size: 13px;
        color: #374151;
        min-width: 60px;
      }
      
      .ability-bar {
        flex: 1;
        height: 6px;
        background: #e5e7eb;
        border-radius: 3px;
        overflow: hidden;
        
        .ability-fill {
          height: 100%;
          background: #52b4b4;
          border-radius: 3px;
          transition: width 0.3s ease;
        }
      }
      
      .ability-value {
        font-size: 12px;
        color: #6b7280;
        min-width: 35px;
        text-align: right;
      }
    }
  }
}
</style>
