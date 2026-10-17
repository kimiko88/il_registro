<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">{{ t('classFormation.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('classFormation.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          outline
          color="primary"
          icon="cloud_upload"
          :label="t('classFormation.sidiImport')"
          class="rounded-lg q-px-md"
          @click="showImportDialog = true"
        />
        <q-btn
          color="primary"
          icon="auto_awesome"
          :label="t('classFormation.runSolver')"
          class="rounded-lg q-px-md shadow-soft"
          @click="showSolveDialog = true"
        />
        <q-btn
          v-if="currentDraft && !currentDraft.is_finalized"
          color="positive"
          icon="task_alt"
          :label="t('common.confirm')"
          class="rounded-lg q-px-md shadow-soft"
          :loading="finalizing"
          @click="finalizeCurrentDraft"
        />
      </div>
    </div>

    <!-- Normative Alert Banner -->
    <q-card class="bg-blue-50 border-blue-200 text-blue-9 q-mb-md rounded-xl p-3">
      <div class="row items-center no-wrap">
        <q-icon name="gavel" size="24px" class="q-mr-sm text-blue-7" />
        <div class="text-caption leading-relaxed">
          <strong>Vincolo Normativo D.P.R. 81/2009 (art. 5 c. 2):</strong> Le classi con presenza di alunni con disabilità certificata (L. 104/92) sono costituite, di norma, con <strong>non più di 20 alunni</strong>. Per le restanti classi il limite ordinario è compreso tra 18 e 27 alunni.
        </div>
      </div>
    </q-card>

    <!-- Draft Selection Bar -->
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center q-gutter-sm">
        <span class="text-weight-bold text-slate-700">{{ t('classFormation.newDraft') }}:</span>
        <q-select
          v-model="selectedDraftId"
          :options="draftList"
          option-value="id"
          option-label="title"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          style="min-width: 320px;"
          @update:model-value="loadDraftDetails"
        />
      </div>

      <div v-if="currentDraft" class="row items-center q-gutter-xs">
        <q-badge :color="currentDraft.is_finalized ? 'positive' : 'amber-9'" class="text-caption q-py-xs q-px-sm">
          {{ currentDraft.is_finalized ? t('classFormation.draftPublished') : t('classFormation.newDraft') }}
        </q-badge>
      </div>
    </div>


    <!-- KANBAN BOARD -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner color="primary" size="40px" />
      <div class="text-grey-6 q-mt-md">{{ t('common.loading') }}</div>
    </div>

    <div v-else-if="!currentDraft || !currentDraft.assignments?.classes?.length" class="text-center q-pa-xl bg-white rounded-xl shadow-soft border-slate-100">
      <q-icon name="dashboard_customize" size="64px" color="grey-4" />
      <div class="text-h6 text-slate-700 q-mt-md">{{ t('classFormation.unassigned') }}</div>
      <p class="text-slate-400">{{ t('classFormation.dragDropHelp') }}</p>
    </div>

    <div v-else class="row q-col-gutter-md no-wrap overflow-x-auto q-pb-lg">
      <div
        v-for="(col, colIdx) in currentDraft.assignments.classes"
        :key="col.class_name"
        class="col-12 col-md-4"
        style="min-width: 340px; max-width: 380px;"
      >
        <q-card class="glass-card shadow-soft border-slate-200 rounded-xl overflow-hidden flex flex-col h-full bg-slate-100">
          <!-- Column Header -->
          <div class="q-pa-md bg-white border-b border-slate-200">
            <div class="row items-center justify-between">
              <div class="text-h6 text-weight-bold text-slate-800">{{ col.class_name }}</div>
              <q-badge
                :color="isClassNormativeWarning(col) ? 'negative' : 'primary'"
                class="text-weight-bold text-caption q-px-sm q-py-xs"
              >
                {{ col.students?.length || 0 }} {{ t('classFormation.totalStudents') }}
              </q-badge>
            </div>

            <!-- Normative warning badge if violated -->
            <div v-if="isClassNormativeWarning(col)" class="text-negative text-caption text-weight-bold q-mt-xs row items-center">
              <q-icon name="warning" size="14px" class="q-mr-xs" />
              {{ t('classFormation.l104Warning') }}
            </div>

            <!-- Class Stats Row -->
            <div class="row justify-between text-caption text-grey-7 q-mt-sm q-pt-sm border-t border-slate-100">
              <span>{{ t('classFormation.males') }}: <strong>{{ col.males_count || 0 }}</strong> / {{ t('classFormation.females') }}: <strong>{{ col.females_count || 0 }}</strong></span>
              <span>{{ t('classFormation.avgGrade') }}: <strong>{{ (col.average_grade || 0).toFixed(1) }}</strong></span>
              <span v-if="col.l104_count > 0" class="text-amber-9 font-bold">L.104: <strong>{{ col.l104_count }}</strong></span>
              <span v-if="col.dsa_count > 0" class="text-blue-8">DSA: <strong>{{ col.dsa_count }}</strong></span>
            </div>
          </div>

          <!-- Student Cards List -->
          <div class="q-pa-sm flex-1 scroll" style="max-height: 65vh; overflow-y: auto;">
            <div
              v-for="student in col.students"
              :key="student.id"
              class="q-mb-sm p-3 bg-white rounded-lg shadow-xs border border-slate-200"
            >
              <div class="row items-center justify-between no-wrap">
                <div class="text-subtitle2 text-weight-bold text-slate-800 ellipsis">
                  {{ student.student_last_name }} {{ student.student_first_name }}
                </div>
                <div class="row q-gutter-xs">
                  <q-badge :color="student.gender === 'M' ? 'blue-6' : 'pink-6'" text-color="white">
                    {{ student.gender }}
                  </q-badge>
                  <q-badge color="grey-3" text-color="grey-9">
                    {{ student.middle_school_grade }}
                  </q-badge>
                </div>
              </div>

              <div class="row q-gutter-xs q-mt-xs">
                <q-badge v-if="student.has_disability_l104" color="orange-9" label="L. 104/92" />
                <q-badge v-if="student.has_dsa" color="indigo-7" label="DSA" />
                <q-badge v-if="student.second_language" outline color="grey-7" :label="student.second_language" />
              </div>

              <!-- Move Student Selector -->
              <div v-if="!currentDraft.is_finalized && currentDraft.assignments.classes.length > 1" class="text-right q-mt-xs">
                <q-btn flat dense size="xs" color="grey-7" :label="t('common.edit')" icon="swap_horiz">
                  <q-menu auto-close>
                    <q-list dense style="min-width: 100px">
                      <template v-for="(target, tIdx) in currentDraft.assignments.classes" :key="target.class_name">
                        <q-item
                          v-if="tIdx !== colIdx"
                          clickable
                          @click="moveStudent(student, colIdx, tIdx)"
                        >
                          <q-item-section>{{ target.class_name }}</q-item-section>
                        </q-item>
                      </template>
                    </q-list>
                  </q-menu>
                </q-btn>
              </div>
            </div>
          </div>
        </q-card>
      </div>
    </div>

    <!-- SIDI Import Dialog -->
    <q-dialog v-model="showImportDialog">
      <q-card style="width: 500px; max-width: 90vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">{{ t('classFormation.sidiImport') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <p class="text-caption text-grey-7">
            {{ t('classFormation.subtitle') }}
          </p>
          <q-file
            v-model="importFile"
            :label="t('classFormation.sidiImport')"
            outlined
            accept=".csv, .txt"
          >
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            :label="t('classFormation.sidiImport')"
            :loading="importing"
            :disable="!importFile"
            @click="handleImportSIDI"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Solve / Generate Dialog -->
    <q-dialog v-model="showSolveDialog">
      <q-card style="width: 480px; max-width: 90vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">{{ t('classFormation.runSolver') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="solveForm.title"
            :label="t('classFormation.generateDraftTitle')"
            outlined
            dense
          />
          <q-input
            v-model.number="solveForm.target_class_count"
            :label="t('classFormation.numClasses')"
            type="number"
            outlined
            dense
            min="1"
            max="12"
          />
          <q-input
            v-model.number="solveForm.max_l104_per_class"
            :label="t('classFormation.l104Students')"
            type="number"
            outlined
            dense
            min="1"
            max="3"
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            :label="t('classFormation.runSolver')"
            :loading="solving"
            @click="handleSolveFormation"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { enrollmentService } from '@/services/enrollmentService'

const $q = useQuasar()
const { t } = useI18n()


const loading = ref(false)
const draftList = ref([])
const selectedDraftId = ref(null)
const currentDraft = ref(null)

const showImportDialog = ref(false)
const importFile = ref(null)
const importing = ref(false)

const showSolveDialog = ref(false)
const solving = ref(false)
const finalizing = ref(false)

const solveForm = reactive({
  title: 'Bozza Classi Prime 2026/2027',
  target_class_count: 3,
  max_l104_per_class: 1,
  balance_gender: true,
  balance_grades: true
})

onMounted(fetchDrafts)

async function fetchDrafts() {
  loading.value = true
  try {
    const res = await enrollmentService.listDrafts()
    draftList.value = res.data || []
    if (draftList.value.length > 0 && !selectedDraftId.value) {
      selectedDraftId.value = draftList.value[0].id
      currentDraft.value = draftList.value[0]
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

async function loadDraftDetails(draftId) {
  if (!draftId) return
  loading.value = true
  try {
    const res = await enrollmentService.getDraft(draftId)
    currentDraft.value = res.data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function isClassNormativeWarning(col) {
  const count = col.students?.length || 0
  const l104 = col.l104_count || 0
  if (l104 > 0 && count > 20) return true
  if (count > 27) return true
  return false
}

function moveStudent(student, fromIdx, toIdx) {
  if (!currentDraft.value || !currentDraft.value.assignments?.classes) return
  const fromClass = currentDraft.value.assignments.classes[fromIdx]
  const toClass = currentDraft.value.assignments.classes[toIdx]

  fromClass.students = fromClass.students.filter(s => s.id !== student.id)
  toClass.students.push(student)

  // Recalculate metrics
  recalculateClassMetrics(fromClass)
  recalculateClassMetrics(toClass)

  // Save updated draft assignments
  enrollmentService.updateDraftAssignments(currentDraft.value.id, currentDraft.value.assignments)
    .catch(() => {
      $q.notify({ type: 'negative', message: 'Errore salvataggio spostamento' })
    })
}

function recalculateClassMetrics(cls) {
  cls.total_students = cls.students.length
  let males = 0
  let females = 0
  let l104 = 0
  let dsa = 0
  let totalGrade = 0

  for (const s of cls.students) {
    if (s.gender === 'M') males++
    else females++
    if (s.has_disability_l104) l104++
    if (s.has_dsa) dsa++
    totalGrade += (s.middle_school_grade || 7)
  }

  cls.males_count = males
  cls.females_count = females
  cls.l104_count = l104
  cls.dsa_count = dsa
  cls.average_grade = cls.total_students > 0 ? (totalGrade / cls.total_students) : 0
}

async function handleImportSIDI() {
  if (!importFile.value) return
  importing.value = true
  try {
    const fd = new FormData()
    fd.append('file', importFile.value)
    const res = await enrollmentService.importSIDI(fd)
    $q.notify({ type: 'positive', message: res.data?.message || t('classFormation.draftSaved') })
    showImportDialog.value = false
    importFile.value = null
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    importing.value = false
  }
}

async function handleSolveFormation() {
  solving.value = true
  try {
    const payload = {
      title: solveForm.title,
      academic_year: '2026/2027',
      parameters: {
        target_class_count: solveForm.target_class_count,
        max_l104_per_class: solveForm.max_l104_per_class,
        balance_gender: solveForm.balance_gender,
        balance_grades: solveForm.balance_grades
      }
    }
    const res = await enrollmentService.generateFormationDraft(payload)
    $q.notify({ type: 'positive', message: t('classFormation.draftSaved') })
    showSolveDialog.value = false
    await fetchDrafts()
    selectedDraftId.value = res.data.id
    currentDraft.value = res.data
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    solving.value = false
  }
}

async function finalizeCurrentDraft() {
  if (!currentDraft.value) return
  finalizing.value = true
  try {
    await enrollmentService.finalizeDraft(currentDraft.value.id)
    $q.notify({ type: 'positive', message: t('classFormation.draftPublished') })
    await loadDraftDetails(currentDraft.value.id)
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    finalizing.value = false
  }
}

</script>

<style scoped>
.opacity-30 {
  opacity: 0.3;
}
</style>
