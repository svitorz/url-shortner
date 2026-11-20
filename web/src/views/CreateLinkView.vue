<script setup>
import { ref } from 'vue'
import api, { errorMessage } from '../services/API'
import { useCopy } from '../composables/useCopy'

const targetUrl = ref('')
const loading = ref(false)
const error = ref('')
const created = ref(null) // { slug, url }
const { copy, label } = useCopy()

async function createLink() {
  loading.value = true
  error.value = ''
  created.value = null
  try {
    const res = await api.post('/links', { target_url: targetUrl.value.trim() })
    created.value = { slug: res.data.slug, url: res.data.url, target: targetUrl.value.trim() }
    targetUrl.value = ''
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao criar link')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-xl space-y-6 p-6">
    <h1 class="text-xl font-semibold">Criar link curto</h1>

    <form class="space-y-4" @submit.prevent="createLink">
      <div>
        <label for="target" class="block text-sm font-medium text-gray-700">URL de destino</label>
        <input id="target" v-model="targetUrl" type="url" required maxlength="2048"
          placeholder="https://exemplo.com/pagina-longa"
          class="mt-1 block w-full rounded border border-gray-300 px-3 py-2" />
        <p class="mt-1 text-xs text-gray-500">Só http:// ou https://.</p>
      </div>
      <button type="submit" :disabled="loading || !targetUrl"
        class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50">
        {{ loading ? 'Criando...' : 'Encurtar' }}
      </button>
      <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
    </form>

    <div v-if="created" class="space-y-3 rounded border border-green-200 bg-green-50 p-4">
      <p class="text-sm text-green-800">Link criado para <span class="break-all">{{ created.target }}</span></p>
      <div class="flex items-center gap-2">
        <input :value="created.url" readonly aria-label="URL curta"
          class="block w-full rounded border border-gray-300 bg-white px-3 py-2 font-mono text-sm"
          @focus="$event.target.select()" />
        <button type="button"
          class="shrink-0 rounded bg-gray-800 px-3 py-2 text-sm font-semibold text-white hover:bg-gray-700"
          @click="copy(created.url)">
          {{ label(created.url) }}
        </button>
      </div>
      <RouterLink :to="{ name: 'links-details', params: { slug: created.slug } }"
        class="inline-block text-sm text-indigo-600 hover:underline">
        Ver detalhes
      </RouterLink>
    </div>
  </div>
</template>
