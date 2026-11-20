<script setup>
import { BoltIcon, ChartBarIcon, QrCodeIcon, ShieldCheckIcon, ArrowRightIcon } from '@heroicons/vue/24/outline'
import { useAuth } from '../composables/useAuth'

const { isAuthenticated } = useAuth()

const features = [
  { icon: BoltIcon, title: 'Redirect rápido', text: 'Links curtos servidos com cache em memória: o banco nem é consultado na maioria dos acessos.' },
  { icon: ChartBarIcon, title: 'Cliques por link', text: 'Veja quantas vezes cada link foi aberto, sem precisar de ferramenta externa.' },
  { icon: QrCodeIcon, title: 'QR code pronto', text: 'Gere e baixe o QR code de qualquer link curto para colocar em material impresso.' },
  { icon: ShieldCheckIcon, title: 'Sob seu controle', text: 'Desative um link a qualquer momento: quem acessar vê que ele foi desativado.' },
]

const steps = [
  { n: 1, title: 'Crie sua conta', text: 'Leva alguns segundos e é gratuito.' },
  { n: 2, title: 'Cole a URL longa', text: 'Receba na hora um endereço curto, fácil de compartilhar.' },
  { n: 3, title: 'Acompanhe', text: 'Copie, gere o QR code e veja os cliques no seu painel.' },
]
</script>

<template>
  <div>
    <!-- hero -->
    <section class="relative overflow-hidden">
      <div class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(60rem_30rem_at_50%_-10%,var(--color-brand-200),transparent)] opacity-60 dark:bg-[radial-gradient(60rem_30rem_at_50%_-10%,var(--color-brand-900),transparent)]" />
      <div class="mx-auto max-w-4xl px-4 pb-16 pt-20 text-center sm:px-6 sm:pt-28">
        <span class="inline-flex items-center gap-2 rounded-full border border-brand-200 bg-white/70 px-3 py-1 text-xs font-medium text-brand-700 dark:border-brand-800 dark:bg-slate-900/70 dark:text-brand-300">
          Links curtos, cliques e QR code
        </span>
        <h1 class="mt-6 text-4xl font-bold tracking-tight text-balance sm:text-6xl">
          Encurte. Compartilhe. <span class="text-brand-600 dark:text-brand-400">Acompanhe.</span>
        </h1>
        <p class="mx-auto mt-5 max-w-2xl text-lg text-pretty text-slate-600 dark:text-slate-400">
          Transforme URLs longas em links curtos, acompanhe quantas vezes foram abertos e desative quando quiser.
        </p>
        <div class="mt-8 flex flex-wrap justify-center gap-3">
          <template v-if="isAuthenticated">
            <RouterLink to="/links" class="btn-primary px-5 py-3 text-base">
              Ir para meus links <ArrowRightIcon class="size-5" aria-hidden="true" />
            </RouterLink>
          </template>
          <template v-else>
            <RouterLink to="/register" class="btn-primary px-5 py-3 text-base">
              Criar conta grátis <ArrowRightIcon class="size-5" aria-hidden="true" />
            </RouterLink>
            <RouterLink to="/login" class="btn-secondary px-5 py-3 text-base">Já tenho conta</RouterLink>
          </template>
        </div>

        <!-- exemplo visual -->
        <div class="card mx-auto mt-14 max-w-xl p-4 text-left">
          <p class="truncate text-sm text-slate-400 line-through decoration-slate-300 dark:decoration-slate-600">
            https://www.exemplo.com.br/blog/2026/11/um-artigo-com-um-endereco-muito-longo?utm_source=newsletter
          </p>
          <div class="mt-2 flex items-center justify-between gap-3">
            <p class="font-mono text-base font-semibold text-brand-700 sm:text-lg dark:text-brand-300">seu.dominio/r/aZ3k9Qp</p>
            <span class="badge-active shrink-0 whitespace-nowrap">128 cliques</span>
          </div>
        </div>
      </div>
    </section>

    <!-- recursos -->
    <section class="mx-auto max-w-6xl px-4 py-16 sm:px-6">
      <h2 class="text-center text-2xl font-bold tracking-tight">O que você ganha</h2>
      <div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div v-for="f in features" :key="f.title" class="card p-5">
          <div class="flex size-10 items-center justify-center rounded-xl bg-brand-50 dark:bg-brand-500/10">
            <component :is="f.icon" class="size-5 text-brand-600 dark:text-brand-400" aria-hidden="true" />
          </div>
          <h3 class="mt-4 font-semibold">{{ f.title }}</h3>
          <p class="mt-1 text-sm text-slate-600 dark:text-slate-400">{{ f.text }}</p>
        </div>
      </div>
    </section>

    <!-- como funciona -->
    <section class="mx-auto max-w-4xl px-4 pb-20 sm:px-6">
      <h2 class="text-center text-2xl font-bold tracking-tight">Como funciona</h2>
      <ol class="mt-10 grid gap-6 sm:grid-cols-3">
        <li v-for="s in steps" :key="s.n" class="text-center">
          <span class="mx-auto flex size-10 items-center justify-center rounded-full bg-brand-600 font-bold text-white">{{ s.n }}</span>
          <h3 class="mt-3 font-semibold">{{ s.title }}</h3>
          <p class="mt-1 text-sm text-slate-600 dark:text-slate-400">{{ s.text }}</p>
        </li>
      </ol>
    </section>
  </div>
</template>
