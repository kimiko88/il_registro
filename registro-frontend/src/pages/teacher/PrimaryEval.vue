<template>
  <q-page class="q-pa-md" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-grey-1 text-dark'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bolder row items-center text-primary">
          <q-icon name="menu_book" class="q-mr-sm" size="28px" />
          {{ t('primaryEval.title') }}
        </div>
        <div class="text-caption" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">
          {{ t('primaryEval.subtitle') }}
        </div>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-btn
          color="primary"
          icon="add_task"
          :label="t('primaryEval.addObjective')"
          unelevated
          dense
          class="q-px-sm rounded-lg text-weight-bold"
          @click="openAddObjectiveDialog"
          :disable="!selectedSubjectId"
        />
        <q-btn
          color="secondary"
          icon="file_download"
          :label="t('primaryEval.exportReport')"
          outline
          dense
          class="q-px-sm rounded-lg text-weight-bold"
          @click="exportCSV"
          :disable="!matrix || matrix.students.length === 0"
        />
        <q-btn
          color="positive"
          icon="save"
          :label="t('primaryEval.saveBatch')"
          unelevated
          dense
          class="q-px-md rounded-lg text-weight-bold shadow-soft"
          :loading="saving"
          @click="saveAllEvaluations"
          :disable="!hasUnsavedChanges || !selectedClassId || !selectedSubjectId"
        >
          <q-badge v-if="hasUnsavedChanges" color="orange-9" floating rounded>•</q-badge>
        </q-btn>
      </div>
    </div>

    <!-- Filters Bar -->
    <q-card flat bordered class="rounded-xl q-mb-md shadow-1" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-white'">
      <q-card-section class="row q-col-gutter-sm items-center q-py-sm">
        <div class="col-12 col-sm-4 col-md-3">
          <q-select
            v-model="selectedClassId"
            :options="primaryClassOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            dense outlined
            :label="t('primaryEval.selectClass')"
            :loading="classesStore.loading"
            @update:model-value="onClassChanged"
          >
            <template v-slot:prepend>
              <q-icon name="class" color="primary" />
            </template>
          </q-select>
        </div>

        <div class="col-12 col-sm-4 col-md-3">
          <q-select
            v-model="selectedSubjectId"
            :options="subjectOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            dense outlined
            :label="t('primaryEval.selectSubject')"
            :loading="loadingSubjects"
            @update:model-value="loadMatrixData"
          >
            <template v-slot:prepend>
              <q-icon name="subject" color="primary" />
            </template>
          </q-select>
        </div>

        <div class="col-12 col-sm-4 col-md-2">
          <q-select
            v-model="selectedSemester"
            :options="[
              { label: t('primaryEval.semester1'), value: 1 },
              { label: t('primaryEval.semester2'), value: 2 }
            ]"
            emit-value
            map-options
            dense outlined
            :label="t('primaryEval.semester')"
            @update:model-value="loadMatrixData"
          >
            <template v-slot:prepend>
              <q-icon name="event" color="primary" />
            </template>
          </q-select>
        </div>

        <div class="col-12 col-sm-6 col-md-2">
          <q-input
            v-model="evaluationDate"
            dense outlined
            type="date"
            :label="t('primaryEval.evalDate')"
          />
        </div>

        <div class="col-12 col-sm-6 col-md-2 text-right">
          <q-btn-toggle
            v-model="activeView"
            spread
            dense
            unelevated
            toggle-color="primary"
            :options="[
              { label: t('primaryEval.viewMatrix'), value: 'matrix', icon: 'grid_on' },
              { label: t('primaryEval.viewObjectives'), value: 'objectives', icon: 'format_list_bulleted' }
            ]"
          />
        </div>
      </q-card-section>
    </q-card>

    <!-- 4 Ministerial Levels Legend Card (O.M. 172/2020) -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-6 col-md-3" v-for="lvl in ministerialLevels" :key="lvl.value">
        <q-card flat bordered class="rounded-xl shadow-xs" :class="lvl.bgClass">
          <q-card-section class="q-pa-sm">
            <div class="row items-center justify-between no-wrap">
              <span class="text-weight-bold" :class="lvl.textClass">{{ lvl.label }}</span>
              <q-badge :color="lvl.color" text-color="white" class="text-caption font-mono">
                {{ lvl.code }}
              </q-badge>
            </div>
            <div class="text-caption text-grey-8 q-mt-xs ellipsis-2-lines">
              {{ lvl.desc }}
              <q-tooltip anchor="top middle" self="bottom middle" max-width="320px" class="bg-slate-900">
                <strong>{{ lvl.label }}</strong>: {{ lvl.desc }}
              </q-tooltip>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Area: Matrix vs Objectives -->
    <div v-if="loadingMatrix" class="text-center q-pa-xl">
      <q-spinner color="primary" size="48px" />
      <div class="text-subtitle2 text-grey-6 q-mt-md">{{ t('primaryEval.loadingMatrix') }}</div>
    </div>

    <!-- Empty Class State -->
    <div v-else-if="!selectedClassId" class="text-center q-pa-xl bg-white rounded-xl shadow-1">
      <q-icon name="school" size="64px" color="grey-5" />
      <div class="text-h6 text-grey-7 q-mt-md">{{ t('primaryEval.promptSelectClass') }}</div>
      <div class="text-caption text-grey-5">{{ t('primaryEval.promptSelectClassDesc') }}</div>
    </div>

    <!-- Objectives Management View -->
    <div v-else-if="activeView === 'objectives'">
      <q-card flat bordered class="rounded-xl shadow-1">
        <q-card-section class="row items-center justify-between bg-slate-50 border-bottom">
          <div class="text-subtitle1 text-weight-bold text-slate-800">
            <q-icon name="playlist_add_check" color="primary" size="22px" class="q-mr-xs" />
            {{ t('primaryEval.objectivesListTitle') }} ({{ objectives.length }})
          </div>
          <q-btn
            color="primary"
            icon="add"
            :label="t('primaryEval.newObjective')"
            dense unelevated
            class="q-px-sm rounded-lg"
            @click="openAddObjectiveDialog"
          />
        </q-card-section>

        <q-list separator>
          <q-item v-for="(obj, idx) in objectives" :key="obj.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="indigo-1" text-color="indigo-9" size="36px" class="text-weight-bold">
                {{ idx + 1 }}
              </q-avatar>
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-bold text-subtitle2">{{ obj.title }}</q-item-label>
              <q-item-label caption v-if="obj.description" class="text-grey-7">
                {{ obj.description }}
              </q-item-label>
              <div class="row q-gutter-x-sm q-mt-xs">
                <q-badge color="grey-3" text-color="grey-9">
                  {{ t('primaryEval.gradeYear') }} {{ obj.year_grade }}°
                </q-badge>
                <q-badge color="grey-3" text-color="grey-9" v-if="obj.academic_year">
                  A.S. {{ obj.academic_year }}
                </q-badge>
              </div>
            </q-item-section>

            <q-item-section side>
              <q-btn
                flat round dense
                color="negative"
                icon="delete_outline"
                @click="confirmDeleteObjective(obj)"
              >
                <q-tooltip>{{ t('common.delete') }}</q-tooltip>
              </q-btn>
            </q-item-section>
          </q-item>

          <q-item v-if="objectives.length === 0" class="text-center q-pa-lg text-grey-6">
            <q-item-section>
              <q-icon name="lightbulb" size="32px" color="amber-8" class="q-mx-auto q-mb-sm" />
              <div>{{ t('primaryEval.noObjectivesFound') }}</div>
              <div class="text-caption">{{ t('primaryEval.noObjectivesHint') }}</div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- Matrix View (Students x Objectives) -->
    <div v-else>
      <div v-if="objectives.length === 0" class="q-pa-xl text-center bg-white rounded-xl shadow-1">
        <q-icon name="auto_stories" size="64px" color="indigo-4" />
        <div class="text-h6 text-grey-8 q-mt-md">{{ t('primaryEval.noObjectivesForMatrix') }}</div>
        <div class="text-caption text-grey-6 q-mb-md">{{ t('primaryEval.noObjectivesForMatrixDesc') }}</div>
        <q-btn
          color="primary"
          icon="add"
          :label="t('primaryEval.addFirstObjective')"
          unelevated
          @click="openAddObjectiveDialog"
        />
      </div>

      <q-card v-else flat bordered class="rounded-xl shadow-1 overflow-hidden">
        <!-- Matrix Toolbar -->
        <q-toolbar class="bg-slate-100 text-slate-800 q-px-md border-bottom">
          <q-icon name="table_chart" color="primary" class="q-mr-sm" size="20px" />
          <span class="text-subtitle2 text-weight-bold">
            {{ t('primaryEval.matrixGridTitle') }} — {{ currentClassName }} ({{ matrix?.students?.length || 0 }} alunni)
          </span>
          <q-space />
          <q-input
            v-model="studentFilter"
            dense outlined
            bg-color="white"
            placeholder="Cerca alunno..."
            class="q-mr-sm"
            style="max-width: 180px;"
          >
            <template v-slot:append>
              <q-icon name="search" size="xs" />
            </template>
          </q-input>
        </q-toolbar>

        <!-- Interactive Matrix Table -->
        <div class="table-responsive">
          <table class="primary-matrix-table">
            <thead>
              <tr>
                <th class="sticky-col student-header text-left">
                  {{ t('primaryEval.studentName') }}
                </th>
                <th v-for="(obj, oIdx) in objectives" :key="obj.id" class="objective-header text-center">
                  <div class="text-weight-bold text-caption ellipsis-2-lines" :title="obj.title">
                    {{ oIdx + 1 }}. {{ obj.title }}
                  </div>
                  <q-tooltip max-width="300px" class="bg-slate-900">
                    <strong>Obiettivo {{ oIdx + 1 }}:</strong><br>{{ obj.title }}
                  </q-tooltip>
                </th>
                <th class="summary-col text-center">
                  {{ t('primaryEval.prevalentLevel') }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="student in filteredStudents" :key="student.student_id" class="student-row">
                <td class="sticky-col student-name-cell">
                  <div class="row items-center no-wrap">
                    <q-avatar size="28px" color="indigo-1" text-color="indigo-9" class="q-mr-xs text-weight-bold">
                      {{ student.student_name ? student.student_name.charAt(0) : '?' }}
                    </q-avatar>
                    <span class="text-weight-bold ellipsis" :title="student.student_name">
                      {{ student.student_name }}
                    </span>
                  </div>
                </td>

                <!-- Cells per objective -->
                <td v-for="obj in objectives" :key="obj.id" class="matrix-cell text-center">
                  <q-btn-dropdown
                    dense
                    unelevated
                    no-caps
                    :color="getCellColor(student.student_id, obj.id)"
                    :text-color="getCellTextColor(student.student_id, obj.id)"
                    :label="getCellLabel(student.student_id, obj.id)"
                    class="level-dropdown-btn text-weight-bold"
                  >
                    <q-list dense style="min-width: 240px;">
                      <q-item-label header class="text-weight-bold text-caption text-primary">
                        {{ t('primaryEval.chooseLevel') }}
                      </q-item-label>
                      <q-item
                        v-for="lvl in ministerialLevels"
                        :key="lvl.value"
                        clickable
                        v-close-popup
                        @click="setStudentLevel(student.student_id, obj.id, lvl.value)"
                        :class="{ 'bg-indigo-1': isCurrentLevel(student.student_id, obj.id, lvl.value) }"
                      >
                        <q-item-section avatar>
                          <q-badge :color="lvl.color" text-color="white">{{ lvl.code }}</q-badge>
                        </q-item-section>
                        <q-item-section>
                          <q-item-label class="text-weight-bold">{{ lvl.label }}</q-item-label>
                          <q-item-label caption class="ellipsis">{{ lvl.desc }}</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-separator />
                      <q-item clickable v-close-popup @click="openDimensionsDialog(student, obj)">
                        <q-item-section avatar>
                          <q-icon name="tune" color="indigo" />
                        </q-item-section>
                        <q-item-section>
                          <q-item-label class="text-weight-bold text-indigo">
                            {{ t('primaryEval.editDimensions') }}
                          </q-item-label>
                          <q-item-label caption>{{ t('primaryEval.editDimensionsCaption') }}</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item clickable v-close-popup @click="clearStudentLevel(student.student_id, obj.id)">
                        <q-item-section avatar>
                          <q-icon name="clear" color="grey" />
                        </q-item-section>
                        <q-item-section class="text-grey-7">{{ t('primaryEval.clearCell') }}</q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                </td>

                <!-- Prevalent Level Summary Cell -->
                <td class="summary-col text-center">
                  <q-badge
                    v-if="computeStudentSummary(student.student_id)"
                    :color="computeStudentSummary(student.student_id).color"
                    text-color="white"
                    class="q-pa-xs text-weight-bold"
                  >
                    {{ computeStudentSummary(student.student_id).label }}
                  </q-badge>
                  <span v-else class="text-caption text-grey-4">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </q-card>
    </div>

    <!-- Dimensions Details Dialog (4 ministerial dimensions) -->
    <q-dialog v-model="showDimensionsDialog">
      <q-card style="width: min(540px, 95vw)" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div>
            <div class="text-subtitle1 text-weight-bold">
              {{ t('primaryEval.dialogDimensionsTitle') }}
            </div>
            <div class="text-caption text-indigo-1">
              {{ activeDimensionStudent?.student_name }} — {{ activeDimensionObjective?.title }}
            </div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <q-select
            v-model="dimensionForm.level"
            :options="ministerialLevels"
            option-value="value"
            option-label="label"
            emit-value
            map-options
            outlined dense
            :label="t('primaryEval.dialogLevelLabel')"
          />

          <!-- Dimension 1: Autonomia -->
          <div>
            <div class="text-caption text-weight-bold q-mb-xs">{{ t('primaryEval.dimAutonomy') }}</div>
            <q-btn-toggle
              v-model="dimensionForm.dimension_autonomy"
              spread dense unelevated
              toggle-color="primary"
              :options="[
                { label: t('primaryEval.autonomous'), value: 'autonomo' },
                { label: t('primaryEval.guided'), value: 'con_guida' }
              ]"
            />
          </div>

          <!-- Dimension 2: Continuità -->
          <div>
            <div class="text-caption text-weight-bold q-mb-xs">{{ t('primaryEval.dimContinuity') }}</div>
            <q-btn-toggle
              v-model="dimensionForm.dimension_continuity"
              spread dense unelevated
              toggle-color="primary"
              :options="[
                { label: t('primaryEval.continuous'), value: 'continuo' },
                { label: t('primaryEval.nonContinuous'), value: 'non_continuo' }
              ]"
            />
          </div>

          <!-- Dimension 3: Familiarità situazione -->
          <div>
            <div class="text-caption text-weight-bold q-mb-xs">{{ t('primaryEval.dimFamiliarity') }}</div>
            <q-btn-toggle
              v-model="dimensionForm.dimension_familiarity"
              spread dense unelevated
              toggle-color="primary"
              :options="[
                { label: t('primaryEval.familiar'), value: 'nota' },
                { label: t('primaryEval.unfamiliar'), value: 'non_nota' }
              ]"
            />
          </div>

          <!-- Dimension 4: Risorse impiegate -->
          <div>
            <div class="text-caption text-weight-bold q-mb-xs">{{ t('primaryEval.dimResources') }}</div>
            <q-btn-toggle
              v-model="dimensionForm.dimension_resources"
              spread dense unelevated
              toggle-color="primary"
              :options="[
                { label: t('primaryEval.ownResources'), value: 'risorse_proprie' },
                { label: t('primaryEval.providedResources'), value: 'risorse_fornite' }
              ]"
            />
          </div>

          <!-- Notes -->
          <q-input
            v-model="dimensionForm.notes"
            outlined dense autogrow
            :label="t('primaryEval.dialogNotes')"
            placeholder="Eventuali annotazioni descrittive individuali..."
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            :label="t('common.confirm')"
            unelevated
            class="q-px-md rounded-lg font-bold"
            @click="saveDimensionForm"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Create Learning Objective Dialog -->
    <q-dialog v-model="showAddObjectiveDialog">
      <q-card style="width: min(500px, 95vw)" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold">{{ t('primaryEval.addObjective') }}</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-sm">
          <q-input
            v-model="newObjective.title"
            :label="t('primaryEval.objectiveTitle') + ' *'"
            outlined dense
            type="textarea"
            rows="3"
            placeholder="Es. Eseguire addizioni e sottrazioni con i numeri naturali con consapevolezza del significato del calcolo."
            autofocus
          />
          <q-input
            v-model="newObjective.description"
            :label="t('primaryEval.objectiveDesc')"
            outlined dense
            type="textarea"
            rows="2"
            placeholder="Nucleo tematico / traguardi di competenza collegati..."
          />
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select
                v-model="newObjective.year_grade"
                :options="[1, 2, 3, 4, 5]"
                :label="t('primaryEval.gradeYear') + ' (1-5) *'"
                outlined dense
              />
            </div>
            <div class="col-6">
              <q-input
                v-model="newObjective.academic_year"
                :label="t('primaryEval.academicYear') + ' *'"
                outlined dense
                placeholder="2025/2026"
              />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            :label="t('common.save')"
            unelevated
            :loading="creatingObjective"
            :disable="!newObjective.title.trim()"
            @click="createObjective"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import { useSchoolStore } from '@/stores/schools'
