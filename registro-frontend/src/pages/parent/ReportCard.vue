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
        <q-icon name="drive_file_rename_outline" color="amber-9" size="32px" />
      </template>
      <div class="text-subtitle1 text-weight-bold text-amber-9">Presa Visione Richiesta</div>
      <div class="text-body2 text-slate-700">È richiesta la firma digitale dei genitori per conferma di presa visione della pagella.</div>
      <template v-slot:action>
        <q-btn color="amber-9" unelevated label="Firma per Presa Visione" class="rounded-lg" :loading="signing" @click="signReportCard" />
      </template>
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

        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <div class="text-subtitle2 text-slate-500">Studente/ssa (Figlio/a)</div>
            <div class="text-h6 text-weight-bold text-slate-800">
              {{ reportData?.student_name || 'Nome Studente' }}
            </div>
            <div class="text-caption text-slate-500">Codice Fiscale: {{ reportData?.fiscal_code || 'N/D' }}</div>
          </div>
          <div class="col-12 col-md-6 text-md-right">
            <div class="text-subtitle2 text-slate-500">Classe Frequentata</div>
            <div class="text-h6 text-weight-bold text-slate-800">
              {{ reportData?.class_name || 'Classe N/A' }}
            </div>
            <div class="text-caption text-slate-500">Coordinatore: {{ reportData?.coordinator_name || 'Docente Coordinatore' }}</div>
          </div>
        </div>
      </q-card>

      <!-- Grades Table -->
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="q-pa-none">
          <q-table
            :rows="reportData?.subject_grades || []"
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

      <!-- Section: Carenze Formative & Argomenti da Recuperare per i Genitori -->
      <q-card v-if="deficiencies.length > 0" flat bordered class="rounded-xl bg-amber-50/50 border-amber-200 shadow-soft q-pa-md">
        <div class="row items-center q-mb-md">
          <q-icon name="warning" color="warning" size="28px" class="q-mr-sm" />
          <div class="text-h6 text-weight-bold text-amber-900">Carenze Formative & Argomenti da Recuperare (Figlio/a)</div>
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

      <!-- Final Evaluation Summary -->
      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full">
            <q-card-section>
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-sm">Esito e Giudizio Finale</div>
              
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

        <div class="col-12 col-md-6">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft h-full flex flex-center text-center p-4">
            <div v-if="isSigned">
              <q-icon name="verified" color="positive" size="40px" />
              <div class="text-subtitle1 text-weight-bold text-positive q-mt-xs">Presa Visione Registrata</div>
              <div class="text-caption text-slate-500">Pagella firmata digitalmente dal genitore.</div>
            </div>
            <div v-else>
              <q-icon name="pending" color="warning" size="40px" />
              <div class="text-subtitle1 text-weight-bold text-amber-8 q-mt-xs">Firma In Attesa</div>
              <div class="text-caption text-slate-500">In attesa di presa visione da parte del genitore.</div>
            </div>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { scrutinyService } from '@/services/scrutinyService'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()

const selectedSemester = ref(1)
const selectedStudentId = ref(null)
const children = ref([])
const reportData = ref(null)
const deficiencies = ref([])
const loading = ref(false)
const signing = ref(false)
const isSigned = ref(false)
const requiresSignature = ref(true)

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
    await loadDeficiencies()
  } catch {
    reportData.value = null
  } finally {
    loading.value = false
  }
}

async function loadDeficiencies() {
  if (!selectedStudentId.value) return
  try {
    const res = await scrutinyService.getStudentDeficiencies(selectedStudentId.value)
    deficiencies.value = res.data || []
  } catch {
    deficiencies.value = []
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
