<script setup lang="ts">
import type { Customer, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { formatDate } = useFormatters()

const page = ref(1)
const q = ref('')
const { data, refresh } = await useAsyncData('admin-customers', () =>
  useApi()<PaginatedResponse<Customer>>('/admin/customers', { query: { page: page.value, page_size: 20, q: q.value || undefined } }),
  { watch: [page, q] }
)
const customers = computed<Customer[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const editing = ref(false)
const form = reactive<any>(blank())
function blank() { return { id: '', name: '', company: '', status: 'active', phone: '', phone_alt: '', address: '', city: '', notes: '', password: '' } }
const newCustomer = () => { Object.assign(form, blank()); editing.value = true }
const save = async () => {
  const body = { ...form }; delete body.id
  if (form.id) await useApi()(`/admin/customers/${form.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/customers', { method: 'POST', body })
  editing.value = false; await refresh()
}
const remove = async (c: Customer) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/customers/${c.id}`, { method: 'DELETE' }); await refresh()
}
const doExport = async () => {
  const res = await useApi().raw('/admin/export/customers', { responseType: 'blob' })
  downloadBlob(res._data as Blob, 'customers.xlsx')
}
useHead({ title: () => `${t('admin.customers')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-2">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.customers') }}</h1>
      <div class="flex gap-2"><button class="btn-outline btn-sm" @click="doExport">{{ t('admin.export') }}</button><button class="btn-primary btn-sm" @click="newCustomer">+ {{ t('admin.new_customer') }}</button></div>
    </div>
    <div v-if="editing" class="card p-5 mb-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div><label class="label">{{ t('common.name') }}</label><input v-model="form.name" class="input" /></div>
        <div><label class="label">{{ t('admin.company') }}</label><input v-model="form.company" class="input" /></div>
        <div><label class="label">{{ t('common.status') }}</label><select v-model="form.status" class="select"><option value="active">active</option><option value="inactive">inactive</option><option value="blocked">blocked</option></select></div>
        <div><label class="label">{{ t('admin.phone') }}</label><input v-model="form.phone" class="input" /></div>
        <div><label class="label">{{ t('admin.phone') }} 2</label><input v-model="form.phone_alt" class="input" /></div>
        <div><label class="label">{{ t('seller.city') }}</label><input v-model="form.city" class="input" /></div>
        <div class="md:col-span-2"><label class="label">{{ t('auth.address') }}</label><input v-model="form.address" class="input" /></div>
        <div><label class="label">{{ t('admin.password') }}</label><input v-model="form.password" class="input" /></div>
      </div>
      <div class="flex gap-2"><button class="btn-primary" @click="save">{{ t('common.save') }}</button><button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button></div>
    </div>
    <div class="card p-4 mb-4"><input v-model="q" type="search" class="input max-w-md" :placeholder="t('common.search')" /></div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500"><tr><th class="p-3">{{ t('admin.company') }}</th><th class="p-3">{{ t('seller.city') }}</th><th class="p-3 hidden md:table-cell">{{ t('common.date') }}</th><th class="p-3"></th></tr></thead>
        <tbody>
          <tr v-for="c in customers" :key="c.id" class="border-t border-gray-100">
            <td class="p-3 font-medium">{{ c.company || '—' }}</td>
            <td class="p-3">{{ c.city || '—' }}</td>
            <td class="p-3 hidden md:table-cell">{{ formatDate((c as any).created_at) }}</td>
            <td class="p-3 text-right whitespace-nowrap"><button class="text-red-600" @click="remove(c)">{{ t('common.delete') }}</button></td>
          </tr>
          <tr v-if="!customers.length"><td colspan="4" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
