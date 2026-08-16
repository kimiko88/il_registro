<template>
  <div class="flex flex-center q-pa-md">
    <q-spinner
      color="primary"
      size="3em"
      :thickness="10"
    />
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>
