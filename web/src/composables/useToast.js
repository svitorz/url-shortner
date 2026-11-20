// Avisos rápidos (canto da tela) para sucesso e erro de ações.
import { ref } from 'vue'

const toasts = ref([])
let nextId = 1

function push(kind, message, ms = 3500) {
  const id = nextId++
  toasts.value.push({ id, kind, message })
  setTimeout(() => dismiss(id), ms)
}

function dismiss(id) {
  toasts.value = toasts.value.filter((t) => t.id !== id)
}

export function useToast() {
  return {
    toasts,
    dismiss,
    success: (m) => push('success', m),
    error: (m) => push('error', m, 6000),
    info: (m) => push('info', m),
  }
}
