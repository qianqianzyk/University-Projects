<script setup lang="ts">
import { onErrorCaptured, ref } from 'vue'

const error = ref<Error | null>(null)

onErrorCaptured((err) => {
  console.error('应用错误:', err)
  error.value = err
  return false
})
</script>

<template>
  <div id="app">
    <div v-if="error" class="error-container">
      <h2>应用出现错误</h2>
      <p>{{ error.message }}</p>
      <button @click="error = null">重试</button>
    </div>
    <router-view v-else />
  </div>
</template>

<style>
html, body, #app {
  width: 100vw;
  height: 100vh;
  min-width: 100vw;
  min-height: 100vh;
  margin: 0;
  padding: 0;
  overflow: hidden !important;
  box-sizing: border-box;
}

#app {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  text-align: center;
  padding: 20px;
}

.error-container h2 {
  color: #f56c6c;
  margin-bottom: 10px;
}

.error-container button {
  margin-top: 20px;
  padding: 10px 20px;
  background: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.error-container button:hover {
  background: #337ecc;
}
</style>
