<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const { t } = useI18n()
const config = useRuntimeConfig()

const email = ref('')
const resetToken = ref('')
const newPassword = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)
const stage = ref<'request' | 'reset'>('request')

const requestReset = async () => {
  error.value = ''
  success.value = ''
  loading.value = true
  try {
    const res: any = await $fetch('/auth/forgot-password', {
      baseURL: useApiBase(),
      method: 'POST',
      body: { email: email.value },
    })
    // For dev/staging, the API returns the token directly so admins can deliver it.
    if (res?.data?.reset_token) {
      resetToken.value = res.data.reset_token
      success.value = 'Token issued — paste it below to set a new password.'
      stage.value = 'reset'
    } else {
      success.value = 'If an account exists, a reset has been initiated.'
    }
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Request failed'
  } finally {
    loading.value = false
  }
}

const doReset = async () => {
  error.value = ''
  loading.value = true
  try {
    await $fetch('/auth/reset-password', {
      baseURL: useApiBase(),
      method: 'POST',
      body: { token: resetToken.value, new_password: newPassword.value },
    })
    success.value = t('auth.reset_done')
    stage.value = 'request'
    email.value = ''
    resetToken.value = ''
    newPassword.value = ''
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Reset failed'
  } finally {
    loading.value = false
  }
}

useHead({ title: () => `${t('auth.reset_request')} — AliStroy` })
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-ink-900 text-center mb-6">{{ t('auth.reset_request') }}</h1>

    <form v-if="stage === 'request'" class="space-y-3" @submit.prevent="requestReset">
      <div>
        <label class="label">{{ t('auth.email') }}</label>
        <input v-model="email" type="email" required class="input" />
      </div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? t('common.loading') : t('common.submit') }}
      </button>
    </form>

    <form v-else class="space-y-3" @submit.prevent="doReset">
      <div>
        <label class="label">{{ t('auth.reset_token') }}</label>
        <input v-model="resetToken" required class="input" />
      </div>
      <div>
        <label class="label">{{ t('auth.new_password') }}</label>
        <input v-model="newPassword" type="password" required minlength="8" class="input" />
      </div>
      <button type="submit" class="btn-primary w-full" :disabled="loading">{{ t('common.save') }}</button>
    </form>

    <div v-if="error" class="mt-3 text-sm text-red-600">{{ error }}</div>
    <div v-if="success" class="mt-3 text-sm text-emerald-600">{{ success }}</div>

    <div class="mt-6 text-center text-sm text-gray-500">
      <NuxtLink to="/login" class="text-brand-600 hover:text-brand-700">← {{ t('nav.login') }}</NuxtLink>
    </div>
  </div>
</template>
