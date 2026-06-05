<script setup lang="ts">
const props = defineProps<{
  modelValue: { url: string; alt?: string }[]
  endpoint?: string // e.g. '/seller/upload' or '/admin/upload'
  subdir?: string
  multiple?: boolean
}>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: { url: string; alt?: string }[]): void
}>()

const api = useApi()
const { resolve } = useImageUrl()
const { t } = useI18n()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
const error = ref('')

const onPick = () => fileInput.value?.click()

const onFile = async (e: Event) => {
  const input = e.target as HTMLInputElement
  if (!input.files || !input.files.length) return
  uploading.value = true
  error.value = ''
  const next = [...props.modelValue]
  try {
    for (const file of Array.from(input.files)) {
      const fd = new FormData()
      fd.append('file', file)
      fd.append('subdir', props.subdir || 'products')
      const res = await api<{ data: { url: string } }>(props.endpoint || '/seller/upload', {
        method: 'POST',
        body: fd,
      })
      next.push({ url: res.data.url, alt: file.name })
    }
    emit('update:modelValue', next)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Upload failed'
  } finally {
    uploading.value = false
    input.value = ''
  }
}

const remove = (i: number) => {
  const next = [...props.modelValue]
  next.splice(i, 1)
  emit('update:modelValue', next)
}

const move = (i: number, dir: -1 | 1) => {
  const j = i + dir
  if (j < 0 || j >= props.modelValue.length) return
  const next = [...props.modelValue]
  ;[next[i], next[j]] = [next[j], next[i]]
  emit('update:modelValue', next)
}
</script>

<template>
  <div>
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div v-for="(img, i) in modelValue" :key="i" class="relative aspect-square rounded-lg overflow-hidden border border-gray-200 bg-gray-50">
        <img :src="resolve(img.url)" alt="" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-black/0 hover:bg-black/30 transition-colors flex items-end justify-between gap-1 p-1">
          <div class="flex gap-1">
            <button type="button" class="bg-white/90 rounded px-1.5 text-xs" @click="move(i, -1)">←</button>
            <button type="button" class="bg-white/90 rounded px-1.5 text-xs" @click="move(i, 1)">→</button>
          </div>
          <button type="button" class="bg-red-600 text-white rounded px-1.5 text-xs" @click="remove(i)">×</button>
        </div>
        <span v-if="i === 0" class="absolute top-1 left-1 badge bg-brand-500 text-white text-[10px]">★</span>
      </div>

      <button type="button" class="aspect-square rounded-lg border-2 border-dashed border-gray-300 hover:border-brand-500 text-gray-400 hover:text-brand-600 flex flex-col items-center justify-center text-sm" @click="onPick" :disabled="uploading">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 4v16m8-8H4"/></svg>
        <span class="mt-1">{{ t('seller_panel.upload_image') }}</span>
      </button>
    </div>
    <input ref="fileInput" type="file" :multiple="multiple !== false" accept="image/*" hidden @change="onFile" />
    <div v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</div>
    <div v-if="uploading" class="mt-2 text-sm text-gray-500">{{ t('common.loading') }}</div>
  </div>
</template>
