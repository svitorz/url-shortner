// Copiar texto com feedback. Tenta navigator.clipboard (exige contexto seguro
// e permissão) e, se falhar, o execCommand antigo. Se os dois falharem, avisa
// para copiar à mão em vez de fingir que copiou.
import { ref } from 'vue'

function legacyCopy(text) {
  const el = document.createElement('textarea')
  el.value = text
  el.setAttribute('readonly', '')
  el.style.position = 'fixed'
  el.style.opacity = '0'
  document.body.appendChild(el)
  el.select()
  try {
    return document.execCommand('copy')
  } catch {
    return false
  } finally {
    document.body.removeChild(el)
  }
}

export function useCopy() {
  const copied = ref('') // texto copiado com sucesso
  const failed = ref('') // texto que não deu para copiar

  async function copy(text) {
    let ok = false
    if (navigator.clipboard && window.isSecureContext) {
      try {
        await navigator.clipboard.writeText(text)
        ok = true
      } catch {
        ok = false
      }
    }
    if (!ok) ok = legacyCopy(text)

    copied.value = ok ? text : ''
    failed.value = ok ? '' : text
    setTimeout(() => {
      if (copied.value === text) copied.value = ''
      if (failed.value === text) failed.value = ''
    }, 2500)
    return ok
  }

  // Rótulo do botão para um texto.
  const label = (text) => (copied.value === text ? 'Copiado!' : failed.value === text ? 'Copie manualmente' : 'Copiar')

  return { copied, failed, copy, label }
}
