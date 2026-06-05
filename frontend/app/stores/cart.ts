// Note: AliStroy is not an online-payment marketplace. We use this lightweight
// "favorites" store on the client purely to keep heart icons reactive in
// product cards. The source of truth lives on the server (Favorite resource).

import { defineStore } from 'pinia'

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({
    ids: new Set<string>(),
    loaded: false,
  }),

  getters: {
    has: (s) => (id: string) => s.ids.has(id),
  },

  actions: {
    async loadFromServer() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) {
        this.ids = new Set()
        this.loaded = true
        return
      }
      const api = useApi()
      try {
        const res = await api<{ data: Array<{ product_id: string }> }>('/favorites?page_size=100')
        this.ids = new Set((res.data || []).map((f) => f.product_id))
      } catch {
        this.ids = new Set()
      }
      this.loaded = true
    },

    async toggle(productId: string) {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) {
        // Defer to caller — they typically redirect to /login.
        return false
      }
      const api = useApi()
      if (this.ids.has(productId)) {
        await api(`/favorites/${productId}`, { method: 'DELETE' })
        this.ids.delete(productId)
      } else {
        await api(`/favorites/${productId}`, { method: 'POST' })
        this.ids.add(productId)
      }
      return true
    },
  },
})
