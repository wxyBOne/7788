<template>
  <div class="sidebar">
      <div class="logo-section">
        <div class="logo">seven</div>
      </div>
    
    <nav class="nav-menu">
      <div class="nav-item" :class="{ active: activeSection === 'home' }" @click="$emit('setActiveSection', 'home')">
        <div class="nav-icon home-icon"></div>
      </div>
      <div class="nav-item" :class="{ active: activeSection === 'messages' }" @click="$emit('setActiveSection', 'messages')">
        <div class="nav-icon messages-icon"></div>
      </div>
      <div class="nav-item" :class="{ active: activeSection === 'settings' }" @click="$emit('setActiveSection', 'settings')">
        <div class="nav-icon settings-icon"></div>
      </div>
    </nav>
    
    <div class="user-section">
      <div class="user-avatar">
        <img src="/src/img/DefaultUserAvatar.jpg" alt="User" />
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  activeSection: {
    type: String,
    default: 'messages'
  }
})

defineEmits(['setActiveSection'])
</script>

<style lang="scss" scoped>
.sidebar {
  width: 60px;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
  border-right: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40px;
}

.logo {
  color: white;
  background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
  font-weight: bold;
  font-size: 16px;
  padding: 10px 14px 14px;
  border-radius: 12px;
  margin-bottom: 8px;
  width: fit-content;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 88;
  box-shadow: 0 4px 12px rgba(49, 49, 49, 0.14);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(255, 132, 132, 0.4);
  }
}

.logo-text {
  text-align: center;
  font-size: 10px;
  color: #64748b;
}

.logo-title {
  font-weight: 600;
  margin-bottom: 2px;
}

.logo-subtitle {
  font-size: 8px;
}

.nav-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #64748b;
  margin-bottom: 2px;
  width: 40px;
  
  &:hover {
    background: #e2e8f0;
  }
  
  &.active {
    color: transparent;
    background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    
    .nav-icon {
      background: linear-gradient(135deg, #ff8484ea 0%, #52b4b4 100%);
      -webkit-mask: var(--icon-mask) no-repeat center;
      mask: var(--icon-mask) no-repeat center;
      -webkit-mask-size: contain;
      mask-size: contain;
    }
  }
}

.nav-icon {
  width: 20px;
  height: 20px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.home-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z'/%3E%3Cpolyline points='9,22 9,12 15,12 15,22'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z'/%3E%3Cpolyline points='9,22 9,12 15,12 15,22'/%3E%3C/svg%3E");
}

.messages-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z'/%3E%3C/svg%3E");
}

.settings-icon {
  --icon-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='3'/%3E%3Cpath d='M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z'/%3E%3C/svg%3E");
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='12' cy='12' r='3'/%3E%3Cpath d='M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z'/%3E%3C/svg%3E");
}


.user-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: auto;
  width: 40px;
  
  &:hover {
    background: #e2e8f0;
  }
}

.user-avatar {
  width: 37px;
  height: 37px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 6px;
  border: 2px solid #2632449c;
  opacity: 0.9;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}



</style>
