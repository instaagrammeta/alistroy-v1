<script setup lang="ts">
import type { Banner } from '~/types/api'

defineProps<{ banner: Banner | undefined; ratio?: string }>()
const { resolve } = useImageUrl()
</script>

<template>
  <component
    v-if="banner"
    :is="banner.link_url ? 'a' : 'div'"
    :href="banner.link_url || undefined"
    class="block overflow-hidden rounded-2xl group"
    :class="ratio || 'aspect-[16/9]'"
  >
    <picture>
      <source v-if="banner.mobile_url" :srcset="resolve(banner.mobile_url)" media="(max-width: 640px)" />
      <img :src="resolve(banner.desktop_url || banner.mobile_url)" :alt="banner.title_ru || banner.title_tj" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
    </picture>
  </component>
</template>