import primaryEvalService from '@/services/primaryEvalService'
import { useOfflineSync } from '@/composables/useOfflineSync'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()
const classesStore = useClassesStore()
const schoolStore = useSchoolStore()
const { executeWithOfflineQueue } = useOfflineSync()

// State
const selectedClassId = ref(null)
const selectedSubjectId = ref(null)
const selectedSemester = ref(1)
const evaluationDate = ref(new Date().toISOString().substring(0, 10))
const activeView = ref('matrix') // 'matrix' | 'objectives'
const studentFilter = ref('')

const loadingMatrix = ref(false)
const loadingSubjects = ref(false)
const saving = ref(false)
const creatingObjective = ref(false)

const matrix = ref(null)
const objectives = ref([])
const subjectsList = ref([])

// In-memory changes: key = `${studentId}_${objectiveId}` -> cell evaluation object
const localEvaluations = reactive({})
const hasUnsavedChanges = ref(false)

// Modals
const showAddObjectiveDialog = ref(false)
const showDimensionsDialog = ref(false)
const activeDimensionStudent = ref(null)
const activeDimensionObjective = ref(null)

const newObjective = reactive({
  title: '',
  description: '',
  year_grade: 1,
  academic_year: '2025/2026'
})

const dimensionForm = reactive({
  level: 'intermedio',
  dimension_autonomy: 'autonomo',
  dimension_continuity: 'continuo',
  dimension_familiarity: 'nota',
  dimension_resources: 'risorse_proprie',
  notes: ''
})

