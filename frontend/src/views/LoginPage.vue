<template>
  <div class="login-page">
    <!-- 背景光影效果 - 完全照搬首屏 -->
    <div class="background-effects">
      <div class="base-layer"></div>
      <div class="backdrop-contrast"></div>
      <div class="emotion-layer">
        <div class="light-field">
          <div class="cluster-mask">
            <div class="swirl">
              <div class="gradient-waves"></div>
              <div class="orbit orbit-a"></div>
              <div class="orbit orbit-b"></div>
              <div class="light-node" style="--index: 0; --d: 1.6s; top: 20%; left: 30%; width: 8px; height: 8px; background: radial-gradient(circle, rgba(255, 255, 255, 0.8) 0%, rgba(255, 255, 255, 0.4) 50%, transparent 100%);"></div>
              <div class="light-node" style="--index: 1; --d: 2.1s; top: 60%; left: 20%; width: 6px; height: 6px; background: radial-gradient(circle, rgba(255, 255, 255, 0.7) 0%, rgba(255, 255, 255, 0.3) 50%, transparent 100%);"></div>
              <div class="light-node" style="--index: 2; --d: 1.8s; top: 40%; left: 70%; width: 10px; height: 10px; background: radial-gradient(circle, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.5) 50%, transparent 100%);"></div>
              <div class="light-node" style="--index: 3; --d: 2.3s; top: 80%; left: 60%; width: 7px; height: 7px; background: radial-gradient(circle, rgba(255, 255, 255, 0.6) 0%, rgba(255, 255, 255, 0.2) 50%, transparent 100%);"></div>
              <div class="light-node" style="--index: 4; --d: 1.9s; top: 15%; left: 80%; width: 5px; height: 5px; background: radial-gradient(circle, rgba(255, 255, 255, 0.8) 0%, rgba(255, 255, 255, 0.4) 50%, transparent 100%);"></div>
            </div>
          </div>
        </div>
      </div>
      <div class="heartbeat"></div>
    </div>

    <!-- 主内容区域 -->
    <div class="login-container">
      <!-- 左侧欢迎区域 -->
      <div class="welcome-section">
        <div class="slogan-container">
          <div class="slogan-item slogan-1">与AI角色建立不可定义的关系</div>
          <div class="slogan-item slogan-2">深度情感交流与专属养成体验</div>
          <div class="slogan-item slogan-3">记录每一个美好互动的瞬间</div>
          <div class="slogan-item slogan-4">探索AI角色的无限可能性</div>
          <div class="slogan-item slogan-5">开启你的专属AI伙伴之旅</div>
        </div>
      </div>

      <!-- 右侧登录表单 -->
      <div class="form-section">
        <div class="login-form">
          <div class="form-header">
            <h2>欢迎来到 Seven</h2>
            <p>开启你的AI角色养成之旅</p>
          </div>

          <form @submit.prevent="handleQuickLogin">
            <div class="input-group">
              <input
                v-model="formData.email"
                type="text"
                placeholder="请输入邮箱"
                class="form-input"
              />
            </div>

            <div class="input-group">
              <input
                v-model="formData.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                class="form-input"
              />
              <button
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
              >
                <svg v-if="showPassword" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>

          <div class="form-options">
            <label class="remember-me">
              <input type="checkbox" v-model="formData.rememberMe">
              <span>记住我</span>
            </label>
            <a href="#" class="forgot-password" @click.prevent="openForgotPasswordModal">忘记密码？</a>
          </div>

            <button type="submit" class="login-btn" :disabled="isLoading">
              <span v-if="!isLoading">一键登录</span>
              <span v-else class="loading-spinner"></span>
            </button>
          </form>

          <div class="login-tip">
            <span>没有账号？系统将自动为您创建</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 忘记密码弹窗组件 -->
    <ForgotPasswordModal 
      :visible="showForgotPasswordModal" 
      @close="closeForgotPasswordModal" 
    />
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import ForgotPasswordModal from '@/components/ForgotPasswordModal.vue'
import chatService from '@/services/chatService.js'

const router = useRouter()

// 响应式数据
const isLoading = ref(false)
const showPassword = ref(false)
const showForgotPasswordModal = ref(false)

