import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import AccountView from '../views/AccountView.vue'
import CreateLinkView from '../views/CreateLinkView.vue'
import LinkDetailsView from '../views/LinkDetailsView.vue'
import MyLinksView from '../views/MyLinksView.vue'
import NotFoundView from '../views/NotFoundView.vue'

const { isAuthenticated } = useAuth()

// Web history (URLs sem #): o nginx devolve o index.html para qualquer caminho
// que não seja arquivo, /api/ ou /r/ (try_files).
const router = createRouter({
  history: createWebHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    { path: '/', name: 'home', component: HomeView, meta: { title: 'Encurtador de links' } },
    { path: '/login', name: 'login', component: LoginView, meta: { guestOnly: true, title: 'Entrar' } },
    { path: '/register', name: 'register', component: RegisterView, meta: { guestOnly: true, title: 'Criar conta' } },
    { path: '/account', name: 'account', component: AccountView, meta: { requiresAuth: true, title: 'Minha conta' } },
    { path: '/links', name: 'links', component: MyLinksView, meta: { requiresAuth: true, title: 'Meus links' } },
    { path: '/links/new', name: 'links-new', component: CreateLinkView, meta: { requiresAuth: true, title: 'Criar link' } },
    {
      path: '/links/:slug',
      name: 'links-details',
      component: LinkDetailsView,
      meta: { requiresAuth: true, title: 'Detalhes do link' },
    },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView, meta: { title: 'Página não encontrada' } },
  ],
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guestOnly && isAuthenticated.value) {
    return { name: 'links' }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · URL Shortner` : 'URL Shortner'
})

export default router
