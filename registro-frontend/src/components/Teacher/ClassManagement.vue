<template>
  <div class="q-pa-md">
    <div class="text-h6">{{ $t('udaPage.classLabel') || 'Gestione Classe' }}</div>
    <p>{{ $t('common.info') || 'Gestione della classe' }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
