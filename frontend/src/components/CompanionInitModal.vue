<template>
  <div v-if="show" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>创建你的AI伙伴</h3>
        <button class="close-btn" @click="closeModal">×</button>
      </div>
      
      <div class="modal-body">
        <div class="form-group">
          <label>给你的AI伙伴起个名字：</label>
          <input 
            v-model="companionName" 
            type="text" 
            placeholder="例如：小智、小爱、小助手..."
            maxlength="20"
            class="name-input"
          />
        </div>
        
        <div class="form-group">
          <label>选择成长周期：</label>
          <div class="growth-mode-options">
            <label class="option-item">
              <input 
                v-model="growthMode" 
                type="radio" 
                value="short" 
                name="growthMode"
              />
              <div class="option-content">
                <div class="option-title">快速成长</div>
                <div class="option-desc">3-7天快速形成个性，适合想要快速体验的用户</div>
              </div>
            </label>
            
            <label class="option-item">
              <input 
                v-model="growthMode" 
                type="radio" 
                value="long" 
                name="growthMode"
              />
              <div class="option-content">
                <div class="option-title">长期养成</div>
                <div class="option-desc">30-90天深度成长，适合想要长期陪伴的用户</div>
              </div>
            </label>
          </div>
        </div>
        
        <div class="preview-section">
          <div class="preview-title">预览：</div>
          <div class="preview-content">
            <div class="preview-avatar">
              <ParticleAvatar
                :emotion="'好奇'"
                :intensity="0.8"
                :color="'#52b4b4'"
                :brightness="0.7"
                :particle-speed="0.6"
                :growth-percentage="0"
                size="medium"
              />
            </div>
            <div class="preview-info">
              <div class="preview-name">{{ companionName || '未命名' }}</div>
              <div class="preview-mode">{{ growthMode === 'short' ? '快速成长' : '长期养成' }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button 
          class="create-btn" 
          @click="createCompanion"
          :disabled="!companionName.trim() || !growthMode"
        >
          创建AI伙伴
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import ParticleAvatar from './ParticleAvatar.vue'
import api from '@/services/api.js'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'created'])

const companionName = ref('')
const growthMode = ref('long') // 默认长期养成

const closeModal = () => {
  emit('close')
}

const createCompanion = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await api.companion.createCompanion(token, {
      name: companionName.value.trim(),
      growth_mode: growthMode.value
    })
    
    if (response.success) {
      emit('created', response.companion)
      closeModal()
    } else {
      alert('创建失败：' + response.message)
    }
  } catch (error) {
    console.error('创建AI伙伴失败:', error)
    alert('创建失败，请重试')
  }
}

// 重置表单
watch(() => props.show, (newVal) => {
  if (newVal) {
    companionName.value = ''
    growthMode.value = 'long'
  }
})
</script>

<style scoped>
.modal-overlay {
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
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e9ebed;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #666;
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.name-input {
  width: 100%;
  padding: 12px;
  border: 2px solid #e9ebed;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s;
}

.name-input:focus {
  outline: none;
  border-color: #52b4b4;
}

.growth-mode-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-item {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  border: 2px solid #e9ebed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.option-item:hover {
  border-color: #52b4b4;
  background: #f8f9fa;
}

.option-item input[type="radio"] {
  margin-right: 12px;
  margin-top: 2px;
}

.option-item input[type="radio"]:checked + .option-content {
  color: #52b4b4;
}

.option-content {
  flex: 1;
}

.option-title {
  font-weight: 600;
  margin-bottom: 4px;
  color: #333;
}

.option-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.4;
}

.preview-section {
  margin-top: 24px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.preview-title {
  font-weight: 600;
  margin-bottom: 16px;
  color: #333;
}

.preview-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.preview-avatar {
  flex-shrink: 0;
}

.preview-info {
  flex: 1;
}

.preview-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.preview-mode {
  font-size: 14px;
  color: #666;
}

.modal-footer {
  padding: 20px 24px;
  border-top: 1px solid #e9ebed;
  display: flex;
  justify-content: flex-end;
}

.create-btn {
  background: #52b4b4;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.create-btn:hover:not(:disabled) {
  background: #429a9a;
}

.create-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>
