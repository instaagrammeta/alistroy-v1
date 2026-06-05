<script setup lang="ts">
import type { Product } from '~/types/api'

const props = defineProps<{ product: Product; compact?: boolean }>()
const { productTitle, categoryTitle } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()

const cover = computed(() => {
  const first = props.product.images?.[0]
  return first ? resolve(first.url) : ''
})
</script>

<template>
  <NuxtLink :to="`/products/${product.slug}`" class="group block card overflow-hidden hover:shadow-card transition-shadow">
    <div class="aspect-square bg-gray-50 relative overflow-hidden">
      <img
        v-if="cover"
        :src="cover"
        :alt="productTitle(product)"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-300">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="3" y="3" width="18" height="18" rx="2"/>
          <circle cx="9" cy="9" r="2"/>
          <path d="m21 15-3.1-3.1a2 2 0 0 0-2.81.01L6 21"/>
        </svg>
      </div>
      <span v-if="product.is_featured" class="absolute top-2 left-2 badge bg-brand-500 text-white">★</span>
      <span v-if="!product.is_available" class="absolute top-2 right-2 badge bg-red-500 text-white text-[10px]">{{ $t('catalog.out_of_stock') }}</span>
    </div>
    <div class="p-3 sm:p-4">
      <div class="text-xs text-gray-400 truncate">
        {{ categoryTitle(product.category) }}
      </div>
      <h3 class="mt-1 font-semibold text-ink-900 line-clamp-2 group-hover:text-brand-600 min-h-[3rem]">
        {{ productTitle(product) }}
      </h3>
      <div class="mt-2 flex items-end justify-between">
        <div class="text-lg font-bold text-brand-600">
          {{ formatPrice(product.price, product.currency) }}
        </div>
        <div v-if="!compact && product.seller" class="text-xs text-gray-500 truncate max-w-[50%]">
          {{ product.seller.name }}
        </div>
      </div>
    </div>
  </NuxtLink>
</template>
