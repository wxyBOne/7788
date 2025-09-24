<template>
  <div class="chat-page">
    <!-- 左侧导航栏 -->
    <Sidebar @showBlankAI="showBlankAI = true" @showDiary="showDiary = true" @showSettings="showSettings = true" />

      <!-- 主内容区域 -->
      <div class="main-content">
        <!-- 光影效果 -->
        <div class="red-glow-deep"></div>
        <div class="blue-glow-deep"></div>
        <div class="yellow-glow-deep"></div>
        <div class="red-glow-light"></div>
        <div class="blue-glow-light"></div>
        <div class="yellow-glow-light"></div>
        <div class="red-glow-extra"></div>
        <div class="blue-glow-extra"></div>
        
        <!-- 顶部导航栏 -->
        <div class="top-header">
        <div class="header-left">
          <div class="logo">
            <span class="logo-text">Seven</span>
          </div>
        </div>
        <div class="header-center">
          <h1 class="page-title">赫敏</h1>
        </div>
        <div class="header-right">
          <button class="upgrade-btn">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M6.62 10.79c1.44 2.83 3.76 5.14 6.59 6.59l2.2-2.2c.27-.27.67-.36 1.02-.24 1.12.37 2.33.57 3.57.57.55 0 1 .45 1 1V20c0 .55-.45 1-1 1-9.39 0-17-7.61-17-17 0-.55.45-1 1-1h3.5c.55 0 1 .45 1 1 0 1.25.2 2.45.57 3.57.11.35.03.74-.25 1.02l-2.2 2.2z"/>
            </svg>
            通话
          </button>
        </div>
      </div>


      <!-- 底部区域 -->
      <div class="bottom-section">
        
        <!-- 输入区域 -->
        <div class="input-area">
        <div class="input-container">
            <div class="input-icon-wrapper">
              <svg class="input-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
              </svg>
            </div>
            <input v-model="currentMessage" @keyup.enter="sendMessage" 
                   placeholder='例如："用简单的话解释量子计算"' />
            <div class="input-buttons">
              <button @click="startVoiceCall" class="voice-btn">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 14c1.66 0 2.99-1.34 2.99-3L15 5c0-1.66-1.34-3-3-3S9 3.34 9 5v6c0 1.66 1.34 3 3 3zm5.3-3c0 3-2.54 5.1-5.3 5.1S6.7 14 6.7 11H5c0 3.41 2.72 6.23 6 6.72V21h2v-3.28c3.28-.48 6-3.3 6-6.72h-1.7z"/>
                </svg>
              </button>
              <button @click="sendMessage" class="send-btn">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
                </svg>
              </button>
            </div>
            </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Sidebar from '@/components/Sidebar.vue'
import { ref } from 'vue'

const showBlankAI = ref(false)
const showDiary = ref(false)
const showSettings = ref(false)
const currentMessage = ref('')

const sendMessage = () => {
  if (!currentMessage.value.trim()) return
  console.log('发送消息:', currentMessage.value)
  currentMessage.value = ''
}

const startVoiceCall = () => {
  console.log('开始语音通话')
}

</script>

