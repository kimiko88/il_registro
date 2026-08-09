<template>
  <router-view />
</template>

<script setup>
import { onMounted, watch } from 'vue'
import { useAuthStore } from 'src/stores/auth'
import { useWebSocketStore } from 'src/stores/websocket'
import { useThemeStore } from 'src/stores/theme'

const authStore = useAuthStore()
const wsStore = useWebSocketStore()
const themeStore = useThemeStore()

onMounted(() => {
  themeStore.initTheme()
  if (authStore.isAuthenticated) {
    wsStore.connect()
  }
})

watch(() => authStore.isAuthenticated, (val) => {
  if (val) {
    wsStore.connect()
  } else {
    wsStore.disconnect()
  }
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Atkinson+Hyperlegible:ital,wght@0,400;0,700;1,400;1,700&display=swap');

@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/downloads/OpenDyslexic-Regular.otf') format('opentype');
  font-weight: normal;
  font-style: normal;
}
@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/downloads/OpenDyslexic-Bold.otf') format('opentype');
  font-weight: bold;
  font-style: normal;
}

/* Accessibility DSA Font */
body.dsa-font-active,
body.dsa-font-active * {
  font-family: 'OpenDyslexic', 'Atkinson Hyperlegible', 'Trebuchet MS', sans-serif !important;
  letter-spacing: 0.05em !important;
  word-spacing: 0.12em !important;
  line-height: 1.65 !important;
}

/* High Contrast Mode */
body.high-contrast-active {
  filter: contrast(135%);
}
body.high-contrast-active .q-card,
body.high-contrast-active .q-btn,
body.high-contrast-active .q-table {
  border: 2px solid #0f172a !important;
}
</style>
