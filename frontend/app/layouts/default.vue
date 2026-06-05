<script setup lang="ts">
const settings = useSettingsStore()
onMounted(async () => {
  if (!settings.loaded) await settings.load()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-50">
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>
