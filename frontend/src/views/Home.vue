<template>
  <div class="chat-page">
    <!-- 左侧导航栏 -->
    <ChatSidebar v-show="sidebarVisible" @showBlankAI="showBlankAI = true" @showDiary="showDiary = true" @showSettings="showSettings = true" />

    <!-- 收缩按钮 -->
    <div class="collapse-btn" @click="toggleSidebar">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M3 6h18M3 12h18M3 18h18"/>
      </svg>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content" :class="{ 'sidebar-open': sidebarVisible }">
      
      <!-- 顶部导航栏 -->
      <div class="top-header">
        <div class="header-left">
          <div class="logo">
            <span class="logo-text">Seven</span>
          </div>
        </div>
        <div class="header-center">
          <!-- 标题已移除 -->
        </div>
        <div class="header-right">
          <div class="search-container">
            <input type="text" placeholder="搜索你喜欢的角色..." class="search-input" />
            <svg class="search-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
            </svg>
          </div>
        </div>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <!-- 欢迎区域 - 居中 -->
        <div class="welcome-section">
          <div class="welcome-content">
            <h2>欢迎来到 Seven AI 世界</h2>
            <div class="ai-bot">
              <div class="bot-avatar">
                <div class="bot-face">
                  <div class="bot-eyes">
                    <div class="eye left"></div>
                    <div class="eye right"></div>
                  </div>
                  <div class="bot-mouth"></div>
                </div>
              </div>
              <div class="speech-bubble">
                <p>你好！我是Seven AI助手，你可以搜索你喜欢的角色开始对话，或者培养属于你的专属AI伙伴！</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 功能方块区域 - 按照图片布局 -->
        <div class="feature-blocks">
          <!-- 左上角 - 小方块 -->
          <div class="feature-block small" 
               :class="{ active: hoveredCard === 0 }"
               @mouseenter="handleMouseEnter(0)" 
               @mouseleave="handleMouseLeave">
            <div class="block-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
              </svg>
            </div>
          </div>

          <!-- 左中 - 大方块 -->
          <div class="feature-block large" 
               :class="{ active: hoveredCard === 1 }"
               @mouseenter="handleMouseEnter(1)" 
               @mouseleave="handleMouseLeave">
            <div class="block-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
          </div>

          <!-- 左下角 - 小方块 -->
          <div class="feature-block small" 
               :class="{ active: hoveredCard === 2 }"
               @mouseenter="handleMouseEnter(2)" 
               @mouseleave="handleMouseLeave">
            <div class="block-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/>
              </svg>
            </div>
          </div>

          <!-- 右上角 - 小方块 -->
          <div class="feature-block small" 
               :class="{ active: hoveredCard === 3 }"
               @mouseenter="handleMouseEnter(3)" 
               @mouseleave="handleMouseLeave">
            <div class="block-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
            </div>
          </div>

          <!-- 右下角 - 小方块 -->
          <div class="feature-block large" 
               :class="{ active: hoveredCard === 4 }"
               @mouseenter="handleMouseEnter(4)" 
               @mouseleave="handleMouseLeave">
            <div class="block-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm5 11H7v-2h10v2z"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 功能卡片 -->
        <div class="feature-cards">
          <!-- 角色对话 -->
          <div class="feature-card" 
               :class="{ visible: hoveredCard === 0 }"
               @click="$router.push('/chat')"
               @mouseenter="handleMouseEnter(0)"
               @mouseleave="handleMouseLeave">
            <div class="card-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
              </svg>
            </div>
            <h3>角色对话</h3>
            <p>与林黛玉、孙悟空、李白、赫敏等经典角色深度对话</p>
          </div>

          <!-- AI养成 -->
          <div class="feature-card" 
               :class="{ visible: hoveredCard === 1 }"
               @click="showBlankAI = true"
               @mouseenter="handleMouseEnter(1)"
               @mouseleave="handleMouseLeave">
            <div class="card-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
            <h3>AI养成</h3>
            <p>培养属于你的专属AI伙伴，从空白到成熟的成长之旅</p>
          </div>

          <!-- AI日记 -->
          <div class="feature-card" 
               :class="{ visible: hoveredCard === 2 }"
               @click="showDiary = true"
               @mouseenter="handleMouseEnter(2)"
               @mouseleave="handleMouseLeave">
            <div class="card-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/>
              </svg>
            </div>
            <h3>AI日记</h3>
            <p>查看AI伙伴的成长记录，感受每一次对话的温暖回忆</p>
          </div>

          <!-- 语音通话 -->
          <div class="feature-card" 
               :class="{ visible: hoveredCard === 3 }"
               @mouseenter="handleMouseEnter(3)"
               @mouseleave="handleMouseLeave">
            <div class="card-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
            </div>
            <h3>语音通话</h3>
            <p>与AI伙伴进行实时语音对话，体验更真实的交流</p>
          </div>

          <!-- 图像识别 -->
          <div class="feature-card" 
               :class="{ visible: hoveredCard === 4 }"
               @click="showImageRecognition = true"
               @mouseenter="handleMouseEnter(4)"
               @mouseleave="handleMouseLeave">
            <div class="card-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm5 11H7v-2h10v2z"/>
              </svg>
            </div>
            <h3>图像识别</h3>
            <p>上传图片与AI伙伴分享，让AI理解你的世界</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ChatSidebar from '@/components/ChatSidebar.vue'
