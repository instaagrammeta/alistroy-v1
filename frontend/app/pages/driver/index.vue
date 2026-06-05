<script setup lang="ts">
import type { Driver, Order, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'driver', middleware: 'driver' })
const { t } = useI18n()
const { formatPrice, formatDate } = useFormatters()

const tab = ref<'new' | 'assigned' | 'completed'>('assigned')
const statusFilter = computed(() => (tab.value === 'new' ? 'assigned' : tab.value === 'assigned' ? 'on_delivery' : 'completed'))

const { data: meRes } = await useAsyncData('driver-me', () => useApi()<{ data: Driver }>('/driver/me'))
const driver = computed<Driver | null>(() => (meRes.value as any)?.data || null)

const { data, refresh } = await useAsyncData('driver-orders', () =>
  useApi()<PaginatedResponse<Order>>('/driver/orders', { query: { status: statusFilter.value, page_size: 50 } }),
  { watch: [tab] }
)
const orders = computed<Order[]>(() => data.value?.data || [])

const advance = async (o: Order, status: string) => {
  await useApi()(`/driver/orders/${o.id}/status`, { method: 'POST', body: { status } })
  await refresh()
}

const tabs = [
  { key: 'assigned', label: t('driver_panel.assigned') },
  { key: 'new', label: t('driver_panel.new_orders') },
  { key: 'completed', label: t('driver_panel.completed') },
] as const

useHead({ title: () => `${t('nav.driver_panel')} — AliStroy` })
</script>

<template>
  <div>
    <div v-if="driver" class="card p-4 mb-4">
      <div class="font-semibold text-ink-900">{{ driver.full_name }}</div>
      <div class="text-sm text-gray-500">{{ driver.vehicle || t('driver_panel.vehicle') }} · {{ driver.phone }}</div>
    </div>

    <div class="flex gap-2 mb-4">
      <button v-for="tb in tabs" :key="tb.key" class="badge cursor-pointer" :class="tab === tb.key ? 'bg-brand-500 text-white' : 'bg-gray-100 text-gray-700'" @click="tab = tb.key as any">{{ tb.label }}</button>
    </div>

    <EmptyState v-if="!orders.length" />
    <div v-else class="space-y-3">
      <div v-for="o in orders" :key="o.id" class="card p-4">
        <div class="flex items-center justify-between">
          <div class="font-semibold text-ink-900">{{ o.number }}</div>
          <StatusBadge :status="o.status" />
        </div>
        <div class="text-sm text-gray-600 mt-1">{{ o.customer_name }} · {{ o.customer_phone }}</div>
        <div class="text-sm text-gray-500">{{ o.delivery_address }}</div>
        <div class="text-sm mt-1"><strong class="text-brand-600">{{ formatPrice(o.total, o.currency) }}</strong> · {{ formatDate(o.created_at) }}</div>
        <div class="mt-3 flex gap-2">
          <button v-if="o.status === 'assigned'" class="btn-primary btn-sm" @click="advance(o, 'on_delivery')">{{ t('driver_panel.mark_on_delivery') }}</button>
          <button v-if="o.status === 'on_delivery'" class="btn-success btn-sm" @click="advance(o, 'completed')">{{ t('driver_panel.mark_completed') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
