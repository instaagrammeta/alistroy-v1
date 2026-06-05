<script setup lang="ts">
import type { Category, Subcategory } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { categoryName } = useLocaleField()
const { resolve } = useImageUrl()

const { data, refresh } = await useAsyncData('admin-cats', () => useApi()<{ data: Category[] }>('/admin/categories'))
const cats = computed<Category[]>(() => (data.value as any)?.data || [])

// category form
const catForm = reactive<any>(blankCat())
const catEditing = ref(false)
function blankCat() { return { id: '', name_tj: '', name_ru: '', slug: '', icon_url: '', sort_order: 0, active: true } }
const newCat = () => { Object.assign(catForm, blankCat()); catEditing.value = true }
const editCat = (c: Category) => { Object.assign(catForm, c); catEditing.value = true }
const iconArr = computed({ get: () => (catForm.icon_url ? [{ url: catForm.icon_url }] : []), set: (v: any[]) => (catForm.icon_url = v[0]?.url || '') })
const saveCat = async () => {
  const body = { ...catForm }; delete body.id; delete body.subcategories
  if (catForm.id) await useApi()(`/admin/categories/${catForm.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/categories', { method: 'POST', body })
  catEditing.value = false; await refresh()
}
const delCat = async (c: Category) => {
  if (!confirm(t('common.confirm_delete'))) return
  try { await useApi()(`/admin/categories/${c.id}`, { method: 'DELETE' }); await refresh() }
  catch (e: any) { alert(e?.data?.error?.message || 'Error') }
}

// subcategory form
const subForm = reactive<any>(blankSub())
const subEditing = ref(false)
function blankSub() { return { id: '', category_id: '', name_tj: '', name_ru: '', slug: '', icon_url: '', sort_order: 0, active: true } }
const newSub = (catId: string) => { Object.assign(subForm, blankSub()); subForm.category_id = catId; subEditing.value = true }
const editSub = (s: Subcategory) => { Object.assign(subForm, s); subEditing.value = true }
const saveSub = async () => {
  const body = { ...subForm }; delete body.id
  if (subForm.id) await useApi()(`/admin/subcategories/${subForm.id}`, { method: 'PATCH', body })
  else await useApi()('/admin/subcategories', { method: 'POST', body })
  subEditing.value = false; await refresh()
}
const delSub = async (s: Subcategory) => {
  if (!confirm(t('common.confirm_delete'))) return
  try { await useApi()(`/admin/subcategories/${s.id}`, { method: 'DELETE' }); await refresh() }
  catch (e: any) { alert(e?.data?.error?.message || 'Error') }
}

useHead({ title: () => `${t('admin.categories')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.categories') }}</h1>
      <button class="btn-primary btn-sm" @click="newCat">+ {{ t('admin.new_category') }}</button>
    </div>

    <div v-if="catEditing" class="card p-5 mb-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div><label class="label">{{ t('seller_panel.name_tj') }}</label><input v-model="catForm.name_tj" class="input" /></div>
        <div><label class="label">{{ t('seller_panel.name_ru') }}</label><input v-model="catForm.name_ru" class="input" /></div>
        <div><label class="label">Slug</label><input v-model="catForm.slug" class="input" /></div>
        <div><label class="label">{{ t('admin.sort_order') }}</label><input v-model.number="catForm.sort_order" type="number" class="input" /></div>
      </div>
      <div><label class="label">Icon (500×500)</label><ImageUploader v-model="iconArr" endpoint="/admin/upload" subdir="categories" :multiple="false" /></div>
      <label class="inline-flex items-center gap-2"><input v-model="catForm.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      <div class="flex gap-2"><button class="btn-primary" @click="saveCat">{{ t('common.save') }}</button><button class="btn-outline" @click="catEditing = false">{{ t('common.cancel') }}</button></div>
    </div>

    <div v-if="subEditing" class="card p-5 mb-4 space-y-3 border-l-4 border-brand-300">
      <h3 class="font-semibold">{{ t('admin.subcategories') }}</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div><label class="label">{{ t('seller_panel.name_tj') }}</label><input v-model="subForm.name_tj" class="input" /></div>
        <div><label class="label">{{ t('seller_panel.name_ru') }}</label><input v-model="subForm.name_ru" class="input" /></div>
      </div>
      <label class="inline-flex items-center gap-2"><input v-model="subForm.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}</label>
      <div class="flex gap-2"><button class="btn-primary" @click="saveSub">{{ t('common.save') }}</button><button class="btn-outline" @click="subEditing = false">{{ t('common.cancel') }}</button></div>
    </div>

    <div class="space-y-3">
      <div v-for="c in cats" :key="c.id" class="card p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-brand-50 flex items-center justify-center overflow-hidden"><img v-if="c.icon_url" :src="resolve(c.icon_url)" alt="" class="w-7 h-7 object-contain" /><span v-else class="text-brand-600 font-bold">{{ categoryName(c).charAt(0) }}</span></div>
          <div class="flex-1 min-w-0">
            <div class="font-semibold text-ink-900">{{ categoryName(c) }} <span class="text-xs text-gray-400">/{{ c.slug }}</span></div>
          </div>
          <span class="badge" :class="c.active ? 'badge-approved' : 'badge-draft'">{{ c.active ? t('common.active') : t('common.inactive') }}</span>
          <button class="text-sky-600 text-sm" @click="newSub(c.id)">+ {{ t('admin.subcategories') }}</button>
          <button class="text-brand-600 text-sm" @click="editCat(c)">{{ t('common.edit') }}</button>
          <button class="text-red-600 text-sm" @click="delCat(c)">{{ t('common.delete') }}</button>
        </div>
        <div v-if="c.subcategories && c.subcategories.length" class="mt-3 flex flex-wrap gap-2 pl-13">
          <span v-for="s in c.subcategories" :key="s.id" class="badge bg-gray-100 text-gray-700 flex items-center gap-1">
            {{ categoryName(s) }}
            <button class="text-brand-600" @click="editSub(s)">✎</button>
            <button class="text-red-500" @click="delSub(s)">✕</button>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
