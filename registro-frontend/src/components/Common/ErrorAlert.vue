<template>
  <q-banner v-if="error" dense rounded class="bg-negative text-white q-mb-md">
    {{ error }}
    <template v-slot:action>
      <q-btn flat color="white" :label="t('common.close') || 'Chiudi'" @click="$emit('dismiss')" />
    </template>
  </q-banner>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
  error: {
    type: String,
    default: null
  }
})

defineEmits(['dismiss'])
</script>
