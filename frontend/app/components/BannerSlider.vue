<script setup lang="ts">
import type { Banner } from '~/types/api'

const props = defineProps<{ banners: Banner[]; interval?: number; rounded?: boolean; ratio?: string }>()
const { resolve } = useImageUrl()

const index = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const go = (i: number) => {
  if (!props.banners.length) return
  index.value = (i + props.banners.length) % props.banners.length
}
const next = () => go(index.value + 1)
const prev = () => go(index.value - 1)

const start = () => {
  stop()
  if (props.banners.length > 1) timer = setInterval(next, props.interval || 5000)
}
const stop = () => {
  if (timer) clearInterval(timer)
  timer = null
}

onMounted(start)
onBeforeUnmount(stop)
watch(() => props.banners.length, start)

const ratioClass = computed(() => props.ratio || 'aspect-[21/9]')
</script>

<template>
  <div v-if="banners.length" class="relative overflow-hidden group" :class="[rounded !== false ? 'rounded-2xl' : '', ratioClass]" @mouseenter="stop" @mouseleave="start">
    <TransitionGroup name="fade">
      <template v-for="(b, i) in banners" :key="b.id">
        <component
          :is="b.link_url ? 'a' : 'div'"
          v-show="i === index"
          :href="b.link_url || undefined"
          class="absolute inset-0 block animate-zoom-in"
        >
          <picture>
            <source v-if="b.mobile_url" :srcset="resolve(b.mobile_url)" media="(max-width: 640px)" />
            <source v-if="b.tablet_url" :srcset="resolve(b.tablet_url)" media="(max-width: 1024px)" />
            <img :src="resolve(b.desktop_url || b.tablet_url || b.mobile_url)" :alt="b.title_ru || b.title_tj" class="w-full h-full object-cover" />
          </picture>
        </component>
      </template>
    </TransitionGroup>

    <button v-if="banners.length > 1" class="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/80 hover:bg-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition" @click="prev">‹</button>
    <button v-if="banners.length > 1" class="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/80 hover:bg-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition" @click="next">›</button>

    <div v-if="banners.length > 1" class="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-2">
      <button v-for="(b, i) in banners" :key="b.id" class="w-2.5 h-2.5 rounded-full transition" :class="i === index ? 'bg-white w-6' : 'bg-white/60'" @click="go(i)" />
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.6s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
