<template>
  <div class="chat-area">
    <!-- 通话页面 -->
    <VoiceCallPage
      v-if="isStreamingCall"
      :character="selectedChat"
      :is-listening="!isProcessingAudio"
      :is-talking="isProcessingAudio"
      :volume="currentVolume"
      @hangup="endStreamingCall"
    />
    
    <!-- 聊天页面 -->
    <div v-else-if="selectedChat" class="chat-conversation">
      <div class="chat-header">
        <div class="chat-user-info" @click="$emit('toggleProfile')">
          <div class="chat-user-avatar">
            <!-- AI伙伴使用粒子小球头像 -->
            <ParticleAvatar
              v-if="isCompanion"
              :emotion="companionEmotion.emotion"
              :intensity="companionEmotion.intensity"
              :color="companionEmotion.color"
              :brightness="companionEmotion.brightness"
              :particle-speed="companionEmotion.particle_speed"
              :growth-percentage="selectedChat.growth_percentage || 0"
              size="small"
            />
            <!-- 普通角色使用普通头像 -->
            <img v-else :src="selectedChat.avatar_url || selectedChat.avatar" :alt="selectedChat.name" />
          </div>
          <div class="chat-user-details">
            <div class="chat-user-name">{{ selectedChat.name }}</div>
            <div class="chat-user-status">{{ selectedChat.is_online ? '在线' : '离线' }}</div>
          </div>
        </div>
          <div class="chat-actions">
            <!-- AI伙伴不显示通话按钮 -->
            <button 
              v-if="!isCompanion"
              class="action-btn streaming-phone-btn" 
              @click="toggleStreamingVoiceCall" 
              :class="{ active: isStreamingCall }"
            ></button>
            <button class="action-btn more-btn"></button>
          </div>
      </div>
      
      <div class="messages-container" ref="messagesContainer">
          <template v-for="message in messages" :key="message.id">
            <ReceivedMessage 
              v-if="message.message_type !== 'user' && !message.user_message" 
              :message="message" 
              :character="selectedChat"
            />
            <SentMessage 
              v-else-if="message.message_type === 'user' || message.user_message" 
              :message="message" 
              :user-avatar="userAvatar"
            />
          </template>
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
          :disabled="!selectedChat"
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
import { ref, watch, nextTick, onMounted, onUnmounted, computed } from 'vue'
import ReceivedMessage from './ReceivedMessage.vue'
import SentMessage from './SentMessage.vue'
import VoiceCallPage from './VoiceCallPage.vue'
import ParticleAvatar from './ParticleAvatar.vue'
import chatService from '@/services/chatService.js'
import api from '@/services/api.js'

const props = defineProps({
  selectedChat: {
    type: Object,
    default: null
  }
})

defineEmits(['toggleProfile', 'showEmojiPicker', 'emotionUpdate'])

// 判断是否为AI伙伴
const isCompanion = computed(() => {
  return props.selectedChat?.name === '空白AI' || props.selectedChat?.type === 'companion'
})

// 响应式数据
const inputMessage = ref('')
const messages = ref([])
const messagesContainer = ref(null)
const userAvatar = ref('/src/img/default-avatar.png') // 默认用户头像

// AI伙伴情绪状态
const companionEmotion = ref({
  emotion: '平静',
  intensity: 0.5,
  color: '#52b4b4',
  brightness: 0.7,
  particle_speed: 0.5
})
  

// 流式语音通话相关
const isStreamingCall = ref(false)
const websocket = ref(null)
const streamingSessionId = ref('')
const isProcessingAudio = ref(false)
const currentVolume = ref(0) // 当前音量，用于通话页面显示
const vadThreshold = ref(0.3) // 语音活动检测阈值，降低到0.3让检测更敏感
const silenceTimeout = ref(1000) // 静音超时时间(毫秒) - 1秒
const audioBuffer = ref([])
const isVoiceActive = ref(false)
const voiceStartTime = ref(0) // 语音开始时间
const voiceError = ref('') // 语音错误信息
const audioContext = ref(null)
const analyser = ref(null)
const processor = ref(null)
const source = ref(null)
const microphone = ref(null)
const silenceTimer = ref(null)

