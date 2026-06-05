<script setup lang="ts">
import type { ChatMessage, PaginatedResponse } from '~/types/api'

definePageMeta({ middleware: 'auth' })
const { t } = useI18n()
const { formatDateTime } = useFormatters()
const { resolve } = useImageUrl()
const auth = useAuthStore()

const messages = ref<ChatMessage[]>([])
const body = ref('')
const attachments = ref<{ url: string; mime_type: string; size_bytes: number }[]>([])
const sending = ref(false)

const load = async () => {
  const res = await useApi()<PaginatedResponse<ChatMessage>>('/customer/chat/messages', { query: { page: 1, page_size: 100 } })
  messages.value = (res.data || []).slice().reverse()
  await nextTick()
  scrollBottom()
}

const send = async () => {
  if (!body.value.trim() && !attachments.value.length) return
  sending.value = true
  try {
    await useApi()('/customer/chat/messages', { method: 'POST', body: { body: body.value, attachments: attachments.value } })
    body.value = ''; attachments.value = []
    await load()
  } finally {
    sending.value = false
  }
}

const container = ref<HTMLElement | null>(null)
const scrollBottom = () => {
  if (container.value) container.value.scrollTop = container.value.scrollHeight
}

let socket: WebSocket | null = null
onMounted(async () => {
  await load()
  // realtime updates
  try {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    const base = useRuntimeConfig().public.apiBase as string
    socket = new WebSocket(`${proto}://${location.host}${base}/customer/chat/socket?token=${auth.accessToken}`)
    socket.onmessage = () => load()
  } catch {}
})
onBeforeUnmount(() => socket?.close())

const onUploaded = (v: { url: string; alt?: string }[]) => {
  attachments.value = v.map((x) => ({ url: x.url, mime_type: '', size_bytes: 0 }))
}

useHead({ title: () => `${t('admin.chat')} — AliStroy` })
</script>

<template>
  <div class="container-page py-6 max-w-3xl">
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.chat') }}</h1>
    <div class="card flex flex-col h-[70vh]">
      <div ref="container" class="flex-1 overflow-auto p-4 space-y-3">
        <div v-for="m in messages" :key="m.id" class="flex" :class="m.sender_role === 'customer' ? 'justify-end' : 'justify-start'">
          <div class="max-w-[75%] rounded-2xl px-3 py-2 text-sm" :class="m.sender_role === 'customer' ? 'bg-brand-500 text-white' : 'bg-gray-100 text-ink-900'">
            <p v-if="m.body" class="whitespace-pre-line">{{ m.body }}</p>
            <div v-for="a in m.attachments" :key="a.id" class="mt-1">
              <img v-if="a.mime_type.startsWith('image') || /\.(jpg|jpeg|png|gif|webp)$/i.test(a.url)" :src="resolve(a.url)" class="rounded-lg max-h-48" />
              <video v-else-if="a.mime_type.startsWith('video') || /\.(mp4|webm|mov)$/i.test(a.url)" :src="resolve(a.url)" controls class="rounded-lg max-h-48"></video>
              <a v-else :href="resolve(a.url)" target="_blank" class="underline">file</a>
            </div>
            <div class="text-[10px] opacity-70 mt-1">{{ formatDateTime(m.created_at) }}</div>
          </div>
        </div>
        <EmptyState v-if="!messages.length" />
      </div>
      <div class="border-t border-gray-100 p-3">
        <ImageUploader :model-value="attachments.map(a => ({ url: a.url }))" endpoint="/customer/upload" subdir="chat" @update:model-value="onUploaded" />
        <div class="flex gap-2 mt-2">
          <input v-model="body" class="input flex-1" :placeholder="t('product.comment')" @keyup.enter="send" />
          <button class="btn-primary" :disabled="sending" @click="send">→</button>
        </div>
      </div>
    </div>
  </div>
</template>
