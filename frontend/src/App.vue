<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
// App.vue 使用 script setup
</script>

<style>
* {
  cursor: default;
}

button, .nav-icon, select, a {
  cursor: pointer;
}

input, textarea {
  cursor: text;
}

input[type="checkbox"], input[type="radio"] {
  cursor: pointer;
}

#app {
  width: 100%;
  height: 100vh;
  position: relative;
}
</style>
