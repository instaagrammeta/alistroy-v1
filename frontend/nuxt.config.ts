// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  ssr: true,
  srcDir: 'app/',

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt', '@nuxtjs/i18n'],

  imports: {
    dirs: ['stores'],
  },

  css: ['~/assets/css/main.css'],

  app: {
    head: {
      htmlAttrs: { lang: 'tg' },
      title: 'AliStroy',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=5' },
        { name: 'theme-color', content: '#FF661A' },
        { name: 'mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'description', content: 'AliStroy — бозори маводи сохтмонии Тоҷикистон' },
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'AliStroy' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap',
        },
      ],
    },
  },

  runtimeConfig: {
    // server-only: SSR fetches hit the backend container directly
    apiBaseInternal: process.env.NUXT_API_BASE_INTERNAL || 'http://backend:8080/api/v1',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/v1',
      wsBase: process.env.NUXT_PUBLIC_WS_BASE || '',
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:3000',
      uploadsBase: process.env.NUXT_PUBLIC_UPLOADS_BASE || '/uploads',
      defaultLocale: process.env.NUXT_PUBLIC_DEFAULT_LOCALE || 'tg',
    },
  },

  tailwindcss: {
    cssPath: '~/assets/css/main.css',
    configPath: 'tailwind.config.ts',
  },

  i18n: {
    strategy: 'no_prefix',
    defaultLocale: 'tg',
    locales: [
      { code: 'tg', name: 'Тоҷикӣ', file: 'tg.json' },
      { code: 'ru', name: 'Русский', file: 'ru.json' },
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
      alwaysRedirect: false,
    },
    bundle: { optimizeTranslationDirective: false },
  },

  nitro: { compressPublicAssets: true },

  typescript: { strict: true, typeCheck: false },
})
