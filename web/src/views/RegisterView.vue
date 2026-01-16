<script setup>
import { ref } from 'vue';
import api from '../services/API';

const form = ref({ full_name: '', email: '', password: '' });
const loading = ref(false);
const success = ref('');
const error = ref('');

async function submit() {
  loading.value = true;
  success.value = '';
  error.value = '';

  try {
    const res = await api.post('/user', {
      full_name: form.value.full_name.trim(),
      email: form.value.email.trim(),
      password: form.value.password,
    });
    success.value = res.data?.message ?? 'Usuário criado com sucesso';
    form.value = { full_name: '', email: '', password: '' };
  } catch (e) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao criar usuário';
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="max-w-md mx-auto p-6">
    <h1 class="text-xl font-semibold mb-4">Criar conta</h1>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Nome completo</label>
        <input v-model="form.full_name" type="text" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Email</label>
        <input v-model="form.email" type="email" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Senha</label>
        <input v-model="form.password" type="password" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
      </div>

      <button
        class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
        :disabled="loading"
        @click="submit"
      >
        {{ loading ? 'Criando...' : 'Criar conta' }}
      </button>

      <p v-if="success" class="text-green-600 text-sm">{{ success }}</p>
      <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
    </div>
  </main>
</template>