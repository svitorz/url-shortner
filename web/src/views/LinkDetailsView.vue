<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api, { errorMessage } from '../services/API'
import { useCopy } from '../composables/useCopy'
import { formatDate } from '../utils/format'

const route = useRoute()
const loading = ref(true)
const error = ref('')
const notFound = ref(false)
const link = ref(null)
const deactivating = ref(false)
const actionMsg = ref('')
const { copy, label } = useCopy()

async function fetchDetails() {
  loading.value = true
  error.value = ''
  notFound.value = false
  try {
    const res = await api.get(`/links/${encodeURIComponent(route.params.slug)}`)
    link.value = res.data.link
  } catch (e) {
    link.value = null
    // A API responde 404 também para link de outro usuário.
    if (e?.response?.status === 404) notFound.value = true
    else error.value = errorMessage(e, 'Erro ao carregar o link')
  } finally {
    loading.value = false
  }
}

async function deactivate() {
  if (!window.confirm('Desativar este link? Quem acessar vai ver "link desativado". Não dá para desfazer.')) return
  deactivating.value = true
  actionMsg.value = ''
  error.value = ''
  try {
    await api.post(`/links/${encodeURIComponent(link.value.slug)}/deactivate`)
    actionMsg.value = 'Link desativado.'
    await fetchDetails()
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao desativar')
  } finally {
    deactivating.value = false
  }
}

onMounted(fetchDetails)
watch(() => route.params.slug, (s) => s && fetchDetails())
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-6 p-6">
    <RouterLink to="/links" class="text-sm text-indigo-600 hover:underline">&larr; Meus links</RouterLink>
    <h1 class="text-xl font-semibold">Detalhes do link</h1>

    <p v-if="loading" class="text-gray-600">Carregando...</p>

    <div v-else-if="notFound" class="rounded border p-4 text-gray-700">
      Link não encontrado.
    </div>

    <div v-else-if="error && !link" class="rounded border border-red-200 bg-red-50 p-4 text-sm text-red-700" role="alert">
      {{ error }}
      <button type="button" class="ml-2 underline" @click="fetchDetails">Tentar de novo</button>
    </div>

    <template v-else-if="link">
      <dl class="space-y-3 rounded border p-4 text-sm">
        <div>
          <dt class="font-medium text-gray-700">URL curta</dt>
          <dd class="mt-1 flex items-center gap-2">
            <a :href="link.url" target="_blank" rel="noopener noreferrer"
              class="break-all font-mono text-indigo-600 hover:underline">{{ link.url }}</a>
            <button type="button" class="shrink-0 rounded border px-2.5 py-1 text-gray-700 hover:bg-gray-50"
              @click="copy(link.url)">
              {{ label(link.url) }}
            </button>
          </dd>
        </div>
        <div>
          <dt class="font-medium text-gray-700">Destino</dt>
          <dd class="mt-1 break-all">{{ link.target_url }}</dd>
        </div>
        <div>
          <dt class="font-medium text-gray-700">Status</dt>
          <dd class="mt-1">
            <span :class="link.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'"
              class="rounded-full px-2 py-0.5 text-xs font-medium">
              {{ link.is_active ? 'Ativo' : 'Desativado' }}
            </span>
          </dd>
        </div>
        <div>
          <dt class="font-medium text-gray-700">Última alteração</dt>
          <dd class="mt-1">{{ formatDate(link.last_change) }}</dd>
        </div>
      </dl>

      <button type="button" :disabled="!link.is_active || deactivating"
        class="rounded bg-red-600 px-3 py-2 text-sm font-semibold text-white hover:bg-red-500 disabled:opacity-40"
        @click="deactivate">
        {{ deactivating ? 'Desativando...' : link.is_active ? 'Desativar link' : 'Link desativado' }}
      </button>
      <p v-if="actionMsg" class="text-sm text-green-700" role="status">{{ actionMsg }}</p>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
    </template>
  </div>
</template>