import { ref } from 'vue'

// 界面控制
const showBlankAI = ref(false)
const showDiary = ref(false)
const showSettings = ref(false)
const showImageRecognition = ref(false)
const hoveredCard = ref(-1)
const sidebarVisible = ref(false)

// 优化的事件处理函数
let hideTimer = null
let showTimer = null
let isHovering = false

const handleMouseEnter = (index) => {
  // 清除隐藏定时器
  if (hideTimer) {
    clearTimeout(hideTimer)
    hideTimer = null
  }
  
  // 清除显示定时器
  if (showTimer) {
    clearTimeout(showTimer)
    showTimer = null
  }
  
  isHovering = true
  
  // 如果当前已经有卡片显示且是同一个，直接返回
  if (hoveredCard.value === index) {
    return
  }
  
  // 延迟显示卡片
  showTimer = setTimeout(() => {
    if (isHovering) {
      hoveredCard.value = index
    }
  }, 200)
}

const handleMouseLeave = () => {
  isHovering = false
  
  // 清除显示定时器
  if (showTimer) {
    clearTimeout(showTimer)
    showTimer = null
  }
  
  // 延迟隐藏卡片
  hideTimer = setTimeout(() => {
    if (!isHovering) {
      hoveredCard.value = -1
    }
  }, 100)
}

const toggleSidebar = () => {
  sidebarVisible.value = !sidebarVisible.value
}
</script>

<style lang="scss" scoped>
.chat-page {
  display: flex;
  height: 100vh;
  background: #ffffff;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Arial, sans-serif;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    bottom: -100%;
    left: -4.5vw;
    width: 110vw;
    height: 82vw;
      background: radial-gradient(ellipse, 
        rgba(255, 0, 0, 0.6) 0%, 
        rgba(255, 165, 0, 0.5) 35%, 
        rgba(255, 255, 0, 0.4) 45%, 
        rgba(0, 255, 0, 0.3) 55%, 
        rgba(0, 255, 255, 0.3) 64%, 
        rgba(0, 0, 255, 0.2) 80%, 
        rgba(128, 0, 128, 0.2) 90%, 
        transparent 100%);
    border-radius: 50%;
    filter: blur(28px);
    pointer-events: none;
    z-index: 0;
    transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  }
  
  .main-content.sidebar-open &::before {
    left: calc(-4.5vw + 120px);
    width: calc(110vw - 120px);
    height: 82vw;
    filter: blur(28px);
  }
}


