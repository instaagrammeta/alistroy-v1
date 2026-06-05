<script setup lang="ts">
import type { PaginatedResponse, Product } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const { productTitle } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()

const status = ref(typeof route.query.status === 'string' ? route.query.status : '')
const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const page = ref(typeof route.query.page === 'string' ? parseInt(route.query.page, 10) || 1 : 1)
const sort = ref('newest')

const { data, refresh } = await useAsyncData(
  'admin-products',
  () =>
    api<PaginatedResponse<Product>>('/admin/products', {
      query: {
        page: page.value,
        page_size: 20,
        status: status.value || undefined,
        q: q.value || undefined,
        sort: sort.value,
      },
    }),
  { watch: [page, status, q, sort] }
)

const items = computed<Product[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const onPage = (p: number) => {
  page.value = p
  router.replace({ query: { ...route.query, page: p } })
}

const onDelete = async (p: Product) => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  await api(`/admin/products/${p.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('admin.products')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.products') }}</h1>

    <div class="card p-4 mb-4 flex flex-wrap items-end gap-3">
      <div class="flex-1 min-w-[200px]">
        <label class="label">{{ t('common.search') }}</label>
        <input v-model="q" class="input" type="search" />
      </div>
      <div class="w-48">
        <label class="label">Status</label>
        <select v-model="status" class="select">
          <option value="">{{ t('common.all') }}</option>
          <option value="draft">{{ t('seller_panel.status_draft') }}</option>
          <option value="pending">{{ t('seller_panel.status_pending') }}</option>
          <option value="approved">{{ t('seller_panel.status_approved') }}</option>
          <option value="rejected">{{ t('seller_panel.status_rejected') }}</option>
        </select>
      </div>
      <div class="w-48">
        <label class="label">{{ t('catalog.sort') }}</label>
        <select v-model="sort" class="select">
          <option value="newest">{{ t('catalog.sort_newest') }}</option>
          <option value="oldest">{{ t('catalog.sort_oldest') }}</option>
          <option value="price_asc">{{ t('catalog.sort_price_asc') }}</option>
          <option value="price_desc">{{ t('catalog.sort_price_desc') }}</option>
          <option value="popular">{{ t('catalog.sort_popular') }}</option>
        </select>
      </div>
    </div>

    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr>
            <th class="p-3">{{ t('catalog.title') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('catalog.seller') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('catalog.price') }}</th>
            <th class="p-3">Status</th>
            <th class="p-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in items" :key="p.id" class="border-t border-gray-100">
            <td class="p-3">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 rounded-lg bg-gray-50 overflow-hidden flex-shrink-0">
                  <img v-if="p.images?.[0]" :src="resolve(p.images[0].url)" alt="" class="w-full h-full object-cover" />
                </div>
                <div class="min-w-0">
                  <div class="font-medium text-ink-900 truncate">{{ productTitle(p) }}</div>
                  <div class="text-xs text-gray-400 truncate">{{ p.sku }}</div>
                </div>
              </div>
            </td>
            <td class="p-3 hidden md:table-cell">{{ p.seller?.name || '—' }}</td>
            <td class="p-3 hidden md:table-cell whitespace-nowrap">{{ formatPrice(p.price, p.currency) }}</td>
            <td class="p-3"><StatusBadge :status="p.status" /></td>
            <td class="p-3 text-right whitespace-nowrap">
              <NuxtLink :to="`/admin/products/${p.id}`" class="text-brand-600 hover:text-brand-700 mr-3">{{ t('common.view') }}</NuxtLink>
              <button class="text-red-600 hover:text-red-700" @click="onDelete(p)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
  </div>
</template>
