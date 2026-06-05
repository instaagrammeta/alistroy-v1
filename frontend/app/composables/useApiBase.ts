/**
 * Returns the right API base for the current context:
 * - SSR (server): absolute URL to the backend container (apiBaseInternal)
 * - Client (browser): relative "/api/v1" proxied by Nginx
 */
export const useApiBase = (): string => {
  const config = useRuntimeConfig()
  if (import.meta.server) {
    return (config.apiBaseInternal as string) || (config.public.apiBase as string)
  }
  return config.public.apiBase as string
}
