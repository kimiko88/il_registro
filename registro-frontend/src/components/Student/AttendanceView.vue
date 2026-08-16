<template>
  <div class="q-pa-md">
    <div class="text-h6">{{ $t('studentAttendance.title') || 'Le Mie Assenze' }}</div>
    <p>{{ $t('studentAttendance.subtitle') || 'Storico presenze e assenze' }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
