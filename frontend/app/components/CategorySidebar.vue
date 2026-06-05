<script setup lang="ts">
import type { Category } from '~/types/api'

const props = defineProps<{ categories: Category[] }>()
const { categoryName } = useLocaleField()
const { resolve } = useImageUrl()
const activeId = ref<string | null>(null)
</script>

<template>
  <div class="card overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-100 font-semibold text-ink-900 flex items-center gap-2">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M3 12h18M3 18h18"/></svg>
      {{ $t('nav.categories') }}
    </div>
    <ul class="py-1 max-h-[460px] overflow-auto no-scrollbar">
      <li
        v-for="cat in categories"
        :key="cat.id"
        class="relative"
        @mouseenter="activeId = cat.id"
        @mouseleave="activeId = null"
      >
        <NuxtLink :to="`/categories/${cat.slug}`" class="flex items-center gap-3 px-4 py-2.5 hover:bg-brand-50 hover:text-brand-700 text-sm text-ink-900">
          <span class="w-6 h-6 flex items-center justify-center text-brand-500 shrink-0">
            <img v-if="cat.icon_url" :src="resolve(cat.icon_url)" alt="" class="w-5 h-5 object-contain" />
            <span v-else>•</span>
          </span>
          <span class="flex-1 truncate">{{ categoryName(cat) }}</span>
          <svg v-if="cat.subcategories && cat.subcategories.length" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </NuxtLink>

        <!-- Flyout subcategories (desktop) -->
        <div
          v-if="activeId === cat.id && cat.subcategories && cat.subcategories.length"
          class="hidden lg:block absolute left-full top-0 z-30 w-64 bg-white border border-gray-100 rounded-xl shadow-card p-3 ml-1"
        >
          <div class="text-xs uppercase text-gray-400 px-2 mb-1">{{ categoryName(cat) }}</div>
          <NuxtLink
            v-for="sub in cat.subcategories"
            :key="sub.id"
            :to="`/categories/${cat.slug}?subcategory=${sub.slug}`"
            class="block px-2 py-1.5 rounded-lg text-sm text-ink-900 hover:bg-brand-50 hover:text-brand-700"
          >
            {{ categoryName(sub) }}
          </NuxtLink>
        </div>
      </li>
    </ul>
  </div>
</template>
