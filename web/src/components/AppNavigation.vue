<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import {
  Bars3Icon,
  XMarkIcon,
  SunIcon,
  MoonIcon,
  ComputerDesktopIcon,
  ArrowRightStartOnRectangleIcon,
} from '@heroicons/vue/24/outline'
import { useAuth } from '../composables/useAuth'
import { useTheme } from '../composables/useTheme'

const { isAuthenticated, logout } = useAuth()
const { mode, cycle } = useTheme()
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

const themeIcon = computed(() => ({ system: ComputerDesktopIcon, light: SunIcon, dark: MoonIcon })[mode.value])
const themeLabel = computed(
  () => ({ system: 'Tema: do sistema', light: 'Tema: claro', dark: 'Tema: escuro' })[mode.value],
)

const isActive = (to) => route.path === to

function signOut(close) {
  logout()
  close?.()
  router.push({ name: 'login' })
}
</script>

<template>
  <Disclosure as="header" v-slot="{ open, close }"
    class="sticky top-0 z-30 border-b border-slate-200 bg-white/80 backdrop-blur dark:border-slate-800 dark:bg-slate-950/80">
    <nav class="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6">
      <RouterLink to="/" class="flex items-center gap-2 font-bold tracking-tight">
        <img src="/logo.svg" alt="" class="size-8" />
        <span class="text-slate-900 dark:text-white">URL<span class="text-brand-600 dark:text-brand-400">Shortner</span></span>
      </RouterLink>

      <div class="hidden items-center gap-1 sm:flex">
        <RouterLink v-for="item in links" :key="item.to" :to="item.to"
          :class="isActive(item.to)
            ? 'bg-brand-50 text-brand-700 dark:bg-brand-500/10 dark:text-brand-300'
            : 'text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800'"
          class="rounded-lg px-3 py-2 text-sm font-medium">
          {{ item.label }}
        </RouterLink>
        <button type="button" class="btn-ghost px-2.5" :title="themeLabel" :aria-label="themeLabel" @click="cycle">
          <component :is="themeIcon" class="size-5" aria-hidden="true" />
        </button>
        <button v-if="isAuthenticated" type="button" class="btn-ghost" @click="signOut()">
          <ArrowRightStartOnRectangleIcon class="size-5" aria-hidden="true" />
          Sair
        </button>
      </div>

      <div class="flex items-center gap-1 sm:hidden">
        <button type="button" class="btn-ghost px-2.5" :aria-label="themeLabel" @click="cycle">
          <component :is="themeIcon" class="size-5" aria-hidden="true" />
        </button>
        <DisclosureButton class="btn-ghost px-2.5">
          <span class="sr-only">Abrir menu</span>
          <Bars3Icon v-if="!open" class="size-6" aria-hidden="true" />
          <XMarkIcon v-else class="size-6" aria-hidden="true" />
        </DisclosureButton>
      </div>
    </nav>

    <DisclosurePanel class="border-t border-slate-200 px-4 pb-3 pt-2 sm:hidden dark:border-slate-800">
      <RouterLink v-for="item in links" :key="item.to" :to="item.to" @click="close()"
        :class="isActive(item.to) ? 'bg-brand-50 text-brand-700 dark:bg-brand-500/10 dark:text-brand-300' : 'text-slate-700 dark:text-slate-200'"
        class="block rounded-lg px-3 py-2 text-base font-medium">
        {{ item.label }}
      </RouterLink>
      <button v-if="isAuthenticated" type="button"
        class="block w-full rounded-lg px-3 py-2 text-left text-base font-medium text-slate-700 dark:text-slate-200"
        @click="signOut(close)">
        Sair
      </button>
    </DisclosurePanel>
  </Disclosure>
</template>
