import { defineStore } from 'pinia'
import type { TokenPair, User } from '~/types/api'

interface State {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  accessExpiresAt: string | null
  refreshExpiresAt: string | null
}

const STORAGE_KEY = 'alistroy.auth'

function loadFromStorage(): State {
  if (typeof localStorage === 'undefined') {
    return emptyState()
  }
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return emptyState()
    const parsed = JSON.parse(raw) as State
    return parsed
  } catch {
    return emptyState()
  }
}

function emptyState(): State {
  return {
    user: null,
    accessToken: null,
    refreshToken: null,
    accessExpiresAt: null,
    refreshExpiresAt: null,
  }
}

export const useAuthStore = defineStore('auth', {
  state: (): State => emptyState(),

  getters: {
    isAuthenticated: (s) => !!s.accessToken,
    isAdmin: (s) => s.user?.role === 'admin',
    isSeller: (s) => s.user?.role === 'seller',
  },

  actions: {
    /** Restore from localStorage on the client. */
    hydrate() {
      if (typeof window === 'undefined') return
      const restored = loadFromStorage()
      this.user = restored.user
      this.accessToken = restored.accessToken
      this.refreshToken = restored.refreshToken
      this.accessExpiresAt = restored.accessExpiresAt
      this.refreshExpiresAt = restored.refreshExpiresAt
    },

    persist() {
      if (typeof window === 'undefined') return
      const payload: State = {
        user: this.user,
        accessToken: this.accessToken,
        refreshToken: this.refreshToken,
        accessExpiresAt: this.accessExpiresAt,
        refreshExpiresAt: this.refreshExpiresAt,
      }
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(payload))
      } catch {
        /* ignore quota errors */
      }
    },

    setTokens(pair: TokenPair) {
      this.user = pair.user
      this.accessToken = pair.access_token
      this.refreshToken = pair.refresh_token
      this.accessExpiresAt = pair.access_expires_at
      this.refreshExpiresAt = pair.refresh_expires_at
      this.persist()
    },

    setUser(user: User | null) {
      this.user = user
      this.persist()
    },

    async login(email: string, password: string) {
      const config = useRuntimeConfig()
      const res = await $fetch<{ data: TokenPair }>('/auth/login', {
        baseURL: useApiBase(),
        method: 'POST',
        body: { email, password },
      })
      this.setTokens(res.data)
      return res.data
    },

    async register(payload: {
      email: string
      password: string
      name: string
      phone?: string
      role?: 'customer' | 'seller'
      seller_name?: string
      city?: string
      locale?: 'tg' | 'ru'
    }) {
      const config = useRuntimeConfig()
      const res = await $fetch<{ data: TokenPair }>('/auth/register', {
        baseURL: useApiBase(),
        method: 'POST',
        body: payload,
      })
      this.setTokens(res.data)
      return res.data
    },

    /** Returns true if refresh succeeded. */
    async tryRefresh(): Promise<boolean> {
      if (!this.refreshToken) return false
      try {
        const config = useRuntimeConfig()
        const res = await $fetch<{ data: TokenPair }>('/auth/refresh', {
          baseURL: useApiBase(),
          method: 'POST',
          body: { refresh_token: this.refreshToken },
        })
        this.setTokens(res.data)
        return true
      } catch {
        return false
      }
    },

    async fetchMe() {
      if (!this.accessToken) return null
      const api = useApi()
      try {
        const res = await api<{ data: User }>('/me')
        this.setUser(res.data)
        return res.data
      } catch {
        return null
      }
    },

    async logout(opts?: { silent?: boolean }) {
      const api = useApi()
      if (this.accessToken && !opts?.silent) {
        try {
          await api('/me/logout', { method: 'POST' })
        } catch {
          /* swallow */
        }
      }
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      this.accessExpiresAt = null
      this.refreshExpiresAt = null
      if (typeof window !== 'undefined') {
        localStorage.removeItem(STORAGE_KEY)
      }
    },
  },
})
