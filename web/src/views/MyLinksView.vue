<script setup>
import { computed, onMounted, ref } from 'vue'
import {
  LinkIcon,
  CheckCircleIcon,
  NoSymbolIcon,
  CursorArrowRaysIcon,
  MagnifyingGlassIcon,
  QrCodeIcon,
  ArrowPathIcon,
} from '@heroicons/vue/24/outline'
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'
import api, { errorMessage } from '../services/API'
import { useToast } from '../composables/useToast'
import { useConfirm } from '../composables/useConfirm'
import { formatDate, formatNumber } from '../utils/format'
import QuickCreate from '../components/QuickCreate.vue'
import StatCard from '../components/StatCard.vue'
import StatusBadge from '../components/StatusBadge.vue'
import CopyButton from '../components/CopyButton.vue'
import QrCode from '../components/QrCode.vue'

const toast = useToast()
const { confirm } = useConfirm()

const loading = ref(true)
const error = ref('')
const links = ref([])
const busy = ref('')
const highlight = ref('') // link recém-criado, destacado na lista
const qrFor = ref(null)

const query = ref('')
const status = ref('all') // all | active | inactive
const sort = ref('recent') // recent | clicks | slug

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

const stats = computed(() => ({
  total: links.value.length,
  active: links.value.filter((l) => l.is_active).length,
  inactive: links.value.filter((l) => !l.is_active).length,
  clicks: links.value.reduce((n, l) => n + (l.clicks ?? 0), 0),
}))

const visible = computed(() => {
  const q = query.value.trim().toLowerCase()
  const list = links.value.filter((l) => {
    if (status.value === 'active' && !l.is_active) return false
    if (status.value === 'inactive' && l.is_active) return false
    return !q || l.slug.toLowerCase().includes(q) || l.target_url.toLowerCase().includes(q)
  })
  const by = {
    recent: (a, b) => new Date(b.last_change) - new Date(a.last_change),
    clicks: (a, b) => (b.clicks ?? 0) - (a.clicks ?? 0),
    slug: (a, b) => a.slug.localeCompare(b.slug),
  }[sort.value]
  return [...list].sort(by)
})

async function onCreated(link) {
  highlight.value = link.slug
  toast.success('Link criado. Já está no topo da lista.')
  sort.value = 'recent'
  status.value = 'all'
  query.value = ''
  await fetchLinks()
}

async function deactivate(link) {
  const ok = await confirm({
    title: 'Desativar link?',
    message: `Quem acessar ${link.url} vai ver "link desativado". Não dá para reativar.`,
    confirmLabel: 'Desativar',
    danger: true,
  })
  if (!ok) return
  busy.value = link.slug
  try {
    await api.post(`/links/${link.slug}/deactivate`)
    link.is_active = false
    link.last_change = new Date().toISOString()
    toast.success(`Link ${link.slug} desativado.`)
  } catch (e) {
    toast.error(errorMessage(e, 'Erro ao desativar link'))
  } finally {
    busy.value = ''
  }
}

// Separa o domínio do resto do destino para destacar o domínio na lista.
function splitUrl(url) {
  try {
    const u = new URL(url)
    return { host: u.host, rest: url.slice(url.indexOf(u.host) + u.host.length) }
  } catch {
    return { host: url, rest: '' }
  }
}

