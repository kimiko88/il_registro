<template>
  <q-page class="q-pa-md bg-indigo-10">
    <div class="row items-center q-mb-xl text-white">
      <div class="text-h3 text-weight-bolder glass-header q-pa-md rounded-borders shadow-10">
        <q-icon name="auto_awesome" color="amber" class="q-mr-sm" />
        Scrutinio Accademico
      </div>
      <q-space />
      <div class="row q-gutter-md glass-card q-pa-md rounded-borders">
        <q-select
          v-model="selectedClassId"
          :options="classOptions"
          label="Seleziona Classe"
          dark outlined dense
          style="min-width: 250px"
          emit-value map-options
          bg-color="white-1"
        />
        <q-btn-toggle
          v-model="semester"
          toggle-color="amber"
          flat dark
          :options="[{label: '1° Quad', value: 1}, {label: '2° Quad', value: 2}]"
        />
      </div>
    </div>

    <div v-if="!selectedClassId" class="flex flex-center" style="height: 50vh">
        <q-card class="glass-card text-center q-pa-xl text-white">
            <q-icon name="rocket_launch" size="80px" color="amber-2" class="q-mb-md" />
            <div class="text-h4 text-weight-light">Pronto per il Consiglio di Classe?</div>
            <div class="text-subtitle1 opacity-70">Seleziona una classe per iniziare il processo di scrutinio.</div>
        </q-card>
    </div>

    <q-card v-else class="glass-card shadow-24 rounded-borders overflow-hidden">
      <q-table
        :rows="matrix.students || []"
        :columns="columns"
        row-key="student_id"
        flat dark
        class="bg-transparent"
        :loading="loading"
        hide-bottom
        :pagination="{ rowsPerPage: 0 }"
      >
        <!-- Header: Subjects -->
        <template v-slot:header="props">
          <q-tr :props="props" class="bg-indigo-9">
            <q-th rowspan="2" align="left" class="text-h6">Studente</q-th>
            <q-th colspan="3" align="center" class="bg-indigo-7">Presenze</q-th>
            <q-th v-for="sub in matrix.subjects" :key="sub.id" align="center" class="subject-header">
              {{ sub.name }}
            </q-th>
            <q-th rowspan="2" align="center" class="bg-amber-9 text-black">Condotta</q-th>
            <q-th rowspan="2" align="center" class="bg-green-9">Esito</q-th>
            <q-th rowspan="2" align="center">Azioni</q-th>
          </q-tr>
          <q-tr :props="props" class="bg-indigo-8">
            <q-th align="center" class="text-caption text-indigo-2">Ass.</q-th>
            <q-th align="center" class="text-caption text-indigo-2">Rit.</q-th>
            <q-th align="center" class="text-caption text-indigo-2">Usc.</q-th>
            <q-th v-for="sub in matrix.subjects" :key="sub.id" align="center" class="text-caption text-indigo-2">
              Media Aggregata
            </q-th>
          </q-tr>
        </template>

        <!-- Body -->
        <template v-slot:body="props">
          <q-tr :props="props" class="hover-row">
            <q-td class="text-weight-bold text-h6 text-amber-1">{{ props.row.student_name }}</q-td>
            
            <!-- Attendance Stats -->
            <q-td align="center" class="text-indigo-2 text-weight-bold">{{ props.row.attendance_stats?.absences || 0 }}</q-td>
            <q-td align="center" class="text-indigo-2">{{ props.row.attendance_stats?.lates || 0 }}</q-td>
            <q-td align="center" class="text-indigo-2">{{ props.row.attendance_stats?.early_exits || 0 }}</q-td>

            <!-- Subject Averages -->
            <q-td v-for="sub in matrix.subjects" :key="sub.id" align="center">
              <!-- Average Reference -->
              <div class="text-caption text-grey-5 mb-1">
                Avg: {{ props.row.subject_data[sub.id]?.average?.toFixed(1) || '-' }}
              </div>
              
              <!-- Final Grade Input -->
              <q-input
                v-if="scrutinyData[props.row.student_id]"
                v-model.number="scrutinyData[props.row.student_id].grades[sub.id]"
                type="number"
                dense dark outlined
                input-class="text-center text-weight-bold"
                style="width: 45px"
                :bg-color="getAverageColor(scrutinyData[props.row.student_id].grades[sub.id])"
              />

              <div class="text-caption text-indigo-2 q-mt-xs" v-if="props.row.subject_data[sub.id]?.grade_count > 0">
                <q-icon name="grade" size="10px" /> {{ props.row.subject_data[sub.id].grade_count }}
              </div>
            </q-td>

            <!-- Conduct -->
            <q-td align="center" class="conduct-cell">
              <div class="text-caption text-grey-7 q-mb-xs">Condotta</div>
              <q-input
                v-if="scrutinyData[props.row.student_id]"
                v-model.number="scrutinyData[props.row.student_id].conduct_grade"
                type="number"
                dense
                outlined
                rounded
                input-class="text-center text-weight-bold"
                class="conduct-input"
                :bg-color="scrutinyData[props.row.student_id].conduct_grade < 6 ? 'red-1' : 'amber-1'"
              />
            </q-td>

            <!-- Decision -->
            <q-td align="center">
              <q-select
                v-if="scrutinyData[props.row.student_id]"
                v-model="scrutinyData[props.row.student_id].final_decision"
                :options="['Ammesso', 'Non Ammesso', 'Sospeso', 'Promosso', 'Respinto']"
                dense dark outlined
                options-dense
                class="decision-select"
              />
            </q-td>

            <q-td align="center">
              <q-btn fab-mini icon="save" color="amber" text-color="black" @click="saveStudentScrutiny(props.row.student_id)">
                <q-tooltip>Salva Singolo</q-tooltip>
              </q-btn>
            </q-td>
          </q-tr>
        </template>
      </q-table>
      
      <q-separator />
      <q-card-actions align="right" class="q-pa-md">
        <q-btn label="Salva Tutto" color="primary" icon="done_all" @click="saveAll" :loading="saving" />
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

