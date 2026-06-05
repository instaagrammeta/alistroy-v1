<script setup lang="ts">
import type { Category } from '~/types/api'

defineProps<{ category: Category }>()
const { categoryName } = useLocaleField()
const { resolve } = useImageUrl()
</script>

<template>
  <NuxtLink :to="`/categories/${category.slug}`" class="group block card p-4 hover:shadow-card transition-shadow text-center">
    <div class="mx-auto w-16 h-16 rounded-xl bg-brand-50 flex items-center justify-center text-brand-600 overflow-hidden">
      <img v-if="category.icon_url" :src="resolve(category.icon_url)" alt="" class="w-12 h-12 object-contain" />
      <span v-else class="text-2xl font-bold">{{ categoryName(category).charAt(0).toUpperCase() }}</span>
    </div>
    <div class="mt-3 font-medium text-sm text-ink-900 group-hover:text-brand-600 line-clamp-2">{{ categoryName(category) }}</div>
  </NuxtLink>
</template>
