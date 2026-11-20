import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import AccountView from '../views/AccountView.vue'
import CreateLinkView from '../views/CreateLinkView.vue'
import LinkDetailsView from '../views/LinkDetailsView.vue'
import MyLinksView from '../views/MyLinksView.vue'

const { isAuthenticated } = useAuth()

// Web history (URLs sem #): o nginx devolve o index.html para qualquer caminho
// que não seja arquivo, /api/ ou /r/ (try_files).
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: () => (isAuthenticated.value ? '/links' : '/login') },
    { path: '/login', name: 'login', component: LoginView, meta: { guestOnly: true } },
    { path: '/register', name: 'register', component: RegisterView, meta: { guestOnly: true } },
    { path: '/account', name: 'account', component: AccountView, meta: { requiresAuth: true } },
    { path: '/links', name: 'links', component: MyLinksView, meta: { requiresAuth: true } },
    { path: '/links/new', name: 'links-new', component: CreateLinkView, meta: { requiresAuth: true } },
    { path: '/links/:slug', name: 'links-details', component: LinkDetailsView, meta: { requiresAuth: true } },
    { path: '/:pathMatch(.*)*', redirect: '/' },
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

export default router
