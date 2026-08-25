<template>
  <q-card class="bg-blue-1 text-center q-py-md">
    <div class="text-h6">{{ t('roleDashboards.attendanceAverage') || 'Percentuale Presenze' }}</div>
    <q-circular-progress
      show-value
      :value="stats?.percentage || 0"
      size="80px"
      :thickness="0.2"
      color="primary"
      track-color="blue-2"
      class="q-ma-md"
    >
      {{ stats?.percentage || 0 }}%
    </q-circular-progress>
    <div class="row justify-center q-gutter-md text-caption">
        <div><q-badge color="green" /> {{ stats?.present || 0 }} {{ t('attendance.present') || 'Presente' }}</div>
        <div><q-badge color="red" /> {{ stats?.absent || 0 }} {{ t('attendance.absent') || 'Assente' }}</div>
        <div><q-badge color="orange" /> {{ stats?.late || 0 }} {{ t('attendance.late') || 'In Ritardo' }}</div>
    </div>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
  stats: {
    type: Object,
    default: () => ({ percentage: 100, present: 0, absent: 0, late: 0 })
  }
});
</script>
