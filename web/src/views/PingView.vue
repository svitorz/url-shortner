<script setup>
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { SlashIcon } from '@heroicons/vue/24/outline';

const info = ref("");
const loading = ref(true);
const hasError = ref(false);

async function testPing() {
  loading.value = true;
  hasError.value = false;

  try {
    const res = await axios.get('http://localhost:8000/api/ping');
    info.value = res.data?.message ?? '';
    console.log("response:",res)
    console.log("info:",info.value)
  } catch (err) {
    console.error(err);
    hasError.value = true;
  } finally {
    loading.value = false;
  }
}

onMounted(testPing);
</script>

<template>
  <main>
    <div class="justify-center items-center text-center p-5 w-full">
      <div class="w-1/2">
        <button class="rounded bg-indigo-600 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600" 
                @click="testPing"
        >
          Teste ping
        </button>
      </div>
      <div class="w-1/2">
        <section v-if="hasError">
          <p>Perdão, ocorreu um erro na nossa API.</p>
        </section>

        <section v-else>
          <div v-if="loading">
            <SlashIcon class="animate-spin block size-6" />
          <p>Carregando...</p>
          </div>

          <div v-else>
            <p>Resposta: {{ info }}</p>
          </div>
        </section>
      </div>
    </div>
  </main>
</template>