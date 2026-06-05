import { defineStore } from 'pinia'
import type { CartItem } from '~/types/api'

export const useCartStore = defineStore('cart', {
  state: () => ({ items: [] as CartItem[], loaded: false }),
  getters: {
    count: (s) => s.items.reduce((n, i) => n + i.quantity, 0),
    distinctCount: (s) => s.items.length,
    subtotal: (s) => s.items.reduce((sum, i) => sum + (i.product?.sale_price || 0) * i.quantity, 0),
    has: (s) => (productId: string) => s.items.some((i) => i.product_id === productId),
    qtyOf: (s) => (productId: string) => s.items.find((i) => i.product_id === productId)?.quantity || 0,
  },
  actions: {
    async load() {
      const auth = useAuthStore()
      if (!auth.isCustomer) {
        this.items = []
        this.loaded = true
        return
      }
      try {
        const res = await useApi()<{ data: CartItem[] }>('/customer/cart')
        this.items = res.data || []
      } catch {
        this.items = []
      }
      this.loaded = true
    },
    async set(productId: string, quantity: number) {
      await useApi()('/customer/cart', { method: 'POST', body: { product_id: productId, quantity } })
      await this.load()
    },
    async remove(productId: string) {
      await useApi()(`/customer/cart/${productId}`, { method: 'DELETE' })
      await this.load()
    },
    async clear() {
      await useApi()('/customer/cart', { method: 'DELETE' })
      this.items = []
    },
  },
})
