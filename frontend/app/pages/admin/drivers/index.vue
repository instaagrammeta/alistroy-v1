<script setup lang="ts">
import type { Driver, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { resolve } = useImageUrl()

const page = ref(1)
const { data, refresh } = await useAsyncData('admin-drivers', () =>
  useApi()<PaginatedResponse<Driver>>('/admin/drivers', { query: { page: page.value, page_size: 20 } }),
  { watch: [page] }
)
const drivers = computed<Driver[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: 20, total: 0, total_pages: 1 })

const editing = ref(false)
const form = reactive<any>(blank())
function blank() { return { id: '', full_name: '', age: 0, phone: '', phone_alt: '', whatsapp: '', telegram: '', vehicle: '', photo_url: '', notes: '', login: '', password: '', active: true, on_duty: true } }
const newDriver = () => { Object.assign(form, blank()); editing.value = true }
const editDriver = (d: Driver) => { Object.assign(form, d, { password: '' }); editing.value = true }
const photoArr = computed({ get: () => (form.photo_url ? [{ url: form.photo_url }] : []), set: (v: any[]) => (form.photo_url = v[0]?.url || '') })
const save = async () => {
  const body = { ...form }; delete body.id
  if (form.id) await useApi()(`/admin/drivers/${form.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/drivers', { method: 'POST', body })
  editing.value = false; await refresh()
}
const remove = async (d: Driver) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/drivers/${d.id}`, { method: 'DELETE' }); await refresh()
}
useHead({ title: () => `${t('admin.drivers')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.drivers') }}</h1>
      <button class="btn-primary btn-sm" @click="newDriver">+ {{ t('admin.new_driver') }}</button>
    </div>
    <div v-if="editing" class="card p-5 mb-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div><label class="label">{{ t('admin.full_name') }}</label><input v-model="form.full_name" class="input" /></div>
        <div><label class="label">{{ t('admin.age') }}</label><input v-model.number="form.age" type="number" class="input" /></div>
        <div><label class="label">{{ t('driver_panel.vehicle') }}</label><input v-model="form.vehicle" class="input" /></div>
        <div><label class="label">{{ t('admin.phone') }}</label><input v-model="form.phone" class="input" /></div>
        <div><label class="label">{{ t('admin.phone') }} 2</label><input v-model="form.phone_alt" class="input" /></div>
        <div><label class="label">WhatsApp</label><input v-model="form.whatsapp" class="input" /></div>
        <div><label class="label">Telegram</label><input v-model="form.telegram" class="input" /></div>
        <div><label class="label">{{ t('admin.login') }}</label><input v-model="form.login" class="input" /></div>
        <div><label class="label">{{ t('admin.password') }}</label><input v-model="form.password" class="input" /></div>
        <label class="inline-flex items-center gap-2 mt-7"><input v-model="form.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      </div>
      <div><label class="label">{{ t('common.name') }} (photo)</label><ImageUploader v-model="photoArr" endpoint="/admin/upload" subdir="drivers" :multiple="false" /></div>
      <div class="flex gap-2"><button class="btn-primary" @click="save">{{ t('common.save') }}</button><button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button></div>
    </div>
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500"><tr><th class="p-3">{{ t('admin.full_name') }}</th><th class="p-3">{{ t('admin.phone') }}</th><th class="p-3 hidden md:table-cell">{{ t('driver_panel.vehicle') }}</th><th class="p-3">{{ t('common.status') }}</th><th class="p-3"></th></tr></thead>
        <tbody>
          <tr v-for="d in drivers" :key="d.id" class="border-t border-gray-100">
            <td class="p-3"><div class="flex items-center gap-3"><div class="w-10 h-10 rounded-full bg-gray-50 overflow-hidden"><img v-if="d.photo_url" :src="resolve(d.photo_url)" alt="" class="w-full h-full object-cover" /></div><span class="font-medium">{{ d.full_name }}</span></div></td>
            <td class="p-3">{{ d.phone }}</td>
            <td class="p-3 hidden md:table-cell">{{ d.vehicle || '—' }}</td>
            <td class="p-3"><span class="badge" :class="d.active ? 'badge-approved' : 'badge-draft'">{{ d.active ? t('common.active') : t('common.inactive') }}</span></td>
            <td class="p-3 text-right whitespace-nowrap"><button class="text-brand-600 mr-3" @click="editDriver(d)">{{ t('common.edit') }}</button><button class="text-red-600" @click="remove(d)">{{ t('common.delete') }}</button></td>
          </tr>
          <tr v-if="!drivers.length"><td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td></tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
