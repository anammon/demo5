import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // 允许局域网和外部访问
    port: 5173, // 你的开发端口
    allowedHosts: [
      '.ngrok-free.dev', // ✅ 允许所有 ngrok 免费域名
      'knife-metro-wonder-manuals.trycloudflare.com',
    ],
  },
})
