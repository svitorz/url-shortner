<script setup>
import { ref } from 'vue';
import { MinusIcon } from '@heroicons/vue/24/outline';
import axios from 'axios';
import router from '@/router';

const user = ref({
  email: '',
  password: '',
});

const loading = ref(false);
const response = ref('');
const hasError = ref(false);
const errorMessage = ref('');

const validationErrors = ref([]);
const validationSuccess = ref(true);

async function submitForm() {
  loading.value = true;
  hasError.value = false;
  response.value = '';
  errorMessage.value = '';

  try {
    const payload = {
      email: user.value.email.trim(),
      password: user.value.password,
    };

    const res = await axios.post('http://localhost:8000/api/login', payload);
    response.value = res.data?.message ?? 'User created';
    localStorage.setItem('token', res.data?.token);
  } catch (err) {
    hasError.value = true;
    errorMessage.value = err?.response?.data?.message ?? err?.message ?? 'Unknown error';
    console.error(err);
  } finally {
    loading.value = false;
    router.push('/')
  }
}

function resetFields() {
  user.value.email = '';
  user.value.password = '';

  loading.value = false;
  response.value = '';
  hasError.value = false;
  errorMessage.value = '';
  validationErrors.value = [];
  validationSuccess.value = true;
}
</script>


<template>
    <div class="p-10 m-12 w-full">
      <div class="grid grid-cols-1 gap-x-8 gap-y-8 py-10 md:grid-cols-3">
        <div class="px-4 sm:px-0">
          <h2 class="text-base/7 font-semibold text-gray-900">Personal Information</h2> 
          <p class="mt-1 text-sm/6 text-gray-600">Use a permanent address where you can access application.</p>
        </div>
        <div v-if="loading">
            <div class="animate-spin">
              <MinusIcon />
            </div>
        </div>
        <div v-else class="bg-white shadow-sm outline outline-1 outline-gray-900/5 sm:rounded-xl md:col-span-2">
          <div v-if="!validationSuccess">
              <div v-for="validation in validationErrors">
                <p class="text-red-600"> {{  validation  }} </p>
              </div>
          </div>
          <div v-if="hasError">
              <p class="text-red-600"> {{  error }} </p>
          </div>
          <div class="px-4 py-6 sm:p-8" v-else>
            <div class="grid max-w-2xl grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
              <div class="sm:col-span-4">
                <label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
                <div class="mt-2">
                  <input id="email" v-model="user.email" name="email" type="email" autocomplete="email" class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" />
                </div>
              </div>
              <div class="sm:col-span-4">
                <label for="email" class="block text-sm/6 font-medium text-gray-900">Password</label>
                <div class="mt-2">
                  <input id="password" v-model="user.password" name="password" type="password" class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" />
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center justify-end gap-x-6 border-t border-gray-900/10 px-4 py-4 sm:px-8">
            <button type="button" @click="resetFields" class="text-sm/6 font-semibold text-gray-900 border border-2 px-2 py-1 rounded-md hover:bg-gray-300 border-gray-300">Cancel</button>
            <button type="button" @click="submitForm" class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Log in</button>
          </div>
        </div>
      </div>

      <div v-if="response" class="mt-4 text-green-700">{{ response }}</div>
    </div>
</template>