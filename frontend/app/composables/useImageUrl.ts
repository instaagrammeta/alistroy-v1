export const useImageUrl = () => {
  const config = useRuntimeConfig()
  const base = config.public.uploadsBase as string

  const resolve = (url?: string | null): string => {
    if (!url) return ''
    if (/^https?:\/\//i.test(url)) return url
    if (url.startsWith('/')) return url
    return `${base}/${url}`
  }
  return { resolve }
}
