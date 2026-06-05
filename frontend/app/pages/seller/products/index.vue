<script setup lang="ts">
import type { PaginatedResponse, Product } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })
const { t } = useI18n()
const { productName } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()

const page = ref(1)
const status = ref('')

const { data, refresh } = await useAsyncData('seller-products', () =>
  useApi()<PaginatedResponse<Product>>('/seller/products', { query: { page: page.value, page_size: 20, status: status.value || undefined } }),
  { watch: [page, status] }
)
const products = computed<Product[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const onDelete = async (p: Product) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/seller/products/${p.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('seller_panel.products')} — AliStroy` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('seller_panel.products') }}</h1>
      <NuxtLink to="/seller/products/new" class="btn-primary">{{ t('seller_panel.add_product') }}</NuxtLink>
    </div>
    <div class="card p-4 mb-4 w-56">
      <select v-model="status" class="select">
        <option value="">{{ t('common.all') }}</option>
        <option value="draft">{{ t('seller_panel.status_draft') }}</option>
        <option value="pending">{{ t('seller_panel.status_pending') }}</option>
        <option value="approved">{{ t('seller_panel.status_approved') }}</option>
        <option value="rejected">{{ t('seller_panel.status_rejected') }}</option>
      </select>
    </div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr><th class="p-3">{{ t('common.name') }}</th><th class="p-3 hidden md:table-cell">{{ t('common.price') }}</th><th class="p-3 hidden md:table-cell">{{ t('seller_panel.views') }}</th><th class="p-3">{{ t('common.status') }}</th><th class="p-3"></th></tr>
        </thead>
        <tbody>
          <tr v-for="p in products" :key="p.id" class="border-t border-gray-100">
            <td class="p-3">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 rounded-lg bg-gray-50 overflow-hidden shrink-0"><img v-if="p.images?.[0]" :src="resolve(p.images[0].url)" alt="" class="w-full h-full object-cover" /></div>
                <div class="min-w-0"><div class="font-medium text-ink-900 truncate">{{ productName(p) }}</div><div class="text-xs text-gray-400">{{ p.sku }}</div></div>
              </div>
            </td>
            <td class="p-3 hidden md:table-cell whitespace-nowrap">{{ formatPrice(p.sale_price, p.currency) }}</td>
            <td class="p-3 hidden md:table-cell">{{ p.views_count }}</td>
            <td class="p-3"><StatusBadge :status="p.status" /></td>
            <td class="p-3 text-right whitespace-nowrap">
              <NuxtLink :to="`/seller/products/${p.id}/edit`" class="text-brand-600 hover:text-brand-700 mr-3">{{ t('common.edit') }}</NuxtLink>
              <button class="text-red-600 hover:text-red-700" @click="onDelete(p)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!products.length"><td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
