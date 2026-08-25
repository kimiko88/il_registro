<template>
  <q-page padding class="bg-slate-50 print-container">
    <!-- Header Controls (Hidden during print) -->
    <div class="row items-center justify-between q-mb-lg no-print">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="description" color="primary" class="q-mr-sm" />
          Pagella Valutazione Finale
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Scheda di valutazione formale di fine periodo scolastico
        </p>
      </div>

      <div class="row items-center q-gutter-md">
        <!-- Semester Selector -->
        <q-btn-toggle
          v-model="selectedSemester"
          toggle-color="primary"
          unelevated
          no-caps
          class="bg-white border border-slate-200 rounded-xl"
          :options="[
            { label: '1° Quadrimestre', value: 1 },
            { label: '2° Quadrimestre', value: 2 }
          ]"
          @update:model-value="loadReport"
        />

        <!-- Export PDF Button -->
        <q-btn
          color="primary"
          unelevated
          icon="picture_as_pdf"
          label="Esporta PDF"
          class="rounded-lg q-px-md"
          no-caps
          @click="exportPDF"
        />
      </div>
    </div>

    <!-- Printable Report Card Paper Document -->
    <div class="printable-card space-y-6">
      <!-- School & Student Info Card -->
      <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md">
        <div class="row items-center justify-between border-b border-slate-200 q-pb-md q-mb-md">
          <div class="row items-center">
            <q-avatar color="primary" text-color="white" icon="school" size="48px" class="q-mr-md" />
            <div>
              <div class="text-h6 text-weight-bold text-slate-800">{{ authStore.user?.school_name || reportData?.school_name || 'ISTITUTO SCOLASTICO' }}</div>
              <div class="text-caption text-slate-500">{{ t('common.year') || 'Anno Scolastico' }} {{ reportData?.school_year || '2025/2026' }}</div>
            </div>
          </div>
          <div class="text-right">
            <q-badge color="indigo-1" text-color="indigo-8" class="text-subtitle2 q-px-md q-py-xs text-weight-bold">
              {{ selectedSemester }}° QUADRIMESTRE
            </q-badge>
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <div class="text-subtitle2 text-slate-500">{{ t('usersPage.roleStudents') || 'Studente/ssa' }}</div>
            <div class="text-h6 text-weight-bold text-slate-800">
              {{ reportData?.student_name || `${authStore.user?.first_name || ''} ${authStore.user?.last_name || ''}`.trim() || 'Studente' }}
            </div>
            <div class="text-caption text-slate-500">Codice Fiscale: {{ authStore.user?.fiscal_code || reportData?.fiscal_code || 'N/D' }}</div>
          </div>
          <div class="col-12 col-md-6 text-md-right">
            <div class="text-subtitle2 text-slate-500">{{ t('classRegister.selectClass') || 'Classe Frequentata' }}</div>
            <div class="text-h6 text-weight-bold text-slate-800">
              {{ reportData?.class_name || 'Classe N/A' }}
            </div>
            <div class="text-caption text-slate-500">{{ t('classRegister.coordinator') || 'Coordinatore' }}: {{ reportData?.coordinator_name || 'Docente Coordinatore' }}</div>
          </div>
        </div>
      </q-card>

      <!-- Grades Table -->
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="q-pa-none">
          <q-table
            :rows="reportData?.subjects || reportData?.subject_grades || []"
            :columns="columns"
            row-key="subject"
            flat
            hide-bottom
            :pagination="{ rowsPerPage: 0 }"
            class="bg-transparent"
          >
            <template v-slot:body-cell-average="props">
              <q-td :props="props" align="center" :class="getAverageClass(props.row.subject_average)">
                {{ props.row.subject_average ? props.row.subject_average.toFixed(2) : '-' }}
              </q-td>
            </template>

            <template v-slot:body-cell-final_grade="props">
              <q-td :props="props" align="center">
                <q-chip
                  :color="props.row.final_grade >= 6 ? 'positive' : 'negative'"
                  text-color="white"
                  class="text-weight-bold"
                >
                  {{ props.row.final_grade ?? '-' }}
                </q-chip>
              </q-td>
            </template>
          </q-table>
        </q-card-section>
      </q-card>

      <!-- Section: Carenze Formative & Argomenti da Recuperare -->
      <q-card v-if="deficiencies.length > 0" flat bordered class="rounded-xl bg-amber-50/50 border-amber-200 shadow-soft q-pa-md">
        <div class="row items-center q-mb-md">
          <q-icon name="warning" color="warning" size="28px" class="q-mr-sm" />
          <div class="text-h6 text-weight-bold text-amber-900">Carenze Formative & Argomenti da Recuperare</div>
        </div>

        <div class="space-y-3">
          <div v-for="def in deficiencies" :key="def.id" class="bg-white p-3 rounded-lg border border-amber-200 shadow-sm">
            <div class="row items-center justify-between">
              <div class="text-subtitle1 text-weight-bold text-slate-800">{{ def.subject_name || 'Materia' }}</div>
              <q-badge :color="def.status === 'recuperato' ? 'positive' : 'warning'" class="q-px-sm q-py-xs text-weight-bold">
                {{ def.status === 'recuperato' ? 'RECUPERATO' : 'DA RECUPERARE' }}
              </q-badge>
            </div>
            <div class="text-body2 text-slate-700 q-mt-xs">
              <strong>Argomenti della carenza:</strong> {{ def.topics }}
            </div>
            <div class="row items-center justify-between text-caption text-slate-500 q-mt-xs">
              <span>Modalità: {{ def.recovery_mode }}</span>
              <span v-if="def.recovery_grade">Voto prova recupero: <strong>{{ def.recovery_grade }}</strong></span>
            </div>
          </div>
        </div>
      </q-card>

      <!-- Final Evaluation Summary & Attendance -->
      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-sm">{{ t('attendancePage.tableTitle') || 'Quadro Presenze' }}</div>
              <div class="row q-col-gutter-sm text-center">
                <div class="col-4">
                  <div class="bg-slate-50 p-3 rounded-lg border border-slate-200">
                    <div class="text-h6 text-weight-bold text-primary">{{ attendanceStats.absences }}</div>
                    <div class="text-caption text-slate-500">{{ t('classRegister.absent') || 'Ore Assenza' }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="bg-slate-50 p-3 rounded-lg border border-slate-200">
                    <div class="text-h6 text-weight-bold text-warning">{{ attendanceStats.lates }}</div>
                    <div class="text-caption text-slate-500">{{ t('classRegister.late') || 'Ritardi' }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="bg-slate-50 p-3 rounded-lg border border-slate-200">
                    <div class="text-h6 text-weight-bold text-info">{{ attendanceStats.earlyExits }}</div>
                    <div class="text-caption text-slate-500">{{ t('classRegister.earlyExit') || 'Uscite Ant.' }}</div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-sm">{{ t('reportCardPage.finalEvaluation') || 'Esito e Giudizio Finale' }}</div>
              
              <div class="row items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-200 q-mb-sm">
                <span class="text-weight-bold text-slate-700">Voto di Comportamento:</span>
                <q-chip color="primary" text-color="white" class="text-weight-bold">
                  {{ reportData?.behavior_grade ?? '—' }} / 10
                </q-chip>
              </div>

              <div class="row items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-200">
                <span class="text-weight-bold text-slate-700">Esito Periodo:</span>
                <q-badge
                  :color="(reportData?.promoted === 'SÌ' || reportData?.promoted === true) ? 'positive' : 'warning'"
                  class="text-weight-bold q-px-md q-py-xs"
                >
                  {{ (reportData?.promoted === 'SÌ' || reportData?.promoted === true) ? 'PROMOSSO / REGOLARE' : 'CON GIUDIZIO SOSPESO' }}
                </q-badge>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useGradesStore } from '@/stores/grades'
import { useAuthStore } from '@/stores/auth'
import { scrutinyService } from '@/services/scrutinyService'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const gradesStore = useGradesStore()
const authStore = useAuthStore()

const selectedSemester = ref(1)
const reportData = ref(null)
const deficiencies = ref([])

const attendanceStats = reactive({
  absences: 0,
  lates: 0,
  earlyExits: 0
})

const columns = computed(() => [
  { name: 'subject', label: t('agendaPage.subject') || 'Materia', align: 'left', field: 'subject', sortable: true },
  { name: 'teacher', label: t('agendaPage.teacher') || 'Docente', align: 'left', field: 'teacher' },
  { name: 'grades', label: t('gradesPage.periodGrades') || 'Voti Periodo', align: 'center', field: 'grade_count' },
  { name: 'average', label: t('roleDashboards.averageGrade') || 'Media Voti', align: 'center', field: 'subject_average', sortable: true },
  { name: 'final_grade', label: t('gradesPage.finalGrade') || 'Voto Finale', align: 'center', field: 'final_grade', sortable: true },
  { name: 'notes', label: t('common.notes') || 'Note', align: 'left', field: 'notes' }
])

onMounted(async () => {
  await loadReport()
  await loadAttendanceStats()
  await loadDeficiencies()
})

async function loadReport() {
  try {
    reportData.value = await gradesStore.fetchSemesterReport(selectedSemester.value)
  } catch (e) {
    reportData.value = null
    $q.notify({ type: 'negative', message: 'Errore nel caricamento della pagella' })
  }
}

async function loadDeficiencies() {
  try {
    if (authStore.user?.id) {
      const res = await scrutinyService.getStudentDeficiencies(authStore.user.id)
      deficiencies.value = res.data || []
    }
  } catch {
    deficiencies.value = []
  }
}

async function loadAttendanceStats() {
  try {
    const res = await api.get('/attendance/my-attendance/summary')
    if (res.data) {
      attendanceStats.absences = res.data.total_absences ?? res.data.absences ?? 0
      attendanceStats.lates = res.data.total_lates ?? res.data.lates ?? 0
      attendanceStats.earlyExits = res.data.total_early_exits ?? res.data.early_exits ?? 0
    } else if (reportData.value?.total_absence_days != null) {
      attendanceStats.absences = reportData.value.total_absence_days
    }
  } catch {
    if (reportData.value?.total_absence_days != null) {
      attendanceStats.absences = reportData.value.total_absence_days
    }
  }
}

function getAverageClass(avg) {
  if (!avg) return 'text-slate-400'
  return avg >= 6.0 ? 'text-positive font-bold' : 'text-negative font-bold'
}

async function exportPDF() {
  try {
    await gradesStore.downloadReportCardPDF(selectedSemester.value)
    $q.notify({ type: 'positive', message: 'PDF della pagella scaricato con successo' })
  } catch (e) {
    window.print()
  }
}
</script>

<style>
@media print {
  body * {
    visibility: hidden;
  }
  .no-print {
    display: none !important;
  }
  .print-container, .print-container * {
    visibility: visible;
  }
  .print-container {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
  }
}
</style>
