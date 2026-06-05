<script setup lang="ts">
const props = defineProps<{ size?: 'sm' | 'md' | 'lg'; inverted?: boolean }>()
const settings = useSettingsStore()
const { resolve } = useImageUrl()
const sizeClass = computed(() => (props.size === 'lg' ? 'h-12' : props.size === 'sm' ? 'h-7' : 'h-9'))
const textColor = computed(() => (props.inverted ? 'text-white' : 'text-ink-900'))
</script>

<template>
  <NuxtLink to="/" class="flex items-center gap-2 select-none shrink-0">
    <img v-if="settings.logoUrl" :src="resolve(settings.logoUrl)" alt="AliStroy" :class="sizeClass" />
    <span v-else class="flex items-center justify-center rounded-lg bg-brand-500 text-white font-extrabold w-9 h-9">A</span>
    <span class="font-extrabold tracking-tight text-lg" :class="textColor">AliStroy</span>
  </NuxtLink>
</template>
