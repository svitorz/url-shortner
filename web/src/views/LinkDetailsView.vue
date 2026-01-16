<script setup>
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import api from '../services/API';

const route = useRoute();
const slug = ref(route.params.slug);
const loading = ref(false);
const err = ref('');
const link = ref(null);
const deactivating = ref(false);
const deactivateMsg = ref('');

async function fetchDetails() {
  loading.value = true; err.value = ''; link.value = null;
  try {
    const res = await api.get(`/links/${slug.value}`);
    link.value = res.data ?? {};
  } catch (e) {
    err.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao buscar detalhes';
  } finally {
    loading.value = false;
  }
}

async function deactivateLink() {
  deactivating.value = true; deactivateMsg.value = ''; err.value = '';
  try {
    const res = await api.post(`/links/${slug.value}/deactivate`);
    deactivateMsg.value = res.data?.message ?? 'Link desativado';
    await fetchDetails();
  } catch (e) {
    err.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao desativar';
  } finally {
    deactivating.value = false;
  }
}

onMounted(fetchDetails);
</script>

<template>
  <main class="max-w-2xl mx-auto p-6 space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold">Detalhes do link</h1>
      <div class="flex gap-2">
        <button
          class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-semibold text-white hover:bg-indigo-500"
          @click="fetchDetails"
        >
          Atualizar
        </button>
        <button
          class="rounded bg-red-600 px-3 py-1.5 text-sm font-semibold text-white hover:bg-red-500 disabled:opacity-50"
          :disabled="deactivating"
          @click="deactivateLink"
        >
          {{ deactivating ? 'Desativando...' : 'Desativar' }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-gray-600">Carregando...</div>
    <p v-if="err" class="text-red-600">{{ err }}</p>

    <div v-if="link" class="rounded border p-4 space-y-2">
      <p><span class="font-medium">Slug:</span> <span class="font-mono">{{ slug }}</span></p>
      <p><span class="font-medium">URL destino:</span> {{ link.target_url ?? link.targetUrl }}</p>
      <p><span class="font-medium">Ativo:</span> <span :class="(link.active ?? link.is_active) ? 'text-green-700' : 'text-gray-500'">
        {{ (link.active ?? link.is_active) ? 'Sim' : 'Não' }}
      </span></p>
    </div>

    <p v-if="deactivateMsg" class="text-green-600">{{ deactivateMsg }}</p>
  </main>
</template>