// Map: StudentID -> { conduct_grade, final_decision, grades: { subject_id: final_grade } }
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
    
    // Initialize state BEFORE updating matrix.value to avoid render race conditions
    const newScrutinyData = {}
    const students = res.data.students || []
    const subjects = res.data.subjects || []

    students.forEach(s => {
      newScrutinyData[s.student_id] = {
        conduct_grade: s.record?.conduct_grade || 8,
        final_decision: s.record?.final_decision || 'Ammesso',
        grades: {}
      }
      // Populate final grades from record if exists, otherwise from averages
      subjects.forEach(sub => {
        const existing = s.record?.grades?.find(g => g.subject_id === sub.id)
        newScrutinyData[s.student_id].grades[sub.id] = existing ? existing.final_grade : Math.round(s.subject_data[sub.id]?.average || 6)
      })
    })

    // Clear old data and batch update state
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
    $q.notify({ type: 'positive', message: 'Dati salvati per lo studente', position: 'top' })
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
    $q.notify({ type: 'positive', message: 'Tutto salvato con successo!' })
  } finally {
    saving.value = false
  }
}

const getAverageColor = (avg) => {
  if (!avg) return 'bg-grey-7'
  if (avg < 5.5) return 'bg-insufficient'
  if (avg < 6) return 'bg-warning-grade'
  if (avg < 8) return 'bg-sufficient'
  return 'bg-excellent'
}
</script>

<style scoped>
.glass-header {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}
.glass-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(15px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
.grade-badge {
  width: 45px;
  height: 45px;
  line-height: 45px;
  border-radius: 12px;
  font-weight: 800;
  font-size: 1.1rem;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}
.bg-insufficient { background: linear-gradient(135deg, #f44336 0%, #d32f2f 100%); color: white; }
.bg-warning-grade { background: linear-gradient(135deg, #ff9800 0%, #f57c00 100%); color: white; }
.bg-sufficient { background: linear-gradient(135deg, #4caf50 0%, #388e3c 100%); color: white; }
.bg-excellent { background: linear-gradient(135deg, #2196f3 0%, #1976d2 100%); color: white; }

.hover-row:hover {
  background: rgba(255, 255, 255, 0.05) !important;
  transition: all 0.3s ease;
}
.subject-header {
  min-width: 100px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 0.7rem;
}
.decision-select {
  min-width: 130px;
}
.conduct-cell {
  background: rgba(255, 193, 7, 0.05);
  border-left: 2px solid #ffc107;
  min-width: 80px;
}
.conduct-input {
  width: 65px;
  margin: 0 auto;
}
.conduct-input :deep(.q-field__control) {
  height: 40px;
  background: white !important;
}
.opacity-70 { opacity: 0.7; }
</style>
