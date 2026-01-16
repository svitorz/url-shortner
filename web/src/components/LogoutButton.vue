<script setup>
import api from '../services/API';
import { ref } from 'vue';
import router from '@/router';

const loading = ref(false);
const error = ref('');

async function logout() {
  loading.value = true;
  error.value = '';
  try {
    await api.post('/logout');
  } catch (e) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao deslogar';
  } finally {
    localStorage.removeItem('token');
    loading.value = false;
    router.push({ name: 'login' });
  }
}
</script>

<template>
  <div class="flex items-center gap-3">
    <button
      class="rounded bg-gray-800 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-gray-700 disabled:opacity-50"
      :disabled="loading"
      @click="logout"
    >
      {{ loading ? 'Saindo...' : 'Sair' }}
    </button>
    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
  </div>
</template>