<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('pnrrProjects.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('pnrrProjects.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          color="primary"
          icon="add_box"
          :label="$t('pnrrProjects.newProjectBtn')"
          class="rounded-lg shadow-sm"
          @click="showNewProjectDialog = true"
        />
        <q-btn
          outline
          color="secondary"
          icon="file_download"
          :label="$t('pnrrProjects.exportFuturaBtn')"
          class="rounded-lg shadow-sm"
          :disable="!selectedProject"
          @click="exportFuturaData"
        />
      </div>
    </div>

    <!-- Mandatory 75% Requirement Banner -->
    <q-banner rounded class="bg-indigo-900 text-white shadow-soft q-mb-lg">
      <template #avatar>
        <q-icon name="info" size="md" color="indigo-200" />
      </template>
      <div class="text-weight-bold text-subtitle1">
        Requisito Vincolante PNRR &amp; Piattaforma Futura MIM (Target 75%)
      </div>
      <div class="text-caption text-indigo-100">
        {{ $t('pnrrProjects.target75Notice') }} Ciascun corsista deve completare almeno il 75% delle ore del modulo per consentire la rendicontazione e liquidazione del costo standard unitario (UCS) al Ministero.
      </div>
    </q-banner>

    <!-- Projects Selection and Overview -->
    <div class="row q-col-gutter-lg">
      <!-- Left: Projects List -->
      <div class="col-12 col-lg-4">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">
              Progetti PNRR Attivi ({{ projects.length }})
            </div>

            <div v-if="projects.length === 0" class="q-pa-lg text-center text-slate-400">
              Nessun progetto PNRR registrato per l'A.S. corrente.
            </div>

            <q-list separator v-else class="rounded-lg">
              <q-item
                v-for="p in projects"
                :key="p.id"
                clickable
                v-ripple
                :active="selectedProject && selectedProject.id === p.id"
                active-class="bg-indigo-50 text-primary text-weight-bold"
                class="rounded-lg q-my-xs"
                @click="selectProject(p)"
              >
                <q-item-section>
                  <q-item-label class="text-subtitle2 text-weight-bold">
                    {{ p.title }}
                  </q-item-label>
                  <q-item-label caption>
                    CUP: <span class="text-weight-bold text-slate-700">{{ p.cup }}</span>
                  </q-item-label>
                  <q-item-label caption class="text-indigo-7">
                    {{ p.notice_type }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-badge color="primary">
                    {{ p.total_hours }}h
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Right: Project Management Tabs -->
      <div class="col-12 col-lg-8">
        <div v-if="!selectedProject" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
          <q-icon name="account_tree" size="64px" color="slate-300" class="q-mb-sm" />
          <div class="text-h6 text-weight-bold text-slate-600">Seleziona un progetto PNRR</div>
          <div class="text-body2 text-slate-400">
            Visualizza corsisti, registro firme lezioni e verifica la soglia vincolante del 75%
          </div>
        </div>

        <q-card v-else class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div>
                <div class="text-h5 text-weight-bold text-slate-800">
                  {{ selectedProject.title }}
                </div>
                <div class="text-caption text-slate-500">
                  CUP: <strong class="text-slate-800">{{ selectedProject.cup }}</strong> |
                  Monte Ore Totale: <strong class="text-primary">{{ selectedProject.total_hours }}h</strong> |
                  Target Corsisti: <strong class="text-slate-800">{{ selectedProject.target_students }}</strong>
                </div>
              </div>

              <div class="row q-gutter-xs">
                <q-btn
                  color="primary"
                  icon="person_add"
                  label="Iscrivi Corsista"
                  dense
                  class="rounded-lg"
                  @click="showEnrollDialog = true"
                />
                <q-btn
                  color="secondary"
                  icon="edit_calendar"
                  label="Nuova Lezione"
                  dense
                  class="rounded-lg"
                  @click="showLessonDialog = true"
                />
              </div>
            </div>

            <q-tabs
              v-model="activeTab"
              dense
              class="text-grey-7"
              active-color="primary"
              indicator-color="primary"
              align="left"
            >
              <q-tab name="students" icon="people" label="Corsisti &amp; Vincolo 75%" />
              <q-tab name="lessons" icon="history_edu" label="Registro Firme Lezioni" />
            </q-tabs>

            <q-separator class="q-mb-md" />

            <q-tab-panels v-model="activeTab" animated>
              <!-- Students Tab -->
              <q-tab-panel name="students" class="q-pa-none">
                <q-table
                  :rows="projectStudents"
                  :columns="studentColumns"
                  row-key="id"
                  flat
                  dense
                  :pagination="{ rowsPerPage: 15 }"
                  no-data-label="Nessun corsista iscritto a questo percorso"
                >
                  <template #body-cell-progress="props">
                    <q-td :props="props" style="min-width: 140px">
                      <div class="row items-center q-gutter-x-sm">
                        <q-linear-progress
                          :value="props.row.attendance_percentage / 100"
                          size="8px"
                          rounded
                          :color="props.row.target_75_met ? 'positive' : (props.row.attendance_percentage >= 50 ? 'warning' : 'negative')"
                          class="col"
                        />
                        <span class="text-caption text-weight-bold">
                          {{ Math.round(props.row.attendance_percentage) }}%
                        </span>
                      </div>
                    </q-td>
                  </template>

                  <template #body-cell-target="props">
                    <q-td :props="props" align="center">
                      <q-badge
                        :color="props.row.target_75_met ? 'positive' : 'negative'"
                        class="text-weight-bold"
                      >
                        {{ props.row.target_75_met ? 'CONVALIDATO (>=75%)' : 'SOTTO SOGLIA' }}
                      </q-badge>
                    </q-td>
                  </template>
                </q-table>
              </q-tab-panel>

              <!-- Lessons Tab -->
              <q-tab-panel name="lessons" class="q-pa-none">
                <q-table
                  :rows="projectLessons"
                  :columns="lessonColumns"
                  row-key="id"
                  flat
                  dense
                  :pagination="{ rowsPerPage: 15 }"
                  no-data-label="Nessuna lezione registrata"
                />
              </q-tab-panel>
            </q-tab-panels>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Dialog: Create Project -->
    <q-dialog v-model="showNewProjectDialog" persistent>
      <q-card style="min-width: 500px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            {{ $t('pnrrProjects.newProjectBtn') }}
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newProject.title"
            outlined
            dense
            label="Titolo Modulo Formativo / Progetto *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newProject.cup"
            outlined
            dense
            label="Codice Unico di Progetto (CUP) *"
            hint="Es. B84D23000550006"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-select
            v-model="newProject.notice_type"
            :options="noticeTypeOptions"
            outlined
            dense
            label="Avviso / Investimento PNRR *"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input
                v-model.number="newProject.total_hours"
                type="number"
                outlined
                dense
                label="Ore Totali Percorso *"
              />
            </div>
            <div class="col-6">
              <q-input
                v-model.number="newProject.target_students"
                type="number"
                outlined
                dense
                label="Target Minimo Corsisti"
              />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="save"
            label="Crea Progetto"
            :disable="!newProject.title || !newProject.cup"
            @click="createProject"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Enroll Student -->
    <q-dialog v-model="showEnrollDialog" persistent>
      <q-card style="min-width: 400px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            Iscrizione Corsista al Progetto
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="enrollStudentId"
            outlined
            dense
            label="ID Studente / Codice Fiscale Corsista *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            label="Iscrivi Corsista"
            :disable="!enrollStudentId"
            @click="enrollStudent"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Record Lesson -->
    <q-dialog v-model="showLessonDialog" persistent>
      <q-card style="min-width: 520px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            {{ $t('pnrrProjects.recordLessonBtn') }}
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input
                v-model="newLesson.lesson_date"
                type="date"
                outlined
                dense
                label="Data Lezione *"
              />
            </div>
            <div class="col-6">
              <q-input
                v-model.number="newLesson.hours"
                type="number"
                outlined
                dense
                label="Ore Lezione *"
              />
            </div>
          </div>

          <q-input
            v-model="newLesson.topic"
            outlined
            dense
            label="Argomento / Attività Svolta *"
          />

          <div class="text-subtitle2 text-slate-700 q-mt-sm">
            Seleziona Corsisti Presenti ({{ newLesson.attended_student_ids.length }} / {{ projectStudents.length }}):
          </div>

          <q-list dense bordered class="rounded-lg max-h-48 overflow-y-auto">
            <q-item v-for="s in projectStudents" :key="s.id" tag="label" v-ripple>
              <q-item-section side top>
                <q-checkbox v-model="newLesson.attended_student_ids" :val="s.student_id" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ s.student_id }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="verified"
            label="Firma e Salva Lezione"
            :disable="!newLesson.topic || newLesson.hours <= 0"
            @click="saveLesson"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const projects = ref([])
const selectedProject = ref(null)
const activeTab = ref('students')

const projectStudents = ref([])
const projectLessons = ref([])

const showNewProjectDialog = ref(false)
const newProject = ref({
  title: '',
  cup: '',
  notice_type: 'D.M. 65 STEM & Multilinguismo',
  total_hours: 30,
  target_students: 15,
  academic_year: '2025/2026'
})

const noticeTypeOptions = [
  'D.M. 65 STEM & Multilinguismo',
  'D.M. 66 Transizione Digitale',
  'D.M. 19 Riduzione Divari & Dispersione',
  'Piano Scuola 4.0 - Next Generation Classrooms'
]

const showEnrollDialog = ref(false)
const enrollStudentId = ref('')

const showLessonDialog = ref(false)
const newLesson = ref({
  lesson_date: new Date().toISOString().substring(0, 10),
  hours: 2,
  topic: '',
  attended_student_ids: []
})

const studentColumns = [
  { name: 'student_id', label: 'ID / Corsista', field: 'student_id', align: 'left', sortable: true },
  { name: 'attended_hours', label: 'Ore Frequentate', field: 'attended_hours', align: 'center', sortable: true },
  { name: 'progress', label: 'Avanzamento Target', align: 'center' },
  { name: 'target', label: 'Stato Target 75%', align: 'center' }
]

const lessonColumns = [
  { name: 'lesson_date', label: 'Data', field: 'lesson_date', align: 'left', sortable: true },
  { name: 'hours', label: 'Ore', field: 'hours', align: 'center' },
  { name: 'topic', label: 'Argomento', field: 'topic', align: 'left' }
]

async function loadProjects() {
  try {
    const res = await api.get('/pnrr-projects')
    projects.value = res.data?.projects || []
    if (projects.value.length > 0 && !selectedProject.value) {
      selectProject(projects.value[0])
    }
  } catch (err) {
    console.error('Errore caricamento progetti PNRR:', err)
  }
}

async function selectProject(p) {
  selectedProject.value = p
  await Promise.all([loadStudents(p.id), loadLessons(p.id)])
}

async function loadStudents(projectId) {
  try {
    const res = await api.get(`/pnrr-projects/${projectId}/students`)
    projectStudents.value = res.data?.students || []
  } catch (err) {
    console.error('Errore caricamento corsisti:', err)
  }
}

async function loadLessons(projectId) {
  try {
    const res = await api.get(`/pnrr-projects/${projectId}/lessons`)
    projectLessons.value = res.data?.lessons || []
  } catch (err) {
    console.error('Errore caricamento lezioni:', err)
  }
}

async function createProject() {
  try {
    await api.post('/pnrr-projects', newProject.value)
    $q.notify({ type: 'positive', message: 'Progetto PNRR registrato con successo' })
    showNewProjectDialog.value = false
    await loadProjects()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore creazione progetto' })
  }
}

async function enrollStudent() {
  if (!selectedProject.value || !enrollStudentId.value) return
  try {
    await api.post(`/pnrr-projects/${selectedProject.value.id}/enroll`, {
      student_id: enrollStudentId.value
    })
    $q.notify({ type: 'positive', message: 'Corsista iscritto al progetto' })
    showEnrollDialog.value = false
    enrollStudentId.value = ''
    await loadStudents(selectedProject.value.id)
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore iscrizione corsista' })
  }
}

async function saveLesson() {
  if (!selectedProject.value) return
  try {
    await api.post(`/pnrr-projects/${selectedProject.value.id}/lessons`, {
      ...newLesson.value,
      signatures_json: JSON.stringify({ certified_by: 'docente', timestamp: new Date().toISOString() })
    })
    $q.notify({ type: 'positive', message: 'Lezione registrata e presenze certificate' })
    showLessonDialog.value = false
    newLesson.value.topic = ''
    newLesson.value.attended_student_ids = []
    await Promise.all([
      loadStudents(selectedProject.value.id),
      loadLessons(selectedProject.value.id)
    ])
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore registrazione lezione' })
  }
}

function exportFuturaData() {
  if (!selectedProject.value) return
  const data = {
    cup: selectedProject.value.cup,
    titolo: selectedProject.value.title,
    avviso: selectedProject.value.notice_type,
    ore_totali: selectedProject.value.total_hours,
    corsisti: projectStudents.value.map(s => ({
      corsista_id: s.student_id,
      ore_frequentate: s.attended_hours,
      percentuale_frequenza: `${Math.round(s.attendance_percentage)}%`,
      target_liquidabile: s.target_75_met ? 'SI' : 'NO'
    })),
    data_esportazione: new Date().toISOString()
  }
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `PNRR_FUTURA_${selectedProject.value.cup}.json`
  a.click()
  URL.revokeObjectURL(url)
  $q.notify({ type: 'positive', message: 'File di rendicontazione Piattaforma Futura esportato!' })
}

onMounted(() => {
  loadProjects()
})
</script>
