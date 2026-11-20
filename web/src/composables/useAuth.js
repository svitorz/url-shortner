// Estado de autenticação compartilhado (módulo único = mesma instância em
// todos os componentes). O token fica no localStorage.
import { computed, ref } from 'vue'

const token = ref(localStorage.getItem('token'))

// Lê o user_id do payload do JWT só para a interface (ex.: rota /user/:id).
// Não valida assinatura: quem valida é a API em toda requisição.
function userIdFrom(t) {
  try {
    const payload = t.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(payload)).user_id ?? null
  } catch {
    return null
  }
}

const isAuthenticated = computed(() => !!token.value)
const userId = computed(() => (token.value ? userIdFrom(token.value) : null))

function login(newToken) {
  localStorage.setItem('token', newToken)
  token.value = newToken
}

function logout() {
  localStorage.removeItem('token')
  token.value = null
}

export function useAuth() {
  return { isAuthenticated, userId, login, logout }
}
