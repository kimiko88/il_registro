<template>
  <div class="skeleton-table" :style="{ '--rows': rows, '--cols': cols }">
    <!-- Header row -->
    <div class="skeleton-header row q-gutter-xs q-mb-sm">
      <div
        v-for="c in cols"
        :key="`h-${c}`"
        class="skeleton-cell skeleton-header-cell"
        :style="{ flex: c === 1 ? '2' : '1' }"
      />
    </div>
    <!-- Data rows -->
    <div
      v-for="r in rows"
      :key="`r-${r}`"
      class="skeleton-row row q-gutter-xs q-mb-xs"
      :style="{ opacity: 1 - (r / (rows + 2)) * 0.4 }"
    >
      <div
        v-for="c in cols"
        :key="`c-${c}`"
        class="skeleton-cell"
        :style="{ flex: c === 1 ? '2' : '1' }"
      />
    </div>
  </div>
</template>

<script setup>
defineProps({
  rows: { type: Number, default: 8 },
  cols: { type: Number, default: 5 }
})
</script>

<style scoped>
.skeleton-table {
  padding: 0.5rem;
}

.skeleton-row,
.skeleton-header {
  display: flex;
  gap: 8px;
  align-items: center;
}

.skeleton-cell {
  border-radius: 6px;
  height: 18px;
  background: linear-gradient(90deg, #e8eaed 25%, #f8f9fa 50%, #e8eaed 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  min-width: 40px;
}

.skeleton-header-cell {
  height: 14px;
  background: linear-gradient(90deg, #d0d4da 25%, #e8eaed 50%, #d0d4da 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  border-radius: 4px;
}

@keyframes skeleton-shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.body--dark .skeleton-cell {
  background: linear-gradient(90deg, #2d3748 25%, #374151 50%, #2d3748 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
}
.body--dark .skeleton-header-cell {
  background: linear-gradient(90deg, #374151 25%, #4a5568 50%, #374151 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
}
</style>
