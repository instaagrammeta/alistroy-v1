<script setup lang="ts">
import type { Product } from '~/types/api'

const props = defineProps<{ product: Product; compact?: boolean }>()
const { productName, categoryName } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()
const auth = useAuthStore()
const favorites = useFavoritesStore()
const router = useRouter()
const route = useRoute()

const cover = computed(() => {
  const f = props.product.images?.[0]
  return f ? resolve(f.url) : ''
})
const isFav = computed(() => favorites.has(props.product.id))
const discounted = computed(() => props.product.discount_percent > 0)
const oldPrice = computed(() =>
  discounted.value ? props.product.sale_price / (1 - props.product.discount_percent / 100) : 0
)

const toggleFav = async () => {
  if (!auth.isAuthenticated) {
    router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }
  await favorites.toggle(props.product.id)
}
</script>

<template>
  <div class="group card overflow-hidden hover:shadow-card transition-shadow flex flex-col">
    <NuxtLink :to="`/products/${product.slug}`" class="block aspect-square bg-gray-50 relative overflow-hidden">
      <img v-if="cover" :src="cover" :alt="productName(product)" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" loading="lazy" />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-300">
        <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.1-3.1a2 2 0 0 0-2.81.01L6 21"/></svg>
      </div>
      <span v-if="product.is_featured" class="absolute top-2 left-2 badge bg-brand-500 text-white">★</span>
      <span v-if="discounted" class="absolute top-2 left-2 mt-7 badge bg-red-500 text-white">-{{ Math.round(product.discount_percent) }}%</span>
      <button class="absolute top-2 right-2 w-8 h-8 rounded-full bg-white/90 flex items-center justify-center shadow" :class="isFav ? 'text-red-500' : 'text-gray-400'" @click.prevent="toggleFav">
        <svg width="16" height="16" viewBox="0 0 24 24" :fill="isFav ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
      </button>
      <span v-if="!product.is_available" class="absolute bottom-2 right-2 badge bg-gray-800/80 text-white text-[10px]">{{ $t('catalog.out_of_stock') }}</span>
    </NuxtLink>
    <div class="p-3 sm:p-4 flex flex-col flex-1">
      <div class="text-xs text-gray-400 truncate">{{ categoryName(product.category) }}</div>
      <NuxtLink :to="`/products/${product.slug}`" class="mt-1 font-semibold text-ink-900 line-clamp-2 group-hover:text-brand-600 min-h-[2.8rem]">{{ productName(product) }}</NuxtLink>
      <div class="mt-2 flex items-end justify-between">
        <div>
          <div class="text-lg font-bold text-brand-600">{{ formatPrice(product.sale_price, product.currency) }}</div>
          <div v-if="discounted" class="text-xs text-gray-400 line-through">{{ formatPrice(oldPrice, product.currency) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
