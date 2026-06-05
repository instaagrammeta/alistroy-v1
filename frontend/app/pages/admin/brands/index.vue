<script setup lang="ts">
import type { Brand } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { resolve } = useImageUrl()

const { data, refresh } = await useAsyncData('admin-brands', () => useApi()<{ data: Brand[] }>('/admin/brands?all=true'))
const brands = computed<Brand[]>(() => (data.value as any)?.data || [])

const editing = ref(false)
const form = reactive<any>(blank())
function blank() { return { id: '', name: '', slug: '', logo_url: '', sort_order: 0, active: true } }
const newBrand = () => { Object.assign(form, blank()); editing.value = true }
const editBrand = (b: Brand) => { Object.assign(form, b); editing.value = true }
const logoArr = computed({ get: () => (form.logo_url ? [{ url: form.logo_url }] : []), set: (v: any[]) => (form.logo_url = v[0]?.url || '') })
const save = async () => {
  const body = { ...form }; delete body.id
  if (form.id) await useApi()(`/admin/brands/${form.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/brands', { method: 'POST', body })
  editing.value = false; await refresh()
}
const remove = async (b: Brand) => {
  if (!confirm(t('common.confirm_delete'))) return
  try { await useApi()(`/admin/brands/${b.id}`, { method: 'DELETE' }); await refresh() }
  catch (e: any) { alert(e?.data?.error?.message || 'Error') }
}

useHead({ title: () => `${t('admin.brands')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.brands') }}</h1>
      <button class="btn-primary btn-sm" @click="newBrand">+ {{ t('common.create') }}</button>
    </div>
    <div v-if="editing" class="card p-5 mb-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div><label class="label">{{ t('common.name') }}</label><input v-model="form.name" class="input" /></div>
        <div><label class="label">{{ t('admin.sort_order') }}</label><input v-model.number="form.sort_order" type="number" class="input" /></div>
      </div>
      <div><label class="label">Logo (≥500×500)</label><ImageUploader v-model="logoArr" endpoint="/admin/upload" subdir="brands" :multiple="false" /></div>
      <label class="inline-flex items-center gap-2"><input v-model="form.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      <div class="flex gap-2"><button class="btn-primary" @click="save">{{ t('common.save') }}</button><button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button></div>
    </div>
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
      <div v-for="b in brands" :key="b.id" class="card p-4 text-center">
        <div class="aspect-square rounded-lg bg-gray-50 flex items-center justify-center overflow-hidden mb-2"><img v-if="b.logo_url" :src="resolve(b.logo_url)" alt="" class="w-full h-full object-contain" /><span v-else class="text-2xl font-bold text-brand-600">{{ b.name.charAt(0) }}</span></div>
        <div class="font-medium text-sm truncate">{{ b.name }}</div>
        <div class="mt-2 flex justify-center gap-2 text-sm"><button class="text-brand-600" @click="editBrand(b)">{{ t('common.edit') }}</button><button class="text-red-600" @click="remove(b)">{{ t('common.delete') }}</button></div>
      </div>
    </div>
  </div>
</template>
