<script setup lang="ts">
import type { Category } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const api = useApi()
const { t } = useI18n()
const { categoryTitle } = useLocaleField()

const { data, refresh } = await useAsyncData('admin-cats', () =>
  api<{ data: Category[] }>('/admin/categories')
)
const items = computed<Category[]>(() => (data.value as any)?.data || [])

const form = reactive({
  id: '' as string | null,
  title_tj: '',
  title_ru: '',
  slug: '',
  sort_order: 0,
  is_active: true,
})
const editing = ref(false)
const message = ref('')
const error = ref('')

const startNew = () => {
  editing.value = true
  form.id = null
  form.title_tj = ''
  form.title_ru = ''
  form.slug = ''
  form.sort_order = 0
  form.is_active = true
}

const startEdit = (c: Category) => {
  editing.value = true
  form.id = c.id
  form.title_tj = c.title_tj
  form.title_ru = c.title_ru
  form.slug = c.slug
  form.sort_order = c.sort_order
  form.is_active = c.is_active
}

const save = async () => {
  message.value = ''
  error.value = ''
  try {
    const body = {
      title_tj: form.title_tj,
      title_ru: form.title_ru,
      slug: form.slug,
      sort_order: form.sort_order,
      is_active: form.is_active,
    }
    if (form.id) {
      await api(`/admin/categories/${form.id}`, { method: 'PATCH', body })
    } else {
      await api('/admin/categories', { method: 'POST', body })
    }
    message.value = t('common.save')
    editing.value = false
    await refresh()
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
  }
}

const onDelete = async (c: Category) => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  try {
    await api(`/admin/categories/${c.id}`, { method: 'DELETE' })
    await refresh()
  } catch (e: any) {
    alert(e?.data?.error?.message || 'Delete failed')
  }
}

useHead({ title: () => `${t('admin.categories')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.categories') }}</h1>
      <button class="btn-primary" @click="startNew">+ New</button>
    </div>

    <div v-if="editing" class="card p-5 mb-6">
      <h2 class="font-semibold mb-3">{{ form.id ? t('common.edit') : '+ New' }}</h2>
      <form class="grid grid-cols-1 md:grid-cols-2 gap-3" @submit.prevent="save">
        <div>
          <label class="label">{{ t('seller_panel.title_tj') }}</label>
          <input v-model="form.title_tj" class="input" required />
        </div>
        <div>
          <label class="label">{{ t('seller_panel.title_ru') }}</label>
          <input v-model="form.title_ru" class="input" required />
        </div>
        <div>
          <label class="label">Slug</label>
          <input v-model="form.slug" class="input" />
        </div>
        <div>
          <label class="label">Order</label>
          <input v-model.number="form.sort_order" type="number" class="input" />
        </div>
        <label class="inline-flex items-center gap-2 col-span-2">
          <input v-model="form.is_active" type="checkbox" class="h-4 w-4" /> Active
        </label>
        <div v-if="error" class="col-span-2 text-sm text-red-600">{{ error }}</div>
        <div v-if="message" class="col-span-2 text-sm text-emerald-600">{{ message }}</div>
        <div class="col-span-2 flex gap-2">
          <button type="submit" class="btn-primary">{{ t('common.save') }}</button>
          <button type="button" class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button>
        </div>
      </form>
    </div>

    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 text-left text-xs uppercase text-gray-500">
          <tr>
            <th class="p-3">{{ t('catalog.title') }}</th>
            <th class="p-3">Slug</th>
            <th class="p-3">Order</th>
            <th class="p-3">Active</th>
            <th class="p-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in items" :key="c.id" class="border-t border-gray-100">
            <td class="p-3 font-medium">{{ categoryTitle(c) }}</td>
            <td class="p-3 text-gray-500">{{ c.slug }}</td>
            <td class="p-3">{{ c.sort_order }}</td>
            <td class="p-3">{{ c.is_active ? '✓' : '✗' }}</td>
            <td class="p-3 text-right whitespace-nowrap">
              <button class="text-brand-600 hover:text-brand-700 mr-3" @click="startEdit(c)">{{ t('common.edit') }}</button>
              <button class="text-red-600 hover:text-red-700" @click="onDelete(c)">{{ t('common.delete') }}</button>
            </td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="5" class="p-8 text-center text-gray-400">{{ t('common.no_results') }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
