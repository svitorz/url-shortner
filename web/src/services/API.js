// Cliente HTTP único do front. Todas as chamadas à API passam por aqui.
//
// baseURL relativa (padrão /api): o front é servido pela mesma origem que a
// API (nginx em produção, proxy do Vite em dev), então não há URL absoluta
// nem CORS. VITE_API_BASE é lida no build.
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api',
  timeout: 15000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Quem trata o 401 (limpar a sessão e ir ao login) é registrado no main.js:
// assim este módulo não importa o router (evita import circular).
let onUnauthorized = () => {}
export function setUnauthorizedHandler(fn) {
  onUnauthorized = fn
}

api.interceptors.response.use(
  (res) => res,
  (error) => {
    // 401 no próprio login é "senha errada", não sessão expirada.
    const isLogin = error.config?.url === '/login'
    if (error.response?.status === 401 && !isLogin) {
      onUnauthorized()
    }
    return Promise.reject(error)
  },
)

// Mensagem para mostrar ao usuário a partir de um erro do axios.
export function errorMessage(e, fallback = 'Algo deu errado') {
  if (e?.response?.data?.error) return e.response.data.error
  if (e?.response?.status === 429) return 'Muitas requisições. Tente de novo em instantes.'
  if (e?.code === 'ECONNABORTED') return 'O servidor demorou para responder.'
  if (e?.request && !e?.response) return 'Sem conexão com o servidor.'
  return fallback
}

export default api
