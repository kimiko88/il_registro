<template>
  <q-card class="bg-blue-1 text-center q-py-md">
    <div class="text-h6">Attendance Rate</div>
    <q-circular-progress
      show-value
      :value="stats.percentage"
      size="80px"
      :thickness="0.2"
      color="primary"
      track-color="blue-2"
      class="q-ma-md"
    >
      {{ stats.percentage }}%
    </q-circular-progress>
    <div class="row justify-center q-gutter-md text-caption">
        <div><q-badge color="green" /> {{ stats.present }} Present</div>
        <div><q-badge color="red" /> {{ stats.absent }} Absent</div>
        <div><q-badge color="orange" /> {{ stats.late }} Late</div>
    </div>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps(['stats']);
</script>
