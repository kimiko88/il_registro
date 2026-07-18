<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Scrutinio Accademico</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Gestione voti finali e deliberazioni del Consiglio di Classe</p>
      </div>
      <div class="col-auto">
        <div class="row q-gutter-md glass-card q-pa-sm rounded-xl border-slate-200">
          <q-select
            v-model="selectedClassId"
            :options="classOptions"
            label="Classe"
            outlined dense
            style="min-width: 200px"
            emit-value map-options
            class="rounded-lg"
          />
          <q-btn-toggle
            v-model="semester"
            toggle-color="primary"
            flat
            class="rounded-lg border-slate-200"
            :options="[{label: '1° Quad', value: 1}, {label: '2° Quad', value: 2}]"
          />
        </div>
      </div>
    </div>

    <div v-if="!selectedClassId" class="flex flex-center" style="height: 60vh">
        <q-card class="glass-card text-center q-pa-xl rounded-2xl border-slate-100 shadow-soft">
            <q-icon name="rocket_launch" size="80px" color="primary" class="q-mb-md opacity-80" />
            <div class="text-h4 text-weight-bold text-slate-800">Consiglio di Classe</div>
            <div class="text-subtitle1 text-slate-500 q-mt-sm">Seleziona una classe per iniziare il processo di scrutinio.</div>
        </q-card>
    </div>

    <q-card v-else class="glass-card shadow-soft border-slate-100 overflow-hidden">
      <div class="scroll-container overflow-auto">
        <q-table
          :rows="matrix.students || []"
          :columns="columns"
          row-key="student_id"
          flat
          class="bg-transparent scrutiny-table"
          :loading="loading"
          hide-bottom
          :pagination="{ rowsPerPage: 0 }"
        >
          <!-- Header: Subjects -->
          <template v-slot:header="props">
            <q-tr :props="props" class="bg-slate-50">
              <q-th rowspan="2" align="left" class="text-weight-bold text-slate-700 sticky-col">Studente</q-th>
              <q-th colspan="3" align="center" class="bg-indigo-50 text-indigo-900 border-x">Presenze</q-th>
              <q-th v-for="sub in matrix.subjects" :key="sub.id" align="center" class="subject-header text-weight-bold text-slate-600">
                {{ sub.name }}
              </q-th>
              <q-th rowspan="2" align="center" class="bg-amber-50 text-amber-900 border-l text-weight-bold">Condotta</q-th>
              <q-th rowspan="2" align="center" class="bg-emerald-50 text-emerald-900 text-weight-bold">Esito</q-th>
              <q-th rowspan="2" align="center" class="text-slate-500">Azioni</q-th>
            </q-tr>
            <q-tr :props="props" class="bg-slate-50">
              <q-th align="center" class="text-caption text-indigo-400 border-l">Ass.</q-th>
              <q-th align="center" class="text-caption text-indigo-400">Rit.</q-th>
              <q-th align="center" class="text-caption text-indigo-400 border-r">Usc.</q-th>
              <q-th v-for="sub in matrix.subjects" :key="sub.id" align="center" class="text-caption text-slate-400">
                Media
              </q-th>
            </q-tr>
          </template>

          <!-- Body -->
          <template v-slot:body="props">
            <q-tr :props="props" class="hover-row">
              <q-td class="text-weight-bold text-slate-800 sticky-col bg-white">
                {{ props.row.student_name }}
              </q-td>
              
              <!-- Attendance Stats -->
              <q-td align="center" class="text-indigo-700 text-weight-medium border-l">{{ props.row.attendance_stats?.absences || 0 }}</q-td>
              <q-td align="center" class="text-slate-500">{{ props.row.attendance_stats?.lates || 0 }}</q-td>
              <q-td align="center" class="text-slate-500 border-r">{{ props.row.attendance_stats?.early_exits || 0 }}</q-td>

              <!-- Subject Averages -->
              <q-td v-for="sub in matrix.subjects" :key="sub.id" align="center" class="subject-cell">
                <div class="text-caption text-slate-400 q-mb-xs">
                  {{ props.row.subject_data[sub.id]?.average?.toFixed(1) || '-' }}
                </div>
                
                <q-input
                  v-if="scrutinyData[props.row.student_id]"
                  v-model.number="scrutinyData[props.row.student_id].grades[sub.id]"
                  type="number"
                  dense outlined
                  input-class="text-center text-weight-bold"
                  class="grade-input rounded-lg overflow-hidden"
                  :class="getGradeClass(scrutinyData[props.row.student_id].grades[sub.id])"
                />
              </q-td>

              <!-- Conduct -->
              <q-td align="center" class="bg-amber-50 border-l">
                <q-input
                  v-if="scrutinyData[props.row.student_id]"
                  v-model.number="scrutinyData[props.row.student_id].conduct_grade"
                  type="number"
                  dense outlined
                  input-class="text-center text-weight-bold"
                  class="conduct-input rounded-lg"
                  :class="scrutinyData[props.row.student_id].conduct_grade < 6 ? 'bg-negative text-white' : 'bg-amber-100 text-amber-900'"
                />
              </q-td>

              <!-- Decision -->
              <q-td align="center">
                <q-select
                  v-if="scrutinyData[props.row.student_id]"
                  v-model="scrutinyData[props.row.student_id].final_decision"
                  :options="['Ammesso', 'Non Ammesso', 'Sospeso', 'Promosso', 'Respinto']"
                  dense outlined
                  options-dense
                  class="decision-select rounded-lg"
                />
              </q-td>

              <q-td align="center">
                <q-btn flat round dense icon="save" color="primary" @click="saveStudentScrutiny(props.row.student_id)">
                  <q-tooltip>Salva Singolo</q-tooltip>
                </q-btn>
              </q-td>
            </q-tr>
          </template>
        </q-table>
      </div>
      
      <q-separator />
      <q-card-actions align="right" class="q-pa-md bg-transparent">
        <q-btn label="Salva Scrutinio Finale" color="primary" icon="done_all" class="q-px-lg rounded-lg shadow-sm" @click="saveAll" :loading="saving" />
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { scrutinyService } from 'src/services/scrutinyService'
import { useClassesStore } from 'src/stores/classes'
import { useAuthStore } from 'src/stores/auth'

