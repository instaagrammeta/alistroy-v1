import { defineStore } from 'pinia'
import type { SettingsMap } from '~/types/api'

interface State {
  loaded: boolean
  values: SettingsMap
}

export const useSettingsStore = defineStore('settings', {
  state: (): State => ({
    loaded: false,
    values: {},
  }),

  getters: {
    get: (s) => (key: string, fallback = '') => s.values[key] ?? fallback,
    logoUrl: (s) => s.values.logo_url || '',
    faviconUrl: (s) => s.values.favicon_url || '',
    siteName: (s) => (locale: 'tg' | 'ru') =>
      (locale === 'ru' ? s.values.site_name_ru : s.values.site_name_tj) || 'AliStroy',
    tagline: (s) => (locale: 'tg' | 'ru') =>
      (locale === 'ru' ? s.values.tagline_ru : s.values.tagline_tj) || '',
    heroTitle: (s) => (locale: 'tg' | 'ru') =>
      (locale === 'ru' ? s.values.hero_title_ru : s.values.hero_title_tj) || '',
    heroSubtitle: (s) => (locale: 'tg' | 'ru') =>
      (locale === 'ru' ? s.values.hero_subtitle_ru : s.values.hero_subtitle_tj) || '',
    seoDescription: (s) => (locale: 'tg' | 'ru') =>
      (locale === 'ru' ? s.values.seo_description_ru : s.values.seo_description_tj) || '',
    marketplacePhone: (s) => s.values.marketplace_phone || '',
    marketplaceWhatsApp: (s) => s.values.marketplace_whatsapp || '',
  },

  actions: {
    async load(force = false) {
      if (this.loaded && !force) return
      const config = useRuntimeConfig()
      try {
        const res = await $fetch<{ data: SettingsMap }>('/settings/public', {
          baseURL: useApiBase(),
        })
        this.values = res.data || {}
        this.loaded = true
      } catch {
        this.values = {}
        this.loaded = true
      }
    },
    setLocal(values: SettingsMap) {
      this.values = { ...this.values, ...values }
    },
  },
})
