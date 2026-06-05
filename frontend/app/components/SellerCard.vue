<script setup lang="ts">
import type { Seller } from '~/types/api'

defineProps<{ seller: Seller }>()
const { resolve } = useImageUrl()
</script>

<template>
  <NuxtLink :to="`/sellers/${seller.slug}`" class="group block card p-5 hover:shadow-card transition-shadow">
    <div class="flex items-center gap-4">
      <div class="w-14 h-14 rounded-xl bg-gray-100 flex items-center justify-center overflow-hidden shrink-0">
        <img v-if="seller.logo_url" :src="resolve(seller.logo_url)" :alt="seller.full_name" class="w-full h-full object-cover" />
        <span v-else class="text-brand-600 text-xl font-bold">{{ (seller.market_name || seller.full_name).charAt(0).toUpperCase() }}</span>
      </div>
      <div class="min-w-0">
        <div class="font-semibold text-ink-900 truncate group-hover:text-brand-600">{{ seller.market_name || seller.company_name || seller.full_name }}</div>
        <div class="text-xs text-gray-400 truncate">{{ seller.city || ' ' }}</div>
      </div>
    </div>
  </NuxtLink>
</template>
