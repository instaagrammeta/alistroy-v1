<script setup lang="ts">
import type { PaginatedResponse, Seller } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const { t } = useI18n()
const { resolve } = useImageUrl()

const page = ref(1)
const q = ref('')
const status = ref('')

const { data, refresh } = await useAsyncData(
  'admin-sellers',
  () =>
    api<PaginatedResponse<Seller>>('/admin/sellers', {
      query: { page: page.value, page_size: 20, q: q.value || undefined, status: status.value || undefined },
    }),
  { watch: [page, q, status] }
)

const items = computed<Seller[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const onDelete = async (s: Seller) => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  await api(`/admin/sellers/${s.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('admin.sellers')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.sellers') }}</h1>

    <div class="card p-4 mb-4 flex flex-wrap items-end gap-3">
      <div class="flex-1 min-w-[200px]">
        <label class="label">{{ t('common.search') }}</label>
        <input v-model="q" class="input" type="search" />
      </div>
      <div class="w-48">
        <label class="label">Status</label>
        <select v-model="status" class="select">
          <option value="">{{ t('common.all') }}</option>
          <option value="pending">{{ t('seller_panel.status_pending') }}</option>
          <option value="approved">{{ t('seller_panel.status_approved') }}</option>
          <option value="blocked">Blocked</option>
        </select>
      </div>
    </div>

    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr>
            <th class="p-3">{{ t('catalog.seller') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('seller.city') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('admin.phone_number') }}</th>
            <th class="p-3">Status</th>
            <th class="p-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in items" :key="s.id" class="border-t border-gray-100">
            <td class="p-3">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded bg-gray-50 overflow-hidden flex items-center justify-center">
                  <img v-if="s.logo_url" :src="resolve(s.logo_url)" alt="" class="w-full h-full object-cover" />
                  <span v-else class="text-brand-600 font-bold">{{ s.name.charAt(0).toUpperCase() }}</span>
                </div>
                <div>
                  <div class="font-medium">{{ s.name }}</div>
                  <div class="text-xs text-gray-400">{{ s.slug }}</div>
                </div>
              </div>
            </td>
            <td class="p-3 hidden md:table-cell">{{ s.city || '—' }}</td>
            <td class="p-3 hidden md:table-cell">{{ s.phone || '—' }}</td>
            <td class="p-3"><span class="badge bg-gray-100 text-gray-800">{{ s.status }}</span></td>
            <td class="p-3 text-right whitespace-nowrap">
              <NuxtLink :to="`/admin/sellers/${s.id}`" class="text-brand-600 hover:text-brand-700 mr-3">{{ t('common.edit') }}</NuxtLink>
              <button class="text-red-600 hover:text-red-700" @click="onDelete(s)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
