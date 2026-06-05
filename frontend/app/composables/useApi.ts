import type { FetchContext, FetchOptions } from 'ofetch'

/**
 * Returns a `$fetch` instance pre-configured with the API base URL and
 * an Authorization header sourced from the auth store.
 *
 * On a 401 response we attempt a refresh-token flow once, then retry.
 */
export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  // Lazy access to avoid circular imports during SSR.
  const getStores = async () => {
    const { useAuthStore } = await import('~/stores/auth')
    return { auth: useAuthStore() }
  }

  let isRefreshing = false

  const onRequest = async ({ options }: FetchContext) => {
    const { auth } = await getStores()
    if (auth.accessToken) {
      const headers = new Headers(options.headers as HeadersInit | undefined)
      headers.set('Authorization', `Bearer ${auth.accessToken}`)
      options.headers = headers
    }
  }

  const onResponseError = async ({ request, response, options }: FetchContext) => {
    if (!response) return
    if (response.status !== 401) return
    if (isRefreshing) return
    const { auth } = await getStores()
    if (!auth.refreshToken) return

    isRefreshing = true
    try {
      const ok = await auth.tryRefresh()
      if (!ok) {
        await auth.logout({ silent: true })
        return
      }
      // Retry the original request once with the new token.
      const headers = new Headers(options.headers as HeadersInit | undefined)
      headers.set('Authorization', `Bearer ${auth.accessToken}`)
      // Mutate the response so callers see the retried result.
      const retried = await $fetch.raw(request as string, {
        ...options,
        headers,
        baseURL,
      } as FetchOptions)
      ;(response as any)._data = retried._data
      ;(response as any).status = retried.status
    } finally {
      isRefreshing = false
    }
  }

  return $fetch.create({
    baseURL,
    onRequest,
    onResponseError,
  })
}
