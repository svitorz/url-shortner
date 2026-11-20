<script setup>
import { CheckCircleIcon, ExclamationTriangleIcon, InformationCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useToast } from '../composables/useToast'

const { toasts, dismiss } = useToast()
const icons = { success: CheckCircleIcon, error: ExclamationTriangleIcon, info: InformationCircleIcon }
const colors = {
  success: 'text-emerald-600 dark:text-emerald-400',
  error: 'text-rose-600 dark:text-rose-400',
  info: 'text-brand-600 dark:text-brand-400',
}
</script>

<template>
  <div aria-live="polite" class="pointer-events-none fixed inset-x-0 bottom-0 z-50 flex flex-col items-center gap-2 p-4 sm:items-end">
    <TransitionGroup enter-from-class="translate-y-2 opacity-0" enter-active-class="transition duration-200"
      leave-to-class="opacity-0" leave-active-class="transition duration-150">
      <div v-for="t in toasts" :key="t.id" :role="t.kind === 'error' ? 'alert' : 'status'"
        class="card pointer-events-auto flex w-full max-w-sm items-start gap-3 p-3 shadow-lg">
        <component :is="icons[t.kind]" class="size-5 shrink-0" :class="colors[t.kind]" aria-hidden="true" />
        <p class="flex-1 text-sm">{{ t.message }}</p>
        <button type="button" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200" @click="dismiss(t.id)">
          <span class="sr-only">Fechar</span>
          <XMarkIcon class="size-4" aria-hidden="true" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
