import { createRouter, createWebHashHistory } from 'vue-router'
import PingView from '../views/PingView.vue'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  linkActiveClass: 'text-gray-900 border-indigo-600 ',
  history: createWebHashHistory(),
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
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
  ],
})

export default router
