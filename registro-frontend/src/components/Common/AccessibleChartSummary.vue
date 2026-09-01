<template>
  <div class="accessible-chart-summary q-my-md">
    <!-- Narrative Text Summary for Screen Readers & General Users -->
    <q-banner dense rounded class="bg-blue-1 text-blue-10 q-mb-sm border-blue" role="region" aria-label="Sintesi accessibile dei dati">
      <template #avatar>
        <q-icon name="analytics" color="primary" size="24px" />
      </template>
      <div class="text-subtitle2 text-bold">Sintesi Dati Accessibile (WCAG 1.1.1)</div>
      <div class="text-body2 text-weight-medium q-mt-xs">
        {{ narrativeSummary }}
      </div>

      <template #action>
        <q-btn
          flat
          dense
          color="primary"
          :icon="showTable ? 'expand_less' : 'table_view'"
          :label="showTable ? 'Nascondi Tabella Dati' : 'Mostra Tabella Accessibile'"
          :aria-expanded="showTable"
          @click="showTable = !showTable"
        />
      </template>
    </q-banner>

    <!-- Accessible High-Contrast Data Table -->
    <q-slide-transition>
      <div v-show="showTable">
        <q-table
          flat
          bordered
          dense
          :rows="rows"
          :columns="columns"
          row-key="id"
          :pagination="{ rowsPerPage: 10 }"
          class="accessible-data-table"
          aria-label="Tabella dati grafici accessibile"
        >
          <template #body-cell="props">
            <q-td :props="props" class="text-weight-medium">
              {{ props.value }}
            </q-td>
          </template>
        </q-table>
      </div>
    </q-slide-transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: 'Andamento'
  },
  items: {
    type: Array,
    default: () => []
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  valueKey: {
    type: String,
    default: 'value'
  }
})

const showTable = ref(false)

const columns = computed(() => [
  { name: 'label', align: 'left', label: 'Voce / Periodo', field: props.labelKey, sortable: true },
  { name: 'value', align: 'right', label: 'Valore / Esito', field: props.valueKey, sortable: true }
])

const rows = computed(() => {
  return props.items.map((item, idx) => ({
    id: idx,
    ...item
  }))
})

const narrativeSummary = computed(() => {
  if (!props.items || props.items.length === 0) {
    return `${props.title}: Nessun dato registrato.`
  }

  const numericValues = props.items
    .map(i => parseFloat(i[props.valueKey]))
    .filter(v => !isNaN(v))

  if (numericValues.length === 0) {
    return `${props.title}: ${props.items.length} elementi registrati.`
  }

  const sum = numericValues.reduce((a, b) => a + b, 0)
  const avg = (sum / numericValues.length).toFixed(1)
  const max = Math.max(...numericValues)
  const min = Math.min(...numericValues)

  return `${props.title}: ${numericValues.length} elementi analizzati. Media complessiva: ${avg}. Valore massimo registrato: ${max}. Valore minimo registrato: ${min}.`
})
</script>

<style scoped>
.border-blue {
  border: 1px solid #90caf9;
}

.accessible-data-table :deep(th) {
  font-weight: 700;
  background-color: #f5f5f5;
  font-size: 0.95rem;
}
</style>
