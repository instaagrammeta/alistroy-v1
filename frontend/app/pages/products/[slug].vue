<script setup lang="ts">
import type { Product, Review } from '~/types/api'

const route = useRoute()
const router = useRouter()
const slug = route.params.slug as string

const auth = useAuthStore()
const cart = useCartStore()
const favorites = useFavoritesStore()
const { productName, productDesc, categoryName } = useLocaleField()
const { formatPrice, formatDate } = useFormatters()
const { resolve } = useImageUrl()
const { t } = useI18n()

const { data: bundle, error } = await useAsyncData(`product-${slug}`, () =>
  $fetch<{ data: { product: Product; avg_rating: number; review_count: number } }>(`/products/${slug}`, { baseURL: useApiBase() })
)
if (error.value) throw createError({ statusCode: 404, statusMessage: 'Product not found' })

const product = computed<Product>(() => (bundle.value as any)?.data?.product)
const avgRating = computed<number>(() => (bundle.value as any)?.data?.avg_rating ?? 0)
const reviewCount = computed<number>(() => (bundle.value as any)?.data?.review_count ?? 0)

const { data: relatedRes } = await useAsyncData(`product-${slug}-related`, () =>
  $fetch(`/catalog/product/${product.value?.id}/related`, { baseURL: useApiBase(), query: { limit: 10 } })
)
const related = computed<Product[]>(() => (relatedRes.value as any)?.data || [])

const { data: reviewsRes, refresh: refreshReviews } = await useAsyncData(`product-${slug}-reviews`, () =>
  $fetch(`/catalog/product/${product.value?.id}/reviews`, { baseURL: useApiBase(), query: { page_size: 20 } })
)
const reviews = computed<Review[]>(() => (reviewsRes.value as any)?.data || [])

const activeImage = ref(0)

useHead(() => ({
  title: `${productName(product.value)} — AliStroy`,
  meta: [
    { name: 'description', content: productDesc(product.value).slice(0, 200) },
    { property: 'og:title', content: productName(product.value) },
    { property: 'og:image', content: product.value?.images?.[0]?.url ? resolve(product.value.images[0].url) : '' },
  ],
}))

onMounted(async () => {
  if (!product.value) return
  try {
    await $fetch(`/catalog/product/${product.value.id}/track`, { baseURL: useApiBase(), method: 'POST', body: { event: 'view' } })
  } catch {}
})

const isFav = computed(() => product.value && favorites.has(product.value.id))
const inCart = computed(() => product.value && cart.has(product.value.id))

const toggleFav = async () => {
  if (!auth.isAuthenticated) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  await favorites.toggle(product.value.id)
}
const addToCart = async () => {
  if (!auth.isAuthenticated) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  if (!auth.isCustomer) return
  await cart.set(product.value.id, (cart.qtyOf(product.value.id) || 0) + 1)
}

const newReview = reactive({ rating: 5, comment: '' })
const reviewSent = ref(false)
const submitReview = async () => {
  if (!auth.isAuthenticated) return router.push({ path: '/login', query: { redirect: route.fullPath } })
  try {
    await useApi()(`/catalog/product/${product.value.id}/reviews`, { method: 'POST', body: { rating: newReview.rating, comment: newReview.comment } })
    reviewSent.value = true
    newReview.comment = ''; newReview.rating = 5
    await refreshReviews()
  } catch {}
}
</script>

