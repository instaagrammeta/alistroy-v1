<script setup lang="ts">
import { emptyProductForm, formToPayload } from '~/composables/useProductForm'

definePageMeta({ layout: 'seller', middleware: 'seller' })
const { t } = useI18n()
const router = useRouter()

const form = reactive(emptyProductForm())
const submitting = ref(false)
const error = ref('')
const success = ref('')

const submit = async () => {
  error.value = ''; success.value = ''; submitting.value = true
  try {
    await useApi()('/seller/products', { method: 'POST', body: formToPayload(form) })
    success.value = t('seller_panel.submitted')
    setTimeout(() => router.push('/seller/products'), 700)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    submitting.value = false
  }
}

useHead({ title: () => `${t('seller_panel.add_product')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('seller_panel.add_product') }}</h1>
    <div v-if="error" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700">{{ error }}</div>
    <div v-if="success" class="card p-3 mb-4 border-l-4 border-emerald-500 bg-emerald-50 text-sm text-emerald-700">{{ success }}</div>
    <div class="card p-6">
      <ProductForm v-model="form" upload-endpoint="/seller/upload" :submitting="submitting" @submit="submit">
        <template #actions>
          <NuxtLink to="/seller/products" class="btn-outline">{{ t('common.cancel') }}</NuxtLink>
        </template>
      </ProductForm>
    </div>
  </div>
</template>
