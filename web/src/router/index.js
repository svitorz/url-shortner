import { createRouter, createWebHashHistory } from 'vue-router'
import PingView from '../views/PingView.vue'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'

import AccountView from '../views/AccountView.vue'
import CreateLinkView from '../views/CreateLinkView.vue'
import LinkDetailsView from '../views/LinkDetailsView.vue'
import RedirectTesterView from '../views/RedirectTesterView.vue'

const router = createRouter({
  linkActiveClass: 'text-gray-900 border-indigo-600 ',
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'home', component: PingView },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/account', name: 'account', component: AccountView, meta: { requiresAuth: true } },
    { path: '/links/new', name: 'links-new', component: CreateLinkView, meta: { requiresAuth: true } },
    { path: '/links/:slug', name: 'links-details', component: LinkDetailsView, meta: { requiresAuth: true } },
    { path: '/test-redirect', name: 'redirect-tester', component: RedirectTesterView },
  ],
})

router.beforeEach((to, from, next) => {
  if (to.meta?.requiresAuth && !localStorage.getItem('token')) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router