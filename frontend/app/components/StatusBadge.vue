<script setup lang="ts">
const props = defineProps<{ status: string }>()
const { t } = useI18n()
const cls = computed(() => {
  switch (props.status) {
    case 'approved': case 'completed': return 'badge-approved'
    case 'pending': case 'processing': case 'assigned': case 'on_delivery': case 'new': return 'badge-pending'
    case 'rejected': case 'cancelled': return 'badge-rejected'
    default: return 'badge-draft'
  }
})
const label = computed(() => {
  const orderStatuses = ['new', 'processing', 'assigned', 'on_delivery', 'completed', 'cancelled']
  if (orderStatuses.includes(props.status)) return t(`order_status.${props.status}`)
  return t(`seller_panel.status_${props.status}`)
})
</script>

<template>
  <span :class="cls">{{ label }}</span>
</template>
