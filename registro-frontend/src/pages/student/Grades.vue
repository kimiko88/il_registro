<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">I Miei Voti</div>
       <q-btn icon="download" label="Scarica Report" color="primary" @click="downloadReport" />
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Sidebar Controls -->
        <div class="col-12 col-md-3">
            <q-card class="q-mb-md">
                <q-card-section>
                    <div class="text-h6">Filtri</div>
                    <q-select v-model="filters.semester" :options="[1, 2]" label="Quadrimestre" outlined class="q-mb-sm" />
                    <q-select v-model="filters.period" :options="['Tutti', 'Ultimo Mese', 'Ultima Settimana']" label="Periodo" outlined class="q-mb-sm" />
                </q-card-section>
            </q-card>

            <q-card>
               <q-card-section>
                   <div class="text-h6">Andamento</div>
                   <!-- Simple CSS Bar Chart Fallback/Placeholder if ChartJS setup is complex in one file -->
                   <div v-for="sub in subjectAverages" :key="sub.name" class="q-mb-sm">
                       <div class="row justify-between text-caption">
                           <span>{{ sub.name }}</span>
                           <span :class="{'text-green': sub.avg>=6, 'text-red': sub.avg<6}">{{ sub.avg }}</span>
                       </div>
                       <q-linear-progress :value="sub.avg/10" :color="sub.avg>=6?'green':'red'" />
                   </div>
               </q-card-section>
            </q-card>
        </div>

        <!-- Main Grade Table -->
        <div class="col-12 col-md-9">
            <q-table
              title="Registro Voti"
              :rows="filteredGrades"
              :columns="columns"
              row-key="id"
              :pagination="{ rowsPerPage: 10 }"
              flat bordered
            >
                <template v-slot:body-cell-value="props">
                    <q-td :props="props">
                        <q-badge :color="getGradeColor(props.value)" class="text-body2 q-px-sm">
                            {{ props.value }}
                        </q-badge>
                    </q-td>
                </template>
            </q-table>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const filters = ref({
    semester: 1,
    period: 'Tutti'
})

const columns = [
    { name: 'date', label: 'Data', align: 'left', field: 'date', sortable: true },
    { name: 'subject', label: 'Materia', align: 'left', field: 'subject', sortable: true },
    { name: 'type', label: 'Tipo', align: 'left', field: 'type' },
    { name: 'value', label: 'Voto', align: 'center', field: 'value', sortable: true },
    { name: 'desc', label: 'Argomento', align: 'left', field: 'description' }
]

const grades = ref([
    { id: 1, date: '2024-10-15', subject: 'Matematica', type: 'Scritto', value: 8.5, description: 'Equazioni', semester: 1 },
    { id: 2, date: '2024-10-18', subject: 'Storia', type: 'Orale', value: 7, description: 'Interrogazione', semester: 1 },
    { id: 3, date: '2024-11-02', subject: 'Fisica', type: 'Pratico', value: 5.5, description: 'Lab', semester: 1 },
    { id: 4, date: '2025-02-10', subject: 'Matematica', type: 'Scritto', value: 9, description: 'Funzioni', semester: 2 },
])

const filteredGrades = computed(() => {
    return grades.value.filter(g => g.semester === filters.value.semester)
})

const subjectAverages = computed(() => {
    // simplified calculation
    const sums = {}
    const counts = {}
    filteredGrades.value.forEach(g => {
        if (!sums[g.subject]) { sums[g.subject] = 0; counts[g.subject] = 0; }
        sums[g.subject] += g.value;
        counts[g.subject]++;
    });
    return Object.keys(sums).map(sub => ({
        name: sub,
        avg: (sums[sub] / counts[sub]).toFixed(1)
    }))
})

const getGradeColor = (val) => {
    if (val >= 8) return 'green'
    if (val >= 6) return 'orange'
    return 'red'
}

const downloadReport = () => {
    $q.notify({ type: 'positive', message: 'Report PDF scaricato (simulato)' })
}
</script>
