/** Resolves an upload-relative URL into a full URL the browser can load. */
export const useImageUrl = () => {
  const config = useRuntimeConfig()
  const base = config.public.uploadsBase

  const resolve = (url?: string | null): string => {
    if (!url) return ''
    if (/^https?:\/\//i.test(url)) return url
    if (url.startsWith('/')) return url
    return `${base}/${url}`
  }
  return { resolve }
}
