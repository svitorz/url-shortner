<script setup>
// Formulário de encurtar usado no topo de "Meus links" e na tela "Criar link".
import { ref } from 'vue'
import { LinkIcon } from '@heroicons/vue/24/outline'
import api, { errorMessage } from '../services/API'

const emit = defineEmits(['created'])
defineProps({ autofocus: { type: Boolean, default: false } })

const target = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const url = target.value.trim()
    const res = await api.post('/links', { target_url: url })
    emit('created', { slug: res.data.slug, url: res.data.url, target_url: url })
    target.value = ''
  } catch (e) {
    error.value = errorMessage(e, 'Erro ao criar link')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="space-y-2" @submit.prevent="submit">
    <label for="quick-target" class="sr-only">URL de destino</label>
    <div class="flex flex-col gap-2 sm:flex-row">
      <div class="relative flex-1">
        <LinkIcon class="pointer-events-none absolute left-3 top-1/2 size-5 -translate-y-1/2 text-slate-400" aria-hidden="true" />
        <input id="quick-target" v-model="target" type="url" required maxlength="2048" :autofocus="autofocus"
          placeholder="Cole uma URL longa: https://..." class="input py-2.5 pl-10" />
      </div>
      <button type="submit" class="btn-primary py-2.5" :disabled="loading || !target">
        {{ loading ? 'Encurtando...' : 'Encurtar' }}
      </button>
    </div>
    <p v-if="error" class="text-sm text-rose-600 dark:text-rose-400" role="alert">{{ error }}</p>
  </form>
</template>
