<script setup lang="ts">
import type { Seller } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })
const { t } = useI18n()

const { data: meRes, refresh } = await useAsyncData('seller-profile', () => useApi()<{ data: Seller }>('/seller/me'))
const seller = (meRes.value as any)?.data as Seller

const form = reactive({
  full_name: seller.full_name, company_name: seller.company_name, market_name: seller.market_name,
  phone: seller.phone, phone_alt: seller.phone_alt, whatsapp: seller.whatsapp,
  telegram: seller.telegram, telegram_username: seller.telegram_username,
  address: seller.address, city: seller.city, notes: seller.notes, logo_url: seller.logo_url,
})
const msg = ref(''); const err = ref(''); const saving = ref(false)

const save = async () => {
  msg.value = ''; err.value = ''; saving.value = true
  try {
    await useApi()('/seller/me', { method: 'PATCH', body: form })
    msg.value = t('common.saved')
    await refresh()
  } catch (e: any) {
    err.value = e?.data?.error?.message || 'Error'
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
      <form class="space-y-4" @submit.prevent="save">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div><label class="label">{{ t('admin.full_name') }}</label><input v-model="form.full_name" class="input" /></div>
          <div><label class="label">{{ t('admin.market') }}</label><input v-model="form.market_name" class="input" /></div>
          <div><label class="label">{{ t('admin.company') }}</label><input v-model="form.company_name" class="input" /></div>
          <div><label class="label">{{ t('seller.city') }}</label><input v-model="form.city" class="input" /></div>
          <div><label class="label">{{ t('seller.address') }}</label><input v-model="form.address" class="input" /></div>
          <div><label class="label">{{ t('admin.phone') }}</label><input v-model="form.phone" class="input" /></div>
          <div><label class="label">WhatsApp</label><input v-model="form.whatsapp" class="input" /></div>
          <div><label class="label">Telegram</label><input v-model="form.telegram" class="input" /></div>
        </div>
        <div><label class="label">Logo</label><ImageUploader v-model="logoArr" endpoint="/seller/upload" subdir="logos" :multiple="false" /></div>
        <div v-if="msg" class="text-sm text-emerald-600">{{ msg }}</div>
        <div v-if="err" class="text-sm text-red-600">{{ err }}</div>
        <button type="submit" class="btn-primary" :disabled="saving">{{ t('common.save') }}</button>
      </form>
    </div>
  </div>
</template>
