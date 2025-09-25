<template>
  <div v-if="visible" class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>重置密码</h3>
        <button class="close-btn" @click="closeModal">×</button>
      </div>
      
      <div class="modal-body">
        <div class="input-group">
          <input
            v-model="formData.email"
            type="text"
            placeholder="请输入邮箱"
            class="form-input"
          />
          <button type="button" class="send-code-btn" @click="sendVerificationCode">发送验证码</button>
        </div>
        
        <div class="input-group">
          <input
            v-model="formData.verificationCode"
            type="text"
            placeholder="请输入验证码"
            class="form-input"
          />
        </div>
        
        <div class="input-group">
          <input
            v-model="formData.newPassword"
            type="password"
            placeholder="请输入新密码"
            class="form-input"
          />
        </div>
        
        <div class="input-group">
          <input
            v-model="formData.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            class="form-input"
          />
        </div>
      </div>
      
      <div class="modal-footer">
        <button type="button" class="reset-btn" @click="resetPassword">确认重置</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

// 表单数据
const formData = reactive({
  email: '',
  verificationCode: '',
  newPassword: '',
  confirmPassword: ''
})

// 关闭弹窗
const closeModal = () => {
  // 重置表单数据
  Object.assign(formData, {
    email: '',
    verificationCode: '',
    newPassword: '',
    confirmPassword: ''
  })
  emit('close')
}

// 发送验证码
const sendVerificationCode = async () => {
  if (!formData.email) {
    alert('请输入邮箱地址')
    return
  }
  
  try {
    console.log('发送验证码到:', formData.email)
    // TODO: 调用发送验证码的API
    alert('验证码已发送到您的邮箱')
  } catch (error) {
    console.error('发送验证码失败:', error)
    alert('发送验证码失败，请重试')
  }
}

// 重置密码
const resetPassword = async () => {
  if (!formData.email || !formData.verificationCode || 
      !formData.newPassword || !formData.confirmPassword) {
    alert('请填写完整信息')
    return
  }
  
  if (formData.newPassword !== formData.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }
  
  try {
    console.log('重置密码数据:', formData)
    // TODO: 调用重置密码的API
    alert('密码重置成功，请重新登录')
    closeModal()
  } catch (error) {
    console.error('重置密码失败:', error)
    alert('重置密码失败，请重试')
  }
}
</script>

<style lang="scss" scoped>
// 忘记密码弹窗样式
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(80px);
  border-radius: 20px;
  padding: 0;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  
  h3 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: #8d99a9; // 标题颜色调淡
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 24px;
    color: #6b7280;
    cursor: pointer;
    padding: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
    
    &:hover {
      background: rgba(0, 0, 0, 0.05);
      color: #374151;
    }
  }
}

.modal-body {
  padding: 32px;
  
  .input-group {
    margin-bottom: 20px;
    position: relative;
    display: flex;
    gap: 12px;
    
    .form-input {
      flex: 1;
      width: 100%;
      padding: 12px 16px;
      border: 1px solid #e2e8f0;
      border-radius: 8px;
      font-size: 14px;
      color: #8d99a9;
      background: rgba(255, 255, 255, 0.8);
      transition: all 0.2s ease;
      
      &::placeholder {
        color: #a0aec0;
      }
      
      &:focus {
        outline: none;
        border-color: #cacaca;
        background: rgba(255, 255, 255, 1);
      }
    }
    
    .send-code-btn {
      padding: 12px 20px;
      background: #f3f4f6;
      border: 1px solid #d1d5db;
      border-radius: 8px;
      color: #374151;
      font-size: 14px;
      cursor: pointer;
      transition: all 0.2s ease;
      white-space: nowrap;
      
      &:hover {
        background: #e5e7eb;
        border-color: #9ca3af;
      }
    }
  }
}

.modal-footer {
  padding: 24px 32px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  
  .reset-btn {
    width: 100%;
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
}

// 响应式设计
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 20px;
  }
  
  .modal-header, .modal-body, .modal-footer {
    padding: 20px;
  }
  
  .modal-body .input-group {
    flex-direction: column;
    gap: 8px;
    
    .send-code-btn {
      width: 100%;
    }
  }
}

@media (max-width: 480px) {
  .modal-header, .modal-body, .modal-footer {
    padding: 15px;
  }
  
  .modal-header h3 {
    font-size: 18px;
  }
  
  .reset-btn {
    padding: 14px 20px;
    font-size: 14px;
  }
}
</style>
