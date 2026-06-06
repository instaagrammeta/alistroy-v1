<script setup lang="ts">
import type { Banner } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t, locale } = useI18n()
const { resolve } = useImageUrl()

interface PositionDef { key: string; labelTJ: string; labelRU: string; ratio: string; subtitle: string }
const positions: PositionDef[] = [
  { key: 'hero', labelTJ: 'Банери калони боло', labelRU: 'Большой верхний баннер', ratio: 'aspect-[21/9]', subtitle: '1920×820 (PC) / 720×400 (mobile)' },
  { key: 'features', labelTJ: 'Зер-банер бо иконкаҳо', labelRU: 'Под-баннер с иконками', ratio: 'aspect-square', subtitle: '128×128 icon (PC + mobile)' },
  { key: 'category_row', labelTJ: 'Банерҳои категория', labelRU: 'Баннеры категорий', ratio: 'aspect-square', subtitle: '1000×1000 (PC + mobile)' },
  { key: 'mid_large', labelTJ: 'Банери калони поён (рост)', labelRU: 'Большой нижний (справа)', ratio: 'aspect-[16/9]', subtitle: '1280×720 (PC) / 720×420 (mobile)' },
  { key: 'mid_small', labelTJ: 'Банерҳои хурди поён (чап)', labelRU: 'Малые нижние (слева)', ratio: 'aspect-[16/6]', subtitle: '1280×480 (PC) / 720×280 (mobile)' },
  { key: 'side', labelTJ: 'Банерҳои паҳлуӣ', labelRU: 'Боковые баннеры', ratio: 'aspect-[3/4]', subtitle: '480×640 (PC + mobile)' },
  { key: 'footer', labelTJ: 'Банерҳои зерин', labelRU: 'Нижние баннеры', ratio: 'aspect-[16/4]', subtitle: '1280×320 (PC + mobile)' },
]
const activePos = ref(positions[0])

const { data, refresh } = await useAsyncData('admin-banners', () => useApi()<{ data: Banner[] }>('/admin/banners'))
const banners = computed<Banner[]>(() => (data.value as any)?.data || [])
const filtered = computed(() => banners.value.filter((b) => b.position === activePos.value.key))

const editing = ref(false)
const form = reactive<any>(blank())
function blank() {
  return { id: '', position: activePos.value.key, title_tj: '', title_ru: '', description_tj: '', description_ru: '', desktop_url: '', tablet_url: '', mobile_url: '', link_url: '', sort_order: 0, active: true }
}
const startNew = () => { Object.assign(form, blank()); form.position = activePos.value.key; editing.value = true }
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
const toggleActive = async (b: Banner) => {
  await useApi()(`/admin/banners/${b.id}`, { method: 'PATCH', body: { ...b, active: !b.active } })
  await refresh()
}

const posLabel = (p: PositionDef) => (locale.value === 'ru' ? p.labelRU : p.labelTJ)

useHead({ title: () => `${t('admin.banners')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.banners') }}</h1>
      <button class="btn-primary inline-flex items-center gap-1" @click="startNew">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        {{ t('admin.new_banner') }}
      </button>
    </div>

    <!-- Position tabs -->
    <div class="card p-2 mb-5 overflow-x-auto no-scrollbar">
      <div class="flex gap-2 min-w-max">
        <button
          v-for="p in positions" :key="p.key"
          class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap"
          :class="activePos.key === p.key ? 'bg-brand-500 text-white' : 'text-ink-900 hover:bg-gray-50'"
          @click="activePos = p"
        >{{ posLabel(p) }}</button>
      </div>
    </div>
    <div class="text-xs text-gray-400 mb-4 -mt-3">{{ activePos.subtitle }}</div>

    <!-- Editor -->
    <div v-if="editing" class="card p-5 mb-6 space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="label">{{ t('admin.position') }}</label>
          <select v-model="form.position" class="select">
            <option v-for="p in positions" :key="p.key" :value="p.key">{{ posLabel(p) }}</option>
          </select>
        </div>
        <div><label class="label">{{ t('admin.sort_order') }}</label><input v-model.number="form.sort_order" type="number" class="input" /></div>
        <label class="inline-flex items-center gap-2 mt-7">
          <input v-model="form.active" type="checkbox" class="h-4 w-4" /> {{ t('common.active') }}
        </label>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div><label class="label">{{ t('admin.link') }}</label><input v-model="form.link_url" class="input" placeholder="/categories/cement" /></div>
        <div>
          <label class="label">{{ t('common.name') }} (TJ / RU)</label>
          <div class="grid grid-cols-2 gap-2"><input v-model="form.title_tj" class="input" placeholder="TJ" /><input v-model="form.title_ru" class="input" placeholder="RU" /></div>
        </div>
        <div class="md:col-span-2">
          <label class="label">{{ t('product.description') }} (TJ / RU)</label>
          <div class="grid grid-cols-2 gap-2"><textarea v-model="form.description_tj" class="textarea" rows="2" placeholder="TJ"></textarea><textarea v-model="form.description_ru" class="textarea" rows="2" placeholder="RU"></textarea></div>
        </div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div><label class="label">{{ t('admin.desktop_image') }}</label><ImageUploader v-model="desktop" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
        <div><label class="label">{{ t('admin.tablet_image') }}</label><ImageUploader v-model="tablet" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
        <div><label class="label">{{ t('admin.mobile_image') }}</label><ImageUploader v-model="mobile" endpoint="/admin/upload" subdir="banners" :multiple="false" /></div>
      </div>
      <div class="flex gap-2">
        <button class="btn-primary" @click="save">{{ t('common.save') }}</button>
        <button class="btn-outline" @click="editing = false">{{ t('common.cancel') }}</button>
      </div>
    </div>

    <!-- Banner grid for active position -->
    <EmptyState v-if="!filtered.length" />
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="b in filtered" :key="b.id" class="card overflow-hidden">
        <div class="bg-gray-100 overflow-hidden" :class="activePos.ratio">
          <img v-if="b.desktop_url || b.mobile_url" :src="resolve(b.desktop_url || b.mobile_url)" alt="" class="w-full h-full object-cover" />
        </div>
        <div class="p-3 flex items-center gap-2">
          <div class="flex-1 min-w-0">
            <div class="font-medium text-ink-900 truncate">{{ b.title_ru || b.title_tj || b.link_url || '—' }}</div>
            <div class="text-xs text-gray-400 truncate">{{ b.link_url || '—' }}</div>
          </div>
          <IconButton :variant="b.active ? 'toggle-on' : 'toggle-off'" :title="b.active ? t('common.active') : t('common.inactive')" @click="toggleActive(b)" />
          <IconButton variant="edit" :title="t('common.edit')" @click="startEdit(b)" />
          <IconButton variant="delete" :title="t('common.delete')" @click="remove(b)" />
        </div>
      </div>
    </div>
  </div>
</template>
