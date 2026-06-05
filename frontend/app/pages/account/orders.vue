<script setup lang="ts">
import type { Order, PaginatedResponse } from '~/types/api'

definePageMeta({ middleware: 'auth' })
const { t } = useI18n()
const { formatPrice, formatDate } = useFormatters()
const page = ref(1)

const { data } = await useAsyncData('my-orders', () =>
  useApi()<PaginatedResponse<Order>>('/customer/orders', { query: { page: page.value, page_size: 20 } }),
  { watch: [page] }
)
const orders = computed<Order[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

useHead({ title: () => `${t('nav.orders')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6">
    <h1 class="text-2xl md:text-3xl font-bold text-ink-900 mb-6">{{ t('nav.orders') }}</h1>
    <EmptyState v-if="!orders.length" />
    <div v-else class="space-y-3">
      <div v-for="o in orders" :key="o.id" class="card p-4">
        <div class="flex items-center justify-between">
          <div class="font-semibold text-ink-900">{{ o.number }}</div>
          <StatusBadge :status="o.status" />
        </div>
        <div class="text-sm text-gray-500 mt-1">{{ formatDate(o.created_at) }}</div>
        <div class="mt-2 text-sm text-gray-600">{{ o.items?.length || 0 }} × — <strong class="text-brand-600">{{ formatPrice(o.total, o.currency) }}</strong></div>
      </div>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
