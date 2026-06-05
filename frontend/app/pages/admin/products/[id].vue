<script setup lang="ts">
import type { Product } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const { productName, productDesc, categoryName } = useLocaleField()
const { formatPrice } = useFormatters()
const { resolve } = useImageUrl()
const settings = useSettingsStore()
const id = route.params.id as string

const { data, refresh } = await useAsyncData(`admin-product-${id}`, () => useApi()<{ data: Product }>(`/admin/products/${id}`))
const product = computed<Product | null>(() => (data.value as any)?.data || null)
if (!product.value) throw createError({ statusCode: 404, statusMessage: 'Not found' })

const decision = reactive({
  status: (product.value.status === 'approved' ? 'approved' : 'approved') as 'approved' | 'rejected' | 'pending',
  sale_price: product.value.sale_price || product.value.cost_price,
  contact_owner: product.value.contact_owner as 'admin' | 'seller',
  contact_phone: product.value.contact_phone || (product.value.contact_owner === 'admin' ? settings.marketplacePhone : product.value.seller?.phone || ''),
  contact_whatsapp: product.value.contact_whatsapp || (product.value.contact_owner === 'admin' ? settings.marketplaceWhatsApp : product.value.seller?.whatsapp || ''),
  contact_telegram: product.value.contact_telegram || (product.value.contact_owner === 'admin' ? settings.marketplaceTelegram : product.value.seller?.telegram || ''),
  rejection_note: product.value.rejection_note || '',
  is_featured: product.value.is_featured,
})

watch(() => decision.contact_owner, (v) => {
  if (v === 'admin') {
    decision.contact_phone = settings.marketplacePhone
    decision.contact_whatsapp = settings.marketplaceWhatsApp
    decision.contact_telegram = settings.marketplaceTelegram
  } else if (product.value?.seller) {
    decision.contact_phone = product.value.seller.phone
    decision.contact_whatsapp = product.value.seller.whatsapp
    decision.contact_telegram = product.value.seller.telegram
  }
})

const msg = ref(''); const err = ref(''); const submitting = ref(false)

const moderate = async () => {
  msg.value = ''; err.value = ''; submitting.value = true
  try {
    await useApi()(`/admin/products/${id}/moderate`, { method: 'POST', body: { ...decision } })
    msg.value = t('common.saved')
    await refresh()
  } catch (e: any) {
    err.value = e?.data?.error?.message || 'Error'
  } finally {
    submitting.value = false
  }
}

const onDelete = async () => {
  if (!confirm(t('common.confirm_delete'))) return
  await useApi()(`/admin/products/${id}`, { method: 'DELETE' })
  router.push('/admin/products')
}

useHead({ title: () => `${productName(product.value)} — Admin` })
</script>

<template>
  <div v-if="product">
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-ink-900">{{ productName(product) }}</h1>
      <StatusBadge :status="product.status" />
    </div>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card p-5">
        <div class="grid grid-cols-3 gap-2 mb-3">
          <div v-for="img in product.images" :key="img.id" class="aspect-square rounded overflow-hidden bg-gray-50"><img :src="resolve(img.url)" alt="" class="w-full h-full object-cover" /></div>
        </div>
        <div class="space-y-1.5 text-sm">
          <div><span class="text-gray-400">{{ t('catalog.category') }}:</span> {{ categoryName(product.category) }}</div>
          <div><span class="text-gray-400">{{ t('catalog.seller') }}:</span> {{ product.seller?.market_name || product.seller?.full_name }}</div>
          <div><span class="text-gray-400">{{ t('seller_panel.cost') }}:</span> {{ formatPrice(product.cost_price, product.currency) }}</div>
          <div><span class="text-gray-400">{{ t('admin.sale_price') }}:</span> {{ formatPrice(product.sale_price, product.currency) }}</div>
          <div><span class="text-gray-400">{{ t('seller_panel.stock') }}:</span> {{ product.stock_quantity }} {{ product.unit }}</div>
        </div>
        <hr class="my-3" />
        <p class="text-sm text-gray-700 whitespace-pre-line">{{ productDesc(product) || '—' }}</p>
      </div>

      <div class="card p-5">
        <h2 class="font-semibold mb-4">{{ t('admin.moderate') }}</h2>
        <form class="space-y-4" @submit.prevent="moderate">
          <div>
            <label class="label">{{ t('common.status') }}</label>
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
          <template v-else>
            <div>
              <label class="label">{{ t('admin.sale_price') }}</label>
              <input v-model.number="decision.sale_price" type="number" min="0" step="0.01" class="input" />
            </div>
            <div>
              <label class="label">{{ t('admin.contact_routing') }}</label>
              <div class="flex gap-2">
                <button type="button" class="px-3 py-2 rounded-lg text-sm font-medium border flex-1" :class="decision.contact_owner === 'admin' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'" @click="decision.contact_owner = 'admin'">{{ t('admin.contact_admin') }}</button>
                <button type="button" class="px-3 py-2 rounded-lg text-sm font-medium border flex-1" :class="decision.contact_owner === 'seller' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'" @click="decision.contact_owner = 'seller'">{{ t('admin.contact_seller') }}</button>
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div><label class="label">{{ t('admin.phone') }}</label><input v-model="decision.contact_phone" class="input" /></div>
              <div><label class="label">WhatsApp</label><input v-model="decision.contact_whatsapp" class="input" /></div>
              <div><label class="label">Telegram</label><input v-model="decision.contact_telegram" class="input" /></div>
            </div>
            <label class="inline-flex items-center gap-2"><input v-model="decision.is_featured" type="checkbox" class="h-4 w-4" /><span class="text-sm">Featured</span></label>
          </template>
          <div v-if="err" class="text-sm text-red-600">{{ err }}</div>
          <div v-if="msg" class="text-sm text-emerald-600">{{ msg }}</div>
          <div class="flex gap-3 pt-1">
            <button type="submit" class="btn-primary" :disabled="submitting">{{ t('common.save') }}</button>
            <button type="button" class="btn-danger" @click="onDelete">{{ t('common.delete') }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