// 表单数据
const formData = reactive({
  email: '',
  password: '',
  rememberMe: false
})

// 一键登录方法（自动注册或登录）
const handleQuickLogin = async () => {
  if (!formData.email || !formData.password) {
    alert('请输入邮箱和密码')
    return
  }

  isLoading.value = true
  try {
    // 使用chatService的一键登录功能
    const response = await chatService.quickLogin(formData.email, formData.password)
    
    if (response.success) {
      // 登录成功，跳转到聊天页面
      router.push('/home')
    } else {
      alert('登录失败：' + (response.error || '未知错误'))
    }
  } catch (error) {
    console.error('一键登录失败:', error)
    alert('登录失败：' + error.message)
  } finally {
    isLoading.value = false
  }
}

// 打开忘记密码弹窗
const openForgotPasswordModal = () => {
  showForgotPasswordModal.value = true
}

// 关闭忘记密码弹窗
const closeForgotPasswordModal = () => {
  showForgotPasswordModal.value = false
}
</script>

<style lang="scss" scoped>
.login-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

// 背景光影效果 - 完全照搬首屏
.background-effects {
  position: absolute;
  inset: 0;
  z-index: 0;
}

// 基础层：淡色渐变背景
.base-layer { 
  position: absolute;
  inset: 0;
  z-index: 0;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 50%, #f1f5f9 100%);
}

// 背板对比层：红蓝黄色调与心跳同步
.backdrop-contrast {
  position: absolute;
  inset: -10% -10% -10% -10%;
  z-index: 0;
  background:
    radial-gradient(50% 70% at 30% 45%, rgba(255, 0, 0, 0.3) 0%, rgba(255, 0, 0, 0.1) 50%, rgba(0,0,0,0) 80%),
    radial-gradient(55% 75% at 70% 55%, rgba(0, 0, 255, 0.1) 0%, rgba(0, 0, 255, 0.05) 50%, rgba(0,0,0,0) 80%),
    radial-gradient(35% 50% at 50% 50%, rgba(255, 255, 0, 0.2) 0%, rgba(255, 255, 0, 0.05) 60%, rgba(0,0,0,0) 85%);
  filter: blur(32px) saturate(1.05);
}

// 情绪光影层：心脏跳动效果
.emotion-layer { 
  position: absolute;
  inset: 0;
  z-index: 1;
  transition: all 0.7s cubic-bezier(0.4, 0, 0.2, 1);
  
  &.transitioning {
    transform: scale(1.8);
    opacity: 0;
    filter: blur(20px);
  }
}

