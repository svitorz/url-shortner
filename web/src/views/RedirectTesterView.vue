<script setup>
import { ref } from 'vue';

const slug = ref('');
const base = 'http://localhost:8000/api';

function go() {
  if (!slug.value) return;
  // Navega para o endpoint de redirecionamento; o backend deve responder com 3xx
  window.location.href = `${base}/r/${slug.value}`;
}
</script>

<template>
  <main class="max-w-md mx-auto p-6 space-y-4">
    <h1 class="text-xl font-semibold">Testar redirecionamento</h1>
    <div>
      <label class="block text-sm font-medium text-gray-700">Slug</label>
      <input v-model="slug" type="text" class="mt-1 block w-full rounded border-gray-300 px-3 py-2" />
    </div>
    <button
      class="w-full rounded bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
      :disabled="!slug"
      @click="go"
    >
      Ir
    </button>
  </main>
</template>