// API服务 - 连接前端和后端
const API_BASE_URL = 'http://localhost:8080/api/v1';

// 用户认证相关
export const authAPI = {
  // 用户注册
  register: async (userData) => {
    const response = await fetch(`${API_BASE_URL}/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || '注册失败');
    }
    
    return response.json();
  },

  // 用户登录
  login: async (credentials) => {
    try {
      const response = await fetch(`${API_BASE_URL}/users/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });
      
      console.log('Login response status:', response.status);
      console.log('Login response ok:', response.ok);
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '登录失败');
      }
      
      const data = await response.json();
      console.log('Login response data:', data);
      return data;
    } catch (error) {
      console.error('Login API error:', error);
      throw error;
    }
  },

  // 获取用户信息
  getProfile: async (token) => {
    const response = await fetch(`${API_BASE_URL}/users/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },
};

// 角色相关API
export const characterAPI = {
  // 获取所有角色
  getAllCharacters: async () => {
    const response = await fetch(`${API_BASE_URL}/characters`);
    return response.json();
  },

  // 搜索角色
  searchCharacters: async (keyword) => {
    const response = await fetch(`${API_BASE_URL}/characters/search?keyword=${encodeURIComponent(keyword)}`);
    return response.json();
  },

  // 获取角色详情
  getCharacter: async (id) => {
    const response = await fetch(`${API_BASE_URL}/characters/${id}`);
    return response.json();
  },
};

// 好友关系相关API
export const friendshipAPI = {
  // 获取用户好友列表
  getUserFriends: async (token) => {
    const response = await fetch(`${API_BASE_URL}/friendships`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 搜索可添加的角色
  searchAvailableCharacters: async (token, keyword) => {
    const response = await fetch(`${API_BASE_URL}/friendships/search?keyword=${encodeURIComponent(keyword)}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 添加好友
  addFriend: async (token, characterId) => {
    const response = await fetch(`${API_BASE_URL}/friendships/add`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ character_id: characterId }),
    });
    return response.json();
  },

  // 移除好友
  removeFriend: async (token, characterId) => {
    const response = await fetch(`${API_BASE_URL}/friendships/${characterId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

};

// 对话相关API
export const conversationAPI = {
  // 发送文本消息
  sendMessage: async (token, messageData) => {
    const response = await fetch(`${API_BASE_URL}/conversations/chat`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(messageData),
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }
    
    return response.json();
  },

  // 发送语音消息
  sendVoiceMessage: async (token, voiceData) => {
    const response = await fetch(`${API_BASE_URL}/conversations/voice-chat`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(voiceData),
    });
    return response.json();
  },

  // 发送图片消息
  sendImageMessage: async (token, imageData) => {
    const response = await fetch(`${API_BASE_URL}/conversations/image-chat`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(imageData),
    });
    return response.json();
  },

  // 获取对话历史
  getHistory: async (token, characterId, limit = 50) => {
    const userId = getUserIdFromToken(token);
    
    const response = await fetch(`${API_BASE_URL}/conversations/history?character_id=${characterId}&limit=${limit}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': userId,
      },
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }
    
    return response.json();
  },

  // 获取会话历史
  getSessionHistory: async (token, sessionId) => {
    const response = await fetch(`${API_BASE_URL}/conversations/sessions/${sessionId}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },
};

// 工具函数
function getUserIdFromToken(token) {
  try {
    if (!token || typeof token !== 'string') {
      console.error('Invalid token:', token);
      return null;
    }
    
    // 处理mock token的情况
    if (token === 'mock-jwt-token') {
      // 从localStorage获取用户信息
      const userStr = localStorage.getItem('user');
      if (userStr) {
        const user = JSON.parse(userStr);
        return user.id;
      }
      return null;
    }
    
    // 简单的JWT解析（实际项目中应该使用更安全的方法）
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.user_id || payload.sub;
  } catch (error) {
    console.error('Failed to parse token:', error);
    return null;
  }
}

// AI伙伴相关API
export const companionAPI = {
  // 创建AI伙伴
  createCompanion: async (token, companionData) => {
    const response = await fetch(`${API_BASE_URL}/companions`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(companionData),
    });
    return response.json();
  },

  // 获取用户的AI伙伴列表
  getUserCompanions: async (token) => {
    const response = await fetch(`${API_BASE_URL}/companions`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 获取AI伙伴详情
  getCompanion: async (token, companionId) => {
    const response = await fetch(`${API_BASE_URL}/companions/${companionId}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 更新AI伙伴信息
  updateCompanion: async (token, companionId, updateData) => {
    const response = await fetch(`${API_BASE_URL}/companions/${companionId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updateData),
    });
    return response.json();
  },

  // 获取AI伙伴成长状态
  getGrowthStatus: async (token, companionId) => {
    const response = await fetch(`${API_BASE_URL}/companions/${companionId}/growth`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 获取AI伙伴日记
  getDiary: async (token, companionId, limit = 10) => {
    const response = await fetch(`${API_BASE_URL}/companions/${companionId}/diary?limit=${limit}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },

  // 获取AI伙伴情绪状态
  getEmotionState: async (token, companionId) => {
    const response = await fetch(`${API_BASE_URL}/companions/${companionId}/emotion`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': getUserIdFromToken(token),
      },
    });
    return response.json();
  },
};

// 默认导出所有API
export default {
  auth: authAPI,
  character: characterAPI,
  friendship: friendshipAPI,
  conversation: conversationAPI,
  companion: companionAPI,
};
