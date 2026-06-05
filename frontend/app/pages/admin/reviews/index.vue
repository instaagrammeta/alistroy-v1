<script setup lang="ts">
import type { PaginatedResponse, Review } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const { t } = useI18n()
const { productTitle } = useLocaleField()
const { formatDate } = useFormatters()

const status = ref('')
const page = ref(1)

const { data, refresh } = await useAsyncData(
  'admin-reviews',
  () =>
    api<PaginatedResponse<Review>>('/admin/reviews', {
      query: { page: page.value, page_size: 20, status: status.value || undefined },
    }),
  { watch: [status, page] }
)

const items = computed<Review[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const moderate = async (r: Review, decision: 'approved' | 'rejected') => {
  await api(`/admin/reviews/${r.id}/moderate`, { method: 'POST', body: { status: decision } })
  await refresh()
}

const onDelete = async (r: Review) => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  await api(`/admin/reviews/${r.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('admin.reviews')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.reviews') }}</h1>

    <div class="card p-4 mb-4 flex items-end gap-3">
      <div class="w-48">
        <label class="label">Status</label>
        <select v-model="status" class="select">
          <option value="">{{ t('common.all') }}</option>
          <option value="pending">pending</option>
          <option value="approved">approved</option>
          <option value="rejected">rejected</option>
        </select>
      </div>
    </div>

    <div class="space-y-3">
      <div v-for="r in items" :key="r.id" class="card p-4">
        <div class="flex items-center justify-between">
          <div class="font-medium">{{ r.user?.name || '—' }} <span class="text-gray-400 text-xs">{{ r.user?.email }}</span></div>
          <div class="text-amber-400">{{ '★'.repeat(r.rating) }}<span class="text-gray-200">{{ '★'.repeat(5 - r.rating) }}</span></div>
        </div>
        <div class="text-xs text-gray-400 mt-1">{{ formatDate(r.created_at) }} · {{ r.status }}</div>
        <div v-if="r.product" class="text-sm mt-1 text-gray-500">{{ productTitle(r.product) }}</div>
        <p v-if="r.comment" class="mt-2 text-gray-700">{{ r.comment }}</p>
        <div class="mt-3 flex gap-2">
          <button class="btn-success" @click="moderate(r, 'approved')">{{ t('admin.approve') }}</button>
          <button class="btn-outline" @click="moderate(r, 'rejected')">{{ t('admin.reject') }}</button>
          <button class="btn-danger ml-auto" @click="onDelete(r)">{{ t('common.delete') }}</button>
        </div>
      </div>
      <EmptyState v-if="!items.length" />
    </div>

    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
