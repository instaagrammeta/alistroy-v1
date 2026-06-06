<script setup lang="ts">
import type { Banner, Category, Product } from '~/types/api'

const { t, locale } = useI18n()
const settings = useSettingsStore()

const fetchApi = (url: string, query?: Record<string, any>) => $fetch(url, { baseURL: useApiBase(), query })

const { data: bannersRes } = await useAsyncData('home-banners', () => fetchApi('/banners'))
const { data: catsRes } = await useAsyncData('home-cats', () => fetchApi('/categories'))
const { data: featuredRes } = await useAsyncData('home-featured', () => fetchApi('/products', { featured: true, page_size: 12 }))
const { data: latestRes } = await useAsyncData('home-latest', () => fetchApi('/products', { page_size: 12, sort: 'newest' }))

const banners = computed<Record<string, Banner[]>>(() => (bannersRes.value as any)?.data || {})
const heroBanners = computed(() => banners.value.hero || [])
const featureBanners = computed(() => banners.value.features || [])
const categoryTiles = computed(() => banners.value.category_row || [])
const midLarge = computed(() => banners.value.mid_large || [])
const midSmall = computed(() => banners.value.mid_small || [])

const categories = computed<Category[]>(() => (catsRes.value as any)?.data || [])
const featured = computed<Product[]>(() => (featuredRes.value as any)?.data || [])
const latest = computed<Product[]>(() => (latestRes.value as any)?.data || [])

useHead(() => ({
  title: `${settings.siteName(locale.value as 'tg' | 'ru')} — ${settings.tagline(locale.value as 'tg' | 'ru')}`,
}))
</script>

<template>
  <div class="pb-12">
    <!-- 1. Hero slider + categories sidebar (PC) -->
    <section class="container-page pt-4">
      <div class="grid lg:grid-cols-[260px_1fr] gap-4">
        <div class="hidden lg:block">
          <CategorySidebar :categories="categories" />
        </div>
        <div>
          <BannerSlider v-if="heroBanners.length" :banners="heroBanners" :interval="5000" ratio="aspect-[21/9]" />
          <div v-else class="aspect-[21/9] rounded-2xl bg-gradient-to-br from-brand-400 to-brand-600 flex items-center justify-center text-center text-white p-8">
            <div>
              <h1 class="text-2xl md:text-4xl font-extrabold">{{ settings.get(locale === 'ru' ? 'hero_title_ru' : 'hero_title_tj') || t('site.tagline') }}</h1>
              <NuxtLink to="/products" class="btn bg-white text-brand-600 hover:bg-gray-100 mt-5 px-6 py-3">{{ t('nav.catalog') }}</NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 2. Sub-banner: feature row with icons (admin-managed in /admin/banners?position=features) -->
    <section v-if="featureBanners.length" class="container-page mt-5">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <FeatureBanner v-for="b in featureBanners" :key="b.id" :banner="b" />
      </div>
    </section>

    <!-- 3. Category tiles (1000x1000 admin-managed) -->
    <section v-if="categoryTiles.length" class="container-page mt-6">
      <div class="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-6 gap-3 sm:gap-4">
        <CategoryTile v-for="b in categoryTiles" :key="b.id" :banner="b" />
      </div>
    </section>

    <!-- 4. Mid promo: 2 small (left) + 1 large (right) -->
    <section v-if="midLarge.length || midSmall.length" class="container-page mt-8">
      <div class="grid lg:grid-cols-2 gap-4">
        <div class="grid grid-rows-2 gap-4">
          <BannerTile v-for="b in midSmall.slice(0, 2)" :key="b.id" :banner="b" ratio="aspect-[16/6]" />
        </div>
        <BannerTile :banner="midLarge[0]" ratio="aspect-[16/9] lg:h-full" />
      </div>
    </section>

    <!-- 5. Featured products -->
    <section v-if="featured.length" class="container-page mt-10">
      <div class="flex items-end justify-between mb-4">
        <h2 class="text-xl md:text-2xl font-bold text-ink-900">{{ t('home.featured_products') }}</h2>
        <NuxtLink to="/products?featured=true" class="text-sm font-medium text-brand-600 hover:text-brand-700">{{ t('common.all') }} →</NuxtLink>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
        <ProductCard v-for="p in featured" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- 6. Latest products -->
    <section class="container-page mt-10">
      <div class="flex items-end justify-between mb-4">
        <h2 class="text-xl md:text-2xl font-bold text-ink-900">{{ t('home.latest_products') }}</h2>
        <NuxtLink to="/products" class="text-sm font-medium text-brand-600 hover:text-brand-700">{{ t('common.all') }} →</NuxtLink>
      </div>
      <EmptyState v-if="!latest.length" />
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
        <ProductCard v-for="p in latest" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- 7. Contact -->
    <section class="bg-ink-900 text-white mt-12">
      <div class="container-page py-12 text-center">
        <h2 class="text-2xl md:text-3xl font-bold">{{ t('home.contact_title') }}</h2>
        <p class="mt-2 text-gray-300 max-w-xl mx-auto">{{ t('home.contact_text') }}</p>
        <div class="mt-6 flex flex-wrap justify-center gap-3">
          <a v-if="settings.marketplacePhone" :href="`tel:${settings.marketplacePhone.replace(/[^0-9+]/g, '')}`" class="btn-primary">{{ settings.marketplacePhone }}</a>
          <a v-if="settings.marketplaceWhatsApp" :href="`https://wa.me/${settings.marketplaceWhatsApp.replace(/[^0-9]/g, '')}`" target="_blank" class="btn-success">WhatsApp</a>
        </div>
      </div>
    </section>
  </div>
</template>
