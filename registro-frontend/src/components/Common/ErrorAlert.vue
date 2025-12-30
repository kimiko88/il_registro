<template>
  <q-banner v-if="error" dense rounded class="bg-negative text-white q-mb-md">
    {{ error }}
    <template v-slot:action>
      <q-btn flat color="white" label="Dismiss" @click="$emit('dismiss')" />
    </template>
  </q-banner>
</template>

<script setup>
defineProps({
  error: {
    type: String,
    default: null
  }
})

defineEmits(['dismiss'])
</script>
