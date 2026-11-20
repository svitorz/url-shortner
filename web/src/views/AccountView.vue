<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'
import { useTheme } from '../composables/useTheme'
import { useToast } from '../composables/useToast'

const router = useRouter()
const { userId, logout } = useAuth()
const { mode } = useTheme()
const toast = useToast()

// Não há GET /me na API: os campos começam vazios e o PUT exige os dois.
const form = ref({ full_name: '', email: '' })
const loading = ref(false)
const error = ref('')

async function updateUser() {
  loading.value = true
  error.value = ''
  try {
    await api.put(`/user/${userId.value}`, {
      full_name: form.value.full_name.trim(),
      email: form.value.email.trim(),
    })
    toast.success('Dados atualizados.')
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
  <div class="mx-auto max-w-2xl space-y-6 px-4 py-10 sm:px-6">
    <h1 class="text-2xl font-bold tracking-tight">Minha conta</h1>

    <section class="card p-6">
      <h2 class="font-semibold">Dados pessoais</h2>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">Informe nome e e-mail (os dois são obrigatórios).</p>
      <form class="mt-5 space-y-4" @submit.prevent="updateUser">
        <div>
          <label for="name" class="label">Nome completo</label>
          <input id="name" v-model="form.full_name" type="text" required autocomplete="name" class="input mt-1" />
        </div>
        <div>
          <label for="email" class="label">E-mail</label>
          <input id="email" v-model="form.email" type="email" required autocomplete="email" class="input mt-1" />
        </div>
        <button type="submit" :disabled="loading || !userId" class="btn-primary">
          {{ loading ? 'Salvando...' : 'Salvar' }}
        </button>
        <p v-if="error" class="text-sm text-rose-600 dark:text-rose-400" role="alert">{{ error }}</p>
      </form>
    </section>

    <section class="card p-6">
      <h2 class="font-semibold">Aparência</h2>
      <fieldset class="mt-4">
        <legend class="sr-only">Tema</legend>
        <div class="grid grid-cols-3 gap-2">
          <label v-for="opt in [{ v: 'system', l: 'Sistema' }, { v: 'light', l: 'Claro' }, { v: 'dark', l: 'Escuro' }]" :key="opt.v"
            :class="mode === opt.v ? 'border-brand-500 ring-2 ring-brand-500/30' : 'border-slate-300 dark:border-slate-700'"
            class="flex cursor-pointer items-center justify-center rounded-lg border px-3 py-2 text-sm font-medium">
            <input v-model="mode" type="radio" name="theme" :value="opt.v" class="sr-only" />
            {{ opt.l }}
          </label>
        </div>
      </fieldset>
    </section>

    <section class="card flex items-center justify-between p-6">
      <div>
        <h2 class="font-semibold">Sessão</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Sair remove o acesso deste navegador.</p>
      </div>
      <button type="button" class="btn-secondary" @click="signOut">Sair</button>
    </section>
  </div>
</template>
