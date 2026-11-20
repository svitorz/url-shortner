<script setup>
import { ClipboardDocumentIcon, CheckIcon } from '@heroicons/vue/24/outline'
import { useCopy } from '../composables/useCopy'
import { useToast } from '../composables/useToast'

const props = defineProps({
  text: { type: String, required: true },
  compact: { type: Boolean, default: false },
})

const { copy, copied } = useCopy()
const toast = useToast()

async function onClick() {
  if (!(await copy(props.text))) {
    toast.error('O navegador bloqueou a cópia. Selecione o link e copie manualmente.')
  }
}
</script>

<template>
  <button type="button" class="btn-secondary" :class="compact ? 'px-2.5 py-1.5' : ''"
    :aria-label="`Copiar ${text}`" @click="onClick">
    <CheckIcon v-if="copied === text" class="size-4 text-emerald-600 dark:text-emerald-400" aria-hidden="true" />
    <ClipboardDocumentIcon v-else class="size-4" aria-hidden="true" />
    <span v-if="!compact">{{ copied === text ? 'Copiado!' : 'Copiar' }}</span>
  </button>
</template>
