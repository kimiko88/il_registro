<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Scrutinio Accademico & Differito</h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">Gestione voti finali, argomenti delle carenze e recupero debiti formativi</p>
      </div>
      <div class="col-auto">
        <div class="row q-gutter-md glass-card q-pa-sm rounded-xl border-slate-200 items-center">
          <q-btn v-if="selectedClassId" color="negative" icon="lock" label="Chiudi Scrutinio" unelevated @click="closeScrutiny" />
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
            v-model="period"
            toggle-color="primary"
            flat
            class="rounded-lg border-slate-200"
            :options="periodOptions"
          />
        </div>
      </div>
    </div>

    <!-- Banner Sessione Scrutinio Differito -->
    <q-banner v-if="period === 3" rounded class="bg-deep-orange-1 text-deep-orange-10 border border-deep-orange-3 q-mb-md q-pa-md shadow-xs">
      <template v-slot:avatar>
        <q-icon name="event_repeat" color="deep-orange-8" size="md" />
      </template>
      <div class="text-weight-bold text-subtitle1">{{ $t('help.teacher.scrutiny.deferredBannerTitle') }}</div>
      <div class="text-body2">{{ $t('help.teacher.scrutiny.deferredBannerBody') }}</div>
    </q-banner>

    <div v-if="!selectedClassId" class="flex flex-center" style="height: 60vh">
        <q-card class="glass-card text-center q-pa-xl rounded-2xl border-slate-100 shadow-soft">
            <q-icon name="rocket_launch" size="80px" color="primary" class="q-mb-md opacity-80" />
            <div class="text-h4 text-weight-bold text-slate-800">Consiglio di Classe</div>
            <div class="text-subtitle1 text-slate-500 q-mt-sm">
              {{ classOptions.length === 0 ? "Non risulti coordinatore di alcuna classe per lo scrutinio." : "Seleziona una classe per iniziare il processo di scrutinio." }}
            </div>
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
                <q-badge v-if="['Sospeso', 'Giudizio Sospeso'].includes(props.row.record?.final_decision)" color="deep-orange" class="q-ml-xs">
                  Sospeso
                </q-badge>
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
                  :options="['Ammesso', 'Non Ammesso', 'Sospeso', 'Promosso', 'Respinto', 'Promosso con debiti saldati']"
                  dense outlined
                  options-dense
                  class="decision-select rounded-lg"
                />
              </q-td>

              <q-td align="center" class="q-gutter-xs">
                <q-btn flat round dense icon="save" color="primary" @click="saveStudentScrutiny(props.row.student_id)">
                  <q-tooltip>Salva Singolo</q-tooltip>
                </q-btn>
                <q-btn flat round dense icon="warning" color="amber-9" @click="openDeficiencyModal(props.row)">
                  <q-tooltip>Argomenti Carenze & Recuperi</q-tooltip>
                </q-btn>
                <q-btn v-if="period === 2 || period === 3 || ['Sospeso', 'Giudizio Sospeso'].includes(props.row.record?.final_decision)" flat round dense icon="event_repeat" color="deep-orange" @click="openDeferredModal(props.row)">
                  <q-tooltip>Scrutinio Differito (Esami Recupero)</q-tooltip>
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

    <!-- Dialog: Argomenti Carenze & Recuperi -->
    <q-dialog v-model="showDeficiencyModal">
      <q-card style="min-width: 550px" class="rounded-xl">
        <q-card-section class="bg-amber-700 text-white row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">Argomenti Carenze — {{ selectedStudent?.student_name }}</div>
            <div class="text-caption">Indicazione delle lacune da recuperare (visibile a studente e genitori)</div>
          </div>
          <q-btn flat round icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-4">
          <q-select
            v-model="deficiencyForm.subject_id"
            :options="matrix.subjects"
            option-label="name"
            option-value="id"
            label="Materia della Carenza"
            outlined dense emit-value map-options
          />

          <q-input
            v-model="deficiencyForm.topics"
            type="textarea"
            rows="3"
            outlined
            label="Argomenti della Carenza / Lacune Specifiche *"
            hint="Es. Equazioni di 2° grado, Sintassi del periodo, Verbi irregolari"
          />

          <q-select
            v-model="deficiencyForm.recovery_mode"
            :options="[
              { label: 'Studio Individuale', value: 'studio_individuale' },
              { label: 'Corso di Recupero Estivo', value: 'corso_recupero' },
              { label: 'Sportello Didattico', value: 'sportello_didattico' }
            ]"
            label="Modalità di Recupero"
            outlined dense emit-value map-options
          />

          <q-select
            v-model="deficiencyForm.status"
            :options="[
              { label: 'Da Recuperare', value: 'da_recuperare' },
              { label: 'In Corso', value: 'in_corso' },
              { label: 'Recuperato', value: 'recuperato' },
              { label: 'Non Recuperato', value: 'non_recuperato' }
            ]"
            label="Stato Recupero"
            outlined dense emit-value map-options
          />

          <div class="row q-col-gutter-md">
            <div class="col-6">
              <q-input v-model.number="deficiencyForm.recovery_grade" type="number" step="0.5" label="Voto Prova di Recupero" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="deficiencyForm.recovery_date" type="date" label="Data Prova Recupero" outlined dense stack-label />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="amber-9" label="Salva Carenza" unelevated @click="saveDeficiency" :loading="savingDeficiency" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Scrutinio Differito (Saldo Debiti Formativi) -->
    <q-dialog v-model="showDeferredModal">
      <q-card style="min-width: 600px" class="rounded-xl">
        <q-card-section class="bg-deep-orange-8 text-white row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">{{ $t('help.teacher.scrutiny.deferredModalTitle', { name: selectedStudent?.student_name }) }}</div>
            <div class="text-caption">{{ $t('help.teacher.scrutiny.deferredModalSubtitle') }}</div>
          </div>
          <q-btn flat round icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-4">
          <div v-if="studentDeficienciesList.length === 0" class="text-slate-500 text-center q-pa-md">
            Nessuna carenza o debito salvato per questo studente.
          </div>

          <div v-for="def in studentDeficienciesList" :key="def.id" class="p-3 border rounded-lg bg-slate-50 space-y-2">
            <div class="row items-center justify-between">
              <div class="text-weight-bold text-slate-800">{{ def.subject_name || 'Materia' }}</div>
              <q-badge :color="def.status === 'recuperato' ? 'positive' : 'negative'">{{ def.status }}</q-badge>
            </div>
            <div class="text-caption text-slate-600"><strong>Argomenti:</strong> {{ def.topics }}</div>

            <div class="row q-col-gutter-md q-pt-xs">
              <div class="col-6">
                <q-select
                  v-model="def.status"
                  :options="[
                    { label: 'Recuperato (Debito Saldato)', value: 'recuperato' },
                    { label: 'Non Recuperato', value: 'non_recuperato' }
                  ]"
                  label="Esito Verifica"
                  outlined dense emit-value map-options
                />
              </div>
              <div class="col-6">
                <q-input v-model.number="def.recovery_grade" type="number" step="0.5" label="Voto Prova Recupero" outlined dense />
              </div>
            </div>
          </div>

          <q-separator />

          <q-select
            v-model="deferredForm.final_decision"
            :options="[
              { label: $t('help.teacher.scrutiny.promotedDebtsCleared'), value: 'promosso_con_debiti_saldati' },
              { label: $t('help.teacher.scrutiny.notPromotedDebtsNotCleared'), value: 'non_promosso' }
            ]"
            label="Delibera Finale Scrutinio Differito *"
            outlined dense emit-value map-options
          />

          <q-input v-model="deferredForm.notes" type="textarea" rows="2" :label="$t('help.teacher.scrutiny.deferredNotes')" outlined />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="deep-orange-8" :label="$t('help.teacher.scrutiny.deliberateDeferred')" unelevated @click="saveDeferredScrutiny" :loading="savingDeferred" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { scrutinyService } from '@/services/scrutinyService'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'
