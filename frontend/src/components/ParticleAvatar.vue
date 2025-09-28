<template>
  <div class="particle-avatar" :style="avatarStyle">
    <!-- 粒子小球主体 -->
    <div class="particle-core" :style="coreStyle">
      <!-- 粒子效果 -->
      <div class="particles" :style="particlesStyle">
        <div 
          v-for="i in particleCount" 
          :key="i" 
          class="particle" 
          :style="getParticleStyle(i)"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  emotion: {
    type: String,
    default: '平静'
  },
  intensity: {
    type: Number,
    default: 0.5
  },
  color: {
    type: String,
    default: '#52b4b4'
  },
  brightness: {
    type: Number,
    default: 0.7
  },
  particleSpeed: {
    type: Number,
    default: 0.5
  },
  size: {
    type: String,
    default: 'medium' // small, medium, large
  }
})

// 粒子数量
const particleCount = computed(() => {
  return Math.floor(props.intensity * 15) + 5 // 5-20个粒子
})

// 情绪颜色映射
const emotionColor = computed(() => {
  const colorMap = {
    '开心': '#ffd700', // 暖黄
    '好奇': '#00bfff', // 闪烁蓝
    '孤单': '#dda0dd', // 淡紫
    '悲伤': '#87ceeb', // 天蓝
    '兴奋': '#ff6347', // 番茄红
    '平静': '#52b4b4', // 柔绿
    '依赖': '#ffb6c1', // 浅粉
    '愤怒': '#ff4500', // 橙红
    '恐惧': '#9370db', // 中紫
    '惊讶': '#ffa500'  // 橙色
  }
  return colorMap[props.emotion] || props.color
})

// 头像容器样式
const avatarStyle = computed(() => {
  const sizeMap = {
    small: '40px',
    medium: '60px',
    large: '80px'
  }
  return {
    width: sizeMap[props.size],
    height: sizeMap[props.size]
  }
})

// 核心样式
const coreStyle = computed(() => {
  return {
    background: `radial-gradient(circle, ${emotionColor.value} 0%, ${emotionColor.value}80 70%, transparent 100%)`,
    opacity: props.brightness,
    boxShadow: `0 0 ${20 * props.intensity}px ${emotionColor.value}40`
  }
})

// 粒子容器样式
const particlesStyle = computed(() => {
  return {
    animationDuration: `${2 / props.particleSpeed}s`
  }
})

// 获取单个粒子样式
const getParticleStyle = (index) => {
  const angle = (360 / particleCount.value) * index
  const distance = 12 + Math.random() * 15 // 增加距离变化
  const size = 1.5 + Math.random() * 4 // 增加大小变化
  
  // 根据情绪调整动画参数
  let animationDuration = 1.5 + Math.random() * 2
  let animationDelay = Math.random() * 3
  
  if (props.emotion === '兴奋') {
    animationDuration = 0.8 + Math.random() * 1
    animationDelay = Math.random() * 1.5
  } else if (props.emotion === '孤单') {
    animationDuration = 2 + Math.random() * 3
    animationDelay = Math.random() * 4
  }
  
  return {
    transform: `rotate(${angle}deg) translateX(${distance}px)`,
    width: `${size}px`,
    height: `${size}px`,
    backgroundColor: emotionColor.value,
    opacity: 0.4 + Math.random() * 0.6, // 增加透明度变化
    animationDelay: `${animationDelay}s`,
    animationDuration: `${animationDuration}s`,
    // 添加随机颜色变化
    filter: `hue-rotate(${Math.random() * 30 - 15}deg) brightness(${0.8 + Math.random() * 0.4})`
  }
}

// 动画控制
const animationId = ref(null)

onMounted(() => {
  // 启动粒子动画
  animateParticles()
})

onUnmounted(() => {
  if (animationId.value) {
    cancelAnimationFrame(animationId.value)
  }
})

const animateParticles = () => {
  // 这里可以添加更复杂的粒子动画逻辑
  animationId.value = requestAnimationFrame(animateParticles)
}
</script>

<style lang="scss" scoped>
.particle-avatar {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  overflow: hidden; // 改为hidden，防止粒子超出边界
}

.particle-core {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.particles {
  position: absolute;
  width: 100%;
  height: 100%;
  animation: rotate 3s linear infinite, float 2s ease-in-out infinite;
}

.particle {
  position: absolute;
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform-origin: 0 0;
  animation: pulse 1.5s ease-in-out infinite, drift 3s ease-in-out infinite;
}


@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.3);
    opacity: 1;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-3px);
  }
}

@keyframes drift {
  0%, 100% {
    transform: translateX(0px) translateY(0px);
  }
  25% {
    transform: translateX(1px) translateY(-1px);
  }
  50% {
    transform: translateX(-1px) translateY(1px);
  }
  75% {
    transform: translateX(1px) translateY(1px);
  }
}

// 情绪特定动画
.particle-avatar[data-emotion="开心"] .particle-core {
  animation: bounce 0.8s ease-in-out infinite, glow 2s ease-in-out infinite;
}

.particle-avatar[data-emotion="开心"] .particles {
  animation: rotate 2s linear infinite, sparkle 1s ease-in-out infinite;
}

.particle-avatar[data-emotion="好奇"] .particles {
  animation: rotate 1.5s linear infinite, flicker 0.3s ease-in-out infinite alternate, pulse 1s ease-in-out infinite;
}

.particle-avatar[data-emotion="孤单"] .particle-core {
  animation: fade 2s ease-in-out infinite, gentle-sway 4s ease-in-out infinite;
}

.particle-avatar[data-emotion="兴奋"] .particle-core {
  animation: shake 0.3s ease-in-out infinite, vibrate 0.5s ease-in-out infinite;
}

.particle-avatar[data-emotion="兴奋"] .particles {
  animation: rotate 1s linear infinite, rapid-pulse 0.4s ease-in-out infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-5px);
  }
}

@keyframes flicker {
  0% {
    opacity: 0.7;
  }
  100% {
    opacity: 1;
  }
}

@keyframes fade {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 0.8;
  }
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-2px);
  }
  75% {
    transform: translateX(2px);
  }
}

@keyframes glow {
  0%, 100% {
    box-shadow: 0 0 10px rgba(255, 215, 0, 0.3);
  }
  50% {
    box-shadow: 0 0 20px rgba(255, 215, 0, 0.6);
  }
}

@keyframes sparkle {
  0%, 100% {
    opacity: 0.7;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.1);
  }
}

@keyframes gentle-sway {
  0%, 100% {
    transform: translateX(0) rotate(0deg);
  }
  25% {
    transform: translateX(-1px) rotate(-1deg);
  }
  75% {
    transform: translateX(1px) rotate(1deg);
  }
}

@keyframes vibrate {
  0%, 100% {
    transform: translateX(0) translateY(0);
  }
  10% {
    transform: translateX(-1px) translateY(-1px);
  }
  20% {
    transform: translateX(1px) translateY(1px);
  }
  30% {
    transform: translateX(-1px) translateY(1px);
  }
  40% {
    transform: translateX(1px) translateY(-1px);
  }
  50% {
    transform: translateX(-1px) translateY(-1px);
  }
  60% {
    transform: translateX(1px) translateY(1px);
  }
  70% {
    transform: translateX(-1px) translateY(1px);
  }
  80% {
    transform: translateX(1px) translateY(-1px);
  }
  90% {
    transform: translateX(-1px) translateY(-1px);
  }
}

@keyframes rapid-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.4);
    opacity: 1;
  }
}
</style>
