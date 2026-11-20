<script setup>
import { onMounted, ref } from 'vue'
import api, { errorMessage } from '../services/API'
import { useCopy } from '../composables/useCopy'
import { formatDate } from '../utils/format'

const loading = ref(true)
const error = ref('')
const links = ref([])
const busy = ref('') // slug sendo desativado
const actionMsg = ref('')
const { copy, label } = useCopy()

async function fetchLinks() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('/links')
    links.value = res.data.links ?? []
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao carregar seus links')
  } finally {
    loading.value = false
  }
}

async function deactivate(link) {
  // Não há como reativar pela API: confirmar antes.
  if (!window.confirm(`Desativar ${link.url}? Quem acessar vai ver "link desativado". Não dá para desfazer.`)) return
  busy.value = link.slug
  actionMsg.value = ''
  error.value = ''
  try {
    await api.post(`/links/${link.slug}/deactivate`)
    link.is_active = false
    link.last_change = new Date().toISOString()
    actionMsg.value = `Link ${link.slug} desativado.`
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao desativar link')
  } finally {
    busy.value = ''
  }
}

onMounted(fetchLinks)
</script>

<template>
  <div class="mx-auto max-w-5xl space-y-6 p-6">
    <div class="flex items-center justify-between gap-4">
      <h1 class="text-2xl font-semibold">Meus links</h1>
      <div class="flex gap-2">
        <button type="button" :disabled="loading"
          class="rounded border px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
          @click="fetchLinks">
          Atualizar
        </button>
        <RouterLink to="/links/new"
          class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-semibold text-white hover:bg-indigo-500">
          Novo link
        </RouterLink>
      </div>
    </div>

    <p v-if="actionMsg" class="text-sm text-green-700" role="status">{{ actionMsg }}</p>

    <!-- carregando -->
    <p v-if="loading" class="text-gray-600">Carregando...</p>

    <!-- erro -->
    <div v-else-if="error && links.length === 0" class="rounded border border-red-200 bg-red-50 p-4 text-sm text-red-700" role="alert">
      {{ error }}
      <button type="button" class="ml-2 underline" @click="fetchLinks">Tentar de novo</button>
    </div>

    <!-- vazio -->
    <div v-else-if="links.length === 0" class="rounded border p-6 text-center text-gray-700">
      Você ainda não tem links.
      <RouterLink to="/links/new" class="text-indigo-600 hover:underline">Crie o primeiro</RouterLink>.
    </div>

    <!-- lista -->
    <template v-else>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
      <div class="overflow-x-auto rounded border">
        <table class="min-w-full divide-y divide-gray-200 text-sm">
          <thead class="bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
            <tr>
              <th class="px-4 py-3">URL curta</th>
              <th class="px-4 py-3">Destino</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3">Última alteração</th>
              <th class="px-4 py-3 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white">
            <tr v-for="link in links" :key="link.slug">
              <td class="px-4 py-3 font-mono">
                <RouterLink :to="{ name: 'links-details', params: { slug: link.slug } }"
                  class="text-indigo-600 hover:underline">{{ link.url }}</RouterLink>
              </td>
              <td class="max-w-xs truncate px-4 py-3" :title="link.target_url">{{ link.target_url }}</td>
              <td class="px-4 py-3">
                <span :class="link.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'"
                  class="rounded-full px-2 py-0.5 text-xs font-medium">
                  {{ link.is_active ? 'Ativo' : 'Desativado' }}
                </span>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-gray-700">{{ formatDate(link.last_change) }}</td>
              <td class="px-4 py-3 text-right">
                <div class="inline-flex gap-2">
                  <button type="button" class="rounded border px-2.5 py-1 text-gray-700 hover:bg-gray-50"
                    @click="copy(link.url)">
                    {{ label(link.url) }}
                  </button>
                  <button type="button" :disabled="!link.is_active || busy === link.slug"
                    class="rounded bg-red-600 px-2.5 py-1 font-semibold text-white hover:bg-red-500 disabled:opacity-40"
                    @click="deactivate(link)">
                    {{ busy === link.slug ? 'Desativando...' : 'Desativar' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
