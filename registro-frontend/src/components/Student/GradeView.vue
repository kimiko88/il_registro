<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">My Grades</div>

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
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const { data: grades, loading, fetch } = useApi('/grades')

const columns = [
  { name: 'subject', label: 'Subject', field: row => row.subject.name, sortable: true },
  { name: 'date', label: 'Date', field: 'date', sortable: true },
  { name: 'value', label: 'Grade', field: 'value', sortable: true },
  { name: 'type', label: 'Type', field: 'type' }
]

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
