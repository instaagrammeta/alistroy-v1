<script setup lang="ts">
import type { PaginatedResponse, Seller } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { resolve } = useImageUrl()

const page = ref(1)
const q = ref('')
const { data, refresh } = await useAsyncData('admin-sellers', () =>
  useApi()<PaginatedResponse<Seller>>('/admin/sellers', { query: { page: page.value, page_size: 20, q: q.value || undefined } }),
  { watch: [page, q] }
)
const sellers = computed<Seller[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const editing = ref(false)
const form = reactive<any>(blank())
function blank() { return { id: '', full_name: '', company_name: '', market_name: '', phone: '', phone_alt: '', whatsapp: '', telegram: '', telegram_username: '', address: '', city: '', notes: '', logo_url: '', login: '', password: '', active: true } }
const newSeller = () => { Object.assign(form, blank()); editing.value = true }
const editSeller = (s: Seller) => { Object.assign(form, s, { password: '' }); editing.value = true }
const logoArr = computed({ get: () => (form.logo_url ? [{ url: form.logo_url }] : []), set: (v: any[]) => (form.logo_url = v[0]?.url || '') })
const save = async () => {
  const body = { ...form }; delete body.id
  if (form.id) await useApi()(`/admin/sellers/${form.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/sellers', { method: 'POST', body })
  editing.value = false; await refresh()
}
const remove = async (s: Seller) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/sellers/${s.id}`, { method: 'DELETE' }); await refresh()
}
const doExport = async () => {
  const res = await useApi().raw('/admin/export/sellers', { responseType: 'blob' })
  downloadBlob(res._data as Blob, 'sellers.xlsx')
}
useHead({ title: () => `${t('admin.sellers')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-2">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.sellers') }}</h1>
      <div class="flex gap-2"><button class="btn-outline btn-sm" @click="doExport">{{ t('admin.export') }}</button><button class="btn-primary btn-sm" @click="newSeller">+ {{ t('admin.new_seller') }}</button></div>
    </div>

    <div v-if="editing" class="card p-5 mb-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div><label class="label">{{ t('admin.full_name') }}</label><input v-model="form.full_name" class="input" /></div>
        <div><label class="label">{{ t('admin.company') }}</label><input v-model="form.company_name" class="input" /></div>
        <div><label class="label">{{ t('admin.market') }}</label><input v-model="form.market_name" class="input" /></div>
        <div><label class="label">{{ t('admin.phone') }}</label><input v-model="form.phone" class="input" /></div>
        <div><label class="label">{{ t('admin.phone') }} 2</label><input v-model="form.phone_alt" class="input" /></div>
        <div><label class="label">WhatsApp</label><input v-model="form.whatsapp" class="input" /></div>
        <div><label class="label">Telegram</label><input v-model="form.telegram" class="input" /></div>
        <div><label class="label">{{ t('seller.city') }}</label><input v-model="form.city" class="input" /></div>
        <div><label class="label">{{ t('seller.address') }}</label><input v-model="form.address" class="input" /></div>
        <div><label class="label">{{ t('admin.login') }}</label><input v-model="form.login" class="input" /></div>
        <div><label class="label">{{ t('admin.password') }}</label><input v-model="form.password" type="text" class="input" :placeholder="form.id ? '••••' : ''" /></div>
        <label class="inline-flex items-center gap-2 mt-7"><input v-model="form.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      </div>
      <div><label class="label">Logo</label><ImageUploader v-model="logoArr" endpoint="/admin/upload" subdir="logos" :multiple="false" /></div>
      <div class="flex gap-2"><button class="btn-primary" @click="save">{{ t('common.save') }}</button><button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button></div>
    </div>

    <div class="card p-4 mb-4"><input v-model="q" type="search" class="input max-w-md" :placeholder="t('common.search')" /></div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500"><tr><th class="p-3">{{ t('catalog.seller') }}</th><th class="p-3 hidden md:table-cell">{{ t('seller.city') }}</th><th class="p-3 hidden md:table-cell">{{ t('admin.phone') }}</th><th class="p-3">{{ t('common.status') }}</th><th class="p-3"></th></tr></thead>
        <tbody>
          <tr v-for="s in sellers" :key="s.id" class="border-t border-gray-100">
            <td class="p-3"><div class="flex items-center gap-3"><div class="w-10 h-10 rounded bg-gray-50 overflow-hidden flex items-center justify-center"><img v-if="s.logo_url" :src="resolve(s.logo_url)" alt="" class="w-full h-full object-cover" /><span v-else class="text-brand-600 font-bold">{{ (s.market_name || s.full_name).charAt(0) }}</span></div><div><div class="font-medium">{{ s.market_name || s.full_name }}</div><div class="text-xs text-gray-400">{{ s.slug }}</div></div></div></td>
            <td class="p-3 hidden md:table-cell">{{ s.city || '—' }}</td>
            <td class="p-3 hidden md:table-cell">{{ s.phone || '—' }}</td>
            <td class="p-3"><span class="badge" :class="s.active ? 'badge-approved' : 'badge-draft'">{{ s.active ? t('common.active') : t('common.inactive') }}</span></td>
            <td class="p-3 text-right whitespace-nowrap"><button class="text-brand-600 mr-3" @click="editSeller(s)">{{ t('common.edit') }}</button><button class="text-red-600" @click="remove(s)">{{ t('common.delete') }}</button></td>
          </tr>
          <tr v-if="!sellers.length"><td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
