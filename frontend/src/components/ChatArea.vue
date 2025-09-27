<template>
  <div class="chat-area">
    <div v-if="selectedChat" class="chat-conversation">
      <div class="chat-header">
        <div class="chat-user-info" @click="$emit('toggleProfile')">
          <div class="chat-user-avatar">
            <img :src="selectedChat.avatar_url || selectedChat.avatar" :alt="selectedChat.name" />
          </div>
          <div class="chat-user-details">
            <div class="chat-user-name">{{ selectedChat.name }}</div>
            <div class="chat-user-status">{{ selectedChat.is_online ? '在线' : '离线' }}</div>
          </div>
        </div>
        <div class="chat-actions">
          <button class="action-btn phone-btn" @click="toggleVoiceCall"></button>
          <button class="action-btn more-btn"></button>
        </div>
      </div>
      
      <div class="messages-container" ref="messagesContainer">
        <template v-for="message in messages" :key="message.id">
          <ReceivedMessage 
            v-if="message.message_type !== 'user'" 
            :message="message" 
            :character="selectedChat"
          />
          <SentMessage 
            v-else 
            :message="message" 
          />
        </template>
      </div>
      
      <!-- 语音通话状态指示器 -->
      <div v-if="isInCall" class="voice-call-indicator">
        <div class="call-status-container">
          <div v-if="isRecording" class="recording-status">
            <div class="pulse-dot"></div>
            <span>正在录音... 请说话</span>
          </div>
          <div v-else class="listening-status">
            <div class="listening-dot"></div>
            <span>点击下方按钮开始录音</span>
          </div>
          <button @click.stop="endCall" class="end-call-btn">📞 挂断</button>
          <button @click.stop="manualStartRecording" class="manual-record-btn" v-if="!isRecording">🎤 开始录音</button>
        </div>
      </div>
      
      <!-- 通话结束提示 -->
      <div v-if="callEnded" class="call-ended-indicator">
        <div class="call-ended-icon">📞</div>
        <span>通话结束</span>
      </div>
      
      <!-- 权限请求提示 -->
      <div v-if="voiceError && voiceError.includes('正在请求')" class="permission-request-indicator">
        <div class="permission-icon">🎤</div>
        <span>{{ voiceError }}</span>
      </div>
      
      <!-- 语音错误提示 -->
      <div v-if="voiceError && !voiceError.includes('正在请求')" class="voice-error-message">
        {{ voiceError }}
        <div v-if="voiceError.includes('权限被拒绝')" class="permission-help">
          <p>💡 解决方法：</p>
          <ol>
            <li>点击地址栏左侧的锁图标</li>
            <li>选择"允许"麦克风权限</li>
            <li>刷新页面重试</li>
          </ol>
          <button @click="forceRequestPermission" class="retry-permission-btn">
            🔄 强制请求权限
          </button>
        </div>
      </div>
      
      <div class="message-input">
        <input 
          v-model="inputMessage"
          type="text" 
          placeholder="输入消息..."
          @keyup.enter="sendMessage"
        />
        <div class="input-actions">
          <button class="input-btn attach-btn"></button>
          <button class="input-btn emoji-btn" @click="$emit('showEmojiPicker')"></button>
          <button class="send-btn" @click="sendMessage" :disabled="!inputMessage.trim()"></button>
        </div>
      </div>
    </div>
    
    <div v-else class="no-chat">
      <div class="no-chat-content">
        <div class="no-chat-icon">💬</div>
        <div class="no-chat-text">选择一个聊天开始对话</div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import ReceivedMessage from './ReceivedMessage.vue'
import SentMessage from './SentMessage.vue'
import chatService from '@/services/chatService.js'
import api from '@/services/api.js'

const props = defineProps({
  selectedChat: {
    type: Object,
    default: null
  }
})

defineEmits(['toggleProfile', 'showEmojiPicker'])

// 响应式数据
const inputMessage = ref('')
const messages = ref([])
const messagesContainer = ref(null)

// 语音通话相关
const isRecording = ref(false)
const isInCall = ref(false)
const callEnded = ref(false)
const mediaRecorder = ref(null)
const audioChunks = ref([])
const aiAudioUrl = ref('')
const voiceError = ref('')
const sessionId = ref('')
const audioContext = ref(null)
const analyser = ref(null)
const microphone = ref(null)
const silenceTimer = ref(null)

