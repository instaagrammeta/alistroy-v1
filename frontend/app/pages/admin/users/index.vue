<script setup lang="ts">
import type { PaginatedResponse, User } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const { t } = useI18n()

const page = ref(1)
const q = ref('')
const role = ref('')

const { data, refresh } = await useAsyncData(
  'admin-users',
  () =>
    api<PaginatedResponse<User>>('/admin/users', {
      query: { page: page.value, page_size: 20, q: q.value || undefined, role: role.value || undefined },
    }),
  { watch: [page, q, role] }
)

const items = computed<User[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const onDelete = async (u: User) => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  await api(`/admin/users/${u.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('admin.users')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.users') }}</h1>

    <div class="card p-4 mb-4 flex flex-wrap items-end gap-3">
      <div class="flex-1 min-w-[200px]">
        <label class="label">{{ t('common.search') }}</label>
        <input v-model="q" type="search" class="input" />
      </div>
      <div class="w-48">
        <label class="label">Role</label>
        <select v-model="role" class="select">
          <option value="">{{ t('common.all') }}</option>
          <option value="customer">customer</option>
          <option value="seller">seller</option>
          <option value="admin">admin</option>
        </select>
      </div>
    </div>

    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr>
            <th class="p-3">{{ t('auth.email') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('auth.name') }}</th>
            <th class="p-3 hidden md:table-cell">{{ t('auth.phone') }}</th>
            <th class="p-3">Role</th>
            <th class="p-3">Active</th>
            <th class="p-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in items" :key="u.id" class="border-t border-gray-100">
            <td class="p-3 font-medium">{{ u.email }}</td>
            <td class="p-3 hidden md:table-cell">{{ u.name }}</td>
            <td class="p-3 hidden md:table-cell">{{ u.phone || '—' }}</td>
            <td class="p-3"><span class="badge bg-gray-100 text-gray-800">{{ u.role }}</span></td>
            <td class="p-3">{{ u.is_active ? '✓' : '✗' }}</td>
            <td class="p-3 text-right whitespace-nowrap">
              <NuxtLink :to="`/admin/users/${u.id}`" class="text-brand-600 hover:text-brand-700 mr-3">{{ t('common.edit') }}</NuxtLink>
              <button class="text-red-600 hover:text-red-700" @click="onDelete(u)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="6" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
