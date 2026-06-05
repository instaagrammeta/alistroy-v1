<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()
const open = ref(false)
const root = ref<HTMLElement | null>(null)

const choose = async (code: string) => {
  await setLocale(code as 'tg' | 'ru')
  open.value = false
}
const onOutside = (e: MouseEvent) => {
  if (root.value && !root.value.contains(e.target as Node)) open.value = false
}
onMounted(() => document.addEventListener('click', onOutside))
onBeforeUnmount(() => document.removeEventListener('click', onOutside))
</script>

<template>
  <div ref="root" class="relative">
    <button type="button" class="inline-flex items-center gap-1.5 text-sm rounded-lg px-2.5 py-2 border border-gray-200 hover:bg-gray-50" @click="open = !open">
      <span class="font-semibold uppercase">{{ locale }}</span>
      <svg width="12" height="12" viewBox="0 0 20 20" fill="currentColor"><path d="M5.23 7.21a.75.75 0 011.06.02L10 11.06l3.71-3.83a.75.75 0 111.08 1.04l-4.25 4.39a.75.75 0 01-1.08 0L5.21 8.27a.75.75 0 01.02-1.06z"/></svg>
    </button>
    <div v-if="open" class="absolute right-0 mt-2 w-36 rounded-xl border border-gray-100 bg-white shadow-card z-50">
      <button v-for="l in (locales as any[])" :key="l.code" class="flex items-center justify-between w-full px-3 py-2 text-sm hover:bg-gray-50" :class="{ 'text-brand-600 font-medium': l.code === locale }" @click="choose(l.code)">
        <span>{{ l.name }}</span>
        <span class="uppercase text-xs text-gray-400">{{ l.code }}</span>
      </button>
    </div>
  </div>
</template>