// 监听选中聊天变化
watch(() => props.selectedChat, async (newChat) => {
  if (newChat) {
    // 异步加载消息，不阻塞UI
    loadMessages()
    // 如果是AI伙伴，加载情绪状态
    if (isCompanion.value) {
      loadCompanionEmotion()
    }
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

// 加载AI伙伴情绪状态
const loadCompanionEmotion = async () => {
  if (!isCompanion.value || !props.selectedChat?.id) return
  
  try {
    const token = localStorage.getItem('token')
    const response = await api.companion.getEmotionState(token, props.selectedChat.id)
    if (response.success) {
      companionEmotion.value = response.emotion
    }
  } catch (error) {
    console.error('加载AI伙伴情绪状态失败:', error)
  }
}

// 更新AI伙伴情绪状态
const updateCompanionEmotion = async (userMessage) => {
  if (!isCompanion.value || !props.selectedChat?.id) return
  
  try {
    // 简单的情绪分析（基于关键词）
    const emotion = analyzeEmotion(userMessage)
    
    // 更新本地情绪状态
    companionEmotion.value = {
      emotion: emotion,
      intensity: 0.7, // 增加强度
      color: getEmotionColor(emotion),
      brightness: 0.8,
      particle_speed: 0.6
    }
    
    console.log('AI伙伴情绪更新为:', emotion, '基于消息:', userMessage)
    
    // 通知父组件情绪状态更新
    emit('emotionUpdate', companionEmotion.value)
  } catch (error) {
    console.error('更新AI伙伴情绪状态失败:', error)
  }
}

// 分析消息情绪
const analyzeEmotion = (message) => {
  const lowerMessage = message.toLowerCase()
  
  // 开心情绪关键词
  const happyWords = ['开心', '高兴', '快乐', '哈哈', '😊', '😄', '😁', '好', '棒', '赞', '喜欢', '爱']
  if (happyWords.some(word => lowerMessage.includes(word))) {
    return '开心'
  }
  
  // 好奇情绪关键词
  const curiousWords = ['什么', '为什么', '怎么', '如何', '?', '？', '好奇', '想知道', '不明白']
  if (curiousWords.some(word => lowerMessage.includes(word))) {
    return '好奇'
  }
  
  // 孤单情绪关键词
  const lonelyWords = ['孤单', '寂寞', '一个人', '没人', '无聊', '😢', '😔', '难过', '伤心']
  if (lonelyWords.some(word => lowerMessage.includes(word))) {
    return '孤单'
  }
  
  // 兴奋情绪关键词
  const excitedWords = ['兴奋', '激动', '太棒了', '!', '！', '哇', '厉害', 'amazing', 'awesome']
  if (excitedWords.some(word => lowerMessage.includes(word))) {
    return '兴奋'
  }
  
  // 默认情绪
  return '平静'
}

// 获取情绪对应的颜色
const getEmotionColor = (emotion) => {
  const colorMap = {
    '开心': '#ffd700', // 暖黄
    '好奇': '#00bfff', // 闪烁蓝
    '孤单': '#dda0dd', // 淡紫
    '兴奋': '#ff6347', // 番茄红
    '平静': '#52b4b4'  // 柔绿
  }
  return colorMap[emotion] || '#52b4b4'
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
      
      // 如果是AI伙伴，更新情绪状态
      if (isCompanion.value) {
        await updateCompanionEmotion(messageText)
      }
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
  console.log('🎤 语音通话完成处理:', {
    userText: voiceData.userText,
    aiText: voiceData.aiText,
    hasUserText: !!(voiceData.userText && voiceData.userText.trim() !== ''),
    hasAiText: !!(voiceData.aiText && voiceData.aiText.trim() !== '')
  })
  
  // 只有当有用户文本时才添加用户消息
  if (voiceData.userText && voiceData.userText.trim() !== '') {
    console.log('🎤 添加用户消息:', voiceData.userText)
    const userMessage = {
      id: `voice_user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      user_message: voiceData.userText,
      ai_response: '',
      timestamp: new Date().toISOString(),
      message_type: 'voice'
    }
    messages.value.push(userMessage)
    console.log('🎤 用户消息已添加到数组，当前消息数量:', messages.value.length)
  } else {
    console.log('🎤 跳过用户消息（文本为空）')
  }
  
  // 添加AI语音回复
  if (voiceData.aiText && voiceData.aiText.trim() !== '') {
    console.log('🎤 添加AI消息:', voiceData.aiText)
    const aiMessage = {
      id: `voice_ai_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      user_message: '',
      ai_response: voiceData.aiText,
      timestamp: new Date().toISOString(),
      message_type: 'voice'
    }
    messages.value.push(aiMessage)
    console.log('🎤 AI消息已添加到数组，当前消息数量:', messages.value.length)
  } else {
    console.log('🎤 跳过AI消息（文本为空）')
  }
  
  // 滚动到底部
  nextTick(() => {
    console.log('🎤 强制触发响应式更新，当前消息数组:', messages.value.map(m => ({ id: m.id, type: m.message_type, user: m.user_message, ai: m.ai_response })))
    scrollToBottom()
  })
}


// 切换流式语音通话
const toggleStreamingVoiceCall = async () => {
  if (isStreamingCall.value) {
    await endStreamingCall()
  } else {
    // 显示确认提示
    const confirmed = confirm('开始流式语音通话需要访问您的麦克风，是否继续？')
    if (confirmed) {
      await startStreamingCall()
    } else {
      // 用户取消了语音通话请求
    }
  }
}

// 发送第一次流式通话请求，让AI主动打招呼
const sendFirstStreamingCall = async () => {
  try {
    console.log('发送第一次流式通话请求，让AI主动打招呼')
    
    const userId = chatService.currentUser?.id || 2 // 获取当前用户ID
    const response = await fetch('http://localhost:8080/api/v1/streaming-voice-calls/first-call', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-User-ID': userId.toString()
      },
      body: JSON.stringify({
        character_id: props.selectedChat?.character_id || 1,
        session_id: streamingSessionId.value
      })
    })
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const result = await response.json()
    console.log('AI主动打招呼完成:', result)
    
    // 播放AI的打招呼音频
    if (result.audio_response) {
      await playAIAudio(result.audio_response)
    }
    
  } catch (err) {
    console.error('AI主动打招呼失败:', err)
    voiceError.value = 'AI主动打招呼失败: ' + err.message
  }
}

// 开始流式语音通话
const startStreamingCall = async () => {
  try {
    voiceError.value = ''
    
    // 检查浏览器兼容性
    if (!checkBrowserCompatibility()) {
      isStreamingCall.value = false
      return
    }
    
    // 检查权限状态
    const hasPermission = await checkMicrophonePermission()
    if (!hasPermission) {
      isStreamingCall.value = false
      return
    }
    
    isStreamingCall.value = true
    
    // 生成会话ID
    streamingSessionId.value = 'streaming_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
    
    // 建立WebSocket连接
    await connectWebSocket()
    
    // 让AI主动打招呼（模拟第一次通话）
    await sendFirstStreamingCall()
    
    // 开始录音
    await startStreamingRecording()
    
  } catch (err) {
    voiceError.value = '流式通话连接失败: ' + err.message
    isStreamingCall.value = false
  }
}

// 建立WebSocket连接
const connectWebSocket = async () => {
  return new Promise((resolve, reject) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    // 修复端口问题：前端3000端口，后端8080端口
    const backendHost = window.location.hostname === 'localhost' ? 'localhost:8080' : window.location.host
    // 添加token和用户ID参数
    const token = 'test-token-123' // 临时token，实际应该从localStorage或用户状态获取
    const userId = chatService.currentUser?.id || 2 // 获取当前用户ID
    const wsUrl = `${protocol}//${backendHost}/api/v1/streaming-voice-calls/ws?token=${token}&user_id=${userId}`
    
    console.log('尝试连接WebSocket:', wsUrl)
    websocket.value = new WebSocket(wsUrl)
    
    websocket.value.onopen = () => {
      console.log('WebSocket连接已建立')
      
      // 发送开始通话消息
      const startMessage = {
        type: 'start_call',
        session_id: streamingSessionId.value,
        data: {
          user_id: 1, // 这里应该从用户状态获取
          character_id: props.selectedChat?.character_id || 1
        }
      }
      
      websocket.value.send(JSON.stringify(startMessage))
      resolve()
    }
    
    websocket.value.onmessage = async (event) => {
      try {
        const message = JSON.parse(event.data)
        await handleWebSocketMessage(message)
      } catch (err) {
        console.error('解析WebSocket消息失败:', err)
      }
    }
    
    websocket.value.onerror = (error) => {
      console.error('WebSocket错误:', error)
      reject(error)
    }
    
    websocket.value.onclose = () => {
      console.log('WebSocket连接已关闭')
      isStreamingCall.value = false
    }
  })
}

