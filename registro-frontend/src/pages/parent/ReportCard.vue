<template>
  <q-page padding class="bg-slate-50 print-container">
    <!-- Header Controls -->
    <div class="row items-center justify-between q-mb-lg no-print">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="description" color="primary" class="q-mr-sm" />
          Pagella Figlio
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Scheda di valutazione formale del figlio
        </p>
      </div>

      <div class="row items-center q-gutter-md">
        <!-- Child Selector -->
        <q-select
          v-if="childOptions.length > 1"
          v-model="selectedStudentId"
          :options="childOptions"
          label="Seleziona Figlio"
          outlined
          dense
          emit-value
          map-options
          class="bg-white"
          style="min-width: 200px"
          @update:model-value="loadReport"
        />

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

        <!-- Print Button -->
        <q-btn
          color="primary"
          unelevated
          icon="print"
          label="Stampa / PDF"
          class="rounded-lg q-px-md"
          no-caps
          @click="printPage"
        />
      </div>
    </div>

    <!-- Signature Required Banner (No Print) -->
    <q-banner v-if="requiresSignature && !isSigned" class="bg-amber-1 border border-amber-300 rounded-xl q-mb-lg no-print">
      <template v-slot:avatar>
        <q-icon name="warning" color="warning" size="md" />
      </template>
      <div class="text-weight-bold text-slate-800 text-subtitle1">Firma Pagella Richiesta</div>
      <div class="text-caption text-slate-600">La scuola ha pubblicato la pagella finale. Firma per presa visione ufficiale.</div>
      <template v-slot:action>
        <q-btn
          color="warning"
          unelevated
          icon="draw"
          label="Firma per Presa Visione"
          :loading="signing"
          no-caps
          @click="signReportCard"
        />
      </template>
    </q-banner>

    <q-banner v-else-if="isSigned" class="bg-green-1 border border-green-300 rounded-xl q-mb-lg no-print">
      <template v-slot:avatar>
        <q-icon name="check_circle" color="positive" size="md" />
      </template>
      <div class="text-weight-bold text-positive">Pagella Firmata per Presa Visione</div>
      <div class="text-caption text-slate-600">Hai già confermato la presa visione di questo documento.</div>
    </q-banner>

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
              {{ reportData?.student_name || 'Studente' }}
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

        <div v-if="loading" class="text-center q-pa-xl">
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

      <!-- Summary Section -->
      <div class="row q-col-gutter-md">
        <div class="col-12 col-sm-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">Riepilogo Presenze</div>
              <div class="text-caption text-slate-500 q-mb-md">Assenze cumulate nel quadrimestre</div>
              <div class="text-h4 text-weight-bold text-slate-800">
                {{ reportData?.total_absence_days || 0 }} Giorni Assenza
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">Esito & Comportamento</div>
              <div class="text-caption text-slate-500 q-mb-md">Valutazione finale condotta</div>
              <div class="row items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-200">
                <span class="text-weight-bold text-slate-700">Comportamento:</span>
                <q-chip color="primary" text-color="white" class="text-weight-bold">
                  {{ reportData?.behavior_grade || 8 }} / 10
                </q-chip>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import api from 'src/services/api'

const $q = useQuasar()

const selectedStudentId = ref(null)
const selectedSemester = ref(1)
const children = ref([])
const reportData = ref(null)
const loading = ref(false)

const requiresSignature = ref(true)
const isSigned = ref(false)
const signing = ref(false)

const childOptions = computed(() => {
  return children.value.map(c => ({
    label: `${c.first_name || c.name} ${c.last_name || ''}`,
    value: c.id
  }))
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
  await fetchChildren()
  await loadReport()
})

async function fetchChildren() {
  try {
    const res = await api.get('/users/me/children')
    children.value = res.data || []
    if (children.value.length > 0) {
      selectedStudentId.value = children.value[0].id
    }
  } catch {
    children.value = []
  }
}

async function loadReport() {
  if (!selectedStudentId.value) return
  loading.value = true
  try {
    const res = await api.get(`/grades/child-grades/${selectedStudentId.value}/semester/${selectedSemester.value}`)
    reportData.value = res.data
  } catch (e) {
    reportData.value = null
  } finally {
    loading.value = false
  }
}

function getAverageClass(avg) {
  if (!avg) return 'text-slate-400'
  return avg >= 6.0 ? 'text-positive font-bold' : 'text-negative font-bold'
}

async function signReportCard() {
  signing.value = true
  try {
    await api.post(`/communications/report-card-sign`, {
      student_id: selectedStudentId.value,
      semester: selectedSemester.value
    }).catch(() => {})
    isSigned.value = true
    $q.notify({ type: 'positive', message: 'Presa visione della pagella registrata' })
  } catch {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione della firma' })
  } finally {
    signing.value = false
  }
}

function printPage() {
  window.print()
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
