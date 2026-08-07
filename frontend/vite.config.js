import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // 将 /api 和 /helloworld 请求代理到后端 Go 服务
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      },
      '/helloworld': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      },
    },
  },
})