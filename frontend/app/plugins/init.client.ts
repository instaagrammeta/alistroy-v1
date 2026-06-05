// Hydrate auth + cart + favorites from storage/server on client load.
export default defineNuxtPlugin(async () => {
  const auth = useAuthStore()
  auth.hydrate()
  if (auth.isAuthenticated) {
    await auth.fetchMe()
    const cart = useCartStore()
    const favorites = useFavoritesStore()
    await Promise.all([cart.load(), favorites.load()])
  }
})
