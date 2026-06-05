<script setup lang="ts">
import type { Category } from '~/types/api'

const props = defineProps<{ category: Category }>()
const { categoryTitle } = useLocaleField()
const { resolve } = useImageUrl()
</script>

<template>
  <NuxtLink :to="`/categories/${category.slug}`" class="group block">
    <div class="card p-5 flex items-center gap-4 hover:shadow-card transition-shadow">
      <div class="w-12 h-12 rounded-lg bg-brand-50 flex items-center justify-center text-brand-600 text-lg font-bold shrink-0">
        <img v-if="category.icon_url" :src="resolve(category.icon_url)" alt="" class="w-8 h-8 object-contain" />
        <span v-else>{{ categoryTitle(category).charAt(0).toUpperCase() }}</span>
      </div>
      <div class="min-w-0">
        <div class="font-semibold text-ink-900 truncate group-hover:text-brand-600">{{ categoryTitle(category) }}</div>
        <div class="text-xs text-gray-400 truncate">{{ category.slug }}</div>
      </div>
    </div>
  </NuxtLink>
</template>
