<script setup lang="ts">
import type { ChatMessage, ChatRoom, PaginatedResponse } from '~/types/api'

definePageMeta({ layout: 'admin', middleware: 'admin' })
const { t } = useI18n()
const { formatDateTime } = useFormatters()
const { resolve } = useImageUrl()

const { data: roomsRes, refresh: refreshRooms } = await useAsyncData('admin-chat-rooms', () => useApi()<PaginatedResponse<ChatRoom>>('/admin/chat/rooms', { query: { page_size: 100 } }))
const rooms = computed<ChatRoom[]>(() => (roomsRes.value as any)?.data || [])

const activeRoom = ref<ChatRoom | null>(null)
const messages = ref<ChatMessage[]>([])
const body = ref('')
const container = ref<HTMLElement | null>(null)
let socket: WebSocket | null = null

const openRoom = async (room: ChatRoom) => {
  activeRoom.value = room
  await loadMessages()
  connectSocket(room.id)
}
const loadMessages = async () => {
  if (!activeRoom.value) return
  const res = await useApi()<PaginatedResponse<ChatMessage>>(`/admin/chat/rooms/${activeRoom.value.id}/messages`, { query: { page_size: 100 } })
  messages.value = (res.data || []).slice().reverse()
  await nextTick(); if (container.value) container.value.scrollTop = container.value.scrollHeight
}
const send = async () => {
  if (!activeRoom.value || !body.value.trim()) return
  await useApi()(`/admin/chat/rooms/${activeRoom.value.id}/messages`, { method: 'POST', body: { body: body.value } })
  body.value = ''
  await loadMessages()
  await refreshRooms()
}
const connectSocket = (roomId: string) => {
  socket?.close()
  try {
    const auth = useAuthStore()
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    const base = useRuntimeConfig().public.apiBase as string
    socket = new WebSocket(`${proto}://${location.host}${base}/admin/chat/rooms/${roomId}/socket?token=${auth.accessToken}`)
    socket.onmessage = () => loadMessages()
  } catch {}
}
onBeforeUnmount(() => socket?.close())

useHead({ title: () => `${t('admin.chat')} — Admin` })
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-ink-900 mb-4">{{ t('admin.chat') }}</h1>
    <div class="grid lg:grid-cols-[300px_1fr] gap-4 h-[70vh]">
      <div class="card overflow-auto">
        <button v-for="room in rooms" :key="room.id" class="w-full text-left px-4 py-3 border-b border-gray-50 hover:bg-gray-50 flex items-center justify-between" :class="{ 'bg-brand-50': activeRoom?.id === room.id }" @click="openRoom(room)">
          <div>
            <div class="font-medium text-ink-900">{{ room.customer?.company || room.customer_id?.slice?.(0, 8) || 'Customer' }}</div>
            <div class="text-xs text-gray-400">{{ formatDateTime(room.last_message_at) }}</div>
          </div>
          <span v-if="room.unread_admin > 0" class="badge bg-brand-500 text-white">{{ room.unread_admin }}</span>
        </button>
        <EmptyState v-if="!rooms.length" />
      </div>

      <div class="card flex flex-col">
        <template v-if="activeRoom">
          <div ref="container" class="flex-1 overflow-auto p-4 space-y-3">
            <div v-for="m in messages" :key="m.id" class="flex" :class="m.sender_role === 'admin' ? 'justify-end' : 'justify-start'">
              <div class="max-w-[75%] rounded-2xl px-3 py-2 text-sm" :class="m.sender_role === 'admin' ? 'bg-brand-500 text-white' : 'bg-gray-100 text-ink-900'">
                <p v-if="m.body" class="whitespace-pre-line">{{ m.body }}</p>
                <div v-for="a in m.attachments" :key="a.id" class="mt-1">
                  <img v-if="/\.(jpg|jpeg|png|gif|webp)$/i.test(a.url)" :src="resolve(a.url)" class="rounded-lg max-h-48" />
                  <video v-else-if="/\.(mp4|webm|mov)$/i.test(a.url)" :src="resolve(a.url)" controls class="rounded-lg max-h-48"></video>
                </div>
                <div class="text-[10px] opacity-70 mt-1">{{ formatDateTime(m.created_at) }}</div>
              </div>
            </div>
          </div>
          <div class="border-t border-gray-100 p-3 flex gap-2">
            <input v-model="body" class="input flex-1" @keyup.enter="send" />
            <button class="btn-primary" @click="send">→</button>
          </div>
        </template>
        <div v-else class="flex-1 flex items-center justify-center text-gray-400">{{ t('common.select') }}</div>
      </div>
    </div>
  </div>
</template>
