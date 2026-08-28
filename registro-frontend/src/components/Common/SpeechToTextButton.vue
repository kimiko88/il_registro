<template>
  <q-btn
    flat
    round
    dense
    :color="isListening ? 'negative' : color"
    :icon="isListening ? 'mic' : 'mic_none'"
    :class="{ 'stt-pulse': isListening }"
    :aria-label="isListening ? 'Interrompi dettatura vocale' : 'Avvia dettatura vocale (Alt+D)'"
    :aria-pressed="isListening"
    :disabled="!isSupported"
    @click="handleToggle"
  >
    <q-tooltip anchor="top middle" self="bottom middle">
      {{ isSupported ? (isListening ? 'Interrompi dettatura' : 'Dettatura Vocale (Alt+D)') : 'Dettatura vocale non supportata' }}
    </q-tooltip>

    <q-badge v-if="isListening" floating color="red" rounded class="pulse-badge" />
  </q-btn>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useSpeechToText } from '@/composables/useSpeechToText'

const props = defineProps({
  color: {
    type: String,
    default: 'primary'
  },
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'result'])

const { isListening, isSupported, start, stop, toggle } = useSpeechToText({
  onResult: (text) => {
    let newText = text
    if (props.modelValue && props.modelValue.trim()) {
      newText = `${props.modelValue.trim()} ${text}`
    }
    emit('update:modelValue', newText)
    emit('result', text)
  }
})

function handleToggle() {
  toggle()
}

function handleKeyDown(e) {
  if (e.altKey && (e.key === 'd' || e.key === 'D')) {
    e.preventDefault()
    handleToggle()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.stt-pulse {
  animation: pulse-border 1.5s infinite;
}

@keyframes pulse-border {
  0% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(239, 68, 68, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0);
  }
}

.pulse-badge {
  animation: pulse-opacity 1s infinite alternate;
}

@keyframes pulse-opacity {
  from { opacity: 0.4; }
  to { opacity: 1; }
}
</style>
