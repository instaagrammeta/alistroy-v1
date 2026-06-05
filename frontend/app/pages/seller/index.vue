<script setup lang="ts">
import type { SellerTotals, Seller } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })

const api = useApi()
const { t } = useI18n()

const { data: stats } = await useAsyncData('seller-stats', () =>
  api<{ data: SellerTotals }>('/seller/stats')
)
const { data: meRes } = await useAsyncData('seller-me', () =>
  api<{ data: Seller }>('/seller/me')
)
const totals = computed<SellerTotals | null>(() => (stats.value as any)?.data || null)
const seller = computed<Seller | null>(() => (meRes.value as any)?.data || null)

useHead({ title: () => `${t('seller_panel.dashboard')} — AliStroy` })

const cards = computed(() => [
  { label: t('seller_panel.stats_total_products'), value: totals.value?.total_products ?? 0 },
  { label: t('seller_panel.stats_approved'), value: totals.value?.approved_products ?? 0 },
  { label: t('seller_panel.stats_pending'), value: totals.value?.pending_products ?? 0 },
  { label: t('seller_panel.stats_views'), value: totals.value?.total_views ?? 0 },
  { label: t('seller_panel.stats_phone'), value: totals.value?.total_phone_clicks ?? 0 },
  { label: t('seller_panel.stats_whatsapp'), value: totals.value?.total_whatsapp_clicks ?? 0 },
])
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-ink-900">{{ seller?.name || t('seller_panel.dashboard') }}</h1>
        <div class="text-sm text-gray-500">{{ seller?.city }}</div>
      </div>
      <NuxtLink to="/seller/products/new" class="btn-primary">{{ t('seller_panel.add_product') }}</NuxtLink>
    </div>

    <div v-if="seller && seller.status !== 'approved'" class="card p-4 mb-6 border-l-4 border-amber-400 bg-amber-50">
      <div class="font-medium text-amber-900">
        {{ t('seller_panel.status_' + seller.status) }}
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <div v-for="c in cards" :key="c.label" class="card p-5">
        <div class="text-sm text-gray-500">{{ c.label }}</div>
        <div class="mt-1 text-3xl font-extrabold text-ink-900">{{ c.value }}</div>
      </div>
    </div>
  </div>
</template>
