<script setup lang="ts">
import type { Brand, Category, PaginatedResponse, Product } from '~/types/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { categoryName } = useLocaleField()

const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const categorySlug = ref(typeof route.query.category === 'string' ? route.query.category : '')
const minPrice = ref(typeof route.query.min_price === 'string' ? route.query.min_price : '')
const maxPrice = ref(typeof route.query.max_price === 'string' ? route.query.max_price : '')
const selectedBrands = ref<string[]>([])
const sort = ref(typeof route.query.sort === 'string' ? route.query.sort : 'newest')
const page = ref(typeof route.query.page === 'string' ? parseInt(route.query.page, 10) || 1 : 1)
const pageSize = 20

const buildQuery = () => {
  const x: Record<string, any> = { page: page.value, page_size: pageSize, sort: sort.value }
  if (q.value) x.q = q.value
  if (categorySlug.value) x.category = categorySlug.value
  if (minPrice.value) x.min_price = minPrice.value
  if (maxPrice.value) x.max_price = maxPrice.value
  if (selectedBrands.value.length) x.brand = selectedBrands.value
  return x
}

const { data: list } = await useAsyncData('products-list', () =>
  $fetch<PaginatedResponse<Product>>('/products', { baseURL: useApiBase(), query: buildQuery() }),
  { watch: [page, sort] }
)
const { data: cats } = await useAsyncData('products-cats', () => $fetch<{ data: Category[] }>('/categories', { baseURL: useApiBase() }))
const { data: brands } = await useAsyncData('products-brands', () => $fetch<{ data: Brand[] }>('/brands', { baseURL: useApiBase() }))

const products = computed<Product[]>(() => list.value?.data || [])
const pagination = computed(() => list.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })
const categories = computed<Category[]>(() => cats.value?.data || [])
const brandList = computed<Brand[]>(() => brands.value?.data || [])

const refreshList = async () => {
  page.value = 1
  router.push({ path: '/products', query: buildQuery() })
  const res = await $fetch<PaginatedResponse<Product>>('/products', { baseURL: useApiBase(), query: buildQuery() })
  list.value = res
}

const reset = () => {
  q.value = ''; categorySlug.value = ''; minPrice.value = ''; maxPrice.value = ''
  selectedBrands.value = []; sort.value = 'newest'; page.value = 1
  refreshList()
}

const onPage = async (p: number) => {
  page.value = p
  router.push({ path: '/products', query: { ...buildQuery(), page: p } })
  const res = await $fetch<PaginatedResponse<Product>>('/products', { baseURL: useApiBase(), query: buildQuery() })
  list.value = res
  if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
}

useHead({ title: () => `${t('catalog.title')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6">
    <h1 class="text-2xl md:text-3xl font-bold text-ink-900 mb-5">{{ t('catalog.title') }}</h1>
    <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-6">
      <aside class="card p-4 h-fit lg:sticky lg:top-20">
        <form class="space-y-4" @submit.prevent="refreshList">
          <div>
            <label class="label">{{ t('common.search') }}</label>
            <input v-model="q" type="search" class="input" :placeholder="t('common.search_placeholder')" />
          </div>
          <div>
            <label class="label">{{ t('catalog.category') }}</label>
            <select v-model="categorySlug" class="select">
              <option value="">{{ t('common.all') }}</option>
              <option v-for="c in categories" :key="c.id" :value="c.slug">{{ categoryName(c) }}</option>
            </select>
          </div>
          <div>
            <label class="label">{{ t('catalog.price_range') }} ({{ t('common.currency') }})</label>
            <div class="grid grid-cols-2 gap-2">
              <input v-model="minPrice" type="number" class="input" :placeholder="t('common.from')" min="0" />
              <input v-model="maxPrice" type="number" class="input" :placeholder="t('common.to')" min="0" />
            </div>
          </div>
          <div v-if="brandList.length">
            <label class="label">{{ t('catalog.brand') }}</label>
            <div class="max-h-44 overflow-auto space-y-1 pr-1">
              <label v-for="b in brandList" :key="b.id" class="flex items-center gap-2 text-sm">
                <input v-model="selectedBrands" type="checkbox" :value="b.id" class="h-4 w-4 rounded border-gray-300" />
                <span>{{ b.name }}</span>
              </label>
            </div>
          </div>
          <div class="flex gap-2">
            <button type="submit" class="btn-primary flex-1">{{ t('common.apply') }}</button>
            <button type="button" class="btn-outline" @click="reset">{{ t('common.reset') }}</button>
          </div>
        </form>
      </aside>

      <section>
        <div class="flex items-center justify-between mb-3 gap-3">
          <div class="text-sm text-gray-500">{{ t('common.results_count', { count: pagination.total }) }}</div>
          <select v-model="sort" class="select max-w-[200px]">
            <option value="newest">{{ t('catalog.sort_newest') }}</option>
            <option value="oldest">{{ t('catalog.sort_oldest') }}</option>
            <option value="price_asc">{{ t('catalog.sort_price_asc') }}</option>
            <option value="price_desc">{{ t('catalog.sort_price_desc') }}</option>
            <option value="popular">{{ t('catalog.sort_popular') }}</option>
          </select>
        </div>
        <EmptyState v-if="!products.length" />
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4">
          <ProductCard v-for="p in products" :key="p.id" :product="p" />
        </div>
        <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
      </section>
    </div>
  </div>
</template>
