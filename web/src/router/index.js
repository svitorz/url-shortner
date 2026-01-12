import { createRouter, createWebHistory } from 'vue-router'
import PingView from '../views/PingView.vue'
import RegisterView from '../views/RegisterView.vue'

const router = createRouter({
  linkActiveClass: 'text-gray-900',
  // link: 'text-gray-500 hover:border-gray-300 hover:text-gray-700',
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: PingView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
  ],
})

export default router
