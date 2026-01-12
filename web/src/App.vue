<script setup>
import Navigation from './components/Navigation.vue';
</script>

<template>
  <nav>
    <Navigation />
  </nav>
  <main>
      <RouterView />
  </main>
</template>