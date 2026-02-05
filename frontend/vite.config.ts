import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // 代理 API 请求到后端
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 代理健康检查请求
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  // 确保 React 等依赖被正确预构建
  optimizeDeps: {
    include: ['react', 'react-dom', 'react-router-dom'],
  },
})
