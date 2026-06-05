<template>
  <div class="min-h-screen bg-gradient-to-br from-brand-50 via-white to-white flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-md">
      <div class="text-center mb-6">
        <NuxtLink to="/" class="inline-flex items-center gap-2 text-2xl font-extrabold text-ink-900">
          <span class="flex items-center justify-center rounded-lg bg-brand-500 text-white font-bold w-9 h-9">A</span>
          AliStroy
        </NuxtLink>
      </div>
      <div class="card p-6">
        <slot />
      </div>
    </div>
  </div>
</template>