// 4 Ministerial Levels (O.M. 172/2020)
const ministerialLevels = [
  {
    value: 'avanzato',
    label: 'Avanzato',
    code: 'AV',
    color: 'positive',
    bgClass: 'bg-emerald-50 border-emerald-300',
    textClass: 'text-emerald-900',
    desc: 'L’alunno porta a termine compiti in situazioni note e non note, mobilitando una varietà di risorse con continuità e autonomia.'
  },
  {
    value: 'intermedio',
    label: 'Intermedio',
    code: 'INT',
    color: 'primary',
    bgClass: 'bg-blue-50 border-blue-300',
    textClass: 'text-blue-900',
    desc: 'L’alunno porta a termine compiti in situazioni note in modo autonomo e continuo; in situazioni non note utilizza risorse fornite.'
  },
  {
    value: 'base',
    label: 'Base',
    code: 'BASE',
    color: 'warning',
    bgClass: 'bg-amber-50 border-amber-300',
    textClass: 'text-amber-900',
    desc: 'L’alunno porta a termine compiti solo in situazioni note e con il supporto del docente o con risorse fornite appositamente.'
  },
  {
    value: 'in_via_di_prima_acquisizione',
    label: 'In via di 1ª acq.',
    code: '1ª ACQ',
    color: 'negative',
    bgClass: 'bg-red-50 border-red-300',
    textClass: 'text-red-900',
    desc: 'L’alunno porta a termine compiti solo se guidato continuativamente e unicamente con risorse fornite.'
  }
]

