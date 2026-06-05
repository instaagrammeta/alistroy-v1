<script setup lang="ts">
import type { PaginatedResponse, Product, Seller } from '~/types/api'

const route = useRoute()
const slug = route.params.slug as string
const { t } = useI18n()
const { resolve } = useImageUrl()
const { formatDate } = useFormatters()

const sort = ref('newest')
const page = ref(1)
const pageSize = 12

const { data: sellerRes, error } = await useAsyncData(`seller-${slug}`, () => $fetch<{ data: Seller }>(`/sellers/${slug}`, { baseURL: useApiBase() }))
if (error.value) throw createError({ statusCode: 404, statusMessage: 'Seller not found' })
const seller = computed<Seller | null>(() => sellerRes.value?.data || null)

const { data: list } = await useAsyncData(`seller-${slug}-products`, () =>
  $fetch<PaginatedResponse<Product>>('/products', { baseURL: useApiBase(), query: { seller: slug, page: page.value, page_size: pageSize, sort: sort.value } }),
  { watch: [page, sort] }
)
const products = computed<Product[]>(() => list.value?.data || [])
const pagination = computed(() => list.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

useHead(() => ({ title: `${seller.value?.market_name || seller.value?.full_name || ''} — AliStroy` }))
</script>

<template>
  <div v-if="seller" class="container-page py-6">
    <div class="card p-6 flex flex-col sm:flex-row items-start sm:items-center gap-5">
      <div class="w-20 h-20 rounded-2xl bg-gray-100 overflow-hidden flex items-center justify-center">
        <img v-if="seller.logo_url" :src="resolve(seller.logo_url)" :alt="seller.full_name" class="w-full h-full object-cover" />
        <span v-else class="text-3xl font-bold text-brand-600">{{ (seller.market_name || seller.full_name).charAt(0).toUpperCase() }}</span>
      </div>
      <div class="min-w-0 flex-1">
        <h1 class="text-2xl font-bold text-ink-900">{{ seller.market_name || seller.company_name || seller.full_name }}</h1>
        <div class="text-sm text-gray-500 mt-1">
          <span v-if="seller.city">{{ seller.city }}</span>
          <span v-if="seller.address" class="ml-2">• {{ seller.address }}</span>
        </div>
        <div class="text-xs text-gray-400 mt-1">{{ t('seller.since', { date: formatDate(seller.created_at) }) }}</div>
      </div>
      <div class="flex flex-col items-start sm:items-end gap-1 text-sm">
        <a v-if="seller.phone" :href="`tel:${seller.phone.replace(/[^0-9+]/g, '')}`" class="text-brand-600 hover:text-brand-700">{{ seller.phone }}</a>
        <a v-if="seller.whatsapp" :href="`https://wa.me/${seller.whatsapp.replace(/[^0-9]/g, '')}`" target="_blank" class="text-emerald-600 hover:text-emerald-700">WhatsApp</a>
      </div>
    </div>

    <div class="flex items-center justify-between mt-8 mb-4">
      <h2 class="text-xl font-bold">{{ t('seller.products') }}</h2>
      <select v-model="sort" class="select max-w-[200px]">
        <option value="newest">{{ t('catalog.sort_newest') }}</option>
        <option value="price_asc">{{ t('catalog.sort_price_asc') }}</option>
        <option value="price_desc">{{ t('catalog.sort_price_desc') }}</option>
      </select>
    </div>
    <EmptyState v-if="!products.length" />
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
