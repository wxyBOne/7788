// 聊天状态管理服务
import api from './api.js';
import { reactive } from 'vue';

class ChatService {
  constructor() {
    this.currentUser = null;
    this.friends = reactive([]);
    this.currentChat = null;
    this.messages = reactive([]);
    this.searchResults = reactive([]);
    this.isSearching = false;
    this.searchKeyword = '';
    this.isAddFriendMode = false;
    
    // 初始化时检查localStorage中的登录状态
    this.initializeFromStorage();
  }

  // 从localStorage初始化状态
  initializeFromStorage() {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    
    if (token && userStr) {
      try {
        this.currentUser = JSON.parse(userStr);
        console.log('Restored user from localStorage:', this.currentUser);
      } catch (error) {
        console.error('Failed to parse user from localStorage:', error);
        this.logout(); // 清除无效数据
      }
    }
  }

  // 用户认证
  async login(credentials) {
    try {
      const response = await api.auth.login(credentials);
      if (response.success) {
        this.currentUser = response.data;
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.data));
        await this.loadUserFriends();
        return response;
      }
      throw new Error(response.error || '登录失败');
    } catch (error) {
      console.error('Login error:', error);
      // 检查是否是网络错误
      if (error.message.includes('Failed to fetch')) {
        throw new Error('网络连接失败，请检查后端服务是否启动');
      }
      // 如果是用户不存在，直接抛出错误让quickLogin处理
      if (error.message.includes('用户不存在')) {
        throw new Error('用户不存在');
      }
      // 如果是401错误，说明密码错误
      if (error.message.includes('401') || error.message.includes('邮箱或密码错误') || error.message.includes('密码错误')) {
        throw new Error('密码错误，请检查密码是否正确');
      }
      throw new Error('登录失败：' + error.message);
    }
  }

  async register(userData) {
    try {
      const response = await api.auth.register(userData);
      if (response.success) {
        this.currentUser = response.data;
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.data));
        await this.loadUserFriends();
        return response;
      }
      throw new Error(response.error || '注册失败');
    } catch (error) {
      console.error('Register error:', error);
      // 检查是否是网络错误
      if (error.message.includes('Failed to fetch')) {
        throw new Error('网络连接失败，请检查后端服务是否启动');
      }
      // 如果是400错误，可能是邮箱已存在
      if (error.message.includes('400') || error.message.includes('用户已存在') || error.message.includes('邮箱已存在')) {
        throw new Error('邮箱已存在，请直接登录');
      }
      throw new Error('注册失败：' + error.message);
    }
  }

  // 登录（自动注册）
  async quickLogin(email, password) {
    try {
      console.log('Starting quickLogin with:', { email, password: '***' });
      
      // 直接调用登录接口，后端会自动处理注册
      const response = await api.auth.login({
        email: email,
        password: password
      });
      
      console.log('API response:', response);
      
      if (response.success) {
        this.currentUser = response.data;
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.data));
        
        // 尝试加载好友列表，如果失败也不影响登录
        try {
          await this.loadUserFriends();
        } catch (error) {
          console.warn('Failed to load friends, but login successful:', error);
        }
        
        return response;
      }
      throw new Error(response.error || '登录失败');
    } catch (error) {
      console.error('Quick login error:', error);
      // 检查是否是网络错误
      if (error.message.includes('Failed to fetch')) {
        throw new Error('网络连接失败，请检查后端服务是否启动');
      }
      throw new Error('登录失败：' + error.message);
    }
  }

  // 加载用户好友列表
  async loadUserFriends() {
    try {
      const token = localStorage.getItem('token');
      const response = await api.friendship.getUserFriends(token);
      if (response.success) {
        // 清空数组并添加新数据，保持响应式
        this.friends.splice(0, this.friends.length, ...response.data);
        console.log('好友列表已更新为响应式:', this.friends);
        return response.data;
      }
      throw new Error(response.error || '加载好友列表失败');
    } catch (error) {
      console.error('Load friends error:', error);
      throw error;
    }
  }

  // 搜索好友
  async searchFriends(keyword) {
    try {
      const token = localStorage.getItem('token');
      const response = await api.friendship.searchAvailableCharacters(token, keyword);
      if (response.success) {
        // 确保response.data是数组，如果是null或undefined则使用空数组
        const data = response.data || [];
        this.searchResults.splice(0, this.searchResults.length, ...data);
        return data;
      }
      throw new Error(response.error || '搜索失败');
    } catch (error) {
      console.error('Search friends error:', error);
      throw error;
    }
  }

  // 添加好友
  async addFriend(characterId) {
    try {
      const token = localStorage.getItem('token');
      const response = await api.friendship.addFriend(token, characterId);
      if (response.success) {
        // 重新加载好友列表
        await this.loadUserFriends();
        return response;
      }
      throw new Error(response.error || '添加好友失败');
    } catch (error) {
      console.error('Add friend error:', error);
      throw error;
    }
  }

  // 加载聊天记录
  async loadMessages(characterId) {
    try {
      const token = localStorage.getItem('token');
      
          // 如果是AI伙伴（character_id = 5），需要特殊处理
          if (characterId === 5) {
        console.log('加载AI伙伴消息，characterId:', characterId);
        // 使用现有的conversation API，后端需要支持character_id = 5
        const response = await api.conversation.getHistory(token, 5);
        const rawMessages = response.data || [];
        
        // 处理消息格式
        const messages = [];
        rawMessages.forEach(record => {
          if (record.user_message && record.user_message.trim()) {
            const userMessage = {
              id: record.id + '_user',
              user_message: record.user_message,
              ai_response: '',
              message_type: 'user',
              created_at: record.created_at
            };
            messages.push(userMessage);
          }
          
          if (record.ai_response && record.ai_response.trim()) {
            const aiMessage = {
              id: record.id + '_ai',
              user_message: '',
              ai_response: record.ai_response,
              message_type: record.message_type || 'text',
              created_at: record.created_at
            };
            messages.push(aiMessage);
          }
        });
        
        this.messages.splice(0, this.messages.length, ...messages);
        return this.messages;
      }
      
      // 普通角色的处理逻辑
      const response = await api.conversation.getHistory(token, characterId);
      // 直接使用response.data，不依赖success字段
      const rawMessages = response.data || [];
      
      // 将后端返回的对话记录拆分成单独的消息
      const messages = [];
      rawMessages.forEach(record => {
        // 如果有用户消息，添加用户消息
        if (record.user_message && record.user_message.trim()) {
          const userMessage = {
            id: record.id + '_user',
            user_message: record.user_message,
            ai_response: '',
            message_type: 'user',
            created_at: record.created_at
          };
          messages.push(userMessage);
        }
        
        // 如果有AI回复，添加AI回复
        if (record.ai_response && record.ai_response.trim()) {
          const aiMessage = {
            id: record.id + '_ai',
            user_message: '',
            ai_response: record.ai_response,
            message_type: record.message_type || 'text',
            created_at: record.created_at
          };
          messages.push(aiMessage);
        }
      });
      
      this.messages.splice(0, this.messages.length, ...messages);
      return this.messages;
    } catch (error) {
      console.error('Load messages error:', error);
      throw error;
    }
  }

  // 发送消息
  async sendMessage(message, characterId) {
    try {
      console.log('发送消息:', { message, characterId });
      const token = localStorage.getItem('token');
      console.log('Token:', token);
      const response = await api.conversation.sendMessage(token, {
        character_id: characterId,
        message: message
      });
      console.log('API响应:', response);
      if (response.response) {
        // 不重新加载消息，让前端处理
        // 只重新加载好友列表（更新最后消息）
        await this.loadUserFriends();
        return response;
      }
      throw new Error('API响应格式错误');
    } catch (error) {
      console.error('Send message error:', error);
      // 提供更详细的错误信息
      if (error.message.includes('HTTP 500')) {
        throw new Error('服务器内部错误，请稍后重试');
      } else if (error.message.includes('HTTP 401')) {
        throw new Error('登录已过期，请重新登录');
      } else if (error.message.includes('HTTP 400')) {
        throw new Error('请求参数错误');
      } else {
        throw new Error(`发送消息失败: ${error.message}`);
      }
    }
  }

  // 切换聊天
  async switchChat(friend) {
    this.currentChat = friend;
    await this.loadMessages(friend.character_id);
  }

  // 工具函数
  truncateText(text, maxLength = 20) {
    if (!text) return '';
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
  }

  formatTime(timestamp) {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) { // 1分钟内
      return '刚刚';
    } else if (diff < 3600000) { // 1小时内
      return Math.floor(diff / 60000) + '分钟前';
    } else if (diff < 86400000) { // 24小时内
      return Math.floor(diff / 3600000) + '小时前';
    } else {
      return date.toLocaleDateString();
    }
  }

  // 获取用户头像
  getUserAvatar() {
    return this.currentUser?.avatar_url || '/src/img/DefaultUserAvatar.jpg';
  }

  // 获取完整的角色信息（包含skills字段）
  async getFullCharacterData(characterId) {
    try {
      const response = await api.character.getCharacter(characterId);
      return response.character;
    } catch (error) {
      console.error('获取完整角色信息失败:', error);
      throw error;
    }
  }

  // 登出
  logout() {
    this.currentUser = null;
    this.friends.splice(0);
    this.currentChat = null;
    this.messages.splice(0);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  // 工具方法
  truncateMessage(message, maxLength = 30) {
    if (!message) return '';
    if (message.length <= maxLength) return message;
    return message.substring(0, maxLength) + '...';
  }

  formatTime(timestamp) {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    // 如果是今天
    if (date.toDateString() === now.toDateString()) {
      return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
    }
    
    // 如果是昨天
    const yesterday = new Date(now);
    yesterday.setDate(yesterday.getDate() - 1);
    if (date.toDateString() === yesterday.toDateString()) {
      return '昨天';
    }
    
    // 如果是本周
    const weekAgo = new Date(now);
    weekAgo.setDate(weekAgo.getDate() - 7);
    if (date > weekAgo) {
      const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
      return days[date.getDay()];
    }
    
    // 更早的日期
    return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' });
  }

  // 强制清除登录状态
  forceLogout() {
    this.currentUser = null;
    this.friends.splice(0);
    this.currentChat = null;
    this.messages.splice(0);
    this.searchResults.splice(0);
    this.isSearching = false;
    this.searchKeyword = '';
    this.isAddFriendMode = false;
    
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    console.log('Force logout completed');
  }
}

export default new ChatService();
