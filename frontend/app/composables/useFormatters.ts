/** Shared formatting helpers used in templates. */
export const useFormatters = () => {
  const formatPrice = (value: number, currency = 'TJS') => {
    if (typeof value !== 'number' || Number.isNaN(value)) return ''
    const formatter = new Intl.NumberFormat('ru-RU', {
      maximumFractionDigits: 2,
      minimumFractionDigits: 0,
    })
    const sym = currency === 'TJS' ? 'сом.' : currency
    return `${formatter.format(value)} ${sym}`
  }

  const formatDate = (value: string | Date) => {
    if (!value) return ''
    const d = typeof value === 'string' ? new Date(value) : value
    if (Number.isNaN(d.getTime())) return ''
    return d.toLocaleDateString('ru-RU', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }

  const phoneLink = (raw: string) => `tel:${(raw || '').replace(/[^0-9+]/g, '')}`
  const whatsappLink = (raw: string) =>
    `https://wa.me/${(raw || '').replace(/[^0-9]/g, '')}`

  return { formatPrice, formatDate, phoneLink, whatsappLink }
}
