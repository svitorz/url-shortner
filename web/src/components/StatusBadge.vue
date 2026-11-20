<script setup>
defineProps({ active: { type: Boolean, required: true } })
</script>

<template>
  <span :class="active ? 'badge-active' : 'badge-inactive'">
    <span class="size-1.5 rounded-full" :class="active ? 'bg-emerald-500' : 'bg-slate-400'" aria-hidden="true" />
    {{ active ? 'Ativo' : 'Desativado' }}
  </span>
</template>
