<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { login } = useAuth()

const form = ref({ full_name: '', email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.post('/user', {
      full_name: form.value.full_name.trim(),
      email: form.value.email.trim(),
      password: form.value.password,
    })
    // A API já devolve um token no cadastro: entra direto.
    login(res.data.token)
    router.push({ name: 'links' })
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao criar conta')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md p-6">
    <h1 class="mb-4 text-xl font-semibold">Criar conta</h1>

    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="name" class="block text-sm font-medium text-gray-700">Nome completo</label>
        <input id="name" v-model="form.full_name" type="text" required autocomplete="name"
          class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
      </div>
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">E-mail</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email"
          class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
      </div>
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Senha</label>
        <!-- 72 bytes é o limite do bcrypt na API -->
        <input id="password" v-model="form.password" type="password" required minlength="8" maxlength="72"
          autocomplete="new-password" class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
        <p class="mt-1 text-xs text-gray-500">Entre 8 e 72 caracteres.</p>
      </div>
      <button type="submit" :disabled="loading"
        class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50">
        {{ loading ? 'Criando...' : 'Criar conta' }}
      </button>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <p class="text-sm text-gray-600">
        Já tem conta? <RouterLink to="/login" class="text-indigo-600 hover:underline">Entrar</RouterLink>
      </p>
    </form>
  </div>
</template>
