<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import { Bars3Icon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAuth } from '../composables/useAuth'
import logo from '@/assets/logo.svg'

const { isAuthenticated, logout } = useAuth()
const route = useRoute()
const router = useRouter()

const links = computed(() =>
  isAuthenticated.value
    ? [
        { label: 'Meus links', to: '/links' },
        { label: 'Criar link', to: '/links/new' },
        { label: 'Minha conta', to: '/account' },
      ]
    : [
        { label: 'Entrar', to: '/login' },
        { label: 'Criar conta', to: '/register' },
      ],
)

const isActive = (to) => route.path === to

function signOut(close) {
  logout()
  close?.()
  router.push({ name: 'login' })
}
</script>

<template>
  <Disclosure as="nav" class="bg-white shadow" v-slot="{ open, close }">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-between">
        <RouterLink to="/" class="flex items-center gap-2 font-semibold text-gray-900">
          <img :src="logo" alt="" class="h-7 w-auto" />
          <span>URL Shortner</span>
        </RouterLink>

        <div class="hidden items-center gap-6 sm:flex">
          <RouterLink v-for="item in links" :key="item.to" :to="item.to"
            :class="isActive(item.to) ? 'border-indigo-600 text-gray-900' : 'border-transparent text-gray-500 hover:text-gray-700'"
            class="border-b-2 py-5 text-sm font-medium">
            {{ item.label }}
          </RouterLink>
          <button v-if="isAuthenticated" type="button"
            class="text-sm font-medium text-gray-500 hover:text-gray-700" @click="signOut()">
            Sair
          </button>
        </div>

        <DisclosureButton
          class="inline-flex items-center justify-center rounded-md p-2 text-gray-500 hover:bg-gray-100 sm:hidden">
          <span class="sr-only">Abrir menu</span>
          <Bars3Icon v-if="!open" class="size-6" aria-hidden="true" />
          <XMarkIcon v-else class="size-6" aria-hidden="true" />
        </DisclosureButton>
      </div>
    </div>

    <DisclosurePanel class="border-t border-gray-200 pb-3 pt-2 sm:hidden">
      <RouterLink v-for="item in links" :key="item.to" :to="item.to" @click="close()"
        :class="isActive(item.to) ? 'border-indigo-600 bg-indigo-50 text-indigo-700' : 'border-transparent text-gray-600 hover:bg-gray-50'"
        class="block border-l-4 py-2 pl-3 pr-4 text-base font-medium">
        {{ item.label }}
      </RouterLink>
      <button v-if="isAuthenticated" type="button"
        class="block w-full border-l-4 border-transparent py-2 pl-3 pr-4 text-left text-base font-medium text-gray-600 hover:bg-gray-50"
        @click="signOut(close)">
        Sair
      </button>
    </DisclosurePanel>
  </Disclosure>
</template>
