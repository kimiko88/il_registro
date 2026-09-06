<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="fact_check" color="primary" class="q-mr-sm" />
          {{ t('rubricsPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('rubricsPage.subtitle') }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          color="primary"
          unelevated
          icon="add"
          :label="t('rubricsPage.newRubric')"
          class="rounded-lg q-px-md"
          no-caps
          @click="openCreateRubricDialog"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="rubricsStore.loading" @click="loadRubrics" />
      </div>
    </div>

    <!-- Tabs Container -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-slate-600 bg-slate-100 border-b border-slate-200"
        active-color="primary"
        indicator-color="primary"
        align="left"
        no-caps
      >
        <q-tab name="rubrics" icon="menu_book" :label="t('rubricsPage.title')" />
        <q-tab name="assessments" icon="history_edu" :label="t('nav.competencies')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-white">
        <!-- TAB 1: MY RUBRICS -->
        <q-tab-panel name="rubrics" class="q-pa-lg">
          <div v-if="rubricsStore.loading" class="text-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>

          <div v-else-if="rubricsStore.rubrics.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="format_list_bulleted" size="64px" class="q-mb-md opacity-40" />
            <div class="text-h6">Nessuna rubrica definita</div>
            <div class="text-caption">Crea la tua prima rubrica valutativa per le tue materie.</div>
          </div>

          <div v-else class="row q-col-gutter-md">
            <div v-for="rub in rubricsStore.rubrics" :key="rub.id" class="col-12 col-md-6 col-lg-4">
              <q-card bordered flat class="rounded-xl h-full flex column justify-between shadow-sm hover:shadow-md transition-shadow">
                <q-card-section>
                  <div class="row items-center justify-between q-mb-xs">
                    <q-chip size="xs" color="indigo-1" text-color="indigo-8" class="text-weight-bold">
                      {{ rub.subject_id || 'Materia' }}
                    </q-chip>
                    <q-badge color="primary">
                      {{ rub.criteria ? rub.criteria.length : 0 }} Criteri
                    </q-badge>
                  </div>

                  <div class="text-h6 text-weight-bold text-slate-800 q-mt-xs">{{ rub.title }}</div>
                  <div class="text-caption text-slate-600 q-mt-xs">{{ rub.description || 'Nessuna descrizione.' }}</div>
                </q-card-section>

                <q-card-actions align="between" class="bg-slate-50 border-t border-slate-100 q-px-md q-py-sm">
                  <q-btn flat dense color="negative" icon="delete" size="sm" @click="confirmDeleteRubric(rub.id)" />
                  <q-btn
                    color="positive"
                    unelevated
                    size="sm"
                    icon="assignment_turned_in"
                    label="Valuta Studente"
                    no-caps
                    @click="openAssessmentDialog(rub)"
                  />
                </q-card-actions>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB 2: ASSESSMENTS HISTORY -->
        <q-tab-panel name="assessments" class="q-pa-lg">
          <!-- Filters -->
          <div class="row q-col-gutter-md q-mb-md items-center">
            <div class="col-12 col-sm-6 col-md-4">
              <q-select
                v-model="historyClassId"
                :options="classOptions"
                label="Filtra per Classe"
                outlined dense
                emit-value map-options
                @update:model-value="loadAssessments"
              />
            </div>
          </div>

          <div v-if="rubricsStore.loading" class="text-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>

          <div v-else-if="rubricsStore.assessments.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="find_in_page" size="64px" class="q-mb-md opacity-40" />
            <div class="text-h6">Nessuna valutazione registrata</div>
            <div class="text-caption">Seleziona una classe o valuta uno studente usando una delle tue rubriche.</div>
          </div>

          <q-list v-else separator class="rounded-lg border-slate-200">
            <q-item v-for="a in rubricsStore.assessments" :key="a.id" class="q-py-md">
              <q-item-section avatar>
                <q-avatar color="green-1" text-color="positive" icon="stars" size="48px" />
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800">
                  {{ a.student_name || a.student_id }}
                </q-item-label>
                <q-item-label caption class="text-slate-600">
                  Rubrica: <strong class="text-slate-800">{{ a.rubric_title || a.rubric_id }}</strong> · Data: {{ formatDate(a.date) }}
                </q-item-label>
                <q-item-label caption class="text-slate-500" v-if="a.notes">
                  Note: {{ a.notes }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-chip color="positive" text-color="white" class="text-weight-bold text-subtitle2">
                  Punteggio: {{ a.total_score }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog 1: Create Rubric -->
    <q-dialog v-model="createRubricDialog">
      <q-card style="width: min(700px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden">
        <q-form @submit="saveRubric" greedy>
          <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
            <div class="text-h6 text-weight-bold">Nuova Rubrica Valutativa</div>
            <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
          </q-card-section>

          <q-card-section class="q-pa-md space-y-4 max-h-70vh overflow-y-auto">
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-sm-8">
                <q-input v-model="rubricForm.title" :label="t('common.title') || 'Titolo Rubrica *'" outlined dense :rules="[val => !!val || t('common.requiredField') || 'Titolo obbligatorio']" />
              </div>
              <div class="col-12 col-sm-4">
                <q-select
                  v-model="rubricForm.subject_id"
                  :options="subjectOptions"
                  option-value="value"
                  option-label="label"
                  emit-value
                  map-options
                  :label="t('common.subject') || 'Materia *'"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField') || 'Materia obbligatoria']"
                />
              </div>
            </div>

            <q-input v-model="rubricForm.description" :label="t('common.description') || 'Descrizione Rubrica'" outlined dense type="textarea" rows="2" />

            <!-- Dynamic Criteria Section -->
            <div class="border border-slate-200 rounded-xl q-pa-md bg-slate-50">
              <div class="row items-center justify-between q-mb-sm">
                <div class="text-subtitle1 text-weight-bold text-slate-800">Criteri di Valutazione</div>
                <q-btn color="primary" size="sm" icon="add" label="Aggiungi Criterio" no-caps @click="addCriterion" />
              </div>

              <div v-for="(crit, cIdx) in rubricForm.criteria" :key="cIdx" class="bg-white p-3 rounded-lg border border-slate-200 q-mb-sm">
                <div class="row items-center justify-between q-mb-xs">
                  <div class="text-weight-bold text-slate-700">Criterio {{ cIdx + 1 }}</div>
                  <q-btn flat dense icon="delete" color="negative" size="xs" @click="removeCriterion(cIdx)" v-if="rubricForm.criteria.length > 1" />
                </div>

                <div class="row q-col-gutter-xs q-mb-xs">
                  <div class="col-8">
                    <q-input v-model="crit.name" label="Nome Criterio *" dense outlined :rules="[val => !!val || t('common.requiredField') || 'Nome criterio obbligatorio']" />
                  </div>
                  <div class="col-4">
                    <q-input v-model.number="crit.max_score" type="number" label="Punti Max" dense outlined />
                  </div>
                </div>

                <!-- Levels for Criterion -->
                <div class="q-pl-sm border-l-2 border-primary mt-2 space-y-1">
                  <div class="row items-center justify-between text-caption text-slate-500">
                    <span>Livelli di prestazione</span>
                    <q-btn flat size="xs" color="primary" icon="add" label="Livello" @click="addLevel(cIdx)" />
                  </div>
                  <div v-for="(lvl, lIdx) in crit.levels" :key="lIdx" class="row q-col-gutter-xs items-center">
                    <div class="col-5">
                      <q-input v-model="lvl.label" label="Label (es. Avanzato)" dense outlined />
                    </div>
                    <div class="col-4">
                      <q-input v-model.number="lvl.score" type="number" label="Punti" dense outlined />
                    </div>
                    <div class="col-3 text-right">
                      <q-btn flat dense icon="close" color="grey" size="xs" @click="removeLevel(cIdx, lIdx)" v-if="crit.levels.length > 1" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
            <q-btn color="primary" type="submit" label="Crea Rubrica" :loading="savingRubric" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>

    <!-- Dialog 2: Assess Student -->
    <q-dialog v-model="assessmentDialog">
      <q-card style="width: min(650px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden">
        <q-form @submit="saveAssessment" greedy>
          <q-card-section class="bg-positive text-white row items-center justify-between q-py-md">
            <div>
              <div class="text-h6 text-weight-bold">Valutazione con Rubrica</div>
              <div class="text-caption">{{ targetRubric?.title }}</div>
            </div>
            <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
          </q-card-section>

          <q-card-section class="q-pa-md space-y-4 max-h-70vh overflow-y-auto">
            <!-- Class & Student Select -->
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-sm-6">
                <q-select
                  v-model="assessmentForm.class_id"
                  :options="classOptions"
                  :label="t('common.class') || 'Classe *'"
                  outlined dense
                  emit-value map-options
                  :rules="[val => !!val || t('common.requiredField') || 'Seleziona una classe']"
                  @update:model-value="onAssessmentClassChange"
                />
              </div>
              <div class="col-12 col-sm-6">
                <q-select
                  v-model="assessmentForm.student_id"
                  :options="assessmentStudentOptions"
                  :label="t('common.student') || 'Studente *'"
                  outlined dense
                  emit-value map-options
                  :rules="[val => !!val || t('common.requiredField') || 'Seleziona uno studente']"
                />
              </div>
            </div>

            <!-- Loading indicator -->
            <div v-if="loadingRubricDetails" class="text-center q-pa-lg">
              <q-spinner-dots color="positive" size="40px" />
            </div>

            <!-- Dynamic Criteria Assessment Toggles -->
            <template v-else>
              <div v-for="crit in targetRubric?.criteria" :key="crit.id" class="border border-slate-200 rounded-xl q-pa-md bg-slate-50">
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">{{ crit.name }}</div>
                <div class="text-caption text-slate-500 q-mb-sm" v-if="crit.description">{{ crit.description }}</div>

                <div v-if="!crit.levels || crit.levels.length === 0" class="text-caption text-amber-900 bg-amber-100 q-pa-sm rounded-lg">
                  ⚠️ Nessun livello di valutazione definito per questo criterio.
                </div>
                <div v-else class="row q-gutter-xs">
                  <q-btn
                    v-for="lvl in crit.levels" :key="lvl.id"
                    :unelevated="getSelectedLevelId(crit.id) === lvl.id"
                    :outline="getSelectedLevelId(crit.id) !== lvl.id"
                    :color="getLevelBtnColor(lvl.score, crit.max_score)"
                    no-caps size="sm" class="q-px-sm"
                    @click="selectCriterionLevel(crit.id, lvl.id, lvl.score)"
                  >
                    {{ lvl.label }} ({{ lvl.score }} pt)
                  </q-btn>
                </div>
              </div>
            </template>

            <!-- Total Score Summary Banner -->
            <div class="bg-green-50 border border-green-200 rounded-xl q-pa-md row items-center justify-between">
              <span class="text-weight-bold text-slate-800">Punteggio Totale Calcolato:</span>
              <span class="text-h4 text-weight-bold text-positive">{{ computedTotalScore }} pt</span>
            </div>

            <!-- Notes -->
            <q-input v-model="assessmentForm.notes" label="Note ed Osservazioni" outlined dense type="textarea" rows="2" />
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
            <q-btn color="positive" type="submit" label="Salva Valutazione" :loading="savingAssessment" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useRubricsStore } from '@/stores/rubrics'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import api from 'src/services/api'

const $q = useQuasar()
const { t } = useI18n()
const rubricsStore = useRubricsStore()
const classesStore = useClassesStore()
const gradesStore = useGradesStore()

const activeTab = ref('rubrics')
const createRubricDialog = ref(false)
const assessmentDialog = ref(false)
const savingRubric = ref(false)
const savingAssessment = ref(false)
const loadingRubricDetails = ref(false)

const historyClassId = ref(null)
const targetRubric = ref(null)
const assessmentStudents = ref([])
const subjectsList = ref([])

const rubricForm = reactive({
  title: '',
  subject_id: '',
  description: '',
  criteria: []
})

const assessmentForm = reactive({
  class_id: '',
  student_id: '',
  date: new Date().toISOString().substring(0, 10),
  scores: {}, // map criterionId -> { levelId, score }
  notes: ''
})

const classOptions = computed(() => {
  return classesStore.classes.map(c => ({
    label: c.name || `Classe ${c.id}`,
    value: c.id
  }))
})

const subjectOptions = computed(() => {
  const storeSubjects = gradesStore.subjects || []
  const list = storeSubjects.length > 0 ? storeSubjects : subjectsList.value
  return list.map(s => ({
    label: s.name || s.subject_name || s.code || s.id,
    value: s.subject_id || s.id || s.code
  }))
})

const assessmentStudentOptions = computed(() => {
  return assessmentStudents.value.map(s => ({
    label: `${s.first_name || s.name} ${s.last_name || ''}`,
    value: s.id
  }))
})

const computedTotalScore = computed(() => {
  return Object.values(assessmentForm.scores).reduce((acc, curr) => acc + (curr.score || 0), 0)
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''

onMounted(async () => {
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  try {
    const res = await api.get('/subjects')
    subjectsList.value = res.data || []
  } catch (e) {
    console.warn('Could not load subjects list for rubrics', e)
  }
  await loadRubrics()
})

async function loadRubrics() {
  await rubricsStore.fetchRubrics().catch(() => {})
}

async function loadAssessments() {
  if (historyClassId.value) {
    await rubricsStore.fetchClassAssessments(historyClassId.value).catch(() => {})
  }
}

function openCreateRubricDialog() {
  rubricForm.title = ''
  rubricForm.subject_id = subjectOptions.value[0]?.value || 'MAT'
  rubricForm.description = ''
  rubricForm.criteria = [
    {
      name: 'Comprensione del problema',
      max_score: 4,
      levels: [
        { label: 'Avanzato', score: 4 },
        { label: 'Intermedio', score: 3 },
        { label: 'Base', score: 2 }
      ]
    }
  ]
  createRubricDialog.value = true
}

function addCriterion() {
  rubricForm.criteria.push({
    name: '',
    max_score: 4,
    levels: [{ label: 'Avanzato', score: 4 }]
  })
}

function removeCriterion(idx) {
  rubricForm.criteria.splice(idx, 1)
}

function addLevel(cIdx) {
  rubricForm.criteria[cIdx].levels.push({ label: 'Livello', score: 1 })
}

function removeLevel(cIdx, lIdx) {
  rubricForm.criteria[cIdx].levels.splice(lIdx, 1)
}

async function saveRubric() {
  savingRubric.value = true
  try {
    await rubricsStore.createRubric({
      title: rubricForm.title,
      subject_id: rubricForm.subject_id,
      description: rubricForm.description,
      criteria: rubricForm.criteria
    })
    $q.notify({ type: 'positive', message: 'Rubrica creata con successo' })
    createRubricDialog.value = false
    await loadRubrics()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore creazione rubrica' })
  } finally {
    savingRubric.value = false
  }
}

async function confirmDeleteRubric(id) {
  $q.dialog({
    title: 'Elimina Rubrica',
    message: 'Sei sicuro di voler eliminare questa rubrica?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await rubricsStore.deleteRubric(id)
      $q.notify({ type: 'positive', message: 'Rubrica eliminata' })
      await loadRubrics()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore eliminazione' })
    }
  })
}

async function openAssessmentDialog(rubric) {
  targetRubric.value = null
  loadingRubricDetails.value = true
  assessmentForm.class_id = classOptions.value[0]?.value || ''
  assessmentForm.student_id = ''
  assessmentForm.scores = {}
  assessmentForm.notes = ''
  assessmentDialog.value = true

  try {
    targetRubric.value = await rubricsStore.getRubric(rubric.id)
  } catch {
    targetRubric.value = rubric
  } finally {
    loadingRubricDetails.value = false
  }

  await onAssessmentClassChange()
}

async function onAssessmentClassChange() {
  if (!assessmentForm.class_id) {
    assessmentStudents.value = []
    return
  }
  try {
    const res = await api.get('/users', { params: { role: 'student', class_id: assessmentForm.class_id } })
    assessmentStudents.value = res.data?.users || res.data || []
    if (assessmentStudents.value.length > 0) {
      assessmentForm.student_id = assessmentStudents.value[0].id
    }
  } catch {
    assessmentStudents.value = []
  }
}

function selectCriterionLevel(critId, levelId, score) {
  assessmentForm.scores[critId] = { levelId, score }
}

function getSelectedLevelId(critId) {
  return assessmentForm.scores[critId]?.levelId || null
}

function getLevelBtnColor(score, maxScore) {
  const ratio = maxScore > 0 ? score / maxScore : 1
  if (ratio >= 0.75) return 'positive'
  if (ratio >= 0.5) return 'warning'
  return 'negative'
}

async function saveAssessment() {
  savingAssessment.value = true
  try {
    const scoresArr = Object.entries(assessmentForm.scores).map(([critId, val]) => ({
      criterion_id: critId,
      level_id: val.levelId,
      score: val.score
    }))

    await rubricsStore.assessStudent(targetRubric.value.id, {
      student_id: assessmentForm.student_id,
      class_id: assessmentForm.class_id,
      date: assessmentForm.date,
      scores: scoresArr,
      notes: assessmentForm.notes
    })

    $q.notify({ type: 'positive', message: 'Valutazione registrata con successo' })
    assessmentDialog.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore salvataggio valutazione' })
  } finally {
    savingAssessment.value = false
  }
}
</script>

<style scoped>
.max-h-70vh {
  max-height: 70vh;
}
</style>
