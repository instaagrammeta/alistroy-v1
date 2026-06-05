<script setup lang="ts">
import type { Favorite, PaginatedResponse } from '~/types/api'

definePageMeta({ middleware: 'auth' })

const { t } = useI18n()
const api = useApi()

const page = ref(1)
const pageSize = 24

const { data, refresh } = await useAsyncData(
  'my-favorites',
  () =>
    api<PaginatedResponse<Favorite>>(`/favorites`, {
      query: { page: page.value, page_size: pageSize },
    }),
  { watch: [page] }
)

const favorites = computed<Favorite[]>(() => data.value?.data || [])
const pagination = computed(() => data.value?.pagination || { page: 1, page_size: pageSize, total: 0, total_pages: 1 })

const onPage = (p: number) => (page.value = p)

useHead({ title: () => `${t('favorites.title')} — AliStroy` })
</script>

<template>
  <div class="container-page py-8">
    <h1 class="text-3xl font-bold text-ink-900 mb-6">{{ t('favorites.title') }}</h1>
    <EmptyState v-if="!favorites.length" :title="t('favorites.empty')" />
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      <ProductCard v-for="f in favorites" :key="f.id" :product="(f.product as any)" />
    </div>
    <Pagination :page="pagination.page" :total-pages="pagination.total_pages" @change="onPage" />
  </div>
</template>
