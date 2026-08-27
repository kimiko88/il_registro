<template>
  <q-page class="q-pa-md recovery-page">
    <!-- Header Banner -->
    <div class="page-header q-mb-lg flex justify-between items-center wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="school" color="primary" />
          {{ $t('recovery.title') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ $t('recovery.subtitle') }}
        </p>
      </div>

      <div class="flex gap-sm">
        <q-btn
          color="secondary"
          icon="fact_check"
          :label="$t('recovery.recordTestBtn')"
          unelevated
          no-caps
          rounded
          @click="openRecordTestDialog"
        />
        <q-btn
          color="primary"
          icon="add_circle"
          :label="$t('recovery.newCourseBtn')"
          unelevated
          no-caps
          rounded
          @click="openCreateCourseDialog"
        />
      </div>
    </div>

    <!-- Filter & Stats Bar -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-4">
        <q-card class="stat-card shadow-1 rounded-borders">
          <q-card-section class="flex items-center justify-between">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">{{ $t('recovery.activeCourses') }}</div>
              <div class="text-h5 text-weight-bold text-primary">{{ courses.length }}</div>
            </div>
            <q-avatar color="primary" text-color="white" icon="auto_stories" font-size="24px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card class="stat-card shadow-1 rounded-borders">
          <q-card-section class="flex items-center justify-between">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">{{ $t('recovery.completedTests') }}</div>
              <div class="text-h5 text-weight-bold text-secondary">{{ tests.length }}</div>
            </div>
            <q-avatar color="secondary" text-color="white" icon="verified" font-size="24px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card class="stat-card shadow-1 rounded-borders">
          <q-card-section class="flex items-center justify-between">
            <div>
              <div class="text-caption text-grey-7 text-uppercase">{{ $t('recovery.resolvedDebts') }}</div>
              <div class="text-h5 text-weight-bold text-positive">{{ resolvedDebtsCount }}</div>
            </div>
            <q-avatar color="positive" text-color="white" icon="check_circle" font-size="24px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Tabs: Corsi di Recupero & Prove di Settembre -->
    <q-card class="shadow-2 rounded-borders">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="courses" icon="cast_for_education" :label="$t('recovery.tabs.courses')" no-caps />
        <q-tab name="tests" icon="assignment_turned_in" :label="$t('recovery.tabs.tests')" no-caps />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Tab Corsi -->
        <q-tab-panel name="courses">
          <q-table
            :rows="courses"
            :columns="courseColumns"
            row-key="id"
            :loading="loadingCourses"
            no-data-label="Nessun corso di recupero pianificato"
            flat
          >
            <template #body-cell-status="props">
              <q-td :props="props">
                <q-badge
                  :color="getStatusColor(props.value)"
                  :label="$t(`recovery.status.${props.value}`)"
                  rounded
                />
              </q-td>
            </template>

            <template #body-cell-actions="props">
              <q-td :props="props" align="right">
                <q-btn
                  flat
                  dense
                  round
                  icon="visibility"
                  color="primary"
                  @click="viewCourseDetails(props.row)"
                >
                  <q-tooltip>{{ $t('common.details') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab Prove Integrative -->
        <q-tab-panel name="tests">
          <q-table
            :rows="tests"
            :columns="testColumns"
            row-key="id"
            :loading="loadingTests"
            no-data-label="Nessuna prova integrativa registrata"
            flat
          >
            <template #body-cell-grade="props">
              <q-td :props="props">
                <span :class="props.value >= 6.0 ? 'text-positive text-weight-bold' : 'text-negative text-weight-bold'">
                  {{ props.value.toFixed(1) }}
                </span>
              </q-td>
            </template>

            <template #body-cell-outcome="props">
              <q-td :props="props">
                <q-badge
                  :color="props.value === 'recuperato' ? 'positive' : 'negative'"
                  :label="$t(`recovery.outcomes.${props.value}`)"
                  rounded
                />
              </q-td>
            </template>

            <template #body-cell-final_deliberation="props">
              <q-td :props="props">
                <q-chip
                  dense
                  :color="props.value === 'Ammesso' ? 'green-1' : 'red-1'"
                  :text-color="props.value === 'Ammesso' ? 'positive' : 'negative'"
                >
                  {{ props.value }}
                </q-chip>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog Creazione Corso -->
    <q-dialog v-model="createCourseDialog" persistent>
      <q-card style="min-width: 550px; max-width: 700px;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">{{ $t('recovery.dialog.newCourseTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitCreateCourse" class="q-gutter-md">
            <q-input
              v-model="courseForm.title"
              :label="$t('recovery.form.title')"
              outlined
              dense
              required
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="courseForm.subject_id"
                  :options="subjectOptions"
                  emit-value
                  map-options
                  :label="$t('recovery.form.subject')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="courseForm.period"
                  :options="periodOptions"
                  emit-value
                  map-options
                  :label="$t('recovery.form.period')"
                  outlined
                  dense
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="courseForm.total_hours"
                  type="number"
                  :label="$t('recovery.form.hours')"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model="courseForm.room"
                  :label="$t('recovery.form.room')"
                  outlined
                  dense
                />
              </div>
            </div>

            <q-input
              v-model="courseForm.description"
              type="textarea"
              rows="2"
              :label="$t('recovery.form.description')"
              outlined
              dense
            />

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="primary" :label="$t('common.save')" type="submit" :loading="savingCourse" unelevated no-caps />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Registrazione Prova Settembre -->
    <q-dialog v-model="recordTestDialog" persistent>
      <q-card style="min-width: 500px; max-width: 600px;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-secondary">{{ $t('recovery.dialog.recordTestTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitRecordTest" class="q-gutter-md">
            <q-select
              v-model="testForm.student_id"
              :options="studentOptions"
              emit-value
              map-options
              :label="$t('recovery.form.student')"
              outlined
              dense
              required
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="testForm.subject_id"
                  :options="subjectOptions"
                  emit-value
                  map-options
                  :label="$t('recovery.form.subject')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="testForm.class_id"
                  :options="classOptions"
                  emit-value
                  map-options
                  :label="$t('recovery.form.class')"
                  outlined
                  dense
                  required
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-input
                  v-model="testForm.test_date"
                  type="date"
                  :label="$t('recovery.form.testDate')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="testForm.test_type"
                  :options="testTypeOptions"
                  emit-value
                  map-options
                  :label="$t('recovery.form.testType')"
                  outlined
                  dense
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-input
                  v-model.number="testForm.grade"
                  type="number"
                  step="0.1"
                  min="1"
                  max="10"
                  :label="$t('recovery.form.grade')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model="testForm.verbale_number"
                  :label="$t('recovery.form.verbale')"
                  placeholder="Es. VERB-01/2025"
                  outlined
                  dense
                />
              </div>
            </div>

            <!-- Anteprima Calcolo Scioglimento Debito -->
            <q-banner dense rounded :class="testForm.grade >= 6.0 ? 'bg-green-1 text-positive' : 'bg-red-1 text-negative'">
              <template #avatar>
                <q-icon :name="testForm.grade >= 6.0 ? 'check_circle' : 'warning'" />
              </template>
              <div class="text-weight-bold">
                {{ testForm.grade >= 6.0 ? $t('recovery.preview.resolved') : $t('recovery.preview.unresolved') }}
              </div>
              <div class="text-caption">
                {{ testForm.grade >= 6.0 ? $t('recovery.preview.promoted') : $t('recovery.preview.notPromoted') }}
              </div>
            </q-banner>

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="secondary" :label="$t('common.save')" type="submit" :loading="savingTest" unelevated no-caps />
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
import recoveryService from 'src/services/recoveryService'
import api from 'src/services/api'

const $q = useQuasar()

const activeTab = ref('courses')
const courses = ref([])
const tests = ref([])
const loadingCourses = ref(false)
const loadingTests = ref(false)

const createCourseDialog = ref(false)
const savingCourse = ref(false)
const recordTestDialog = ref(false)
const savingTest = ref(false)

const subjectOptions = ref([])
const studentOptions = ref([])
const classOptions = ref([])

const periodOptions = [
  { label: 'Estivo (Luglio/Agosto)', value: 'summer' },
  { label: 'Intermedio (1° Quadrimestre)', value: 'intermedio' },
  { label: 'Pomeridiano Continuativo', value: 'pomeridiano' }
]

const testTypeOptions = [
  { label: 'Scritto', value: 'written' },
  { label: 'Orale', value: 'oral' },
  { label: 'Pratico / Laboratorio', value: 'practical' },
  { label: 'Misto (Scritto + Orale)', value: 'mixed' }
]

const courseForm = ref({
  title: '',
  subject_id: '',
  period: 'summer',
  total_hours: 10,
  room: '',
  description: '',
  academic_year: '2024/2025',
  teacher_id: '',
  sessions: [],
  student_ids: []
})

const testForm = ref({
  student_id: '',
  subject_id: '',
  class_id: '',
  test_date: new Date().toISOString().split('T')[0],
  test_type: 'written',
  grade: 6.0,
  verbale_number: '',
  notes: ''
})

const resolvedDebtsCount = computed(() => {
  return tests.value.filter(t => t.outcome === 'recuperato').length
})

const courseColumns = [
  { name: 'title', label: 'Titolo Corso', field: 'title', align: 'left', sortable: true },
  { name: 'subject_name', label: 'Materia', field: 'subject_name', align: 'left' },
  { name: 'period', label: 'Periodo', field: 'period', align: 'center' },
  { name: 'total_hours', label: 'Ore Totali', field: 'total_hours', align: 'center' },
  { name: 'room', label: 'Aula', field: 'room', align: 'center' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', field: 'actions', align: 'right' }
]

const testColumns = [
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'class_name', label: 'Classe', field: 'class_name', align: 'center' },
  { name: 'subject_name', label: 'Materia', field: 'subject_name', align: 'left' },
  { name: 'test_date', label: 'Data Prova', field: 'test_date', align: 'center', sortable: true },
  { name: 'grade', label: 'Voto Prova', field: 'grade', align: 'center' },
  { name: 'outcome', label: 'Scioglimento Debito', field: 'outcome', align: 'center' },
  { name: 'final_deliberation', label: 'Esito Scrutinio', field: 'final_deliberation', align: 'center' },
  { name: 'verbale_number', label: 'N° Verbale', field: 'verbale_number', align: 'center' }
]

function getStatusColor(status) {
  switch (status) {
    case 'in_progress': return 'warning'
    case 'completed': return 'positive'
    case 'cancelled': return 'negative'
    default: return 'primary'
  }
}

async function loadData() {
  loadingCourses.value = true
  loadingTests.value = true
  try {
    const [coursesRes, testsRes, classesRes, subjectsRes, studentsRes] = await Promise.all([
      recoveryService.listCourses(),
      recoveryService.listTests(),
      api.get('/classes'),
      api.get('/subjects'),
      api.get('/users/search?role=student')
    ])
    courses.value = coursesRes.data || []
    tests.value = testsRes.data || []
    classOptions.value = (classesRes.data || []).map(c => ({ label: c.name, value: c.id }))
    subjectOptions.value = (subjectsRes.data || []).map(s => ({ label: s.name, value: s.id }))
    studentOptions.value = (studentsRes.data || []).map(s => ({ label: `${s.last_name} ${s.first_name}`, value: s.id }))
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel caricamento dei corsi di recupero' })
  } finally {
    loadingCourses.value = false
    loadingTests.value = false
  }
}

function openCreateCourseDialog() {
  createCourseDialog.value = true
}

function openRecordTestDialog() {
  recordTestDialog.value = true
}

function viewCourseDetails(course) {
  $q.dialog({
    title: course.title,
    message: `Materia: ${course.subject_name || 'N/A'}\nOre: ${course.total_hours}h\nAula: ${course.room || 'N/A'}\nDescrizione: ${course.description || 'Nessuna descrizione'}`,
    ok: 'Chiudi'
  })
}

async function submitCreateCourse() {
  savingCourse.value = true
  try {
    await recoveryService.createCourse(courseForm.value)
    $q.notify({ type: 'positive', message: 'Corso di recupero creato con successo' })
    createCourseDialog.value = false
    loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante la creazione del corso' })
  } finally {
    savingCourse.value = false
  }
}

async function submitRecordTest() {
  savingTest.value = true
  try {
    await recoveryService.recordTest(testForm.value)
    $q.notify({ type: 'positive', message: 'Prova integrativa e scioglimento debito registrati' })
    recordTestDialog.value = false
    loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione della prova' })
  } finally {
    savingTest.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.recovery-page {
  max-width: 1300px;
  margin: 0 auto;
}
.stat-card {
  border-left: 4px solid var(--q-primary);
  transition: transform 0.2s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
}
</style>
