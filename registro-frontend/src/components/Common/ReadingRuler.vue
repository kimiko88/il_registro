<template>
  <div v-if="themeStore.readingRuler" class="reading-ruler-container" aria-hidden="true">
    <!-- Top Mask -->
    <div
      class="ruler-mask ruler-mask-top"
      :style="{
        height: `${maskTopHeight}px`,
        backgroundColor: `rgba(0, 0, 0, ${themeStore.readingRulerOpacity})`
      }"
    />

    <!-- Active Line Highlight -->
    <div
      class="ruler-active-line"
      :style="{
        top: `${maskTopHeight}px`,
        height: `${themeStore.readingRulerHeight}px`
      }"
    />

    <!-- Bottom Mask -->
    <div
      class="ruler-mask ruler-mask-bottom"
      :style="{
        top: `${maskTopHeight + themeStore.readingRulerHeight}px`,
        backgroundColor: `rgba(0, 0, 0, ${themeStore.readingRulerOpacity})`
      }"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useThemeStore } from 'src/stores/theme'

const themeStore = useThemeStore()
const currentY = ref(200)
const maskTopHeight = ref(180)

function onMouseMove(e) {
  currentY.value = e.clientY
  updateRulerPosition()
}

function onKeyDown(e) {
  if (!themeStore.readingRuler) return
  if (e.key === 'ArrowDown' && e.altKey) {
    currentY.value += 30
    updateRulerPosition()
  } else if (e.key === 'ArrowUp' && e.altKey) {
    currentY.value = Math.max(50, currentY.value - 30)
    updateRulerPosition()
  }
}

function updateRulerPosition() {
  const halfHeight = (themeStore.readingRulerHeight || 40) / 2
  maskTopHeight.value = Math.max(0, currentY.value - halfHeight)
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('mousemove', onMouseMove, { passive: true })
    window.addEventListener('keydown', onKeyDown, { passive: true })
    updateRulerPosition()
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('keydown', onKeyDown)
  }
})
</script>

<style scoped>
.reading-ruler-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 9999;
  overflow: hidden;
}

.ruler-mask {
  position: absolute;
  left: 0;
  width: 100%;
  pointer-events: none;
  transition: background-color 0.2s ease;
}

.ruler-mask-top {
  top: 0;
}

.ruler-mask-bottom {
  bottom: 0;
}

.ruler-active-line {
  position: absolute;
  left: 0;
  width: 100%;
  border-top: 2px solid rgba(245, 158, 11, 0.8);
  border-bottom: 2px solid rgba(245, 158, 11, 0.8);
  background-color: rgba(254, 240, 138, 0.15);
  box-shadow: 0 0 15px rgba(245, 158, 11, 0.2);
  pointer-events: none;
}
</style>
