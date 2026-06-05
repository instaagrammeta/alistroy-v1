/**
 * Returns the correct API base URL for the current execution context.
 *
 * - On the server (SSR): a full URL pointing at the backend container
 *   directly, e.g. http://backend:8080/api/v1. Required because relative
 *   paths like "/api/v1" have no host when fetched server-side.
 * - On the client: the relative path "/api/v1" so requests are proxied
 *   by Nginx and benefit from same-origin cookies.
 */
export const useApiBase = (): string => {
  const config = useRuntimeConfig()
  if (import.meta.server) {
    return (config.apiBaseInternal as string) || (config.public.apiBase as string)
  }
  return config.public.apiBase as string
}
