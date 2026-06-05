<script setup lang="ts">
import type { Brand, Category, Subcategory } from '~/types/api'

interface ImageItem { url: string; alt?: string }
export interface ProductFormState {
  category_id: string
  subcategory_id: string
  brand_id: string
  name_tj: string
  name_ru: string
  description_tj: string
  description_ru: string
  sku: string
  unit: string
  currency: string
  sale_price: number
  cost_price: number
  discount_percent: number
  stock_quantity: number
  minimum_stock: number
  is_available: boolean
  images: ImageItem[]
}

const props = defineProps<{
  modelValue: ProductFormState
  uploadEndpoint?: string
  submitting?: boolean
  isAdmin?: boolean
}>()
const emit = defineEmits<{ (e: 'update:modelValue', v: ProductFormState): void; (e: 'submit'): void }>()

const { t } = useI18n()
const { categoryName } = useLocaleField()

const { data: catsRes } = await useAsyncData('form-cats', () => $fetch<{ data: Category[] }>('/categories', { baseURL: useApiBase() }))
const { data: brandsRes } = await useAsyncData('form-brands', () => $fetch<{ data: Brand[] }>('/brands', { baseURL: useApiBase() }))
const cats = computed<Category[]>(() => catsRes.value?.data || [])
const brands = computed<Brand[]>(() => brandsRes.value?.data || [])

const local = computed({ get: () => props.modelValue, set: (v) => emit('update:modelValue', v) })

const subcats = computed<Subcategory[]>(() => {
  const c = cats.value.find((x) => x.id === local.value.category_id)
  return c?.subcategories || []
})

const onImages = (next: ImageItem[]) => emit('update:modelValue', { ...local.value, images: next })
</script>

<template>
  <form class="space-y-5" @submit.prevent="emit('submit')">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div><label class="label">{{ t('seller_panel.name_tj') }}</label><input v-model="local.name_tj" class="input" required /></div>
      <div><label class="label">{{ t('seller_panel.name_ru') }}</label><input v-model="local.name_ru" class="input" required /></div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div><label class="label">{{ t('seller_panel.desc_tj') }}</label><textarea v-model="local.description_tj" class="textarea" rows="4"></textarea></div>
      <div><label class="label">{{ t('seller_panel.desc_ru') }}</label><textarea v-model="local.description_ru" class="textarea" rows="4"></textarea></div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div>
        <label class="label">{{ t('seller_panel.category') }}</label>
        <select v-model="local.category_id" class="select" required>
          <option value="" disabled>{{ t('common.select') }}</option>
          <option v-for="c in cats" :key="c.id" :value="c.id">{{ categoryName(c) }}</option>
        </select>
      </div>
      <div>
        <label class="label">{{ t('seller_panel.subcategory') }}</label>
        <select v-model="local.subcategory_id" class="select">
          <option value="">—</option>
          <option v-for="s in subcats" :key="s.id" :value="s.id">{{ categoryName(s) }}</option>
        </select>
      </div>
      <div>
        <label class="label">{{ t('seller_panel.brand') }}</label>
        <select v-model="local.brand_id" class="select">
          <option value="">—</option>
          <option v-for="b in brands" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
      </div>
    </div>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div><label class="label">{{ isAdmin ? t('seller_panel.price') : t('admin.sale_price') }}</label><input v-model.number="local.sale_price" type="number" min="0" step="0.01" class="input" /></div>
      <div v-if="isAdmin"><label class="label">{{ t('seller_panel.cost') }}</label><input v-model.number="local.cost_price" type="number" min="0" step="0.01" class="input" /></div>
      <div><label class="label">{{ t('seller_panel.currency') }}</label><input v-model="local.currency" class="input" maxlength="8" /></div>
      <div><label class="label">{{ t('seller_panel.discount') }}</label><input v-model.number="local.discount_percent" type="number" min="0" max="100" class="input" /></div>
    </div>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div><label class="label">{{ t('seller_panel.sku') }}</label><input v-model="local.sku" class="input" maxlength="64" /></div>
      <div><label class="label">{{ t('seller_panel.unit') }}</label><input v-model="local.unit" class="input" maxlength="20" /></div>
      <div><label class="label">{{ t('seller_panel.stock') }}</label><input v-model.number="local.stock_quantity" type="number" min="0" class="input" /></div>
      <div><label class="label">{{ t('seller_panel.min_stock') }}</label><input v-model.number="local.minimum_stock" type="number" min="0" class="input" /></div>
    </div>
    <label class="inline-flex items-center gap-2">
      <input v-model="local.is_available" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
      <span class="text-sm">{{ t('seller_panel.available') }}</span>
    </label>
    <div>
      <label class="label">{{ t('seller_panel.images') }}</label>
      <ImageUploader :model-value="local.images" :endpoint="uploadEndpoint || '/seller/upload'" subdir="products" @update:model-value="onImages" />
    </div>
    <div class="flex gap-3 pt-1">
      <button type="submit" class="btn-primary" :disabled="submitting">{{ submitting ? t('common.loading') : t('seller_panel.save_product') }}</button>
      <slot name="actions" />
    </div>
  </form>
</template>
