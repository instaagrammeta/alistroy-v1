<script setup lang="ts">
import type { FinancialSummary } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { formatPrice } = useFormatters()

const preset = ref('month')
const from = ref('')
const to = ref('')

const buildQuery = () => {
  if (preset.value === 'custom') return { from: from.value || undefined, to: to.value || undefined }
  return { preset: preset.value }
}

const { data, refresh } = await useAsyncData('admin-report', () => useApi()<{ data: FinancialSummary }>('/admin/reports/summary', { query: buildQuery() }))
const summary = computed<FinancialSummary | null>(() => (data.value as any)?.data || null)

watch([preset, from, to], () => refresh())

const doExport = async () => {
  const res = await useApi().raw('/admin/export/report', { query: buildQuery(), responseType: 'blob' })
  downloadBlob(res._data as Blob, 'report.xlsx')
}

const presets = [
  { key: 'today', label: t('admin.preset_today') },
  { key: 'week', label: t('admin.preset_week') },
  { key: 'month', label: t('admin.preset_month') },
  { key: 'custom', label: t('admin.preset_custom') },
]

useHead({ title: () => `${t('admin.reports')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-2 flex-wrap">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.reports') }}</h1>
      <button class="btn-outline btn-sm" @click="doExport">{{ t('admin.export') }}</button>
    </div>

    <div class="card p-4 mb-6 flex flex-wrap items-end gap-3">
      <div class="flex gap-2">
        <button v-for="p in presets" :key="p.key" class="badge cursor-pointer" :class="preset === p.key ? 'bg-brand-500 text-white' : 'bg-gray-100 text-gray-700'" @click="preset = p.key">{{ p.label }}</button>
      </div>
      <template v-if="preset === 'custom'">
        <div><label class="label">{{ t('common.from') }}</label><input v-model="from" type="date" class="input" /></div>
        <div><label class="label">{{ t('common.to') }}</label><input v-model="to" type="date" class="input" /></div>
      </template>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.income') }}</div><div class="mt-1 text-2xl font-extrabold text-emerald-600">{{ formatPrice(summary?.income ?? 0) }}</div></div>
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.purchase') }}</div><div class="mt-1 text-2xl font-extrabold text-amber-600">{{ formatPrice(summary?.purchase ?? 0) }}</div></div>
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.expense') }}</div><div class="mt-1 text-2xl font-extrabold text-red-600">{{ formatPrice(summary?.expense ?? 0) }}</div></div>
      <div class="card p-5"><div class="text-sm text-gray-500">{{ t('admin.profit') }}</div><div class="mt-1 text-2xl font-extrabold text-brand-600">{{ formatPrice(summary?.profit ?? 0) }}</div></div>
    </div>
  </div>
</template>