// 监听选中聊天变化
watch(() => props.selectedChat, async (newChat) => {
  if (newChat) {
    await loadMessages()
  }
}, { immediate: true })

// 加载消息
const loadMessages = async () => {
  if (!props.selectedChat) return
  
  try {
    await chatService.loadMessages(props.selectedChat.character_id)
    messages.value = chatService.messages
    // 滚动到底部
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('加载消息失败:', error)
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim() || !props.selectedChat) return
  
  const messageText = inputMessage.value.trim()
  inputMessage.value = ''
  
  // 立即显示用户消息
  const userMessage = {
    id: Date.now(), // 临时ID
    user_message: messageText,
    ai_response: '',
    message_type: 'user', // 改为'user'类型
    created_at: new Date().toISOString()
  }
  messages.value.push(userMessage)
  
  // 滚动到底部显示用户消息
  await nextTick()
  scrollToBottom()
  
  try {
    const response = await chatService.sendMessage(messageText, props.selectedChat.character_id)
    // 直接添加AI回复，而不是重新加载所有消息
    if (response.response) {
      const aiMessage = {
        id: response.message_id || Date.now() + 1,
        user_message: '',
        ai_response: response.response,
        message_type: 'text',
        created_at: new Date().toISOString()
      }
      messages.value.push(aiMessage)
      
      // 滚动到底部显示AI回复
      await nextTick()
      scrollToBottom()
    }
  } catch (error) {
    console.error('发送消息失败:', error)
    // 显示用户友好的错误提示
    alert(`发送消息失败: ${error.message}`)
    // 移除失败的用户消息
    const index = messages.value.findIndex(msg => msg.id === userMessage.id)
    if (index > -1) {
      messages.value.splice(index, 1)
    }
  }
}

// 处理表情选择
const handleEmojiSelect = (emoji) => {
  inputMessage.value += emoji.code
}

// 暴露方法给父组件
defineExpose({
  handleEmojiSelect
})

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 语音通话完成处理
const onVoiceCallComplete = (voiceData) => {
  // 添加用户语音消息
  if (voiceData.userText) {
    const userMessage = {
      id: Date.now(),
      user_message: voiceData.userText,
      ai_response: '',
      message_type: 'voice',
      created_at: new Date().toISOString()
    }
    messages.value.push(userMessage)
  }
  
  // 添加AI语音回复
  if (voiceData.aiText) {
    const aiMessage = {
      id: Date.now() + 1,
      user_message: '',
      ai_response: voiceData.aiText,
      message_type: 'voice',
      created_at: new Date().toISOString()
    }
    messages.value.push(aiMessage)
  }
  
  // 滚动到底部
  nextTick(() => {
    scrollToBottom()
  })
}

// 切换语音通话
const toggleVoiceCall = async () => {
  if (isInCall.value) {
    await endCall()
  } else {
    // 显示确认提示
    const confirmed = confirm('开始语音通话需要访问您的麦克风，是否继续？')
    if (confirmed) {
      await startCall()
    } else {
      // 用户取消了语音通话请求
    }
  }
}

// 检查浏览器兼容性
const checkBrowserCompatibility = () => {
  if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
    voiceError.value = '您的浏览器不支持麦克风访问功能，请使用Chrome、Firefox或Edge浏览器'
    return false
  }
  
  // 检查是否为HTTPS或localhost
  if (location.protocol !== 'https:' && location.hostname !== 'localhost' && location.hostname !== '127.0.0.1') {
    voiceError.value = '麦克风访问需要HTTPS协议，请使用https://访问或使用localhost'
    return false
  }
  
  if (!navigator.permissions) {
    // 浏览器不支持Permissions API，将直接测试getUserMedia
  }
  
  return true
}

