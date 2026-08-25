<template>
  <q-card class="q-mb-md">
    <q-card-section class="text-subtitle1">{{ t('roleDashboards.generalTrend') || 'Andamento Generale' }}</q-card-section>
    <q-card-section>
       <!-- Simple Visualization -->
       <div v-for="(avg, subject) in averages" :key="subject" class="row items-center q-mb-sm">
         <div class="col-3 text-caption ellipsis">{{ subject }}</div>
         <div class="col-9">
             <q-linear-progress :value="getProgressValue(avg)" size="15px" :color="getProgressColor(avg)" rounded>
                 <div class="absolute-full flex flex-center">
                     <q-badge color="transparent" text-color="white" :label="formatAvgLabel(avg)" />
                 </div>
             </q-linear-progress>
         </div>
       </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
  averages: {
    type: Object,
    default: () => ({})
  }
});

function getProgressValue(avg) {
  if (typeof avg === 'number' && !isNaN(avg)) {
    return Math.max(0, Math.min(1, avg / 10));
  }
  const parsed = parseFloat(avg);
  if (!isNaN(parsed)) {
    return Math.max(0, Math.min(1, parsed / 10));
  }
  return 0;
}

function getProgressColor(avg) {
  const val = typeof avg === 'number' ? avg : parseFloat(avg);
  if (isNaN(val)) return 'grey';
  return val >= 6 ? 'green' : 'red';
}

function formatAvgLabel(avg) {
  if (avg === undefined || avg === null || avg === '') return '-';
  return String(avg);
}
</script>
