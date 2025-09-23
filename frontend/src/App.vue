<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
* {
  cursor: default;
}

button, .nav-icon, input, select, textarea, a {
  cursor: pointer;
}

#app {
  width: 100%;
  height: 100vh;
  position: relative;
}
</style>
