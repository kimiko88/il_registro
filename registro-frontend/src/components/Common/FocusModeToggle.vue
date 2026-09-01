<template>
  <q-btn
    flat
    round
    dense
    :color="themeStore.isFocusMode ? 'amber-8' : 'primary'"
    :icon="themeStore.isFocusMode ? 'visibility_off' : 'center_focus_strong'"
    :aria-label="themeStore.isFocusMode ? 'Disattiva Modalità Lettura Pulita / Focus ADHD' : 'Attiva Modalità Lettura Pulita / Focus ADHD'"
    :aria-pressed="themeStore.isFocusMode"
    @click="themeStore.toggleFocusMode()"
  >
    <q-tooltip anchor="top middle" self="bottom middle">
      {{ themeStore.isFocusMode ? 'Disattiva Modalità Lettura Pulita (ADHD/DSA)' : 'Modalità Lettura Pulita / Focus ADHD' }}
    </q-tooltip>
  </q-btn>
</template>

<script setup>
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
</script>
