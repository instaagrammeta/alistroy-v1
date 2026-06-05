<script setup lang="ts">
import type { SettingsMap } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { resolve } = useImageUrl()
const settings = useSettingsStore()

const { data, refresh } = await useAsyncData('admin-settings', () => useApi()<{ data: SettingsMap }>('/admin/settings'))
const initial = (data.value as any)?.data || {}

const form = reactive<Record<string, string>>({
  logo_url: initial.logo_url || '', favicon_url: initial.favicon_url || '',
  site_name_tj: initial.site_name_tj || '', site_name_ru: initial.site_name_ru || '',
  tagline_tj: initial.tagline_tj || '', tagline_ru: initial.tagline_ru || '',
  hero_title_tj: initial.hero_title_tj || '', hero_title_ru: initial.hero_title_ru || '',
  hero_subtitle_tj: initial.hero_subtitle_tj || '', hero_subtitle_ru: initial.hero_subtitle_ru || '',
  seo_description_tj: initial.seo_description_tj || '', seo_description_ru: initial.seo_description_ru || '',
  marketplace_phone: initial.marketplace_phone || '', marketplace_whatsapp: initial.marketplace_whatsapp || '',
  marketplace_telegram: initial.marketplace_telegram || '', footer_email: initial.footer_email || '', footer_address: initial.footer_address || '',
})

const logoArr = computed({ get: () => (form.logo_url ? [{ url: form.logo_url }] : []), set: (v: any[]) => (form.logo_url = v[0]?.url || '') })
const faviconArr = computed({ get: () => (form.favicon_url ? [{ url: form.favicon_url }] : []), set: (v: any[]) => (form.favicon_url = v[0]?.url || '') })

const msg = ref(''); const err = ref(''); const saving = ref(false)
const save = async () => {
  msg.value = ''; err.value = ''; saving.value = true
  try {
    await useApi()('/admin/settings', { method: 'PATCH', body: { items: form } })
    settings.setLocal(form)
    msg.value = t('admin.settings_saved')
    await refresh()
  } catch (e: any) {
    err.value = e?.data?.error?.message || 'Error'
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.settings')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('admin.settings') }}</h1>
    <form class="space-y-6" @submit.prevent="save">
      <div class="card p-5 grid grid-cols-1 md:grid-cols-2 gap-5">
        <div><label class="label">{{ t('admin.logo') }}</label><ImageUploader v-model="logoArr" endpoint="/admin/upload" subdir="brand" :multiple="false" /><img v-if="form.logo_url" :src="resolve(form.logo_url)" alt="" class="h-12 mt-3" /></div>
        <div><label class="label">{{ t('admin.favicon') }}</label><ImageUploader v-model="faviconArr" endpoint="/admin/upload" subdir="brand" :multiple="false" /><img v-if="form.favicon_url" :src="resolve(form.favicon_url)" alt="" class="w-8 h-8 mt-3" /></div>
      </div>
      <div class="card p-5 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div><label class="label">{{ t('admin.site_name') }} (TJ)</label><input v-model="form.site_name_tj" class="input" /></div>
        <div><label class="label">{{ t('admin.site_name') }} (RU)</label><input v-model="form.site_name_ru" class="input" /></div>
        <div><label class="label">{{ t('admin.tagline') }} (TJ)</label><input v-model="form.tagline_tj" class="input" /></div>
        <div><label class="label">{{ t('admin.tagline') }} (RU)</label><input v-model="form.tagline_ru" class="input" /></div>
        <div><label class="label">Hero Title (TJ)</label><input v-model="form.hero_title_tj" class="input" /></div>
        <div><label class="label">Hero Title (RU)</label><input v-model="form.hero_title_ru" class="input" /></div>
        <div><label class="label">Hero Subtitle (TJ)</label><textarea v-model="form.hero_subtitle_tj" class="textarea" rows="2"></textarea></div>
        <div><label class="label">Hero Subtitle (RU)</label><textarea v-model="form.hero_subtitle_ru" class="textarea" rows="2"></textarea></div>
        <div><label class="label">SEO (TJ)</label><textarea v-model="form.seo_description_tj" class="textarea" rows="2"></textarea></div>
        <div><label class="label">SEO (RU)</label><textarea v-model="form.seo_description_ru" class="textarea" rows="2"></textarea></div>
      </div>
      <div class="card p-5 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div><label class="label">{{ t('admin.marketplace_phone') }}</label><input v-model="form.marketplace_phone" class="input" /></div>
        <div><label class="label">{{ t('admin.marketplace_whatsapp') }}</label><input v-model="form.marketplace_whatsapp" class="input" /></div>
        <div><label class="label">{{ t('admin.marketplace_telegram') }}</label><input v-model="form.marketplace_telegram" class="input" /></div>
        <div><label class="label">Footer email</label><input v-model="form.footer_email" class="input" /></div>
        <div><label class="label">Footer address</label><input v-model="form.footer_address" class="input" /></div>
      </div>
      <div v-if="msg" class="text-sm text-emerald-600">{{ msg }}</div>
      <div v-if="err" class="text-sm text-red-600">{{ err }}</div>
      <button type="submit" class="btn-primary" :disabled="saving">{{ t('admin.save_settings') }}</button>
    </form>
  </div>
</template>