// 检查麦克风权限状态
const checkMicrophonePermission = async () => {
  try {
    voiceError.value = '正在请求麦克风权限，请在弹出的对话框中点击"允许"'
    
    // 直接测试getUserMedia，这是最可靠的方法
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true
        } 
      })
      
      // 立即停止流，我们只需要测试权限
      stream.getTracks().forEach(track => track.stop())
      
      voiceError.value = ''
      return true
    } catch (testError) {
      // 检查权限API状态作为参考
      try {
        const permission = await navigator.permissions.query({ name: 'microphone' })
        // 权限API状态
      } catch (permError) {
        // 无法查询权限API状态
      }
      
      if (testError.name === 'NotAllowedError') {
        voiceError.value = '麦克风权限被拒绝，请点击地址栏左侧的锁图标，选择"允许"麦克风权限'
      } else if (testError.name === 'NotFoundError') {
        voiceError.value = '未找到麦克风设备，请检查设备连接'
      } else if (testError.name === 'NotSupportedError') {
        voiceError.value = '浏览器不支持语音通话功能'
      } else {
        voiceError.value = '无法访问麦克风，请检查权限设置'
      }
      return false
    }
  } catch (error) {
    // 权限检查过程出错
    voiceError.value = '无法检查麦克风权限，请刷新页面重试'
    return false
  }
}

// 强制请求麦克风权限
const forceRequestPermission = async () => {
  voiceError.value = '正在请求麦克风权限，请在弹出的对话框中点击"允许"'
  
  // 直接调用权限检查，它会测试getUserMedia
  return await checkMicrophonePermission()
}

// 发送第一次通话请求，让AI主动打招呼
const sendFirstCallRequest = async () => {
  try {
    const response = await api.voiceCall.processVoiceCall({
      character_id: props.selectedChat.character_id,
      audio_data: '', // 第一次请求不需要音频数据
      session_id: sessionId.value,
      is_first_call: true
    })
    
    if (response.text_response && response.text_response.trim() !== '') {
      // 播放AI打招呼语音
      if (response.audio_response && response.audio_response.trim() !== '') {
        await playAIAudio(response.audio_response)
      }
      
      // 触发完成事件，显示文字消息
      onVoiceCallComplete({
        userText: '',
        aiText: response.text_response,
        aiAudio: response.audio_response || ''
      })
    }
  } catch (err) {
    voiceError.value = '通话连接失败: ' + err.message
  }
}

// 开始通话
const startCall = async () => {
  try {
    voiceError.value = ''
    
    // 检查浏览器兼容性
    if (!checkBrowserCompatibility()) {
      isInCall.value = false
      return
    }
    
    // 检查权限状态
    const hasPermission = await checkMicrophonePermission()
    if (!hasPermission) {
      isInCall.value = false
      return
    }
    
    isInCall.value = true
    
    // 请求麦克风权限并获取流
    const stream = await navigator.mediaDevices.getUserMedia({ 
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true
      } 
    })
    
    // 设置音频上下文用于语音检测
    audioContext.value = new (window.AudioContext || window.webkitAudioContext)()
    analyser.value = audioContext.value.createAnalyser()
    microphone.value = audioContext.value.createMediaStreamSource(stream)
    
    microphone.value.connect(analyser.value)
    analyser.value.fftSize = 256
    
    // 发送第一次通话请求，让AI主动打招呼
    await sendFirstCallRequest()
    
  } catch (err) {
    console.error('❌ 开始通话失败:', err)
    if (err.name === 'NotAllowedError') {
      voiceError.value = '麦克风权限被拒绝，请点击地址栏左侧的锁图标，选择"允许"麦克风权限'
    } else if (err.name === 'NotFoundError') {
      voiceError.value = '未找到麦克风设备，请检查设备连接'
    } else if (err.name === 'NotSupportedError') {
      voiceError.value = '浏览器不支持语音通话功能'
    } else {
      voiceError.value = '无法访问麦克风，请检查权限设置'
    }
    isInCall.value = false
  }
}

