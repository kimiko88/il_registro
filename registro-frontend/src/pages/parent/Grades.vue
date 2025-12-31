<template>
  <q-page class="q-pa-md bg-grey-1">
    <div v-if="!selectedChild" class="text-center q-pa-xl text-grey">
        <q-icon name="face" size="64px" />
        <div class="text-h6">Seleziona un figlio dalla Dashboard per vedere i voti.</div>
        <q-btn label="Vai alla Dashboard" color="primary" flat to="/parent" />
    </div>

    <div v-else>
        <div class="row items-center justify-between q-mb-md">
           <div class="text-h4">Voti di {{ selectedChild.name }}</div>
           <q-btn-toggle
              v-model="period"
              toggle-color="primary"
              :options="[{label: '1° Quad', value: 1}, {label: '2° Quad', value: 2}]"
              rounded unelevated
              class="bg-white border-primary"
           />
        </div>

        <!-- Grade List Grouped -->
        <div v-for="(subject, name) in gradesBySubject" :key="name" class="q-mb-md">
            <q-expansion-item
                class="shadow-1 overflow-hidden bg-white"
                style="border-radius: 8px"
                icon="book"
                :label="name"
                :caption="'Media: ' + subject.average"
                header-class="text-weight-medium"
            >
                <q-table
                    :rows="subject.grades"
                    :columns="columns"
                    hide-bottom
                    flat dense
                >
                    <template v-slot:body-cell-value="props">
                        <q-td :props="props">
                            <q-badge :color="getGradeColor(props.value)" class="text-subtitle2 q-pa-xs">
                                {{ props.value }}
                            </q-badge>
                        </q-td>
                    </template>
                </q-table>
            </q-expansion-item>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useParentStore } from 'src/stores/parent'

const parentStore = useParentStore()
const selectedChild = computed(() => parentStore.selectedChild)
const period = ref(1)

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
    { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
    { name: 'value', label: 'Voto', field: 'value', align: 'center', sortable: true },
    { name: 'notes', label: 'Note', field: 'notes', align: 'left' }
]

// Mock data
const rawGrades = [
    { id: 1, subject: 'Matematica', value: 8.5, date: '2025-01-10', type: 'Scritto', notes: '', semester: 1 },
    { id: 2, subject: 'Matematica', value: 7, date: '2025-01-20', type: 'Orale', notes: '', semester: 1 },
    { id: 3, subject: 'Storia', value: 6, date: '2025-01-15', type: 'Orale', notes: 'Interrogazione', semester: 1 },
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
    
    Object.keys(grouped).forEach(k => {
        grouped[k].average = (grouped[k].total / grouped[k].count).toFixed(1)
    })
    return grouped
})

const getGradeColor = (val) => {
    if (val < 6) return 'red'
    if (val < 8) return 'orange'
    return 'green'
}
</script>
