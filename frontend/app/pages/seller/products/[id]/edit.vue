<script setup lang="ts">
import type { PaginatedResponse, Product } from '~/types/api'
import { productToForm, formToPayload } from '~/composables/useProductForm'

definePageMeta({ layout: 'seller', middleware: 'seller' })
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const id = route.params.id as string

const { data } = await useAsyncData(`seller-product-${id}`, async () => {
  const res = await useApi()<PaginatedResponse<Product>>('/seller/products', { query: { page_size: 200 } })
  return (res.data || []).find((p) => p.id === id) || null
})
if (!data.value) throw createError({ statusCode: 404, statusMessage: 'Not found' })
const product = data.value as Product

const form = reactive(productToForm(product))
const submitting = ref(false)
const error = ref('')
const success = ref('')

const submit = async () => {
  error.value = ''; success.value = ''; submitting.value = true
  try {
    await useApi()(`/seller/products/${id}`, { method: 'PATCH', body: formToPayload(form) })
    success.value = t('seller_panel.submitted')
    setTimeout(() => router.push('/seller/products'), 700)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    submitting.value = false
  }
}

useHead({ title: () => `${t('seller_panel.edit_product')} — AliStroy` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('seller_panel.edit_product') }}</h1>
      <StatusBadge :status="product.status" />
    </div>
    <div v-if="product.status === 'rejected' && product.rejection_note" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700"><strong>{{ t('admin.rejection_note') }}:</strong> {{ product.rejection_note }}</div>
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