// 结束通话
const endCall = async () => {
  console.log('📞 开始挂断通话...')
  
  isInCall.value = false
  isRecording.value = false
  callEnded.value = true
  
  // 强制停止所有音频播放
  forceStopAllAudio()
  
  // 重置处理状态
  isProcessingRecording.value = false
  
  // 清理定时器
  if (silenceTimer.value) {
    clearTimeout(silenceTimer.value)
    silenceTimer.value = null
  }
  
  // 停止录音
  if (mediaRecorder.value && mediaRecorder.value.state === 'recording') {
    mediaRecorder.value.stop()
  }
  
  // 清理录音数据
  audioChunks.value = []
  
  // 清理音频资源
  if (microphone.value) {
    microphone.value.disconnect()
  }
  if (audioContext.value) {
    audioContext.value.close()
  }
  if (aiAudioUrl.value) {
    URL.revokeObjectURL(aiAudioUrl.value)
    aiAudioUrl.value = ''
  }
  
  // 清理MediaRecorder
  if (mediaRecorder.value) {
    mediaRecorder.value = null
  }
  
  console.log('✅ 通话已挂断')
  
  // 立即隐藏通话结束提示，返回正常聊天界面
  setTimeout(() => {
    callEnded.value = false
  }, 1000)
}

// 手动开始录音
const manualStartRecording = async () => {
  if (isRecording.value) return
  
  try {
    await startRecording()
  } catch (err) {
    voiceError.value = '录音失败，请重试'
  }
}

// 开始录音
const startRecording = async () => {
  if (isRecording.value) return
  
  try {
    isRecording.value = true
    
    // 清理之前的定时器
    if (silenceTimer.value) {
      clearTimeout(silenceTimer.value)
      silenceTimer.value = null
    }
    
    // 请求麦克风权限
    const stream = await navigator.mediaDevices.getUserMedia({ 
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true
      } 
    })
    
    // 创建MediaRecorder - 使用webm格式，更稳定
    mediaRecorder.value = new MediaRecorder(stream, {
      mimeType: 'audio/webm;codecs=opus'
    })
    
    audioChunks.value = []
    
    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }
    
    mediaRecorder.value.onstop = () => {
      // 防止重复处理
      if (!isProcessingRecording.value && isInCall.value) {
        processRecording()
      }
      // 只停止录音轨道，不影响AI音频播放
      stream.getTracks().forEach(track => {
        if (track.kind === 'audio') {
          track.stop()
        }
      })
    }
    
    mediaRecorder.value.start()
    
    // 设置静音检测定时器（5秒后自动停止录音，用户一般5秒就能说完话）
    silenceTimer.value = setTimeout(() => {
      if (isRecording.value) {
        console.log('⏰ 录音超时，自动停止')
        stopRecording()
      }
    }, 5000)
    
  } catch (err) {
    console.error('录音失败:', err)
    voiceError.value = '录音失败，请重试'
    isRecording.value = false
  }
}

// 停止录音
const stopRecording = async () => {
  if (mediaRecorder.value && mediaRecorder.value.state === 'recording') {
    console.log('🛑 停止录音')
    mediaRecorder.value.stop()
    isRecording.value = false
    
    // 清除定时器
    if (silenceTimer.value) {
      clearTimeout(silenceTimer.value)
      silenceTimer.value = null
    }
  }
}

// 处理录音数据
const processRecording = async () => {
  // 立即设置处理状态，防止重复调用
  if (isProcessingRecording.value) {
    console.log('⚠️ 正在处理录音，跳过重复请求')
    return
  }
  
  isProcessingRecording.value = true
  
  try {
    // 检查是否还在通话中
    if (!isInCall.value) {
      console.log('⚠️ 通话已结束，停止处理录音')
      return
    }
    
    // 检查是否有录音数据
    if (audioChunks.value.length === 0) {
      console.log('⚠️ 没有录音数据，跳过处理')
      return
    }
    
    // 检查MediaRecorder状态
    if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
      console.log('⚠️ MediaRecorder还在运行，跳过处理')
      return
    }
    
    const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
    
    // 转换为base64
    const base64Audio = await blobToBase64(audioBlob)
    
    // 发送到后端处理
    const response = await api.voiceCall.processVoiceCall({
      character_id: props.selectedChat.character_id,
      audio_data: base64Audio,
      session_id: sessionId.value,
      is_first_call: false
    })
    
    // 再次检查是否还在通话中
    if (!isInCall.value) {
      isProcessingRecording.value = false
      return
    }
    
    // 播放AI回复 - 播放音频并显示文字
    if (response.text_response && response.text_response.trim() !== '') {
      // 播放AI语音回复
      if (response.audio_response && response.audio_response.trim() !== '') {
        await playAIAudio(response.audio_response)
      }
      
      // 触发完成事件，显示文字消息
      onVoiceCallComplete({
        userText: response.user_text || '',
        aiText: response.text_response,
        aiAudio: response.audio_response || ''
      })
    }
    
    console.log('✅ 录音处理完成')
    
    // 清理录音数据，防止重复处理
    audioChunks.value = []
    
  } catch (err) {
    console.error('语音处理失败:', err)
    voiceError.value = '语音处理失败，请重试'
  } finally {
    isProcessingRecording.value = false
    // 确保清理录音数据
    audioChunks.value = []
  }
}

