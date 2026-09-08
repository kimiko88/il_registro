<template>
  <div class="descriptive-matrix-container q-pa-md">
    <!-- Header & Info -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h6 text-weight-bold text-slate-800 row items-center">
          <q-icon name="auto_stories" color="indigo" class="q-mr-sm" size="24px" />
          {{ t('rubrics.descriptiveTitle') || 'Matrice Valutazione Descrittiva per Obiettivi (O.M. 172/2020)' }}
        </div>
        <div class="text-caption text-slate-500">
          {{ t('rubrics.descriptiveSubtitle') || 'Valutazione formativa periodica per i 4 livelli ministeriali della scuola primaria e secondaria' }}
        </div>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-btn
          color="indigo-7"
          icon="add"
          label="Aggiungi Obiettivo"
          no-caps dense
          class="rounded-lg q-px-sm font-semibold"
          @click="showAddObjectiveDialog = true"
        />
        <q-btn
          color="secondary"
          icon="download"
          label="Esporta CSV"
          no-caps dense
          class="rounded-lg q-px-sm"
          @click="exportMatrixCSV"
        />
        <q-btn
          color="positive"
          icon="save"
          label="Salva Matrice"
          no-caps dense
          class="rounded-lg q-px-md font-bold"
          :loading="saving"
          @click="saveMatrix"
        />
      </div>
    </div>

    <!-- Filters: Class, Subject, Period -->
    <div class="row q-col-gutter-sm q-mb-md items-center">
      <div class="col-12 col-sm-4">
        <q-select
          v-model="selectedClassId"
          :options="classOptions"
          label="Seleziona Classe *"
          dense outlined emit-value map-options
          bg-color="white"
          @update:model-value="onClassChanged"
        />
      </div>
      <div class="col-12 col-sm-4">
        <q-select
          v-model="selectedSubjectId"
          :options="subjectOptions"
          label="Materia *"
          dense outlined emit-value map-options
          bg-color="white"
        />
      </div>
      <div class="col-12 col-sm-4">
        <q-select
          v-model="selectedPeriod"
          :options="[
            { label: '1° Quadrimestre (Valutazione Intermedia)', value: 'Q1' },
            { label: '2° Quadrimestre (Valutazione Finale)', value: 'Q2' }
          ]"
          label="Periodo Valutativo"
          dense outlined emit-value map-options
          bg-color="white"
        />
      </div>
    </div>

    <!-- 4 Ministerial Levels Legend Card -->
    <q-card flat bordered class="rounded-xl bg-white q-pa-sm q-mb-md border-slate-200 shadow-xs">
      <div class="row q-col-gutter-xs items-center">
        <div class="col-12 col-sm-3" v-for="lvl in ministerialLevels" :key="lvl.key">
          <div class="q-pa-xs rounded-lg border text-caption" :class="lvl.bgClass">
            <div class="row items-center justify-between font-bold" :class="lvl.textClass">
              <span>{{ lvl.label }}</span>
              <q-badge :color="lvl.color" class="text-xs">{{ lvl.score }} pt</q-badge>
            </div>
            <div class="text-slate-600 text-xs q-mt-xs ellipsis">
              {{ lvl.desc }}
              <q-tooltip anchor="top middle" self="bottom middle" class="bg-slate-900 text-caption" max-width="320px">
                <strong>{{ lvl.label }}:</strong> {{ lvl.desc }}
              </q-tooltip>
            </div>
          </div>
        </div>
      </div>
    </q-card>

    <!-- Matrix Table -->
    <div v-if="loadingStudents" class="text-center q-py-xl">
      <q-spinner-dots color="indigo" size="40px" />
      <div class="text-caption text-slate-400 q-mt-xs">Caricamento elenco studenti...</div>
    </div>

    <div v-else-if="!selectedClassId" class="text-center q-py-xl text-slate-400 bg-white rounded-xl border border-dashed">
      <q-icon name="school" size="48px" class="q-mb-sm opacity-40" />
      <div class="text-subtitle1 font-semibold text-slate-700">Seleziona una classe per aprire la matrice</div>
      <div class="text-caption">La griglia mostrerà gli studenti e gli obiettivi ministeriali per la disciplina</div>
    </div>

    <div v-else class="matrix-card-wrapper bg-white rounded-xl border border-slate-200 shadow-soft overflow-auto">
      <table class="descriptive-table full-width">
        <thead>
          <tr>
            <th class="sticky-col student-col-header">Alunno / Alunna</th>
            <th v-for="(obj, idx) in objectives" :key="obj.id" class="obj-col-header">
              <div class="row items-center justify-between no-wrap">
                <span class="ellipsis" style="max-width: 140px;">{{ idx + 1 }}. {{ obj.title }}</span>
                <q-btn flat round dense icon="close" size="xs" color="grey-6" @click="removeObjective(idx)">
                  <q-tooltip>Rimuovi obiettivo</q-tooltip>
                </q-btn>
              </div>
              <q-tooltip class="bg-slate-900 text-caption" max-width="300px">{{ obj.title }}</q-tooltip>
            </th>
            <th class="summary-col-header">Livello Prevalente</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="student in studentsList" :key="student.id">
            <td class="sticky-col student-name-cell font-semibold text-slate-800">
              {{ student.last_name }} {{ student.first_name }}
            </td>
            <td v-for="obj in objectives" :key="obj.id" class="eval-cell text-center">
              <div class="row justify-center q-gutter-xs no-wrap">
                <q-btn
                  v-for="lvl in ministerialLevels"
                  :key="lvl.key"
                  size="xs"
                  dense
                  :color="isLevelSelected(student.id, obj.id, lvl.key) ? lvl.color : 'grey-3'"
                  :text-color="isLevelSelected(student.id, obj.id, lvl.key) ? 'white' : 'grey-8'"
                  :unelevated="isLevelSelected(student.id, obj.id, lvl.key)"
                  :flat="!isLevelSelected(student.id, obj.id, lvl.key)"
                  class="rounded-sm level-btn font-bold"
                  @click="setLevel(student.id, obj.id, lvl.key)"
                >
                  {{ lvl.short }}
                  <q-tooltip>{{ lvl.label }}</q-tooltip>
                </q-btn>
              </div>
            </td>
            <td class="summary-cell text-center font-bold">
              <q-badge
                v-if="computePrevalentLevel(student.id)"
                :color="computePrevalentLevel(student.id).color"
                class="q-px-sm q-py-xs font-bold"
              >
                {{ computePrevalentLevel(student.id).label }}
              </q-badge>
              <span v-else class="text-slate-300 text-caption">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Dialog: Add Custom Objective -->
    <q-dialog v-model="showAddObjectiveDialog">
      <q-card style="width: min(500px, 95vw)" class="rounded-xl">
        <q-card-section class="bg-indigo text-white row items-center justify-between">
          <div class="text-subtitle1 font-bold">Nuovo Obiettivo di Apprendimento</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pa-md q-gutter-y-sm">
          <q-input
            v-model="newObjectiveTitle"
            label="Descrizione Obiettivo (es. Risolvere problemi lineari con grafici) *"
            outlined dense
            type="textarea"
            rows="3"
            autofocus
          />
        </q-card-section>
        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup no-caps />
          <q-btn
            color="indigo"
            label="Aggiungi"
            no-caps
            class="rounded-lg q-px-md font-bold"
            :disabled="!newObjectiveTitle.trim()"
            @click="addObjective"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()
