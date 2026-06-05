<script setup lang="ts">
const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()

const onLogout = async () => {
  await auth.logout()
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <AppHeader />
    <div class="flex-1">
      <div class="container-page py-6">
        <div class="flex items-center justify-between mb-4">
          <h1 class="text-2xl font-bold text-ink-900">{{ t('nav.driver_panel') }}</h1>
          <button class="btn-ghost text-red-600" @click="onLogout">{{ t('nav.logout') }}</button>
        </div>
        <slot />
      </div>
    </div>
  </div>
</template>
