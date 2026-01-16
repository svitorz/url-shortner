<script setup>
import { ref } from 'vue';
import api from '../services/API';
import LogoutButton from '@/components/LogoutButton.vue';

const form = ref({ id: '', full_name: '', email: '' });
const loadingUpdate = ref(false);
const loadingDelete = ref(false);
const msg = ref('');
const err = ref('');

async function updateUser() {
  loadingUpdate.value = true;
  msg.value = ''; err.value = '';
  if (!form.value.id) {
    err.value = 'Informe o ID do usuário';
    loadingUpdate.value = false;
    return;
  }
  try {
    const res = await api.put(`/user/${form.value.id}`, {
      full_name: form.value.full_name.trim(),
      email: form.value.email.trim(),
    });
    msg.value = res.data?.message ?? 'Usuário atualizado';
  } catch (e) {
    err.value = e?.response?.data?.error ?? e?.message ?? 'Erro ao atualizar usuário';
  } finally {
    loadingUpdate.value = false;
  }
}

</script>

<template>
  <main class="max-w-lg mx-auto p-6 space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold">Minha conta</h1>
      <LogoutButton />
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Nome completo</label>
        <input v-model="form.full_name" type="text" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Email</label>
        <input v-model="form.email" type="email" class="mt-1 block w-full rounded border border-blue-800 px-3 py-2" />
      </div>

      <div class="flex gap-3">
        <button
          class="rounded bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:opacity-50"
          :disabled="loadingUpdate"
          @click="updateUser"
        >
          {{ loadingUpdate ? 'Atualizando...' : 'Atualizar' }}
        </button>
      </div>

      <p v-if="msg" class="text-green-600 text-sm">{{ msg }}</p>
      <p v-if="err" class="text-red-600 text-sm">{{ err }}</p>
    </div>
  </main>
</template>