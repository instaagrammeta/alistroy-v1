<script setup lang="ts">
import type { Banner } from '~/types/api'

defineProps<{ banner: Banner }>()
const { resolve } = useImageUrl()
const { locale } = useI18n()
const title = (b: Banner) => (locale.value === 'ru' ? b.title_ru : b.title_tj) || b.title_tj
</script>

<template>
  <component
    :is="banner.link_url ? 'a' : 'div'"
    :href="banner.link_url || undefined"
    class="block group"
  >
    <div class="aspect-square rounded-2xl overflow-hidden bg-gray-100 ring-1 ring-gray-100 group-hover:ring-brand-300 transition-all">
      <picture>
        <source v-if="banner.mobile_url" :srcset="resolve(banner.mobile_url)" media="(max-width: 640px)" />
        <img
          v-if="banner.desktop_url || banner.mobile_url"
          :src="resolve(banner.desktop_url || banner.mobile_url)"
          :alt="title(banner)"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
        />
        <div v-else class="w-full h-full bg-gradient-to-br from-brand-100 to-brand-300"></div>
      </picture>
    </div>
    <div v-if="title(banner)" class="mt-2 text-center text-sm font-medium text-ink-900 truncate">{{ title(banner) }}</div>
  </component>
</template>
