<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">{{ $t('nav.dashboard') }}</div>
    
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-3">
        <q-btn class="full-width" color="primary" icon="school" :label="$t('nav.schools')" to="/admin/schools" size="lg" />
      </div>
      <div class="col-12 col-md-3">
        <q-btn class="full-width" color="secondary" icon="analytics" :label="$t('nav.monitoring')" to="/admin/monitoring" size="lg" />
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <q-card>
          <q-card-section>
            <div class="text-h6">{{ $t('common.stats') || 'Quick Stats' }}</div>
          </q-card-section>
          <q-card-section>
             <!-- Mini stats -->
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
