<script setup lang="ts">
import type { Seller } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const api = useApi()
const { t } = useI18n()
const id = route.params.id as string

const { data: meRes, refresh } = await useAsyncData(`admin-seller-${id}`, () =>
  api<{ data: Seller }>(`/admin/sellers/${id}`)
)
const seller = (meRes.value as any)?.data as Seller

const form = reactive({
  name: seller.name,
  description_tj: seller.description_tj,
  description_ru: seller.description_ru,
  logo_url: seller.logo_url,
  phone: seller.phone,
  whatsapp: seller.whatsapp,
  address: seller.address,
  city: seller.city,
  status: seller.status,
  is_featured: seller.is_featured,
})

const message = ref('')
const error = ref('')
const saving = ref(false)

const submit = async () => {
  message.value = ''
  error.value = ''
  saving.value = true
  try {
    await api(`/admin/sellers/${id}`, { method: 'PATCH', body: form })
    message.value = t('common.save')
    await refresh()
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
  } finally {
    saving.value = false
  }
}

const logoArr = computed({
  get: () => (form.logo_url ? [{ url: form.logo_url }] : []),
  set: (v: { url: string }[]) => (form.logo_url = v[0]?.url || ''),
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ seller.name }}</h1>

    <div class="card p-6">
      <form class="space-y-4" @submit.prevent="submit">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('auth.seller_name') }}</label>
            <input v-model="form.name" class="input" />
          </div>
          <div>
            <label class="label">Status</label>
            <select v-model="form.status" class="select">
              <option value="pending">pending</option>
              <option value="approved">approved</option>
              <option value="blocked">blocked</option>
            </select>
          </div>
          <div>
            <label class="label">{{ t('seller.city') }}</label>
            <input v-model="form.city" class="input" />
          </div>
          <div>
            <label class="label">{{ t('seller.address') }}</label>
            <input v-model="form.address" class="input" />
          </div>
          <div>
            <label class="label">{{ t('admin.phone_number') }}</label>
            <input v-model="form.phone" class="input" />
          </div>
          <div>
            <label class="label">{{ t('admin.whatsapp_number') }}</label>
            <input v-model="form.whatsapp" class="input" />
          </div>
          <div>
            <label class="label">{{ t('seller_panel.desc_tj') }}</label>
            <textarea v-model="form.description_tj" class="textarea" rows="4"></textarea>
          </div>
          <div>
            <label class="label">{{ t('seller_panel.desc_ru') }}</label>
            <textarea v-model="form.description_ru" class="textarea" rows="4"></textarea>
          </div>
        </div>

        <div>
          <label class="label">Logo</label>
          <ImageUploader v-model="logoArr" endpoint="/admin/upload" subdir="logos" :multiple="false" />
        </div>

        <label class="inline-flex items-center gap-2">
          <input v-model="form.is_featured" type="checkbox" class="h-4 w-4" />
          <span class="text-sm">Featured</span>
        </label>

        <div v-if="message" class="text-sm text-emerald-600">{{ message }}</div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>

        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? t('common.loading') : t('common.save') }}</button>
      </form>
    </div>
  </div>
</template>
