// Hydrates the auth store from localStorage on app load (client only).
export default defineNuxtPlugin(async () => {
  const auth = useAuthStore()
  auth.hydrate()
  if (auth.isAuthenticated) {
    await auth.fetchMe()
  }
})