// 处理WebSocket消息
const handleWebSocketMessage = async (message) => {
  switch (message.type) {
    case 'call_started':
      console.log('流式通话已开始:', message.data)
      break
      
    case 'call_stopped':
      console.log('流式通话已停止:', message.data)
      isStreamingCall.value = false
      break
      
    case 'ai_response':
      console.log('收到AI回复:', message.data)
      await handleAIResponse(message.data)
      break
      
    case 'error':
      console.error('WebSocket错误:', message.data)
      
      // 处理不同类型的错误消息
      let errorMsg = message.data.error || '未知错误'
      if (errorMsg.includes('panic recovered')) {
        errorMsg = '语音识别服务连接超时，请稍后重试'
      } else if (errorMsg.includes('ASR')) {
        errorMsg = '语音识别服务暂时不可用，请稍后重试'
      } else if (errorMsg.includes('connection')) {
        errorMsg = '网络连接异常，请检查网络后重试'
      }
      
      voiceError.value = errorMsg
      
      // 如果是ASR连接错误，自动结束流式通话
      if (message.data.error && (
        message.data.error.includes('panic recovered') ||
        message.data.error.includes('ASR') ||
        message.data.error.includes('connection')
      )) {
        console.log('检测到连接错误，自动结束流式通话')
        isStreamingCall.value = false
        if (websocket.value) {
          websocket.value.close()
          websocket.value = null
        }
      }
      break
      
    case 'pong':
      // 心跳响应
      break
      
    default:
      console.log('未知消息类型:', message.type, message.data)
  }
}

