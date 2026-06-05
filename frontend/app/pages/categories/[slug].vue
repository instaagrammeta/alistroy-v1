<script setup lang="ts">
import type { Category, PaginatedResponse, Product } from '~/types/api'

const route = useRoute()
const router = useRouter()
const slug = route.params.slug as string
const { t } = useI18n()
const { categoryName } = useLocaleField()

const subcategory = ref(typeof route.query.subcategory === 'string' ? route.query.subcategory : '')
const sort = ref('newest')
const page = ref(1)
const pageSize = 24

const { data: cat } = await useAsyncData(`cat-${slug}`, () => $fetch<{ data: Category }>(`/categories/${slug}`, { baseURL: useApiBase() }))
const category = computed<Category | null>(() => cat.value?.data || null)

const fetchList = () => $fetch<PaginatedResponse<Product>>('/products', {
  baseURL: useApiBase(),
  query: { category: slug, subcategory: subcategory.value || undefined, page: page.value, page_size: pageSize, sort: sort.value },
})
const { data: list } = await useAsyncData(`cat-products-${slug}`, fetchList, { watch: [page, sort, subcategory] })

const products = computed<Product[]>(() => list.value?.data || [])
const pagination = computed(() => list.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

const onPage = (p: number) => {
  page.value = p
  if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
}
const pickSub = (slugVal: string) => {
  subcategory.value = subcategory.value === slugVal ? '' : slugVal
  page.value = 1
  router.replace({ query: subcategory.value ? { subcategory: subcategory.value } : {} })
}

useHead(() => ({ title: `${category.value ? categoryName(category.value) : t('nav.categories')} — AliStroy` }))
</script>

<template>
  <div class="container-page py-6">
    <NuxtLink to="/categories" class="text-sm text-brand-600 hover:text-brand-700">← {{ t('nav.categories') }}</NuxtLink>
    <h1 class="mt-2 text-2xl md:text-3xl font-bold text-ink-900">{{ category ? categoryName(category) : '' }}</h1>

    <div v-if="category?.subcategories && category.subcategories.length" class="mt-4 flex flex-wrap gap-2">
      <button class="badge" :class="!subcategory ? 'bg-brand-500 text-white' : 'bg-gray-100 text-gray-700'" @click="pickSub('')">{{ t('common.all') }}</button>
      <button v-for="sub in category.subcategories" :key="sub.id" class="badge" :class="subcategory === sub.slug ? 'bg-brand-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-brand-50'" @click="pickSub(sub.slug)">{{ categoryName(sub) }}</button>
    </div>

    <div class="flex items-center justify-between mt-6 mb-3">
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
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
  </div>
</template>
