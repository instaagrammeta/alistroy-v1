<script setup lang="ts">
import type { Category } from '~/types/api'

interface ImageItem { url: string; alt?: string }

interface FormState {
  category_id: string
  title_tj: string
  title_ru: string
  description_tj: string
  description_ru: string
  price: number
  currency: string
  unit: string
  sku: string
  stock_quantity: number
  is_available: boolean
  images: ImageItem[]
}

const props = defineProps<{
  modelValue: FormState
  uploadEndpoint?: string
  submitting?: boolean
  submitLabel?: string
}>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: FormState): void
  (e: 'submit'): void
}>()

const config = useRuntimeConfig()
const { t } = useI18n()
const { categoryTitle } = useLocaleField()

const { data: catsRes } = await useAsyncData('form-cats', () =>
  $fetch<{ data: Category[] }>('/categories', { baseURL: config.public.apiBase })
)
const cats = computed<Category[]>(() => catsRes.value?.data || [])

const local = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const onChangeImages = (next: ImageItem[]) => {
  emit('update:modelValue', { ...local.value, images: next })
}
</script>

<template>
  <form class="space-y-5" @submit.prevent="emit('submit')">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label">{{ t('seller_panel.title_tj') }}</label>
        <input v-model="local.title_tj" class="input" required minlength="2" />
      </div>
      <div>
        <label class="label">{{ t('seller_panel.title_ru') }}</label>
        <input v-model="local.title_ru" class="input" required minlength="2" />
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label">{{ t('seller_panel.desc_tj') }}</label>
        <textarea v-model="local.description_tj" class="textarea" rows="5"></textarea>
      </div>
      <div>
        <label class="label">{{ t('seller_panel.desc_ru') }}</label>
        <textarea v-model="local.description_ru" class="textarea" rows="5"></textarea>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="label">{{ t('seller_panel.category') }}</label>
        <select v-model="local.category_id" class="select" required>
          <option value="" disabled>{{ t('common.required') }}</option>
          <option v-for="c in cats" :key="c.id" :value="c.id">{{ categoryTitle(c) }}</option>
        </select>
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="label">{{ t('seller_panel.price') }}</label>
          <input v-model.number="local.price" type="number" min="0" step="0.01" class="input" required />
        </div>
        <div>
          <label class="label">{{ t('seller_panel.currency') }}</label>
          <input v-model="local.currency" class="input" maxlength="8" />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div>
        <label class="label">{{ t('seller_panel.sku') }}</label>
        <input v-model="local.sku" class="input" maxlength="64" />
      </div>
      <div>
        <label class="label">{{ t('seller_panel.unit') }}</label>
        <input v-model="local.unit" class="input" maxlength="20" />
      </div>
      <div>
        <label class="label">{{ t('seller_panel.stock') }}</label>
        <input v-model.number="local.stock_quantity" type="number" min="0" class="input" />
      </div>
    </div>

    <label class="inline-flex items-center gap-2">
      <input v-model="local.is_available" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
      <span class="text-sm">{{ t('seller_panel.is_available') }}</span>
    </label>

    <div>
      <label class="label">{{ t('seller_panel.images') }}</label>
      <ImageUploader
        :model-value="local.images"
        :endpoint="uploadEndpoint || '/seller/upload'"
        subdir="products"
        @update:model-value="onChangeImages"
      />
    </div>

    <div class="flex items-center gap-3 pt-2">
      <button type="submit" class="btn-primary" :disabled="submitting">
        {{ submitting ? t('common.loading') : (submitLabel || t('seller_panel.save_product')) }}
      </button>
      <NuxtLink to="/seller/products" class="btn-outline">{{ t('common.cancel') }}</NuxtLink>
    </div>
  </form>
</template>
