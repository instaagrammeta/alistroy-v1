// vue-i18n runtime config (used by @nuxtjs/i18n v9)
export default defineI18nConfig(() => ({
  legacy: false,
  fallbackLocale: 'tg',
  numberFormats: {
    tg: {
      currency: {
        style: 'currency',
        currency: 'TJS',
      },
    },
    ru: {
      currency: {
        style: 'currency',
        currency: 'TJS',
      },
    },
  },
}))