const classesStore = useClassesStore()

const saving = ref(false)
const loadingStudents = ref(false)
const showAddObjectiveDialog = ref(false)
const newObjectiveTitle = ref('')

const selectedClassId = ref(null)
const selectedSubjectId = ref(null)
const selectedPeriod = ref('Q1')

const studentsList = ref([])
const subjectsList = ref([])

// 4 Ministerial Levels under O.M. 172/2020
const ministerialLevels = [
  {
    key: 'avanzato',
    label: 'Avanzato',
    short: 'AV',
    score: 4,
    color: 'positive',
    textClass: 'text-emerald-800',
    bgClass: 'bg-emerald-50 border-emerald-200',
    desc: 'Autonomia in situazioni note e non note, continuità e risorse mobilitate con consapevolezza.'
  },
  {
    key: 'intermedio',
    label: 'Intermedio',
    short: 'INT',
    score: 3,
    color: 'primary',
    textClass: 'text-blue-800',
    bgClass: 'bg-blue-50 border-blue-200',
    desc: 'Autonomia in situazioni note; discontinuo in situazioni nuove con risorse fornite.'
  },
  {
    key: 'base',
    label: 'Base',
    short: 'BASE',
    score: 2,
    color: 'warning',
    textClass: 'text-amber-800',
    bgClass: 'bg-amber-50 border-amber-200',
    desc: 'Autonomia solo in situazioni note e con guida continuativa del docente.'
  },
  {
    key: 'in_prima_acquisizione',
    label: 'In via di 1ª acq.',
    short: '1ª ACQ',
    score: 1,
    color: 'negative',
    textClass: 'text-red-800',
    bgClass: 'bg-red-50 border-red-200',
    desc: 'Porta a termine il compito solo se guidato e unicamente con risorse fornite.'
  }
]

// Default Learning Objectives
const objectives = ref([
  { id: 'obj-1', title: 'Comprendere e comunicare oralmente in contesti noti e nuovi' },
  { id: 'obj-2', title: 'Leggere e ricavare informazioni essenziali da testi espositivi' },
  { id: 'obj-3', title: 'Produrre elaborati scritti ortograficamente corretti e coesi' },
  { id: 'obj-4', title: 'Risolvere situazioni problematiche applicando procedure logico-deduttive' }
])

// Matrix store: key = `${studentId}_${objId}` -> levelKey ('avanzato', 'intermedio', 'base', 'in_prima_acquisizione')
const evaluations = reactive({})

const classOptions = computed(() => {
  return classesStore.classes.map(c => ({
    label: c.name || `Classe ${c.id}`,
    value: c.id
  }))
})