import { useSchoolYearStore } from '@/stores/schoolYear'

const $q = useQuasar()
const { t } = useI18n()
const authStore = useAuthStore()
const classesStore = useClassesStore()
const schoolYearStore = useSchoolYearStore()

const selectedClassId = ref(null)
const period = ref(1)
const periodOptions = computed(() => [
  { label: t('help.teacher.scrutiny.period1'), value: 1 },
  { label: t('help.teacher.scrutiny.period2'), value: 2 },
  { label: t('help.teacher.scrutiny.deferredScrutiny'), value: 3 }
])
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

// Deficiency & Deferred Scrutiny state
const showDeficiencyModal = ref(false)
const showDeferredModal = ref(false)
const selectedStudent = ref(null)
const savingDeficiency = ref(false)
const savingDeferred = ref(false)
const studentDeficienciesList = ref([])

const deficiencyForm = reactive({
  id: '',
  student_id: '',
  class_id: '',
  subject_id: '',
  semester: 1,
  topics: '',
  recovery_mode: 'studio_individuale',
  status: 'da_recuperare',
  recovery_grade: null,
  recovery_date: '',
  notes: ''
})

const deferredForm = reactive({
  student_id: '',
  class_id: '',
  final_decision: 'promosso_con_debiti_saldati',
  notes: ''
})

