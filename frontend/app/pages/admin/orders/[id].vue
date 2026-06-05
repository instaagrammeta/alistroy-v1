<script setup lang="ts">
import type { Driver, Order, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const route = useRoute()
const { t } = useI18n()
const { formatPrice, formatDateTime } = useFormatters()
const id = route.params.id as string

const { data, refresh } = await useAsyncData(`admin-order-${id}`, () => useApi()<{ data: Order }>(`/admin/orders/${id}`))
const order = computed<Order | null>(() => (data.value as any)?.data || null)
if (!order.value) throw createError({ statusCode: 404, statusMessage: 'Not found' })

const { data: driversRes } = await useAsyncData('order-drivers', () => useApi()<PaginatedResponse<Driver>>('/admin/drivers', { query: { page_size: 200, active: true } }))
const drivers = computed<Driver[]>(() => (driversRes.value as any)?.data || [])

const status = ref(order.value.status)
const driverId = ref(order.value.driver_id || '')
const saving = ref(false)

const update = async () => {
  saving.value = true
  try {
    await useApi()(`/admin/orders/${id}/status`, { method: 'POST', body: { status: status.value, driver_id: driverId.value || undefined } })
    await refresh()
  } finally {
    saving.value = false
  }
}

const receiptUrl = computed(() => `${useApiBase()}/admin/orders/${id}/receipt`)
const openReceipt = async () => {
  const res = await useApi().raw(`/admin/orders/${id}/receipt`, { responseType: 'blob' })
  const url = URL.createObjectURL(res._data as Blob)
  const w = window.open(url, '_blank')
  if (w) w.onload = () => setTimeout(() => w.print(), 400)
}

useHead({ title: () => `${order.value?.number} — Admin` })
</script>

<template>
  <div v-if="order">
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ order.number }}</h1>
      <div class="flex items-center gap-2"><StatusBadge :status="order.status" /><button class="btn-outline btn-sm" @click="openReceipt">🖨 {{ t('admin.receipt') }}</button></div>
    </div>
    <div class="grid lg:grid-cols-[1fr_320px] gap-6">
      <div class="card p-5">
        <div class="grid grid-cols-2 gap-2 text-sm mb-4">
          <div><span class="text-gray-400">{{ t('common.name') }}:</span> {{ order.customer_name }}</div>
          <div><span class="text-gray-400">{{ t('admin.phone') }}:</span> {{ order.customer_phone }}</div>
          <div class="col-span-2"><span class="text-gray-400">{{ t('cart.delivery_address') }}:</span> {{ order.delivery_address }}</div>
        </div>
        <table class="w-full text-sm">
          <thead class="text-left text-xs uppercase text-gray-500 border-b"><tr><th class="py-2">{{ t('common.name') }}</th><th class="py-2 text-right">{{ t('common.quantity') }}</th><th class="py-2 text-right">{{ t('common.price') }}</th><th class="py-2 text-right">{{ t('common.total') }}</th></tr></thead>
          <tbody>
            <tr v-for="it in order.items" :key="it.id" class="border-b border-gray-50"><td class="py-2">{{ it.name_snapshot }}</td><td class="py-2 text-right">{{ it.quantity }}</td><td class="py-2 text-right">{{ formatPrice(it.sale_price, order.currency) }}</td><td class="py-2 text-right">{{ formatPrice(it.line_total, order.currency) }}</td></tr>
          </tbody>
        </table>
        <div class="mt-4 text-right space-y-1">
          <div class="text-sm text-gray-500">{{ t('cart.subtotal') }}: {{ formatPrice(order.subtotal, order.currency) }}</div>
          <div class="text-lg font-bold text-brand-600">{{ t('common.total') }}: {{ formatPrice(order.total, order.currency) }}</div>
          <div class="text-sm text-emerald-600">{{ t('admin.profit') }}: {{ formatPrice(order.profit, order.currency) }}</div>
        </div>
      </div>

      <div class="card p-5 h-fit space-y-3">
        <div><label class="label">{{ t('common.status') }}</label><select v-model="status" class="select"><option v-for="s in ['new','processing','assigned','on_delivery','completed','cancelled']" :key="s" :value="s">{{ t(`order_status.${s}`) }}</option></select></div>
        <div><label class="label">{{ t('admin.assign_driver') }}</label><select v-model="driverId" class="select"><option value="">—</option><option v-for="d in drivers" :key="d.id" :value="d.id">{{ d.full_name }}</option></select></div>
        <button class="btn-primary w-full" :disabled="saving" @click="update">{{ t('common.save') }}</button>
      </div>
    </div>
  </div>
</template>
