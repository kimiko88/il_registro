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
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { gradeService } from 'src/services/gradeService'
import adminService from 'src/services/adminService'
import { useStudentStore } from 'src/stores/student'

const $q = useQuasar()
const studentStore = useStudentStore()

const filters = ref({
    semester: 1,
    period: 'Tutti'
})

const columns = [
    { name: 'date', label: 'Data', align: 'left', field: 'date', sortable: true },
    { name: 'subject', label: 'Materia', align: 'left', field: 'subject', sortable: true },
    { name: 'evalType', label: 'Tipo Prova', align: 'left', field: 'evalType' },
    { name: 'type', label: 'Categoria', align: 'left', field: 'type' },
    { name: 'value', label: 'Voto', align: 'center', field: 'value', sortable: true },
    { name: 'desc', label: 'Argomento', align: 'left', field: 'description' }
]

const grades = ref([])
const subjectsMap = ref({})

onMounted(async () => {
    await studentStore.fetchProfile()
    await fetchSubjects()
    fetchMyGrades()
})

const fetchSubjects = async () => {
    try {
        const schoolId = studentStore.profile?.school_id
        if (!schoolId) return
        const { data } = await adminService.getSubjects(schoolId)
        if (data) {
            const map = {}
            data.forEach(s => map[s.id] = s.name)
            subjectsMap.value = map
        }
    } catch (e) {
        console.error('Error loading subjects', e)
    }
}

const mapEvalType = (type) => {
    const map = { Written: 'Scritto', Oral: 'Orale', Practical: 'Pratico' }
    return map[type] || type || '-'
}

const fetchMyGrades = async () => {
    try {
        const res = await gradeService.getMyGrades()
        const all = []
        if (res.data && res.data.semesters) {
            res.data.semesters.forEach(s => {
                if (s.grades) {
                   s.grades.forEach(g => {
                       all.push({
                           id: g.id,
                           date: g.date.split('T')[0],
                           subject: subjectsMap.value[g.subject_id] || g.subject_id,
                           evalType: mapEvalType(g.evaluation_type),
                           type: g.grade_type,
                           value: g.grade_value === -1 ? 'A' : g.grade_value,
                           description: g.description,
                           semester: g.semester
                       })
                   })
                }
            })
        }
        grades.value = all
    } catch (e) {
        console.error(e)
    }
}

const filteredGrades = computed(() => {
    let list = grades.value.filter(g => g.semester === filters.value.semester)
    if (filters.value.period === 'Ultimo Mese') {
        const monthAgo = new Date(); monthAgo.setMonth(monthAgo.getMonth() - 1);
        list = list.filter(g => new Date(g.date) >= monthAgo)
    }
    return list
})

const subjectAverages = computed(() => {
    const sums = {}
    const counts = {}
    filteredGrades.value.forEach(g => {
        if (g.value === 'A') return;
        if (!sums[g.subject]) { sums[g.subject] = 0; counts[g.subject] = 0; }
        sums[g.subject] += g.value;
        counts[g.subject]++;
    });
    return Object.keys(sums).map(sub => ({
        name: sub,
        avg: counts[sub] > 0 ? (sums[sub] / counts[sub]).toFixed(1) : '-'
    }))
})

const getGradeColor = (val) => {
    if (val === 'A') return 'grey'
    if (val >= 8) return 'green'
    if (val >= 6) return 'orange'
    return 'red'
}

const downloadReport = () => {
    $q.notify({ type: 'positive', message: 'Report PDF scaricato (simulato)' })
}
</script>
