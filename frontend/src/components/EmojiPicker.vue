<template>
  <div v-if="visible" class="emoji-picker-overlay" @click="handleOverlayClick">
    <div class="emoji-picker-container" @click.stop>
      <div class="emoji-picker-header">
        <h4>选择表情</h4>
        <button class="close-btn" @click="$emit('close')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <div class="emoji-picker-body">
        <div class="emoji-categories">
          <button 
            v-for="category in categories"
            :key="category.name"
            class="category-btn"
            :class="{ active: selectedCategory === category.name }"
            @click="selectedCategory = category.name"
          >
            {{ category.icon }}
          </button>
        </div>

        <div class="emoji-grid">
          <button
            v-for="emoji in filteredEmojis"
            :key="emoji.code"
            class="emoji-btn"
            @click="selectEmoji(emoji)"
            :title="emoji.name"
          >
            {{ emoji.code }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'select'])

// 表情分类
const categories = ref([
  { name: 'all', icon: '😀' },
  { name: 'smileys', icon: '😊' },
  { name: 'gestures', icon: '👍' },
  { name: 'hearts', icon: '❤️' },
  { name: 'objects', icon: '🎉' }
])

// 表情数据
const emojis = ref([
  // 笑脸类
  { code: '😀', name: '大笑', category: 'smileys' },
  { code: '😃', name: '开心', category: 'smileys' },
  { code: '😄', name: '眯眼笑', category: 'smileys' },
  { code: '😁', name: '露齿笑', category: 'smileys' },
  { code: '😆', name: '大笑', category: 'smileys' },
  { code: '😅', name: '苦笑', category: 'smileys' },
  { code: '🤣', name: '笑翻', category: 'smileys' },
  { code: '😂', name: '笑哭', category: 'smileys' },
  { code: '🙂', name: '微笑', category: 'smileys' },
  { code: '🙃', name: '倒微笑', category: 'smileys' },
  { code: '😉', name: '眨眼', category: 'smileys' },
  { code: '😊', name: '害羞笑', category: 'smileys' },
  { code: '😇', name: '天使笑', category: 'smileys' },
  { code: '🥰', name: '爱心眼', category: 'smileys' },
  { code: '😍', name: '花痴', category: 'smileys' },
  { code: '🤩', name: '星星眼', category: 'smileys' },
  { code: '😘', name: '飞吻', category: 'smileys' },
  { code: '😗', name: '亲吻', category: 'smileys' },
  { code: '😚', name: '闭眼吻', category: 'smileys' },
  { code: '😙', name: '微笑吻', category: 'smileys' },
  { code: '😋', name: '美味', category: 'smileys' },
  { code: '😛', name: '吐舌', category: 'smileys' },
  { code: '😜', name: '眨眼吐舌', category: 'smileys' },
  { code: '🤪', name: '疯狂', category: 'smileys' },
  { code: '😝', name: '闭眼吐舌', category: 'smileys' },
  { code: '🤑', name: '金钱眼', category: 'smileys' },
  { code: '🤗', name: '拥抱', category: 'smileys' },
  { code: '🤭', name: '捂嘴笑', category: 'smileys' },
  { code: '🤫', name: '嘘', category: 'smileys' },
  { code: '🤔', name: '思考', category: 'smileys' },
  { code: '🤐', name: '拉链嘴', category: 'smileys' },
  { code: '🤨', name: '挑眉', category: 'smileys' },
  { code: '😐', name: '面无表情', category: 'smileys' },
  { code: '😑', name: '无语', category: 'smileys' },
  { code: '😶', name: '无嘴', category: 'smileys' },
  { code: '😏', name: '得意', category: 'smileys' },
  { code: '😒', name: '无聊', category: 'smileys' },
  { code: '🙄', name: '翻白眼', category: 'smileys' },
  { code: '😬', name: '龇牙', category: 'smileys' },
  { code: '🤥', name: '长鼻子', category: 'smileys' },
  { code: '😔', name: '沉思', category: 'smileys' },
  { code: '😪', name: '困倦', category: 'smileys' },
  { code: '🤤', name: '流口水', category: 'smileys' },
  { code: '😴', name: '睡觉', category: 'smileys' },
  { code: '😷', name: '口罩', category: 'smileys' },
  { code: '🤒', name: '发烧', category: 'smileys' },
  { code: '🤕', name: '受伤', category: 'smileys' },
  { code: '🤢', name: '恶心', category: 'smileys' },
  { code: '🤮', name: '呕吐', category: 'smileys' },
  { code: '🤧', name: '打喷嚏', category: 'smileys' },
  { code: '🥵', name: '热', category: 'smileys' },
  { code: '🥶', name: '冷', category: 'smileys' },
  { code: '🥴', name: '眩晕', category: 'smileys' },
  { code: '😵', name: '晕倒', category: 'smileys' },
  { code: '🤯', name: '爆炸头', category: 'smileys' },
  { code: '🤠', name: '牛仔', category: 'smileys' },
  { code: '🥳', name: '派对', category: 'smileys' },
  { code: '😎', name: '墨镜', category: 'smileys' },
  { code: '🤓', name: '书呆子', category: 'smileys' },
  { code: '🧐', name: '单镜片', category: 'smileys' },

  // 手势类
  { code: '👍', name: '点赞', category: 'gestures' },
  { code: '👎', name: '踩', category: 'gestures' },
  { code: '👌', name: 'OK', category: 'gestures' },
  { code: '✌️', name: '胜利', category: 'gestures' },
  { code: '🤞', name: '交叉手指', category: 'gestures' },
  { code: '🤟', name: '爱你', category: 'gestures' },
  { code: '🤘', name: '摇滚', category: 'gestures' },
  { code: '🤙', name: '电话', category: 'gestures' },
  { code: '👈', name: '左指', category: 'gestures' },
  { code: '👉', name: '右指', category: 'gestures' },
  { code: '👆', name: '上指', category: 'gestures' },
  { code: '🖕', name: '中指', category: 'gestures' },
  { code: '👇', name: '下指', category: 'gestures' },
  { code: '☝️', name: '食指', category: 'gestures' },
  { code: '👋', name: '挥手', category: 'gestures' },
  { code: '🤚', name: '举手', category: 'gestures' },
  { code: '🖐️', name: '张开手', category: 'gestures' },
  { code: '✋', name: '停止', category: 'gestures' },
  { code: '🖖', name: '瓦肯手势', category: 'gestures' },
  { code: '👏', name: '鼓掌', category: 'gestures' },
  { code: '🙌', name: '举双手', category: 'gestures' },
  { code: '👐', name: '张开双手', category: 'gestures' },
  { code: '🤲', name: '掌心向上', category: 'gestures' },
  { code: '🤝', name: '握手', category: 'gestures' },
  { code: '🙏', name: '祈祷', category: 'gestures' },
  { code: '✍️', name: '写字', category: 'gestures' },
  { code: '💅', name: '指甲', category: 'gestures' },
  { code: '🤳', name: '自拍', category: 'gestures' },
  { code: '💪', name: '肌肉', category: 'gestures' },
  { code: '🦾', name: '机械臂', category: 'gestures' },
  { code: '🦿', name: '机械腿', category: 'gestures' },
  { code: '🦵', name: '腿', category: 'gestures' },
  { code: '🦶', name: '脚', category: 'gestures' },
  { code: '👂', name: '耳朵', category: 'gestures' },
  { code: '🦻', name: '助听器', category: 'gestures' },
  { code: '👃', name: '鼻子', category: 'gestures' },
  { code: '🧠', name: '大脑', category: 'gestures' },
  { code: '🦷', name: '牙齿', category: 'gestures' },
  { code: '🦴', name: '骨头', category: 'gestures' },
  { code: '👀', name: '眼睛', category: 'gestures' },
  { code: '👁️', name: '眼睛', category: 'gestures' },
  { code: '👅', name: '舌头', category: 'gestures' },
  { code: '👄', name: '嘴巴', category: 'gestures' },

  // 爱心类
  { code: '❤️', name: '红心', category: 'hearts' },
  { code: '🧡', name: '橙心', category: 'hearts' },
  { code: '💛', name: '黄心', category: 'hearts' },
  { code: '💚', name: '绿心', category: 'hearts' },
  { code: '💙', name: '蓝心', category: 'hearts' },
  { code: '💜', name: '紫心', category: 'hearts' },
  { code: '🖤', name: '黑心', category: 'hearts' },
  { code: '🤍', name: '白心', category: 'hearts' },
  { code: '🤎', name: '棕心', category: 'hearts' },
  { code: '💔', name: '心碎', category: 'hearts' },
  { code: '❣️', name: '感叹心', category: 'hearts' },
  { code: '💕', name: '双心', category: 'hearts' },
  { code: '💞', name: '旋转心', category: 'hearts' },
  { code: '💓', name: '跳动心', category: 'hearts' },
  { code: '💗', name: '长大心', category: 'hearts' },
  { code: '💖', name: '闪亮心', category: 'hearts' },
  { code: '💘', name: '箭穿心', category: 'hearts' },
  { code: '💝', name: '礼物心', category: 'hearts' },
  { code: '💟', name: '装饰心', category: 'hearts' },

  // 物品类
  { code: '🎉', name: '庆祝', category: 'objects' },
  { code: '🎊', name: '彩带', category: 'objects' },
  { code: '🎈', name: '气球', category: 'objects' },
  { code: '🎁', name: '礼物', category: 'objects' },
  { code: '🎀', name: '蝴蝶结', category: 'objects' },
  { code: '🎂', name: '蛋糕', category: 'objects' },
  { code: '🍰', name: '切片蛋糕', category: 'objects' },
  { code: '🧁', name: '纸杯蛋糕', category: 'objects' },
  { code: '🍭', name: '棒棒糖', category: 'objects' },
  { code: '🍬', name: '糖果', category: 'objects' },
  { code: '🍫', name: '巧克力', category: 'objects' },
  { code: '🍩', name: '甜甜圈', category: 'objects' },
  { code: '🍪', name: '饼干', category: 'objects' },
  { code: '🍯', name: '蜂蜜', category: 'objects' },
  { code: '🍮', name: '布丁', category: 'objects' },
  { code: '🍨', name: '冰淇淋', category: 'objects' },
  { code: '🍧', name: '刨冰', category: 'objects' },
  { code: '🍦', name: '软冰淇淋', category: 'objects' },
  { code: '🥧', name: '派', category: 'objects' },
  { code: '🍕', name: '披萨', category: 'objects' },
  { code: '🍔', name: '汉堡', category: 'objects' },
  { code: '🌭', name: '热狗', category: 'objects' },
  { code: '🥪', name: '三明治', category: 'objects' },
  { code: '🌮', name: '玉米饼', category: 'objects' },
  { code: '🌯', name: '卷饼', category: 'objects' },
  { code: '🥙', name: '口袋面包', category: 'objects' },
  { code: '🥚', name: '鸡蛋', category: 'objects' },
  { code: '🍳', name: '煎蛋', category: 'objects' },
  { code: '🥞', name: '薄饼', category: 'objects' },
  { code: '🧇', name: '华夫饼', category: 'objects' },
  { code: '🥓', name: '培根', category: 'objects' },
  { code: '🥩', name: '肉', category: 'objects' },
  { code: '🍗', name: '鸡腿', category: 'objects' },
  { code: '🍖', name: '骨头肉', category: 'objects' },
  { code: '🦴', name: '骨头', category: 'objects' },
  { code: '🌽', name: '玉米', category: 'objects' },
  { code: '🍅', name: '番茄', category: 'objects' },
  { code: '🍄', name: '蘑菇', category: 'objects' },
  { code: '🥕', name: '胡萝卜', category: 'objects' },
  { code: '🌶️', name: '辣椒', category: 'objects' },
  { code: '🌶', name: '辣椒', category: 'objects' },
  { code: '🫑', name: '甜椒', category: 'objects' },
  { code: '🥒', name: '黄瓜', category: 'objects' },
  { code: '🥬', name: '绿叶菜', category: 'objects' },
  { code: '🥦', name: '西兰花', category: 'objects' },
  { code: '🧄', name: '大蒜', category: 'objects' },
  { code: '🧅', name: '洋葱', category: 'objects' },
  { code: '🥜', name: '坚果', category: 'objects' },
  { code: '🌰', name: '栗子', category: 'objects' },
  { code: '🍞', name: '面包', category: 'objects' },
  { code: '🥐', name: '牛角包', category: 'objects' },
  { code: '🥖', name: '法棍', category: 'objects' },
  { code: '🍞', name: '面包', category: 'objects' },
  { code: '🥨', name: '椒盐卷饼', category: 'objects' },
  { code: '🥯', name: '贝果', category: 'objects' },
  { code: '🥞', name: '薄饼', category: 'objects' },
  { code: '🧇', name: '华夫饼', category: 'objects' },
  { code: '🧀', name: '奶酪', category: 'objects' },
  { code: '🍖', name: '骨头肉', category: 'objects' },
  { code: '🍗', name: '鸡腿', category: 'objects' },
  { code: '🥩', name: '肉', category: 'objects' },
  { code: '🥓', name: '培根', category: 'objects' },
  { code: '🍔', name: '汉堡', category: 'objects' },
  { code: '🌭', name: '热狗', category: 'objects' },
  { code: '🥪', name: '三明治', category: 'objects' },
  { code: '🌮', name: '玉米饼', category: 'objects' },
  { code: '🌯', name: '卷饼', category: 'objects' },
  { code: '🥙', name: '口袋面包', category: 'objects' },
  { code: '🥚', name: '鸡蛋', category: 'objects' },
  { code: '🍳', name: '煎蛋', category: 'objects' },
  { code: '🥞', name: '薄饼', category: 'objects' },
  { code: '🧇', name: '华夫饼', category: 'objects' },
  { code: '🥓', name: '培根', category: 'objects' },
  { code: '🥩', name: '肉', category: 'objects' },
  { code: '🍗', name: '鸡腿', category: 'objects' },
  { code: '🍖', name: '骨头肉', category: 'objects' },
  { code: '🦴', name: '骨头', category: 'objects' },
  { code: '🌽', name: '玉米', category: 'objects' },
  { code: '🍅', name: '番茄', category: 'objects' },
  { code: '🍄', name: '蘑菇', category: 'objects' },
  { code: '🥕', name: '胡萝卜', category: 'objects' },
  { code: '🌶️', name: '辣椒', category: 'objects' },
  { code: '🌶', name: '辣椒', category: 'objects' },
  { code: '🫑', name: '甜椒', category: 'objects' },
  { code: '🥒', name: '黄瓜', category: 'objects' },
  { code: '🥬', name: '绿叶菜', category: 'objects' },
  { code: '🥦', name: '西兰花', category: 'objects' },
  { code: '🧄', name: '大蒜', category: 'objects' },
  { code: '🧅', name: '洋葱', category: 'objects' },
  { code: '🥜', name: '坚果', category: 'objects' },
  { code: '🌰', name: '栗子', category: 'objects' },
  { code: '🍞', name: '面包', category: 'objects' },
  { code: '🥐', name: '牛角包', category: 'objects' },
  { code: '🥖', name: '法棍', category: 'objects' },
  { code: '🍞', name: '面包', category: 'objects' },
  { code: '🥨', name: '椒盐卷饼', category: 'objects' },
  { code: '🥯', name: '贝果', category: 'objects' },
  { code: '🥞', name: '薄饼', category: 'objects' },
  { code: '🧇', name: '华夫饼', category: 'objects' },
  { code: '🧀', name: '奶酪', category: 'objects' }
])

const selectedCategory = ref('all')

// 计算过滤后的表情
const filteredEmojis = computed(() => {
  if (selectedCategory.value === 'all') {
    return emojis.value
  }
  return emojis.value.filter(emoji => emoji.category === selectedCategory.value)
})

// 选择表情
const selectEmoji = (emoji) => {
  emit('select', emoji)
  emit('close')
}

// 点击遮罩层关闭
const handleOverlayClick = () => {
  emit('close')
}
</script>

<style scoped>
.emoji-picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99999;
}

.emoji-picker-container {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  max-height: 500px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.emoji-picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.emoji-picker-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  color: #6b7280;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

.emoji-picker-body {
  display: flex;
  flex-direction: column;
  height: 400px;
}

.emoji-categories {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
  gap: 8px;
}

.category-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  font-size: 18px;
  transition: all 0.2s;
}

.category-btn:hover {
  background: #e5e7eb;
}

.category-btn.active {
  background: #3b82f6;
  color: white;
}

.emoji-grid {
  flex: 1;
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.emoji-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  font-size: 20px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-btn:hover {
  background: #f3f4f6;
  transform: scale(1.1);
}

.emoji-btn:active {
  transform: scale(0.95);
}
</style>
