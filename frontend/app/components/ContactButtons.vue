<script setup lang="ts">
import type { Product } from '~/types/api'

const props = defineProps<{
  product: Product
  size?: 'sm' | 'md' | 'lg'
}>()
const { t } = useI18n()
const { phoneLink, whatsappLink } = useFormatters()

const trackEvent = async (event: 'phone_click' | 'whatsapp_click') => {
  try {
    const config = useRuntimeConfig()
    await $fetch(`/products/id/${props.product.id}/track`, {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: { event },
    })
  } catch {
    /* tracking is best-effort */
  }
}

const sizeClass = computed(() => {
  switch (props.size || 'md') {
    case 'sm':
      return 'px-3 py-2 text-sm'
    case 'lg':
      return 'px-6 py-3 text-base'
    default:
      return 'px-4 py-2.5 text-sm'
  }
})
</script>

<template>
  <div class="flex flex-wrap gap-3">
    <a
      v-if="product.phone_number"
      :href="phoneLink(product.phone_number)"
      class="btn btn-primary"
      :class="sizeClass"
      @click="trackEvent('phone_click')"
    >
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
        <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13 1.05.37 2.07.72 3.05a2 2 0 0 1-.45 2.11L8.09 10.91a16 16 0 0 0 6 6l2-1.27a2 2 0 0 1 2.11-.45c.98.35 2 .59 3.05.72A2 2 0 0 1 22 16.92z"/>
      </svg>
      {{ t('product.call') }}
    </a>

    <a
      v-if="product.whatsapp_number"
      :href="whatsappLink(product.whatsapp_number)"
      target="_blank"
      rel="noopener"
      class="btn btn-success"
      :class="sizeClass"
      @click="trackEvent('whatsapp_click')"
    >
      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" class="mr-2">
        <path d="M19.05 4.91A9.82 9.82 0 0 0 12.04 2c-5.46 0-9.91 4.45-9.91 9.91 0 1.75.46 3.45 1.32 4.95L2.05 22l5.27-1.38a9.9 9.9 0 0 0 4.72 1.2h.01c5.46 0 9.91-4.45 9.91-9.91a9.83 9.83 0 0 0-2.91-7zm-7.01 15.24h-.01a8.23 8.23 0 0 1-4.19-1.15l-.3-.18-3.13.82.83-3.05-.2-.31a8.22 8.22 0 0 1-1.26-4.37c0-4.55 3.7-8.25 8.25-8.25 2.2 0 4.27.86 5.83 2.42a8.18 8.18 0 0 1 2.42 5.83c0 4.55-3.7 8.24-8.24 8.24zm4.52-6.16c-.25-.13-1.47-.72-1.69-.81-.23-.08-.39-.13-.56.13-.17.25-.64.81-.79.97-.15.17-.29.18-.54.06-.25-.12-1.05-.39-2-1.23-.74-.66-1.23-1.47-1.38-1.72-.14-.25-.02-.39.11-.51.11-.11.25-.29.37-.43.12-.14.16-.25.25-.41.08-.17.04-.31-.02-.43-.06-.13-.56-1.34-.77-1.83-.2-.49-.41-.42-.56-.43h-.48c-.17 0-.43.06-.66.31-.23.25-.86.84-.86 2.06s.88 2.39 1 2.55c.13.17 1.74 2.66 4.21 3.73.59.25 1.05.4 1.41.51.59.19 1.13.16 1.55.1.47-.07 1.47-.6 1.67-1.18.21-.58.21-1.07.14-1.18-.06-.11-.22-.18-.47-.31z"/>
      </svg>
      {{ t('product.whatsapp') }}
    </a>
  </div>
</template>
