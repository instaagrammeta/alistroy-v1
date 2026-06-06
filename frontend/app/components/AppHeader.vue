<script setup lang="ts">
const auth = useAuthStore()
const cart = useCartStore()
const favorites = useFavoritesStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const search = ref(typeof route.query.q === 'string' ? route.query.q : '')
const profileOpen = ref(false)
const root = ref<HTMLElement | null>(null)

const submit = () => {
  const q = search.value.trim()
  if (q) router.push({ path: '/search', query: { q } })
}

const dashboardPath = computed(() => {
  if (auth.isAdmin) return '/admin'
  if (auth.isSeller) return '/seller'
  if (auth.isDriver) return '/driver'
  return '/account'
})

const onLogout = async () => {
  await auth.logout()
  profileOpen.value = false
  router.push('/')
}

const onOutside = (e: MouseEvent) => {
  if (root.value && !root.value.contains(e.target as Node)) profileOpen.value = false
}
onMounted(() => document.addEventListener('click', onOutside))
onBeforeUnmount(() => document.removeEventListener('click', onOutside))
</script>

<template>
  <header class="sticky top-0 z-40 bg-white border-b border-gray-100">
    <div class="container-page flex items-center gap-3 h-16">
      <AppLogo />

      <form class="flex-1 max-w-2xl mx-auto" @submit.prevent="submit">
        <div class="relative">
          <input
            v-model="search"
            type="search"
            class="input pl-11 pr-4 h-11 bg-gray-50 border-gray-100 focus:bg-white"
            :placeholder="t('common.search_placeholder')"
          />
          <button type="submit" class="absolute left-0 top-0 h-11 w-11 flex items-center justify-center text-gray-400 hover:text-brand-600">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
            </svg>
          </button>
        </div>
      </form>

      <div ref="root" class="flex items-center gap-1.5">
        <NuxtLink to="/favorites" class="hidden sm:inline-flex relative btn-ghost p-2.5" :title="t('nav.favorites')">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
          <span v-if="favorites.ids.length > 0" class="absolute -top-0.5 -right-0.5 bg-brand-500 text-white text-[10px] rounded-full min-w-[16px] h-4 px-1 flex items-center justify-center">{{ favorites.ids.length }}</span>
        </NuxtLink>

        <NuxtLink v-if="auth.isCustomer || !auth.isAuthenticated" to="/cart" class="relative btn-ghost p-2.5" :title="t('nav.cart')">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
          <span v-if="cart.count > 0" class="absolute -top-0.5 -right-0.5 bg-brand-500 text-white text-[10px] rounded-full min-w-[16px] h-4 px-1 flex items-center justify-center">{{ cart.count }}</span>
        </NuxtLink>

        <LanguageSwitcher />

        <!-- Profile dropdown -->
        <div class="relative">
          <button class="btn-ghost p-2.5" :title="t('nav.profile')" @click.stop="profileOpen = !profileOpen">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
          </button>
          <div v-if="profileOpen" class="absolute right-0 mt-2 w-56 rounded-xl border border-gray-100 bg-white shadow-card z-50 py-1.5">
            <template v-if="auth.isAuthenticated">
              <div class="px-4 py-2 border-b border-gray-100">
                <div class="text-sm font-semibold text-ink-900 truncate">{{ auth.user?.name }}</div>
                <div class="text-xs text-gray-400 truncate">{{ auth.user?.phone || auth.user?.email }}</div>
              </div>
              <NuxtLink :to="dashboardPath" class="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-50" @click="profileOpen = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
                {{ auth.isAdmin ? t('nav.admin_panel') : auth.isSeller ? t('nav.seller_panel') : auth.isDriver ? t('nav.driver_panel') : t('nav.profile') }}
              </NuxtLink>
              <NuxtLink v-if="auth.isCustomer" to="/account/orders" class="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-50" @click="profileOpen = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                {{ t('nav.orders') }}
              </NuxtLink>
              <NuxtLink v-if="auth.isCustomer" to="/account/chat" class="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-50" @click="profileOpen = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
                {{ t('admin.chat') }}
              </NuxtLink>
              <button class="flex items-center gap-2 w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50" @click="onLogout">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                {{ t('nav.logout') }}
              </button>
            </template>
            <template v-else>
              <NuxtLink to="/login" class="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-50" @click="profileOpen = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                {{ t('nav.login') }}
              </NuxtLink>
              <NuxtLink to="/register" class="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-50" @click="profileOpen = false">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="8.5" cy="7" r="4"/><line x1="20" y1="8" x2="20" y2="14"/><line x1="23" y1="11" x2="17" y2="11"/></svg>
                {{ t('nav.register') }}
              </NuxtLink>
            </template>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
