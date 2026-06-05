<script setup lang="ts">
const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const search = ref(typeof route.query.q === 'string' ? route.query.q : '')

const submit = () => {
  const q = search.value.trim()
  if (!q) return
  router.push({ path: '/search', query: { q } })
}

const mobileOpen = ref(false)

const dashboardPath = computed(() => {
  if (auth.isAdmin) return '/admin'
  if (auth.isSeller) return '/seller'
  return '/me'
})

const onLogout = async () => {
  await auth.logout()
  router.push('/')
}
</script>

<template>
  <header class="sticky top-0 z-30 bg-white border-b border-gray-100">
    <div class="container-page flex items-center gap-4 h-16">
      <AppLogo />

      <nav class="hidden lg:flex items-center gap-6 text-sm font-medium text-ink-700 ml-6">
        <NuxtLink to="/" class="hover:text-brand-600">{{ t('nav.home') }}</NuxtLink>
        <NuxtLink to="/products" class="hover:text-brand-600">{{ t('nav.catalog') }}</NuxtLink>
        <NuxtLink to="/categories" class="hover:text-brand-600">{{ t('nav.categories') }}</NuxtLink>
        <NuxtLink to="/sellers" class="hover:text-brand-600">{{ t('nav.sellers') }}</NuxtLink>
      </nav>

      <form class="hidden md:block flex-1 mx-4" @submit.prevent="submit">
        <div class="relative">
          <input
            v-model="search"
            type="search"
            class="input pl-10"
            :placeholder="t('common.search_placeholder')"
          />
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
          </span>
        </div>
      </form>

      <div class="flex items-center gap-2 ml-auto">
        <NuxtLink to="/favorites" class="hidden sm:inline-flex btn-ghost" :title="t('nav.favorites')">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
        </NuxtLink>

        <LanguageSwitcher />

        <template v-if="auth.isAuthenticated">
          <NuxtLink :to="dashboardPath" class="hidden sm:inline-flex btn-outline">
            {{ auth.isAdmin ? t('nav.admin_panel') : auth.isSeller ? t('nav.seller_panel') : t('nav.profile') }}
          </NuxtLink>
          <button class="btn-ghost" @click="onLogout">{{ t('nav.logout') }}</button>
        </template>
        <template v-else>
          <NuxtLink to="/login" class="btn-ghost">{{ t('nav.login') }}</NuxtLink>
          <NuxtLink to="/register" class="hidden sm:inline-flex btn-primary">{{ t('nav.register') }}</NuxtLink>
        </template>

        <button class="lg:hidden btn-ghost" @click="mobileOpen = !mobileOpen">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M3 12h18M3 18h18"/></svg>
        </button>
      </div>
    </div>

    <div v-if="mobileOpen" class="lg:hidden border-t border-gray-100 bg-white">
      <div class="container-page py-3 space-y-2">
        <form @submit.prevent="submit">
          <div class="relative">
            <input v-model="search" type="search" class="input pl-10" :placeholder="t('common.search_placeholder')" />
            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
            </span>
          </div>
        </form>
        <NuxtLink to="/" class="block py-2 text-ink-900">{{ t('nav.home') }}</NuxtLink>
        <NuxtLink to="/products" class="block py-2 text-ink-900">{{ t('nav.catalog') }}</NuxtLink>
        <NuxtLink to="/categories" class="block py-2 text-ink-900">{{ t('nav.categories') }}</NuxtLink>
        <NuxtLink to="/sellers" class="block py-2 text-ink-900">{{ t('nav.sellers') }}</NuxtLink>
        <NuxtLink to="/favorites" class="block py-2 text-ink-900">{{ t('nav.favorites') }}</NuxtLink>
      </div>
    </div>
  </header>
</template>