// Computed: Filter primary classes
const primaryClassOptions = computed(() => {
  return classesStore.classes.filter(c => {
    if (c.school_level === 'primaria') return true
    if (c.name && c.name.toLowerCase().includes('primaria')) return true
    if (schoolStore.hasPrimaryLevels) return true
    return false
  })
})

const currentClassName = computed(() => {
  const c = classesStore.classes.find(item => item.id === selectedClassId.value)
  return c?.name || c?.label || 'Classe'
})

const subjectOptions = computed(() => {
  return subjectsList.value.map(s => ({
    id: s.id || s.subject_id,
    name: s.name || s.subject_name
  }))
})

const filteredStudents = computed(() => {
  if (!matrix.value?.students) return []
  if (!studentFilter.value.trim()) return matrix.value.students
  const q = studentFilter.value.toLowerCase()
  return matrix.value.students.filter(s => s.student_name?.toLowerCase().includes(q))
})

// Lifecycle
onMounted(async () => {
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  if (primaryClassOptions.value.length > 0) {
    selectedClassId.value = primaryClassOptions.value[0].id
  }
  await fetchSubjects()
  if (selectedClassId.value && selectedSubjectId.value) {
    await loadMatrixData()
  }
})

async function fetchSubjects() {
  loadingSubjects.value = true
  try {
    const res = await api.get('/subjects')
    subjectsList.value = res.data || []
    if (subjectsList.value.length > 0) {
      selectedSubjectId.value = subjectsList.value[0].id
    }
  } catch (err) {
    console.error('Error fetching subjects', err)
  } finally {
    loadingSubjects.value = false
  }
}

