<script setup lang="ts">
definePageMeta({ middleware: 'auth' })
const auth = useAuthStore()
const api = useApi()
const { t } = useI18n()

const name = ref(auth.user?.name || '')
const phone = ref(auth.user?.phone || '')
const message = ref('')
const error = ref('')
const loading = ref(false)

const save = async () => {
  error.value = ''
  message.value = ''
  loading.value = true
  try {
    const res = await api<{ data: any }>('/me', {
      method: 'PATCH',
      body: { name: name.value, phone: phone.value },
    })
    auth.setUser(res.data)
    message.value = t('common.save')
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Save failed'
  } finally {
    loading.value = false
  }
}

const oldPassword = ref('')
const newPassword = ref('')
const pwLoading = ref(false)
const pwMessage = ref('')
const pwError = ref('')

const changePassword = async () => {
  pwError.value = ''
  pwMessage.value = ''
  pwLoading.value = true
  try {
    await api('/me/change-password', {
      method: 'POST',
      body: { old_password: oldPassword.value, new_password: newPassword.value },
    })
    pwMessage.value = t('common.save')
    oldPassword.value = ''
    newPassword.value = ''
  } catch (e: any) {
    pwError.value = e?.data?.error?.message || 'Change failed'
  } finally {
    pwLoading.value = false
  }
}

useHead({ title: () => `${t('nav.profile')} — AliStroy` })
</script>

<template>
  <div class="container-page py-8 max-w-2xl">
    <h1 class="text-3xl font-bold text-ink-900 mb-6">{{ t('nav.profile') }}</h1>

    <section class="card p-6">
      <h2 class="font-semibold mb-4">{{ t('auth.name') }} / {{ t('auth.phone') }}</h2>
      <form class="space-y-3" @submit.prevent="save">
        <div>
          <label class="label">{{ t('auth.email') }}</label>
          <input :value="auth.user?.email" disabled class="input bg-gray-50" />
        </div>
        <div>
          <label class="label">{{ t('auth.name') }}</label>
          <input v-model="name" class="input" />
        </div>
        <div>
          <label class="label">{{ t('auth.phone') }}</label>
          <input v-model="phone" class="input" />
        </div>
        <div v-if="message" class="text-sm text-emerald-600">{{ message }}</div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
        <button type="submit" class="btn-primary" :disabled="loading">{{ t('common.save') }}</button>
      </form>
    </section>

    <section class="card p-6 mt-6">
      <h2 class="font-semibold mb-4">{{ t('auth.change_password') }}</h2>
      <form class="space-y-3" @submit.prevent="changePassword">
        <div>
          <label class="label">{{ t('auth.old_password') }}</label>
          <input v-model="oldPassword" type="password" class="input" required />
        </div>
        <div>
          <label class="label">{{ t('auth.new_password') }}</label>
          <input v-model="newPassword" type="password" class="input" required minlength="8" />
        </div>
        <div v-if="pwMessage" class="text-sm text-emerald-600">{{ pwMessage }}</div>
        <div v-if="pwError" class="text-sm text-red-600">{{ pwError }}</div>
        <button type="submit" class="btn-primary" :disabled="pwLoading">{{ t('auth.change_password') }}</button>
      </form>
    </section>
  </div>
</template>
