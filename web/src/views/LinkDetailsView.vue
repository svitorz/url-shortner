<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeftIcon, ArrowTopRightOnSquareIcon, CursorArrowRaysIcon, ClockIcon } from '@heroicons/vue/24/outline'
import api, { errorMessage } from '../services/API'
import { useToast } from '../composables/useToast'
import { useConfirm } from '../composables/useConfirm'
import { formatDate, formatNumber } from '../utils/format'
import CopyButton from '../components/CopyButton.vue'
import QrCode from '../components/QrCode.vue'
import StatusBadge from '../components/StatusBadge.vue'

const route = useRoute()
const toast = useToast()
const { confirm } = useConfirm()

const loading = ref(true)
const error = ref('')
const notFound = ref(false)
const link = ref(null)
const deactivating = ref(false)

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
  const ok = await confirm({
    title: 'Desativar link?',
    message: 'Quem acessar vai ver "link desativado". Não dá para reativar.',
    confirmLabel: 'Desativar',
    danger: true,
  })
  if (!ok) return
  deactivating.value = true
  try {
    await api.post(`/links/${encodeURIComponent(link.value.slug)}/deactivate`)
    toast.success('Link desativado.')
    await fetchDetails()
  } catch (e) {
    toast.error(errorMessage(e, 'Erro ao desativar'))
  } finally {
    deactivating.value = false
  }
}

onMounted(fetchDetails)
watch(() => route.params.slug, (s) => s && fetchDetails())
</script>

<template>
  <div class="mx-auto max-w-4xl space-y-6 px-4 py-8 sm:px-6">
    <RouterLink to="/links" class="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200">
      <ArrowLeftIcon class="size-4" aria-hidden="true" /> Meus links
    </RouterLink>

    <!-- carregando -->
    <div v-if="loading" class="grid gap-6 md:grid-cols-[1fr_auto]" aria-busy="true">
      <div class="card h-56 animate-pulse bg-slate-100 dark:bg-slate-800/60" />
      <div class="card h-56 w-60 animate-pulse bg-slate-100 dark:bg-slate-800/60" />
      <p class="sr-only" role="status">Carregando...</p>
    </div>

    <div v-else-if="notFound" class="card p-10 text-center">
      <h1 class="font-semibold">Link não encontrado</h1>
      <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">Ele não existe ou pertence a outra conta.</p>
    </div>

    <div v-else-if="error" class="card p-6 text-sm text-rose-700 dark:text-rose-300" role="alert">
      {{ error }}
      <button type="button" class="ml-2 font-semibold underline" @click="fetchDetails">Tentar de novo</button>
    </div>

    <div v-else-if="link" class="grid gap-6 md:grid-cols-[1fr_auto]">
      <div class="card space-y-6 p-6">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="min-w-0">
            <p class="text-sm text-slate-500 dark:text-slate-400">Link curto</p>
            <h1 class="break-all font-mono text-2xl font-bold text-brand-700 dark:text-brand-300">
              {{ link.url.replace(/^https?:\/\//, '') }}
            </h1>
          </div>
          <StatusBadge :active="link.is_active" />
        </div>

        <div class="flex flex-wrap gap-2">
          <CopyButton :text="link.url" />
          <a :href="link.url" target="_blank" rel="noopener noreferrer" class="btn-secondary">
            <ArrowTopRightOnSquareIcon class="size-4" aria-hidden="true" /> Abrir
          </a>
        </div>

        <div>
          <p class="text-sm text-slate-500 dark:text-slate-400">Destino</p>
          <a :href="link.target_url" target="_blank" rel="noopener noreferrer"
            class="break-all text-slate-800 hover:underline dark:text-slate-200">{{ link.target_url }}</a>
        </div>

        <dl class="grid grid-cols-2 gap-4">
          <div class="rounded-xl bg-slate-50 p-4 dark:bg-slate-800/50">
            <dt class="flex items-center gap-1.5 text-sm text-slate-500 dark:text-slate-400">
              <CursorArrowRaysIcon class="size-4" aria-hidden="true" /> Cliques
            </dt>
            <dd class="mt-1 text-3xl font-semibold tabular-nums">{{ formatNumber(link.clicks ?? 0) }}</dd>
            <p class="text-xs text-slate-400">atualiza a cada ~10 s</p>
          </div>
          <div class="rounded-xl bg-slate-50 p-4 dark:bg-slate-800/50">
            <dt class="flex items-center gap-1.5 text-sm text-slate-500 dark:text-slate-400">
              <ClockIcon class="size-4" aria-hidden="true" /> Última alteração
            </dt>
            <dd class="mt-1 text-lg font-semibold">{{ formatDate(link.last_change) }}</dd>
          </div>
        </dl>

        <div class="border-t border-slate-200 pt-4 dark:border-slate-800">
          <button type="button" class="btn-danger" :disabled="!link.is_active || deactivating" @click="deactivate">
            {{ deactivating ? 'Desativando...' : link.is_active ? 'Desativar link' : 'Link desativado' }}
          </button>
        </div>
      </div>

      <div class="card flex flex-col items-center justify-center p-6">
        <p class="mb-3 text-sm font-medium text-slate-600 dark:text-slate-300">QR code</p>
        <QrCode :value="link.url" :filename="`qr-${link.slug}`" />
      </div>
    </div>
  </div>
</template>