async function onClassChanged() {
  await loadMatrixData()
}

async function loadMatrixData() {
  if (!selectedClassId.value || !selectedSubjectId.value) return
  loadingMatrix.value = true
  try {
    const res = await primaryEvalService.getMatrix({
      class_id: selectedClassId.value,
      subject_id: selectedSubjectId.value,
      semester: selectedSemester.value
    })
    matrix.value = res.data?.data || res.data
    objectives.value = matrix.value.objectives || []

    // Populate localEvaluations cache
    Object.keys(localEvaluations).forEach(k => delete localEvaluations[k])
    if (matrix.value.students) {
      for (const st of matrix.value.students) {
        if (st.evaluations) {
          for (const [objId, cell] of Object.entries(st.evaluations)) {
            localEvaluations[`${st.student_id}_${objId}`] = { ...cell }
          }
        }
      }
    }
    hasUnsavedChanges.value = false
  } catch (err) {
    console.error('Error loading primary matrix', err)
    $q.notify({
      type: 'negative',
      message: t('primaryEval.matrixFetchError')
    })
  } finally {
    loadingMatrix.value = false
  }
}

// Cell helpers
function getCellKey(studentId, objId) {
  return `${studentId}_${objId}`
}

function isCurrentLevel(studentId, objId, levelVal) {
  const cell = localEvaluations[getCellKey(studentId, objId)]
  return cell?.level === levelVal
}

