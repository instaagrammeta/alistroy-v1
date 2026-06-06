<script setup lang="ts">
import type { Banner } from '~/types/api'

defineProps<{ banner: Banner }>()
const { resolve } = useImageUrl()
const { locale } = useI18n()
const title = (b: Banner) => (locale.value === 'ru' ? b.title_ru : b.title_tj) || b.title_tj
const desc = (b: Banner) => (locale.value === 'ru' ? b.description_ru : b.description_tj) || b.description_tj
</script>

<template>
  <component
    :is="banner.link_url ? 'a' : 'div'"
    :href="banner.link_url || undefined"
    class="card p-4 sm:p-5 flex items-center gap-3 hover:shadow-card transition-shadow"
  >
    <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shrink-0 overflow-hidden">
      <img v-if="banner.desktop_url || banner.mobile_url" :src="resolve(banner.desktop_url || banner.mobile_url)" alt="" class="w-9 h-9 sm:w-10 sm:h-10 object-contain" />
      <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2 L15 8 L22 9 L17 14 L18 21 L12 18 L6 21 L7 14 L2 9 L9 8 Z"/></svg>
    </div>
    <div class="min-w-0">
      <div class="font-semibold text-ink-900 text-sm sm:text-base truncate">{{ title(banner) }}</div>
      <div v-if="desc(banner)" class="text-xs sm:text-sm text-gray-500 line-clamp-2">{{ desc(banner) }}</div>
    </div>
  </component>
</template>
