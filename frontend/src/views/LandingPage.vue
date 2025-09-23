<template>
  <div class="landing-page" @click="handleClick">
    <!-- 基础层：淡色渐变背景 -->
    <div class="base-layer"></div>
    
    <!-- 背板对比层：冷暖色调与心跳同步 -->
    <div class="backdrop-contrast" :style="contrastStyle"></div>

    <!-- 情绪光影层：心脏跳动效果 -->
    <div class="emotion-layer" :class="{ transitioning: isTransitioning }">
      <div class="light-field">
        <div class="cluster-mask">
          <div class="gradient-waves"></div>
          <div class="swirl">
            <div class="beam beam-x"></div>
            <div class="beam beam-y"></div>
            <div class="orbit orbit-a"></div>
            <div class="orbit orbit-b"></div>
            <div v-for="i in 16" :key="i" class="light-node" :style="getNodeStyle(i)"></div>
          </div>
        </div>
      </div>
      <div class="heartbeat">
        <div class="ring"></div>
        <div class="ring d1"></div>
        <div class="ring d2"></div>
      </div>
    </div>

    <!-- 长虹玻璃层 -->
    <div class="frosted-layer" :class="{ transitioning: isTransitioning }">
      <div class="glass-cover">
        <div class="prism-effect"></div>
        <div class="vertical-stripes"></div>
        <div class="grain"></div>
        <div class="sheen"></div>
      </div>
    </div>

    <!-- 标题层 -->
    <div class="title-layer" :class="{ transitioning: isTransitioning }">
      <h1 class="brand">Seven</h1>
      <p class="tagline">AI 角色扮演 · 养成伙伴</p>
    </div>

    <!-- 白色过渡层 -->
    <div class="transition-overlay" :class="{ active: isTransitioning }"></div>

    <!-- SVG 滤镜：玻璃折射效果 -->
    <svg class="filters" width="0" height="0" aria-hidden="true" focusable="false">
      <defs>
        <filter id="glass-displace" x="-10%" y="-10%" width="120%" height="120%">
          <feTurbulence type="fractalNoise" baseFrequency="0.015" numOctaves="2" seed="3" result="noise"/>
          <feDisplacementMap in="SourceGraphic" in2="noise" scale="6" xChannelSelector="R" yChannelSelector="G"/>
        </filter>
      </defs>
    </svg>
  </div>
</template>

