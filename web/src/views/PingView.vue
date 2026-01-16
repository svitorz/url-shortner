<script setup>
import axios from 'axios';
import { onMounted, ref } from 'vue';

const info = ref("");
const loading = ref(true);
const hasError = ref(false);

async function testPing() {
  loading.value = true;
  hasError.value = false;

  try {
    const res = await axios.get('http://localhost:8000/api/ping');
    info.value = res.data?.message ?? '';
  } catch (err) {
    hasError.value = true;
  } finally {
    loading.value = false;
  }
}

onMounted(testPing);
</script>

<template>
  <main class="max-w-2xl mx-auto p-6">
    <div class="space-y-4 text-center">
      <h1 class="text-2xl font-bold">Teste de Ping</h1>
      <button
        class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
        @click="testPing"
      >
        Re-testar ping
      </button>

      <div v-if="loading" class="text-gray-600">Carregando...</div>
      <div v-else-if="hasError" class="text-red-600">Falha no ping</div>
      <div v-else class="text-green-700 font-medium">
        {{ info }}
      </div>
    </div>
  </main>
</template>