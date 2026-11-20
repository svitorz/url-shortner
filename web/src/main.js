import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { setUnauthorizedHandler } from './services/API'
import { useAuth } from './composables/useAuth'

// Token expirado ou inválido em qualquer chamada: limpa a sessão e manda ao
// login, voltando depois para onde o usuário estava.
const { logout } = useAuth()
setUnauthorizedHandler(() => {
  logout()
  const current = router.currentRoute.value
  if (current.name !== 'login') {
    router.push({ name: 'login', query: { redirect: current.fullPath, expired: '1' } })
  }
})

createApp(App).use(router).mount('#app')
