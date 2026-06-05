<script setup lang="ts">
import type { Customer, PaginatedResponse, Product } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const router = useRouter()
const { formatPrice } = useFormatters()
const { productName } = useLocaleField()

const { data: customersRes } = await useAsyncData('order-new-customers', () => useApi()<PaginatedResponse<Customer>>('/admin/customers', { query: { page_size: 200 } }))
const customers = computed<Customer[]>(() => (customersRes.value as any)?.data || [])

const form = reactive<any>({ customer_id: '', customer_name: '', customer_phone: '', delivery_address: '', delivery_date: '', discount_percent: 0, status: 'new', notes: '' })
const items = ref<{ product: Product; quantity: number }[]>([])

const search = ref('')
const results = ref<Product[]>([])
let timer: any = null
watch(search, (v) => {
  clearTimeout(timer)
  if (!v.trim()) { results.value = []; return }
  timer = setTimeout(async () => {
    const res = await useApi()<PaginatedResponse<Product>>('/admin/products', { query: { q: v, status: 'approved', page_size: 10 } })
    results.value = res.data || []
  }, 300)
})
const addItem = (p: Product) => {
  if (!items.value.find((i) => i.product.id === p.id)) items.value.push({ product: p, quantity: 1 })
  search.value = ''; results.value = []
}
const removeItem = (idx: number) => items.value.splice(idx, 1)

const subtotal = computed(() => items.value.reduce((s, i) => s + i.product.sale_price * i.quantity, 0))
const total = computed(() => subtotal.value * (1 - (form.discount_percent || 0) / 100))

const error = ref(''); const saving = ref(false)
const save = async () => {
  error.value = ''
  if (!items.value.length) { error.value = t('admin.add_product_to_order'); return }
  saving.value = true
  try {
    const body = {
      customer_id: form.customer_id || undefined,
      customer_name: form.customer_name, customer_phone: form.customer_phone,
      delivery_address: form.delivery_address, delivery_date: form.delivery_date || undefined,
      discount_percent: form.discount_percent, status: form.status, notes: form.notes,
      items: items.value.map((i) => ({ product_id: i.product.id, quantity: i.quantity })),
    }
    const res = await useApi()<{ data: { id: string } }>('/admin/orders', { method: 'POST', body })
    router.push(`/admin/orders/${res.data.id}`)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    saving.value = false
  }
}

useHead({ title: () => `${t('admin.new_order')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-6">{{ t('admin.new_order') }}</h1>
    <div class="grid lg:grid-cols-[1fr_320px] gap-6">
      <div class="space-y-4">
        <div class="card p-5 grid grid-cols-1 md:grid-cols-2 gap-3">
          <div>
            <label class="label">{{ t('admin.select_customer') }}</label>
            <select v-model="form.customer_id" class="select"><option value="">— {{ t('common.name') }} —</option><option v-for="c in customers" :key="c.id" :value="c.id">{{ c.company || c.id.slice(0, 8) }}</option></select>
          </div>
          <div><label class="label">{{ t('common.name') }}</label><input v-model="form.customer_name" class="input" /></div>
          <div><label class="label">{{ t('admin.phone') }}</label><input v-model="form.customer_phone" class="input" /></div>
          <div><label class="label">{{ t('cart.delivery_address') }}</label><input v-model="form.delivery_address" class="input" /></div>
          <div><label class="label">{{ t('common.date') }}</label><input v-model="form.delivery_date" type="date" class="input" /></div>
          <div><label class="label">{{ t('seller_panel.discount') }}</label><input v-model.number="form.discount_percent" type="number" min="0" max="100" class="input" /></div>
        </div>

        <div class="card p-5">
          <label class="label">{{ t('admin.search_product') }}</label>
          <input v-model="search" class="input" :placeholder="t('admin.search_product')" />
          <div v-if="results.length" class="mt-2 border border-gray-100 rounded-lg divide-y">
            <button v-for="p in results" :key="p.id" class="w-full text-left px-3 py-2 hover:bg-gray-50 flex justify-between" @click="addItem(p)">
              <span>{{ productName(p) }}</span><span class="text-brand-600">{{ formatPrice(p.sale_price, p.currency) }}</span>
            </button>
          </div>
          <div class="mt-4 space-y-2">
            <div v-for="(it, idx) in items" :key="it.product.id" class="flex items-center gap-3">
              <span class="flex-1 truncate">{{ productName(it.product) }}</span>
              <input v-model.number="it.quantity" type="number" min="1" class="input w-20" />
              <span class="w-24 text-right">{{ formatPrice(it.product.sale_price * it.quantity, it.product.currency) }}</span>
              <button class="text-red-500" @click="removeItem(idx)">✕</button>
            </div>
            <EmptyState v-if="!items.length" :title="t('admin.add_product_to_order')" />
          </div>
        </div>
      </div>

      <div class="card p-5 h-fit">
        <div class="flex justify-between text-sm mb-1"><span>{{ t('cart.subtotal') }}</span><strong>{{ formatPrice(subtotal) }}</strong></div>
        <div class="flex justify-between text-lg font-bold mb-3"><span>{{ t('common.total') }}</span><span class="text-brand-600">{{ formatPrice(total) }}</span></div>
        <div><label class="label">{{ t('common.status') }}</label><select v-model="form.status" class="select"><option v-for="s in ['new','processing','assigned','on_delivery','completed','cancelled']" :key="s" :value="s">{{ t(`order_status.${s}`) }}</option></select></div>
        <div class="mt-2"><label class="label">{{ t('cart.notes') }}</label><textarea v-model="form.notes" class="textarea" rows="2"></textarea></div>
        <div v-if="error" class="text-sm text-red-600 mt-2">{{ error }}</div>
        <button class="btn-primary w-full mt-3" :disabled="saving" @click="save">{{ t('common.save') }}</button>
      </div>
    </div>
  </div>
</template>
