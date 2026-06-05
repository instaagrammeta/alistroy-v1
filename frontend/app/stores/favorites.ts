import { defineStore } from 'pinia'

export const useFavoritesStore = defineStore('favorites', {
  state: () => ({ ids: [] as string[], loaded: false }),
  getters: {
    has: (s) => (id: string) => s.ids.includes(id),
  },
  actions: {
    async load() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) {
        this.ids = []
        this.loaded = true
        return
      }
      try {
        const res = await useApi()<{ data: Array<{ product_id: string }> }>('/favorites?page_size=200')
        this.ids = (res.data || []).map((f) => f.product_id)
      } catch {
        this.ids = []
      }
      this.loaded = true
    },
    async toggle(productId: string): Promise<boolean> {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return false
      if (this.ids.includes(productId)) {
        await useApi()(`/favorites/${productId}`, { method: 'DELETE' })
        this.ids = this.ids.filter((x) => x !== productId)
      } else {
        await useApi()(`/favorites/${productId}`, { method: 'POST' })
        this.ids.push(productId)
      }
      return true
    },
  },
})