// 音频管理状态
const currentAudio = ref(null)
const isPlayingAudio = ref(false)
const audioPlayPromise = ref(null)
const isProcessingRecording = ref(false)

// 强制停止所有音频
const forceStopAllAudio = () => {
  // 停止当前音频
  if (currentAudio.value) {
    currentAudio.value.pause()
    currentAudio.value.currentTime = 0
    currentAudio.value = null
  }
  
  // 重置状态
  isPlayingAudio.value = false
  
  // 清理音频URL
  if (aiAudioUrl.value) {
    URL.revokeObjectURL(aiAudioUrl.value)
    aiAudioUrl.value = ''
  }
  
  // 停止所有Audio元素
  const audioElements = document.querySelectorAll('audio')
  audioElements.forEach(audio => {
    audio.pause()
    audio.currentTime = 0
  })
  
  // 停止所有HTMLAudioElement
  const allAudioElements = document.getElementsByTagName('audio')
  for (let i = 0; i < allAudioElements.length; i++) {
    allAudioElements[i].pause()
    allAudioElements[i].currentTime = 0
  }
}

// 工具函数：Blob转Base64
const blobToBase64 = (blob) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const base64 = reader.result.split(',')[1] // 移除data:audio/webm;base64,前缀
      resolve(base64)
    }
    reader.onerror = reject
    reader.readAsDataURL(blob)
  })
}

// 播放AI音频
const playAIAudio = (base64Audio) => {
  return new Promise((resolve, reject) => {
    try {
      // 强制停止所有音频，防止重叠
      forceStopAllAudio()
      
      // 检查是否还在通话中
      if (!isInCall.value) {
        console.log('通话已结束，不播放音频')
        resolve()
        return
      }
      
      // 检查是否有音频数据
      if (!base64Audio || base64Audio.trim() === '') {
        console.log('没有音频数据，跳过播放')
        resolve()
        return
      }
      
      console.log('🎵 准备播放AI音频')
      
      // 将base64转换为Blob URL
      const audioBlob = base64ToBlob(base64Audio, 'audio/mp3')
      aiAudioUrl.value = URL.createObjectURL(audioBlob)
      
      // 创建音频元素
      const audio = new Audio(aiAudioUrl.value)
      currentAudio.value = audio
      
      // 预加载音频
      audio.preload = 'auto'
      
      audio.oncanplaythrough = () => {
        // 再次检查是否还在通话中
        if (!isInCall.value) {
          resolve()
          return
        }
        
        // 音频可以播放时才开始播放
        isPlayingAudio.value = true
        audio.play().catch(err => {
          isPlayingAudio.value = false
          reject(err)
        })
      }
      
      audio.onended = () => {
        console.log('🎵 AI音频播放完成')
        // 清理音频URL
        URL.revokeObjectURL(aiAudioUrl.value)
        aiAudioUrl.value = ''
        isPlayingAudio.value = false
        currentAudio.value = null
        resolve()
      }
      
      audio.onerror = (err) => {
        console.error('🎵 AI音频播放失败:', err)
        voiceError.value = '音频播放失败'
        isPlayingAudio.value = false
        currentAudio.value = null
        reject(err)
      }
      
      audio.onpause = () => {
        console.log('🎵 AI音频被暂停')
      }
      
      audio.onstop = () => {
        console.log('🎵 AI音频被停止')
      }
      
      // 加载音频
      audio.load()
      
    } catch (err) {
      console.error('音频播放失败:', err)
      voiceError.value = '音频播放失败'
      isPlayingAudio.value = false
      currentAudio.value = null
      reject(err)
    }
  })
}

