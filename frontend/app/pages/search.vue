<script setup lang="ts">
import type { PaginatedResponse, Product } from '~/types/api'

const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const { t } = useI18n()

const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const page = ref(typeof route.query.page === 'string' ? parseInt(route.query.page, 10) || 1 : 1)
const pageSize = 24

const { data } = await useAsyncData(
  'search',
  () =>
    $fetch<PaginatedResponse<Product>>('/products', {
      baseURL: useApiBase(),
      query: { q: q.value || undefined, page: page.value, page_size: pageSize },
    }),
  { watch: [q, page] }
)

const products = computed<Product[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

watch(q, () => {
  router.replace({ path: '/search', query: { q: q.value, page: 1 } })
  page.value = 1
})

const onPage = (p: number) => {
  page.value = p
  router.push({ path: '/search', query: { q: q.value, page: p } })
}

useHead(() => ({ title: () => `${t('search.title')} — AliStroy` }))
</script>

<template>
  <div class="container-page py-8">
    <h1 class="text-3xl font-bold text-ink-900 mb-2">{{ t('search.title') }}</h1>
    <div v-if="q" class="text-gray-500 mb-6">{{ t('search.for', { q }) }}</div>

    <input v-model="q" type="search" class="input max-w-xl mb-6" :placeholder="t('common.search_placeholder')" />

    <div class="text-sm text-gray-500 mb-3">{{ t('common.results_count', { count: pagination.total }) }}</div>
    <EmptyState v-if="!products.length" />
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
  </div>
</template>
