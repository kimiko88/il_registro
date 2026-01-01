<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5 text-weight-bold text-slate-800">
        Voti: {{ selectedChild?.firstName || '...' }}
      </div>
      <q-btn flat icon="download" label="Scarica Pagella" color="primary" @click="downloadReport" />
    </div>

    <!-- Filters -->
    <div class="row q-gutter-sm q-mb-md">
       <q-select 
         dense 
         outlined 
         v-model="period" 
         :options="['Primo Quadrimestre', 'Secondo Quadrimestre']" 
         label="Periodo" 
         class="bg-white" 
         style="width: 200px" 
       />
    </div>

    <!-- Grades Table -->
    <q-card class="shadow-sm rounded-lg">
      <q-table
        :rows="currentGrades"
        :columns="columns"
        row-key="id"
        flat
        bordered
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:body-cell-value="props">
          <q-td :props="props">
            <q-badge :color="getGradeColor(props.value)" class="text-subtitle2 q-pa-xs">
              {{ props.value }}
            </q-badge>
          </q-td>
        </template>
      </q-table>
    </q-card>
    
    <!-- Child Selector Warning -->
    <div v-if="!selectedChild" class="fixed-bottom q-pa-md bg-warning text-white text-center">
      Seleziona un figlio dalla Dashboard per vedere i dati corretti.
    </div>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { gradeService } from 'src/services/gradeService'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const period = ref('Primo Quadrimestre')

const columns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'subject', label: 'Materia', field: 'subject', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'value', label: 'Voto', field: 'value', align: 'center', sortable: true },
  { name: 'notes', label: 'Note', field: 'notes', align: 'left' }
]

const gradesData = ref(null)

// Compute grades based on selected period
const currentGrades = computed(() => {
    if (!gradesData.value || !gradesData.value.semesters) return []
    
    const semNum = period.value === 'Primo Quadrimestre' ? 1 : 2
    const semData = gradesData.value.semesters.find(s => s.semester === semNum)
    if (!semData || !semData.grades) return []
    
    return semData.grades.map(g => ({
        id: g.id,
        date: g.date.split('T')[0],
        subject: g.subject_id, // TODO: Map to Name
        type: g.grade_type,
        value: g.grade_value,
        notes: g.description
    }))
})

onMounted(() => {
    if (selectedChild.value) {
        fetchGrades()
    }
})

watch(selectedChild, (val) => {
    if (val) fetchGrades()
})

const fetchGrades = async () => {
    try {
        const res = await gradeService.getChildGrades(selectedChild.value.id)
        gradesData.value = res.data
    } catch (e) {
        console.error(e)
        // Optionally notify error
    }
}

function getGradeColor(val) {
  if (val < 6) return 'negative';
  if (val >= 8) return 'positive';
  return 'orange';
}

function downloadReport() {
  $q.notify({ type: 'info', message: 'Download avviato...' })
}
</script>