// 开始流式录音
const startStreamingRecording = async () => {
  try {
    // 请求麦克风权限并获取流
    const stream = await navigator.mediaDevices.getUserMedia({ 
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true,
        sampleRate: 16000
      }
    })
    
    // 创建音频上下文
    audioContext.value = new (window.AudioContext || window.webkitAudioContext)({
      sampleRate: 16000
    })
    
    // 创建分析器节点用于语音活动检测
    analyser.value = audioContext.value.createAnalyser()
    analyser.value.fftSize = 256
    analyser.value.smoothingTimeConstant = 0.8
    
    // 创建麦克风源
    microphone.value = audioContext.value.createMediaStreamSource(stream)
    microphone.value.connect(analyser.value)
    
    // 开始语音活动检测
    startVoiceActivityDetection()
    
    console.log('流式录音已开始')
    
  } catch (err) {
    console.error('开始流式录音失败:', err)
    voiceError.value = '录音失败: ' + err.message
    isStreamingCall.value = false
  }
}

// 语音活动检测函数（全局作用域）
const checkVoiceActivity = () => {
  if (!isStreamingCall.value || !analyser.value) {
    console.log('VAD检测停止: isStreamingCall=', isStreamingCall.value, 'analyser=', !!analyser.value)
    return
  }
  
  const bufferLength = analyser.value.frequencyBinCount
  const dataArray = new Uint8Array(bufferLength)
  
  analyser.value.getByteFrequencyData(dataArray)
  
  // 计算平均音量
  let sum = 0
  for (let i = 0; i < bufferLength; i++) {
    sum += dataArray[i]
  }
  const average = sum / bufferLength
  const normalizedVolume = average / 255
  
  // 更新当前音量，用于通话页面显示
  currentVolume.value = normalizedVolume
  
    // 检测语音活动 - 简单的阈值检测
    const wasVoiceActive = isVoiceActive.value
    isVoiceActive.value = normalizedVolume > vadThreshold.value

    // 添加调试信息
    if (normalizedVolume > 0.001) { // 只在有声音时输出
      console.log(`音量: ${normalizedVolume.toFixed(4)}, 阈值: ${vadThreshold.value}, 语音活动: ${isVoiceActive.value}, 之前语音活动: ${wasVoiceActive}, 静音计时器: ${silenceTimer.value}`)
    }
  
  if (isVoiceActive.value && !wasVoiceActive) {
    // 开始说话
    console.log('检测到语音活动')
    voiceStartTime.value = Date.now() // 记录语音开始时间
    // 清除静音计时器（如果存在）
    if (silenceTimer.value) {
      console.log('检测到语音活动，清除静音计时器，ID:', silenceTimer.value)
      clearTimeout(silenceTimer.value)
      silenceTimer.value = null
    }
    if (!isProcessingAudio.value) {
      startAudioCapture()
    }
  } else if (!isVoiceActive.value && wasVoiceActive) {
    // 检测到语音结束，启动静音计时器
    console.log('检测到语音结束，启动静音计时器')
    if (silenceTimer.value) {
      clearTimeout(silenceTimer.value)
    }
    console.log('创建静音计时器，超时时间:', silenceTimeout.value, 'ms')
    silenceTimer.value = setTimeout(() => {
      console.log('静音计时器触发，检查语音状态:', isVoiceActive.value)
      console.log('静音计时器ID:', silenceTimer.value)
      if (!isVoiceActive.value) { // 再次确认语音确实结束了
        console.log('静音超时，发送语音结束信号')
        stopAudioCapture()
        sendVoiceEndSignal()
      } else {
        console.log('静音计时器触发时仍有语音活动，不停止录音')
      }
      silenceTimer.value = null
    }, silenceTimeout.value) // 1秒静音超时
    console.log('静音计时器已创建，ID:', silenceTimer.value)
  }
  
  requestAnimationFrame(checkVoiceActivity)
}

