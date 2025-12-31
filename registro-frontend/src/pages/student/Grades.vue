<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">I Miei Voti</div>
       <q-btn-toggle
          v-model="period"
          toggle-color="primary"
          :options="[{label: '1° Quad', value: 1}, {label: '2° Quad', value: 2}]"
          rounded
          unelevated
          class="bg-white border-primary"
       />
    </div>

    <!-- Summary Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-md-4">
            <q-card class="bg-primary text-white">
                <q-card-section>
                    <div class="text-subtitle1">Media Generale</div>
                    <div class="text-h2 text-weight-bolder">{{ overallAverage }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-4">
             <q-card class="bg-white text-primary">
                <q-card-section>
                    <div class="text-subtitle1">Voto Più Alto</div>
                    <div class="text-h3 text-weight-bold">{{ highestGrade }}</div>
                    <div class="text-caption text-grey">Matematica - 10/01</div>
                </q-card-section>
            </q-card>
        </div>
    </div>

    <!-- Grades List by Subject -->
    <div v-for="(subject, name) in gradesBySubject" :key="name" class="q-mb-md">
        <q-expansion-item
            class="shadow-1 overflow-hidden"
            style="border-radius: 8px"
            icon="school"
            :label="name"
            :caption="'Media: ' + subject.average"
            header-class="bg-white text-primary text-weight-medium"
            expand-icon-class="text-primary"
            default-opened
        >
            <q-card>
                <q-table
                    :rows="subject.grades"
                    :columns="columns"
                    hide-bottom
                    flat
                    dense
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
        </q-expansion-item>
    </div>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const period = ref(1)

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
    { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
    { name: 'value', label: 'Voto', field: 'value', align: 'center', sortable: true },
    { name: 'notes', label: 'Note', field: 'notes', align: 'left' }
]

// Mock data - would come from useGradesStore().fetchMyGrades()
const rawGrades = [
    { id: 1, subject: 'Matematica', value: 8.5, date: '2025-01-10', type: 'Scritto', notes: 'Compito sui limiti', semester: 1 },
    { id: 2, subject: 'Matematica', value: 7, date: '2025-01-20', type: 'Orale', notes: '', semester: 1 },
    { id: 3, subject: 'Storia', value: 6, date: '2025-01-15', type: 'Orale', notes: 'Interrogazione', semester: 1 },
    { id: 4, subject: 'Inglese', value: 9, date: '2025-01-12', type: 'Scritto', notes: 'Essay', semester: 1 },
]

const gradesBySubject = computed(() => {
    const grouped = {}
    rawGrades.filter(g => g.semester === period.value).forEach(g => {
        if (!grouped[g.subject]) {
            grouped[g.subject] = { grades: [], total: 0, count: 0, average: 0 }
        }
        grouped[g.subject].grades.push(g)
        grouped[g.subject].total += g.value
        grouped[g.subject].count++
    })
    
    // Calculate averages
    Object.keys(grouped).forEach(k => {
        grouped[k].average = (grouped[k].total / grouped[k].count).toFixed(1)
    })
    return grouped
})

const overallAverage = computed(() => {
    const subjects = Object.values(gradesBySubject.value)
    if (subjects.length === 0) return '-'
    const sum = subjects.reduce((acc, curr) => acc + parseFloat(curr.average), 0)
    return (sum / subjects.length).toFixed(2)
})

const highestGrade = computed(() => {
    return Math.max(...rawGrades.map(g => g.value))
})

const getGradeColor = (val) => {
    if (val < 6) return 'red'
    if (val < 8) return 'orange'
    return 'green'
}
</script>
