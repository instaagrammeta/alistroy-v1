<script setup lang="ts">
const props = defineProps<{ error: { statusCode?: number; message?: string } }>()
const goHome = () => clearError({ redirect: '/' })
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-6 bg-gray-50">
    <div class="text-center max-w-lg">
      <div class="text-7xl font-extrabold text-brand-500">{{ props.error?.statusCode || 500 }}</div>
      <h1 class="mt-4 text-2xl font-semibold text-ink-900">{{ props.error?.message || 'Error' }}</h1>
      <button class="btn-primary mt-6" @click="goHome">← AliStroy</button>
    </div>
  </div>
</template>
