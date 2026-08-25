<template>
  <div class="q-pa-md">
    <div class="text-h6">{{ $t('documentsPage.title') || 'I Miei Documenti' }}</div>
    <p>{{ $t('documentsPage.subtitle') || 'Documenti scaricabili' }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
