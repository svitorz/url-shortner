<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/API';

const loading = ref(false);
const error = ref('');
const links = ref([]);
const actionMsg = ref('');

async function fetchLinks() {
  loading.value = true;
  error.value = '';
  actionMsg.value = '';
  try {
    const res = await api.get('/links');
    // Espera resposta no formato { links: [{ slug, target_url, is_active, last_change }...] }
    const data = res.data?.links ?? res.data?.link ?? []; // fallback caso backend use "link"
    links.value = Array.isArray(data) ? data : [];
  } catch (e) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao carregar seus links';
  } finally {
    loading.value = false;
  }
}

async function deactivateLink(slug, idx) {
  if (!slug) return;
  actionMsg.value = '';
  try {
    await api.post(`/links/${slug}/deactivate`);
    // Atualiza estado local
    if (typeof idx === 'number' && links.value[idx]) {
      links.value[idx].is_active = false;
      links.value[idx].last_change = new Date().toISOString();
    }
    actionMsg.value = `Link ${slug} desativado`;
  } catch (e) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao desativar link';
  }
}

function formatDate(d) {
  try {
    const date = new Date(d);
    return isNaN(date) ? String(d ?? '') : date.toLocaleString();
  } catch {
    return String(d ?? '');
  }
}

onMounted(fetchLinks);
</script>

<template>
  <main class="max-w-5xl mx-auto p-6 space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">Meus links</h1>
      <button
        class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
        :disabled="loading"
        @click="fetchLinks"
      >
        {{ loading ? 'Atualizando...' : 'Atualizar' }}
      </button>
    </div>

    <p v-if="error" class="text-red-600">{{ error }}</p>
    <p v-if="actionMsg" class="text-green-700">{{ actionMsg }}</p>

    <div v-if="!loading && links.length === 0" class="rounded border p-4 text-gray-700">
      Nenhum link encontrado. Crie um novo em “Criar link”.
    </div>

    <div class="overflow-x-auto rounded border">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Slug</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">URL destino</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ativo</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Última alteração</th>
            <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Ações</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="(item, idx) in links" :key="item.slug">
            <td class="px-4 py-3 font-mono text-sm text-gray-900">{{ item.slug }}</td>
            <td class="px-4 py-3 text-sm">
              <a
                class="text-indigo-600 hover:text-indigo-500 break-all"
                :href="item.target_url ?? item.targetUrl"
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ item.target_url ?? item.targetUrl }}
              </a>
            </td>
            <td class="px-4 py-3 text-sm">
              <span :class="(item.is_active ?? item.active) ? 'text-green-700' : 'text-gray-500'">
                {{ (item.is_active ?? item.active) ? 'Sim' : 'Não' }}
              </span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-700">{{ formatDate(item.last_change ?? item.updated_at) }}</td>
            <td class="px-4 py-3 text-right">
              <div class="inline-flex gap-2">
                <a
                  class="rounded border px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                  :href="`http://localhost:8000/api/r/${item.slug}`"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Abrir curto
                </a>
                <button
                  class="rounded bg-red-600 px-3 py-1.5 text-sm font-semibold text-white hover:bg-red-500 disabled:opacity-50"
                  :disabled="!(item.is_active ?? item.active)"
                  @click="deactivateLink(item.slug, idx)"
                >
                  Desativar
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>