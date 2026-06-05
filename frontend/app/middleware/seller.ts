export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()
  if (typeof window !== 'undefined') auth.hydrate()
  if (!auth.isAuthenticated) return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  if (auth.user?.role !== 'seller') return navigateTo('/')
})
