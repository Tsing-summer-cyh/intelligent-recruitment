<template>
  <div class="app-container">
    <router-view></router-view>
  </div>
</template>

<style>
body {
  margin: 0;
  /* 苹果风格的系统字体栈 */
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f3f4f6; /* 高级灰背景 */
  color: #1f2937;
  -webkit-font-smoothing: antialiased;
}
.app-container {
  min-height: 100vh;
}
</style>