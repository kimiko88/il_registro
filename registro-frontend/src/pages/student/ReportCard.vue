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
              <div class="text-h6 text-weight-bold text-slate-800">ISTITUTO SCOLASTICO REGISTRO V2</div>
              <div class="text-caption text-slate-500">Anno Scolastico {{ reportData?.school_year || '2025/2026' }}</div>
            </div>
          </div>
          <div class="text-right">
            <q-badge color="indigo-1" text-color="indigo-8" class="text-subtitle2 q-px-md q-py-xs text-weight-bold">
              {{ selectedSemester }}° QUADRIMESTRE
            </q-badge>
          </div>
        </div>

        <div class="row q-col-gutter-md text-body2">
          <div class="col-12 col-sm-4">
            <span class="text-slate-400">Studente:</span>
            <div class="text-weight-bold text-slate-800 text-subtitle1">
              {{ reportData?.student_name || authStore.user?.name || 'Studente' }}
            </div>
          </div>
          <div class="col-12 col-sm-4">
            <span class="text-slate-400">Classe:</span>
            <div class="text-weight-bold text-slate-800 text-subtitle1">
              {{ reportData?.class_name || 'N/D' }}
            </div>
          </div>
          <div class="col-12 col-sm-4">
            <span class="text-slate-400">Media Generale:</span>
            <div class="text-weight-bold text-subtitle1" :class="getAverageClass(reportData?.overall_average)">
              {{ reportData?.overall_average ? reportData.overall_average.toFixed(2) : '0.00' }} / 10
            </div>
          </div>
        </div>
      </q-card>

      <!-- Grades Table Card -->
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">Valutazioni per Materia</div>
        </q-card-section>

        <div v-if="gradesStore.loading" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <q-table
          v-else
          flat
          dense
          :rows="reportData?.subjects || []"
          :columns="columns"
          row-key="subject_id"
          hide-pagination
          :pagination="{ rowsPerPage: 0 }"
          class="no-shadow"
        >
          <!-- Subject Column -->
          <template v-slot:body-cell-subject="props">
            <q-td :props="props" class="text-weight-bold text-slate-800">
              {{ props.row.subject }}
            </q-td>
          </template>

          <!-- Teacher Column -->
          <template v-slot:body-cell-teacher="props">
            <q-td :props="props" class="text-slate-600">
              {{ props.row.teacher || 'Docente' }}
            </q-td>
          </template>

          <!-- Period Grades Column -->
          <template v-slot:body-cell-grades="props">
            <q-td :props="props">
              <span class="text-caption text-slate-500">
                {{ props.row.grade_count || (props.row.grades ? props.row.grades.length : 0) }} voti registrati
              </span>
            </q-td>
          </template>

          <!-- Average Column -->
          <template v-slot:body-cell-average="props">
            <q-td :props="props" class="text-weight-bold" :class="getAverageClass(props.row.subject_average)">
              {{ props.row.subject_average ? props.row.subject_average.toFixed(2) : '-' }}
            </q-td>
          </template>

          <!-- Final Grade Column -->
          <template v-slot:body-cell-final_grade="props">
            <q-td :props="props">
              <q-badge
                size="md"
                :color="props.row.passed || props.row.subject_average >= 6 ? 'positive' : 'negative'"
                class="text-weight-bold text-subtitle2 q-px-sm"
              >
                {{ props.row.final_grade || Math.round(props.row.subject_average) || '-' }}
              </q-badge>
            </q-td>
          </template>

          <!-- Notes Column -->
          <template v-slot:body-cell-notes="props">
            <q-td :props="props" class="text-caption text-slate-500">
              {{ props.row.notes || '-' }}
            </q-td>
          </template>
        </q-table>
      </q-card>

      <!-- Summary & Attendance Section -->
      <div class="row q-col-gutter-md">
        <div class="col-12 col-sm-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">Riepilogo Presenze</div>
              <div class="text-caption text-slate-500 q-mb-md">Assenze e ritardi cumulate nel quadrimestre</div>

              <div class="row q-col-gutter-sm text-center">
                <div class="col-4">
                  <div class="bg-red-50 border border-red-200 rounded-lg p-2">
                    <div class="text-h6 text-weight-bold text-negative">
                      {{ reportData?.total_absence_days || attendanceStats.absences }}
                    </div>
                    <div class="text-caption text-slate-500">Assenze (gg)</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="bg-amber-50 border border-amber-200 rounded-lg p-2">
                    <div class="text-h6 text-weight-bold text-warning">
                      {{ attendanceStats.lates }}
                    </div>
                    <div class="text-caption text-slate-500">Ritardi</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="bg-blue-50 border border-blue-200 rounded-lg p-2">
                    <div class="text-h6 text-weight-bold text-primary">
                      {{ attendanceStats.earlyExits }}
                    </div>
                    <div class="text-caption text-slate-500">Uscite Anticipate</div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">Esito & Giudizio Comportamento</div>
              <div class="text-caption text-slate-500 q-mb-md">Stato finale del percorso di apprendimento</div>

              <div class="row items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-200 q-mb-sm">
                <span class="text-weight-bold text-slate-700">Voto di Comportamento:</span>
                <q-chip color="primary" text-color="white" class="text-weight-bold">
                  {{ reportData?.behavior_grade || 8 }} / 10
                </q-chip>
              </div>

              <div class="row items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-200">
                <span class="text-weight-bold text-slate-700">Esito Periodo:</span>
                <q-badge
                  :color="reportData?.promoted === 'SÌ' || reportData?.overall_average >= 6 ? 'positive' : 'warning'"
                  class="text-weight-bold q-px-md q-py-xs"
                >
                  {{ reportData?.promoted === 'SÌ' || reportData?.overall_average >= 6 ? 'PROMOSSO / REGOLARE' : 'CON GIUDIZIO SOSPESO' }}
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
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useGradesStore } from '@/stores/grades'
import { useAuthStore } from '@/stores/auth'
import api from 'src/services/api'

const $q = useQuasar()
const gradesStore = useGradesStore()
const authStore = useAuthStore()

const selectedSemester = ref(1)
const reportData = ref(null)

const attendanceStats = reactive({
  absences: 0,
  lates: 0,
  earlyExits: 0
})

const columns = [
  { name: 'subject', label: 'Materia', align: 'left', field: 'subject', sortable: true },
  { name: 'teacher', label: 'Docente', align: 'left', field: 'teacher' },
  { name: 'grades', label: 'Voti Periodo', align: 'center', field: 'grade_count' },
  { name: 'average', label: 'Media Voti', align: 'center', field: 'subject_average', sortable: true },
  { name: 'final_grade', label: 'Voto Finale', align: 'center', field: 'final_grade', sortable: true },
  { name: 'notes', label: 'Note', align: 'left', field: 'notes' }
]

onMounted(async () => {
  await loadReport()
  await loadAttendanceStats()
})

async function loadReport() {
  try {
    reportData.value = await gradesStore.fetchSemesterReport(selectedSemester.value)
  } catch (e) {
    reportData.value = null
  }
}

async function loadAttendanceStats() {
  try {
    const res = await api.get('/attendance/my-summary')
    if (res.data) {
      attendanceStats.absences = res.data.absences || 0
      attendanceStats.lates = res.data.lates || 0
      attendanceStats.earlyExits = res.data.early_exits || 0
    }
  } catch {
    // fallback defaults
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
    // Fallback to print
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
