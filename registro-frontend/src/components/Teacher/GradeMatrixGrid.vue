<template>
  <q-card
    class="rounded-xl shadow-xs border"
    :class="$q.dark.isActive ? 'bg-dark border-grey-8 text-white' : 'bg-white text-slate-800'"
  >
    <q-card-section class="row items-center justify-between q-pb-sm">
      <div>
        <div class="text-h6 text-weight-bold row items-center" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="grid_on" color="primary" class="q-mr-sm" />
          {{ t('gradesPage.matrixViewTitle') || 'Inserimento Rapido Voti in Griglia (Matrix View)' }}
        </div>
        <div class="text-caption" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('gradesPage.matrixViewKbdHint') || 'Usa TAB, INVIO o le FRECCE per spostarti velocemente tra gli studenti' }}
        </div>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          color="positive"
          icon="save"
          :label="t('gradesPage.saveAllGrades') || 'Salva Tutti i Voti'"
          unelevated
          class="rounded-lg text-weight-bold"
          :loading="saving"
          @click="saveAllGrades"
        />
      </div>
    </q-card-section>

    <q-separator />

    <q-card-section class="q-pa-none">
      <div class="table-responsive">
        <table class="matrix-table full-width">
          <thead>
            <tr
              class="text-left"
              :class="$q.dark.isActive ? 'bg-grey-9 text-grey-3' : 'bg-slate-100 text-slate-700'"
            >
              <th class="q-pa-sm" style="width: 40px">#</th>
              <th class="q-pa-sm">{{ t('competenciesPage.student') || 'Alunno' }}</th>
              <th class="q-pa-sm" :style="{ width: isReligionSubject ? '180px' : '150px' }">{{ isReligionSubject ? 'Giudizio IRC' : (t('classRegister.tableHeaderGrade') || 'Voto (1-10)') }}</th>
              <th class="q-pa-sm" style="width: 220px">{{ t('gradesPage.besDsaMeasures') || 'Misure BES / DSA' }}</th>
              <th class="q-pa-sm">{{ t('classRegister.tableHeaderGradeNotes') || 'Note / Descrizione' }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(student, idx) in students"
              :key="student.id"
              class="border-b"
              :class="$q.dark.isActive ? 'hover:bg-grey-8 border-grey-8' : 'hover:bg-slate-50 border-slate-200'"
            >
              <td class="q-pa-sm text-weight-bold" :class="$q.dark.isActive ? 'text-grey-5' : 'text-grey-6'">{{ idx + 1 }}</td>
              <td class="q-pa-sm text-weight-bold" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
                {{ student.last_name }} {{ student.first_name }}
              </td>
              <td class="q-pa-sm">
                <!-- Se studente esonerato da religione -->
                <div v-if="isReligionSubject && isStudentExempt(student)" class="row items-center q-gutter-xs">
                  <q-badge
                    :color="student.religion_choice === 'attivita_alternativa' ? 'purple-7' : 'orange-8'"
                    text-color="white"
                    class="text-weight-bold q-py-xs q-px-sm rounded-md"
                  >
                    <q-icon :name="student.religion_choice === 'attivita_alternativa' ? 'swap_horiz' : 'block'" size="14px" class="q-mr-xs" />
                    {{ student.religion_choice === 'attivita_alternativa' ? 'Attività Alternativa' : 'Non Avvalente' }}
                  </q-badge>
                  <q-tooltip>Questo studente non si avvale dell'IRC e non può ricevere valutazioni</q-tooltip>
                </div>
                <!-- Se materia religione: select con i 6 giudizi -->
                <q-select
                  v-else-if="isReligionSubject"
                  v-model="student.religion_judgment"
                  :options="RELIGION_JUDGMENT_OPTIONS"
                  dense
                  outlined
                  placeholder="Seleziona..."
                  class="text-weight-bold"
                  :bg-color="getGradeColor(student.religion_judgment)"
                  hide-bottom-space
                  :ref="el => inputRefs[idx] = el"
                  @keydown.enter.prevent="focusNext(idx)"
                  @keydown.down.prevent="focusNext(idx)"
                  @keydown.up.prevent="focusPrev(idx)"
                />
                <!-- Altrimenti: input numerico standard -->
                <q-input
                  v-else
                  v-model.number="student.grade_value"
                  type="number"
                  step="0.25"
                  min="1"
                  max="10"
                  dense
                  outlined
                  class="text-weight-bold text-center"
                  :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-white'"
                  :input-class="isGradeInvalid(student.grade_value) ? 'text-weight-bold text-negative text-center' : 'text-weight-bold text-primary text-center'"
                  :rules="[
                    val => val === null || val === undefined || val === '' ||
                      (!isNaN(val) && Number(val) >= 1 && Number(val) <= 10) ||
                      (t('gradesPage.gradeRuleError') || '1-10')
                  ]"
                  lazy-rules
                  hide-bottom-space
                  :ref="el => inputRefs[idx] = el"
                  @keydown.enter.prevent="focusNext(idx)"
                  @keydown.down.prevent="focusNext(idx)"
                  @keydown.up.prevent="focusPrev(idx)"
                  @keydown.right="handleGradeArrowRight(idx)"
                >
                  <template v-if="isGradeInvalid(student.grade_value)" #append>
                    <q-icon name="warning" color="negative" size="xs">
                      <q-tooltip class="bg-negative">{{ t('gradesPage.gradeRuleError') || 'Il voto deve essere compreso tra 1 e 10' }}</q-tooltip>
                    </q-icon>
                  </template>
                </q-input>
              </td>
              <td class="q-pa-sm">
                <CompensativeMeasuresSelector v-model="student.compensative_measures" />
              </td>
              <td class="q-pa-sm">
                <q-input
                  v-model="student.notes"
                  dense
                  outlined
                  :placeholder="t('common.optionalNotes') || 'Note facoltative'"
                  :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-white'"
                  :ref="el => noteRefs[idx] = el"
                  @keydown.enter.prevent="focusNextNote(idx)"
                  @keydown.down.prevent="focusNextNote(idx)"
                  @keydown.up.prevent="focusPrevNote(idx)"
                  @keydown.left="handleNoteArrowLeft(idx, $event)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'
