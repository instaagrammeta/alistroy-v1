<script setup lang="ts">
import type { PaginatedResponse, Seller } from '~/types/api'

const route = useRoute()
const { t } = useI18n()
const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const page = ref(1)
const pageSize = 24

const { data } = await useAsyncData('sellers', () =>
  $fetch<PaginatedResponse<Seller>>('/sellers', { baseURL: useApiBase(), query: { page: page.value, page_size: pageSize, q: q.value || undefined } }),
  { watch: [page, q] }
)
const sellers = computed<Seller[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

useHead({ title: () => `${t('nav.sellers')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6">
    <h1 class="text-2xl md:text-3xl font-bold text-ink-900 mb-5">{{ t('nav.sellers') }}</h1>
    <input v-model="q" type="search" class="input max-w-md mb-6" :placeholder="t('common.search')" />
    <EmptyState v-if="!sellers.length" />
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <SellerCard v-for="s in sellers" :key="s.id" :seller="s" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="(p) => (page = p)" />
  </div>
</template>