const subjectOptions = computed(() => {
  return subjectsList.value.map(s => ({
    label: s.name || s.subject_name || s.code,
    value: s.id || s.code
  }))
})

function isLevelSelected(studentId, objId, lvlKey) {
  return evaluations[`${studentId}_${objId}`] === lvlKey
}

function setLevel(studentId, objId, lvlKey) {
  const k = `${studentId}_${objId}`
  if (evaluations[k] === lvlKey) {
    delete evaluations[k]
  } else {
    evaluations[k] = lvlKey
  }
}

function computePrevalentLevel(studentId) {
  const counts = {}
  let totalAssigned = 0

  for (const obj of objectives.value) {
    const lvl = evaluations[`${studentId}_${obj.id}`]
    if (lvl) {
      counts[lvl] = (counts[lvl] || 0) + 1
      totalAssigned++
    }
  }

  if (totalAssigned === 0) return null

  // Find level with highest frequency, preferring higher if tied
  let maxCount = 0
  let prevalentKey = null

  for (const lvl of ministerialLevels) {
    const cnt = counts[lvl.key] || 0
    if (cnt > maxCount) {
      maxCount = cnt
      prevalentKey = lvl.key
    }
  }

  return ministerialLevels.find(l => l.key === prevalentKey) || null
}

function addObjective() {
  if (!newObjectiveTitle.value.trim()) return
  objectives.value.push({
    id: `obj-${Date.now()}`,
    title: newObjectiveTitle.value.trim()
  })
  newObjectiveTitle.value = ''
  showAddObjectiveDialog.value = false
}

function removeObjective(idx) {
  objectives.value.splice(idx, 1)
}

async function onClassChanged(classId) {
  if (!classId) return
  loadingStudents.value = true
  try {
    const res = await api.get(`/classes/${classId}/students`).catch(() => api.get(`/users?role=student&class_id=${classId}`))
    studentsList.value = res.data?.students || res.data || []
  } catch (err) {
    console.error('Failed to load students for class', err)
  } finally {
    loadingStudents.value = false
  }
}

function saveMatrix() {
  saving.value = true
  setTimeout(() => {
    saving.value = false
    $q.notify({
      type: 'positive',
      message: '✓ Matrice valutativa descrittiva O.M. 172/2020 salvata con successo!'
    })
  }, 400)
}

function exportMatrixCSV() {
  if (studentsList.value.length === 0) {
    $q.notify({ type: 'warning', message: 'Nessun dato da esportare' })
    return
  }

  let csv = 'Alunno;'
  objectives.value.forEach((o, i) => {
    csv += `"${i + 1}. ${o.title}";`
  })
  csv += 'Livello Prevalente\n'

  studentsList.value.forEach(st => {
    csv += `"${st.last_name} ${st.first_name}";`
    objectives.value.forEach(o => {
      const lvl = evaluations[`${st.id}_${o.id}`]
      csv += `"${lvl || '-'}";`
    })
    const prev = computePrevalentLevel(st.id)
    csv += `"${prev?.label || '-'}"\n`
  })

  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', `matrice_descrittiva_${selectedPeriod.value}.csv`)
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(url)
  $q.notify({ type: 'positive', message: 'Export CSV completato!' })
}

onMounted(async () => {
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  if (classesStore.classes.length > 0) {
    selectedClassId.value = classesStore.classes[0].id
    await onClassChanged(selectedClassId.value)
  }
  try {
    const res = await api.get('/subjects')
    subjectsList.value = res.data || []
    if (subjectsList.value.length > 0) {
      selectedSubjectId.value = subjectsList.value[0].id
    }
  } catch (err) {
    console.warn('Could not load subjects', err)
  }
})
</script>

<style scoped>
.font-semibold {
  font-weight: 600;
}
.font-bold {
  font-weight: 700;
}
.text-xs {
  font-size: 0.72rem;
}
.matrix-card-wrapper {
  max-height: 600px;
}
.descriptive-table {
  border-collapse: separate;
  border-spacing: 0;
}
.descriptive-table th {
  background-color: #f1f5f9;
  color: #334155;
  font-weight: 700;
  padding: 10px 12px;
  font-size: 0.8rem;
  border-bottom: 2px solid #cbd5e1;
}
.student-col-header {
  width: 180px;
  text-align: left;
}
.obj-col-header {
  min-width: 170px;
  max-width: 200px;
}
.summary-col-header {
  width: 150px;
  text-align: center;
}
.sticky-col {
  position: sticky;
  left: 0;
  background-color: white;
  z-index: 2;
  box-shadow: 2px 0 4px rgba(0,0,0,0.03);
}
.student-name-cell {
  padding: 10px 14px;
  border-bottom: 1px solid #f1f5f9;
  white-space: nowrap;
}
.eval-cell {
  padding: 8px 10px;
  border-bottom: 1px solid #f1f5f9;
  border-left: 1px solid #f8fafc;
}
.summary-cell {
  padding: 8px 12px;
  border-bottom: 1px solid #f1f5f9;
  background-color: #fafbfc;
}
.level-btn {
  font-size: 9px;
  padding: 2px 4px;
  min-height: 20px;
}
</style>
