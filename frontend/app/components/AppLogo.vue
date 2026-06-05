<script setup lang="ts">
const props = defineProps<{
  size?: 'sm' | 'md' | 'lg'
  inverted?: boolean
}>()
const settings = useSettingsStore()
const { resolve } = useImageUrl()
const sizeClass = computed(() => {
  switch (props.size || 'md') {
    case 'sm': return 'h-7'
    case 'lg': return 'h-12'
    default: return 'h-9'
  }
})
const textColor = computed(() => (props.inverted ? 'text-white' : 'text-ink-900'))
</script>

<template>
  <NuxtLink to="/" class="flex items-center gap-2 select-none">
    <img
      v-if="settings.logoUrl"
      :src="resolve(settings.logoUrl)"
      alt="AliStroy"
      :class="[sizeClass]"
    />
    <span v-else class="flex items-center justify-center rounded-lg bg-brand-500 text-white font-bold w-9 h-9">A</span>
    <span class="font-extrabold tracking-tight" :class="[textColor]">AliStroy</span>
  </NuxtLink>
</template>