// 语音活动检测
const startVoiceActivityDetection = () => {
  checkVoiceActivity()
}

// 开始语音检测
const startVoiceDetection = () => {
  if (!isStreamingCall.value || !audioContext.value) return
  
  console.log('开始语音检测')
  
  // 强制重新创建analyser和source
  if (analyser.value) {
    analyser.value.disconnect()
    analyser.value = null
  }
  if (source.value) {
    source.value.disconnect()
    source.value = null
  }
  
  // 重新创建analyser
  analyser.value = audioContext.value.createAnalyser()
  analyser.value.fftSize = 256
  analyser.value.smoothingTimeConstant = 0.8
  console.log('重新创建analyser')
  
  // 重新连接到音频源
  if (microphone.value && microphone.value.mediaStream) {
    source.value = audioContext.value.createMediaStreamSource(microphone.value.mediaStream)
    source.value.connect(analyser.value)
    console.log('重新连接analyser到音频源')
  }
  
  // 重置语音活动状态
  isVoiceActive.value = false
  // 重新开始VAD检测循环
  console.log('开始VAD检测循环')
  checkVoiceActivity()
}

// 开始音频捕获
const startAudioCapture = () => {
  if (isProcessingAudio.value) {
    console.log('音频捕获已在进行中，跳过')
    return
  }
  
  isProcessingAudio.value = true
  audioBuffer.value = []
  
  console.log('开始PCM音频捕获')
  
  const stream = microphone.value.mediaStream
  audioContext.value = new (window.AudioContext || window.webkitAudioContext)({
    sampleRate: 16000
  })
  
  source.value = audioContext.value.createMediaStreamSource(stream)
  analyser.value = audioContext.value.createAnalyser()
  
  analyser.value.fftSize = 2048
  analyser.value.smoothingTimeConstant = 0.8
  
  source.value.connect(analyser.value)
  
  // 创建ScriptProcessorNode来获取PCM数据
  const bufferSize = 4096
  processor.value = audioContext.value.createScriptProcessor(bufferSize, 1, 1)
  
  processor.value.onaudioprocess = (event) => {
    // 只有在录音过程中才发送音频数据
    if (!isProcessingAudio.value) {
      return
    }
    
    const inputBuffer = event.inputBuffer
    const inputData = inputBuffer.getChannelData(0)
    
    // 转换为16位PCM
    const pcmData = new Int16Array(inputData.length)
    for (let i = 0; i < inputData.length; i++) {
      pcmData[i] = Math.max(-32768, Math.min(32767, inputData[i] * 32768))
    }
    
    // 转换为字节数组
    const byteArray = new Uint8Array(pcmData.buffer)
    
    // 发送PCM数据
    if (websocket.value && websocket.value.readyState === WebSocket.OPEN) {
      const audioMessage = {
        type: 'audio_chunk',
        session_id: streamingSessionId.value,
        data: {
          audio_data: Array.from(byteArray)
        }
      }
      websocket.value.send(JSON.stringify(audioMessage))
      console.log('发送音频数据:', byteArray.length, 'bytes')
    } else {
      console.warn('WebSocket连接不可用，无法发送音频数据')
    }
  }
  
  source.value.connect(processor.value)
  processor.value.connect(audioContext.value.destination)
  
  console.log('PCM音频捕获已开始')
}

