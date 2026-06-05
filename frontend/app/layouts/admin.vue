<script setup lang="ts">
const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()

const onLogout = async () => {
  await auth.logout()
  router.push('/')
}

const items = computed(() => [
  { to: '/admin', label: t('admin.dashboard') },
  { to: '/admin/products', label: t('admin.products') },
  { to: '/admin/sellers', label: t('admin.sellers') },
  { to: '/admin/users', label: t('admin.users') },
  { to: '/admin/categories', label: t('admin.categories') },
  { to: '/admin/reviews', label: t('admin.reviews') },
  { to: '/admin/settings', label: t('admin.settings') },
])
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <AppHeader />
    <div class="flex-1">
      <div class="container-page py-6 grid grid-cols-1 lg:grid-cols-[240px_1fr] gap-6">
        <aside class="card p-3 h-fit">
          <div class="text-xs uppercase text-gray-400 px-3 py-2">{{ t('nav.admin_panel') }}</div>
          <NuxtLink
            v-for="item in items"
            :key="item.to"
            :to="item.to"
            class="block px-3 py-2 rounded-lg text-sm font-medium text-ink-900 hover:bg-gray-50"
            active-class="bg-brand-50 text-brand-700"
            :exact="item.to === '/admin'"
          >
            {{ item.label }}
          </NuxtLink>
          <button class="w-full text-left px-3 py-2 rounded-lg text-sm text-red-600 hover:bg-red-50" @click="onLogout">
            {{ t('nav.logout') }}
          </button>
        </aside>
        <section>
          <slot />
        </section>
      </div>
    </div>
    <AppFooter />
  </div>
</template>
