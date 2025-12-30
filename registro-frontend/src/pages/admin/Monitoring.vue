<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">System Monitoring</div>
    
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-3">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">System Health</div>
            <div class="text-h4">{{ store.health?.status || 'Unknown' }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Active Users</div>
            <div class="text-h4">{{ store.activeUsers }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card>
      <q-card-section>
        <div class="text-h6">API Traffic (24h)</div>
      </q-card-section>
      <q-card-section style="height: 300px">
        <MetricsChart :data="store.metrics" v-if="store.metrics.length" />
        <div v-else class="text-center q-pa-md">Loading metrics...</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { useMonitoring } from 'src/composables/useMonitoring';
import MetricsChart from 'src/components/Admin/MetricsChart.vue';

const { store } = useMonitoring();
</script>