// 停止音频捕获
const stopAudioCapture = () => {
  // 只停止音频处理，不清空analyser，保持VAD检测继续
  isProcessingAudio.value = false
  audioBuffer.value = []
  
  // 断开音频处理节点
  if (processor.value) {
    processor.value.disconnect()
    processor.value = null
  }
  if (source.value) {
    source.value.disconnect()
    source.value = null
  }
  
  console.log('PCM音频捕获已停止')
}

// 处理AI回复
const handleAIResponse = async (data) => {
  console.log('🎤 处理AI回复:', data)
  
  // 只有当有用户文本时才添加用户消息
  if (data.user_text && data.user_text.trim() !== '') {
    const userMessage = {
      id: `voice_user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      message_type: 'voice',
      user_message: data.user_text,
      ai_response: '',
      created_at: new Date().toISOString()
    }
    messages.value.push(userMessage)
    console.log('🎤 用户消息已添加')
  }
  
  // 只有当有AI文本时才添加AI消息
  if (data.ai_text && data.ai_text.trim() !== '') {
    const aiMessage = {
      id: `voice_ai_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      message_type: 'voice',
      user_message: '',
      ai_response: data.ai_text,
      created_at: new Date().toISOString()
    }
    messages.value.push(aiMessage)
    console.log('🎤 AI消息已添加')
  }
  
  // 播放AI音频回复
  if (data.audio_data && data.audio_data.length > 0) {
    console.log('🎤 开始播放AI音频，数据长度:', data.audio_data.length)
    await playAIAudio(data.audio_data)
  } else {
    console.log('🎤 没有音频数据，跳过播放')
  }
  
  console.log('🎤 AI回复处理完成')
  
  // AI回复完成后，重新开始录音检测
  if (isStreamingCall.value) {
    console.log('🎤 AI回复完成，重新开始录音检测')
    // 重新开始VAD检测
    startVoiceDetection()
  }
}

// 播放AI音频回复（流式通话专用）
const playStreamingAIAudio = (audioData) => {
  try {
    // 将字节数组转换为Blob
    const audioBlob = new Blob([new Uint8Array(audioData)], { type: 'audio/wav' })
    const audioUrl = URL.createObjectURL(audioBlob)
    
    // 创建音频元素并播放
    const audio = new Audio(audioUrl)
    audio.play().then(() => {
      console.log('流式AI音频开始播放')
    }).catch(err => {
      console.error('流式AI音频播放失败:', err)
    })
    
    // 播放完成后清理URL
    audio.onended = () => {
      URL.revokeObjectURL(audioUrl)
      console.log('流式AI音频播放完成')
    }
  } catch (err) {
    console.error('流式AI音频处理失败:', err)
  }
}

// 发送语音结束信号
const sendVoiceEndSignal = () => {
  if (websocket.value && websocket.value.readyState === WebSocket.OPEN) {
    const endMessage = {
      type: 'voice_end',
      session_id: streamingSessionId.value
    }
    websocket.value.send(JSON.stringify(endMessage))
    console.log('发送语音结束信号')
  }
}

// 处理流式音频 - 已废弃，现在直接发送PCM数据
const processStreamingAudio = async () => {
  // 这个函数不再使用，PCM数据直接通过onaudioprocess发送
  isProcessingAudio.value = false
}

