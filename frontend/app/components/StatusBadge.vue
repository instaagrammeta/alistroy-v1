<script setup lang="ts">
const props = defineProps<{ status: string }>()
const { t } = useI18n()
const cls = computed(() => {
  switch (props.status) {
    case 'approved':
      return 'badge-approved'
    case 'pending':
      return 'badge-pending'
    case 'rejected':
      return 'badge-rejected'
    default:
      return 'badge-draft'
  }
})
const label = computed(() => {
  switch (props.status) {
    case 'approved':
      return t('seller_panel.status_approved')
    case 'pending':
      return t('seller_panel.status_pending')
    case 'rejected':
      return t('seller_panel.status_rejected')
    default:
      return t('seller_panel.status_draft')
  }
})
</script>

<template>
  <span :class="cls">{{ label }}</span>
</template>
