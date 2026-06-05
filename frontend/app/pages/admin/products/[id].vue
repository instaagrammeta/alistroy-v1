<script setup lang="ts">
import type { Product } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const route = useRoute()
const router = useRouter()
const api = useApi()
const { t } = useI18n()
const { productTitle, productDescription, categoryTitle } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()
const settings = useSettingsStore()

const id = route.params.id as string

const { data, refresh } = await useAsyncData(`admin-product-${id}`, () =>
  api<{ data: Product }>(`/admin/products/${id}`)
)
const product = computed<Product | null>(() => (data.value as any)?.data || null)

if (!product.value) {
  throw createError({ statusCode: 404, statusMessage: 'Product not found' })
}

// Moderation form
const decision = reactive({
  status: 'approved' as 'approved' | 'rejected' | 'pending',
  contact_type: product.value.contact_type as 'admin' | 'seller',
  phone_number: product.value.phone_number || (product.value.contact_type === 'admin' ? settings.marketplacePhone : product.value.seller?.phone || ''),
  whatsapp_number: product.value.whatsapp_number || (product.value.contact_type === 'admin' ? settings.marketplaceWhatsApp : product.value.seller?.whatsapp || ''),
  rejection_note: product.value.rejection_note || '',
})

watch(
  () => decision.contact_type,
  (val) => {
    if (val === 'admin') {
      decision.phone_number = settings.marketplacePhone
      decision.whatsapp_number = settings.marketplaceWhatsApp
    } else if (product.value?.seller) {
      decision.phone_number = product.value.seller.phone
      decision.whatsapp_number = product.value.seller.whatsapp
    }
  }
)

const isFeatured = ref(product.value.is_featured)
const message = ref('')
const error = ref('')
const submitting = ref(false)

const moderate = async () => {
  message.value = ''
  error.value = ''
  submitting.value = true
  try {
    await api(`/admin/products/${id}/moderate`, {
      method: 'POST',
      body: { ...decision },
    })
    if (isFeatured.value !== product.value!.is_featured) {
      await api(`/admin/products/${id}`, {
        method: 'PATCH',
        body: { is_featured: isFeatured.value },
      })
    }
    message.value = t('common.save')
    await refresh()
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
  } finally {
    submitting.value = false
  }
}

const onDelete = async () => {
  if (!confirm(t('seller_panel.delete_confirm'))) return
  await api(`/admin/products/${id}`, { method: 'DELETE' })
  router.push('/admin/products')
}

useHead({ title: () => `${productTitle(product.value)} — Admin` })
</script>

<template>
  <div v-if="product">
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ productTitle(product) }}</h1>
      <StatusBadge :status="product.status" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Product info -->
      <div class="card p-5">
        <div class="grid grid-cols-3 gap-2 mb-3">
          <div v-for="img in product.images" :key="img.id" class="aspect-square rounded overflow-hidden bg-gray-50">
            <img :src="resolve(img.url)" alt="" class="w-full h-full object-cover" />
          </div>
        </div>
        <div class="space-y-2 text-sm">
          <div><span class="text-gray-400">{{ t('catalog.category') }}:</span> {{ categoryTitle(product.category) }}</div>
          <div><span class="text-gray-400">{{ t('catalog.seller') }}:</span> {{ product.seller?.name }}</div>
          <div><span class="text-gray-400">{{ t('catalog.price') }}:</span> {{ formatPrice(product.price, product.currency) }}</div>
          <div><span class="text-gray-400">SKU:</span> {{ product.sku }}</div>
          <div><span class="text-gray-400">Stock:</span> {{ product.stock_quantity }} {{ product.unit }}</div>
        </div>
        <hr class="my-4" />
        <div class="text-sm text-gray-700 whitespace-pre-line">{{ productDescription(product) || '—' }}</div>
      </div>

      <!-- Moderation -->
      <div class="card p-5">
        <h2 class="font-semibold mb-4">{{ t('admin.moderate_now') }}</h2>
        <form class="space-y-4" @submit.prevent="moderate">
          <div>
            <label class="label">Status</label>
            <select v-model="decision.status" class="select">
              <option value="approved">{{ t('admin.approve') }}</option>
              <option value="rejected">{{ t('admin.reject') }}</option>
              <option value="pending">{{ t('seller_panel.status_pending') }}</option>
            </select>
          </div>

          <div v-if="decision.status === 'rejected'">
            <label class="label">{{ t('admin.rejection_note') }}</label>
            <textarea v-model="decision.rejection_note" class="textarea" rows="3"></textarea>
          </div>

          <div v-else>
            <label class="label">{{ t('admin.contact_routing') }}</label>
            <div class="flex gap-2">
              <button
                type="button"
                class="px-3 py-2 rounded-lg text-sm font-medium border flex-1"
                :class="decision.contact_type === 'admin' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'"
                @click="decision.contact_type = 'admin'"
              >{{ t('admin.contact_admin') }}</button>
              <button
                type="button"
                class="px-3 py-2 rounded-lg text-sm font-medium border flex-1"
                :class="decision.contact_type === 'seller' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'"
                @click="decision.contact_type = 'seller'"
              >{{ t('admin.contact_seller') }}</button>
            </div>
          </div>

          <div v-if="decision.status !== 'rejected'" class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div>
              <label class="label">{{ t('admin.phone_number') }}</label>
              <input v-model="decision.phone_number" class="input" />
            </div>
            <div>
              <label class="label">{{ t('admin.whatsapp_number') }}</label>
              <input v-model="decision.whatsapp_number" class="input" />
            </div>
          </div>

          <label class="inline-flex items-center gap-2">
            <input v-model="isFeatured" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            <span class="text-sm">Featured</span>
          </label>

          <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
          <div v-if="message" class="text-sm text-emerald-600">{{ message }}</div>

          <div class="flex items-center gap-3 pt-1">
            <button type="submit" class="btn-primary" :disabled="submitting">{{ submitting ? t('common.loading') : t('common.save') }}</button>
            <button type="button" class="btn-danger" @click="onDelete">{{ t('common.delete') }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
