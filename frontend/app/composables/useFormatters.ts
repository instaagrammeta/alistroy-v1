export const useFormatters = () => {
  const formatPrice = (value: number, currency = 'TJS') => {
    if (typeof value !== 'number' || Number.isNaN(value)) return ''
    const f = new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 2, minimumFractionDigits: 0 })
    const sym = currency === 'TJS' ? 'сом.' : currency
    return `${f.format(value)} ${sym}`
  }

  const formatDate = (value?: string | Date | null) => {
    if (!value) return ''
    const d = typeof value === 'string' ? new Date(value) : value
    if (Number.isNaN(d.getTime())) return ''
    return d.toLocaleDateString('ru-RU', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }

  const formatDateTime = (value?: string | Date | null) => {
    if (!value) return ''
    const d = typeof value === 'string' ? new Date(value) : value
    if (Number.isNaN(d.getTime())) return ''
    return d.toLocaleString('ru-RU', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
  }

  const telLink = (raw: string) => `tel:${(raw || '').replace(/[^0-9+]/g, '')}`
  const waLink = (raw: string) => `https://wa.me/${(raw || '').replace(/[^0-9]/g, '')}`
  const tgLink = (raw: string) => {
    const v = (raw || '').trim()
    if (!v) return '#'
    if (v.startsWith('@')) return `https://t.me/${v.slice(1)}`
    if (/^https?:\/\//i.test(v)) return v
    return `https://t.me/${v.replace(/[^0-9]/g, '')}`
  }

  return { formatPrice, formatDate, formatDateTime, telLink, waLink, tgLink }
}
