<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], required: true },
  icon: { type: [Object, Function], required: true },
  tone: { type: String, default: 'brand' },
})
const tones = {
  brand: 'bg-brand-50 text-brand-700 dark:bg-brand-500/10 dark:text-brand-300',
  emerald: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300',
  slate: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-300',
  amber: 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-300',
}
</script>

<template>
  <div class="card flex items-center gap-4 p-4">
    <div class="flex size-11 items-center justify-center rounded-xl" :class="tones[tone]">
      <component :is="icon" class="size-6" aria-hidden="true" />
    </div>
    <div>
      <p class="text-sm text-slate-500 dark:text-slate-400">{{ label }}</p>
      <p class="text-2xl font-semibold tabular-nums">{{ value }}</p>
    </div>
  </div>
</template>
