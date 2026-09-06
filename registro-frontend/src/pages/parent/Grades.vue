<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5 text-weight-bold text-slate-800">
        Voti: {{ selectedChild ? `${selectedChild.first_name} ${selectedChild.last_name}` : '...' }}
      </div>
      <q-btn flat icon="download" label="Scarica Pagella" color="primary" @click="downloadReport" />
    </div>

    <!-- Child Selector Warning -->
    <q-card v-if="parentStore.children.length === 0" class="text-center q-pa-lg bg-warning text-white q-mb-md">
      ⚠️ Nessun alunno associato al tuo profilo. Contatta la segreteria.
    </q-card>
    <q-card v-else-if="!selectedChild" class="text-center q-pa-lg bg-warning text-white q-mb-md">
      Seleziona un figlio dal menu in alto per visualizzare la situazione voti
    </q-card>

    <div v-else class="row q-col-gutter-lg">
        <!-- Sidebar Controls -->
        <div class="col-12 col-md-4">
            <q-card class="q-mb-md shadow-1">
                <q-card-section>
                    <div class="text-h6 text-outfit text-weight-bold q-mb-md">Filtri</div>
                    <q-select 
                      dense 
                      outlined 
                      v-model="period" 
                      :options="['Primo Quadrimestre', 'Secondo Quadrimestre']" 
                      label="Periodo" 
                      class="bg-white" 
                    />
                </q-card-section>
            </q-card>

            <q-card class="q-mb-md shadow-1">
               <q-card-section>
                   <div class="text-h6 text-outfit text-weight-bold q-mb-md">Andamento Medie</div>
                   <div v-for="sub in subjectAverages" :key="sub.name" class="q-mb-sm">
                       <div class="row justify-between items-center text-caption">
                           <span class="text-weight-medium">{{ sub.name }}</span>
                           <div class="row items-center">
                               <q-icon :name="getTrendIcon(sub.trend)" :color="getTrendColor(sub.trend)" size="16px" class="q-mr-xs">
                                   <q-tooltip>Trend: {{ sub.trend === 'up' ? 'In miglioramento' : (sub.trend === 'down' ? 'In calo' : 'Stabile') }}</q-tooltip>
                               </q-icon>
                               <span :class="{'text-positive text-weight-bold': Number(sub.avg)>=6, 'text-negative text-weight-bold': Number(sub.avg)<6 || sub.avg==='-'}">{{ sub.avg }}</span>
                           </div>
                       </div>
                       <q-linear-progress :value="sub.avg !== '-' ? Number(sub.avg)/10 : 0" :color="Number(sub.avg)>=6?'positive':'negative'" class="rounded-borders" />
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
                    <div class="text-caption text-weight-bold text-grey-7 q-mb-sm">SIMULA PROSSIMO VOTO:</div>
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
                            />
                        </div>
                        <div class="col-6 text-center">
                            <div class="text-caption text-grey">Nuova Media</div>
                            <div class="text-h6 text-weight-bold" :class="Number(simulatedAverage) >= 6 ? 'text-green' : 'text-red'">
                                {{ simulatedAverage }}
                            </div>
                        </div>
                    </div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Main Grade Table -->
        <div class="col-12 col-md-8">
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
                    <GradeBadge :value="props.value" />
                  </q-td>
                </template>
              </q-table>
            </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useParentStore } from '@/stores/parent'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { gradeService } from '@/services/gradeService'
import adminService from '@/services/adminService'
import GradeBadge from '@/components/Common/GradeBadge.vue'

const $q = useQuasar()
const { t } = useI18n()
const parentStore = useParentStore()
const authStore = useAuthStore()
const { selectedChild } = storeToRefs(parentStore)

const period = ref('Primo Quadrimestre')

watch(period, () => {
    simSubject.value = null
})

const columns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'subject', label: 'Materia', field: 'subject', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'value', label: 'Voto', field: 'value', align: 'center', sortable: true },
  { name: 'notes', label: 'Note', field: 'notes', align: 'left' }
]

const gradesData = ref(null)
const subjectsMap = ref({})

