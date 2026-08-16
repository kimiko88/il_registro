<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">{{ t('nav.myGrades') }}</div>
       <q-btn icon="download" :label="t('gradesPage.printReport')" color="primary" :loading="downloading" @click="downloadReport" />
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Sidebar Controls -->
        <div class="col-12 col-md-4">
            <q-card class="q-mb-md shadow-1">
                <q-card-section>
                    <div class="text-h6 text-outfit text-weight-bold q-mb-md">Filtri</div>
                    <q-select v-model="filters.semester" :options="[1, 2]" label="Quadrimestre" outlined dense class="q-mb-sm" />
                    <q-select v-model="filters.period" :options="['Tutti', 'Ultimo Mese', 'Ultima Settimana']" label="Periodo" outlined dense class="q-mb-sm" />
                </q-card-section>
            </q-card>

            <q-card class="q-mb-md shadow-1">
               <q-card-section>
                   <div class="text-h6 text-outfit text-weight-bold q-mb-md">Andamento Medie</div>
                   <div v-for="sub in subjectAverages" :key="sub.name" class="q-mb-sm">
                       <div class="row justify-between text-caption">
                           <span class="text-weight-medium">{{ sub.name }}</span>
                           <span :class="{'text-green text-weight-bold': sub.avg !== '-' && Number(sub.avg) >= 6, 'text-red text-weight-bold': sub.avg === '-' || Number(sub.avg) < 6}">{{ sub.avg }}</span>
                       </div>
                       <q-linear-progress :value="sub.avg !== '-' ? Number(sub.avg)/10 : 0" :color="sub.avg !== '-' && Number(sub.avg) >= 6 ? 'green' : 'red'" />
                   </div>
               </q-card-section>
            </q-card>

            <!-- Simulator & Projection Card -->
            <q-card class="shadow-1">
                <q-card-section>
                    <div class="text-h6 text-weight-bold text-outfit row items-center q-mb-md">
                        <q-icon name="calculate" color="primary" class="q-mr-sm" />
                        Simulatore & Proiezioni
                    </div>

                    <!-- Sufficiency Warnings -->
                    <div v-if="subjectsBelowSufficiency.length > 0" class="q-mb-md">
                        <div class="text-caption text-weight-bold text-red-9 q-mb-sm">PER RAGGIUNGERE LA SUFFICIENZA (6.0):</div>
                        <q-list dense separator class="bg-red-1 rounded q-pa-xs">
                            <q-item v-for="item in subjectsBelowSufficiency" :key="item.name">
                                <q-item-section>
                                    <div class="text-caption text-weight-bold">{{ item.name }}</div>
                                    <div class="text-caption text-grey-8">Media attuale: {{ item.avg }}</div>
                                </q-item-section>
                                <q-item-section side>
                                    <q-badge color="negative" class="text-weight-bold q-pa-xs">
                                        Prossimo voto: {{ item.needed }}
                                    </q-badge>
                                </q-item-section>
                            </q-item>
                        </q-list>
                    </div>

                    <!-- Interactive Simulator -->
                    <q-separator class="q-my-md" />
                    <div class="text-caption text-weight-bold text-grey-7 q-mb-sm">SIMULA IL TUO PROSSIMO VOTO:</div>
                    <q-select
                        v-model="simSubject"
                        :options="subjectNames"
                        label="Seleziona Materia"
                        dense outlined
                        class="q-mb-sm"
                    />
                    <div class="row q-col-gutter-sm items-center">
                        <div class="col-6">
                            <q-input
                                v-model.number="simGrade"
                                type="number"
                                label="Voto ipotetico"
                                dense outlined
                                min="1" max="10" step="0.25"
                                aria-describedby="simulated-avg-output"
                                :rules="[v => (v >= 1 && v <= 10) || 'Inserisci un voto tra 1 e 10']"
                            />
                        </div>
                        <div class="col-6 text-center">
                            <div class="text-caption text-grey">Nuova Media</div>
                            <div id="simulated-avg-output" class="text-h6 text-weight-bold" :class="simulatedAverage >= 6 ? 'text-green' : 'text-red'">
                                {{ simulatedAverage }}
                            </div>
                        </div>
                    </div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Main Grade Table -->
        <div class="col-12 col-md-8">
            <q-table
              title="Registro Voti"
              :rows="filteredGrades"
              :columns="columns"
              row-key="id"
              :pagination="{ rowsPerPage: 10 }"
              :loading="gradesLoading"
              loading-label="Caricamento voti..."
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
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { gradeService } from '@/services/gradeService'
import api from '@/services/api'
import { useStudentStore } from '@/stores/student'

