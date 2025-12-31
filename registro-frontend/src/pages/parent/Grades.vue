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
        :rows="mockGrades"
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
import { ref, computed } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

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

const mockGrades = [
  { id: 1, date: '2024-12-20', subject: 'Matematica', type: 'Scritto', value: 5.0, notes: 'Equazioni' },
  { id: 2, date: '2024-12-18', subject: 'Storia', type: 'Orale', value: 7.5, notes: 'Interrogazione' },
  { id: 3, date: '2024-12-15', subject: 'Inglese', type: 'Scritto', value: 8.0, notes: 'Grammar Test' },
  { id: 4, date: '2024-12-10', subject: 'Fisica', type: 'Pratico', value: 6.5, notes: 'Laboratorio' },
  { id: 5, date: '2024-11-28', subject: 'Italiano', type: 'Scritto', value: 6.0, notes: 'Tema' },
]

function getGradeColor(val) {
  if (val < 6) return 'negative';
  if (val >= 8) return 'positive';
  return 'orange';
}

function downloadReport() {
  $q.notify({ type: 'info', message: 'Download avviato...' })
}
</script>
