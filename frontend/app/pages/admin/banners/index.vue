<script setup lang="ts">
import type { Banner } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { resolve } = useImageUrl()

const positions = ['hero', 'category_row', 'mid_large', 'mid_small', 'side', 'footer']

const { data, refresh } = await useAsyncData('admin-banners', () => useApi()<{ data: Banner[] }>('/admin/banners'))
const banners = computed<Banner[]>(() => (data.value as any)?.data || [])

const editing = ref(false)
const form = reactive<any>(blank())
function blank() {
  return { id: '', position: 'hero', title_tj: '', title_ru: '', description_tj: '', description_ru: '', desktop_url: '', tablet_url: '', mobile_url: '', link_url: '', sort_order: 0, active: true }
}
const startNew = () => { Object.assign(form, blank()); editing.value = true }
const startEdit = (b: Banner) => { Object.assign(form, b); editing.value = true }

const arr = (key: 'desktop_url' | 'tablet_url' | 'mobile_url') => computed({
  get: () => (form[key] ? [{ url: form[key] }] : []),
  set: (v: { url: string }[]) => (form[key] = v[0]?.url || ''),
})
const desktop = arr('desktop_url'); const tablet = arr('tablet_url'); const mobile = arr('mobile_url')

const save = async () => {
  const body = { ...form }
  delete body.id
  if (form.id) await useApi()(`/admin/banners/${form.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/banners', { method: 'POST', body })
  editing.value = false
  await refresh()
}
const remove = async (b: Banner) => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/banners/${b.id}`, { method: 'DELETE' })
  await refresh()
}

useHead({ title: () => `${t('admin.banners')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.banners') }}</h1>
      <button class="btn-primary btn-sm" @click="startNew">+ {{ t('admin.new_banner') }}</button>
    </div>

    <div v-if="editing" class="card p-5 mb-6 space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="label">{{ t('admin.position') }}</label>
          <select v-model="form.position" class="select"><option v-for="p in positions" :key="p" :value="p">{{ p }}</option></select>
        </div>
        <div><label class="label">{{ t('admin.sort_order') }}</label><input v-model.number="form.sort_order" type="number" class="input" /></div>
        <label class="inline-flex items-center gap-2 mt-7"><input v-model="form.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div><label class="label">{{ t('admin.link') }}</label><input v-model="form.link_url" class="input" placeholder="/categories/cement" /></div>
        <div><label class="label">Title (TJ / RU)</label><div class="grid grid-cols-2 gap-2"><input v-model="form.title_tj" class="input" /><input v-model="form.title_ru" class="input" /></div></div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div><label class="label">{{ t('admin.desktop_image') }}</label><ImageUploader v-model="desktop" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
        <div><label class="label">{{ t('admin.tablet_image') }}</label><ImageUploader v-model="tablet" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
        <div><label class="label">{{ t('admin.mobile_image') }}</label><ImageUploader v-model="mobile" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
      </div>
      <div class="flex gap-2"><button class="btn-primary" @click="save">{{ t('common.save') }}</button><button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button></div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="b in banners" :key="b.id" class="card overflow-hidden">
        <div class="aspect-[16/6] bg-gray-100"><img v-if="b.desktop_url" :src="resolve(b.desktop_url)" alt="" class="w-full h-full object-cover" /></div>
        <div class="p-3 flex items-center justify-between">
          <div>
            <div class="text-xs uppercase text-gray-400">{{ b.position }}</div>
            <div class="font-medium text-ink-900">{{ b.title_ru || b.title_tj || b.link_url || '—' }}</div>
          </div>
          <div class="flex items-center gap-2">
            <span class="badge" :class="b.active ? 'badge-approved' : 'badge-draft'">{{ b.active ? t('common.active') : t('common.inactive') }}</span>
            <button class="text-brand-600 text-sm" @click="startEdit(b)">{{ t('common.edit') }}</button>
            <button class="text-red-600 text-sm" @click="remove(b)">{{ t('common.delete') }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