.light-field { 
  position: absolute;
  inset: 0;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cluster-mask { 
  position: relative;
  width: min(700px, 80vw);
  height: min(700px, 80vw);
  border-radius: 50%;
  overflow: hidden;
  filter: blur(20px);
}

.swirl { 
  position: absolute;
  inset: 0;
  animation: heartbeatBreath 1.2s ease-in-out infinite;
  transform-origin: center;
}

@keyframes heartbeatBreath {
  0% {
    transform: scale(0.9) rotate(0deg);
    filter: brightness(0.95) blur(20px);
  }
  25% {
    transform: scale(0.78) rotate(-3deg);
    filter: brightness(1.4) blur(16px);
  }
  55% {
    transform: scale(1.18) rotate(4deg);
    filter: brightness(1.15) blur(24px);
  }
  100% {
    transform: scale(0.9) rotate(0deg);
    filter: brightness(0.95) blur(20px);
  }
}

.gradient-waves { 
  position: absolute;
  inset: -25%;
  filter: blur(70px);
  opacity: 0.8;
  mix-blend-mode: screen;
  background: conic-gradient(
    from 0deg at 30% 40%, 
    rgba(255, 0, 0, 0.3), 
    rgba(255, 165, 0, 0.3), 
    rgba(255, 255, 0, 0.3), 
    rgba(0, 255, 0, 0.2),
    rgba(0, 255, 255, 0.2),
    rgba(0, 0, 255, 0.1),
    rgba(128, 0, 128, 0.1),
    rgba(255, 0, 0, 0.4), 
    rgba(100, 200, 255, 0.4)
  );
  animation: rotateWave 20s linear infinite;
}

@keyframes rotateWave { 
  to { transform: rotate(360deg); }
}

.orbit { 
  position: absolute;
  inset: 0;
  border-radius: 50%;
  box-shadow: 0 0 140px rgba(255, 255, 255, 0.2) inset;
  mix-blend-mode: screen;
}

.orbit-a { 
  animation: orbitA 4.8s ease-in-out infinite;
}

.orbit-b { 
  animation: orbitB 5.4s ease-in-out infinite;
}

@keyframes orbitA { 
  0%, 100% { transform: rotate(0deg) scale(0.9); opacity: 0.4; }
  50% { transform: rotate(25deg) scale(1.1); opacity: 0.7; }
}

@keyframes orbitB { 
  0%, 100% { transform: rotate(0deg) scale(1.05); opacity: 0.35; }
  50% { transform: rotate(-22deg) scale(1.15); opacity: 0.65; }
}

.light-node { 
  position: absolute;
  border-radius: 50%;
  animation: heartbeatDrift var(--d, 1.6s) ease-in-out infinite alternate;
  mix-blend-mode: screen;
  animation-delay: calc(var(--index) * 0.1s);
}

@keyframes heartbeatDrift {
  0% {
    transform: translate(-14px, -12px) scale(0.75);
    opacity: 0.6;
  }
  28% {
    transform: translate(0px, 0px) scale(1.4);
    opacity: 1;
  }
  70% {
    transform: translate(24px, 20px) scale(1.2);
    opacity: 0.85;
  }
  100% {
    transform: translate(28px, 22px) scale(1.1);
    opacity: 0.75;
  }
}


.heartbeat { 
  position: absolute;
  left: 50%;
  top: calc(50% + 10px);
  transform: translate(-50%, -50%);
  width: min(280px, 45vw);
  height: min(280px, 45vw);
  border-radius: 50%;
  background: radial-gradient(
    circle at 30% 30%, 
    rgba(255, 255, 255, 0.15) 0%, 
    rgba(255, 255, 255, 0.05) 40%, 
    transparent 70%
  );
  filter: blur(40px);
  mix-blend-mode: screen;
  animation: heartbeatPulse 1.2s ease-in-out infinite;
  z-index: 3;
}

@keyframes heartbeatPulse {
  0%, 100% { 
    transform: translate(-50%, -50%) scale(0.9); 
    opacity: 0.6; 
  }
  25% { 
    transform: translate(-50%, -50%) scale(0.75); 
    opacity: 0.9; 
  }
  55% { 
    transform: translate(-50%, -50%) scale(1.15); 
    opacity: 0.7; 
  }
}


// 主容器
.login-container {
  position: relative;
  z-index: 10;
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  padding: 0;
}

// 左侧欢迎区域
.welcome-section {
  flex: 0.6;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
}


// 标语容器
.slogan-container {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

// 标语项目
.slogan-item {
  position: absolute;
  font-size: 25px;
  font-weight: 800;
  font-style: italic;
  letter-spacing: 0.5em;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', 'Arial', sans-serif;
  background: linear-gradient(135deg, #ffffffbd 0%, #f0f0f000 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  white-space: nowrap;
  letter-spacing: 1px;
  transition: all 0.3s ease;
  animation: floatSlogan 4s ease-in-out infinite;
  
  &:hover {
    background: linear-gradient(135deg, #ffffffaf 0%, #ffffffaf 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    transform: scale(1.1);
    text-shadow: 0 0px 8px rgba(255, 247, 243, 0.89);
    animation-play-state: paused;
  }
}

// 标语浮动动画
@keyframes floatSlogan {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-8px);
  }
}

// 各个标语的弧形定位（左侧半圆弧形）
.slogan-1 {
  top: 8%;
  left: 38%;
}

.slogan-2 {
  font-size: 29px;

  top: 25%;
  left: 22%;
}

.slogan-3 {
  font-size: 40px;

  top: 44%;
  left: 13%;
}

.slogan-4 {
  font-size: 26px;

  top: 67%;
  left: 20%;
}

.slogan-5 {
  font-size: 30px;

  top: 83%;
  left: 35%;
}

// 右侧表单区域
.form-section {
  flex: 0.5;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

// 登录表单
.login-form {
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(80px);
  border-radius: 0;
  padding: 60px;
  box-shadow: none;
  border: none;
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: none;
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
  
  h2 {
    font-size: 31px;
    font-weight: 700;
    background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
    margin: 0 0 8px 0;
  }
  
  p {
    font-size: 15px;
    color: #a0aec0;
    margin: 0;
  }
}

// 输入组
.input-group {
  position: relative;
  margin-bottom: 20px;
}

.form-input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid rgba(226, 232, 240, 0.8);
  border-radius: 12px;
  font-size: 16px;
  color: #8a9ba8;
  background: rgba(255, 255, 255, 0.9);
  transition: all 0.3s ease;
  box-sizing: border-box;
  cursor: text;
  
  &:focus {
    outline: none;
    border-color: #cacaca;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    background: rgba(255, 255, 255, 1);
  }
  
  &::placeholder {
    color: #a0aec0;
  }
}

.password-toggle {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #a0aec0;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: color 0.2s ease;
  
  &:hover {
    color: #757575;
  }
}

// 表单选项
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #7c8794;
  cursor: pointer;
  
  input {
    width: 16px;
    height: 16px;
    cursor: pointer;
    
    // 移除 accent-color，使用自定义样式
    appearance: none;
    -webkit-appearance: none;
    border: 1px solid #c8ced6;
    border-radius: 2px;
    position: relative;
    
    &:checked {
      background-color: #dde4ec; // 选中时的背景色
      border-color: #a9aeb4;    // 选中时的边框色
      
      // 使用伪元素创建白色勾选标记
      &::after {
        content: "✓";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: white;
        font-size: 12px;
        font-weight: bold;
      }
    }
  }
}

.forgot-password {
  font-size: 14px;
  color: #7c8794;
  text-decoration: none;
  transition: color 0.2s ease;
  
  &:hover {
    color: #494949;
  }
}

// 按钮
.login-btn {
  width: 92%;
  padding: 16px 24px;
  border: 2px solid rgba(226, 232, 240, 0.8);
  border-radius: 25px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px auto;
  background: linear-gradient(90deg, #ffffff 0%, #ffffff 66%, #e5e7eb 100%);
  color: #8d99a9;
  position: relative;
  
  // 右侧箭头图标容器
  &::after {
    content: '';
    position: absolute;
    right: 8px;
    top: 50%;
    transform: translateY(-50%);
    width: 40px;
    height: 40px;
    background: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 0 15px rgba(255, 255, 255, 0.8);
    z-index: 2;
  }
  
  // 箭头图标 - 使用SVG
  &::before {
    content: '';
    position: absolute;
    right: 20px;
    top: 50%;
    transform: translateY(-50%);
    width: 16px;
    height: 16px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='3' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M5 12h14M12 5l7 7-7 7'/%3E%3C/svg%3E");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    z-index: 3;
  }
  
  &:hover:not(:disabled) {
    background: #ffffff;
    color: #1f29376d;
    box-shadow: 0 0 15px rgba(0, 0, 0, 0.082);
    
    &::after {
      box-shadow: 0 0 20px rgba(255, 255, 255, 1);
    }
  }
  
  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    
    &::after {
      opacity: 0.6;
    }
  }
}

// 登录提示
.login-tip {
  text-align: center;
  font-size: 14px;
  color: #718096;
  margin-top: 16px;
}

// 加载动画
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}


// 响应式设计
@media (max-width: 768px) {
  .login-container {
    flex-direction: column;
    height: auto;
    min-height: 100vh;
  }
  
  .welcome-section {
    height: auto;
    padding: 40px 20px;
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .form-section {
    height: auto;
    padding: 0;
  }
  
  .slogan-item {
    font-size: 20px;
  }
  
  .login-form {
    height: auto;
    padding: 40px 20px;
  }
}

@media (max-width: 480px) {
  .welcome-section {
    padding: 30px 15px;
  }
  
  .form-section {
    padding: 0;
  }
  
  .slogan-item {
    font-size: 18px;
  }
  
  .login-form {
    padding: 30px 15px;
  }
  
  .form-header h2 {
    font-size: 24px;
  }
}
</style>