const $q = useQuasar()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const selectedClassId = ref(null)
const semester = ref(1)
const loading = ref(false)
const saving = ref(false)
const matrix = ref({})
const classOptions = ref([])

const columns = [
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left' },
  { name: 'absences', label: 'Ass.', align: 'center' },
  { name: 'lates', label: 'Rit.', align: 'center' },
  { name: 'early_exits', label: 'Usc.', align: 'center' },
  { name: 'conduct', label: 'Condotta', align: 'center' },
  { name: 'decision', label: 'Esito', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const scrutinyData = reactive({})

onMounted(async () => {
  if (authStore.userRole === 'admin' || authStore.userRole === 'secretary') {
    await classesStore.fetchClasses()
  } else {
    await classesStore.fetchAssignedClasses()
  }
  
  classOptions.value = classesStore.classes.map(c => ({
    label: `${c.name}${c.section} - ${c.academic_year}`,
    value: c.id
  }))
  
  if (classOptions.value.length > 0) {
    selectedClassId.value = classOptions.value[0].value
  }
})

watch([selectedClassId, semester], () => {
  if (selectedClassId.value) fetchMatrix()
})

const fetchMatrix = async () => {
  loading.value = true
  try {
    const res = await scrutinyService.getMatrix(selectedClassId.value, semester.value)
    
    const newScrutinyData = {}
    const students = res.data.students || []
    const subjects = res.data.subjects || []

    students.forEach(s => {
      newScrutinyData[s.student_id] = {
        conduct_grade: s.record?.conduct_grade || 8,
        final_decision: s.record?.final_decision || 'Ammesso',
        grades: {}
      }
      subjects.forEach(sub => {
        const existing = s.record?.grades?.find(g => g.subject_id === sub.id)
        newScrutinyData[s.student_id].grades[sub.id] = existing ? existing.final_grade : Math.round(s.subject_data[sub.id]?.average || 6)
      })
    })

    for (const key in scrutinyData) delete scrutinyData[key]
    Object.assign(scrutinyData, newScrutinyData)
    matrix.value = res.data
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento matrice' })
  } finally {
    loading.value = false
  }
}

const saveStudentScrutiny = async (studentId) => {
  try {
    const data = scrutinyData[studentId]
    const payload = {
      student_id: studentId,
      class_id: selectedClassId.value,
      semester: semester.value,
      conduct_grade: data.conduct_grade,
      final_decision: data.final_decision,
      grades: Object.keys(data.grades).map(sid => ({
        subject_id: sid,
        final_grade: data.grades[sid]
      }))
    }
    await scrutinyService.save(payload)
    $q.notify({ type: 'positive', message: 'Dati salvati con successo', position: 'top' })
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
  }
}

const saveAll = async () => {
  saving.value = true
  try {
    for (const sid of Object.keys(scrutinyData)) {
      await saveStudentScrutiny(sid)
    }
    $q.notify({ type: 'positive', message: 'Scrutinio salvato correttamente' })
  } finally {
    saving.value = false
  }
}

const getGradeClass = (avg) => {
  if (!avg) return 'bg-slate-100'
  if (avg < 5.5) return 'bg-red-50 text-red-900'
  if (avg < 6) return 'bg-orange-50 text-orange-900'
  if (avg < 8) return 'bg-blue-50 text-blue-900'
  return 'bg-emerald-50 text-emerald-900'
}
</script>

<style scoped>
.scrutiny-table {
  min-width: 1200px;
}
.sticky-col {
  position: sticky;
  left: 0;
  z-index: 2;
  border-right: 1px solid var(--border-color);
}
.border-x {
  border-left: 1px solid var(--border-color);
  border-right: 1px solid var(--border-color);
}
.border-l {
  border-left: 1px solid var(--border-color);
}
.border-r {
  border-right: 1px solid var(--border-color);
}
.subject-header {
  min-width: 100px;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}
.subject-cell {
  min-width: 100px;
  border-left: 1px solid var(--border-color);
}
.grade-input {
  width: 50px;
  margin: 0 auto;
}
.grade-input :deep(.q-field__control) {
  height: 40px;
  padding: 0;
}
.conduct-input {
  width: 60px;
  margin: 0 auto;
}
.conduct-input :deep(.q-field__control) {
  height: 40px;
}
.decision-select {
  min-width: 140px;
}
.hover-row:hover {
  background: var(--bg-primary) !important;
}
.body--dark .sticky-col {
  background-color: var(--bg-secondary) !important;
}
.body--dark .bg-white {
  background-color: var(--bg-secondary) !important;
}
</style>
