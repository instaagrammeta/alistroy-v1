import type { FetchContext, FetchOptions } from 'ofetch'

/**
 * Authenticated $fetch instance. Attaches the bearer token, and on a 401
 * attempts a single refresh-token retry.
 */
export const useApi = () => {
  const baseURL = useApiBase()

  const getAuth = async () => {
    const { useAuthStore } = await import('~/stores/auth')
    return useAuthStore()
  }

  let refreshing = false

  const onRequest = async ({ options }: FetchContext) => {
    const auth = await getAuth()
    if (auth.accessToken) {
      const headers = new Headers(options.headers as HeadersInit | undefined)
      headers.set('Authorization', `Bearer ${auth.accessToken}`)
      options.headers = headers
    }
  }

  const onResponseError = async ({ request, response, options }: FetchContext) => {
    if (!response || response.status !== 401 || refreshing) return
    const auth = await getAuth()
    if (!auth.refreshToken) return
    refreshing = true
    try {
      const ok = await auth.tryRefresh()
      if (!ok) {
        await auth.logout({ silent: true })
        return
      }
      const headers = new Headers(options.headers as HeadersInit | undefined)
      headers.set('Authorization', `Bearer ${auth.accessToken}`)
      const retried = await $fetch.raw(request as string, {
        ...options,
        headers,
        baseURL,
      } as FetchOptions)
      ;(response as any)._data = retried._data
      ;(response as any).status = retried.status
    } finally {
      refreshing = false
    }
  }

  return $fetch.create({ baseURL, onRequest, onResponseError })
}
