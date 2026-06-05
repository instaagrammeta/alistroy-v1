<script setup lang="ts">
import type { Product } from '~/types/api'

const props = defineProps<{ product: Product; size?: 'sm' | 'md' | 'lg' }>()
const { t } = useI18n()
const { telLink, waLink, tgLink } = useFormatters()

const track = async (event: 'phone_click' | 'whatsapp_click' | 'telegram_click') => {
  try {
    await $fetch(`/products/id/${props.product.id}/track`, {
      baseURL: useApiBase(), method: 'POST', body: { event },
    })
  } catch {}
}

const sizeClass = computed(() => (props.size === 'lg' ? 'px-6 py-3 text-base' : props.size === 'sm' ? 'px-3 py-2 text-sm' : 'px-4 py-2.5 text-sm'))
</script>

<template>
  <div class="flex flex-wrap gap-3">
    <a v-if="product.contact_phone" :href="telLink(product.contact_phone)" class="btn btn-primary" :class="sizeClass" @click="track('phone_click')">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13 1.05.37 2.07.72 3.05a2 2 0 0 1-.45 2.11L8.09 10.91a16 16 0 0 0 6 6l2-1.27a2 2 0 0 1 2.11-.45c.98.35 2 .59 3.05.72A2 2 0 0 1 22 16.92z"/></svg>
      {{ t('product.call') }}
    </a>
    <a v-if="product.contact_whatsapp" :href="waLink(product.contact_whatsapp)" target="_blank" rel="noopener" class="btn btn-success" :class="sizeClass" @click="track('whatsapp_click')">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M19.05 4.91A9.82 9.82 0 0 0 12.04 2C6.58 2 2.13 6.45 2.13 11.91c0 1.75.46 3.45 1.32 4.95L2.05 22l5.27-1.38a9.9 9.9 0 0 0 4.72 1.2h.01c5.46 0 9.91-4.45 9.91-9.91a9.83 9.83 0 0 0-2.91-7zm-7.01 15.24a8.23 8.23 0 0 1-4.19-1.15l-.3-.18-3.13.82.83-3.05-.2-.31a8.22 8.22 0 0 1-1.26-4.37c0-4.55 3.7-8.25 8.25-8.25 2.2 0 4.27.86 5.83 2.42a8.18 8.18 0 0 1 2.42 5.83c0 4.55-3.7 8.24-8.24 8.24z"/></svg>
      {{ t('product.whatsapp') }}
    </a>
    <a v-if="product.contact_telegram" :href="tgLink(product.contact_telegram)" target="_blank" rel="noopener" class="btn" :class="['bg-sky-500 hover:bg-sky-600 text-white', sizeClass]" @click="track('telegram_click')">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M22 2 2 10l6 2 2 6 3-4 5 4 4-16zM9.5 13.5 18 6l-7 8.5z"/></svg>
      {{ t('product.telegram') }}
    </a>
  </div>
</template>
