<script setup lang="ts">
const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()

const onLogout = async () => {
  await auth.logout()
  router.push('/')
}

const items = computed(() => [
  { to: '/seller', label: t('seller_panel.dashboard'), exact: true },
  { to: '/seller/products', label: t('seller_panel.products') },
  { to: '/seller/products/new', label: t('seller_panel.add_product') },
  { to: '/seller/profile', label: t('seller_panel.profile') },
  { to: '/seller/stats', label: t('seller_panel.stats') },
])
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <AppHeader />
    <div class="flex-1">
      <div class="container-page py-6 grid grid-cols-1 lg:grid-cols-[230px_1fr] gap-6">
        <aside class="card p-2 h-fit lg:sticky lg:top-20">
          <div class="text-xs uppercase text-gray-400 px-3 py-2">{{ t('nav.seller_panel') }}</div>
          <NuxtLink v-for="item in items" :key="item.to" :to="item.to" class="block nav-link" active-class="nav-link-active" :exact="item.exact">{{ item.label }}</NuxtLink>
          <button class="w-full text-left nav-link text-red-600 hover:bg-red-50" @click="onLogout">{{ t('nav.logout') }}</button>
        </aside>
        <section class="min-w-0">
          <slot />
        </section>
      </div>
    </div>
  </div>
</template>
