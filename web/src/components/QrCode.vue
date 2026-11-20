<script setup>
// QR code gerado no navegador (pacote qrcode), sem serviço externo. Sempre
// escuro sobre branco, mesmo no tema escuro: é o que os leitores leem melhor.
import { ref, watchEffect } from 'vue'
import QRCode from 'qrcode'
import { ArrowDownTrayIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  value: { type: String, required: true },
  size: { type: Number, default: 192 },
  filename: { type: String, default: 'qrcode' },
})

const src = ref('')
const failed = ref(false)

watchEffect(async () => {
  try {
    src.value = await QRCode.toDataURL(props.value, {
      width: props.size * 2, // nítido em telas de alta densidade
      margin: 2,
      errorCorrectionLevel: 'M',
      color: { dark: '#0f172a', light: '#ffffff' },
    })
    failed.value = false
  } catch {
    failed.value = true
  }
})
</script>

<template>
  <div class="flex flex-col items-center gap-3">
    <div class="rounded-xl bg-white p-2 ring-1 ring-slate-200 dark:ring-slate-700">
      <img v-if="src && !failed" :src="src" :width="size" :height="size" :alt="`QR code de ${value}`" />
      <p v-else class="p-4 text-sm text-slate-500" :style="{ width: size + 'px' }">Não foi possível gerar o QR code.</p>
    </div>
    <a v-if="src && !failed" :href="src" :download="`${filename}.png`" class="btn-secondary text-xs">
      <ArrowDownTrayIcon class="size-4" aria-hidden="true" />
      Baixar PNG
    </a>
  </div>
</template>