const loadClasses = async () => {
  if (['admin', 'secretary', 'principal', 'vice_principal'].includes(authStore.userRole)) {
    await classesStore.fetchClasses()
  } else {
    await classesStore.fetchAssignedClasses(schoolYearStore.selectedSchoolYear)
  }
  
  let availableClasses = classesStore.classes || []
  if (authStore.userRole === 'teacher') {
    const currentUserId = authStore.user?.id
    availableClasses = availableClasses.filter(c => c.coordinator_id === currentUserId)
  }

  classOptions.value = availableClasses.map(c => ({
    label: `${c.name}${c.section} - ${c.academic_year || schoolYearStore.selectedSchoolYear}`,
    value: c.id
  }))
  
  if (classOptions.value.length > 0) {
    selectedClassId.value = classOptions.value[0].value
  } else {
    selectedClassId.value = null
  }
}

onMounted(async () => {
  await loadClasses()

  try {
    const res = await api.get('/school-calendar/periods')
    if (res.data && res.data.length > 0) {
      periodOptions.value = res.data.map((p, index) => ({
        label: p.name,
        value: p.period || p.value || (index + 1)
      }))
    }
  } catch { /* fallback */ }
})

watch(() => schoolYearStore.selectedSchoolYear, () => {
  loadClasses()
})

watch([selectedClassId, period], () => {
  if (selectedClassId.value) fetchMatrix()
})