function getCellColor(studentId, objId) {
  const cell = localEvaluations[getCellKey(studentId, objId)]
  if (!cell || !cell.level) return 'grey-3'
  const found = ministerialLevels.find(l => l.value === cell.level)
  return found ? found.color : 'grey-3'
}

function getCellTextColor(studentId, objId) {
  const cell = localEvaluations[getCellKey(studentId, objId)]
  if (!cell || !cell.level) return 'grey-7'
  return 'white'
}

function getCellLabel(studentId, objId) {
  const cell = localEvaluations[getCellKey(studentId, objId)]
  if (!cell || !cell.level) return '—'
  const found = ministerialLevels.find(l => l.value === cell.level)
  return found ? found.code : '—'
}

function setStudentLevel(studentId, objId, levelVal) {
  const k = getCellKey(studentId, objId)
  const existing = localEvaluations[k] || {}
  localEvaluations[k] = {
    ...existing,
    level: levelVal,
    date: evaluationDate.value
  }
  hasUnsavedChanges.value = true
}

function clearStudentLevel(studentId, objId) {
  const k = getCellKey(studentId, objId)
  if (localEvaluations[k]) {
    delete localEvaluations[k]
    hasUnsavedChanges.value = true
  }
}

function computeStudentSummary(studentId) {
  const counts = {}
  let total = 0
  for (const obj of objectives.value) {
    const cell = localEvaluations[getCellKey(studentId, obj.id)]
    if (cell && cell.level) {
      counts[cell.level] = (counts[cell.level] || 0) + 1
      total++
    }
  }
  if (total === 0) return null

  // Find level with max count
  let max = 0
  let chosen = null
  for (const lvl of ministerialLevels) {
    if ((counts[lvl.value] || 0) > max) {
      max = counts[lvl.value]
      chosen = lvl
    }
  }
  return chosen
}

