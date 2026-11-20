<script setup>
import AppNavigation from './components/AppNavigation.vue'
import ToastHost from './components/ToastHost.vue'
import ConfirmHost from './components/ConfirmHost.vue'
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppNavigation />
    <main class="flex-1">
      <RouterView />
    </main>
    <footer class="border-t border-slate-200 py-6 text-center text-xs text-slate-500 dark:border-slate-800 dark:text-slate-500">
      URL Shortner · Go, Vue, PostgreSQL, Redis e nginx
    </footer>
  </div>
  <ToastHost />
  <ConfirmHost />
</template>
