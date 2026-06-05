<script setup lang="ts">
const settings = useSettingsStore()
const { t, locale } = useI18n()
const year = new Date().getFullYear()
</script>

<template>
  <footer class="bg-ink-900 text-gray-300 mt-12">
    <div class="container-page py-12 grid grid-cols-1 md:grid-cols-4 gap-8">
      <div>
        <AppLogo inverted />
        <p class="mt-3 text-sm text-gray-400 leading-relaxed">{{ t('footer.about') }}</p>
      </div>

      <div>
        <h4 class="font-semibold text-white mb-3">{{ t('footer.links') }}</h4>
        <ul class="space-y-2 text-sm">
          <li><NuxtLink to="/products" class="hover:text-white">{{ t('nav.catalog') }}</NuxtLink></li>
          <li><NuxtLink to="/categories" class="hover:text-white">{{ t('nav.categories') }}</NuxtLink></li>
          <li><NuxtLink to="/sellers" class="hover:text-white">{{ t('nav.sellers') }}</NuxtLink></li>
          <li><NuxtLink to="/search" class="hover:text-white">{{ t('nav.search') }}</NuxtLink></li>
        </ul>
      </div>

      <div>
        <h4 class="font-semibold text-white mb-3">{{ t('footer.for_sellers') }}</h4>
        <ul class="space-y-2 text-sm">
          <li><NuxtLink to="/register?role=seller" class="hover:text-white">{{ t('nav.become_seller') }}</NuxtLink></li>
          <li><NuxtLink to="/seller" class="hover:text-white">{{ t('nav.seller_panel') }}</NuxtLink></li>
        </ul>
      </div>

      <div>
        <h4 class="font-semibold text-white mb-3">{{ t('footer.contacts') }}</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li v-if="settings.marketplacePhone">
            <a :href="`tel:${settings.marketplacePhone.replace(/[^0-9+]/g, '')}`" class="hover:text-white">
              {{ settings.marketplacePhone }}
            </a>
          </li>
          <li v-if="settings.marketplaceWhatsApp">
            <a :href="`https://wa.me/${settings.marketplaceWhatsApp.replace(/[^0-9]/g, '')}`" class="hover:text-white">
              WhatsApp: {{ settings.marketplaceWhatsApp }}
            </a>
          </li>
          <li v-if="settings.get('footer_email')">
            <a :href="`mailto:${settings.get('footer_email')}`" class="hover:text-white">{{ settings.get('footer_email') }}</a>
          </li>
          <li v-if="settings.get('footer_address')" class="text-gray-400">
            {{ settings.get('footer_address') }}
          </li>
        </ul>
      </div>
    </div>

    <div class="border-t border-white/10">
      <div class="container-page py-4 text-xs text-gray-500 text-center">
        {{ t('footer.rights', { year }) }}
        <span class="ml-2 text-gray-600">{{ locale }}</span>
      </div>
    </div>
  </footer>
</template>
