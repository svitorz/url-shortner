<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { userId, logout } = useAuth()

// Não há GET /me na API: os campos começam vazios e PUT exige os dois.
const form = ref({ full_name: '', email: '' })
const loading = ref(false)
const msg = ref('')
const error = ref('')

async function updateUser() {
  loading.value = true
  msg.value = ''
  error.value = ''
  try {
    const res = await api.put(`/user/${userId.value}`, {
      full_name: form.value.full_name.trim(),
      email: form.value.email.trim(),
    })
    msg.value = res.data?.message ?? 'Dados atualizados'
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao atualizar')
  } finally {
    loading.value = false
  }
}

function signOut() {
  logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="mx-auto max-w-lg space-y-6 p-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold">Minha conta</h1>
      <button type="button"
        class="rounded bg-gray-800 px-3 py-1.5 text-sm font-semibold text-white hover:bg-gray-700"
        @click="signOut">
        Sair
      </button>
    </div>

    <form class="space-y-4" @submit.prevent="updateUser">
      <p class="text-sm text-gray-600">Atualize seu nome e e-mail (os dois são obrigatórios).</p>
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
      <button type="submit" :disabled="loading || !userId"
        class="rounded bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50">
        {{ loading ? 'Salvando...' : 'Salvar' }}
      </button>
      <p v-if="msg" class="text-sm text-green-600" role="status">{{ msg }}</p>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
    </form>
  </div>
</template>
