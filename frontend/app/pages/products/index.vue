<script setup lang="ts">
import type { Category, PaginatedResponse, Product } from '~/types/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()

const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const categorySlug = ref(typeof route.query.category === 'string' ? route.query.category : '')
const minPrice = ref(typeof route.query.min_price === 'string' ? route.query.min_price : '')
const maxPrice = ref(typeof route.query.max_price === 'string' ? route.query.max_price : '')
const sort = ref(typeof route.query.sort === 'string' ? route.query.sort : 'newest')
const page = ref(typeof route.query.page === 'string' ? parseInt(route.query.page, 10) || 1 : 1)
const pageSize = 20

const buildQuery = () => {
  const x: Record<string, any> = { page: page.value, page_size: pageSize, sort: sort.value }
  if (q.value) x.q = q.value
  if (categorySlug.value) x.category = categorySlug.value
  if (minPrice.value) x.min_price = minPrice.value
  if (maxPrice.value) x.max_price = maxPrice.value
  return x
}

const fetchProducts = () =>
  $fetch<PaginatedResponse<Product>>('/products', {
    baseURL: config.public.apiBase,
    query: buildQuery(),
  })

const fetchCats = () =>
  $fetch<{ data: Category[] }>('/categories', { baseURL: config.public.apiBase })

const { data: list, refresh } = await useAsyncData('products-list', fetchProducts, {
  watch: [page, sort, q, categorySlug, minPrice, maxPrice],
})
const { data: cats } = await useAsyncData('products-cats', fetchCats)

const products = computed<Product[]>(() => list.value?.data || [])
const pagination = computed(() => list.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })
const categories = computed<Category[]>(() => cats.value?.data || [])

const { categoryTitle } = useLocaleField()

const submitFilters = () => {
  page.value = 1
  router.push({ path: '/products', query: buildQuery() })
}

const resetFilters = () => {
  q.value = ''
  categorySlug.value = ''
  minPrice.value = ''
  maxPrice.value = ''
  sort.value = 'newest'
  page.value = 1
  router.push({ path: '/products', query: {} })
}

const onPageChange = (p: number) => {
  page.value = p
  router.push({ path: '/products', query: { ...buildQuery(), page: p } })
  if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
}

useHead({ title: () => `${t('catalog.title')} — AliStroy` })
</script>

<template>
  <div class="container-page py-8">
    <h1 class="text-3xl font-bold text-ink-900 mb-6">{{ t('catalog.title') }}</h1>

    <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-6">
      <aside class="card p-4 h-fit">
        <form class="space-y-4" @submit.prevent="submitFilters">
          <div>
            <label class="label">{{ t('common.search') }}</label>
            <input v-model="q" type="search" class="input" :placeholder="t('common.search_placeholder')" />
          </div>
          <div>
            <label class="label">{{ t('catalog.category') }}</label>
            <select v-model="categorySlug" class="select">
              <option value="">{{ t('common.all') }}</option>
              <option v-for="c in categories" :key="c.id" :value="c.slug">{{ categoryTitle(c) }}</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="label">{{ t('common.from') }}</label>
              <input v-model="minPrice" type="number" class="input" min="0" />
            </div>
            <div>
              <label class="label">{{ t('common.to') }}</label>
              <input v-model="maxPrice" type="number" class="input" min="0" />
            </div>
          </div>
          <div>
            <label class="label">{{ t('catalog.sort') }}</label>
            <select v-model="sort" class="select">
              <option value="newest">{{ t('catalog.sort_newest') }}</option>
              <option value="oldest">{{ t('catalog.sort_oldest') }}</option>
              <option value="price_asc">{{ t('catalog.sort_price_asc') }}</option>
              <option value="price_desc">{{ t('catalog.sort_price_desc') }}</option>
              <option value="popular">{{ t('catalog.sort_popular') }}</option>
            </select>
          </div>
          <div class="flex gap-2">
            <button type="submit" class="btn-primary flex-1">{{ t('common.apply') }}</button>
            <button type="button" class="btn-outline" @click="resetFilters">{{ t('common.reset') }}</button>
          </div>
        </form>
      </aside>

      <section>
        <div class="flex items-center justify-between mb-3">
          <div class="text-sm text-gray-500">
            {{ t('common.results_count', { count: pagination.total }) }}
          </div>
        </div>

        <EmptyState v-if="!products.length" />
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4">
          <ProductCard v-for="p in products" :key="p.id" :product="p" />
        </div>

        <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPageChange" />
      </section>
    </div>
  </div>
</template>
