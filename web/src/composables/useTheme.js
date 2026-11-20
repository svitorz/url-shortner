// Tema: 'system' (segue o sistema operacional), 'light' ou 'dark'. A escolha
// fica no localStorage; o script no index.html aplica a mesma regra antes do
// app carregar.
import { ref, watch } from 'vue'

const media = window.matchMedia('(prefers-color-scheme: dark)')

function read() {
  try {
    const t = localStorage.getItem('theme')
    return t === 'light' || t === 'dark' ? t : 'system'
  } catch {
    return 'system'
  }
}

const mode = ref(read())
const isDark = ref(false)

function apply() {
  isDark.value = mode.value === 'dark' || (mode.value === 'system' && media.matches)
  document.documentElement.classList.toggle('dark', isDark.value)
}

watch(mode, (m) => {
  try {
    if (m === 'system') localStorage.removeItem('theme')
    else localStorage.setItem('theme', m)
  } catch {
    /* sem localStorage: vale só para esta aba */
  }
  apply()
})
media.addEventListener('change', apply)
apply()

const order = ['system', 'light', 'dark']

export function useTheme() {
  // Alterna sistema → claro → escuro → sistema.
  const cycle = () => (mode.value = order[(order.indexOf(mode.value) + 1) % order.length])
  return { mode, isDark, cycle }
}
