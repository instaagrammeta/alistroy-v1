<script setup lang="ts">
import type { Category, PaginatedResponse, Product } from '~/types/api'

const route = useRoute()
const router = useRouter()
const slug = route.params.slug as string
const config = useRuntimeConfig()
const { t } = useI18n()
const { categoryTitle } = useLocaleField()

const page = ref(typeof route.query.page === 'string' ? parseInt(route.query.page, 10) || 1 : 1)
const sort = ref(typeof route.query.sort === 'string' ? route.query.sort : 'newest')
const pageSize = 24

const { data: cat } = await useAsyncData(`cat-${slug}`, () =>
  $fetch<{ data: Category }>(`/categories/${slug}`, { baseURL: config.public.apiBase })
)
const category = computed<Category | null>(() => cat.value?.data || null)

const { data: list } = await useAsyncData(
  `cat-products-${slug}`,
  () =>
    $fetch<PaginatedResponse<Product>>('/products', {
      baseURL: config.public.apiBase,
      query: { category: slug, page: page.value, page_size: pageSize, sort: sort.value },
    }),
  { watch: [page, sort] }
)

const products = computed<Product[]>(() => list.value?.data || [])
const pagination = computed(() => list.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

const onPage = (p: number) => {
  page.value = p
  router.push({ query: { ...route.query, page: p } })
  if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
}

useHead(() => ({
  title: `${category.value ? categoryTitle(category.value) : t('nav.categories')} — AliStroy`,
}))
</script>

<template>
  <div class="container-page py-8">
    <NuxtLink to="/categories" class="text-sm text-brand-600 hover:text-brand-700">← {{ t('nav.categories') }}</NuxtLink>
    <h1 class="mt-2 text-3xl font-bold text-ink-900">{{ category ? categoryTitle(category) : '' }}</h1>

    <div class="flex items-center justify-between mt-6 mb-3">
      <div class="text-sm text-gray-500">{{ t('common.results_count', { count: pagination.total }) }}</div>
      <select v-model="sort" class="select max-w-xs">
        <option value="newest">{{ t('catalog.sort_newest') }}</option>
        <option value="oldest">{{ t('catalog.sort_oldest') }}</option>
        <option value="price_asc">{{ t('catalog.sort_price_asc') }}</option>
        <option value="price_desc">{{ t('catalog.sort_price_desc') }}</option>
        <option value="popular">{{ t('catalog.sort_popular') }}</option>
      </select>
    </div>

    <EmptyState v-if="!products.length" />
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
  </div>
</template>