onMounted(fetchLinks)
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 px-4 py-8 sm:px-6">
    <div class="flex flex-wrap items-end justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">Meus links</h1>
        <p class="text-sm text-slate-500 dark:text-slate-400">Cliques atualizam a cada ~10 segundos.</p>
      </div>
      <button type="button" class="btn-secondary" :disabled="loading" @click="fetchLinks">
        <ArrowPathIcon class="size-4" :class="loading && 'animate-spin'" aria-hidden="true" />
        Atualizar
      </button>
    </div>

    <div class="card p-4 sm:p-5">
      <QuickCreate @created="onCreated" />
    </div>

    <!-- contadores -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <template v-if="loading && links.length === 0">
        <div v-for="i in 4" :key="i" class="card h-[78px] animate-pulse bg-slate-100 dark:bg-slate-800/60" />
      </template>
      <template v-else>
        <StatCard label="Links" :value="formatNumber(stats.total)" :icon="LinkIcon" />
        <StatCard label="Ativos" :value="formatNumber(stats.active)" :icon="CheckCircleIcon" tone="emerald" />
        <StatCard label="Desativados" :value="formatNumber(stats.inactive)" :icon="NoSymbolIcon" tone="slate" />
        <StatCard label="Cliques" :value="formatNumber(stats.clicks)" :icon="CursorArrowRaysIcon" tone="amber" />
      </template>
    </div>

    <!-- erro sem dados -->
    <div v-if="error && links.length === 0"
      class="card border-rose-200 p-6 text-sm text-rose-700 dark:border-rose-500/30 dark:text-rose-300" role="alert">
      {{ error }}
      <button type="button" class="ml-2 font-semibold underline" @click="fetchLinks">Tentar de novo</button>
    </div>

    <!-- vazio -->
    <div v-else-if="!loading && links.length === 0" class="card flex flex-col items-center gap-3 p-10 text-center">
      <div class="flex size-14 items-center justify-center rounded-2xl bg-brand-50 dark:bg-brand-500/10">
        <LinkIcon class="size-7 text-brand-600 dark:text-brand-400" aria-hidden="true" />
      </div>
      <h2 class="font-semibold">Nenhum link ainda</h2>
      <p class="max-w-sm text-sm text-slate-500 dark:text-slate-400">
        Cole uma URL no campo acima para criar seu primeiro link curto.
      </p>
    </div>

    <template v-else>
      <!-- filtros -->
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <div class="relative flex-1">
          <MagnifyingGlassIcon class="pointer-events-none absolute left-3 top-1/2 size-5 -translate-y-1/2 text-slate-400" aria-hidden="true" />
          <label for="search" class="sr-only">Buscar</label>
          <input id="search" v-model="query" type="search" placeholder="Buscar por código ou destino" class="input pl-10" />
        </div>
        <div class="flex gap-2">
          <label for="status" class="sr-only">Status</label>
          <select id="status" v-model="status" class="input w-auto">
            <option value="all">Todos</option>
            <option value="active">Ativos</option>
            <option value="inactive">Desativados</option>
          </select>
          <label for="sort" class="sr-only">Ordenar</label>
          <select id="sort" v-model="sort" class="input w-auto">
            <option value="recent">Mais recentes</option>
            <option value="clicks">Mais cliques</option>
            <option value="slug">Código (A–Z)</option>
          </select>
        </div>
      </div>

      <p v-if="loading" class="sr-only" role="status">Carregando...</p>

      <div v-if="!loading && visible.length === 0" class="card p-6 text-center text-sm text-slate-500 dark:text-slate-400">
        Nenhum link corresponde aos filtros.
      </div>

      <!-- tabela (telas médias e grandes) -->
      <div v-if="visible.length" class="card hidden overflow-hidden md:block">
        <table class="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
          <thead class="bg-slate-50 text-left text-xs font-medium uppercase tracking-wider text-slate-500 dark:bg-slate-900/60 dark:text-slate-400">
            <tr>
              <th class="px-4 py-3">Link</th>
              <th class="px-4 py-3">Status</th>
              <th class="px-4 py-3 text-right">Cliques</th>
              <th class="px-4 py-3">Alterado</th>
              <th class="px-4 py-3 text-right"><span class="sr-only">Ações</span></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
            <tr v-for="link in visible" :key="link.slug"
              :class="link.slug === highlight && 'bg-brand-50/60 dark:bg-brand-500/5'">
              <td class="max-w-md px-4 py-3">
                <RouterLink :to="{ name: 'links-details', params: { slug: link.slug } }"
                  class="font-mono font-medium text-brand-700 hover:underline dark:text-brand-300">
                  {{ link.url.replace(/^https?:\/\//, '') }}
                </RouterLink>
                <p class="truncate text-slate-500 dark:text-slate-400" :title="link.target_url">
                  → {{ splitUrl(link.target_url).host }}<span class="text-slate-400 dark:text-slate-500">{{ splitUrl(link.target_url).rest }}</span>
                </p>
              </td>
              <td class="px-4 py-3"><StatusBadge :active="link.is_active" /></td>
              <td class="px-4 py-3 text-right font-medium tabular-nums">{{ formatNumber(link.clicks ?? 0) }}</td>
              <td class="whitespace-nowrap px-4 py-3 text-slate-500 dark:text-slate-400">{{ formatDate(link.last_change) }}</td>
              <td class="px-4 py-3">
                <div class="flex justify-end gap-1.5">
                  <CopyButton :text="link.url" compact />
                  <button type="button" class="btn-secondary px-2.5 py-1.5" :aria-label="`QR code de ${link.url}`" @click="qrFor = link">
                    <QrCodeIcon class="size-4" aria-hidden="true" />
                  </button>
                  <button type="button" class="btn-secondary px-2.5 py-1.5 text-rose-600 dark:text-rose-400"
                    :disabled="!link.is_active || busy === link.slug" @click="deactivate(link)">
                    Desativar
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- cartões (celular) -->
      <ul v-if="visible.length" class="space-y-3 md:hidden">
        <li v-for="link in visible" :key="link.slug" class="card space-y-3 p-4"
          :class="link.slug === highlight && 'ring-2 ring-brand-500/40'">
          <RouterLink :to="{ name: 'links-details', params: { slug: link.slug } }"
            class="block break-all font-mono font-medium text-brand-700 dark:text-brand-300">
            {{ link.url.replace(/^https?:\/\//, '') }}
          </RouterLink>
          <p class="truncate text-sm text-slate-500 dark:text-slate-400" :title="link.target_url">→ {{ link.target_url }}</p>
          <div class="flex items-center justify-between gap-2 text-sm">
            <span class="flex items-center gap-2 text-slate-500 dark:text-slate-400">
              <StatusBadge :active="link.is_active" />
              <span class="whitespace-nowrap"><strong class="font-semibold text-slate-900 tabular-nums dark:text-slate-100">{{ formatNumber(link.clicks ?? 0) }}</strong> cliques</span>
            </span>
            <div class="flex gap-1.5">
              <CopyButton :text="link.url" compact />
              <button type="button" class="btn-secondary px-2.5 py-1.5" :aria-label="`QR code de ${link.url}`" @click="qrFor = link">
                <QrCodeIcon class="size-4" aria-hidden="true" />
              </button>
              <button type="button" class="btn-secondary px-2.5 py-1.5 text-rose-600 dark:text-rose-400"
                :disabled="!link.is_active || busy === link.slug" @click="deactivate(link)">
                Desativar
              </button>
            </div>
          </div>
        </li>
      </ul>
    </template>

    <!-- QR code em diálogo -->
    <Dialog :open="!!qrFor" class="relative z-50" @close="qrFor = null">
      <div class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm" aria-hidden="true" />
      <div class="fixed inset-0 flex items-center justify-center p-4">
        <DialogPanel v-if="qrFor" class="card w-full max-w-xs p-6 text-center">
          <DialogTitle class="mb-1 font-semibold">QR code</DialogTitle>
          <p class="mb-4 break-all font-mono text-sm text-slate-500 dark:text-slate-400">{{ qrFor.url }}</p>
          <QrCode :value="qrFor.url" :filename="`qr-${qrFor.slug}`" />
          <button type="button" class="btn-ghost mt-4 w-full" @click="qrFor = null">Fechar</button>
        </DialogPanel>
      </div>
    </Dialog>
  </div>
</template>
