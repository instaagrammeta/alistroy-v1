<script setup lang="ts">
import type { PaginatedResponse, Product, Seller } from '~/types/api'
import { productToForm, formToPayload } from '~/composables/useProductForm'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const id = route.params.id as string

const { data: prodRes } = await useAsyncData(`admin-edit-${id}`, () => useApi()<{ data: Product }>(`/admin/products/${id}`))
if (!prodRes.value) throw createError({ statusCode: 404, statusMessage: 'Not found' })
const product = (prodRes.value as any).data as Product

const { data: sellersRes } = await useAsyncData('admin-edit-sellers', () => useApi()<PaginatedResponse<Seller>>('/admin/sellers', { query: { page_size: 200 } }))
const sellers = computed<Seller[]>(() => (sellersRes.value as any)?.data || [])

const form = reactive(productToForm(product))
const sellerId = ref(product.seller_id)
const submitting = ref(false)
const error = ref('')
const success = ref('')

const submit = async () => {
  error.value = ''; success.value = ''; submitting.value = true
  try {
    const body = { ...formToPayload(form), seller_id: sellerId.value }
    await useApi()(`/admin/products/${id}`, { method: 'PATCH', body })
    success.value = t('common.saved')
    setTimeout(() => router.push(`/admin/products/${id}`), 600)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    submitting.value = false
  }
}

useHead({ title: () => `${t('common.edit')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('seller_panel.edit_product') }}</h1>
      <StatusBadge :status="product.status" />
    </div>
    <div v-if="error" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700">{{ error }}</div>
    <div v-if="success" class="card p-3 mb-4 border-l-4 border-emerald-500 bg-emerald-50 text-sm text-emerald-700">{{ success }}</div>
    <div class="card p-6 space-y-4">
      <div>
        <label class="label">{{ t('catalog.seller') }}</label>
        <select v-model="sellerId" class="select">
          <option v-for="s in sellers" :key="s.id" :value="s.id">{{ s.market_name || s.full_name }}</option>
        </select>
      </div>
      <ProductForm v-model="form" upload-endpoint="/admin/upload" :submitting="submitting" :is-admin="true" @submit="submit">
        <template #actions>
          <NuxtLink :to="`/admin/products/${id}`" class="btn-outline">{{ t('common.cancel') }}</NuxtLink>
        </template>
      </ProductForm>
    </div>
  </div>
</template>
