<script setup lang="ts">
import type { Category } from '~/types/api'

const config = useRuntimeConfig()
const { t } = useI18n()

const { data } = await useAsyncData('all-categories', () =>
  $fetch<{ data: Category[] }>('/categories', { baseURL: useApiBase() })
)

const categories = computed<Category[]>(() => data.value?.data || [])

useHead({ title: () => `${t('nav.categories')} — AliStroy` })
</script>

<template>
  <div class="container-page py-8">
    <h1 class="text-3xl font-bold text-ink-900 mb-6">{{ t('nav.categories') }}</h1>
    <EmptyState v-if="!categories.length" />
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <CategoryCard v-for="c in categories" :key="c.id" :category="c" />
    </div>
  </div>
</template>
