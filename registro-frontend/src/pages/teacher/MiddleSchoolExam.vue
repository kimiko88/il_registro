<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">{{ t('middleSchoolExam.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('middleSchoolExam.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-select
          v-model="selectedClassId"
          :options="classesStore.classes"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          :label="t('middleSchoolExam.selectClass')"
          style="min-width: 220px;"
          @update:model-value="loadExamData"
        />

        <q-btn
          v-if="exam"
          color="primary"
          icon="refresh"
          flat
          round
          @click="loadExamData"
        >
          <q-tooltip>{{ t('common.refresh') }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Exam Status Bar -->
    <q-card v-if="exam" class="glass-card shadow-soft border-slate-100 rounded-xl q-mb-md p-3">
      <div class="row items-center justify-between">
        <div class="row items-center q-gutter-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">
            {{ t('middleSchoolExam.subcommission') }} n. {{ exam.subcommission_number }} — {{ currentClassName }}
          </div>
          <div class="text-caption text-grey-7">
            Presidente: <strong>{{ exam.president_name }}</strong>
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <span class="text-caption text-grey-7">{{ t('middleSchoolExam.status') }}:</span>
          <q-badge :color="getStatusColor(exam.status)" class="text-caption q-py-xs q-px-sm">
            {{ getStatusLabel(exam.status) }}
          </q-badge>
        </div>
      </div>
    </q-card>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      class="text-grey-7 bg-white rounded-t-lg shadow-sm"
      active-color="primary"
      indicator-color="primary"
      align="left"
    >
      <q-tab name="admission" icon="how_to_reg" :label="t('middleSchoolExam.admissionGrade')" />
      <q-tab name="tests" icon="edit_note" :label="t('middleSchoolExam.written1') + ' & ' + t('middleSchoolExam.interview')" />
      <q-tab name="final" icon="grade" :label="t('middleSchoolExam.finalGrade')" />
      <q-tab name="diplomas" icon="workspace_premium" :label="t('middleSchoolExam.printDiploma')" />
    </q-tabs>

    <q-separator />

    <!-- TAB PANELS -->
    <q-tab-panels v-model="activeTab" animated class="bg-transparent q-mt-md">
      <!-- PANEL 1: ADMISSION -->
      <q-tab-panel name="admission" class="q-pa-none">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl overflow-hidden">
          <q-table
            :rows="candidates"
            :columns="admissionColumns"
            row-key="id"
            flat
            :loading="loading"
            class="bg-transparent"
            :pagination="{ rowsPerPage: 20 }"
          >
            <template v-slot:body-cell-is_admitted="props">
              <q-td :props="props">
                <q-chip :color="props.value ? 'green-1' : 'red-1'" :text-color="props.value ? 'green-9' : 'red-9'" dense>
                  {{ props.value ? t('middleSchoolExam.passed') : t('middleSchoolExam.notPassed') }}
                </q-chip>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn flat round dense color="primary" icon="edit" @click="openAdmissionDialog(props.row)">
                  <q-tooltip>{{ t('common.edit') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 2: EXAM TESTS -->
      <q-tab-panel name="tests" class="q-pa-none">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl overflow-hidden">
          <q-table
            :rows="candidates"
            :columns="testsColumns"
            row-key="id"
            flat
            :loading="loading"
            class="bg-transparent"
            :pagination="{ rowsPerPage: 20 }"
          >
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn color="primary" dense :label="t('common.edit')" icon="edit" class="q-px-sm text-caption" @click="openEvaluationDialog(props.row)" />
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>


      <!-- PANEL 3: FINAL OUTCOME -->
      <q-tab-panel name="final" class="q-pa-none">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl overflow-hidden">
          <q-table
            :rows="candidates"
            :columns="finalColumns"
            row-key="id"
            flat
            :loading="loading"
            class="bg-transparent"
            :pagination="{ rowsPerPage: 20 }"
          >
            <template v-slot:body-cell-final_grade="props">
              <q-td :props="props">
                <span class="text-weight-bold text-subtitle1" :class="props.value >= 6 ? 'text-primary' : 'text-negative'">
                  {{ props.value > 0 ? props.value + '/10' : '—' }}
                </span>
              </q-td>
            </template>

            <template v-slot:body-cell-has_honors="props">
              <q-td :props="props">
                <q-badge v-if="props.value" color="amber-9" :label="t('middleSchoolExam.lodeGranted')" class="text-weight-bold" />
                <span v-else class="text-grey-4">—</span>
              </q-td>
            </template>

            <template v-slot:body-cell-outcome="props">
              <q-td :props="props">
                <q-chip
                  v-if="props.value"
                  :color="props.value === 'licenziato' ? 'green-1' : 'red-1'"
                  :text-color="props.value === 'licenziato' ? 'green-9' : 'red-9'"
                  dense
                  class="text-weight-bold"
                >
                  {{ props.value === 'licenziato' ? t('middleSchoolExam.passed') : t('middleSchoolExam.notPassed') }}
                </q-chip>
                <span v-else class="text-grey-4">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 4: DIPLOMAS -->
      <q-tab-panel name="diplomas" class="q-pa-none">
        <div class="row q-col-gutter-md">
          <div v-for="cand in licenziatiCandidates" :key="cand.id" class="col-12 col-md-6">
            <q-card class="glass-card shadow-soft border-slate-100 rounded-xl p-3 flex justify-between items-center">
              <div>
                <div class="text-h6 text-weight-bold text-slate-800">{{ cand.student_name }}</div>
                <div class="text-caption text-grey-7">
                  {{ t('middleSchoolExam.finalGrade') }}: <strong>{{ cand.final_grade }}/10</strong>
                  <span v-if="cand.has_honors" class="text-amber-9 font-bold"> {{ t('middleSchoolExam.lodeGranted') }}</span>
                </div>
              </div>
              <div>
                <q-btn
                  color="primary"
                  icon="download"
                  :label="t('middleSchoolExam.printDiploma')"
                  class="rounded-lg shadow-sm"
                  @click="downloadDiploma(cand)"
                />
              </div>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Admission Edit Dialog -->
    <q-dialog v-model="showAdmissionDialog">
      <q-card style="width: 450px; max-width: 90vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">{{ t('middleSchoolExam.admissionGrade') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <div class="text-subtitle1 font-bold">{{ activeCandidate?.student_name }}</div>
          <q-input
            v-model.number="admissionForm.grade"
            :label="t('middleSchoolExam.admissionGrade')"
            type="number"
            min="6"
            max="10"
            outlined
            dense
          />
          <q-input
            v-model="admissionForm.judgment"
            :label="t('religionAlternative.judgment')"
            type="textarea"
            outlined
            rows="3"
          />
          <q-checkbox v-model="admissionForm.admitted" :label="t('middleSchoolExam.passed')" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" @click="saveAdmission" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Evaluation Dialog -->
    <q-dialog v-model="showEvalDialog">
      <q-card style="width: 520px; max-width: 95vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">{{ t('middleSchoolExam.title') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-sm">
          <div class="text-subtitle1 font-bold">{{ activeCandidate?.student_name }}</div>
          <div class="text-caption text-grey-7">{{ t('middleSchoolExam.admissionGrade') }}: <strong>{{ activeCandidate?.admission_grade }}/10</strong></div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input v-model.number="evalForm.grade_italian" :label="t('middleSchoolExam.written1')" type="number" step="0.5" min="4" max="10" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model.number="evalForm.grade_math" :label="t('middleSchoolExam.written2')" type="number" step="0.5" min="4" max="10" outlined dense />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input v-model.number="evalForm.grade_english" :label="t('middleSchoolExam.written3')" type="number" step="0.5" min="4" max="10" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model.number="evalForm.grade_second_lang" :label="t('middleSchoolExam.written3')" type="number" step="0.5" min="4" max="10" outlined dense />
            </div>
          </div>

          <q-input v-model.number="evalForm.grade_interview" :label="t('middleSchoolExam.interview')" type="number" step="0.5" min="4" max="10" outlined dense />

          <q-checkbox v-model="evalForm.proposed_honors" :label="t('middleSchoolExam.proposeLode')" />

          <q-input v-model="evalForm.notes" :label="t('common.notes')" outlined dense />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" @click="saveEvaluation" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useClassesStore } from '@/stores/classes'
import { middleSchoolExamService } from '@/services/middleSchoolExamService'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()


const loading = ref(false)
const selectedClassId = ref(null)
const exam = ref(null)
const candidates = ref([])
const activeTab = ref('admission')

const showAdmissionDialog = ref(false)
const showEvalDialog = ref(false)
const activeCandidate = ref(null)

const admissionForm = reactive({
  grade: 8,
  judgment: '',
  admitted: true
})

const evalForm = reactive({
  grade_italian: 8,
  grade_math: 8,
  grade_english: 8,
  grade_second_lang: 8,
  grade_interview: 8,
  proposed_honors: false,
  notes: ''
})

const currentClassName = computed(() => {
  const c = classesStore.classes?.find(cl => cl.id === selectedClassId.value)
  return c ? c.name : '3A'
})

const licenziatiCandidates = computed(() => {
  return candidates.value.filter(c => c.outcome === 'licenziato')
})

const admissionColumns = computed(() => [
  { name: 'student_name', label: t('middleSchoolExam.candidate'), field: 'student_name', align: 'left', sortable: true },
  { name: 'admission_grade', label: t('middleSchoolExam.admissionGrade'), field: 'admission_grade', align: 'center', sortable: true },
  { name: 'admission_judgment', label: t('religionAlternative.judgment'), field: 'admission_judgment', align: 'left' },
  { name: 'is_admitted', label: t('middleSchoolExam.status'), field: 'is_admitted', align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

const testsColumns = computed(() => [
  { name: 'student_name', label: t('middleSchoolExam.candidate'), field: 'student_name', align: 'left' },
  { name: 'grade_italian', label: t('middleSchoolExam.written1'), field: 'grade_italian', align: 'center' },
  { name: 'grade_math', label: t('middleSchoolExam.written2'), field: 'grade_math', align: 'center' },
  { name: 'grade_english', label: t('middleSchoolExam.written3'), field: 'grade_english', align: 'center' },
  { name: 'grade_second_lang', label: t('middleSchoolExam.written3'), field: 'grade_second_lang', align: 'center' },
  { name: 'grade_interview', label: t('middleSchoolExam.interview'), field: 'grade_interview', align: 'center' },
  { name: 'exam_mean', label: t('middleSchoolExam.finalGrade'), field: 'exam_mean', align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

const finalColumns = computed(() => [
  { name: 'student_name', label: t('middleSchoolExam.candidate'), field: 'student_name', align: 'left' },
  { name: 'admission_grade', label: t('middleSchoolExam.admissionGrade'), field: 'admission_grade', align: 'center' },
  { name: 'exam_mean', label: t('middleSchoolExam.finalGrade'), field: 'exam_mean', align: 'center' },
  { name: 'final_grade', label: t('middleSchoolExam.finalGrade'), field: 'final_grade', align: 'center', sortable: true },
  { name: 'has_honors', label: t('middleSchoolExam.lodeGranted'), field: 'has_honors', align: 'center' },
  { name: 'outcome', label: t('middleSchoolExam.status'), field: 'outcome', align: 'center' }
])

onMounted(async () => {
  await classesStore.fetchClasses()
  if (classesStore.classes?.length > 0 && !selectedClassId.value) {
    selectedClassId.value = classesStore.classes[0].id
    await loadExamData()
  }
})

async function loadExamData() {
  if (!selectedClassId.value) return
  loading.value = true
  try {
    const examRes = await middleSchoolExamService.getOrCreateExam(selectedClassId.value)
    exam.value = examRes.data
    const candRes = await middleSchoolExamService.listCandidates(exam.value.id)
    candidates.value = candRes.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function getStatusColor(status) {
  switch (status) {
    case 'closed': return 'grey-7'
    case 'deliberated': return 'positive'
    case 'in_progress': return 'primary'
    default: return 'amber-9'
  }
}

function getStatusLabel(status) {
  switch (status) {
    case 'closed': return t('common.close')
    case 'deliberated': return t('middleSchoolExam.unanimousDeliberation')
    case 'in_progress': return t('middleSchoolExam.title')
    default: return t('middleSchoolExam.admissionGrade')
  }
}

function openAdmissionDialog(row) {
  activeCandidate.value = row
  admissionForm.grade = row.admission_grade || 8
  admissionForm.judgment = row.admission_judgment || ''
  admissionForm.admitted = row.is_admitted !== false
  showAdmissionDialog.value = true
}

async function saveAdmission() {
  if (!exam.value || !activeCandidate.value) return
  try {
    await middleSchoolExamService.saveAdmission(exam.value.id, activeCandidate.value.student_id, admissionForm)
    $q.notify({ type: 'positive', message: t('middleSchoolExam.saveGradesSuccess') })
    showAdmissionDialog.value = false
    await loadExamData()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  }
}

function openEvaluationDialog(row) {
  activeCandidate.value = row
  evalForm.grade_italian = row.grade_italian || 8
  evalForm.grade_math = row.grade_math || 8
  evalForm.grade_english = row.grade_english || 8
  evalForm.grade_second_lang = row.grade_second_lang || 8
  evalForm.grade_interview = row.grade_interview || 8
  evalForm.proposed_honors = row.has_honors || false
  evalForm.notes = row.notes || ''
  showEvalDialog.value = true
}

async function saveEvaluation() {
  if (!exam.value || !activeCandidate.value) return
  try {
    await middleSchoolExamService.evaluateCandidate(exam.value.id, activeCandidate.value.student_id, {
      grades: {
        grade_italian: evalForm.grade_italian,
        grade_math: evalForm.grade_math,
        grade_english: evalForm.grade_english,
        grade_second_lang: evalForm.grade_second_lang,
        grade_interview: evalForm.grade_interview,
        proposed_honors: evalForm.proposed_honors
      },
      notes: evalForm.notes
    })
    $q.notify({ type: 'positive', message: t('middleSchoolExam.saveGradesSuccess') })
    showEvalDialog.value = false
    await loadExamData()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  }
}

async function downloadDiploma(cand) {
  try {
    const res = await middleSchoolExamService.generateDiploma(cand.id)
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `diploma_${cand.student_id}.txt`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    $q.notify({ type: 'positive', message: t('middleSchoolExam.printDiploma') })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  }
}

</script>
