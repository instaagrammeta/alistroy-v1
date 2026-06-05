<script setup lang="ts">
import type { Product } from '~/types/api'

definePageMeta({ layout: 'seller', middleware: 'seller' })

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const id = route.params.id as string

// Seller-scoped read: load via /seller/products?... is non-trivial for a single id.
// The seller endpoint already returns the same record; for simplicity we use the
// admin/public details: products are loaded by id via GET on the public list,
// but the cleanest way is via the seller list. We instead fetch the full record
// through the public list filtered by id.
const { data: prodRes, error: fetchError } = await useAsyncData(`seller-product-${id}`, async () => {
  // Fetch the seller's products and find the one we need.
  const res = await api<{ data: Product[] }>('/seller/products', { query: { page_size: 100 } })
  return (res.data || []).find((p) => p.id === id) || null
})

if (!prodRes.value) {
  throw createError({ statusCode: 404, statusMessage: 'Product not found' })
}

const product = prodRes.value as Product

const form = reactive({
  category_id: product.category_id,
  title_tj: product.title_tj,
  title_ru: product.title_ru,
  description_tj: product.description_tj,
  description_ru: product.description_ru,
  price: product.price,
  currency: product.currency || 'TJS',
  unit: product.unit || 'pcs',
  sku: product.sku || '',
  stock_quantity: product.stock_quantity,
  is_available: product.is_available,
  images: (product.images || []).map((i) => ({ url: i.url, alt: i.alt })),
})

const submitting = ref(false)
const error = ref('')
const success = ref('')

const submit = async () => {
  error.value = ''
  success.value = ''
  submitting.value = true
  try {
    await api(`/seller/products/${id}`, {
      method: 'PATCH',
      body: { ...form },
    })
    success.value = t('seller_panel.submitted_for_review')
    setTimeout(() => router.push('/seller/products'), 700)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
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
    <div v-if="product.status === 'rejected' && product.rejection_note" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700">
      <strong>{{ t('admin.rejection_note') }}:</strong> {{ product.rejection_note }}
    </div>
    <div v-if="error" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700">{{ error }}</div>
    <div v-if="success" class="card p-3 mb-4 border-l-4 border-emerald-500 bg-emerald-50 text-sm text-emerald-700">{{ success }}</div>
    <div class="card p-6">
      <ProductForm v-model="form" upload-endpoint="/seller/upload" :submitting="submitting" @submit="submit" />
    </div>
  </div>
</template>
