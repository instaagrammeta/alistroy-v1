<script setup lang="ts">
import type { Seller } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })

const api = useApi()
const { t } = useI18n()

const { data: meRes, refresh } = await useAsyncData('seller-profile', () =>
  api<{ data: Seller }>('/seller/me')
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
})

const message = ref('')
const error = ref('')
const saving = ref(false)

const submit = async () => {
  error.value = ''
  message.value = ''
  saving.value = true
  try {
    await api('/seller/me', { method: 'PATCH', body: form })
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

useHead({ title: () => `${t('seller_panel.profile')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('seller_panel.profile') }}</h1>

    <div class="card p-6">
      <form class="space-y-5" @submit.prevent="submit">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('auth.seller_name') }}</label>
            <input v-model="form.name" class="input" required />
          </div>
          <div>
            <label class="label">{{ t('seller.city') }}</label>
            <input v-model="form.city" class="input" />
          </div>
        </div>

        <div>
          <label class="label">{{ t('seller.address') }}</label>
          <input v-model="form.address" class="input" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('admin.phone_number') }}</label>
            <input v-model="form.phone" class="input" />
          </div>
          <div>
            <label class="label">{{ t('admin.whatsapp_number') }}</label>
            <input v-model="form.whatsapp" class="input" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
          <ImageUploader v-model="logoArr" endpoint="/seller/upload" subdir="logos" :multiple="false" />
        </div>

        <div v-if="message" class="text-sm text-emerald-600">{{ message }}</div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>

        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? t('common.loading') : t('common.save') }}</button>
      </form>
    </div>
  </div>
</template>
