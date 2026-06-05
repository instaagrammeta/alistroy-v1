<script setup lang="ts">
import { useSettingsStore } from '~/stores/settings'

const settings = useSettingsStore()
const { locale } = useI18n()

useHead(() => ({
  title: settings.siteName(locale.value as 'tg' | 'ru'),
  htmlAttrs: { lang: locale.value },
  meta: [
    { name: 'description', content: settings.seoDescription(locale.value as 'tg' | 'ru') },
    { property: 'og:title', content: settings.siteName(locale.value as 'tg' | 'ru') },
    { property: 'og:description', content: settings.seoDescription(locale.value as 'tg' | 'ru') },
  ],
  link: settings.faviconUrl
    ? [{ rel: 'icon', href: settings.faviconUrl }]
    : [],
}))
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
