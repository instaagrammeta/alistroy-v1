import { defineStore } from 'pinia'
import type { SettingsMap } from '~/types/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({ loaded: false, values: {} as SettingsMap }),
  getters: {
    get: (s) => (key: string, fallback = '') => s.values[key] ?? fallback,
    logoUrl: (s) => s.values.logo_url || '',
    faviconUrl: (s) => s.values.favicon_url || '',
    siteName: (s) => (l: 'tg' | 'ru') => (l === 'ru' ? s.values.site_name_ru : s.values.site_name_tj) || 'AliStroy',
    tagline: (s) => (l: 'tg' | 'ru') => (l === 'ru' ? s.values.tagline_ru : s.values.tagline_tj) || '',
    seoDescription: (s) => (l: 'tg' | 'ru') => (l === 'ru' ? s.values.seo_description_ru : s.values.seo_description_tj) || '',
    marketplacePhone: (s) => s.values.marketplace_phone || '',
    marketplaceWhatsApp: (s) => s.values.marketplace_whatsapp || '',
    marketplaceTelegram: (s) => s.values.marketplace_telegram || '',
  },
  actions: {
    async load(force = false) {
      if (this.loaded && !force) return
      try {
        const res = await $fetch<{ data: SettingsMap }>('/settings/public', { baseURL: useApiBase() })
        this.values = res.data || {}
      } catch {
        this.values = {}
      }
      this.loaded = true
    },
    setLocal(values: SettingsMap) {
      this.values = { ...this.values, ...values }
    },
  },
})