// Dimensions Dialog
function openDimensionsDialog(student, obj) {
  activeDimensionStudent.value = student
  activeDimensionObjective.value = obj
  const cell = localEvaluations[getCellKey(student.student_id, obj.id)] || {}
  dimensionForm.level = cell.level || 'intermedio'
  dimensionForm.dimension_autonomy = cell.dimension_autonomy || 'autonomo'
  dimensionForm.dimension_continuity = cell.dimension_continuity || 'continuo'
  dimensionForm.dimension_familiarity = cell.dimension_familiarity || 'nota'
  dimensionForm.dimension_resources = cell.dimension_resources || 'risorse_proprie'
  dimensionForm.notes = cell.notes || ''
  showDimensionsDialog.value = true
}

function saveDimensionForm() {
  if (!activeDimensionStudent.value || !activeDimensionObjective.value) return
  const k = getCellKey(activeDimensionStudent.value.student_id, activeDimensionObjective.value.id)
  localEvaluations[k] = {
    ...localEvaluations[k],
    ...dimensionForm,
    date: evaluationDate.value
  }
  hasUnsavedChanges.value = true
  showDimensionsDialog.value = false
  $q.notify({
    type: 'positive',
    message: t('primaryEval.dimensionsUpdated')
  })
}

// Save batch
async function saveAllEvaluations() {
  if (!selectedClassId.value || !selectedSubjectId.value) return
  saving.value = true
  try {
    // Collect batch operations per objective
    for (const obj of objectives.value) {
      const items = []
      for (const st of (matrix.value?.students || [])) {
        const cell = localEvaluations[getCellKey(st.student_id, obj.id)]
        if (cell && cell.level) {
          items.push({
            student_id: st.student_id,
            level: cell.level,
            dimension_autonomy: cell.dimension_autonomy || 'autonomo',
            dimension_continuity: cell.dimension_continuity || 'continuo',
            dimension_familiarity: cell.dimension_familiarity || 'nota',
            dimension_resources: cell.dimension_resources || 'risorse_proprie',
            notes: cell.notes || ''
          })
        }
      }
      if (items.length > 0) {
        const payload = {
          class_id: selectedClassId.value,
          subject_id: selectedSubjectId.value,
          objective_id: obj.id,
          date: evaluationDate.value,
          semester: Number(selectedSemester.value),
          evaluations: items
        }
        await executeWithOfflineQueue(
          {
            url: '/primary/evaluations/batch',
            method: 'post',
            data: payload
          },
          { title: `Valutazioni: ${obj.title}` }
        )
      }
    }
    hasUnsavedChanges.value = false
    $q.notify({
      type: 'positive',
      message: t('primaryEval.evaluationsSavedSuccess')
    })
    await loadMatrixData()
  } catch (err) {
    console.error('Error saving evaluations batch', err)
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('primaryEval.evaluationsSaveError')
    })
  } finally {
    saving.value = false
  }
}