// Compute grades based on selected period
const currentGrades = computed(() => {
    if (!gradesData.value || !gradesData.value.semesters) return []
    
    const semNum = period.value === 'Primo Quadrimestre' ? 1 : 2
    const semData = gradesData.value.semesters.find(s => s.semester === semNum)
    if (!semData || !semData.grades) return []
    
    return semData.grades.map(g => {
        let formattedDate = g.date
        try {
            const parts = g.date.split('T')[0].split('-')
            if (parts.length === 3) {
                formattedDate = `${parts[2].padStart(2, '0')}/${parts[1].padStart(2, '0')}/${parts[0]}`
            }
        } catch {
            formattedDate = g.date.split('T')[0]
        }
        return {
            id: g.id,
            date: formattedDate,
            subject: subjectsMap.value[g.subject_id] || g.subject_id,
            type: g.grade_type,
            value: g.grade_value === -1 ? 'A' : g.grade_value,
            notes: g.description
        }
    })
})

onMounted(() => {
    fetchSubjects()
    if (selectedChild.value) {
        fetchGrades()
    }
})

watch(selectedChild, (newVal) => {
    if (newVal) {
        fetchSubjects()
        fetchGrades()
    }
})

const fetchSubjects = async () => {
    try {
        const schoolId = authStore.user?.school_id ||
                         authStore.user?.schoolId ||
                         selectedChild.value?.school_id
        if (!schoolId) return
        const { data } = await adminService.getSubjects(schoolId)
        if (data) {
            const map = {}
            data.forEach(s => map[s.id] = s.name)
            subjectsMap.value = map
        }
    } catch (e) {
        console.error("Error loading subjects", e)
    }
}

const fetchGrades = async () => {
    try {
        const res = await gradeService.getChildGrades(selectedChild.value.id)
        gradesData.value = res.data
    } catch (e) {
        console.error(e)
    }
}

function getSubjectTrend(subjectName) {
  const grades = currentGrades.value
    .filter(g => g.subject === subjectName && g.value !== 'A')
    .sort((a, b) => new Date(a.date) - new Date(b.date))
  if (grades.length < 2) return 'flat'
  const last = Number(grades[grades.length - 1].value)
  const prev = Number(grades[grades.length - 2].value)
  if (last > prev) return 'up'
  if (last < prev) return 'down'
  return 'flat'
}

function getTrendIcon(trend) {
  if (trend === 'up') return 'trending_up'
  if (trend === 'down') return 'trending_down'
  return 'trending_flat'
}

function getTrendColor(trend) {
  if (trend === 'up') return 'positive'
  if (trend === 'down') return 'negative'
  return 'grey-6'
}

const subjectAverages = computed(() => {
    const sums = {}
    const counts = {}
    currentGrades.value.forEach(g => {
        if (g.value === 'A') return
        if (!sums[g.subject]) { sums[g.subject] = 0; counts[g.subject] = 0; }
        sums[g.subject] += Number(g.value)
        counts[g.subject]++
    })
    return Object.keys(sums).map(sub => ({
        name: sub,
        avg: counts[sub] > 0 ? (sums[sub] / counts[sub]).toFixed(1) : '-',
        trend: getSubjectTrend(sub)
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
    
    currentGrades.value.forEach(g => {
        if (g.subject === targetSub && g.value !== 'A') {
            sum += Number(g.value)
            count++
        }
    })
    
    const clampedGrade = Math.min(10, Math.max(1, Number(simGrade.value)))
    if (count === 0) return clampedGrade.toFixed(2)
    sum += clampedGrade
    count++
    
    return (sum / count).toFixed(2)
})

const subjectsBelowSufficiency = computed(() => {
    return subjectAverages.value.filter(s => s.avg !== '-' && Number(s.avg) < 6.0).map(s => {
        let sum = 0
        let count = 0
        currentGrades.value.forEach(g => {
            if (g.subject === s.name && g.value !== 'A') {
                sum += Number(g.value)
                count++
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


const downloadReport = async () => {
  try {
    const semNum = period.value === 'Primo Quadrimestre' ? 1 : 2
    await gradeService.downloadReportCardPDF(semNum)
    $q.notify({ type: 'positive', message: 'Report PDF scaricato con successo' })
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore nel download del report PDF' })
  }
}
</script>
