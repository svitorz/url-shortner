// Confirmação em diálogo próprio (no lugar do window.confirm). Uso:
//   if (await confirm({ title, message, confirmLabel, danger: true })) { ... }
import { ref } from 'vue'

const state = ref(null) // { title, message, confirmLabel, danger, resolve }

function confirm(opts) {
  return new Promise((resolve) => {
    state.value = { confirmLabel: 'Confirmar', danger: false, ...opts, resolve }
  })
}

function answer(ok) {
  state.value?.resolve(ok)
  state.value = null
}

export function useConfirm() {
  return { state, confirm, answer }
}
