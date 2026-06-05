<script setup lang="ts">
definePageMeta({ middleware: 'auth' })
const { t } = useI18n()
const auth = useAuthStore()

const name = ref(auth.user?.name || '')
const phone = ref(auth.user?.phone || '')
const address = ref(auth.user?.customer?.address || '')
const city = ref(auth.user?.customer?.city || '')
const msg = ref(''); const err = ref(''); const loading = ref(false)

const save = async () => {
  msg.value = ''; err.value = ''; loading.value = true
  try {
    const res = await useApi()<{ data: any }>('/me', { method: 'PATCH', body: { name: name.value, phone: phone.value, address: address.value, city: city.value } })
    auth.setUser(res.data)
    msg.value = t('common.saved')
  } catch (e: any) {
    err.value = e?.data?.error?.message || 'Error'
  } finally {
    loading.value = false
  }
}

const oldPw = ref(''); const newPw = ref(''); const pwMsg = ref(''); const pwErr = ref('')
const changePw = async () => {
  pwMsg.value = ''; pwErr.value = ''
  try {
    await useApi()('/me/change-password', { method: 'POST', body: { old_password: oldPw.value, new_password: newPw.value } })
    pwMsg.value = t('common.saved'); oldPw.value = ''; newPw.value = ''
  } catch (e: any) {
    pwErr.value = e?.data?.error?.message || 'Error'
  }
}

useHead({ title: () => `${t('nav.profile')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6 max-w-2xl">
    <div class="flex items-center gap-3 mb-6">
      <h1 class="text-2xl md:text-3xl font-bold text-ink-900">{{ t('nav.profile') }}</h1>
      <NuxtLink to="/account/orders" class="btn-outline btn-sm ml-auto">{{ t('nav.orders') }}</NuxtLink>
      <NuxtLink to="/account/chat" class="btn-outline btn-sm">{{ t('admin.chat') }}</NuxtLink>
    </div>

    <section class="card p-6">
      <form class="space-y-3" @submit.prevent="save">
        <div><label class="label">{{ t('auth.name') }}</label><input v-model="name" class="input" /></div>
        <div><label class="label">{{ t('auth.phone') }}</label><input v-model="phone" class="input" /></div>
        <div><label class="label">{{ t('auth.address') }}</label><input v-model="address" class="input" /></div>
        <div><label class="label">{{ t('auth.city') }}</label><input v-model="city" class="input" /></div>
        <div v-if="msg" class="text-sm text-emerald-600">{{ msg }}</div>
        <div v-if="err" class="text-sm text-red-600">{{ err }}</div>
        <button type="submit" class="btn-primary" :disabled="loading">{{ t('common.save') }}</button>
      </form>
    </section>

    <section class="card p-6 mt-6">
      <h2 class="font-semibold mb-4">{{ t('auth.change_password') }}</h2>
      <form class="space-y-3" @submit.prevent="changePw">
        <div><label class="label">{{ t('auth.old_password') }}</label><input v-model="oldPw" type="password" class="input" /></div>
        <div><label class="label">{{ t('auth.new_password') }}</label><input v-model="newPw" type="password" minlength="8" class="input" /></div>
        <div v-if="pwMsg" class="text-sm text-emerald-600">{{ pwMsg }}</div>
        <div v-if="pwErr" class="text-sm text-red-600">{{ pwErr }}</div>
        <button type="submit" class="btn-primary">{{ t('auth.change_password') }}</button>
      </form>
    </section>
  </div>
</template>
