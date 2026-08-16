<template>
  <q-card class="q-mb-md">
    <q-card-section class="text-subtitle1">{{ t('roleDashboards.generalTrend') || 'Andamento Generale' }}</q-card-section>
    <q-card-section>
       <!-- Simple Visualization -->
       <div v-for="(avg, subject) in averages" :key="subject" class="row items-center q-mb-sm">
         <div class="col-3 text-caption ellipsis">{{ subject }}</div>
         <div class="col-9">
             <q-linear-progress :value="avg / 10" size="15px" :color="avg >= 6 ? 'green' : 'red'" rounded>
                 <div class="absolute-full flex flex-center">
                     <q-badge color="transparent" text-color="white" :label="avg" />
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
defineProps(['averages']);
</script>
