<script setup lang="ts">
import type { PaginatedResponse, Product, SellerTotals } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })

const api = useApi()
const { t } = useI18n()
const { productTitle } = useLocaleField()
const { resolve } = useImageUrl()

const { data: stats } = await useAsyncData('seller-stats-page', () =>
  api<{ data: SellerTotals }>('/seller/stats')
)
const totals = computed<SellerTotals | null>(() => (stats.value as any)?.data || null)

const { data: list } = await useAsyncData('seller-stats-products', () =>
  api<PaginatedResponse<Product>>('/seller/products', {
    query: { sort: 'popular', page_size: 20 },
  })
)
const products = computed<Product[]>(() => (list.value as any)?.data || [])

useHead({ title: () => `${t('seller_panel.stats')} — AliStroy` })

const cards = computed(() => [
  { label: t('seller_panel.stats_views'), value: totals.value?.total_views ?? 0 },
  { label: t('seller_panel.stats_phone'), value: totals.value?.total_phone_clicks ?? 0 },
  { label: t('seller_panel.stats_whatsapp'), value: totals.value?.total_whatsapp_clicks ?? 0 },
])
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('seller_panel.stats') }}</h1>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <div v-for="c in cards" :key="c.label" class="card p-5">
        <div class="text-sm text-gray-500">{{ c.label }}</div>
        <div class="mt-1 text-3xl font-extrabold text-ink-900">{{ c.value }}</div>
      </div>
    </div>

    <h2 class="font-semibold mb-3">{{ t('catalog.sort_popular') }}</h2>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr>
            <th class="p-3">{{ t('catalog.title') }}</th>
            <th class="p-3">{{ t('seller_panel.stats_views') }}</th>
            <th class="p-3">{{ t('seller_panel.stats_phone') }}</th>
            <th class="p-3">{{ t('seller_panel.stats_whatsapp') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in products" :key="p.id" class="border-t border-gray-100">
            <td class="p-3">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded bg-gray-50 overflow-hidden flex-shrink-0">
                  <img v-if="p.images?.[0]" :src="resolve(p.images[0].url)" alt="" class="w-full h-full object-cover" />
                </div>
                <div class="font-medium truncate">{{ productTitle(p) }}</div>
              </div>
            </td>
            <td class="p-3">{{ p.views_count }}</td>
            <td class="p-3">{{ p.phone_clicks }}</td>
            <td class="p-3">{{ p.whatsapp_clicks }}</td>
          </tr>
          <tr v-if="!products.length">
            <td colspan="4" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
