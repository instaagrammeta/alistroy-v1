<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const identifier = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const redirectByRole = () => {
  const r = (route.query.redirect as string)
  if (r) return r
  if (auth.isAdmin) return '/admin'
  if (auth.isSeller) return '/seller'
  if (auth.isDriver) return '/driver'
  return '/'
}

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.login(identifier.value, password.value)
    const cart = useCartStore(); const fav = useFavoritesStore()
    await Promise.all([cart.load(), fav.load()])
    router.push(redirectByRole())
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

const googleLogin = async () => {
  try {
    const res = await $fetch<{ data: { url: string } }>('/auth/google/url', { baseURL: useApiBase() })
    if (res.data?.url) window.location.href = res.data.url
  } catch {
    error.value = 'Google login unavailable'
  }
}

useHead({ title: () => `${t('auth.login_title')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-ink-900 text-center mb-6">{{ t('auth.login_title') }}</h1>
    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label class="label">{{ t('auth.identifier') }}</label>
        <input v-model="identifier" required class="input" autocomplete="username" />
      </div>
      <div>
        <label class="label">{{ t('auth.password') }}</label>
        <input v-model="password" type="password" required class="input" autocomplete="current-password" />
      </div>
      <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">{{ loading ? t('common.loading') : t('auth.submit_login') }}</button>
    </form>
    <button class="btn-outline w-full mt-3" @click="googleLogin">{{ t('auth.google') }}</button>
    <div class="mt-4 text-sm text-center text-gray-500">
      {{ t('auth.no_account') }} <NuxtLink to="/register" class="text-brand-600 hover:text-brand-700">{{ t('nav.register') }}</NuxtLink>
    </div>
  </div>
</template>
