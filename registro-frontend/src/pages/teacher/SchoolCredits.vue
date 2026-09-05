<template>
  <q-page class="q-pa-md credits-page">
    <!-- Header Banner -->
    <div class="page-header q-mb-lg flex justify-between items-center wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="military_tech" color="primary" />
          {{ $t('credits.title') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ $t('credits.subtitle') }}
        </p>
      </div>

      <div class="flex gap-sm items-center">
        <q-select
          v-model="selectedClassId"
          :options="classOptions"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          style="min-width: 180px;"
          :label="$t('credits.selectClass')"
          @update:model-value="onClassChanged"
        />
        <q-btn
          color="primary"
          icon="calculate"
          :label="$t('credits.calculatorModalBtn')"
          unelevated
          no-caps
          rounded
          @click="openCalculatorDialog"
        />
      </div>
    </div>

    <!-- Ministeriale D.Lgs. 62/2017 Reference Guide Card -->
    <q-expansion-item
      class="shadow-1 rounded-borders bg-white q-mb-lg"
      icon="menu_book"
      :label="$t('credits.referenceGuideTitle')"
      caption="D.Lgs. 62/2017 - Tabella A di conversione crediti scolastici per il triennio (Max 40 punti)"
      header-class="text-weight-bold text-primary"
    >
      <q-card>
        <q-card-section>
          <div class="row q-col-gutter-md">
            <!-- 3° Anno -->
            <div class="col-12 col-md-4">
              <div class="text-weight-bold text-primary q-mb-xs">Classe 3ª (Max 12 pt)</div>
              <q-markup-table dense flat bordered>
                <thead><tr><th>Media (M)</th><th>Punti</th></tr></thead>
                <tbody>
                  <tr><td>M = 6</td><td>7 - 8</td></tr>
                  <tr><td>6 &lt; M ≤ 7</td><td>8 - 9</td></tr>
                  <tr><td>7 &lt; M ≤ 8</td><td>9 - 10</td></tr>
                  <tr><td>8 &lt; M ≤ 9</td><td>10 - 11</td></tr>
                  <tr><td>9 &lt; M ≤ 10</td><td>11 - 12</td></tr>
                </tbody>
              </q-markup-table>
            </div>

            <!-- 4° Anno -->
            <div class="col-12 col-md-4">
              <div class="text-weight-bold text-primary q-mb-xs">Classe 4ª (Max 13 pt)</div>
              <q-markup-table dense flat bordered>
                <thead><tr><th>Media (M)</th><th>Punti</th></tr></thead>
                <tbody>
                  <tr><td>M = 6</td><td>8 - 9</td></tr>
                  <tr><td>6 &lt; M ≤ 7</td><td>9 - 10</td></tr>
                  <tr><td>7 &lt; M ≤ 8</td><td>10 - 11</td></tr>
                  <tr><td>8 &lt; M ≤ 9</td><td>11 - 12</td></tr>
                  <tr><td>9 &lt; M ≤ 10</td><td>12 - 13</td></tr>
                </tbody>
              </q-markup-table>
            </div>

            <!-- 5° Anno -->
            <div class="col-12 col-md-4">
              <div class="text-weight-bold text-primary q-mb-xs">Classe 5ª (Max 15 pt)</div>
              <q-markup-table dense flat bordered>
                <thead><tr><th>Media (M)</th><th>Punti</th></tr></thead>
                <tbody>
                  <tr><td>M = 6</td><td>9 - 10</td></tr>
                  <tr><td>6 &lt; M ≤ 7</td><td>10 - 11</td></tr>
                  <tr><td>7 &lt; M ≤ 8</td><td>11 - 12</td></tr>
                  <tr><td>8 &lt; M ≤ 9</td><td>12 - 13</td></tr>
                  <tr><td>9 &lt; M ≤ 10</td><td>14 - 15</td></tr>
                </tbody>
              </q-markup-table>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-expansion-item>

    <!-- Tabellone Crediti Classe -->
    <q-card class="shadow-2 rounded-borders">
      <div v-if="loadingCredits" class="q-pa-md q-gutter-y-sm" role="status" aria-label="Caricamento crediti">
        <q-skeleton type="rect" height="46px" class="rounded-borders" />
        <q-skeleton type="rect" height="40px" class="rounded-borders" />
        <q-skeleton type="rect" height="40px" class="rounded-borders" />
      </div>
      <q-table
        v-else
        :rows="classCredits"
        :columns="creditColumns"
        row-key="id"
        :no-data-label="$t('credits.noCreditsSelected') || 'Seleziona una classe del triennio (3ª, 4ª o 5ª) per visualizzare i crediti assegnati'"
        flat
      >
        <template #body-cell-grade_average="props">
          <q-td :props="props">
            <span class="text-weight-bold text-primary">{{ props.value ? props.value.toFixed(2) : '-' }}</span>
          </q-td>
        </template>

        <template #body-cell-band="props">
          <q-td :props="props" align="center">
            <q-badge color="grey-7" :label="`${props.row.base_credit_range_min} - ${props.row.base_credit_range_max} pt`" />
          </q-td>
        </template>

        <template #body-cell-assigned_credit="props">
          <q-td :props="props" align="center">
            <q-chip color="primary" text-color="white" class="text-weight-bolder">
              {{ props.value }} pt
            </q-chip>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" align="right">
            <q-btn
              flat
              dense
              round
              icon="edit"
              color="primary"
              @click="editStudentCredit(props.row)"
            >
              <q-tooltip>{{ $t('credits.editCredit') }}</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Calcolatore & Assegnazione Credito -->
    <q-dialog v-model="calculatorDialog" persistent>
      <q-card style="width: min(600px, 95vw); max-width: 95vw;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">{{ $t('credits.calcModalTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitAssignCredit" class="q-gutter-md">
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="calcForm.student_id"
                  :options="studentOptions"
                  emit-value
                  map-options
                  :label="$t('credits.form.student')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model.number="calcForm.grade_level"
                  :options="gradeLevelOptions"
                  emit-value
                  map-options
                  :label="$t('credits.form.gradeLevel')"
                  outlined
                  dense
                  @update:model-value="recalcSuggestion"
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="calcForm.grade_average"
                  type="number"
                  step="0.01"
                  min="6.0"
                  max="10.0"
                  :label="$t('credits.form.average')"
                  outlined
                  dense
                  required
                  @update:model-value="recalcSuggestion"
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="calcForm.conduct_grade"
                  type="number"
                  min="6"
                  max="10"
                  :label="$t('credits.form.conduct')"
                  outlined
                  dense
                  @update:model-value="recalcSuggestion"
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm items-center">
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="calcForm.pcto_hours"
                  type="number"
                  :label="$t('credits.form.pctoHours')"
                  outlined
                  dense
                  @update:model-value="recalcSuggestion"
                />
              </div>
              <div class="col-12 col-md-6">
                <q-checkbox
                  v-model="calcForm.has_extracurricular"
                  :label="$t('credits.form.hasExtracurricular')"
                  @update:model-value="recalcSuggestion"
                />
              </div>
            </div>

            <!-- Suggerimento Algoritmo D.Lgs 62/2017 -->
            <q-card class="bg-blue-1 text-primary q-pa-sm rounded-borders">
              <div class="flex items-center justify-between">
                <span class="text-weight-bold">{{ $t('credits.form.suggestedBand') }}:</span>
                <q-badge color="primary" :label="`${calcResult.base_credit_range_min} - ${calcResult.base_credit_range_max} pt`" />
              </div>
              <div class="text-caption q-mt-xs">
                {{ calcResult.motivation || 'Attribuzione calcolata in conformità con la Tabella A del D.Lgs. 62/2017' }}
              </div>
            </q-card>

            <div class="row q-col-gutter-sm items-center">
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="calcForm.assigned_credit"
                  type="number"
                  :min="calcResult.base_credit_range_min || 0"
                  :max="calcResult.base_credit_range_max || 15"
                  :label="$t('credits.form.assignedCredit')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model="calcForm.deliberation_notes"
                  :label="$t('credits.form.deliberationNotes')"
                  placeholder="Note delibera CdC..."
                  outlined
                  dense
                />
              </div>
            </div>

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="primary" :label="$t('common.save')" type="submit" :loading="savingCredit" unelevated no-caps />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useSchoolYearStore } from '@/stores/schoolYear'
import creditService from 'src/services/creditService'
import api from 'src/services/api'

const $q = useQuasar()
const { t } = useI18n()
const schoolYearStore = useSchoolYearStore()

const selectedClassId = ref('')
const classOptions = ref([])
const studentOptions = ref([])
const classCredits = ref([])
const loadingCredits = ref(false)

const calculatorDialog = ref(false)
const savingCredit = ref(false)

const gradeLevelOptions = computed(() => [
  { label: t('credits.gradeLevels.grade3') || '3° Anno (Classe III)', value: 3 },
  { label: t('credits.gradeLevels.grade4') || '4° Anno (Classe IV)', value: 4 },
  { label: t('credits.gradeLevels.grade5') || '5° Anno (Classe V)', value: 5 }
])

const calcForm = ref({
  student_id: '',
  class_id: '',
  academic_year: schoolYearStore.selectedSchoolYear || '2024/2025',
  grade_level: 3,
  grade_average: 7.5,
  conduct_grade: 8,
  assigned_credit: 9,
  pcto_hours: 0,
  has_extracurricular: false,
  deliberation_notes: ''
})

const calcResult = ref({
  base_credit_range_min: 9,
  base_credit_range_max: 10,
  suggested_credit: 9,
  motivation: ''
})

const creditColumns = computed(() => [
  { name: 'student_name', label: t('credits.columns.student') || 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'grade_level', label: t('credits.columns.year') || 'Anno', field: 'grade_level', align: 'center', format: val => `${val}ª` },
  { name: 'grade_average', label: t('credits.columns.average') || 'Media Voti (M)', field: 'grade_average', align: 'center' },
  { name: 'conduct_grade', label: t('credits.columns.conduct') || 'Condotta', field: 'conduct_grade', align: 'center' },
  { name: 'band', label: t('credits.columns.band') || 'Fascia Ministeriale', field: 'base_credit_range_min', align: 'center' },
  { name: 'assigned_credit', label: t('credits.columns.assignedCredit') || 'Credito Assegnato', field: 'assigned_credit', align: 'center' },
  { name: 'deliberation_notes', label: t('credits.columns.deliberation') || 'Motivazione / Delibera', field: 'deliberation_notes', align: 'left' },
  { name: 'actions', label: t('credits.columns.actions') || 'Azioni', field: 'actions', align: 'right' }
])

function extractList(response) {
  if (!response) return []
  const data = response.data !== undefined ? response.data : response
  if (Array.isArray(data)) return data
  if (data && Array.isArray(data.users)) return data.users
  if (data && Array.isArray(data.classes)) return data.classes
  return []
}

async function loadInitialData() {
  try {
    const [classesRes, studentsRes] = await Promise.allSettled([
      api.get('/classes'),
      api.get('/users/search?role=student')
    ])

    const rawClasses = classesRes.status === 'fulfilled' ? extractList(classesRes.value) : []
    classOptions.value = rawClasses.map(c => {
      let name = c.name || `Classe ${c.id}`
      if (c.section && !name.endsWith(c.section)) {
        name += c.section
      }
      if (c.academic_year) {
        name += ` (${c.academic_year})`
      }
      return { label: name, value: c.id }
    })

    const rawStudents = studentsRes.status === 'fulfilled' ? extractList(studentsRes.value) : []
    studentOptions.value = rawStudents.map(s => ({
      label: `${s.last_name || ''} ${s.first_name || ''}`.trim() || s.email || s.id,
      value: s.id
    }))

    if (classOptions.value.length > 0) {
      selectedClassId.value = classOptions.value[0].value
      loadClassCredits()
    }
  } catch (err) {
    console.error('Error loading initial credits data', err)
  }
}

async function loadClassCredits() {
  if (!selectedClassId.value) return
  loadingCredits.value = true
  try {
    const res = await creditService.listClassCredits(selectedClassId.value)
    classCredits.value = extractList(res)
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel caricamento dei crediti scolastici' })
  } finally {
    loadingCredits.value = false
  }
}

async function onClassChanged() {
  loadClassCredits()
  if (selectedClassId.value) {
    try {
      const res = await api.get(`/users/search?role=student&class_id=${selectedClassId.value}`)
      const classStudents = extractList(res)
      if (classStudents.length > 0) {
        studentOptions.value = classStudents.map(s => ({
          label: `${s.last_name || ''} ${s.first_name || ''}`.trim() || s.email || s.id,
          value: s.id
        }))
      }
    } catch {
      // Keep general student options
    }
  }
}

async function recalcSuggestion() {
  try {
    const res = await creditService.calculateSuggested({
      grade_level: calcForm.value.grade_level,
      average: calcForm.value.grade_average,
      conduct: calcForm.value.conduct_grade,
      pcto_hours: calcForm.value.pcto_hours,
      has_extracurricular: calcForm.value.has_extracurricular
    })
    calcResult.value = res.data
    calcForm.value.assigned_credit = res.data.suggested_credit
  } catch (err) {
    console.error('Calculation error', err)
  }
}

function openCalculatorDialog() {
  calcForm.value.class_id = selectedClassId.value
  recalcSuggestion()
  calculatorDialog.value = true
}

function editStudentCredit(row) {
  calcForm.value = {
    student_id: row.student_id,
    class_id: row.class_id,
    academic_year: row.academic_year,
    grade_level: row.grade_level,
    grade_average: row.grade_average,
    conduct_grade: row.conduct_grade,
    assigned_credit: row.assigned_credit,
    pcto_hours: row.pcto_hours,
    has_extracurricular: row.has_extracurricular,
    deliberation_notes: row.deliberation_notes
  }
  recalcSuggestion()
  calculatorDialog.value = true
}

async function submitAssignCredit() {
  savingCredit.value = true
  try {
    await creditService.assignCredit(calcForm.value)
    $q.notify({ type: 'positive', message: 'Credito scolastico assegnato e registrato con successo' })
    calculatorDialog.value = false
    loadClassCredits()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'assegnazione del credito scolastico' })
  } finally {
    savingCredit.value = false
  }
}

onMounted(() => {
  loadInitialData()
})
</script>

<style scoped>
.credits-page {
  max-width: 1300px;
  margin: 0 auto;
}
</style>
