<script setup lang="ts">
import type { Order, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const route = useRoute()
const { t } = useI18n()
const { formatPrice, formatDate } = useFormatters()

const status = ref(typeof route.query.status === 'string' ? route.query.status : '')
const q = ref('')
const page = ref(1)
const { data } = await useAsyncData('admin-orders', () =>
  useApi()<PaginatedResponse<Order>>('/admin/orders', { query: { page: page.value, page_size: 20, status: status.value || undefined, q: q.value || undefined } }),
  { watch: [status, q, page] }
)
const orders = computed<Order[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const doExport = async () => {
  const res = await useApi().raw('/admin/export/orders', { query: { status: status.value || undefined }, responseType: 'blob' })
  downloadBlob(res._data as Blob, 'orders.xlsx')
}
useHead({ title: () => `${t('admin.orders')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-2">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.orders') }}</h1>
      <div class="flex gap-2"><button class="btn-outline btn-sm" @click="doExport">{{ t('admin.export') }}</button><NuxtLink to="/admin/orders/new" class="btn-primary btn-sm">+ {{ t('admin.new_order') }}</NuxtLink></div>
    </div>
    <div class="card p-4 mb-4 flex flex-wrap gap-3 items-end">
      <div class="flex-1 min-w-[200px]"><label class="label">{{ t('common.search') }}</label><input v-model="q" class="input" type="search" /></div>
      <div class="w-48"><label class="label">{{ t('common.status') }}</label>
        <select v-model="status" class="select"><option value="">{{ t('common.all') }}</option><option v-for="s in ['new','processing','assigned','on_delivery','completed','cancelled']" :key="s" :value="s">{{ t(`order_status.${s}`) }}</option></select>
      </div>
    </div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500"><tr><th class="p-3">№</th><th class="p-3">{{ t('common.name') }}</th><th class="p-3 hidden md:table-cell">{{ t('common.total') }}</th><th class="p-3">{{ t('common.status') }}</th><th class="p-3 hidden md:table-cell">{{ t('common.date') }}</th><th class="p-3"></th></tr></thead>
        <tbody>
          <tr v-for="o in orders" :key="o.id" class="border-t border-gray-100">
            <td class="p-3 font-medium">{{ o.number }}</td>
            <td class="p-3">{{ o.customer_name }}<div class="text-xs text-gray-400">{{ o.customer_phone }}</div></td>
            <td class="p-3 hidden md:table-cell whitespace-nowrap">{{ formatPrice(o.total, o.currency) }}</td>
            <td class="p-3"><StatusBadge :status="o.status" /></td>
            <td class="p-3 hidden md:table-cell">{{ formatDate(o.created_at) }}</td>
            <td class="p-3 text-right"><NuxtLink :to="`/admin/orders/${o.id}`" class="text-brand-600">{{ t('common.view') }}</NuxtLink></td>
          </tr>
          <tr v-if="!orders.length"><td colspan="6" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
