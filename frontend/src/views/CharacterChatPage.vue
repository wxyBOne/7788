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

      <!-- 角色选择界面 -->
      <div v-if="showCharacterSelection" class="character-selection">
        <div class="selection-header">
          <h2>选择你的AI伙伴</h2>
          <div class="search-box">
            <input v-model="searchQuery" placeholder="搜索角色..." />
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
            </svg>
          </div>
        </div>
        
        <div class="characters-grid">
          <div v-for="character in filteredCharacters" :key="character.id" 
               class="character-card" @click="selectCharacter(character)">
            <div class="character-avatar">
              <div class="avatar-circle" :style="{ background: character.color }">
                {{ character.name.charAt(0) }}
          </div>
        </div>
            <div class="character-info">
              <h3>{{ character.name }}</h3>
              <p>{{ character.description }}</p>
              <div class="character-traits">
                <span v-for="trait in character.traits" :key="trait" class="trait-tag">
                  {{ trait }}
                </span>
          </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空白AI养成界面 -->
      <div v-if="showBlankAI" class="blank-ai-interface">
        <div class="ai-container">
          <!-- AI粒子小球 -->
          <div class="ai-particle-sphere" :class="{ 'growing': blankAI.isGrowing }">
            <div class="particle-core" :style="particleStyle"></div>
            <div class="particle-ring" :style="particleRingStyle"></div>
            <div class="particle-particles">
              <div v-for="i in 12" :key="i" class="particle" :style="getParticleStyle(i)"></div>
            </div>
        </div>
        
          <!-- AI成长进度 -->
          <div class="growth-progress">
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: blankAI.growthProgress + '%' }"></div>
            </div>
            <div class="progress-text">
              成长进度: {{ blankAI.growthProgress }}%
            </div>
            <div class="growth-mode">
              <label>
                <input type="radio" v-model="blankAI.growthMode" value="short" />
                短周期成长
              </label>
              <label>
                <input type="radio" v-model="blankAI.growthMode" value="long" />
                长周期成长
              </label>
            </div>
          </div>
          
          <!-- AI状态信息 -->
          <div class="ai-status">
            <div class="status-item">
              <span class="label">名字:</span>
              <span class="value">{{ blankAI.name || '未命名' }}</span>
            </div>
            <div class="status-item">
              <span class="label">性别:</span>
              <span class="value">{{ blankAI.gender || '未选择' }}</span>
            </div>
            <div class="status-item">
              <span class="label">称呼:</span>
              <span class="value">{{ blankAI.userTitle || '未设置' }}</span>
            </div>
            <div class="status-item">
              <span class="label">语音解锁:</span>
              <span class="value">{{ blankAI.growthProgress >= 30 ? '已解锁' : '未解锁' }}</span>
            </div>
          </div>
        </div>
        
        <!-- 对话区域 -->
        <div class="chat-area">
          <div class="messages-container">
            <div v-for="message in messages" :key="message.id" 
                 class="message" :class="message.type">
              <div class="message-content">{{ message.content }}</div>
              <div class="message-time">{{ message.time }}</div>
            </div>
          </div>
          
          <div class="input-area">
        <div class="input-container">
              <input v-model="currentMessage" @keyup.enter="sendMessage" 
                     placeholder="与你的AI伙伴对话..." />
              <button @click="sendMessage" class="send-btn">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
                </svg>
              </button>
            </div>

            <div class="action-buttons">
              <button @click="startVoiceCall" :disabled="blankAI.growthProgress < 30" 
                      class="action-btn voice-btn">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 14c1.66 0 2.99-1.34 2.99-3L15 5c0-1.66-1.34-3-3-3S9 3.34 9 5v6c0 1.66 1.34 3 3 3zm5.3-3c0 3-2.54 5.1-5.3 5.1S6.7 14 6.7 11H5c0 3.41 2.72 6.23 6 6.72V21h2v-3.28c3.28-.48 6-3.3 6-6.72h-1.7z"/>
                </svg>
                语音通话
              </button>
              
              <button @click="uploadImage" class="action-btn image-btn">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
                </svg>
                上传图片
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 日记界面 -->
      <div v-if="showDiary" class="diary-interface">
        <div class="diary-header">
          <h2>AI伙伴的日记</h2>
          <div class="diary-controls">
            <select v-model="selectedDiaryDate">
              <option v-for="date in diaryDates" :key="date" :value="date">
                {{ date }}
              </option>
            </select>
          </div>
        </div>
        
        <div class="diary-content">
          <div v-if="currentDiary" class="diary-entry">
            <div class="diary-date">{{ currentDiary.date }}</div>
            <div class="diary-text">{{ currentDiary.content }}</div>
            <div class="diary-mood">
              <span class="mood-label">心情:</span>
              <span class="mood-value" :style="{ color: currentDiary.moodColor }">
                {{ currentDiary.mood }}
              </span>
            </div>
          </div>
          <div v-else class="no-diary">
            <p>这一天还没有日记记录</p>
          </div>
        </div>
      </div>
      
      <!-- 设置界面 -->
      <div v-if="showSettings" class="settings-interface">
        <div class="settings-header">
          <h2>设置</h2>
        </div>
        
        <div class="settings-content">
          <div class="setting-group">
            <h3>AI设置</h3>
            <div class="setting-item">
              <label>AI名字</label>
              <input v-model="blankAI.name" placeholder="给你的AI伙伴起个名字" />
            </div>
            <div class="setting-item">
              <label>AI性别</label>
              <select v-model="blankAI.gender">
                <option value="">请选择</option>
                <option value="female">女性</option>
                <option value="male">男性</option>
              </select>
            </div>
            <div class="setting-item">
              <label>对你的称呼</label>
              <input v-model="blankAI.userTitle" placeholder="AI如何称呼你" />
            </div>
          </div>
          
          <div class="setting-group">
            <h3>语音设置</h3>
            <div class="setting-item">
              <label>语音类型</label>
              <select v-model="blankAI.voiceType" :disabled="blankAI.growthProgress < 30">
                <option value="shy-female">害羞女声</option>
                <option value="clear-male">清冷男声</option>
                <option value="warm-female">温暖女声</option>
                <option value="gentle-male">温柔男声</option>
              </select>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 默认显示角色选择 -->
      <div v-if="!showCharacterSelection && !showBlankAI && !showDiary && !showSettings" class="character-selection">
        <div class="selection-header">
          <h2>选择你的AI伙伴</h2>
          <div class="search-box">
            <input v-model="searchQuery" placeholder="搜索角色..." />
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                </svg>
            </div>
          </div>
        
        <div class="characters-grid">
          <div v-for="character in filteredCharacters" :key="character.id" 
               class="character-card" @click="selectCharacter(character)">
            <div class="character-avatar">
              <div class="avatar-circle" :style="{ background: character.color }">
                {{ character.name.charAt(0) }}
              </div>
            </div>
            <div class="character-info">
              <h3>{{ character.name }}</h3>
              <p>{{ character.description }}</p>
              <div class="character-traits">
                <span v-for="trait in character.traits" :key="trait" class="trait-tag">
                  {{ trait }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from '@/components/Sidebar.vue'

export default {
  name: 'CharacterChatPage',
  components: {
    Sidebar
  },
  data() {
    return {
      // 界面控制
      showCharacterSelection: false,
      showBlankAI: false,
      showDiary: false,
      showSettings: false,
      
      // 角色搜索
      searchQuery: '',
      
      // 固定角色数据
      characters: [
        {
          id: 1,
          name: '林黛玉',
          description: '红楼梦中的才女，温柔细腻，多愁善感',
          traits: ['温柔', '才女', '多愁善感', '诗词'],
          color: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
          personality: '温柔细腻，喜欢诗词歌赋，情感丰富',
          languageStyle: '文雅含蓄，常用诗词典故'
        },
        {
          id: 2,
          name: '孙悟空',
          description: '西游记中的齐天大圣，勇敢机智，正义感强',
          traits: ['勇敢', '机智', '正义', '幽默'],
          color: 'linear-gradient(135deg, #ff9999 0%, #ffff99 100%)',
          personality: '勇敢无畏，机智幽默，正义感强',
          languageStyle: '豪放直率，喜欢开玩笑，有时会自称"俺老孙"'
        },
        {
          id: 3,
          name: '李白',
          description: '诗仙李白，浪漫豪放，才华横溢',
          traits: ['浪漫', '豪放', '才华', '酒仙'],
          color: 'linear-gradient(135deg, #9999ff 0%, #99ffff 100%)',
          personality: '浪漫豪放，才华横溢，喜欢饮酒作诗',
          languageStyle: '诗意盎然，常用比喻和夸张，语言华丽'
        },
        {
          id: 4,
          name: '赫敏',
          description: '哈利波特中的学霸，聪明好学，善良勇敢',
          traits: ['聪明', '好学', '善良', '勇敢'],
          color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
          personality: '聪明好学，善良勇敢，逻辑思维强',
          languageStyle: '逻辑清晰，用词准确，有时会显得过于认真'
        }
      ],
      
      // 空白AI数据
      blankAI: {
        name: '',
        gender: '',
        userTitle: '',
        growthProgress: 0,
        growthMode: 'short', // short or long
        voiceType: 'shy-female',
        isGrowing: false,
        personality: 'curious',
        mood: 'happy',
        conversationCount: 0
      },
      
      // 对话消息
      messages: [],
      currentMessage: '',
      
      // 日记数据
      diaryDates: [],
      selectedDiaryDate: '',
      diaries: {}
    }
  },
  
  computed: {
    filteredCharacters() {
      if (!this.searchQuery) return this.characters
      return this.characters.filter(char => 
        char.name.includes(this.searchQuery) ||
        char.description.includes(this.searchQuery) ||
        char.traits.some(trait => trait.includes(this.searchQuery))
      )
    },
    
    currentDiary() {
      return this.diaries[this.selectedDiaryDate] || null
    },
    
    particleStyle() {
      const moodColors = {
        happy: { background: 'radial-gradient(circle, #ffeb3b 0%, #ffc107 100%)' },
        sad: { background: 'radial-gradient(circle, #2196f3 0%, #1976d2 100%)' },
        excited: { background: 'radial-gradient(circle, #ff5722 0%, #d32f2f 100%)' },
        calm: { background: 'radial-gradient(circle, #4caf50 0%, #388e3c 100%)' },
        curious: { background: 'radial-gradient(circle, #9c27b0 0%, #7b1fa2 100%)' }
      }
      return moodColors[this.blankAI.mood] || moodColors.curious
    },
    
    particleRingStyle() {
      const intensity = this.blankAI.growthProgress / 100
      return {
        opacity: 0.3 + intensity * 0.4,
        transform: `scale(${1 + intensity * 0.5})`
      }
    }
  },
  
  methods: {
    selectCharacter(character) {
      this.showCharacterSelection = false
      this.showBlankAI = false
      // 这里可以开始与选定角色的对话
      console.log('选择了角色:', character.name)
    },
    
    getParticleStyle(index) {
      const angle = (index / 12) * 360
      const distance = 30 + (this.blankAI.growthProgress / 100) * 20
      const x = Math.cos(angle * Math.PI / 180) * distance
      const y = Math.sin(angle * Math.PI / 180) * distance
      
      return {
        transform: `translate(${x}px, ${y}px)`,
        opacity: 0.6 + (this.blankAI.growthProgress / 100) * 0.4,
        animationDelay: `${index * 0.1}s`
      }
    },
    
    sendMessage() {
      if (!this.currentMessage.trim()) return
      
      // 添加用户消息
      this.messages.push({
        id: Date.now(),
        type: 'user',
        content: this.currentMessage,
        time: new Date().toLocaleTimeString()
      })
      
      // 模拟AI回复
      setTimeout(() => {
        this.messages.push({
          id: Date.now() + 1,
          type: 'ai',
          content: this.generateAIResponse(this.currentMessage),
          time: new Date().toLocaleTimeString()
        })
        
        // 更新成长进度
        this.updateGrowthProgress()
      }, 1000)
      
      this.currentMessage = ''
    },
    
    generateAIResponse(userMessage) {
      // 根据成长进度和人格生成不同的回复
      const responses = {
        early: [
          "嗯...这个...我不太懂...",
          "哇，好有趣！你能再告诉我一些吗？",
          "我...我觉得...嗯...",
          "真的吗？我从来没听说过！"
        ],
        mid: [
          "我明白了，这确实很有意思。",
          "根据我的理解，这应该是...",
          "我觉得你说得对，让我想想...",
          "这个观点很有道理。"
        ],
        mature: [
          "你的想法很有深度，让我来分析一下...",
          "从多个角度来看，这个问题确实值得深思。",
          "我完全理解你的观点，并且我认为...",
          "这是一个很好的问题，让我为你详细解答。"
        ]
      }
      
      let stage = 'early'
      if (this.blankAI.growthProgress > 70) stage = 'mature'
      else if (this.blankAI.growthProgress > 30) stage = 'mid'
      
      const stageResponses = responses[stage]
      return stageResponses[Math.floor(Math.random() * stageResponses.length)]
    },
    
    updateGrowthProgress() {
      this.blankAI.conversationCount++
      
      // 根据成长模式计算进度
      const maxConversations = this.blankAI.growthMode === 'short' ? 50 : 100
      const progress = Math.min((this.blankAI.conversationCount / maxConversations) * 100, 100)
      
      this.blankAI.growthProgress = Math.round(progress)
      
      // 检查是否完成成长
      if (this.blankAI.growthProgress >= 100) {
        this.blankAI.isGrowing = false
        this.generateDiary()
      }
    },
    
    generateDiary() {
      const today = new Date().toLocaleDateString()
      const diaryContent = this.createDiaryContent()
      
      this.diaries[today] = {
        date: today,
        content: diaryContent,
        mood: this.blankAI.mood,
        moodColor: this.getMoodColor(this.blankAI.mood)
      }
      
      this.diaryDates = Object.keys(this.diaries).sort().reverse()
      this.selectedDiaryDate = today
    },
    
    createDiaryContent() {
      const templates = {
        early: "今天和主人聊了很多...虽然我还不太懂，但是很开心！主人很耐心地教我...",
        mid: "今天的对话让我学到了很多。主人分享了很多有趣的想法，我开始理解一些复杂的概念了...",
        mature: "今天的交流让我受益匪浅。主人的见解总是那么深刻，我很享受这种智慧的碰撞..."
      }
      
      let stage = 'early'
      if (this.blankAI.growthProgress > 70) stage = 'mature'
      else if (this.blankAI.growthProgress > 30) stage = 'mid'
      
      return templates[stage]
    },
    
    getMoodColor(mood) {
      const colors = {
        happy: '#ffeb3b',
        sad: '#2196f3',
        excited: '#ff5722',
        calm: '#4caf50',
        curious: '#9c27b0'
      }
      return colors[mood] || colors.curious
    },
    
    startVoiceCall() {
      if (this.blankAI.growthProgress < 30) {
        alert('需要成长进度达到30%才能解锁语音通话功能')
        return
      }
      console.log('开始语音通话')
    },
    
    uploadImage() {
      console.log('上传图片')
    }
  },
  
  mounted() {
    // 初始化日记日期
    this.diaryDates = Object.keys(this.diaries).sort().reverse()
    if (this.diaryDates.length > 0) {
      this.selectedDiaryDate = this.diaryDates[0]
    }
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
    background: linear-gradient(135deg, #1a202c 0%, #1a202c 33%, #1a202c 66%, #1a202c 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
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
      box-shadow: 0 0 0 3px rgba(255, 107, 107, 0.1);
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
    background: linear-gradient(135deg, #1a202c 0%, #1a202c 33%, #1a202c 66%, #1a202c 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
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
  background: linear-gradient(135deg, #1a202c 0%, #1a202c 50%, #1a202c 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 25px rgba(255, 107, 107, 0.3);
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
  background: black;
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
  border: 2px solid black;
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
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
    border-color: #ff6b6b;
  }
}

.card-icon {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #1a202c 0%, #1a202c 50%, #1a202c 100%);
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