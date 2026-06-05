<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    const redirect = (route.query.redirect as string) || (auth.isAdmin ? '/admin' : auth.isSeller ? '/seller' : '/')
    router.push(redirect)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
useHead({ title: () => `${t('auth.login_title')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-ink-900 text-center mb-6">{{ t('auth.login_title') }}</h1>

    <form class="space-y-4" @submit.prevent="submit">
      <div>
        <label class="label">{{ t('auth.email') }}</label>
        <input v-model="email" type="email" required class="input" />
      </div>
      <div>
        <label class="label">{{ t('auth.password') }}</label>
        <input v-model="password" type="password" required class="input" />
      </div>
      <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? t('common.loading') : t('auth.submit_login') }}
      </button>
    </form>

    <div class="mt-4 text-sm text-center">
      <NuxtLink to="/forgot-password" class="text-gray-500 hover:text-brand-600">{{ t('auth.forgot') }}</NuxtLink>
    </div>
    <div class="mt-2 text-sm text-center text-gray-500">
      {{ t('auth.no_account') }}
      <NuxtLink to="/register" class="text-brand-600 hover:text-brand-700">{{ t('nav.register') }}</NuxtLink>
    </div>
  </div>
</template>
