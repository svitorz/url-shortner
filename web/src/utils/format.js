export function formatDate(d) {
  const date = new Date(d)
  return isNaN(date) ? '' : date.toLocaleString('pt-BR', { dateStyle: 'short', timeStyle: 'short' })
}

export function formatNumber(n) {
  return Number(n ?? 0).toLocaleString('pt-BR')
}
