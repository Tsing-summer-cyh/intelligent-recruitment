import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 匹配 /api 开头的请求
      '/api': {
        target: 'http://localhost:8080', // 本地 Go 服务的地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '') // 根据后端路由决定是否去掉前缀
      }
    }
  }
})