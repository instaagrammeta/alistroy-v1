<script setup lang="ts">
import type { AdminTotals } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { formatPrice } = useFormatters()

const { data } = await useAsyncData('admin-totals', () => useApi()<{ data: AdminTotals }>('/admin/dashboard'))
const totals = computed<AdminTotals | null>(() => (data.value as any)?.data || null)

const cards = computed(() => [
  { label: t('admin.products'), value: totals.value?.total_products ?? 0, link: '/admin/products' },
  { label: t('seller_panel.pending'), value: totals.value?.pending_products ?? 0, link: '/admin/products?status=pending' },
  { label: t('admin.orders'), value: totals.value?.total_orders ?? 0, link: '/admin/orders' },
  { label: t('order_status.new'), value: totals.value?.new_orders ?? 0, link: '/admin/orders?status=new' },
  { label: t('admin.sellers'), value: totals.value?.total_sellers ?? 0, link: '/admin/sellers' },
  { label: t('admin.customers'), value: totals.value?.total_customers ?? 0, link: '/admin/customers' },
  { label: t('admin.drivers'), value: totals.value?.total_drivers ?? 0, link: '/admin/drivers' },
  { label: t('seller_panel.views'), value: totals.value?.total_views ?? 0 },
])

useHead({ title: () => `${t('admin.dashboard')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('admin.dashboard') }}</h1>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <component :is="c.link ? resolveComponent('NuxtLink') : 'div'" v-for="c in cards" :key="c.label" :to="c.link" class="card p-5 hover:shadow-card transition-shadow block">
        <div class="text-sm text-gray-500">{{ c.label }}</div>
        <div class="mt-1 text-3xl font-extrabold text-ink-900">{{ c.value }}</div>
      </component>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.income') }}</div><div class="mt-1 text-2xl font-extrabold text-emerald-600">{{ formatPrice(totals?.total_revenue ?? 0) }}</div></div>
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.profit') }}</div><div class="mt-1 text-2xl font-extrabold text-brand-600">{{ formatPrice(totals?.total_profit ?? 0) }}</div></div>
      <div class="card p-5"><div class="text-sm text-gray-500">WhatsApp / Tel</div><div class="mt-1 text-2xl font-extrabold text-ink-900">{{ totals?.total_whatsapp_clicks ?? 0 }} / {{ totals?.total_phone_clicks ?? 0 }}</div></div>
    </div>
  </div>
</template>
