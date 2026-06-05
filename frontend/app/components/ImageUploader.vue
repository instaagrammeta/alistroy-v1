<script setup lang="ts">
const props = defineProps<{
  modelValue: { url: string; alt?: string }[]
  endpoint?: string
  subdir?: string
  multiple?: boolean
}>()
const emit = defineEmits<{ (e: 'update:modelValue', v: { url: string; alt?: string }[]): void }>()

const { resolve } = useImageUrl()
const { t } = useI18n()
const auth = useAuthStore()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
const progress = ref(0)
const error = ref('')

const pick = () => fileInput.value?.click()

// Use XHR for real upload progress percentage.
const uploadOne = (file: File): Promise<{ url: string }> =>
  new Promise((resolve2, reject) => {
    const xhr = new XMLHttpRequest()
    const fd = new FormData()
    fd.append('file', file)
    fd.append('subdir', props.subdir || 'misc')
    xhr.open('POST', `${useApiBase()}${props.endpoint || '/admin/upload'}`)
    if (auth.accessToken) xhr.setRequestHeader('Authorization', `Bearer ${auth.accessToken}`)
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) progress.value = Math.round((e.loaded / e.total) * 100)
    }
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const json = JSON.parse(xhr.responseText)
          resolve2({ url: json.data.url })
        } catch {
          reject(new Error('bad response'))
        }
      } else {
        let msg = 'upload failed'
        try { msg = JSON.parse(xhr.responseText).error.message } catch {}
        reject(new Error(msg))
      }
    }
    xhr.onerror = () => reject(new Error('network error'))
    xhr.send(fd)
  })

const onFile = async (e: Event) => {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  uploading.value = true
  error.value = ''
  const next = [...props.modelValue]
  try {
    for (const file of Array.from(input.files)) {
      progress.value = 0
      const res = await uploadOne(file)
      next.push({ url: res.url, alt: file.name })
    }
    emit('update:modelValue', next)
  } catch (err: any) {
    error.value = err?.message || 'upload failed'
  } finally {
    uploading.value = false
    progress.value = 0
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
        <div class="absolute inset-0 hover:bg-black/30 transition-colors flex items-end justify-between gap-1 p-1">
          <div class="flex gap-1">
            <button type="button" class="bg-white/90 rounded px-1.5 text-xs" @click="move(i, -1)">←</button>
            <button type="button" class="bg-white/90 rounded px-1.5 text-xs" @click="move(i, 1)">→</button>
          </div>
          <button type="button" class="bg-red-600 text-white rounded px-1.5 text-xs" @click="remove(i)">×</button>
        </div>
        <span v-if="i === 0" class="absolute top-1 left-1 badge bg-brand-500 text-white text-[10px]">★</span>
      </div>
      <button type="button" class="aspect-square rounded-lg border-2 border-dashed border-gray-300 hover:border-brand-500 text-gray-400 hover:text-brand-600 flex flex-col items-center justify-center text-sm" :disabled="uploading" @click="pick">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 4v16m8-8H4"/></svg>
        <span class="mt-1">{{ t('seller_panel.upload_image') }}</span>
      </button>
    </div>
    <input ref="fileInput" type="file" :multiple="multiple !== false" accept="image/*,video/*" hidden @change="onFile" />
    <div v-if="uploading" class="mt-2">
      <div class="h-2 rounded-full bg-gray-100 overflow-hidden">
        <div class="h-full bg-brand-500 transition-all" :style="{ width: progress + '%' }"></div>
      </div>
      <div class="text-xs text-gray-500 mt-1">{{ progress }}%</div>
    </div>
    <div v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</div>
  </div>
</template>
