<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">Global Analytics</div>
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-subtitle2">Total Schools</div>
            <div class="text-h3">{{ stats.schools }}</div>
          </q-card-section>
        </q-card>
      </div>
      <!-- More widgets -->
    </div>
  </q-page>
</template>

<script setup>
import { useGlobalAnalytics } from 'src/composables/useGlobalAnalytics';

const { stats } = useGlobalAnalytics();
</script>