const fetchMatrix = async () => {
  loading.value = true
  try {
    const fetchSemester = period.value === 3 ? 2 : period.value
    const res = await scrutinyService.getMatrix(selectedClassId.value, fetchSemester)
    
    const newScrutinyData = {}
    const students = res.data.students || []
    const subjects = res.data.subjects || []

    students.forEach(s => {
      newScrutinyData[s.student_id] = {
        conduct_grade: s.record?.conduct_grade ?? null,
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

const saveStudentScrutiny = async (studentId, silent = false) => {
  try {
    const data = scrutinyData[studentId]
    if (data.conduct_grade == null || data.conduct_grade <= 0) {
      if (!silent) {
        $q.notify({ type: 'warning', message: 'Impostare il voto di condotta prima di salvare.' })
      }
      return false
    }
    const payload = {
      student_id: studentId,
      class_id: selectedClassId.value,
      semester: period.value,
      conduct_grade: data.conduct_grade,
      final_decision: data.final_decision,
      grades: Object.keys(data.grades).map(sid => ({
        subject_id: sid,
        final_grade: data.grades[sid]
      }))
    }
    await scrutinyService.save(payload)
    if (!silent) {
      $q.notify({ type: 'positive', message: 'Dati salvati con successo', position: 'top' })
    }
    return true
  } catch (e) {
    if (!silent) {
      $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
    }
    return false
  }
}

const saveAll = async () => {
  saving.value = true
  let successCount = 0
  let totalCount = 0
  try {
    const studentIds = Object.keys(scrutinyData)
    totalCount = studentIds.length
    for (const sid of studentIds) {
      const ok = await saveStudentScrutiny(sid, true)
      if (ok) successCount++
    }
    if (successCount === totalCount && totalCount > 0) {
      $q.notify({ type: 'positive', message: 'Scrutinio salvato correttamente' })
    } else {
      $q.notify({ type: 'warning', message: `Salvataggio parziale: salvati ${successCount} su ${totalCount} studenti. Verificare la condotta.` })
    }
  } finally {
    saving.value = false
  }
}

const openDeficiencyModal = (row) => {
  selectedStudent.value = row
  deficiencyForm.student_id = row.student_id
  deficiencyForm.class_id = selectedClassId.value
  deficiencyForm.semester = period.value
  deficiencyForm.topics = ''
  deficiencyForm.subject_id = matrix.value.subjects?.[0]?.id || ''
  deficiencyForm.status = 'da_recuperare'
  deficiencyForm.recovery_mode = 'studio_individuale'
  deficiencyForm.recovery_grade = null
  deficiencyForm.recovery_date = ''
  showDeficiencyModal.value = true
}

const saveDeficiency = async () => {
  if (!deficiencyForm.topics) {
    $q.notify({ type: 'warning', message: 'Inserire gli argomenti della carenza' })
    return
  }
  savingDeficiency.value = true
  try {
    await scrutinyService.saveDeficiency(deficiencyForm)
    $q.notify({ type: 'positive', message: 'Argomenti della carenza salvati con successo!' })
    showDeficiencyModal.value = false
  } catch {
    $q.notify({ type: 'negative', message: 'Errore salvataggio carenza' })
  } finally {
    savingDeficiency.value = false
  }
}

const openDeferredModal = async (row) => {
  selectedStudent.value = row
  deferredForm.student_id = row.student_id
  deferredForm.class_id = selectedClassId.value
  deferredForm.final_decision = 'promosso_con_debiti_saldati'
  deferredForm.notes = ''

  try {
    const res = await scrutinyService.getStudentDeficiencies(row.student_id)
    studentDeficienciesList.value = res.data || []
  } catch {
    studentDeficienciesList.value = []
  }

  showDeferredModal.value = true
}

const saveDeferredScrutiny = async () => {
  savingDeferred.value = true
  try {
    const payload = {
      student_id: deferredForm.student_id,
      class_id: deferredForm.class_id,
      final_decision: deferredForm.final_decision,
      notes: deferredForm.notes,
      deficiencies: studentDeficienciesList.value.map(d => ({
        deficiency_id: d.id,
        status: d.status,
        recovery_grade: d.recovery_grade
      }))
    }
    await scrutinyService.saveDeferredScrutiny(payload)
    $q.notify({ type: 'positive', message: t('help.teacher.scrutiny.deferredSaved') })
    showDeferredModal.value = false
    fetchMatrix()
  } catch {
    $q.notify({ type: 'negative', message: 'Errore salvataggio dello scrutinio differito' })
  } finally {
    savingDeferred.value = false
  }
}

const getGradeClass = (avg) => {
  if (!avg) return 'bg-slate-100'
  if (avg < 5.5) return 'bg-red-50 text-red-900'
  if (avg < 6) return 'bg-orange-50 text-orange-900'
  if (avg < 8) return 'bg-blue-50 text-blue-900'
  return 'bg-emerald-50 text-emerald-900'
}

const closeScrutiny = () => {
  if (!selectedClassId.value) return
  $q.dialog({
    title: 'Conferma Chiusura Scrutinio',
    message: 'Sei sicuro di voler chiudere e sigillare lo scrutinio per la classe selezionata? L\'operazione è definitiva.',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await scrutinyService.closeScrutiny(selectedClassId.value, period.value)
      $q.notify({ type: 'positive', message: 'Scrutinio chiuso ufficialmente e sigillato!' })
      fetchMatrix()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore durante la chiusura dello scrutinio' })
    }
  })
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