// 收缩按钮
.collapse-btn {
  position: fixed;
  bottom: 40px;
  left: 40px;
  width: 24px;
  height: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  z-index: 100;
  
  &:hover {
    transform: scale(1.1);
  }
  
  svg {
    width: 24px;
    height: 24px;
    color: #2d3748;
  }
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 40px;
  position: relative;
  overflow: hidden;
  margin-left: 0;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  
  &.sidebar-open {
    margin-left: 20px;
    width: calc(100% - 20px);
  }
}

// 顶部导航栏
.top-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 0 0 20px 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  padding: 8px 0;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: rgba(75, 90, 90, 0.363); /* 降级基础色 */
  
  /* 黑灰渐变（简洁专业版） */
  background-image: linear-gradient(
    130deg,
    hsl(0, 0%, 20%),    /* 深黑色 */
    hsl(0, 0%, 50%)     /* 中灰色 */
  );
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  
  /* 极淡阴影增强质感，不破坏简洁感 */
  text-shadow: 0 1px 1px rgba(0, 0, 0, 0.03);
}

.header-center {
  flex: 1;
  text-align: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.search-container {
  position: relative;
  
  .search-input {
    padding: 12px 16px 12px 40px;
    border: 2px solid #e2e8f0;
    border-radius: 25px;
    font-size: 14px;
    width: 300px;
    transition: all 0.3s ease;
    
    &:focus {
      outline: none;
      border-color: #2d3748;
      box-shadow: 0 0 0 3px rgba(45, 55, 72, 0.1), 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }
  
  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    width: 18px;
    height: 18px;
    color: #a0aec0;
  }
}

// 欢迎区域
.welcome-section {
  margin-bottom: 40px;
  padding: 30px 60px;
}

.welcome-content {
  text-align: center;
  transform: scale(0.8);
  
  h2 {
    font-size: 80px;
    font-weight: 700;
    color: #3d4757cf;
    margin: 0 0 30px 0;
  }
}

.ai-bot {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.bot-avatar {
  width: 80px;
  height: 80px;
  background: black;
  border: 5px solid white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
}

.bot-face {
  position: relative;
  width: 50px;
  height: 50px;
}

.bot-eyes {
  position: absolute;
  top: 15px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
}

.eye {
  width: 8px;
  height: 8px;
  background: #ffffff;
  border-radius: 50%;
  animation: blink 3s infinite;
}

@keyframes blink {
  0%, 90%, 100% { transform: scaleY(1); }
  95% { transform: scaleY(0.1); }
}

.bot-mouth {
  position: absolute;
  bottom: 12px;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 10px;
  border: 2px solid #ffffff;
  border-top: none;
  border-radius: 0 0 20px 20px;
}

.speech-bubble {
  background: white;
  padding: 16px 20px;
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: relative;
  max-width: 400px;
  
  &::before {
    content: '';
    position: absolute;
    left: -10px;
    top: 20px;
    width: 0;
    height: 0;
    border-top: 10px solid transparent;
    border-bottom: 10px solid transparent;
    border-right: 10px solid white;
  }
  
  p {
    margin: 0;
    color: #4a5568;
    font-size: 16px;
    line-height: 1.5;
  }
}

// 主内容区域
.main-content {
  flex: 1;
  position: relative;
}

// 内容区域
.content-area {
  position: relative;
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 80px;
  z-index: 0;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

// 欢迎区域
.welcome-section {
  text-align: center;
  position: relative;
  z-index: 0;
}

// 功能方块区域
.feature-blocks {
  position: absolute;
  top: -70px;
  left: 0;
  width: 100%;
  height: 100%;
  transform: scale(0.9);
  pointer-events: auto;
}

.feature-block {
  position: absolute;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  pointer-events: auto;
  transform: translateY(0) rotate(0deg);
  animation: float 6s ease-in-out infinite;
  z-index: 10;
  
  // 小方块尺寸
  &.small {
    width: 80px;
    height: 80px;
  }
  
  // 大方块尺寸
  &.large {
    width: 120px;
    height: 120px;
  }
  
  // 围绕中心分布 - 按照例图微调
  &:nth-child(1) { // 左上角 - 小方块
    top: 12%;
    left: 18%;
    animation-delay: 0s;
  }
  
  &:nth-child(2) { // 左中 - 大方块
    top: 41%;
    left: 6%;
    animation-delay: 2s;
  }
  
  &:nth-child(3) { // 左下角 - 小方块
    top: 75%;
    left: 18%;
    animation-delay: 4s;
  }
  
  &:nth-child(4) { // 右上角 - 小方块
    top: 17%;
    right: 15%;
    animation-delay: 1s;
    z-index: 101;
  }
  
  &:nth-child(5) { // 右下角 - 小方块
    top: 65%;
    right: 8%;
    animation-delay: 3s;
  }
  
  &:hover {
    transform: translateY(-8px) rotate(5deg);
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.15);
    background: rgba(255, 255, 255, 0.95);
  }
  
  &.active {
    opacity: 0;
    transform: translateY(-20px) rotate(10deg) scale(0.8);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  25% {
    transform: translateY(-10px) rotate(2deg);
  }
  50% {
    transform: translateY(-5px) rotate(-1deg);
  }
  75% {
    transform: translateY(-15px) rotate(1deg);
  }
}

.block-icon {
  color: #4a5568;
  transition: color 0.3s ease;
  
  svg {
    width: 100%;
    height: 100%;
  }
}

.feature-block.small .block-icon {
  width: 40px;
  height: 40px;
}

.feature-block.large .block-icon {
  width: 60px;
  height: 60px;
}

.feature-block:hover .block-icon,
.feature-block.active .block-icon {
  color: #1a202c;
}

// 功能卡片
.feature-cards {
  position: absolute;
  top: -70px;
  left: 0;
  width: 100%;
  height: 100%;
  transform: scale(0.9);
  pointer-events: none;
}

.feature-card {
  position: absolute;
  width: 280px;
  min-height: 200px;
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  opacity: 0;
  transform: translateY(15px) scale(0.98);
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  z-index: 50;
  display: flex;
  flex-direction: column;
  pointer-events: auto;
  
  // 与方块位置完全对应
  &:nth-child(1) { // 左上角
    top: 0%;
    left: 5%;
    transform: translateY(15px) scale(0.98);
  }
  
  &:nth-child(2) { // 左中
    top: 41%;
    left: -3%;
    transform: translateY(15px) scale(0.98);
  }
  
  &:nth-child(3) { // 左下角
    top: 75%;
    left: 5%;
    transform: translateY(15px) scale(0.98);
  }
  
  &:nth-child(4) { // 右上角
    top: 10%;
    right: 0%;
    transform: translateY(15px) scale(0.98);
  }
  
  &:nth-child(5) { // 右下角
    top: 68%;
    right: -3%;
    transform: translateY(15px) scale(0.98);
  }
  
  &.visible {
    opacity: 1;
    transform: translateY(0) scale(1);
    pointer-events: auto;
    z-index: 60;
    transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  }
}

.card-icon {
  width: 60px;
  height: 60px;
  background: #1a272caa;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;

  svg {
    width: 28px;
    height: 28px;
    color: white;
  }
}

.feature-card h3 {
  font-size: 20px;
  font-weight: 600;
  color: #313741c1;
  margin: 0 0 8px 0;
}

.feature-card p {
  color: #4a5568;
  font-size: 14px;
  margin: 0;
  line-height: 1.5;
}

// 响应式设计
@media (max-width: 768px) {
  .chat-page {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    flex-direction: row;
    border-radius: 0;
    margin: 0;
    
    .nav-icons {
      flex-direction: row;
      justify-content: space-around;
      padding: 16px 0;
    }
  }
  
  .main-content {
    margin: 0;
    border-radius: 0;
    padding: 20px;
  }
  
  .welcome-section {
    padding: 20px;
  }
  
  .feature-cards {
    grid-template-columns: 1fr;
    padding: 0 20px;
  }
  
  .search-container .search-input {
    width: 200px;
  }
}
</style>