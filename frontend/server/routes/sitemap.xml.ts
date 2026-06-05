// Dynamic sitemap.xml generator — fetches approved products, sellers, and
// categories from the backend and emits a flat XML sitemap.

interface IdSlug { slug: string; updated_at?: string }

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = useApiBase()
  const siteUrl = (config.public.siteUrl as string).replace(/\/$/, '')

  const safeFetch = async <T>(path: string): Promise<T | null> => {
    try {
      return await $fetch<T>(path, { baseURL: apiBase })
    } catch {
      return null
    }
  }

  const productsRes = await safeFetch<{ data: IdSlug[]; pagination?: any }>(
    '/products?page_size=1000'
  )
  const sellersRes = await safeFetch<{ data: IdSlug[] }>('/sellers?page_size=200')
  const catsRes = await safeFetch<{ data: IdSlug[] }>('/categories')

  const urls: string[] = [
    `${siteUrl}/`,
    `${siteUrl}/products`,
    `${siteUrl}/categories`,
    `${siteUrl}/sellers`,
  ]
  for (const c of catsRes?.data || []) urls.push(`${siteUrl}/categories/${c.slug}`)
  for (const s of sellersRes?.data || []) urls.push(`${siteUrl}/sellers/${s.slug}`)
  for (const p of productsRes?.data || []) urls.push(`${siteUrl}/products/${p.slug}`)

  const xml =
    `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n` +
    urls.map((u) => `  <url><loc>${u}</loc></url>`).join('\n') +
    `\n</urlset>\n`

  setHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  return xml
})
