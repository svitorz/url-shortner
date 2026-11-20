<script setup>
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useConfirm } from '../composables/useConfirm'

const { state, answer } = useConfirm()
</script>

<template>
  <Dialog :open="!!state" class="relative z-50" @close="answer(false)">
    <div class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm" aria-hidden="true" />
    <div class="fixed inset-0 flex items-center justify-center p-4">
      <DialogPanel v-if="state" class="card w-full max-w-md p-6">
        <div class="flex gap-4">
          <div v-if="state.danger"
            class="flex size-10 shrink-0 items-center justify-center rounded-full bg-rose-100 dark:bg-rose-500/15">
            <ExclamationTriangleIcon class="size-5 text-rose-600 dark:text-rose-400" aria-hidden="true" />
          </div>
          <div>
            <DialogTitle class="text-base font-semibold">{{ state.title }}</DialogTitle>
            <p class="mt-2 text-sm text-slate-600 dark:text-slate-400">{{ state.message }}</p>
          </div>
        </div>
        <div class="mt-6 flex justify-end gap-2">
          <button type="button" class="btn-secondary" @click="answer(false)">Cancelar</button>
          <button type="button" :class="state.danger ? 'btn-danger' : 'btn-primary'" @click="answer(true)">
            {{ state.confirmLabel }}
          </button>
        </div>
      </DialogPanel>
    </div>
  </Dialog>
</template>