const $q = useQuasar()
const { t } = useI18n()
const studentStore = useStudentStore()
const gradesLoading = ref(true)

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
        const classId = studentStore.profile?.class_id
        if (!classId) return
        const { data } = await api.get(`/classes/${classId}/subjects`)
        if (data) {
            const map = {}
            data.forEach(s => map[s.id || s.subject_id] = s.name || s.subject_name)
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
                           subject: subjectsMap.value[g.subject_id] || g.subject_name || (g.subject_id && !g.subject_id.includes('-') ? g.subject_id : 'Materia sconosciuta'),
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
    } finally {
        gradesLoading.value = false
    }
}

const filteredGrades = computed(() => {
    let list = grades.value.filter(g => g.semester === filters.value.semester)
    if (filters.value.period === 'Ultimo Mese') {
        const monthAgo = new Date(); monthAgo.setMonth(monthAgo.getMonth() - 1);
        list = list.filter(g => new Date(g.date) >= monthAgo)
    } else if (filters.value.period === 'Ultima Settimana') {
        const weekAgo = new Date(); weekAgo.setDate(weekAgo.getDate() - 7);
        list = list.filter(g => new Date(g.date) >= weekAgo)
    }
    return list
})

const subjectAverages = computed(() => {
    const sums = {}
    const counts = {}
    filteredGrades.value.forEach(g => {
        if (g.value === 'A') return;
        const val = Number(g.value);
        if (isNaN(val)) return;
        if (!sums[g.subject]) { sums[g.subject] = 0; counts[g.subject] = 0; }
        sums[g.subject] += val;
        counts[g.subject]++;
    });
    return Object.keys(sums).map(sub => ({
        name: sub,
        avg: counts[sub] > 0 ? (sums[sub] / counts[sub]).toFixed(1) : '-'
    }))
})

const simSubject = ref(null)
const simGrade = ref(6)

const subjectNames = computed(() => {
    return subjectAverages.value.map(s => s.name)
})

const simulatedAverage = computed(() => {
    if (!simSubject.value || !simGrade.value) return '-'
    const targetSub = simSubject.value
    let sum = 0
    let count = 0
    
    filteredGrades.value.forEach(g => {
        if (g.subject === targetSub && g.value !== 'A') {
            const val = Number(g.value);
            if (!isNaN(val)) {
                sum += val
                count++
            }
        }
    })
    
    if (count === 0) return Number(simGrade.value).toFixed(2)
    sum += Number(simGrade.value)
    count++
    
    return (sum / count).toFixed(2)
})

const subjectsBelowSufficiency = computed(() => {
    return subjectAverages.value.filter(s => s.avg !== '-' && Number(s.avg) < 6.0).map(s => {
        let sum = 0
        let count = 0
        filteredGrades.value.forEach(g => {
            if (g.subject === s.name && g.value !== 'A') {
                const val = Number(g.value);
                if (!isNaN(val)) {
                    sum += val
                    count++
                }
            }
        })
        const needed = 6 * (count + 1) - sum
        return {
            name: s.name,
            avg: s.avg,
            needed: needed > 10 ? 'N/A (>10)' : needed.toFixed(2)
        }
    })
})

const downloading = ref(false)

watch(simSubject, () => {
    simGrade.value = 6
})

const downloadReport = async () => {
    downloading.value = true
    try {
        await gradeService.downloadReportCardPDF(filters.value.semester)
        $q.notify({ type: 'positive', message: 'Report PDF scaricato con successo' })
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore nel download del report PDF' })
    } finally {
        downloading.value = false
    }
}
</script>
