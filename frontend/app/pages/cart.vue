<script setup lang="ts">
const { t } = useI18n()
const auth = useAuthStore()
const cart = useCartStore()
const router = useRouter()
const { formatPrice } = useFormatters()
const { productName } = useLocaleField()
const { resolve } = useImageUrl()

onMounted(() => {
  if (auth.isAuthenticated) cart.load()
})

const address = ref('')
const notes = ref('')
const sending = ref(false)
const sent = ref(false)
const error = ref('')

const setQty = async (productId: string, qty: number) => {
  if (qty < 1) return
  await cart.set(productId, qty)
}
const remove = async (productId: string) => { await cart.remove(productId) }

const checkout = async () => {
  if (!auth.isAuthenticated) return router.push({ path: '/login', query: { redirect: '/cart' } })
  sending.value = true
  error.value = ''
  try {
    await useApi()('/customer/checkout', { method: 'POST', body: { delivery_address: address.value, notes: notes.value } })
    await cart.load()
    sent.value = true
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    sending.value = false
  }
}

useHead({ title: () => `${t('cart.title')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6">
    <h1 class="text-2xl md:text-3xl font-bold text-ink-900 mb-6">{{ t('cart.title') }}</h1>

    <div v-if="sent" class="card p-8 text-center">
      <div class="text-5xl mb-3">✅</div>
      <h2 class="text-xl font-semibold text-ink-900">{{ t('cart.order_sent') }}</h2>
      <NuxtLink to="/products" class="btn-primary mt-5 inline-flex">{{ t('cart.continue') }}</NuxtLink>
    </div>

    <template v-else>
      <EmptyState v-if="!cart.items.length" :title="t('cart.empty')">
        <NuxtLink to="/products" class="btn-primary mt-4 inline-flex">{{ t('cart.continue') }}</NuxtLink>
      </EmptyState>

      <div v-else class="grid lg:grid-cols-[1fr_320px] gap-6">
        <div class="space-y-3">
          <div v-for="item in cart.items" :key="item.id" class="card p-3 flex items-center gap-3">
            <div class="w-16 h-16 rounded-lg bg-gray-50 overflow-hidden shrink-0">
              <img v-if="item.product?.images?.[0]" :src="resolve(item.product.images[0].url)" alt="" class="w-full h-full object-cover" />
            </div>
            <div class="min-w-0 flex-1">
              <NuxtLink :to="`/products/${item.product?.slug}`" class="font-medium text-ink-900 hover:text-brand-600 line-clamp-1">{{ productName(item.product) }}</NuxtLink>
              <div class="text-brand-600 font-semibold text-sm">{{ formatPrice(item.product?.sale_price || 0, item.product?.currency) }}</div>
            </div>
            <div class="flex items-center gap-2">
              <button class="w-7 h-7 rounded border border-gray-200" @click="setQty(item.product_id, item.quantity - 1)">−</button>
              <span class="w-8 text-center">{{ item.quantity }}</span>
              <button class="w-7 h-7 rounded border border-gray-200" @click="setQty(item.product_id, item.quantity + 1)">+</button>
            </div>
            <button class="text-red-500 ml-2" @click="remove(item.product_id)">✕</button>
          </div>
        </div>

        <div class="card p-5 h-fit">
          <div class="flex justify-between text-sm mb-2"><span>{{ t('cart.subtotal') }}</span><strong>{{ formatPrice(cart.subtotal) }}</strong></div>
          <div>
            <label class="label mt-2">{{ t('cart.delivery_address') }}</label>
            <input v-model="address" class="input" />
          </div>
          <div>
            <label class="label mt-2">{{ t('cart.notes') }}</label>
            <textarea v-model="notes" class="textarea" rows="2"></textarea>
          </div>
          <div v-if="error" class="text-sm text-red-600 mt-2">{{ error }}</div>
          <button class="btn-primary w-full mt-3" :disabled="sending" @click="checkout">{{ sending ? t('common.loading') : t('cart.checkout') }}</button>
        </div>
      </div>
    </template>
  </div>
</template>
