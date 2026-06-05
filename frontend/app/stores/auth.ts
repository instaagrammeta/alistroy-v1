import { defineStore } from 'pinia'
import type { TokenPair, User } from '~/types/api'

interface State {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
}

const KEY = 'alistroy.auth'

function empty(): State {
  return { user: null, accessToken: null, refreshToken: null }
}

function load(): State {
  if (typeof localStorage === 'undefined') return empty()
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? (JSON.parse(raw) as State) : empty()
  } catch {
    return empty()
  }
}

export const useAuthStore = defineStore('auth', {
  state: (): State => empty(),

  getters: {
    isAuthenticated: (s) => !!s.accessToken,
    isAdmin: (s) => s.user?.role === 'admin',
    isSeller: (s) => s.user?.role === 'seller',
    isDriver: (s) => s.user?.role === 'driver',
    isCustomer: (s) => s.user?.role === 'customer',
  },

  actions: {
    hydrate() {
      if (typeof window === 'undefined') return
      const s = load()
      this.user = s.user
      this.accessToken = s.accessToken
      this.refreshToken = s.refreshToken
    },
    persist() {
      if (typeof window === 'undefined') return
      try {
        localStorage.setItem(KEY, JSON.stringify({ user: this.user, accessToken: this.accessToken, refreshToken: this.refreshToken }))
      } catch {}
    },
    setTokens(p: TokenPair) {
      this.user = p.user
      this.accessToken = p.access_token
      this.refreshToken = p.refresh_token
      this.persist()
    },
    setUser(u: User | null) {
      this.user = u
      this.persist()
    },
    async login(identifier: string, password: string) {
      const res = await $fetch<{ data: TokenPair }>('/auth/login', {
        baseURL: useApiBase(), method: 'POST', body: { identifier, password },
      })
      this.setTokens(res.data)
      return res.data
    },
    async register(payload: { name: string; phone: string; email?: string; password: string; address?: string; city?: string; locale?: string }) {
      const res = await $fetch<{ data: TokenPair }>('/auth/register', {
        baseURL: useApiBase(), method: 'POST', body: payload,
      })
      this.setTokens(res.data)
      return res.data
    },
    async tryRefresh(): Promise<boolean> {
      if (!this.refreshToken) return false
      try {
        const res = await $fetch<{ data: TokenPair }>('/auth/refresh', {
          baseURL: useApiBase(), method: 'POST', body: { refresh_token: this.refreshToken },
        })
        this.setTokens(res.data)
        return true
      } catch {
        return false
      }
    },
    async fetchMe() {
      if (!this.accessToken) return null
      try {
        const res = await useApi()<{ data: User }>('/me')
        this.setUser(res.data)
        return res.data
      } catch {
        return null
      }
    },
    async logout(opts?: { silent?: boolean }) {
      if (this.accessToken && !opts?.silent) {
        try { await useApi()('/me/logout', { method: 'POST' }) } catch {}
      }
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      if (typeof window !== 'undefined') localStorage.removeItem(KEY)
    },
  },
})
