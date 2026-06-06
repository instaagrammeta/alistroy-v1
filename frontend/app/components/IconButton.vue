<script setup lang="ts">
defineProps<{ title?: string; variant?: 'edit' | 'delete' | 'view' | 'toggle-on' | 'toggle-off' | 'plus' | 'check' | 'x' | 'download' }>()
</script>

<template>
  <button
    type="button"
    :title="title"
    class="inline-flex items-center justify-center w-9 h-9 rounded-lg transition-colors"
    :class="{
      'text-brand-600 hover:bg-brand-50': variant === 'edit',
      'text-red-600 hover:bg-red-50': variant === 'delete' || variant === 'x',
      'text-gray-600 hover:bg-gray-100': variant === 'view' || variant === 'download' || !variant,
      'text-emerald-600 hover:bg-emerald-50': variant === 'toggle-on' || variant === 'check' || variant === 'plus',
      'text-gray-400 hover:bg-gray-100': variant === 'toggle-off',
    }"
  >
    <svg v-if="variant === 'edit'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
    <svg v-else-if="variant === 'delete'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-2 14a2 2 0 0 1-2 2H9a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/></svg>
    <svg v-else-if="variant === 'view'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
    <svg v-else-if="variant === 'toggle-on'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="5" width="22" height="14" rx="7" ry="7" fill="currentColor"/><circle cx="16" cy="12" r="3" fill="white" stroke="none"/></svg>
    <svg v-else-if="variant === 'toggle-off'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="5" width="22" height="14" rx="7" ry="7"/><circle cx="8" cy="12" r="3" fill="currentColor" stroke="none"/></svg>
    <svg v-else-if="variant === 'plus'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
    <svg v-else-if="variant === 'check'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
    <svg v-else-if="variant === 'x'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
    <svg v-else-if="variant === 'download'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
    <slot />
  </button>
</template>
