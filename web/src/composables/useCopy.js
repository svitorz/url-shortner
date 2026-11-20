// Copiar texto com feedback. Tenta navigator.clipboard (exige contexto seguro
// e permissão) e, se falhar, o execCommand antigo. Devolve false se os dois
// falharem, para quem chama avisar (o CopyButton mostra um toast).
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
  const copied = ref('') // texto copiado com sucesso (some após 2,5 s)

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
    if (ok) setTimeout(() => copied.value === text && (copied.value = ''), 2500)
    return ok
  }

  return { copied, copy }
}
