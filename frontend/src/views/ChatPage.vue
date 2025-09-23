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

      <!-- 欢迎区域 -->
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

      <!-- 功能卡片 -->
      <div class="feature-cards">
        <div class="feature-card" @click="$router.push('/chat')">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
          </div>
          <h3>角色对话</h3>
          <p>与林黛玉、孙悟空、李白、赫敏等经典角色深度对话</p>
        </div>
        <div class="feature-card" @click="showBlankAI = true">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
          </div>
          <h3>AI养成</h3>
          <p>培养属于你的专属AI伙伴，从空白到成熟的成长之旅</p>
        </div>
        <div class="feature-card" @click="showDiary = true">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/>
            </svg>
          </div>
          <h3>AI日记</h3>
          <p>查看AI伙伴的成长记录，感受每一次对话的温暖回忆</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from '@/components/Sidebar.vue'

export default {
  name: 'ChatPage',
  components: {
    Sidebar
  },
  data() {
    return {
      // 界面控制
      showBlankAI: false,
      showDiary: false,
      showSettings: false
    }
  },
  methods: {
    // 组件方法
  }
}
</script>

<style lang="scss" scoped>
.chat-page {
  display: flex;
  height: 100vh;
  background-color: #f4feff;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Arial, sans-serif;
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
  padding: 0 0 20px 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
}
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
      border-color: #1a202c;
      box-shadow: 0 0 0 3px rgba(26, 32, 44, 0.1);
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
  
h2 {
  font-size: 32px;
  font-weight: 700;
  color: #1a202c;
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
  background: white;
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
  border: 2px solid white;
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

// 功能卡片
.feature-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  padding: 0 60px;
}

.feature-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
}

.card-icon {
  width: 60px;
  height: 60px;
  background: #1a202c;
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
  color: #1a202c;
  margin: 0 0 8px 0;
}

.feature-card p {
  color: #718096;
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