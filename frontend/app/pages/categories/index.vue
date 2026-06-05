<script setup lang="ts">
import type { Category } from '~/types/api'

const { t } = useI18n()
const { categoryName } = useLocaleField()
const { data } = await useAsyncData('all-categories', () => $fetch<{ data: Category[] }>('/categories', { baseURL: useApiBase() }))
const categories = computed<Category[]>(() => data.value?.data || [])
useHead({ title: () => `${t('nav.categories')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6">
    <h1 class="text-2xl md:text-3xl font-bold text-ink-900 mb-6">{{ t('nav.categories') }}</h1>
    <EmptyState v-if="!categories.length" />
    <div v-else class="space-y-6">
      <div v-for="cat in categories" :key="cat.id" class="card p-5">
        <div class="flex items-center justify-between mb-3">
          <NuxtLink :to="`/categories/${cat.slug}`" class="text-lg font-bold text-ink-900 hover:text-brand-600">{{ categoryName(cat) }}</NuxtLink>
        </div>
        <div v-if="cat.subcategories && cat.subcategories.length" class="flex flex-wrap gap-2">
          <NuxtLink v-for="sub in cat.subcategories" :key="sub.id" :to="`/categories/${cat.slug}?subcategory=${sub.slug}`" class="badge bg-gray-100 text-gray-700 hover:bg-brand-50 hover:text-brand-700">{{ categoryName(sub) }}</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
