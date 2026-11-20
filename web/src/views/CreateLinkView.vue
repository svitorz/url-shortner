<script setup>
import { ref } from 'vue'
import { CheckCircleIcon } from '@heroicons/vue/24/outline'
import QuickCreate from '../components/QuickCreate.vue'
import CopyButton from '../components/CopyButton.vue'
import QrCode from '../components/QrCode.vue'

const created = ref(null) // { slug, url, target_url }
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-6 px-4 py-10 sm:px-6">
    <div>
      <h1 class="text-2xl font-bold tracking-tight">Criar link curto</h1>
      <p class="text-sm text-slate-500 dark:text-slate-400">Só endereços http:// ou https://, até 2048 caracteres.</p>
    </div>

    <div class="card p-5">
      <QuickCreate autofocus @created="created = $event" />
    </div>

    <div v-if="created" class="card overflow-hidden">
      <div class="flex items-center gap-2 border-b border-slate-200 bg-emerald-50 px-5 py-3 text-sm font-medium text-emerald-800 dark:border-slate-800 dark:bg-emerald-500/10 dark:text-emerald-300">
        <CheckCircleIcon class="size-5" aria-hidden="true" />
        Link criado
      </div>
      <div class="grid gap-6 p-5 sm:grid-cols-[1fr_auto] sm:items-center">
        <div class="min-w-0 space-y-3">
          <p class="truncate text-sm text-slate-500 dark:text-slate-400" :title="created.target_url">→ {{ created.target_url }}</p>
          <div class="flex items-center gap-2">
            <input :value="created.url" readonly aria-label="URL curta" class="input font-mono" @focus="$event.target.select()" />
            <CopyButton :text="created.url" />
          </div>
          <div class="flex flex-wrap gap-2">
            <RouterLink :to="{ name: 'links-details', params: { slug: created.slug } }" class="btn-secondary">Ver detalhes</RouterLink>
            <RouterLink to="/links" class="btn-ghost">Meus links</RouterLink>
          </div>
        </div>
        <QrCode :value="created.url" :size="140" :filename="`qr-${created.slug}`" />
      </div>
    </div>
  </div>
</template>
