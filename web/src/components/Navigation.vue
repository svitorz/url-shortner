<script setup>
import { Disclosure, DisclosureButton, Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { Bars3Icon, ChevronDownIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useAuth } from '../composables/useAuth';

const { isAuthenticated, logout } = useAuth();

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
            <img class="h-8 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
              alt="Your Company" />
          </div>
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <RouterLink class="inline-flex items-center px-1 pt-1 text-sm font-medium" to="/">Ping</RouterLink>
            <RouterLink class="inline-flex items-center px-1 pt-1 text-sm font-medium" to="/create">Create a new Link
            </RouterLink>
          </div>
        </div>
        <div class="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
          <div v-if="isAuthenticated">
            <Menu as="div" class="relative ml-3">
              <MenuButton
                class="relative flex rounded-full focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                <span class="absolute -inset-1.5"></span>
                <span class="sr-only">Open user menu</span>
                <ChevronDownIcon class="block size-6" />
              </MenuButton>

              <transition enter-active-class="transition ease-out duration-200"
                enter-from-class="transform opacity-0 scale-95" enter-to-class="transform scale-100"
                leave-active-class="transition ease-in duration-75" leave-from-class="transform scale-100"
                leave-to-class="transform opacity-0 scale-95">
                <MenuItems
                  class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg outline outline-1 outline-black/5">
                  <MenuItem v-slot="{ active }">
                  <a href="#"
                    :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">Your
                    profile</a>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                  <a href="#" @click="logout"
                    :class="[active ? 'bg-gray-100 outline-none' : '', 'block px-4 py-2 text-sm text-gray-700']">Sign
                    out</a>
                  </MenuItem>
                </MenuItems>
              </transition>
            </Menu>
          </div>
          <div v-else>
            <RouterLink class="inline-flex items-center px-1 pt-1 text-sm font-medium" to="/login">Login</RouterLink>
            <RouterLink class="inline-flex items-center px-1 pt-1 text-sm font-medium" to="/register">Register
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>
</template>
