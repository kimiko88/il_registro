<template>
  <router-view />
</template>

<script setup>
import { onMounted, watch } from 'vue'
import { useAuthStore } from 'src/stores/auth'
import { useWebSocketStore } from 'src/stores/websocket'

const authStore = useAuthStore()
const wsStore = useWebSocketStore()

onMounted(() => {
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
