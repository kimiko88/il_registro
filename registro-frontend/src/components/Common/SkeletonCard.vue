<template>
  <q-card :class="cardClass" :style="cardStyle">
    <q-card-section>
      <!-- Title skeleton -->
      <div class="skeleton-line" style="width: 40%; height: 20px; margin-bottom: 16px;" />
      <!-- Content lines -->
      <div v-for="i in lines" :key="i" class="skeleton-line q-mb-sm" :style="{ width: lineWidths[i % lineWidths.length], height: '14px' }" />
      <!-- Footer skeleton (optional) -->
      <div v-if="showFooter" class="row q-mt-md q-gutter-sm">
        <div class="skeleton-line" style="width: 80px; height: 32px; border-radius: 8px;" />
        <div class="skeleton-line" style="width: 60px; height: 32px; border-radius: 8px;" />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
  lines: { type: Number, default: 4 },
  showFooter: { type: Boolean, default: false },
  cardClass: { type: String, default: 'shadow-1' },
  cardStyle: { type: String, default: '' }
})

const lineWidths = ['100%', '85%', '90%', '70%', '95%', '60%']
</script>

<style scoped>
.skeleton-line {
  border-radius: 6px;
  display: block;
  background: linear-gradient(90deg, #e8eaed 25%, #f8f9fa 50%, #e8eaed 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
}

@keyframes skeleton-shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.body--dark .skeleton-line {
  background: linear-gradient(90deg, #2d3748 25%, #374151 50%, #2d3748 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
}
</style>
