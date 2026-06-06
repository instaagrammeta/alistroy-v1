<script setup lang="ts">
const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()

const onLogout = async () => {
  await auth.logout()
  router.push('/')
}

interface NavItem { to: string; label: string; icon: string; exact?: boolean }
const items = computed<NavItem[]>(() => [
  { to: '/admin', label: t('admin.dashboard'), icon: 'dashboard', exact: true },
  { to: '/admin/orders', label: t('admin.orders'), icon: 'cart' },
  { to: '/admin/products', label: t('admin.products'), icon: 'box' },
  { to: '/admin/banners', label: t('admin.banners'), icon: 'image' },
  { to: '/admin/categories', label: t('admin.categories'), icon: 'grid' },
  { to: '/admin/brands', label: t('admin.brands'), icon: 'tag' },
  { to: '/admin/sellers', label: t('admin.sellers'), icon: 'store' },
  { to: '/admin/customers', label: t('admin.customers'), icon: 'users' },
  { to: '/admin/drivers', label: t('admin.drivers'), icon: 'truck' },
  { to: '/admin/reviews', label: t('admin.reviews'), icon: 'star' },
  { to: '/admin/chat', label: t('admin.chat'), icon: 'chat' },
  { to: '/admin/board', label: t('admin.board'), icon: 'tree' },
  { to: '/admin/reports', label: t('admin.reports'), icon: 'chart' },
  { to: '/admin/settings', label: t('admin.settings'), icon: 'settings' },
])
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <AppHeader />
    <div class="flex-1">
      <div class="container-page py-6 grid grid-cols-1 lg:grid-cols-[230px_1fr] gap-6">
        <aside class="card p-2 h-fit lg:sticky lg:top-20">
          <div class="text-xs uppercase text-gray-400 px-3 py-2">{{ t('nav.admin_panel') }}</div>
          <NuxtLink v-for="item in items" :key="item.to" :to="item.to" class="flex items-center gap-3 nav-link" active-class="nav-link-active" :exact="item.exact">
            <svg v-if="item.icon === 'dashboard'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            <svg v-else-if="item.icon === 'cart'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
            <svg v-else-if="item.icon === 'box'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
            <svg v-else-if="item.icon === 'image'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="9" cy="9" r="2"/><path d="m21 15-3.1-3.1a2 2 0 0 0-2.81.01L6 21"/></svg>
            <svg v-else-if="item.icon === 'grid'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            <svg v-else-if="item.icon === 'tag'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41 13.42 20.58a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
            <svg v-else-if="item.icon === 'store'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 3h18l-2 6H5z"/><path d="M5 9v12h14V9"/><path d="M9 21V13h6v8"/></svg>
            <svg v-else-if="item.icon === 'users'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            <svg v-else-if="item.icon === 'truck'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
            <svg v-else-if="item.icon === 'star'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
            <svg v-else-if="item.icon === 'chat'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
            <svg v-else-if="item.icon === 'tree'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="6" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="12" r="3"/><line x1="9" y1="6" x2="15" y2="12"/><line x1="9" y1="18" x2="15" y2="12"/></svg>
            <svg v-else-if="item.icon === 'chart'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>
            <svg v-else-if="item.icon === 'settings'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            <span class="text-sm">{{ item.label }}</span>
          </NuxtLink>
          <button class="w-full flex items-center gap-3 nav-link text-red-600 hover:bg-red-50" @click="onLogout">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
            <span class="text-sm">{{ t('nav.logout') }}</span>
          </button>
        </aside>
        <section class="min-w-0">
          <slot />
        </section>
      </div>
    </div>
  </div>
</template>