// 工具函数：Base64转Blob
const base64ToBlob = (base64, mimeType) => {
  const byteCharacters = atob(base64)
  const byteNumbers = new Array(byteCharacters.length)
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i)
  }
  const byteArray = new Uint8Array(byteNumbers)
  return new Blob([byteArray], { type: mimeType })
}

// 初始化
onMounted(() => {
  sessionId.value = `voice_session_${Date.now()}`
})

// 清理资源
onUnmounted(() => {
  if (isInCall.value) {
    endCall()
  }
})
</script>

<style lang="scss" scoped>
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);

  min-width: 0;
}

.chat-conversation {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.chat-user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 10px 8px 8px;
  border-radius: 15px;
  transition: all 0.2s ease;
  
  &:hover {
    background: #f8fafc;
  }
  
  .chat-user-avatar {
    cursor: pointer;
  }
  
  .chat-user-details {
    cursor: pointer;
  }
}

.chat-user-avatar {
  width: 46px;
  height: 46px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 12px;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    cursor: pointer;

  }
}

.chat-user-details {
  .chat-user-name {
    font-size: 16px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 2px;
    cursor: pointer;
  }
  
  .chat-user-status {
    font-size: 12px;
    color: #10b981;
    cursor: pointer;
  }
}

.chat-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8fafc;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background-size: 18px;
  background-repeat: no-repeat;
  background-position: center;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: #e2e8f0;
    background-size: 16px; // hover时图标缩小
  }
}

.phone-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3C/svg%3E");
  }
}

.video-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolygon points='23 7 16 12 23 17 23 7'/%3E%3Crect x='1' y='5' width='15' height='14' rx='2' ry='2'/%3E%3C/svg%3E");
}

.more-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='1'/%3E%3Ccircle cx='19' cy='12' r='1'/%3E%3Ccircle cx='5' cy='12' r='1'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='1'/%3E%3Ccircle cx='19' cy='12' r='1'/%3E%3Ccircle cx='5' cy='12' r='1'/%3E%3C/svg%3E");
  }
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: flex-start;
  
  // 确保消息容器正确布局
  width: 100%;
  box-sizing: border-box;
  
  // 隐藏滚动条但保持滚动功能
  scrollbar-width: none;
  -ms-overflow-style: none;
  
  &::-webkit-scrollbar {
    display: none;
  }
}

// 为发送的消息添加右对齐
.message.sent {
  align-self: flex-end;
  align-items: flex-end;
  
  .message-content {
    align-self: flex-end;
  }
}

.message.received {
  align-self: flex-start;
  align-items: flex-start;
  
  .message-content {
    align-self: flex-start;
  }
}


.message-image {
  margin-bottom: 8px;
  
  img {
    width: 200px;
    height: 120px;
    object-fit: cover;
    border-radius: 8px;
  }
}

.message-audio {
  display: flex;
  align-items: center;
  gap: 8px;
}

.audio-waveform {
  font-size: 16px;
}

.audio-duration {
  font-size: 12px;
  color: #64748b;
}

.message-file {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-icon {
  font-size: 20px;
}

.file-info {
  .file-name {
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 2px;
  }
  
  .file-size {
    font-size: 12px;
    color: #64748b;
  }
}

.message-input {
  display: flex;
  align-items: center;
  padding: 20px;
  border-top: 1px solid #e2e8f0;
  gap: 12px;
  
  input {
    flex: 1;
    border: none;
    outline: none;
    padding: 12px 16px;
    background: #f8fafc;
    border-radius: 24px;
    font-size: 14px;
    cursor: text;
    
    &::placeholder {
      color: #94a3b8;
    }
  }
}

.input-actions {
  display: flex;
  gap: 8px;
}

.input-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f8fafc;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background-size: 18px;
  background-repeat: no-repeat;
  background-position: center;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: #e2e8f0;
    background-size: 16px; // hover时图标缩小
  }
}

.attach-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66L9.64 16.2a2 2 0 0 1-2.83-2.83l8.49-8.49'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66L9.64 16.2a2 2 0 0 1-2.83-2.83l8.49-8.49'/%3E%3C/svg%3E");
  }
}

