<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'
import AuthCard from '../components/AuthCard.vue'

const route = useRoute()
const router = useRouter()
const { login } = useAuth()

const form = ref({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

// Só aceita redirect interno (começa com uma única "/"), nunca outro site.
function safeRedirect(r) {
  return typeof r === 'string' && r.startsWith('/') && !r.startsWith('//') ? r : '/links'
}

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.post('/login', {
      email: form.value.email.trim(),
      password: form.value.password,
    })
    login(res.data.token)
    router.push(safeRedirect(route.query.redirect))
  } catch (e) {
    error.value = e?.response?.status === 401 ? 'E-mail ou senha incorretos.' : errorMessage(e, 'Erro ao entrar')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthCard title="Entrar" subtitle="Acesse seus links curtos.">
    <p v-if="route.query.expired" class="mb-4 rounded-lg bg-amber-50 p-3 text-sm text-amber-800 dark:bg-amber-500/10 dark:text-amber-300">
      Sua sessão expirou. Entre de novo.
    </p>

    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="email" class="label">E-mail</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email" class="input mt-1" />
      </div>
      <div>
        <label for="password" class="label">Senha</label>
        <input id="password" v-model="form.password" type="password" required autocomplete="current-password" class="input mt-1" />
      </div>
      <button type="submit" :disabled="loading" class="btn-primary w-full py-2.5">
        {{ loading ? 'Entrando...' : 'Entrar' }}
      </button>
      <p v-if="error" class="text-sm text-rose-600 dark:text-rose-400" role="alert">{{ error }}</p>
    </form>

    <template #footer>
      Não tem conta? <RouterLink to="/register" class="font-medium text-brand-600 hover:underline dark:text-brand-400">Criar conta</RouterLink>
    </template>
  </AuthCard>
</template>
