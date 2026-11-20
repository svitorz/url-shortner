<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useAuth } from '../composables/useAuth'
import { useToast } from '../composables/useToast'
import AuthCard from '../components/AuthCard.vue'

const router = useRouter()
const { login } = useAuth()
const toast = useToast()

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
    toast.success('Conta criada. Bem-vindo!')
    router.push({ name: 'links' })
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao criar conta')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthCard title="Criar conta" subtitle="Comece a encurtar links em segundos.">
    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label for="name" class="label">Nome completo</label>
        <input id="name" v-model="form.full_name" type="text" required autocomplete="name" class="input mt-1" />
      </div>
      <div>
        <label for="email" class="label">E-mail</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email" class="input mt-1" />
      </div>
      <div>
        <label for="password" class="label">Senha</label>
        <!-- 8 é o mínimo da API; 72 bytes, o limite do bcrypt -->
        <input id="password" v-model="form.password" type="password" required minlength="8" maxlength="72"
          autocomplete="new-password" class="input mt-1" />
        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">Mínimo de 8 caracteres.</p>
      </div>
      <button type="submit" :disabled="loading" class="btn-primary w-full py-2.5">
        {{ loading ? 'Criando...' : 'Criar conta' }}
      </button>
      <p v-if="error" class="text-sm text-rose-600 dark:text-rose-400" role="alert">{{ error }}</p>
    </form>

    <template #footer>
      Já tem conta? <RouterLink to="/login" class="font-medium text-brand-600 hover:underline dark:text-brand-400">Entrar</RouterLink>
    </template>
  </AuthCard>
</template>