import { useOfflineSync } from '@/composables/useOfflineSync'
import { useGradesStore } from '@/stores/grades'
import { RELIGION_JUDGMENT_OPTIONS, RELIGION_JUDGMENT_MAP, getGradeColor } from '@/utils/gradeUtils'
import CompensativeMeasuresSelector from './CompensativeMeasuresSelector.vue'

const { t } = useI18n()
const props = defineProps({
  studentsList: {
    type: Array,
    default: () => []
  },
  subjectId: {
    type: String,
    required: true
  },
  classId: {
    type: String,
    required: true
  },
  isReligion: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['saved'])
const $q = useQuasar()
const { executeWithOfflineQueue } = useOfflineSync()

const saving = ref(false)
const inputRefs = ref([])
const noteRefs = ref([])
const students = ref([])

function isGradeInvalid(val) {
  if (val === null || val === undefined || val === '') return false
  const num = Number(val)
  return isNaN(num) || num < 1 || num > 10
}

const gradesStore = useGradesStore()

const isReligionSubject = computed(() => {
  if (props.isReligion) return true
  if (gradesStore.grades?.is_religion_subject) return true
  const currentSub = (gradesStore.subjects || []).find(s =>
    String(s.id) === String(props.subjectId) ||
    String(s.subject_id) === String(props.subjectId)
  )
  if (currentSub?.is_religion) return true
  if (currentSub?.subject_name && currentSub.subject_name.toLowerCase().includes('religione')) return true
  return false
})

function isStudentExempt(student) {
  return student?.religion_choice === 'non_avvalente' || student?.religion_choice === 'attivita_alternativa'
}

watch(() => props.studentsList, (val) => {
  students.value = (val || []).map(s => ({
    id: s.id,
    first_name: s.first_name,
    last_name: s.last_name,
    religion_choice: s.religion_choice,
    grade_value: null,
    religion_judgment: null,
    compensative_measures: [],
    notes: ''
  }))
}, { immediate: true })

function focusNext(idx) {
  if (idx < students.value.length - 1 && inputRefs.value[idx + 1]) {
    inputRefs.value[idx + 1].focus()
  }
}

function focusPrev(idx) {
  if (idx > 0 && inputRefs.value[idx - 1]) {
    inputRefs.value[idx - 1].focus()
  }
}

function handleGradeArrowRight(idx) {
  if (noteRefs.value[idx]) {
    noteRefs.value[idx].focus()
  }
}

function handleNoteArrowLeft(idx, event) {
  const target = event?.target
  if (!target || target.selectionStart === 0) {
    if (inputRefs.value[idx]) {
      inputRefs.value[idx].focus()
    }
  }
}

function focusNextNote(idx) {
  if (idx < students.value.length - 1 && noteRefs.value[idx + 1]) {
    noteRefs.value[idx + 1].focus()
  }
}

function focusPrevNote(idx) {
  if (idx > 0 && noteRefs.value[idx - 1]) {
    noteRefs.value[idx - 1].focus()
  }
}

async function saveAllGrades() {
  if (isReligionSubject.value) {
    const gradesToSave = students.value.filter(s =>
      s.religion_judgment &&
      !isStudentExempt(s)
    )
    if (gradesToSave.length === 0) {
      $q.notify({ type: 'warning', message: 'Inserisci almeno un giudizio in griglia per studenti avvalenti' })
      return
    }

    saving.value = true
    try {
      const payload = {
        subject_id: props.subjectId,
        class_id: props.classId,
        semester: 1,
        date: new Date().toISOString().split('T')[0],
        grades: gradesToSave.map(s => ({
          student_id: s.id,
          grade_value: RELIGION_JUDGMENT_MAP[s.religion_judgment] ?? 6,
          grade_type: 'judgment',
          description: s.religion_judgment + (s.notes ? ' - ' + s.notes : ''),
          compensative_measures: s.compensative_measures
        }))
      }
      const res = await executeWithOfflineQueue(
        { url: '/grades/bulk', method: 'post', data: payload },
        { title: 'Salvataggio giudizi IRC' }
      )
      if (!res?.offline && !res?.enqueued) {
        $q.notify({ type: 'positive', message: 'Giudizi IRC salvati con successo' })
      }
      emit('saved')
    } catch {
      $q.notify({ type: 'negative', message: t('common.error') })
    } finally {
      saving.value = false
    }
    return
  }

  const gradesToSave = students.value.filter(s => s.grade_value !== null && s.grade_value !== '' && !isNaN(s.grade_value))
  if (gradesToSave.length === 0) {
    $q.notify({ type: 'warning', message: t('gradesPage.insertAtLeastOneGrade') || 'Inserisci almeno un voto in griglia' })
    return
  }

  // Pre-save range check (1 - 10)
  const hasInvalid = gradesToSave.some(s => isGradeInvalid(s.grade_value))
  if (hasInvalid) {
    $q.notify({
      type: 'warning',
      message: t('gradesPage.gradeRuleError') || 'Tutti i voti inseriti devono essere compresi tra 1 e 10'
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      subject_id: props.subjectId,
      class_id: props.classId,
      semester: 1,
      date: new Date().toISOString().split('T')[0],
      grades: gradesToSave.map(s => ({
        student_id: s.id,
        grade_value: Number(s.grade_value),
        description: s.notes || (t('gradesPage.matrixGradeDesc') || 'Valutazione in griglia'),
        compensative_measures: s.compensative_measures
      }))
    }
    const res = await executeWithOfflineQueue(
      { url: '/grades/bulk', method: 'post', data: payload },
      { title: t('gradesPage.saveAllGrades') || 'Salvataggio voti griglia' }
    )
    if (!res?.offline && !res?.enqueued) {
      $q.notify({ type: 'positive', message: t('common.success') })
    }
    emit('saved')
  } catch {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.matrix-table {
  border-collapse: collapse;
}
.matrix-table th, .matrix-table td {
  border: 1px solid rgba(226, 232, 240, 0.8);
}
body.body--dark .matrix-table th,
body.body--dark .matrix-table td {
  border: 1px solid rgba(71, 85, 105, 0.6);
}
kbd {
  font-size: 11px;
  font-family: monospace;
}
</style>