<style lang="scss" scoped>
.chat-page {
  display: flex;
  height: 100vh;
  background-color: #ffffff;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Arial, sans-serif;
  position: relative;
  
  // 背景光影效果
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: 
      radial-gradient(circle at 20% 30%, rgba(100, 200, 255, 0.03) 0%, transparent 50%),
      radial-gradient(circle at 80% 20%, rgba(255, 200, 100, 0.02) 0%, transparent 50%),
      radial-gradient(circle at 60% 70%, rgba(255, 100, 200, 0.025) 0%, transparent 50%),
      radial-gradient(circle at 30% 80%, rgba(200, 255, 100, 0.02) 0%, transparent 50%);
    pointer-events: none;
    z-index: 0;
  }
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 40px;
  background: white;
  margin: 20px 40px 20px 20px;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
  overflow: hidden;
  
  // 光影效果
  .red-glow-deep {
    position: absolute;
    top: -20%;
    right: 0;
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, 
      rgba(255, 100, 100, 0.12) 0%, 
      rgba(255, 100, 100, 0.086) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(80px);
    pointer-events: none;
    z-index: 1;
  }
  
  .blue-glow-deep {
    position: absolute;
    top: 50%;
    right: 0%;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, 
      rgba(100, 200, 255, 0.13) 0%, 
      rgba(100, 201, 255, 0.151) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(80px);
    pointer-events: none;
    z-index: 1;
  }
  
  .yellow-glow-deep {
    position: absolute;
    top: 40%;
    right: -10%;
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, 
      rgba(255, 200, 0, 0.11) 0%, 
      rgba(255, 200, 0, 0.116) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(80px);
    pointer-events: none;
    z-index: 1;
  }
  
  .red-glow-light {
    position: absolute;
    top: 25%;
    right: 50%;
    width: 180px;
    height: 180px;
    background: radial-gradient(circle, 
      rgba(255, 100, 100, 0.08) 0%, 
      rgba(255, 100, 100, 0.04) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(35px);
    pointer-events: none;
    z-index: 1;
  }
  
  .blue-glow-light {
    position: absolute;
    bottom: 30%;
    right: 45%;
    width: 160px;
    height: 160px;
    background: radial-gradient(circle, 
      rgba(100, 200, 255, 0.07) 0%, 
      rgba(100, 200, 255, 0.03) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(32px);
    pointer-events: none;
    z-index: 1;
  }
  
  .yellow-glow-light {
    position: absolute;
    top: 60%;
    right: 30%;
    width: 140px;
    height: 140px;
    background: radial-gradient(circle, 
      rgba(255, 200, 0, 0.06) 0%, 
      rgba(255, 200, 0, 0.03) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(30px);
    pointer-events: none;
    z-index: 1;
  }
  
  .red-glow-extra {
    position: absolute;
    top: 12%;
    right: 5%;
    width: 120px;
    height: 120px;
    background: radial-gradient(circle, 
      rgba(255, 100, 100, 0.05) 0%, 
      rgba(255, 100, 100, 0.02) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(25px);
    pointer-events: none;
    z-index: 1;
  }
  
  .blue-glow-extra {
    position: absolute;
    top: 35%;
    right: 60%;
    width: 100px;
    height: 100px;
    background: radial-gradient(circle, 
      rgba(100, 200, 255, 0.04) 0%, 
      rgba(100, 200, 255, 0.02) 50%, 
      transparent 100%);
    border-radius: 50%;
    filter: blur(20px);
    pointer-events: none;
    z-index: 1;
  }
  
  // 确保内容在光影之上
  > * {
    position: relative;
    z-index: 2;
  }
}

// 顶部导航栏
.top-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 0 120px 20px 120px;
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
  
  .page-title {
    font-size: 22px;
  letter-spacing: 0.1em;
    font-weight: 500;
    background: linear-gradient(135deg, #1a202c 0%, #1a202c 33%, #1a202c 66%, #1a202c 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin: 0;
  }
}

.header-right {
  display: flex;
  align-items: center;
}

.upgrade-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: #1a202c;
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  
  &:hover {
    background: #2d3748;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(26, 32, 44, 0.3);
  }
  
  svg {
    width: 14px;
    height: 14px;
  }
}


// 底部区域
.bottom-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 40px 0;
}


// 输入区域
.input-area {
  width: 100%;
  max-width: 1000px;
  padding: 0 20px;
}

.input-container {
  display: flex;
  align-items: center;
  background: white;
  border: none;
  border-radius: 24px;
  padding: 24px 20px;
  margin-bottom: 24px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  
  &:focus-within {
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  }
  
  .input-icon-wrapper {
    width: 40px;
    height: 40px;
    background: #f5f5f5;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 16px;
    
    .input-icon {
      width: 20px;
      height: 20px;
      color: #666;
    }
  }
  
  input {
    flex: 1;
    border: none;
    outline: none;
    font-size: 16px;
    color: #1a202c;
    background: transparent;
    
    &::placeholder {
      color: #999;
    }
  }

  .input-buttons {
    display: flex;
    gap: 12px;
    margin-left: 16px;
  }
  
  .voice-btn {
    width: 40px;
    height: 40px;
    background: #f5f5f5;
    border: none;
    border-radius: 50%;
    color: #666;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    
    &:hover {
      background: #e0e0e0;
    }
    
    svg {
      width: 20px;
      height: 20px;
    }
  }
  
  .send-btn {
    width: 40px;
    height: 40px;
    background: #1a202c;
    border: none;
    border-radius: 50%;
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    
    &:hover {
      background: #2d3748;
      transform: scale(1.05);
    }
    
    svg {
      width: 20px;
      height: 20px;
    }
  }
}


// 响应式设计
@media (max-width: 768px) {
  .chat-page {
    flex-direction: column;
  }
  
  .main-content {
    margin: 0;
    border-radius: 0;
    padding: 20px;
  }
  
  .input-area {
    padding: 0 10px;
  }
  
  .input-container {
    padding: 20px 24px;
  }
  
  .input-icon-wrapper {
    width: 36px;
    height: 36px;
    margin-right: 12px;
    
    .input-icon {
      width: 18px;
      height: 18px;
    }
  }
  
  .voice-btn, .send-btn {
    width: 36px;
    height: 36px;
    
    svg {
      width: 18px;
      height: 18px;
    }
  }
  
  input {
    font-size: 14px;
  }
  
  .search-container .search-input {
    width: 200px;
  }
}
</style>