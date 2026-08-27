<template>
  <q-page class="q-pa-md support-page">
    <!-- Header Banner -->
    <div class="page-header q-mb-lg flex justify-between items-center wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="favorite" color="primary" />
          {{ $t('support.title') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ $t('support.subtitle') }}
        </p>
      </div>

      <div class="flex gap-sm">
        <q-btn
          color="secondary"
          icon="track_changes"
          :label="$t('support.newPeiGoalBtn')"
          unelevated
          no-caps
          rounded
          @click="openNewGoalDialog"
        />
        <q-btn
          color="primary"
          icon="edit_note"
          :label="$t('support.newDiaryEntryBtn')"
          unelevated
          no-caps
          rounded
          @click="openNewDiaryDialog"
        />
      </div>
    </div>

    <!-- Student Filter Selection -->
    <q-card class="shadow-1 rounded-borders q-mb-md">
      <q-card-section class="flex items-center gap-md wrap">
        <div class="text-weight-bold text-grey-8">{{ $t('support.filterStudent') }}:</div>
        <q-select
          v-model="selectedStudentId"
          :options="studentFilterOptions"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          style="min-width: 250px;"
          @update:model-value="onStudentChanged"
        />
      </q-card-section>
    </q-card>

    <!-- Tabs: Diario di Bordo Orario & Piano Obiettivi PEI -->
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
        <q-tab name="diary" icon="menu_book" :label="$t('support.tabs.diary')" no-caps />
        <q-tab name="pei" icon="flag" :label="$t('support.tabs.peiGoals')" no-caps />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Tab Diario di Bordo -->
        <q-tab-panel name="diary">
          <q-table
            :rows="diaryEntries"
            :columns="diaryColumns"
            row-key="id"
            :loading="loadingDiary"
            no-data-label="Nessuna voce nel diario di sostegno registrata"
            flat
          >
            <template #body-cell-activity_type="props">
              <q-td :props="props">
                <q-chip dense color="blue-1" text-color="primary">
                  {{ $t(`support.activityTypes.${props.value}`) }}
                </q-chip>
              </q-td>
            </template>

            <template #body-cell-is_shared_with_family="props">
              <q-td :props="props" align="center">
                <q-badge
                  :color="props.value ? 'positive' : 'grey-5'"
                  :label="props.value ? $t('support.shared') : $t('support.internal')"
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
                  @click="viewDiaryDetails(props.row)"
                >
                  <q-tooltip>{{ $t('common.details') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab Obiettivi PEI -->
        <q-tab-panel name="pei">
          <q-table
            :rows="peiGoals"
            :columns="peiColumns"
            row-key="id"
            :loading="loadingPei"
            no-data-label="Nessun obiettivo PEI inserito"
            flat
          >
            <template #body-cell-pei_type="props">
              <q-td :props="props">
                <q-badge
                  :color="props.value === 'equipollente' ? 'primary' : 'purple'"
                  :label="props.value === 'equipollente' ? 'PEI Equipollente' : 'PEI Differenziato'"
                  rounded
                />
              </q-td>
            </template>

            <template #body-cell-axis="props">
              <q-td :props="props">
                <span class="text-weight-bold text-capitalize">{{ props.value }}</span>
              </q-td>
            </template>

            <template #body-cell-progress_status="props">
              <q-td :props="props">
                <q-select
                  v-model="props.row.progress_status"
                  :options="progressStatusOptions"
                  emit-value
                  map-options
                  dense
                  outlined
                  @update:model-value="val => updateGoalStatus(props.row.id, val)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog Nuova Voce Diario -->
    <q-dialog v-model="newDiaryDialog" persistent>
      <q-card style="min-width: 550px; max-width: 700px;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">{{ $t('support.dialog.newDiaryTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitDiaryEntry" class="q-gutter-md">
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="diaryForm.student_id"
                  :options="studentOptions"
                  emit-value
                  map-options
                  :label="$t('support.form.student')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="diaryForm.class_id"
                  :options="classOptions"
                  emit-value
                  map-options
                  :label="$t('support.form.class')"
                  outlined
                  dense
                  required
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-input
                  v-model="diaryForm.entry_date"
                  type="date"
                  :label="$t('support.form.date')"
                  outlined
                  dense
                  required
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="diaryForm.time_slot"
                  :options="[
                    '1ª Ora (08:00 - 09:00)',
                    '2ª Ora (09:00 - 10:00)',
                    '3ª Ora (10:00 - 11:00)',
                    '4ª Ora (11:00 - 12:00)',
                    '5ª Ora (12:00 - 13:00)',
                    '6ª Ora (13:00 - 14:00)'
                  ]"
                  :label="$t('support.form.timeSlot')"
                  outlined
                  dense
                  required
                />
              </div>
            </div>

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="diaryForm.activity_type"
                  :options="activityTypeOptions"
                  emit-value
                  map-options
                  :label="$t('support.form.activityType')"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="diaryForm.co_teacher_id"
                  :options="teacherOptions"
                  emit-value
                  map-options
                  clearable
                  :label="$t('support.form.coTeacher')"
                  outlined
                  dense
                />
              </div>
            </div>

            <q-input
              v-model="diaryForm.topic_and_activities"
              type="textarea"
              rows="3"
              :label="$t('support.form.topicAndActivities')"
              outlined
              dense
              required
            />

            <q-input
              v-model="diaryForm.student_responses"
              type="textarea"
              rows="2"
              :label="$t('support.form.studentResponses')"
              placeholder="Livello di autonomia, concentrazione e partecipazione..."
              outlined
              dense
            />

            <q-input
              v-model="diaryForm.educator_notes"
              type="textarea"
              rows="2"
              :label="$t('support.form.educatorNotes')"
              placeholder="Indicazioni per Educatore OEPA / Assistente all'Autonomia..."
              outlined
              dense
            />

            <q-checkbox
              v-model="diaryForm.is_shared_with_family"
              :label="$t('support.form.shareWithFamily')"
            />

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="primary" :label="$t('common.save')" type="submit" :loading="savingDiary" unelevated no-caps />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Nuovo Obiettivo PEI -->
    <q-dialog v-model="newGoalDialog" persistent>
      <q-card style="min-width: 500px; max-width: 600px;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-secondary">{{ $t('support.dialog.newPeiGoalTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitPeiGoal" class="q-gutter-md">
            <q-select
              v-model="goalForm.student_id"
              :options="studentOptions"
              emit-value
              map-options
              :label="$t('support.form.student')"
              outlined
              dense
              required
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="goalForm.pei_type"
                  :options="[
                    { label: 'Equipollente (Prog. Semplificata)', value: 'equipollente' },
                    { label: 'Differenziato (O.M. 90/2001)', value: 'differenziato' }
                  ]"
                  emit-value
                  map-options
                  :label="$t('support.form.peiType')"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="goalForm.axis"
                  :options="axisOptions"
                  emit-value
                  map-options
                  :label="$t('support.form.axis')"
                  outlined
                  dense
                />
              </div>
            </div>

            <q-input
              v-model="goalForm.title"
              :label="$t('support.form.goalTitle')"
              outlined
              dense
              required
            />

            <q-input
              v-model="goalForm.description"
              type="textarea"
              rows="2"
              :label="$t('support.form.goalDescription')"
              outlined
              dense
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="goalForm.expected_term"
                  :options="[
                    { label: '1° Quadrimestre', value: 'q1' },
                    { label: '2° Quadrimestre', value: 'q2' },
                    { label: 'Annuale', value: 'annuale' }
                  ]"
                  emit-value
                  map-options
                  :label="$t('support.form.expectedTerm')"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="goalForm.progress_status"
                  :options="progressStatusOptions"
                  emit-value
                  map-options
                  :label="$t('support.form.status')"
                  outlined
                  dense
                />
              </div>
            </div>

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="secondary" :label="$t('common.save')" type="submit" :loading="savingGoal" unelevated no-caps />
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
import supportService from 'src/services/supportService'
import api from 'src/services/api'

const $q = useQuasar()

const activeTab = ref('diary')
const selectedStudentId = ref('')
const diaryEntries = ref([])
const peiGoals = ref([])
const loadingDiary = ref(false)
const loadingPei = ref(false)

const newDiaryDialog = ref(false)
const savingDiary = ref(false)
const newGoalDialog = ref(false)
const savingGoal = ref(false)

const studentOptions = ref([])
const teacherOptions = ref([])
const classOptions = ref([])

const studentFilterOptions = computed(() => [
  { label: 'Tutti gli studenti', value: '' },
  ...studentOptions.value
])

const activityTypeOptions = [
  { label: 'In Classe con Docente Curricolare', value: 'in_classe' },
  { label: 'Laboratorio Didattico', value: 'laboratorio' },
  { label: 'Aula Sostegno / Spazio Dedicato', value: 'aula_sostegno' },
  { label: 'Attività Individualizzata 1:1', value: 'individuale' },
  { label: 'Piccolo Gruppo Cooperativo', value: 'piccolo_gruppo' }
]

const axisOptions = [
  { label: 'Autonomia e Cura della Persona', value: 'autonomia' },
  { label: 'Cognitiva, Neuropsicologica e Apprendimento', value: 'cognitiva' },
  { label: 'Comunicazione e Linguaggio', value: 'comunicazionale' },
  { label: 'Relazionale e Socializzazione', value: 'relazionale' },
  { label: 'Sensoriale e Motoria', value: 'sensoriale' }
]

const progressStatusOptions = [
  { label: 'Non Avviato', value: 'non_avviato' },
  { label: 'Iniziale', value: 'iniziale' },
  { label: 'Intermedio', value: 'intermedio' },
  { label: 'Avanzato', value: 'avanzato' },
  { label: 'Raggiunto', value: 'raggiunto' }
]

const diaryForm = ref({
  student_id: '',
  class_id: '',
  entry_date: new Date().toISOString().split('T')[0],
  time_slot: '1ª Ora (08:00 - 09:00)',
  activity_type: 'in_classe',
  co_teacher_id: null,
  topic_and_activities: '',
  student_responses: '',
  educator_notes: '',
  is_shared_with_family: false
})

const goalForm = ref({
  student_id: '',
  pei_type: 'equipollente',
  axis: 'autonomia',
  title: '',
  description: '',
  expected_term: 'annuale',
  progress_status: 'non_avviato'
})

const diaryColumns = [
  { name: 'entry_date', label: 'Data', field: 'entry_date', align: 'center', sortable: true },
  { name: 'time_slot', label: 'Ora', field: 'time_slot', align: 'center' },
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'activity_type', label: 'Tipo Attività', field: 'activity_type', align: 'center' },
  { name: 'co_teacher_name', label: 'Docente Compresenza', field: 'co_teacher_name', align: 'left' },
  { name: 'topic_and_activities', label: 'Argomenti & Attività', field: 'topic_and_activities', align: 'left' },
  { name: 'is_shared_with_family', label: 'Famiglia', field: 'is_shared_with_family', align: 'center' },
  { name: 'actions', label: 'Azioni', field: 'actions', align: 'right' }
]

const peiColumns = [
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'pei_type', label: 'Tipo PEI', field: 'pei_type', align: 'center' },
  { name: 'axis', label: 'Asse di Sviluppo', field: 'axis', align: 'left' },
  { name: 'title', label: 'Obiettivo', field: 'title', align: 'left' },
  { name: 'expected_term', label: 'Termine', field: 'expected_term', align: 'center' },
  { name: 'progress_status', label: 'Avanzamento', field: 'progress_status', align: 'center' }
]

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
    const [studentsRes, teachersRes, classesRes] = await Promise.allSettled([
      api.get('/users/search?role=student'),
      api.get('/users/search?role=teacher'),
      api.get('/classes')
    ])

    const rawStudents = studentsRes.status === 'fulfilled' ? extractList(studentsRes.value) : []
    studentOptions.value = rawStudents.map(s => ({
      label: `${s.last_name || ''} ${s.first_name || ''}`.trim() || s.email || s.id,
      value: s.id
    }))

    const rawTeachers = teachersRes.status === 'fulfilled' ? extractList(teachersRes.value) : []
    teacherOptions.value = rawTeachers.map(t => ({
      label: `${t.last_name || ''} ${t.first_name || ''}`.trim() || t.email || t.id,
      value: t.id
    }))

    const rawClasses = classesRes.status === 'fulfilled' ? extractList(classesRes.value) : []
    classOptions.value = rawClasses.map(c => {
      let name = c.name || `Classe ${c.id}`
      if (c.section && !name.endsWith(c.section)) name += c.section
      return { label: name, value: c.id }
    })

    loadData()
  } catch (err) {
    console.error('Error loading support register initial data', err)
  }
}

async function loadData() {
  loadingDiary.value = true
  loadingPei.value = true
  try {
    const params = selectedStudentId.value ? { student_id: selectedStudentId.value } : {}
    const [diaryRes, peiRes] = await Promise.allSettled([
      supportService.listDiaryEntries(params),
      supportService.listPeiGoals(params)
    ])
    if (diaryRes.status === 'fulfilled') {
      diaryEntries.value = extractList(diaryRes.value)
    }
    if (peiRes.status === 'fulfilled') {
      peiGoals.value = extractList(peiRes.value)
    }
  } catch (err) {
    console.error('Error loading support register data', err)
  } finally {
    loadingDiary.value = false
    loadingPei.value = false
  }
}

function onStudentChanged() {
  loadData()
}

function openNewDiaryDialog() {
  newDiaryDialog.value = true
}

function openNewGoalDialog() {
  newGoalDialog.value = true
}

function viewDiaryDetails(row) {
  $q.dialog({
    title: `${row.student_name} - ${row.entry_date} (${row.time_slot})`,
    message: `Attività: ${row.topic_and_activities}\n\nRisposta Studente: ${row.student_responses || 'N/A'}\n\nNote Educatore: ${row.educator_notes || 'N/A'}`,
    ok: 'Chiudi'
  })
}

async function submitDiaryEntry() {
  savingDiary.value = true
  try {
    await supportService.createDiaryEntry(diaryForm.value)
    $q.notify({ type: 'positive', message: 'Attività registrata nel diario di sostegno' })
    newDiaryDialog.value = false
    loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio della voce di diario' })
  } finally {
    savingDiary.value = false
  }
}

async function submitPeiGoal() {
  savingGoal.value = true
  try {
    await supportService.createPeiGoal(goalForm.value)
    $q.notify({ type: 'positive', message: 'Obiettivo PEI creato con successo' })
    newGoalDialog.value = false
    loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio dell\'obiettivo PEI' })
  } finally {
    savingGoal.value = false
  }
}

async function updateGoalStatus(id, status) {
  try {
    await supportService.updateGoalProgress(id, status)
    $q.notify({ type: 'positive', message: 'Avanzamento obiettivo aggiornato' })
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nell\'aggiornamento dell\'avanzamento' })
  }
}

onMounted(() => {
  loadInitialData()
})
</script>

<style scoped>
.support-page {
  max-width: 1300px;
  margin: 0 auto;
}
</style>
