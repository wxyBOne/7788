// 消息处理相关的组合式函数
import { ref } from 'vue'
import chatService from '@/services/chatService.js'
import api from '@/services/api.js'

export function useMessageHandler() {
  // 消息相关状态
  const messages = ref([])
  const newMessage = ref('')
  const messagesContainer = ref(null)
  const userAvatar = ref('/src/img/default-avatar.png')

  // 加载消息
  const loadMessages = async () => {
    if (!chatService.selectedChat) return
    
    try {
      const response = await api.conversation.getHistory({
        character_id: chatService.selectedChat.character_id || chatService.selectedChat.id,
        user_id: chatService.currentUser?.id || 2
      })
      
      if (response.success) {
        messages.value = response.messages || []
        // 滚动到底部
        await nextTick()
        scrollToBottom()
      }
    } catch (error) {
      console.error('加载消息失败:', error)
    }
  }

  // 发送消息
  const sendMessage = async () => {
    if (!newMessage.value.trim() || !chatService.selectedChat) return

    const messageText = newMessage.value.trim()
    newMessage.value = ''

    // 立即显示用户消息
    const userMessage = {
      id: Date.now(), // 临时ID
      content: messageText,
      message_type: 'user', // 改为'user'类型
      timestamp: new Date().toISOString(),
      character_id: chatService.selectedChat.character_id || chatService.selectedChat.id
    }

    // 滚动到底部显示用户消息
    await nextTick()
    scrollToBottom()

    try {
      // 发送消息到后端
      const response = await api.conversation.sendMessage({
        character_id: chatService.selectedChat.character_id || chatService.selectedChat.id,
        user_id: chatService.currentUser?.id || 2,
        content: messageText,
        message_type: 'text'
      })

      if (response.success) {
        // 直接添加AI回复，而不是重新加载所有消息
        const aiMessage = {
          id: Date.now() + 1,
          content: response.message.content,
          message_type: 'ai',
          timestamp: response.message.timestamp,
          character_id: chatService.selectedChat.character_id || chatService.selectedChat.id
        }
        
        messages.value.push(userMessage)
        messages.value.push(aiMessage)
        
        // 滚动到底部显示AI回复
        await nextTick()
        scrollToBottom()
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      // 显示用户友好的错误提示
      alert('发送消息失败，请重试')
      // 移除失败的用户消息
      messages.value = messages.value.filter(msg => msg.id !== userMessage.id)
    }
  }

  // 处理表情选择
  const handleEmojiSelect = (emoji) => {
    newMessage.value += emoji
  }

  // 滚动到底部
  const scrollToBottom = () => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  }

  // 处理语音通话完成
  const handleVoiceCallComplete = async (data) => {
    // 只有当有用户文本时才添加用户消息
    if (data.user_text) {
      const userMessage = {
        id: Date.now(),
        content: data.user_text,
        message_type: 'user',
        timestamp: new Date().toISOString(),
        character_id: chatService.selectedChat?.character_id || chatService.selectedChat?.id
      }
      messages.value.push(userMessage)
    }

    // 添加AI语音回复
    if (data.ai_text) {
      const aiMessage = {
        id: Date.now() + 1,
        content: data.ai_text,
        message_type: 'ai',
        timestamp: new Date().toISOString(),
        character_id: chatService.selectedChat?.character_id || chatService.selectedChat?.id
      }
      messages.value.push(aiMessage)
    }

    // 滚动到底部
    await nextTick()
    scrollToBottom()
  }

  return {
    // 状态
    messages,
    newMessage,
    messagesContainer,
    userAvatar,
    
    // 方法
    loadMessages,
    sendMessage,
    handleEmojiSelect,
    scrollToBottom,
    handleVoiceCallComplete
  }
}
