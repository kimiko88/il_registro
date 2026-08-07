<template>
  <div class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold text-primary q-my-none">
          <q-icon name="stars" class="q-mr-sm" />
          Valutazione delle Competenze (DM 742/2017)
        </h5>
        <div class="text-caption text-grey-7">Certificazione delle competenze per la scuola secondaria di primo grado</div>
      </div>
      <q-btn color="secondary" icon="picture_as_pdf" label="Scarica Certificato PDF" @click="downloadPdfCertificate" class="glossy" />
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
            label="Seleziona Studente"
            outlined
            dense
            @update:model-value="fetchEvaluations"
          />
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedSemester"
            :options="[
              { label: '1° Quadrimestre', value: 1 },
              { label: '2° Quadrimestre (Certificazione Finale)', value: 2 }
            ]"
            emit-value
            map-options
            label="Quadrimestre"
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
        Griglia Competenze Chiave Europee (D.M. 742/2017)
      </q-card-section>

      <q-markup-table flat separator="cell">
        <thead>
          <tr>
            <th class="text-left" style="width: 30%">Competenza Chiave</th>
            <th class="text-center" style="width: 40%">Livello di Padronanza (DM 742)</th>
            <th class="text-left" style="width: 30%">Descrittore / Note</th>
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
                placeholder="Aggiungi descrittore..."
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
import { ref, onMounted } from 'vue'
import competenciesService from 'src/services/competenciesService'
import { useNotify } from 'src/composables/useNotify'

const notify = useNotify()

const selectedStudentId = ref('stu-demo-1')
const selectedSemester = ref(1)

const studentOptions = ref([
  { id: 'stu-demo-1', name: 'Rossi Mario (2A)' },
  { id: 'stu-demo-2', name: 'Bianchi Giulia (2A)' }
])

const levelOptions = [
  { label: 'A - Avanzato', value: 'A_Avanzato' },
  { label: 'B - Intermedio', value: 'B_Intermedio' },
  { label: 'C - Base', value: 'C_Base' },
  { label: 'D - Iniziale', value: 'D_Iniziale' }
]

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

const fetchEvaluations = async () => {
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

const saveEvaluation = async (comp) => {
  try {
    await competenciesService.saveEvaluation({
      student_id: selectedStudentId.value,
      class_id: '47a05d80-3836-452e-ac91-8cfa3a1999dd',
      semester: selectedSemester.value,
      competence_code: comp.code,
      competence_name: comp.name,
      level: comp.level,
      descriptor: comp.descriptor || ''
    })
    notify.success(`Competenza "${comp.code}" salvata`)
  } catch (err) {
    notify.error('Errore durante il salvataggio')
  }
}

const downloadPdfCertificate = () => {
  notify.success('Certificazione Competenze DM 742 in generazione PDF...')
  setTimeout(() => {
    window.open(`/api/v1/competencies/student/${selectedStudentId.value}`, '_blank')
  }, 500)
}

onMounted(() => {
  fetchEvaluations()
})
</script>
