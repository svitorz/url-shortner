import { ref, onMounted } from 'vue';

const isAuthenticated = ref(false);
const isExpired = ref(false)

function checkAuthStatus() {
  isAuthenticated.value = !!localStorage.getItem('token');
}

function login(token) {
  localStorage.setItem('token', token);
  isAuthenticated.value = true;
}

function logout() {
  localStorage.removeItem('token');
  isAuthenticated.value = false;
}

export function useAuth() {
  onMounted(checkAuthStatus);
  return { isAuthenticated, checkAuthStatus, login, logout };
}
