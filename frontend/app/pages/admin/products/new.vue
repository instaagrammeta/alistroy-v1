<script setup lang="ts">
import type { PaginatedResponse, Seller } from '~/types/api'
import { emptyProductForm, formToPayload } from '~/composables/useProductForm'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const router = useRouter()

const form = reactive(emptyProductForm())
const sellerId = ref('')
const contactOwner = ref<'admin' | 'seller'>('admin')
const submitting = ref(false)
const error = ref('')

const { data: sellersRes } = await useAsyncData('admin-new-sellers', () => useApi()<PaginatedResponse<Seller>>('/admin/sellers', { query: { page_size: 200 } }))
const sellers = computed<Seller[]>(() => (sellersRes.value as any)?.data || [])

const submit = async () => {
  error.value = ''
  if (!sellerId.value) { error.value = t('admin.select_seller'); return }
  submitting.value = true
  try {
    const payload = { ...formToPayload(form), seller_id: sellerId.value, contact_owner: contactOwner.value, status: 'approved' }
    await useApi()('/admin/products', { method: 'POST', body: payload })
    router.push('/admin/products')
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    submitting.value = false
  }
}

useHead({ title: () => `${t('admin.new_product')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('admin.new_product') }}</h1>
    <div v-if="error" class="card p-3 mb-4 border-l-4 border-red-500 bg-red-50 text-sm text-red-700">{{ error }}</div>
    <div class="card p-6 space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="label">{{ t('catalog.seller') }}</label>
          <select v-model="sellerId" class="select" required>
            <option value="" disabled>{{ t('admin.select_seller') }}</option>
            <option v-for="s in sellers" :key="s.id" :value="s.id">{{ s.market_name || s.full_name }}</option>
          </select>
        </div>
        <div>
          <label class="label">{{ t('admin.contact_routing') }}</label>
          <div class="flex gap-2">
            <button type="button" class="px-3 py-2 rounded-lg text-sm font-medium border flex-1" :class="contactOwner === 'admin' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'" @click="contactOwner = 'admin'">{{ t('admin.contact_admin') }}</button>
            <button type="button" class="px-3 py-2 rounded-lg text-sm font-medium border flex-1" :class="contactOwner === 'seller' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'" @click="contactOwner = 'seller'">{{ t('admin.contact_seller') }}</button>
          </div>
        </div>
      </div>
      <ProductForm v-model="form" upload-endpoint="/admin/upload" :submitting="submitting" :is-admin="true" @submit="submit">
        <template #actions><NuxtLink to="/admin/products" class="btn-outline">{{ t('common.cancel') }}</NuxtLink></template>
      </ProductForm>
    </div>
  </div>
</template>
