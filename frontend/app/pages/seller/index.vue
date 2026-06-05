<script setup lang="ts">
import type { Seller, SellerTotals } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })
const { t } = useI18n()

const { data: stats } = await useAsyncData('seller-stats', () => useApi()<{ data: SellerTotals }>('/seller/stats'))
const { data: meRes } = await useAsyncData('seller-me', () => useApi()<{ data: Seller }>('/seller/me'))
const totals = computed<SellerTotals | null>(() => (stats.value as any)?.data || null)
const seller = computed<Seller | null>(() => (meRes.value as any)?.data || null)

const cards = computed(() => [
  { label: t('seller_panel.total_products'), value: totals.value?.total_products ?? 0 },
  { label: t('seller_panel.approved'), value: totals.value?.approved_products ?? 0 },
  { label: t('seller_panel.pending'), value: totals.value?.pending_products ?? 0 },
  { label: t('seller_panel.low_stock'), value: totals.value?.low_stock ?? 0 },
  { label: t('seller_panel.views'), value: totals.value?.total_views ?? 0 },
  { label: t('seller_panel.phone_clicks'), value: totals.value?.total_phone_clicks ?? 0 },
  { label: t('seller_panel.whatsapp_clicks'), value: totals.value?.total_whatsapp_clicks ?? 0 },
  { label: t('seller_panel.telegram_clicks'), value: totals.value?.total_telegram_clicks ?? 0 },
])

useHead({ title: () => `${t('seller_panel.dashboard')} — AliStroy` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-ink-900">{{ seller?.market_name || seller?.full_name || t('seller_panel.dashboard') }}</h1>
        <div class="text-sm text-gray-500">{{ seller?.city }}</div>
      </div>
      <NuxtLink to="/seller/products/new" class="btn-primary">{{ t('seller_panel.add_product') }}</NuxtLink>
    </div>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div v-for="c in cards" :key="c.label" class="card p-5">
        <div class="text-sm text-gray-500">{{ c.label }}</div>
        <div class="mt-1 text-3xl font-extrabold text-ink-900">{{ c.value }}</div>
      </div>
    </div>
  </div>
</template>
