export function formatDate(d) {
  const date = new Date(d)
  return isNaN(date) ? '' : date.toLocaleString('pt-BR')
}
