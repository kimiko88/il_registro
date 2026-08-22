<template>
  <div class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold text-primary q-my-none">
          <q-icon name="stars" class="q-mr-sm" />
          {{ t('competenciesPage.title') || 'Valutazione delle Competenze (DM 742/2017)' }}
        </h5>
        <div class="text-caption text-grey-7">{{ t('competenciesPage.subtitle') || 'Certificazione delle competenze per la scuola secondaria di primo grado' }}</div>
      </div>
      <q-btn color="secondary" icon="picture_as_pdf" :label="t('competenciesPage.downloadPdf') || 'Scarica Certificato PDF'" @click="downloadPdfCertificate" class="glossy" />
    </div>

    <!-- Student Selector -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedStudentId"
            :options="studentOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            :label="t('competenciesPage.student') || 'Seleziona Studente'"
            outlined
            dense
            @update:model-value="fetchEvaluations"
          />
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedSemester"
            :options="[
              { label: t('competenciesPage.q1') || '1° Quadrimestre', value: 1 },
              { label: t('competenciesPage.q2') || '2° Quadrimestre (Certificazione Finale)', value: 2 }
            ]"
            emit-value
            map-options
            :label="t('competenciesPage.semester') || 'Quadrimestre'"
            outlined
            dense
            @update:model-value="fetchEvaluations"
          />
        </div>
      </q-card-section>
    </q-card>

    <!-- Competence Grid Table -->
    <q-card flat bordered>
      <q-card-section class="bg-grey-2 text-weight-bold">
        {{ t('competenciesPage.gridTitle') || 'Griglia Competenze Chiave Europee (D.M. 742/2017)' }}
      </q-card-section>

      <q-markup-table flat separator="cell">
        <thead>
          <tr>
            <th class="text-left" style="width: 30%">{{ t('competenciesPage.colCompetence') || 'Competenza Chiave' }}</th>
            <th class="text-center" style="width: 40%">{{ t('competenciesPage.colLevel') || 'Livello di Padronanza (DM 742)' }}</th>
            <th class="text-left" style="width: 30%">{{ t('competenciesPage.colNotes') || 'Descrittore / Note' }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="comp in competenceList" :key="comp.code">
            <td class="text-weight-medium">
              <div>{{ comp.name }}</div>
              <div class="text-caption text-grey-7">{{ comp.code }}</div>
            </td>
            <td class="text-center">
              <q-btn-toggle
                v-model="comp.level"
                dense
                unelevated
                toggle-color="primary"
                :options="levelOptions"
                @update:model-value="saveEvaluation(comp)"
              />
            </td>
            <td>
              <q-input
                v-model="comp.descriptor"
                dense
                borderless
                :placeholder="t('competenciesPage.addDescriptor') || 'Aggiungi descrittore...'"
                @blur="saveEvaluation(comp)"
              />
            </td>
          </tr>
        </tbody>
      </q-markup-table>
    </q-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import competenciesService from '@/services/competenciesService'
import { useClassesStore } from '@/stores/classes'
import { getBaseURL } from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()

const selectedStudentId = ref('')
const selectedSemester = ref(1)

const studentOptions = computed(() => {
  const students = classesStore.classStudents || []
  if (students.length > 0) {
    return students.map(s => ({
      id: s.id || s.user_id,
      name: `${s.last_name || ''} ${s.first_name || ''}`.trim() || s.name || s.id
    }))
  }
  return [
    { id: 'stu-demo-1', name: 'Rossi Mario (2A)' },
    { id: 'stu-demo-2', name: 'Bianchi Giulia (2A)' }
  ]
})

const competenceList = ref([
  { code: 'COMP_L1_ITA', name: 'Comunicazione nella madrelingua / lingua di istruzione', level: 'A_Avanzato', descriptor: 'Padroneggia la lingua con precisione e ricchezza lessicale.' },
  { code: 'COMP_L2_ENG', name: 'Comunicazione nelle lingue straniere (Inglese/Francese)', level: 'B_Intermedio', descriptor: 'Comprende e produce testi scritti e orali di uso quotidiano.' },
  { code: 'COMP_STEM', name: 'Competenza matematica e competenze di base in scienza e tecnologia', level: 'A_Avanzato', descriptor: 'Risolve problemi complessi applicando modelli matematici.' },
  { code: 'COMP_DIGITAL', name: 'Competenza digitale', level: 'A_Avanzato', descriptor: 'Utilizza con consapevolezza e responsabilità le tecnologie digitali.' },
  { code: 'COMP_LEARNING', name: 'Imparare a imparare', level: 'B_Intermedio', descriptor: 'Organizza in modo autonomo il proprio lavoro di studio.' },
  { code: 'COMP_CIVIC', name: 'Competenze sociali e civiche / Educazione Civica', level: 'A_Avanzato', descriptor: 'Rispetta le regole di convivenza civile e partecipa attivamente.' },
  { code: 'COMP_INITIATIVE', name: 'Spirito di iniziativa e imprenditorialità', level: 'B_Intermedio', descriptor: 'Propone soluzioni originali in contesti operativi.' },
  { code: 'COMP_CULTURE', name: 'Consapevolezza ed espressione culturale', level: 'A_Avanzato', descriptor: 'Esprime idee e sentimenti con diversi linguaggi espressivi.' }
])

async function fetchEvaluations() {
  if (!selectedStudentId.value) return
  try {
    const data = await competenciesService.getStudentEvaluations(selectedStudentId.value, selectedSemester.value)
    if (data && data.length > 0) {
      data.forEach(item => {
        const found = competenceList.value.find(c => c.code === item.competence_code)
        if (found) {
          found.level = item.level
          found.descriptor = item.descriptor
        }
      })
    }
  } catch (err) {
    console.error('Error fetching competencies:', err)
  }
}

watch(studentOptions, (opts) => {
  if (opts.length > 0 && !selectedStudentId.value) {
    selectedStudentId.value = opts[0].id
    fetchEvaluations()
  }
}, { immediate: true })

const saveEvaluation = async (comp) => {
  if (!selectedStudentId.value) return
  try {
    await competenciesService.saveEvaluation({
      student_id: selectedStudentId.value,
      class_id: classesStore.selectedClassId || '',
      semester: selectedSemester.value,
      competence_code: comp.code,
      competence_name: comp.name,
      level: comp.level,
      descriptor: comp.descriptor || ''
    })
    $q.notify({ type: 'positive', message: t('common.success') })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  }
}

const downloadPdfCertificate = () => {
  if (!selectedStudentId.value) return
  $q.notify({ type: 'info', message: t('competenciesPage.generatingPdf') || 'Certificazione Competenze DM 742 in generazione PDF...' })
  setTimeout(() => {
    const base = getBaseURL()
    window.open(`${base}/competencies/student/${selectedStudentId.value}`, '_blank')
  }, 500)
}

onMounted(() => {
  if (studentOptions.value.length > 0 && !selectedStudentId.value) {
    selectedStudentId.value = studentOptions.value[0].id
  }
  fetchEvaluations()
})
</script>
