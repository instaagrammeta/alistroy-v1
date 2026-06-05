<script setup lang="ts">
definePageMeta({ layout: 'auth' })
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()

const error = ref('')
const needsProfile = ref(false)
const phone = ref('')
const address = ref('')
const saving = ref(false)

onMounted(async () => {
  const code = route.query.code as string
  if (!code) {
    error.value = 'Missing code'
    return
  }
  try {
    const res = await $fetch<{ data: { tokens: any; needs_profile: boolean } }>('/auth/google/callback', {
      baseURL: useApiBase(), method: 'POST', body: { code },
    })
    auth.setTokens(res.data.tokens)
    if (res.data.needs_profile) {
      needsProfile.value = true
    } else {
      router.push('/')
    }
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Google login failed'
  }
})

const completeProfile = async () => {
  saving.value = true
  try {
    await useApi()('/me', { method: 'PATCH', body: { phone: phone.value, address: address.value } })
    await auth.fetchMe()
    router.push('/')
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Error'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <div v-if="error" class="text-sm text-red-600 text-center">{{ error }}</div>
    <div v-else-if="needsProfile">
      <h1 class="text-xl font-bold text-ink-900 text-center mb-2">{{ t('auth.complete_profile') }}</h1>
      <p class="text-sm text-gray-500 text-center mb-4">{{ t('auth.complete_profile_text') }}</p>
      <div class="space-y-3">
        <div><label class="label">{{ t('auth.phone') }}</label><input v-model="phone" class="input" placeholder="+992..." /></div>
        <div><label class="label">{{ t('auth.address') }}</label><input v-model="address" class="input" /></div>
        <button class="btn-primary w-full" :disabled="saving" @click="completeProfile">{{ t('common.save') }}</button>
      </div>
    </div>
    <div v-else class="text-center text-gray-500">{{ t('common.loading') }}</div>
  </div>
</template>
