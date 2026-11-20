import { fileURLToPath, URL } from 'node:url'
import process from 'node:process'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  // Dev (npm run dev): /api e /r vão para a API, como faz o nginx em produção.
  // Mesma origem = sem CORS. A API do compose fica em 127.0.0.1:APP_PORT.
  server: {
    proxy: {
      '/api': process.env.VITE_DEV_API || 'http://localhost:8000',
      '/r/': process.env.VITE_DEV_API || 'http://localhost:8000',
    },
  },
})
