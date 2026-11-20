<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'

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
    error.value = errorMessage(e, 'Erro ao entrar')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md p-6">
    <h1 class="mb-4 text-xl font-semibold">Entrar</h1>

    <p v-if="route.query.expired" class="mb-4 rounded bg-amber-50 p-3 text-sm text-amber-800">
      Sua sessão expirou. Entre de novo.
    </p>

    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">E-mail</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email"
          class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
      </div>
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Senha</label>
        <input id="password" v-model="form.password" type="password" required autocomplete="current-password"
          class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
      </div>
      <button type="submit" :disabled="loading"
        class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50">
        {{ loading ? 'Entrando...' : 'Entrar' }}
      </button>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <p class="text-sm text-gray-600">
        Não tem conta? <RouterLink to="/register" class="text-indigo-600 hover:underline">Criar conta</RouterLink>
      </p>
    </form>
  </div>
</template>
