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
/* Accessibility DSA Font */
body.dsa-font-active,
body.dsa-font-active * {
  font-family: 'Atkinson Hyperlegible', 'Trebuchet MS', sans-serif !important;
  letter-spacing: 0.06em !important;
  word-spacing: 0.12em !important;
  line-height: 1.65 !important;
}

/* High Contrast Mode */
body.high-contrast-active {
  filter: contrast(130%);
}
body.high-contrast-active .q-card,
body.high-contrast-active .q-btn {
  border: 2px solid #1e293b !important;
}
</style>
