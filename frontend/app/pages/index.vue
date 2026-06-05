<script setup lang="ts">
import type { Category, Product, Seller } from '~/types/api'

const { t, locale } = useI18n()
const settings = useSettingsStore()
const config = useRuntimeConfig()

const baseFetch = (url: string, query?: Record<string, any>) =>
  $fetch(url, { baseURL: config.public.apiBase, query })

const { data: featured } = await useAsyncData('home-featured', () =>
  baseFetch('/products', { featured: true, page_size: 8, sort: 'newest' })
)
const { data: latest } = await useAsyncData('home-latest', () =>
  baseFetch('/products', { page_size: 8, sort: 'newest' })
)
const { data: cats } = await useAsyncData('home-cats', () =>
  baseFetch('/categories/popular', { limit: 8 })
)
const { data: top } = await useAsyncData('home-top-sellers', () =>
  baseFetch('/sellers/top', { limit: 6 })
)

const featuredProducts = computed<Product[]>(() => (featured.value as any)?.data || [])
const latestProducts = computed<Product[]>(() => (latest.value as any)?.data || [])
const popularCats = computed<Category[]>(() => (cats.value as any)?.data || [])
const topSellers = computed<Seller[]>(() => (top.value as any)?.data || [])

useHead(() => ({
  title: `${settings.siteName(locale.value as 'tg' | 'ru')} — ${settings.tagline(locale.value as 'tg' | 'ru')}`,
}))
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="bg-gradient-to-br from-brand-50 via-white to-white">
      <div class="container-page py-12 md:py-20 grid md:grid-cols-2 gap-10 items-center">
        <div>
          <h1 class="text-3xl md:text-5xl font-extrabold text-ink-900 leading-tight">
            {{ settings.heroTitle(locale as 'tg' | 'ru') || t('home.hero_title') }}
          </h1>
          <p class="mt-4 text-base md:text-lg text-gray-600 max-w-xl">
            {{ settings.heroSubtitle(locale as 'tg' | 'ru') || t('home.hero_subtitle') }}
          </p>
          <div class="mt-7 flex flex-wrap gap-3">
            <NuxtLink to="/products" class="btn-primary">{{ t('home.hero_cta') }}</NuxtLink>
            <NuxtLink to="/register?role=seller" class="btn-outline">{{ t('nav.become_seller') }}</NuxtLink>
          </div>
        </div>
        <div class="hidden md:block">
          <div class="grid grid-cols-3 gap-3">
            <div class="aspect-square rounded-xl bg-brand-100"></div>
            <div class="aspect-square rounded-xl bg-brand-200 mt-6"></div>
            <div class="aspect-square rounded-xl bg-brand-300"></div>
            <div class="aspect-square rounded-xl bg-brand-200"></div>
            <div class="aspect-square rounded-xl bg-brand-300 mt-6"></div>
            <div class="aspect-square rounded-xl bg-brand-100"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories -->
    <section v-if="popularCats.length" class="container-page py-12">
      <div class="flex items-end justify-between mb-6">
        <h2 class="text-2xl font-bold text-ink-900">{{ t('home.popular_categories') }}</h2>
        <NuxtLink to="/categories" class="text-sm font-medium text-brand-600 hover:text-brand-700">
          {{ t('common.all') }} →
        </NuxtLink>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
        <CategoryCard v-for="c in popularCats" :key="c.id" :category="c" />
      </div>
    </section>

    <!-- Featured -->
    <section v-if="featuredProducts.length" class="container-page py-6">
      <div class="flex items-end justify-between mb-6">
        <h2 class="text-2xl font-bold text-ink-900">{{ t('home.featured_products') }}</h2>
        <NuxtLink to="/products?featured=true" class="text-sm font-medium text-brand-600 hover:text-brand-700">
          {{ t('common.all') }} →
        </NuxtLink>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
        <ProductCard v-for="p in featuredProducts" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- Latest -->
    <section class="container-page py-12">
      <div class="flex items-end justify-between mb-6">
        <h2 class="text-2xl font-bold text-ink-900">{{ t('home.latest_products') }}</h2>
        <NuxtLink to="/products" class="text-sm font-medium text-brand-600 hover:text-brand-700">
          {{ t('common.all') }} →
        </NuxtLink>
      </div>
      <EmptyState v-if="!latestProducts.length" />
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
        <ProductCard v-for="p in latestProducts" :key="p.id" :product="p" />
      </div>
    </section>

    <!-- Top Sellers -->
    <section v-if="topSellers.length" class="bg-white border-y border-gray-100">
      <div class="container-page py-12">
        <div class="flex items-end justify-between mb-6">
          <h2 class="text-2xl font-bold text-ink-900">{{ t('home.top_sellers') }}</h2>
          <NuxtLink to="/sellers" class="text-sm font-medium text-brand-600 hover:text-brand-700">
            {{ t('common.all') }} →
          </NuxtLink>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <SellerCard v-for="s in topSellers" :key="s.id" :seller="s" />
        </div>
      </div>
    </section>

    <!-- Advantages -->
    <section class="container-page py-16">
      <h2 class="text-2xl md:text-3xl font-bold text-center text-ink-900 mb-10">{{ t('home.advantages_title') }}</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="i in 4" :key="i" class="card p-6">
          <div class="w-10 h-10 rounded-lg bg-brand-50 text-brand-600 flex items-center justify-center mb-3 font-bold">{{ i }}</div>
          <h3 class="font-semibold text-ink-900">{{ t(`home.advantage_${i}_title`) }}</h3>
          <p class="text-sm text-gray-500 mt-2">{{ t(`home.advantage_${i}_text`) }}</p>
        </div>
      </div>
    </section>

    <!-- Contact -->
    <section class="bg-ink-900 text-white">
      <div class="container-page py-12 text-center">
        <h2 class="text-2xl md:text-3xl font-bold">{{ t('home.contact_title') }}</h2>
        <p class="mt-2 text-gray-300 max-w-xl mx-auto">{{ t('home.contact_text') }}</p>
        <div class="mt-6 flex flex-wrap justify-center gap-3">
          <a v-if="settings.marketplacePhone" :href="`tel:${settings.marketplacePhone.replace(/[^0-9+]/g, '')}`" class="btn-primary">
            {{ settings.marketplacePhone }}
          </a>
          <a v-if="settings.marketplaceWhatsApp" :href="`https://wa.me/${settings.marketplaceWhatsApp.replace(/[^0-9]/g, '')}`" target="_blank" class="btn-success">
            WhatsApp
          </a>
        </div>
      </div>
    </section>
  </div>
</template>
