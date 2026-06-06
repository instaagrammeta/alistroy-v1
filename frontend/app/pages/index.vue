<script setup lang="ts">
import type { Banner, Category, Product, Seller } from '~/types/api'

const { t, locale } = useI18n()
const settings = useSettingsStore()

const fetchApi = (url: string, query?: Record<string, any>) => $fetch(url, { baseURL: useApiBase(), query })

const { data: bannersRes } = await useAsyncData('home-banners', () => fetchApi('/banners'))
const { data: catsRes } = await useAsyncData('home-cats', () => fetchApi('/categories'))
const { data: popularRes } = await useAsyncData('home-popular', () => fetchApi('/catalog/popular-categories', { limit: 12 }))
const { data: featuredRes } = await useAsyncData('home-featured', () => fetchApi('/products', { featured: true, page_size: 12 }))
const { data: latestRes } = await useAsyncData('home-latest', () => fetchApi('/products', { page_size: 12, sort: 'newest' }))
const { data: topRes } = await useAsyncData('home-top', () => fetchApi('/catalog/top-sellers', { limit: 6 }))

const banners = computed<Record<string, Banner[]>>(() => (bannersRes.value as any)?.data || {})
const heroBanners = computed(() => banners.value.hero || [])
const categoryRow = computed(() => banners.value.category_row || [])
const midLarge = computed(() => (banners.value.mid_large || [])[0])
const midSmall = computed(() => banners.value.mid_small || [])

const categories = computed<Category[]>(() => (catsRes.value as any)?.data || [])
const popularCats = computed<Category[]>(() => (popularRes.value as any)?.data || [])
const featured = computed<Product[]>(() => (featuredRes.value as any)?.data || [])
const latest = computed<Product[]>(() => (latestRes.value as any)?.data || [])
const topSellers = computed<Seller[]>(() => (topRes.value as any)?.data || [])

const { resolve } = useImageUrl()

useHead(() => ({
  title: `${settings.siteName(locale.value as 'tg' | 'ru')} — ${settings.tagline(locale.value as 'tg' | 'ru')}`,
}))
</script>

<template>
  <div class="pb-10">
    <!-- HERO: categories (left, PC) + animated slider -->
    <section class="container-page pt-5">
      <div class="grid lg:grid-cols-[260px_1fr] gap-5">
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

    <!-- Category image tiles (admin-managed, clickable) -->
    <section v-if="categoryRow.length" class="container-page mt-5">
      <div class="flex gap-4 overflow-x-auto no-scrollbar pb-2">
        <component
          :is="b.link_url ? 'a' : 'div'"
          v-for="b in categoryRow"
          :key="b.id"
          :href="b.link_url || undefined"
          class="shrink-0 w-40 sm:w-48"
        >
          <div class="aspect-[4/3] rounded-xl overflow-hidden bg-gray-100 group">
            <picture>
              <source v-if="b.mobile_url" :srcset="resolve(b.mobile_url)" media="(max-width: 640px)" />
              <img :src="resolve(b.desktop_url || b.mobile_url)" :alt="b.title_ru || b.title_tj" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
            </picture>
          </div>
          <div v-if="b.title_tj || b.title_ru" class="mt-1.5 text-sm font-medium text-center text-ink-900 truncate">{{ locale === 'ru' ? b.title_ru : b.title_tj }}</div>
        </component>
      </div>
    </section>

    <!-- Popular categories grid (fallback / always shown) -->
    <section v-if="popularCats.length" class="container-page mt-8">
      <div class="flex items-end justify-between mb-4">
        <h2 class="text-xl md:text-2xl font-bold text-ink-900">{{ t('home.popular_categories') }}</h2>
        <NuxtLink to="/categories" class="text-sm font-medium text-brand-600 hover:text-brand-700">{{ t('common.all') }} →</NuxtLink>
      </div>
      <div class="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-6 gap-3">
        <CategoryCard v-for="c in popularCats" :key="c.id" :category="c" />
      </div>
    </section>

    <!-- Promo banners: left 2 small stacked, right 1 large -->
    <section v-if="midLarge || midSmall.length" class="container-page mt-8">
      <div class="grid lg:grid-cols-2 gap-4">
        <div class="grid grid-rows-2 gap-4">
          <BannerTile v-for="b in midSmall.slice(0, 2)" :key="b.id" :banner="b" ratio="aspect-[16/6]" />
        </div>
        <BannerTile :banner="midLarge" ratio="aspect-[16/9] lg:aspect-auto lg:h-full" />
      </div>
    </section>

    <!-- Featured products -->
    <section v-if="featured.length" class="container-page mt-10">
      <div class="flex items-end justify-between mb-4">
        <h2 class="text-xl md:text-2xl font-bold text-ink-900">{{ t('home.featured_products') }}</h2>
        <NuxtLink to="/products?featured=true" class="text-sm font-medium text-brand-600 hover:text-brand-700">{{ t('common.all') }} →</NuxtLink>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
        <ProductCard v-for="p in featured" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- Latest products -->
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

    <!-- Top sellers -->
    <section v-if="topSellers.length" class="container-page mt-10">
      <div class="flex items-end justify-between mb-4">
        <h2 class="text-xl md:text-2xl font-bold text-ink-900">{{ t('home.top_sellers') }}</h2>
        <NuxtLink to="/sellers" class="text-sm font-medium text-brand-600 hover:text-brand-700">{{ t('common.all') }} →</NuxtLink>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <SellerCard v-for="s in topSellers" :key="s.id" :seller="s" />
      </div>
    </section>

    <!-- Advantages -->
    <section class="container-page mt-12">
      <h2 class="text-xl md:text-2xl font-bold text-center text-ink-900 mb-8">{{ t('home.advantages_title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        <div v-for="i in 4" :key="i" class="card p-6">
          <div class="w-10 h-10 rounded-lg bg-brand-50 text-brand-600 flex items-center justify-center mb-3 font-bold">{{ i }}</div>
          <h3 class="font-semibold text-ink-900">{{ t(`home.advantage_${i}_title`) }}</h3>
          <p class="text-sm text-gray-500 mt-2">{{ t(`home.advantage_${i}_text`) }}</p>
        </div>
      </div>
    </section>

    <!-- Contact -->
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
