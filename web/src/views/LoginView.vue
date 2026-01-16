<script setup>
import { ref } from 'vue';
import api from '../services/API';
import router from '@/router';

const form = ref({ email: '', password: '' });
const loading = ref(false);
const error = ref('');

async function submit() {
  loading.value = true;
  error.value = '';
  try {
    const res = await api.post('/login', {
      email: form.value.email.trim(),
      password: form.value.password,
    });
    const token = res.data?.token;
    if (token) {
      localStorage.setItem('token', token);
      router.push({ name: 'account' });
    } else {
      throw new Error('Token não retornado');
    }
  } catch (e) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao autenticar';
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="max-w-md mx-auto p-6">
    <h1 class="text-xl font-semibold mb-4">Login</h1>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Email</label>
        <input v-model="form.email" type="email" class="mt-1 block w-full rounded border-gray-300 px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Senha</label>
        <input v-model="form.password" type="password" class="mt-1 block w-full rounded border-gray-300 px-3 py-2" />
      </div>
      <button
        class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
        :disabled="loading"
        @click="submit"
      >
        {{ loading ? 'Entrando...' : 'Entrar' }}
      </button>
      <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
    </div>
  </main>
</template>