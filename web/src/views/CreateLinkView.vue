<script setup>
import { ref } from 'vue';
import api from '../services/API';

const targetUrl = ref('');
const loading = ref(false);
const msg = ref('');
const err = ref('');
const createdSlug = ref('');

async function createLink() {
  loading.value = true;
  msg.value = ''; err.value = ''; createdSlug.value = '';

  try {
    const res = await api.post('/links', { target_url: targetUrl.value.trim() });
    msg.value = res.data?.message ?? 'Link criado';
    createdSlug.value = res.data?.slug ?? '';
    targetUrl.value = '';
  } catch (e) {
    err.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao criar link';
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="max-w-md mx-auto p-6 space-y-4">
    <h1 class="text-xl font-semibold">Criar link curto</h1>
    <div>
      <label class="block text-sm font-medium text-gray-700">URL de destino</label>
      <input v-model="targetUrl" type="url" placeholder="https://example.com" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
    </div>

    <button
      class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
      :disabled="loading || !targetUrl"
      @click="createLink"
    >
      {{ loading ? 'Criando...' : 'Criar' }}
    </button>

    <div v-if="createdSlug" class="rounded bg-green-50 p-3 text-sm">
      <p class="text-green-700">Slug criado: <span class="font-mono">{{ createdSlug }}</span></p>
      <p class="text-gray-700">Detalhes: vá para <span class="font-mono">/links/{{ createdSlug }}</span></p>
    </div>

    <p v-if="msg" class="text-green-600 text-sm">{{ msg }}</p>
    <p v-if="err" class="text-red-600 text-sm">{{ err }}</p>
  </main>
</template>