.emoji-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10'/%3E%3Cpath d='M8 14s1.5 2 4 2 4-2 4-2'/%3E%3Cline x1='9' y1='9' x2='9.01' y2='9'/%3E%3Cline x1='15' y1='9' x2='15.01' y2='9'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='10'/%3E%3Cpath d='M8 14s1.5 2 4 2 4-2 4-2'/%3E%3Cline x1='9' y1='9' x2='9.01' y2='9'/%3E%3Cline x1='15' y1='9' x2='15.01' y2='9'/%3E%3C/svg%3E");
  }
}

.send-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #52b4b4da;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: all 0.2s ease;
  
  &::after {
    content: '';
    position: absolute;
    width: 18px;
    height: 18px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cline x1='22' y1='2' x2='11' y2='13'/%3E%3Cpolygon points='22,2 15,22 11,13 2,9 22,2'/%3E%3C/svg%3E");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
  }
  
  &:hover {
    background: #4da6a6;
    
    &::after {
      width: 16px;
      height: 16px;
    }
  }
}

.no-chat {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.no-chat-content {
  text-align: center;
  color: #64748b;
}

.no-chat-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-chat-text {
  font-size: 16px;
}

// 语音通话相关样式
.voice-call-indicator {
  padding: 8px 16px;
  margin: 8px 16px;
  border-radius: 6px;
  background: #f0f9ff;
  border: 1px solid #0ea5e9;
}

.call-status-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.recording-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #ff6b6b;
  font-size: 14px;
  font-weight: 500;
}

.listening-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #0ea5e9;
  font-size: 14px;
  font-weight: 500;
}

.end-call-btn {
  margin-left: auto;
  padding: 4px 12px;
  background: #dc2626;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.end-call-btn:hover {
  background: #b91c1c;
}

.manual-record-btn {
  margin-left: 8px;
  padding: 4px 12px;
  background: #059669;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.manual-record-btn:hover {
  background: #047857;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  background: #ff6b6b;
  border-radius: 50%;
  animation: pulse 1.5s infinite;
}

.listening-dot {
  width: 8px;
  height: 8px;
  background: #0ea5e9;
  border-radius: 50%;
  animation: listening 2s infinite;
}

.call-ended-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #059669;
  font-size: 16px;
  font-weight: 600;
  padding: 12px 16px;
  margin: 8px 16px;
  background: #ecfdf5;
  border: 1px solid #10b981;
  border-radius: 8px;
  animation: fadeInOut 3s ease-in-out;
}

.call-ended-icon {
  font-size: 18px;
  animation: phoneHang 1s ease-in-out;
}

.permission-request-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #2563eb;
  font-size: 14px;
  font-weight: 500;
  padding: 12px 16px;
  margin: 8px 16px;
  background: #eff6ff;
  border: 1px solid #3b82f6;
  border-radius: 8px;
  animation: pulse 2s infinite;
}

.permission-icon {
  font-size: 16px;
  animation: bounce 1s infinite;
}

.voice-error-message {
  color: #ff6b6b;
  font-size: 12px;
  text-align: center;
  padding: 8px 16px;
  background: #ffe6e6;
  border-radius: 6px;
  border: 1px solid #ffcccc;
  margin: 8px 16px;
}

.permission-help {
  margin-top: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 4px;
  text-align: left;
  font-size: 11px;
}

.permission-help p {
  margin: 0 0 4px 0;
  font-weight: 600;
  color: #495057;
}

.permission-help ol {
  margin: 0;
  padding-left: 16px;
  color: #6c757d;
}

.permission-help li {
  margin-bottom: 2px;
}

.retry-permission-btn {
  margin-top: 8px;
  padding: 6px 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.retry-permission-btn:hover {
  background: #0056b3;
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.7;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes listening {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.5;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes fadeInOut {
  0% {
    opacity: 0;
    transform: translateY(-10px);
  }
  20% {
    opacity: 1;
    transform: translateY(0);
  }
  80% {
    opacity: 1;
    transform: translateY(0);
  }
  100% {
    opacity: 0;
    transform: translateY(-10px);
  }
}

@keyframes phoneHang {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(-15deg);
  }
  75% {
    transform: rotate(15deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-4px);
  }
  60% {
    transform: translateY(-2px);
  }
}
</style>
