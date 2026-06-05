<script setup lang="ts">
const props = defineProps<{ page: number; totalPages: number }>()
const emit = defineEmits<{ (e: 'change', page: number): void }>()

const pages = computed(() => {
  const total = props.totalPages
  const cur = props.page
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  const arr: (number | '...')[] = [1]
  if (cur > 4) arr.push('...')
  const start = Math.max(2, cur - 1)
  const end = Math.min(total - 1, cur + 1)
  for (let i = start; i <= end; i++) arr.push(i)
  if (cur < total - 3) arr.push('...')
  arr.push(total)
  return arr
})
</script>

<template>
  <nav v-if="totalPages > 1" class="flex items-center justify-center gap-1 mt-8">
    <button
      class="btn-outline px-3 py-1.5 text-sm"
      :disabled="page <= 1"
      @click="emit('change', page - 1)"
    >‹</button>
    <template v-for="(p, i) in pages" :key="i">
      <span v-if="p === '...'" class="px-2 text-gray-400">…</span>
      <button
        v-else
        class="px-3 py-1.5 text-sm rounded-lg border"
        :class="p === page ? 'bg-brand-500 text-white border-brand-500' : 'border-gray-200 hover:bg-gray-50'"
        @click="emit('change', p as number)"
      >{{ p }}</button>
    </template>
    <button
      class="btn-outline px-3 py-1.5 text-sm"
      :disabled="page >= totalPages"
      @click="emit('change', page + 1)"
    >›</button>
  </nav>
</template>
