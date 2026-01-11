dcl<script setup>
import { ref, onMounted } from 'vue';
const response = ref("");

const getPing = async () => {
  try {
    const response = await fetch('http://localhost:8000/ping');
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    response.value = await response.json();
  } catch (error) {
    console.error('There was a problem with the fetch operation:', error);
  }
};

onMounted(() => {
  getPing();
});

</script>

<template>
  <main>
    <div class="justify-center items-center">
      <button @click="getPing">
          ping
      </button>
      <span>
        <p id="response"> {{ response }}</p>
      </span>
    </div>
  </main>
</template>
