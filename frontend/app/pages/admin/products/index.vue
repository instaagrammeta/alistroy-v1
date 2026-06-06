<script setup lang="ts">
import type { PaginatedResponse, Product } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const route = useRoute()
const { t } = useI18n()
const { productName } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()

const status = ref(typeof route.query.status === 'string' ? route.query.status : '')
const q = ref('')
const page = ref(1)

const { data, refresh } = await useAsyncData('admin-products', () =>
  useApi()<PaginatedResponse<Product>>('/admin/products', { query: { page: page.value, page_size: 20, status: status.value || undefined, q: q.value || undefined } }),
  { watch: [page, status, q] }
)
const items = computed<Product[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const onDelete = async (p: Product) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/products/${p.id}`, { method: 'DELETE' })
  await refresh()
}
const toggleAvailable = async (p: Product) => {
  await useApi()(`/admin/products/${p.id}`, { method: 'PATCH', body: { is_available: !p.is_available } })
  await refresh()
}
const doExport = async () => {
  const res = await useApi().raw(`/admin/export/products`, { query: { status: status.value || undefined, q: q.value || undefined }, responseType: 'blob' })
  downloadBlob(res._data as Blob, 'products.xlsx')
}

useHead({ title: () => `${t('admin.products')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-2">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.products') }}</h1>
      <div class="flex gap-2">
        <button class="btn-outline btn-sm inline-flex items-center gap-1" @click="doExport">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          {{ t('admin.export') }}
        </button>
        <NuxtLink to="/admin/products/new" class="btn-primary btn-sm inline-flex items-center gap-1">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          {{ t('admin.new_product') }}
        </NuxtLink>
      </div>
    </div>
    <div class="card p-4 mb-4 flex flex-wrap items-end gap-3">
      <div class="flex-1 min-w-[200px]"><label class="label">{{ t('common.search') }}</label><input v-model="q" type="search" class="input" /></div>
      <div class="w-48"><label class="label">{{ t('common.status') }}</label>
        <select v-model="status" class="select">
          <option value="">{{ t('common.all') }}</option>
          <option value="draft">{{ t('seller_panel.status_draft') }}</option>
          <option value="pending">{{ t('seller_panel.status_pending') }}</option>
          <option value="approved">{{ t('seller_panel.status_approved') }}</option>
          <option value="rejected">{{ t('seller_panel.status_rejected') }}</option>
        </select>
      </div>
    </div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500"><tr><th class="p-3">{{ t('common.name') }}</th><th class="p-3 hidden md:table-cell">{{ t('catalog.seller') }}</th><th class="p-3 hidden md:table-cell">{{ t('common.price') }}</th><th class="p-3">{{ t('common.status') }}</th><th class="p-3 text-right">{{ t('common.actions') }}</th></tr></thead>
        <tbody>
          <tr v-for="p in items" :key="p.id" class="border-t border-gray-100">
            <td class="p-3"><div class="flex items-center gap-3"><div class="w-12 h-12 rounded-lg bg-gray-50 overflow-hidden shrink-0"><img v-if="p.images?.[0]" :src="resolve(p.images[0].url)" alt="" class="w-full h-full object-cover" /></div><div class="min-w-0"><div class="font-medium truncate">{{ productName(p) }}</div><div class="text-xs text-gray-400">{{ p.sku }}</div></div></div></td>
            <td class="p-3 hidden md:table-cell">{{ p.seller?.market_name || p.seller?.full_name || '—' }}</td>
            <td class="p-3 hidden md:table-cell whitespace-nowrap">{{ formatPrice(p.sale_price, p.currency) }}</td>
            <td class="p-3"><StatusBadge :status="p.status" /></td>
            <td class="p-3 text-right">
              <div class="inline-flex items-center gap-1">
                <IconButton :variant="p.is_available ? 'toggle-on' : 'toggle-off'" :title="p.is_available ? t('common.active') : t('common.inactive')" @click="toggleAvailable(p)" />
                <NuxtLink :to="`/admin/products/${p.id}`" title="View"><IconButton variant="view" /></NuxtLink>
                <NuxtLink :to="`/admin/products/${p.id}/edit`" title="Edit"><IconButton variant="edit" /></NuxtLink>
                <IconButton variant="delete" :title="t('common.delete')" @click="onDelete(p)" />
              </div>
            </td>
          </tr>
          <tr v-if="!items.length"><td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
