<script setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { Bars3Icon, ChevronDownIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAuth } from '../composables/useAuth';
import { useRoute } from 'vue-router';
import { computed } from 'vue';

const { isAuthenticated, logout } = useAuth();
const route = useRoute();

function navLinkClass(path) {
  const active = route.path === path;
  return [
    'inline-flex items-center px-1 pt-1 text-sm font-medium border-b-2',
    active ? 'border-indigo-600 text-gray-900' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
  ]
}

const leftLinks = computed(() => ([
  { label: 'Ping', to: '/' },
  { label: 'Criar link', to: '/links/new' },
  { label: 'Testar redirect', to: '/test-redirect' },
  { label: 'Meus links', to: "/links"}
]));
</script>

<template>
  <Disclosure as="nav" class="relative bg-white shadow" v-slot="{ open }">
    <div class="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
      <div class="relative flex h-16 justify-between">
        <div class="absolute inset-y-0 left-0 flex items-center sm:hidden">
          <DisclosureButton
            class="relative inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-600">
            <span class="absolute -inset-0.5"></span>
            <span class="sr-only">Open main menu</span>
            <Bars3Icon v-if="!open" class="block size-6" aria-hidden="true" />
            <XMarkIcon v-else class="block size-6" aria-hidden="true" />
          </DisclosureButton>
        </div>

        <div class="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
          <div class="flex shrink-0 items-center">
            <img class="h-8 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600" alt="Your Company" />
          </div>
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <RouterLink
              v-for="item in leftLinks"
              :key="item.to"
              :to="item.to"
              :class="navLinkClass(item.to)"
            >
              {{ item.label }}
            </RouterLink>
          </div>
        </div>

        <!-- Right side: auth/menu -->
        <div class="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
          <div v-if="isAuthenticated" class="flex items-center">
            <Menu as="div" class="relative ml-3">
              <MenuButton
                class="relative flex items-center gap-2 rounded-full px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                <span class="sr-only">Open user menu</span>
                <span>Conta</span>
                <ChevronDownIcon class="block size-5 text-gray-500" />
              </MenuButton>

              <transition
                enter-active-class="transition ease-out duration-200"
                enter-from-class="transform opacity-0 scale-95"
                enter-to-class="transform opacity-100 scale-100"
                leave-active-class="transition ease-in duration-75"
                leave-from-class="transform opacity-100 scale-100"
                leave-to-class="transform opacity-0 scale-95">
                <MenuItems
                  class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg outline outline-1 outline-black/5">
                  <MenuItem v-slot="{ active }">
                    <RouterLink
                      to="/account"
                      :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">
                      Minha conta
                    </RouterLink>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                    <RouterLink
                      to="/links/new"
                      :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">
                      Criar link
                    </RouterLink>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                    <RouterLink
                      to="/links"
                      :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">
                      Meus links
                    </RouterLink>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                    <RouterLink
                      to="/test-redirect"
                      :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">
                      Testar redirect
                    </RouterLink>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                    <a
                      href="#"
                      @click.prevent="logout"
                      :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">
                      Sair
                    </a>
                  </MenuItem>
                </MenuItems>
              </transition>
            </Menu>
          </div>
          <div v-else class="flex items-center gap-4">
            <RouterLink :class="navLinkClass('/login')" to="/login">Login</RouterLink>
            <RouterLink :class="navLinkClass('/register')" to="/register">Register</RouterLink>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile menu -->
    <DisclosurePanel class="sm:hidden">
      <div class="space-y-1 pb-3 pt-2">
        <RouterLink
          v-for="item in leftLinks"
          :key="item.to"
          :to="item.to"
          class="block border-l-4 py-2 pl-3 pr-4 text-base font-medium"
          :class="route.path === item.to ? 'border-indigo-600 bg-indigo-50 text-indigo-700' : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'"
        >
          {{ item.label }}
        </RouterLink>
      </div>

      <div class="border-t border-gray-200 pb-3 pt-4">
        <div v-if="isAuthenticated" class="space-y-1">
          <RouterLink
            to="/account"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Minha conta
          </RouterLink>
          <RouterLink
            to="/links/new"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Criar link
          </RouterLink>
          <RouterLink
            to="/links"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Meus links
          </RouterLink>
          <RouterLink
            to="/test-redirect"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Testar redirect
          </RouterLink>
          <a
            href="#"
            @click.prevent="logout"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Sair
          </a>
        </div>
        <div v-else class="space-y-1">
          <RouterLink
            to="/login"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Login
          </RouterLink>
          <RouterLink
            to="/register"
            class="block px-4 py-2 text-base font-medium text-gray-700 hover:bg-gray-50">
            Register
          </RouterLink>
        </div>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>