<template>
  <div v-if="product" class="container-page py-6">
    <div class="grid lg:grid-cols-2 gap-8">
      <!-- Gallery -->
      <div>
        <div class="aspect-square card overflow-hidden">
          <img v-if="product.images?.length" :src="resolve(product.images[activeImage]?.url)" :alt="productName(product)" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full bg-gray-100 flex items-center justify-center text-gray-300">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.1-3.1a2 2 0 0 0-2.81.01L6 21"/></svg>
          </div>
        </div>
        <div v-if="product.images && product.images.length > 1" class="mt-3 grid grid-cols-5 gap-2">
          <button v-for="(img, idx) in product.images" :key="img.id" class="aspect-square rounded-lg overflow-hidden border-2" :class="idx === activeImage ? 'border-brand-500' : 'border-transparent'" @click="activeImage = idx">
            <img :src="resolve(img.url)" alt="" class="w-full h-full object-cover" />
          </button>
        </div>
      </div>

      <!-- Info -->
      <div>
        <div class="text-sm text-gray-500 flex flex-wrap items-center gap-2">
          <NuxtLink :to="`/categories/${product.category?.slug}`" class="hover:text-brand-600">{{ categoryName(product.category) }}</NuxtLink>
          <span v-if="product.sku">• {{ t('product.sku') }}: {{ product.sku }}</span>
        </div>
        <h1 class="mt-2 text-2xl md:text-3xl font-bold text-ink-900">{{ productName(product) }}</h1>

        <div class="mt-4 flex items-center gap-3">
          <div class="text-3xl font-extrabold text-brand-600">{{ formatPrice(product.sale_price, product.currency) }}</div>
          <span v-if="product.unit" class="text-sm text-gray-500">/ {{ product.unit }}</span>
        </div>

        <div class="mt-3 flex items-center gap-3 flex-wrap">
          <span v-if="product.is_available" class="badge bg-emerald-50 text-emerald-700 border border-emerald-100">{{ t('catalog.in_stock') }}<span v-if="product.stock_quantity > 0" class="ml-1">({{ product.stock_quantity }})</span></span>
          <span v-else class="badge bg-red-50 text-red-700 border border-red-100">{{ t('catalog.out_of_stock') }}</span>
          <span v-if="reviewCount > 0" class="text-sm text-gray-500">★ {{ avgRating.toFixed(1) }} ({{ reviewCount }})</span>
        </div>

        <div class="mt-6">
          <ContactButtons :product="product" size="lg" />
          <div class="mt-3 flex gap-2 flex-wrap">
            <button v-if="auth.isCustomer || !auth.isAuthenticated" class="btn-outline" @click="addToCart">{{ inCart ? t('product.in_cart') : t('product.add_to_cart') }}</button>
            <button class="btn-outline" @click="toggleFav">{{ isFav ? t('product.favorites_remove') : t('product.favorites_add') }}</button>
          </div>
        </div>

        <NuxtLink v-if="product.seller" :to="`/sellers/${product.seller.slug}`" class="mt-6 card p-4 flex items-center gap-3 hover:shadow-card transition-shadow">
          <div class="w-12 h-12 rounded-xl bg-gray-100 overflow-hidden flex items-center justify-center">
            <img v-if="product.seller.logo_url" :src="resolve(product.seller.logo_url)" :alt="product.seller.full_name" class="w-full h-full object-cover" />
            <span v-else class="text-brand-600 font-bold">{{ (product.seller.market_name || product.seller.full_name).charAt(0).toUpperCase() }}</span>
          </div>
          <div class="min-w-0">
            <div class="font-semibold text-ink-900 truncate">{{ product.seller.market_name || product.seller.company_name || product.seller.full_name }}</div>
            <div class="text-xs text-gray-400 truncate">{{ product.seller.city || ' ' }}</div>
          </div>
          <span class="ml-auto text-brand-600">→</span>
        </NuxtLink>

        <section v-if="productDesc(product)" class="mt-6">
          <h2 class="text-lg font-semibold text-ink-900 mb-2">{{ t('product.description') }}</h2>
          <p class="text-gray-700 whitespace-pre-line leading-relaxed">{{ productDesc(product) }}</p>
        </section>
      </div>
    </div>

    <!-- Reviews -->
    <section class="mt-10">
      <h2 class="text-xl font-bold text-ink-900 mb-4">{{ t('product.reviews') }}</h2>
      <div v-if="auth.isAuthenticated" class="card p-4 mb-6">
        <h3 class="font-semibold mb-3">{{ t('product.leave_review') }}</h3>
        <form class="space-y-3" @submit.prevent="submitReview">
          <div>
            <label class="label">{{ t('product.rating') }}</label>
            <div class="flex gap-1">
              <button v-for="n in 5" :key="n" type="button" class="text-2xl" :class="n <= newReview.rating ? 'text-amber-400' : 'text-gray-300'" @click="newReview.rating = n">★</button>
            </div>
          </div>
          <div>
            <label class="label">{{ t('product.comment') }}</label>
            <textarea v-model="newReview.comment" class="textarea" rows="3"></textarea>
          </div>
          <div class="flex items-center gap-3">
            <button type="submit" class="btn-primary">{{ t('product.submit_review') }}</button>
            <span v-if="reviewSent" class="text-sm text-emerald-600">{{ t('product.review_pending') }}</span>
          </div>
        </form>
      </div>
      <EmptyState v-if="!reviews.length" :title="t('product.no_reviews')" />
      <div v-else class="space-y-3">
        <div v-for="r in reviews" :key="r.id" class="card p-4">
          <div class="flex items-center justify-between">
            <div class="font-medium">{{ r.user?.name || '—' }}</div>
            <div class="text-amber-400">{{ '★'.repeat(r.rating) }}<span class="text-gray-200">{{ '★'.repeat(5 - r.rating) }}</span></div>
          </div>
          <p v-if="r.comment" class="mt-2 text-gray-700">{{ r.comment }}</p>
          <div class="mt-2 text-xs text-gray-400">{{ formatDate(r.created_at) }}</div>
        </div>
      </div>
    </section>

    <!-- Related -->
    <section v-if="related.length" class="mt-12">
      <h2 class="text-xl font-bold text-ink-900 mb-4">{{ t('product.related') }}</h2>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
        <ProductCard v-for="p in related" :key="p.id" :product="p" />
      </div>
    </section>
  </div>
</template>