<script>
export default {
  name: 'LandingPage',
  data() {
    return {
      animationFrame: null,
      heartbeatPhase: 0,
      isTransitioning: false
    }
  },
  mounted() {
    this.startAnimation()
  },
  beforeUnmount() {
    if (this.animationFrame) {
      cancelAnimationFrame(this.animationFrame)
    }
  },
  methods: {
    startAnimation() {
      const startTime = Date.now()
      
      const animate = () => {
        const elapsed = (Date.now() - startTime) / 1000
        this.heartbeatPhase = (elapsed % 0.8) / 0.8
        this.updateParticles()
        this.animationFrame = requestAnimationFrame(animate)
      }
      animate()
    },
    
    updateParticles() {
      const elements = document.querySelectorAll('.light-node')
      elements.forEach(el => {
        this.updateNodeIntensity(el)
      })
    },
    
    updateNodeIntensity(node) {
      const intensity = this.getHeartbeatIntensity()
      node.style.opacity = 0.6 + intensity * 0.4
    },
    
    getHeartbeatIntensity() {
      if (this.heartbeatPhase < 0.375) {
        return 1 - (this.heartbeatPhase / 0.375)
      } else {
        return (this.heartbeatPhase - 0.375) / 0.625
      }
    },
    
    getNodeStyle(index) {
      const hues = [205, 185, 165, 145, 325, 285, 25, 35, 210, 55, 265, 305, 85, 125, 345, 15];
      const hue = hues[index % hues.length]
      const size = 8 + (index % 8)
      
      // 心形路径参数
      const t = (index / 16) * Math.PI * 2
      const heartX = 16 * Math.sin(t) ** 3
      const heartY = 13 * Math.cos(t) - 5 * Math.cos(2 * t) - 2 * Math.cos(3 * t) - Math.cos(4 * t)
      
      // 随机抖动
      const frac = n => n - Math.floor(n)
      const rx = frac(Math.sin(index * 12.9898) * 10000)
      const ry = frac(Math.cos(index * 78.233) * 10000)
      const jitterX = (rx - 0.5) * 8
      const jitterY = (ry - 0.5) * 6
      
      const x = 50 + (heartX + jitterX) * 1.2
      const y = 50 - (heartY + jitterY) * 1.2
      
      const delay = (index % 12) * 0.15
      const dur = 1.9 + (index % 7) * 0.35
      
      return {
        left: x + '%',
        top: y + '%',
        width: size + 'px',
        height: size + 'px',
        filter: `blur(${Math.max(12, 25 - size)}px)`,
        background: `radial-gradient(circle, hsla(${hue}, 99%, 72%, .95) 0%, hsla(${hue}, 96%, 58%, .35) 46%, transparent 70%)`,
        animationDelay: `${delay}s`,
        '--d': `${dur}s`,
        '--index': index
      }
    },
    
    handleClick() {
      if (this.isTransitioning) return
      
      this.isTransitioning = true
      
      // 停止心跳动画
      if (this.animationFrame) {
        cancelAnimationFrame(this.animationFrame)
      }
      
      // 快速过渡到聊天页面
      // 延迟后跳转到首页
      setTimeout(() => {
        this.$router.push('/home')
      }, 1500)
    }
  },
  computed: {
    contrastStyle() {
      const p = this.heartbeatPhase
      const brighten = 0.9 + (p < .5 ? (p/.5) : (1 - (p-.5)/.5)) * 0.25
      const shift = Math.sin(p * Math.PI * 2) * 2
      return {
        filter: `brightness(${brighten})`,
        transform: `translateX(${shift}px)`
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.landing-page {
  width: 100vw;
  height: 100vh;
  position: relative;
  overflow: hidden;
  background: #fff;
  cursor: pointer;
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
    radial-gradient(50% 70% at 30% 45%, rgba(255,100,100,.25) 0%, rgba(255,100,100,.08) 50%, rgba(0,0,0,0) 80%),
    radial-gradient(55% 75% at 70% 55%, rgba(100,200,255,.20) 0%, rgba(100,200,255,.06) 50%, rgba(0,0,0,0) 80%),
    radial-gradient(35% 50% at 50% 50%, rgba(255,200,0,.12) 0%, rgba(255,200,0,.04) 60%, rgba(0,0,0,0) 85%);
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
  width: min(600px, 75vw);
  height: min(600px, 75vw);
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
    rgba(255, 100, 100, 0.25), 
    rgba(100, 200, 255, 0.25), 
    rgba(255, 200, 0, 0.25), 
    rgba(255, 100, 100, 0.25), 
    rgba(100, 200, 255, 0.25)
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

.beam { 
  position: absolute;
  inset: 0;
  opacity: 0.4;
  filter: blur(32px);
  mix-blend-mode: screen;
}

.beam-x { 
  background: linear-gradient(
    90deg, 
    rgba(255, 100, 100, 0.6), 
    rgba(255, 200, 0, 0) 45%, 
    rgba(100, 200, 255, 0.6)
  );
  animation: beamX 3.6s ease-in-out infinite;
}

.beam-y { 
  background: linear-gradient(
    0deg, 
    rgba(255, 200, 0, 0.5), 
    rgba(255, 100, 100, 0) 45%, 
    rgba(100, 200, 255, 0.55)
  );
  animation: beamY 3.2s ease-in-out infinite;
}

@keyframes beamX { 
  0%, 100% { transform: scaleX(0.8); opacity: 0.3; }
  30% { transform: scaleX(0.7); opacity: 0.5; }
  60% { transform: scaleX(1.1); opacity: 0.4; }
}

@keyframes beamY { 
  0%, 100% { transform: scaleY(0.75); opacity: 0.25; }
  30% { transform: scaleY(0.65); opacity: 0.45; }
  60% { transform: scaleY(1.05); opacity: 0.35; }
}

.heartbeat { 
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 240px;
  height: 240px;
  pointer-events: none;
  filter: brightness(1.15);
}

.heartbeat .ring { 
  position: absolute;
  inset: 0;
  border-radius: 50%;
  box-shadow: 0 0 70px rgba(255, 255, 255, 0.55) inset;
  animation: heartbeatPulse 1.2s cubic-bezier(.3,.01,.2,1) infinite;
}

.heartbeat .ring.d1 { animation-delay: 0.2s; }
.heartbeat .ring.d2 { animation-delay: 0.4s; }

@keyframes heartbeatPulse {
  0% { transform: scale(0.65); opacity: 0.95; }
  28% { transform: scale(0.55); opacity: 1; }
  58% { transform: scale(1.18); opacity: 0.42; }
  100% { transform: scale(1.5); opacity: 0; }
}

// 长虹玻璃层
.frosted-layer { 
  z-index: 2;
  pointer-events: none;
  position: absolute;
  inset: 0;
  transition: all 0.7s cubic-bezier(0.4, 0, 0.2, 1);
  
  &.transitioning {
    opacity: 0;
    filter: blur(15px);
  }
}

.glass-cover { 
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(42px) saturate(1.35);
  -webkit-backdrop-filter: blur(42px) saturate(1.35);
  overflow: hidden;
  box-shadow: inset 0 0 0 1px rgba(255,255,255,.22), inset 0 -2px 12px rgba(0,0,0,.06);
}

// 长虹玻璃竖条纹效果
.vertical-stripes {
  position: absolute;
  inset: 0;
  background:
    repeating-linear-gradient(
      90deg,
      rgba(255,255,255,.08) 0px,
      rgba(255,255,255,.08) 1px,
      rgba(0,0,0,.03) 4px,
      rgba(0,0,0,.03) 5px,
      rgba(255,255,255,.05) 7px,
      rgba(255,255,255,.05) 8px,
      rgba(0,0,0,0) 12px
    ),
    repeating-linear-gradient(
      90deg,
      rgba(255,255,255,.04) 0px,
      rgba(255,255,255,.04) 1px,
      rgba(0,0,0,.02) 3px,
      rgba(0,0,0,.02) 4px,
      rgba(255,255,255,.03) 6px,
      rgba(255,255,255,.03) 7px,
      rgba(0,0,0,0) 12px
    );
  opacity: .18;
  mix-blend-mode: overlay;
  animation: stripeShift 6s ease-in-out infinite;
}

@keyframes stripeShift {
  0%, 100% { transform: translateX(0px); }
  50% { transform: translateX(6px); }
}

// 棱镜折射效果
.prism-effect {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    45deg,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.1) 25%,
    rgba(255, 255, 255, 0.05) 50%,
    rgba(255, 255, 255, 0.15) 75%,
    rgba(255, 255, 255, 0) 100%
  );
  mix-blend-mode: soft-light;
  opacity: 0.3;
}

.grain { 
  position: absolute;
  inset: 0;
  background-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200"><filter id="n"><feTurbulence type="fractalNoise" baseFrequency="0.8" numOctaves="3" stitchTiles="stitch"/></filter><rect width="100%" height="100%" filter="url(%23n)" opacity="0.04"/></svg>');
  opacity: 0.4;
  mix-blend-mode: multiply;
}

.sheen { 
  position: absolute;
  inset: 0;
  background: linear-gradient(
    120deg, 
    rgba(255, 255, 255, 0.3) 0%, 
    transparent 40%
  ),
  linear-gradient(
    300deg, 
    rgba(255, 255, 255, 0.2) 10%, 
    transparent 60%
  );
  mix-blend-mode: soft-light;
  opacity: 0.5;
}

// 标题层
.title-layer { 
  z-index: 3;
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  transition: all 0.7s cubic-bezier(0.4, 0, 0.2, 1);
  
  &.transitioning {
    transform: scale(1.5);
    opacity: 0;
  }
}

.brand {
  font-size: clamp(90px, 15vw, 220px);
  font-weight: 800;
  letter-spacing: 0.2em;
  margin: 0 0 16px;
  line-height: 1;
  color: rgba(255, 255, 255, 0.363);
  opacity: 0.98;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.tagline { 
  font-size: clamp(15px, 2.5vw, 24px);
  letter-spacing: 0.3em;
  color: rgba(30, 30, 30, 0.339);
  text-transform: uppercase;
  font-weight: 300;
  padding-right: 50px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

// 白色过渡层
.transition-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: white;
  z-index: 1000;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.7s cubic-bezier(0.4, 0, 0.2, 1);
  
  &.active {
    opacity: 1;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .heartbeat { width: 180px; height: 180px; }
  .cluster-mask { width: min(480px, 85vw); height: min(480px, 85vw); }
  .brand { letter-spacing: 0.15em; }
  .tagline { letter-spacing: 0.25em; }
}
</style>