// 结束流式语音通话
const endStreamingCall = async () => {
  try {
    isStreamingCall.value = false
    
    // 清理静音计时器
    if (silenceTimer.value) {
      clearTimeout(silenceTimer.value)
      silenceTimer.value = null
    }
    
    
    // 关闭音频上下文
    if (audioContext.value) {
      await audioContext.value.close()
      audioContext.value = null
    }
    
    // 停止所有音频轨道
    if (microphone.value && microphone.value.mediaStream) {
      microphone.value.mediaStream.getTracks().forEach(track => track.stop())
    }
    
    // 停止所有正在播放的AI音频
    forceStopAllAudio()
    
    // 发送停止通话消息
    if (websocket.value && websocket.value.readyState === WebSocket.OPEN) {
      const stopMessage = {
        type: 'stop_call',
        session_id: streamingSessionId.value
      }
      websocket.value.send(JSON.stringify(stopMessage))
    }
    
    // 关闭WebSocket连接
    if (websocket.value) {
      websocket.value.close()
      websocket.value = null
    }
    
    console.log('流式语音通话已结束')
    
    // 重新加载消息以显示最新的对话记录
    await loadMessages()
    
    // 刷新好友列表以更新最后消息时间
    try {
      await chatService.loadUserFriends()
      console.log('好友列表已刷新')
    } catch (error) {
      console.error('刷新好友列表失败:', error)
    }
    
  } catch (err) {
    console.error('结束流式通话失败:', err)
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









// 强制停止所有音频
const forceStopAllAudio = () => {
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
      
      // 检查是否还在流式通话中
      if (!isStreamingCall.value) {
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
      const audioUrl = URL.createObjectURL(audioBlob)
      
      // 创建音频元素并添加到DOM
      const audio = new Audio(audioUrl)
      audio.style.display = 'none' // 隐藏音频元素
      document.body.appendChild(audio) // 添加到DOM
      
      // 预加载音频
      audio.preload = 'auto'
      
      // 添加多个事件监听确保音频完全准备好
      let isReady = false
      
      audio.onloadeddata = () => {
        console.log('🎵 音频数据加载完成')
      }
      
      audio.oncanplay = () => {
        console.log('🎵 音频可以开始播放')
      }
      
      audio.oncanplaythrough = () => {
        if (isReady) return
        isReady = true
        
        // 再次检查是否还在流式通话中
        if (!isStreamingCall.value) {
          // 清理音频元素
          document.body.removeChild(audio)
          URL.revokeObjectURL(audioUrl)
          resolve()
          return
        }
        
        // 音频可以播放时才开始播放
        console.log('🎵 开始播放AI音频')
        
        // 增加缓冲时间确保音频完全准备好
        setTimeout(() => {
          // 设置音量渐变，避免突然的音量变化
          audio.volume = 0
          audio.play().then(() => {
            // 音量渐变到1
            const fadeIn = () => {
              if (audio.volume < 1) {
                audio.volume = Math.min(1, audio.volume + 0.1)
                setTimeout(fadeIn, 20)
              }
            }
            fadeIn()
          }).catch(err => {
            console.error('🎵 AI音频播放失败:', err)
            // 清理音频元素
            document.body.removeChild(audio)
            URL.revokeObjectURL(audioUrl)
            reject(err)
          })
        }, 100) // 增加到100ms延迟
      }
      
      audio.onended = () => {
        console.log('🎵 AI音频播放完成')
        // 清理音频元素
        document.body.removeChild(audio)
        URL.revokeObjectURL(audioUrl)
        resolve()
      }
      
      audio.onerror = (err) => {
        console.error('🎵 AI音频播放失败:', err)
        voiceError.value = '音频播放失败'
        // 清理音频元素
        document.body.removeChild(audio)
        URL.revokeObjectURL(audioUrl)
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
  // 初始化完成
})

// 清理资源
onUnmounted(() => {
  // 清理资源
})
</script>

<style lang="scss" scoped>
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
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
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
    border: 4px solid #e9ebed; // 普通角色保留边框

  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    cursor: pointer;
  }
  
  // 粒子小球样式调整
  .particle-avatar {
    width: 100%;
    height: 100%;
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


.streaming-phone-btn {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3Cpath d='M8 2v4'/%3E%3Cpath d='M16 2v4'/%3E%3C/svg%3E");
  
  &:hover {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2364748b' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3Cpath d='M8 2v4'/%3E%3Cpath d='M16 2v4'/%3E%3C/svg%3E");
  }
  
  &.active {
    background-color: #10b981;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z'/%3E%3Cpath d='M8 2v4'/%3E%3Cpath d='M16 2v4'/%3E%3C/svg%3E");
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
