<script setup lang="ts">
defineProps<{ title?: string; description?: string }>()
</script>

<template>
  <div class="text-center py-12 px-6">
    <div class="mx-auto w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center text-gray-400">
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="11" cy="11" r="8"/>
        <path d="m21 21-4.3-4.3"/>
      </svg>
    </div>
    <h3 class="mt-4 text-lg font-semibold text-ink-900">{{ title || $t('common.no_results') }}</h3>
    <p v-if="description" class="mt-1 text-sm text-gray-500">{{ description }}</p>
    <slot />
  </div>
</template>
