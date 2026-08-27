<template>
  <q-btn
    v-if="isSupported"
    :flat="flat"
    :round="round"
    :dense="dense"
    :size="size"
    :color="isPlaying ? 'positive' : color"
    :icon="isPlaying ? (isPaused ? 'play_arrow' : 'pause') : icon"
    :class="['tts-button', { 'tts-active': isPlaying }]"
    @click.stop="handleToggle"
    no-caps
    :aria-label="ariaLabelText"
  >
    <q-tooltip>{{ tooltipText }}</q-tooltip>
  </q-btn>
</template>

<script setup>
import { computed } from 'vue'
import { useSpeechSynthesis } from 'src/composables/useSpeechSynthesis'

const props = defineProps({
  text: {
    type: String,
    required: true
  },
  lang: {
    type: String,
    default: 'it-IT'
  },
  rate: {
    type: Number,
    default: undefined
  },
  color: {
    type: String,
    default: 'primary'
  },
  size: {
    type: String,
    default: 'sm'
  },
  flat: {
    type: Boolean,
    default: true
  },
  round: {
    type: Boolean,
    default: true
  },
  dense: {
    type: Boolean,
    default: true
  },
  icon: {
    type: String,
    default: 'volume_up'
  }
})

const { isSupported, isPlaying, isPaused, toggle, stop } = useSpeechSynthesis()

const tooltipText = computed(() => {
  if (isPlaying.value) {
    return isPaused.value ? 'Riprendi lettura vocale' : 'Metti in pausa lettura'
  }
  return 'Ascolta sintesi vocale (TTS)'
})

const ariaLabelText = computed(() => {
  return isPlaying.value ? 'Interrompi lettura vocale' : 'Ascolta testo con sintesi vocale'
})

function handleToggle() {
  toggle(props.text, { lang: props.lang, rate: props.rate })
}

defineExpose({
  stop
})
</script>

<style scoped>
.tts-button {
  transition: transform 0.2s ease, color 0.2s ease;
}
.tts-active {
  animation: pulse-tts 1.5s infinite;
}
@keyframes pulse-tts {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}
</style>
