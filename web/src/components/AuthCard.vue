<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <div class="flex min-h-[calc(100vh-4rem)] items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
      <div class="mb-6 text-center">
        <img src="/logo.svg" alt="" class="mx-auto size-12" />
        <h1 class="mt-4 text-2xl font-bold tracking-tight">{{ title }}</h1>
        <p v-if="subtitle" class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ subtitle }}</p>
      </div>
      <div class="card p-6 sm:p-8">
        <slot />
      </div>
      <p class="mt-6 text-center text-sm text-slate-500 dark:text-slate-400">
        <slot name="footer" />
      </p>
    </div>
  </div>
</template>
