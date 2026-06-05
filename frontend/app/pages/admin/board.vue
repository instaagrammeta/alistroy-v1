<script setup lang="ts">
import type { BoardCategory } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t, locale } = useI18n()

const { data } = await useAsyncData('admin-board', () => useApi()<{ data: BoardCategory[] }>('/admin/board'))
const tree = computed<BoardCategory[]>(() => (data.value as any)?.data || [])

const search = ref('')
const zoom = ref(1)
const name = (n: { name_tj: string; name_ru: string }) => (locale.value === 'ru' ? n.name_ru : n.name_tj) || n.name_tj

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return tree.value
  return tree.value
    .map((c) => {
      const subs = c.subcategories.map((s) => ({ ...s, products: s.products.filter((p) => name(p).toLowerCase().includes(q)) }))
      const noSub = c.products.filter((p) => name(p).toLowerCase().includes(q))
      const matchCat = name(c).toLowerCase().includes(q)
      if (matchCat) return c
      return { ...c, subcategories: subs.filter((s) => s.products.length || name(s).toLowerCase().includes(q)), products: noSub }
    })
    .filter((c) => name(c).toLowerCase().includes(q) || c.subcategories.length || c.products.length)
})

useHead({ title: () => `${t('admin.board')} — Admin` })
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4 gap-3 flex-wrap">
      <h1 class="text-2xl font-bold text-ink-900">{{ t('admin.board') }}</h1>
      <div class="flex items-center gap-2">
        <input v-model="search" class="input max-w-xs" :placeholder="t('common.search')" />
        <div class="flex items-center gap-1">
          <button class="btn-outline btn-sm" @click="zoom = Math.max(0.5, zoom - 0.1)">−</button>
          <span class="text-sm w-12 text-center">{{ Math.round(zoom * 100) }}%</span>
          <button class="btn-outline btn-sm" @click="zoom = Math.min(1.5, zoom + 0.1)">+</button>
        </div>
      </div>
    </div>

    <div class="card p-6 overflow-auto">
      <div class="origin-top-left transition-transform" :style="{ transform: `scale(${zoom})` }">
        <div class="flex flex-col gap-6 min-w-max">
          <div v-for="cat in filtered" :key="cat.id" class="flex items-start gap-4">
            <div class="shrink-0 w-48 rounded-xl bg-brand-500 text-white p-3 font-semibold shadow-card">{{ name(cat) }}</div>
            <div class="flex flex-col gap-3 pt-1">
              <div v-for="sub in cat.subcategories" :key="sub.id" class="flex items-start gap-3">
                <div class="shrink-0 w-40 rounded-lg bg-brand-100 text-brand-800 p-2 text-sm font-medium">{{ name(sub) }}</div>
                <div class="flex flex-wrap gap-2 max-w-2xl">
                  <NuxtLink v-for="p in sub.products" :key="p.id" :to="`/admin/products/${p.id}`" class="badge bg-gray-100 text-gray-700 hover:bg-brand-50">{{ name(p) }}</NuxtLink>
                </div>
              </div>
              <div v-if="cat.products.length" class="flex flex-wrap gap-2 max-w-2xl">
                <NuxtLink v-for="p in cat.products" :key="p.id" :to="`/admin/products/${p.id}`" class="badge bg-gray-100 text-gray-700 hover:bg-brand-50">{{ name(p) }}</NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
