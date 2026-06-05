<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const role = ref<'customer' | 'seller'>(
  route.query.role === 'seller' ? 'seller' : 'customer'
)
const name = ref('')
const email = ref('')
const password = ref('')
const phone = ref('')
const sellerName = ref('')
const city = ref('')
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.register({
      email: email.value,
      password: password.value,
      name: name.value,
      phone: phone.value,
      role: role.value,
      seller_name: role.value === 'seller' ? sellerName.value : undefined,
      city: role.value === 'seller' ? city.value : undefined,
      locale: locale.value as 'tg' | 'ru',
    })
    router.push(role.value === 'seller' ? '/seller' : '/')
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
useHead({ title: () => `${t('auth.register_title')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-ink-900 text-center mb-6">{{ t('auth.register_title') }}</h1>

    <div class="grid grid-cols-2 gap-2 mb-4">
      <button
        type="button"
        class="px-3 py-2 rounded-lg text-sm font-medium border"
        :class="role === 'customer' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'"
        @click="role = 'customer'"
      >{{ t('auth.as_customer') }}</button>
      <button
        type="button"
        class="px-3 py-2 rounded-lg text-sm font-medium border"
        :class="role === 'seller' ? 'bg-brand-50 text-brand-700 border-brand-200' : 'border-gray-200 text-gray-600'"
        @click="role = 'seller'"
      >{{ t('auth.as_seller') }}</button>
    </div>

    <form class="space-y-3" @submit.prevent="submit">
      <div>
        <label class="label">{{ t('auth.name') }}</label>
        <input v-model="name" required class="input" />
      </div>
      <div>
        <label class="label">{{ t('auth.email') }}</label>
        <input v-model="email" type="email" required class="input" />
      </div>
      <div>
        <label class="label">{{ t('auth.phone') }}</label>
        <input v-model="phone" class="input" />
      </div>
      <div>
        <label class="label">{{ t('auth.password') }}</label>
        <input v-model="password" type="password" required minlength="8" class="input" />
      </div>

      <template v-if="role === 'seller'">
        <div>
          <label class="label">{{ t('auth.seller_name') }}</label>
          <input v-model="sellerName" required class="input" />
        </div>
        <div>
          <label class="label">{{ t('auth.city') }}</label>
          <input v-model="city" class="input" />
        </div>
      </template>

      <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? t('common.loading') : t('auth.submit_register') }}
      </button>
    </form>

    <div class="mt-4 text-sm text-center text-gray-500">
      {{ t('auth.have_account') }}
      <NuxtLink to="/login" class="text-brand-600 hover:text-brand-700">{{ t('nav.login') }}</NuxtLink>
    </div>
  </div>
</template>
