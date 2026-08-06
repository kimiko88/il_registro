<template>
  <q-card class="rounded-xl shadow-xs border bg-white">
    <q-card-section class="row items-center justify-between q-pb-sm">
      <div>
        <div class="text-h6 text-weight-bold text-slate-800 row items-center">
          <q-icon name="grid_on" color="primary" class="q-mr-sm" />
          Inserimento Rapido Voti in Griglia (Matrix View)
        </div>
        <div class="text-caption text-slate-500">
          Usa <kbd class="bg-slate-200 q-px-xs rounded">TAB</kbd>, <kbd class="bg-slate-200 q-px-xs rounded">INVIO</kbd> o le <kbd class="bg-slate-200 q-px-xs rounded">FRECCE</kbd> per spostarti velocemente tra gli studenti
        </div>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          color="positive"
          icon="save"
          label="Salva Tutti i Voti"
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
            <tr class="bg-slate-100 text-slate-700 text-left">
              <th class="q-pa-sm" style="width: 40px">#</th>
              <th class="q-pa-sm">Alunno</th>
              <th class="q-pa-sm" style="width: 140px">Voto (1-10)</th>
              <th class="q-pa-sm" style="width: 220px">Misure BES / DSA</th>
              <th class="q-pa-sm">Note / Descrizione</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(student, idx) in students" :key="student.id" class="border-b hover:bg-slate-50">
              <td class="q-pa-sm text-grey-6 text-weight-bold">{{ idx + 1 }}</td>
              <td class="q-pa-sm text-weight-bold text-slate-800">
                {{ student.last_name }} {{ student.first_name }}
              </td>
              <td class="q-pa-sm">
                <q-input
                  v-model.number="student.grade_value"
                  type="number"
                  step="0.25"
                  min="1" max="10"
                  dense outlined
                  class="bg-white text-weight-bold text-center"
                  input-class="text-weight-bold text-primary text-center"
                  :ref="el => inputRefs[idx] = el"
                  @keydown.enter.prevent="focusNext(idx)"
                  @keydown.down.prevent="focusNext(idx)"
                  @keydown.up.prevent="focusPrev(idx)"
                />
              </td>
              <td class="q-pa-sm">
                <CompensativeMeasuresSelector v-model="student.compensative_measures" />
              </td>
              <td class="q-pa-sm">
                <q-input
                  v-model="student.notes"
                  dense outlined
                  placeholder="Note facoltative"
                  class="bg-white"
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
import { ref, watch } from 'vue'
import { useQuasar } from 'quasar'
import api from '@/services/api'
import CompensativeMeasuresSelector from './CompensativeMeasuresSelector.vue'

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
  }
})

const emit = defineEmits(['saved'])
const $q = useQuasar()

const saving = ref(false)
const inputRefs = ref([])
const students = ref([])

watch(() => props.studentsList, (val) => {
  students.value = (val || []).map(s => ({
    id: s.id,
    first_name: s.first_name,
    last_name: s.last_name,
    grade_value: null,
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

async function saveAllGrades() {
  const gradesToSave = students.value.filter(s => s.grade_value !== null && s.grade_value !== '' && !isNaN(s.grade_value))
  if (gradesToSave.length === 0) {
    $q.notify({ type: 'warning', message: 'Inserisci almeno un voto in griglia' })
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
        grade_value: s.grade_value,
        description: s.notes || 'Valutazione in griglia',
        compensative_measures: s.compensative_measures
      }))
    }
    await api.post('/grades/bulk', payload)
    $q.notify({ type: 'positive', message: `${gradesToSave.length} voti salvati con successo!` })
    emit('saved')
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio dei voti in griglia' })
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
  border: 1px solid #e2e8f0;
}
kbd {
  font-size: 11px;
  font-family: monospace;
}
</style>
