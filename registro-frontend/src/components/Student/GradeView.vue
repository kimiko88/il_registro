<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">{{ t('roleDashboards.tabGrades') || 'I Miei Voti' }}</div>

    <q-table
      :rows="grades"
      :columns="columns"
      row-key="id"
      :loading="loading"
    >
      <template v-slot:body-cell-value="props">
        <q-td :props="props">
          <q-badge :color="getGradeColor(props.value)">
            {{ props.value }}
          </q-badge>
        </q-td>
      </template>
    </q-table>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '@/composables/useApi'

const { t } = useI18n()
const { data: grades, loading } = useApi('/grades')

const columns = computed(() => [
  { name: 'subject', label: t('agendaPage.subject') || 'Materia', field: row => row.subject?.name || row.subject, sortable: true },
  { name: 'date', label: t('classRegister.dateLabel') || 'Data', field: 'date', sortable: true },
  { name: 'value', label: t('classRegister.tableHeaderGrade') || 'Voto', field: 'value', sortable: true },
  { name: 'type', label: t('classRegister.tableHeaderGradeType') || 'Tipo', field: 'type' }
])

onMounted(() => {
  // fetch() // Uncomment when API is ready
  grades.value = [
    { id: 1, subject: { name: 'Math' }, date: '2023-10-01', value: 8, type: 'Oral' },
    { id: 2, subject: { name: 'History' }, date: '2023-10-05', value: 5, type: 'Written' }
  ]
})

function getGradeColor(grade) {
  if (grade >= 6) return 'positive'
  return 'negative'
}
</script>
