<template>
  <div class="sr-announcer-container" aria-hidden="false">
    <div
      class="sr-only"
      role="status"
      aria-live="polite"
      aria-atomic="true"
    >
      {{ politeMessage }}
    </div>
    <div
      class="sr-only"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      {{ assertiveMessage }}
    </div>
  </div>
</template>

<script setup>
import { useA11yAnnouncer } from '@/composables/useA11yAnnouncer'

const { politeMessage, assertiveMessage } = useA11yAnnouncer()
</script>

<style scoped>
.sr-only {
  position: absolute !important;
  width: 1px !important;
  height: 1px !important;
  padding: 0 !important;
  margin: -1px !important;
  overflow: hidden !important;
  clip: rect(0, 0, 0, 0) !important;
  white-space: nowrap !important;
  border: 0 !important;
}
</style>
