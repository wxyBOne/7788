<template>
  <div class="voice-call-page">
    <div class="call-container">
      <!-- 角色头像 -->
      <div class="character-avatar">
        <img :src="character.avatar_url || character.avatar" :alt="character.name" />
      </div>
      
      <!-- 语音状态指示器 -->
      <div class="voice-indicator">
        <div class="audio-waveform">
          <div class="wave-bar" v-for="i in 4" :key="i" :style="{ height: getWaveHeight(i) }"></div>
        </div>
        <div class="status-text">{{ isTalking ? 'AI正在回复...' : isListening ? '正在聆听...' : '通话中...' }}</div>
      </div>
      
      <!-- 挂断按钮 -->
      <button class="hangup-btn" @click="$emit('hangup')">
        <svg class="hangup-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M6.62 10.79c1.44 2.83 3.76 5.14 6.59 6.59l2.2-2.2c.27-.27.67-.36 1.02-.24 1.12.37 2.33.57 3.57.57.55 0 1 .45 1 1V20c0 .55-.45 1-1 1-9.39 0-17-7.61-17-17 0-.55.45-1 1-1h3.5c.55 0 1 .45 1 1 0 1.25.2 2.45.57 3.57.11.35.03.74-.25 1.02l-2.2 2.2z"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  character: {
    type: Object,
    required: true
  },
  isListening: {
    type: Boolean,
    default: false
  },
  isTalking: {
    type: Boolean,
    default: false
  },
  volume: {
    type: Number,
    default: 0
  }
})

defineEmits(['hangup'])

// 根据音量生成波形高度
const getWaveHeight = (index) => {
  const baseHeight = 20
  const variation = props.volume * 30
  
  // 为每个波形条设置不同的基础高度，避免随机性
  const heights = [baseHeight, baseHeight + 5, baseHeight + 10, baseHeight + 15]
  const baseWaveHeight = heights[index - 1] || baseHeight
  
  // 根据音量添加变化，但保持相对稳定
  const volumeEffect = variation * (0.5 + index * 0.2)
  return `${baseWaveHeight + volumeEffect}px`
}
</script>

<style lang="scss" scoped>
.voice-call-page {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.call-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;
  width: 100%;
  max-width: 400px;
  padding: 40px 20px;
}

.character-avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 10px solid #e9ebed;
}

.character-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.voice-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #333;
}

.audio-waveform {
  display: flex;
  align-items: end;
  gap: 8px;
  height: 60px;
}

.wave-bar {
  width: 8px;
  background: #52b4b4da;
  border-radius: 4px;
  min-height: 20px;
  animation: wave 1.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  transform-origin: bottom;
}

.wave-bar:nth-child(1) { animation-delay: 0s; }
.wave-bar:nth-child(2) { animation-delay: 0.2s; }
.wave-bar:nth-child(3) { animation-delay: 0.4s; }
.wave-bar:nth-child(4) { animation-delay: 0.6s; }

.status-text {
  font-size: 15px;
  font-weight: 500;
  color: #5a5a5a;
  text-align: center;
}

.hangup-btn {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: #ef4444bb;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.hangup-btn:hover {
  background: #dc2626b3;
  transform: scale(0.9);
}

.hangup-icon {
  width: 32px;
  height: 32px;
  color: white;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.hangup-btn:hover .hangup-icon {
  transform: rotate(135deg);
}

@keyframes wave {
  0%, 100% {
    transform: scaleY(0.5);
  }
  50% {
    transform: scaleY(1.2);
  }
}
</style>
