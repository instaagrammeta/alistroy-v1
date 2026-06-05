<script setup lang="ts">
import type { AdminTotals } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const { t } = useI18n()

const { data } = await useAsyncData('admin-totals', () =>
  api<{ data: AdminTotals }>('/admin/totals')
)
const totals = computed<AdminTotals | null>(() => (data.value as any)?.data || null)

useHead({ title: () => `${t('admin.dashboard')} — AliStroy` })

const cards = computed(() => [
  { label: t('admin.products'), value: totals.value?.total_products ?? 0, link: '/admin/products' },
  { label: t('seller_panel.stats_approved'), value: totals.value?.total_approved_products ?? 0 },
  { label: t('seller_panel.stats_pending'), value: totals.value?.total_pending_products ?? 0, link: '/admin/products?status=pending' },
  { label: t('admin.sellers'), value: totals.value?.total_sellers ?? 0, link: '/admin/sellers' },
  { label: t('admin.users'), value: totals.value?.total_users ?? 0, link: '/admin/users' },
  { label: t('seller_panel.stats_views'), value: totals.value?.total_views ?? 0 },
  { label: t('seller_panel.stats_phone'), value: totals.value?.total_phone_clicks ?? 0 },
  { label: t('seller_panel.stats_whatsapp'), value: totals.value?.total_whatsapp_clicks ?? 0 },
  { label: t('admin.reviews'), value: totals.value?.total_reviews ?? 0, link: '/admin/reviews' },
])
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('admin.dashboard') }}</h1>
    <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <component
        :is="c.link ? resolveComponent('NuxtLink') : 'div'"
        v-for="c in cards"
        :key="c.label"
        :to="c.link"
        class="card p-5 hover:shadow-card transition-shadow block"
      >
        <div class="text-sm text-gray-500">{{ c.label }}</div>
        <div class="mt-1 text-3xl font-extrabold text-ink-900">{{ c.value }}</div>
      </component>
    </div>
  </div>
</template>