// Objective CRUD
function openAddObjectiveDialog() {
  const cls = classesStore.classes.find(c => c.id === selectedClassId.value)
  // Infer year grade from class name if possible (e.g. 3A -> 3)
  const match = cls?.name?.match(/^(\d)/)
  newObjective.year_grade = match ? parseInt(match[1]) : 1
  newObjective.title = ''
  newObjective.description = ''
  newObjective.academic_year = '2025/2026'
  showAddObjectiveDialog.value = true
}

async function createObjective() {
  if (!newObjective.title.trim() || !selectedSubjectId.value) return
  creatingObjective.value = true
  try {
    await primaryEvalService.createObjective({
      class_id: selectedClassId.value,
      subject_id: selectedSubjectId.value,
      year_grade: Number(newObjective.year_grade),
      title: newObjective.title.trim(),
      description: newObjective.description?.trim() || '',
      academic_year: newObjective.academic_year || '2025/2026'
    })
    showAddObjectiveDialog.value = false
    $q.notify({
      type: 'positive',
      message: t('primaryEval.objectiveCreatedSuccess')
    })
    await loadMatrixData()
  } catch (err) {
    console.error('Error creating objective', err)
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('primaryEval.objectiveCreateError')
    })
  } finally {
    creatingObjective.value = false
  }
}

function confirmDeleteObjective(obj) {
  $q.dialog({
    title: t('primaryEval.confirmDeleteObjTitle'),
    message: t('primaryEval.confirmDeleteObjMsg', { title: obj.title }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await primaryEvalService.deleteObjective(obj.id)
      $q.notify({
        type: 'positive',
        message: t('primaryEval.objectiveDeletedSuccess')
      })
      await loadMatrixData()
    } catch (err) {
      $q.notify({
        type: 'negative',
        message: t('primaryEval.objectiveDeleteError')
      })
    }
  })
}

// Export CSV
function exportCSV() {
  if (!matrix.value || !matrix.value.students) return
  let csv = 'Alunno;'
  objectives.value.forEach((o, i) => {
    csv += `"${i + 1}. ${o.title}";`
  })
  csv += 'Livello Prevalente\n'

  for (const st of matrix.value.students) {
    csv += `"${st.student_name}";`
    for (const obj of objectives.value) {
      const cell = localEvaluations[getCellKey(st.student_id, obj.id)]
      csv += `"${cell?.level || '-'}";`
    }
    const sum = computeStudentSummary(st.student_id)
    csv += `"${sum?.label || '-'}"\n`
  }

  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.setAttribute('download', `matrice_primaria_${currentClassName.value}_Q${selectedSemester.value}.csv`)
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)
  $q.notify({
    type: 'positive',
    message: t('primaryEval.exportCSVSuccess')
  })
}
</script>

<style scoped>
.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.primary-matrix-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.primary-matrix-table th,
.primary-matrix-table td {
  border: 1px solid #e2e8f0;
  padding: 8px 10px;
}

.primary-matrix-table thead tr {
  background-color: #f8fafc;
}

.sticky-col {
  position: sticky;
  left: 0;
  background-color: #ffffff;
  z-index: 2;
  box-shadow: 2px 0 4px rgba(0,0,0,0.05);
}

.primary-matrix-table thead .sticky-col {
  background-color: #f8fafc;
  z-index: 3;
}

.student-header {
  min-width: 200px;
  max-width: 260px;
}

.objective-header {
  min-width: 140px;
  max-width: 200px;
}

.summary-col {
  min-width: 130px;
  background-color: #f8fafc;
}

.student-row:hover {
  background-color: #f1f5f9;
}

.student-row:hover .sticky-col {
  background-color: #f1f5f9;
}

.level-dropdown-btn {
  min-width: 70px;
  border-radius: 8px;
}
</style>
