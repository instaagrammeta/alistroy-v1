<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const { t, locale } = useI18n()
const router = useRouter()
const auth = useAuthStore()

const name = ref('')
const phone = ref('')
const email = ref('')
const password = ref('')
const address = ref('')
const city = ref('')
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.register({
      name: name.value, phone: phone.value, email: email.value || undefined,
      password: password.value, address: address.value, city: city.value,
      locale: locale.value as 'tg' | 'ru',
    })
    router.push('/')
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Registration failed'
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

useHead({ title: () => `${t('auth.register_title')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-ink-900 text-center mb-6">{{ t('auth.register_title') }}</h1>
    <form class="space-y-3" @submit.prevent="submit">
      <div><label class="label">{{ t('auth.name') }}</label><input v-model="name" required class="input" /></div>
      <div><label class="label">{{ t('auth.phone') }}</label><input v-model="phone" required class="input" placeholder="+992..." /></div>
      <div><label class="label">{{ t('auth.email') }} ({{ t('common.optional') }})</label><input v-model="email" type="email" class="input" /></div>
      <div><label class="label">{{ t('auth.address') }}</label><input v-model="address" class="input" /></div>
      <div><label class="label">{{ t('auth.city') }}</label><input v-model="city" class="input" /></div>
      <div><label class="label">{{ t('auth.password') }}</label><input v-model="password" type="password" required minlength="8" class="input" /></div>
      <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">{{ loading ? t('common.loading') : t('auth.submit_register') }}</button>
    </form>
    <button class="btn-outline w-full mt-3" @click="googleLogin">{{ t('auth.google') }}</button>
    <div class="mt-4 text-sm text-center text-gray-500">
      {{ t('auth.have_account') }} <NuxtLink to="/login" class="text-brand-600 hover:text-brand-700">{{ t('nav.login') }}</NuxtLink>
    </div>
  </div>
</template>
