<template>
  <div class="q-pa-md">
    <div class="text-h6">{{ $t('studentAttendance.title') || 'Presenze Studente' }}</div>
    <p>{{ $t('studentAttendance.subtitle') || 'Dettaglio presenze e assenze' }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
