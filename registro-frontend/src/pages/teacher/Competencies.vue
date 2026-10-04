<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="stars" color="primary" class="q-mr-sm" />
          {{ t('competenciesPage.title') || 'Certificazione delle Competenze (DM 742/2017 & DM 14/2024)' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('competenciesPage.subtitle') || 'Valutazione delle 8 Competenze Chiave Europee per il Consiglio di Classe' }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-chip
          :color="isDraft ? 'warning' : 'positive'"
          text-color="white"
          class="text-weight-bold q-px-sm"
        >
          <q-icon :name="isDraft ? 'edit_note' : 'check_circle'" class="q-mr-xs" />
          {{ isDraft ? (t('competenciesPage.draftBadge') || 'Modifiche non salvate') : (t('competenciesPage.syncedBadge') || 'Sincronizzato') }}
        </q-chip>
        <q-btn
          color="primary"
          icon="save"
          :label="t('competenciesPage.saveToDb') || 'Salva Valutazioni'"
          unelevated
          :loading="saving"
          class="rounded-lg text-weight-bold"
          @click="saveCompetencies"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchCompetenciesData" />
      </div>
    </div>

    <!-- Draft Warning Banner -->
    <q-banner v-if="isDraft" class="bg-amber-1 text-amber-10 rounded-xl border border-amber-300 q-mb-lg shadow-soft">
      <template v-slot:avatar>
        <q-icon name="warning" color="warning" size="24px" />
      </template>
      <div class="text-weight-bold text-subtitle2">{{ t('competenciesPage.draftWarningTitle') || 'Modifiche in bozza' }}</div>
      <div class="text-caption">
        {{ t('competenciesPage.draftWarningDesc') || 'Le modifiche ai livelli di competenza non sono ancora state persistite. Clicca su Salva per confermare.' }}
      </div>
    </q-banner>

    <!-- Filters Bar (Trasversale - Consiglio di Classe) -->
    <q-card flat bordered class="rounded-xl bg-white q-mb-lg shadow-soft q-pa-md">
      <div class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedClass"
            :options="classOptions"
            :label="t('competenciesPage.selectClass') || 'Seleziona Classe'"
            outlined dense emit-value map-options
            @update:model-value="fetchCompetenciesData"
          />
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedPeriod"
            :options="[
              { label: t('competenciesPage.q1') || '1° Quadrimestre (Valutazione Intermedia)', value: 'q1' },
              { label: t('competenciesPage.q2') || '2° Quadrimestre', value: 'q2' },
              { label: t('competenciesPage.final') || 'Certificazione Finale (Termine Ciclo)', value: 'final' }
            ]"
            emit-value map-options
            :label="t('competenciesPage.period') || 'Periodo di Riferimento'"
            outlined dense
            @update:model-value="markAsDraft"
          />
        </div>
      </div>
    </q-card>

    <!-- Ministerial Legend Card (DM 742/2017 & DM 14/2024) -->
    <q-card flat bordered class="rounded-xl bg-blue-50 border-blue-200 q-mb-lg shadow-soft q-pa-md">
      <div class="text-subtitle2 text-weight-bold text-primary q-mb-xs">
        <q-icon name="info" class="q-mr-xs" />
        {{ t('competenciesPage.legendTitle') || 'Livelli di Padronanza Ministeriali (D.M. 742/2017 & D.M. 14/2024)' }}
      </div>
      <div class="row q-col-gutter-sm text-caption">
        <div class="col-12 col-sm-3">
          <q-badge color="positive" class="q-mr-xs">A</q-badge>
          <strong>{{ t('competenciesPage.levelAdvanced') || 'Avanzato:' }}</strong>
          {{ t('competenciesPage.descAdvanced') || 'Svolge compiti e risolve problemi complessi mostrando autonomia e padronanza.' }}
        </div>
        <div class="col-12 col-sm-3">
          <q-badge color="primary" class="q-mr-xs">B</q-badge>
          <strong>{{ t('competenciesPage.levelIntermediate') || 'Intermedio:' }}</strong>
          {{ t('competenciesPage.descIntermediate') || 'Risolve problemi in situazioni nuove, compie scelte consapevoli.' }}
        </div>
        <div class="col-12 col-sm-3">
          <q-badge color="warning" class="q-mr-xs">C</q-badge>
          <strong>{{ t('competenciesPage.levelBase') || 'Base:' }}</strong>
          {{ t('competenciesPage.descBase') || 'Svolge compiti semplici anche in situazioni nuove applicando regole fondamentali.' }}
        </div>
        <div class="col-12 col-sm-3">
          <q-badge color="negative" class="q-mr-xs">D</q-badge>
          <strong>{{ t('competenciesPage.levelInitial') || 'Iniziale:' }}</strong>
          {{ t('competenciesPage.descInitial') || 'Se opportunamente guidato/a, svolge compiti semplici in situazioni note.' }}
        </div>
      </div>
    </q-card>

    <!-- Table Grid across all 8 European Competencies -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <div v-if="loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <div v-else-if="studentEvaluations.length === 0" class="text-center q-pa-xl text-slate-400">
        <q-icon name="group_off" size="64px" class="q-mb-md opacity-40" />
        <div class="text-h6">{{ t('competenciesPage.noStudentsFound') || 'Nessuno studente trovato per la classe selezionata' }}</div>
        <div class="text-caption">{{ t('competenciesPage.selectClassPrompt') || 'Seleziona una classe dal menu in alto' }}</div>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="q-table q-table--horizontal-separator full-width">
          <thead>
            <tr class="bg-slate-100 text-slate-700">
              <th class="text-left q-pa-md" style="min-width: 220px">{{ t('competenciesPage.student') || 'Alunno / Alunna' }}</th>
              <th v-for="comp in competencyColumns" :key="comp.id" class="text-center q-pa-sm" style="min-width: 140px">
                <div class="text-weight-bold">{{ comp.title }}</div>
                <div class="text-caption text-grey-6 text-weight-normal">{{ comp.subtitle }}</div>
              </th>
              <th class="text-center q-pa-md" style="min-width: 120px">{{ t('competenciesPage.actions') || 'Certificato' }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in studentEvaluations" :key="student.student_id" class="hover-bg-slate">
              <td class="text-left q-pa-md">
                <div class="row items-center no-wrap">
                  <q-avatar color="indigo-1" text-color="indigo-8" icon="person" size="36px" class="q-mr-sm" />
                  <div>
                    <div class="text-weight-bold text-slate-800">{{ student.student_name }}</div>
                    <div class="text-caption text-grey-6">Matr. {{ student.student_id ? student.student_id.substring(0, 8) : 'N/D' }}</div>
                  </div>
                </div>
              </td>
              <td v-for="comp in competencyColumns" :key="comp.id" class="text-center q-pa-xs">
                <q-btn-toggle
                  v-model="student.evaluations[comp.id]"
                  toggle-color="primary"
                  color="grey-2"
                  text-color="grey-8"
                  dense unelevated
                  :options="[
                    { label: 'A', value: 'A' },
                    { label: 'B', value: 'B' },
                    { label: 'C', value: 'C' },
                    { label: 'D', value: 'D' }
                  ]"
                  @update:model-value="markAsDraft"
                />
              </td>
              <td class="text-center q-pa-sm">
                <q-btn
                  flat
                  round
                  color="primary"
                  icon="picture_as_pdf"
                  @click="openCertificateModal(student)"
                >
                  <q-tooltip>{{ t('competenciesPage.printCertificate') || 'Visualizza Certificato Ministeriale D.M. 742' }}</q-tooltip>
                </q-btn>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </q-card>

    <!-- Dialog Certificato Ministeriale Nazionale (D.M. 742/2017 & D.M. 14/2024) -->
    <q-dialog v-model="certificateModalOpen" max-width="900px" full-width>
      <q-card class="q-pa-lg print-certificate-card">
        <q-card-section class="row items-center justify-between no-print">
          <div class="text-h6 text-weight-bold text-primary">
            {{ t('competenciesPage.certificateModalTitle') || 'Certificazione delle Competenze - Modello Nazionale' }}
          </div>
          <div class="row q-gutter-sm">
            <q-btn color="primary" icon="print" :label="t('common.print') || 'Stampa'" @click="printCertificate" />
            <q-btn flat round icon="close" v-close-popup />
          </div>
        </q-card-section>

        <q-separator class="no-print q-mb-md" />

        <div v-if="activeCertificateStudent" class="q-pa-md official-certificate-body">
          <div class="text-center q-mb-md">
            <div class="text-caption text-weight-bold text-uppercase">Ministero dell'Istruzione e del Merito</div>
            <div class="text-h6 text-weight-bold text-uppercase">MODELLO NAZIONALE DI CERTIFICAZIONE DELLE COMPETENZE</div>
            <div class="text-caption text-grey-8">(adottato con D.M. n. 742/2017 ed integrato con D.M. n. 14/2024)</div>
            <div class="text-subtitle1 text-weight-bold q-mt-sm">
              Alunno/a: <span class="text-primary">{{ activeCertificateStudent.student_name }}</span>
            </div>
            <div class="text-caption">Classe: {{ selectedClassLabel }} &bull; Anno Scolastico 2025/2026</div>
          </div>

          <table class="q-table q-table--cell-separator full-width q-mb-md">
            <thead>
              <tr class="bg-grey-2">
                <th class="text-left" style="width: 35%">Competenza Chiave Europea</th>
                <th class="text-center" style="width: 15%">Livello Attribuito</th>
                <th class="text-left" style="width: 50%">Descrittore Ministeriale del Livello</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="comp in competencyColumns" :key="comp.id">
                <td class="text-weight-bold text-slate-800">
                  <div>{{ comp.title }}</div>
                  <div class="text-caption text-grey-6">{{ comp.subtitle }}</div>
                </td>
                <td class="text-center">
                  <q-badge
                    :color="getLevelBadgeColor(activeCertificateStudent.evaluations[comp.id])"
                    class="text-weight-bolder text-subtitle2 q-px-sm"
                  >
                    {{ activeCertificateStudent.evaluations[comp.id] || 'B' }}
                  </q-badge>
                </td>
                <td class="text-caption">
                  {{ getDescriptorForLevel(comp.id, activeCertificateStudent.evaluations[comp.id] || 'B') }}
                </td>
              </tr>
            </tbody>
          </table>

          <div class="row justify-between q-mt-xl q-pt-md border-top">
            <div class="text-center">
              <div class="text-caption">I Docenti del Consiglio di Classe</div>
              <div class="q-mt-sm">___________________________</div>
            </div>
            <div class="text-center">
              <div class="text-caption">Il Dirigente Scolastico</div>
              <div class="q-mt-sm">___________________________</div>
            </div>
          </div>
        </div>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const selectedClass = ref(null)
const selectedPeriod = ref('q1')
const classOptions = ref([])
const isDraft = ref(false)
const loading = ref(false)
const saving = ref(false)
const certificateModalOpen = ref(false)
const activeCertificateStudent = ref(null)

const competencyColumns = ref([
  { id: 'COMP_L1_ITA', title: 'Comunicazione Madrelingua', subtitle: 'Lingua d\'istruzione' },
  { id: 'COMP_L2_ENG', title: 'Lingue Straniere', subtitle: 'Inglese / Seconda Lingua' },
  { id: 'COMP_STEM', title: 'Competenza STEM', subtitle: 'Matematica, Scienze, Tecnologia' },
  { id: 'COMP_DIGITAL', title: 'Competenza Digitale', subtitle: 'Media e tecnologie' },
  { id: 'COMP_LEARNING', title: 'Imparare a Imparare', subtitle: 'Autonomia di studio' },
  { id: 'COMP_CIVIC', title: 'Sociali e Civiche', subtitle: 'Educazione Civica' },
  { id: 'COMP_INITIATIVE', title: 'Spirito di Iniziativa', subtitle: 'Creatività e problem solving' },
  { id: 'COMP_CULTURE', title: 'Espressione Culturale', subtitle: 'Arte, Musica e Patrimonio' }
])

const studentEvaluations = ref([])

const selectedClassLabel = computed(() => {
  const f = classOptions.value.find(c => c.value === selectedClass.value)
  return f ? f.label : ''
})

const markAsDraft = () => {
  isDraft.value = true
}

const fetchOptions = async () => {
  try {
    const cRes = await api.get('/classes')
    const classes = cRes.data?.classes || cRes.data || []
    classOptions.value = classes.map(c => ({ label: `Classe ${c.name || c.section || c.id}`, value: c.id }))
    if (classOptions.value.length) selectedClass.value = classOptions.value[0].value
  } catch {
    /* ignore options fetch errors */
  }
}

const fetchCompetenciesData = async () => {
  if (!selectedClass.value) return
  loading.value = true
  try {
    const res = await api.get('/competencies', {
      params: { class_id: selectedClass.value }
    })
    const data = res.data?.evaluations || res.data || []
    studentEvaluations.value = data.map(item => {
      const ev = item.evaluations || {}
      // Normalize to 8 competencies with fallback
      const normalizedEv = {}
      competencyColumns.value.forEach(col => {
        normalizedEv[col.id] = ev[col.id] || ev[col.id.replace('COMP_', '').toLowerCase()] || 'B'
      })
      return {
        student_id: item.student_id,
        student_name: item.student_name,
        evaluations: normalizedEv
      }
    })
    isDraft.value = false
  } catch (err) {
    studentEvaluations.value = []
  } finally {
    loading.value = false
  }
}

const saveCompetencies = async () => {
  if (!selectedClass.value) {
    $q.notify({ type: 'warning', message: t('competenciesPage.selectClassPrompt') || 'Selezionare una classe' })
    return
  }
  saving.value = true
  try {
    await api.post('/competencies/batch', {
      class_id: selectedClass.value,
      period: selectedPeriod.value,
      evaluations: studentEvaluations.value
    })
    isDraft.value = false
    $q.notify({ type: 'positive', message: t('competenciesPage.savedSuccess') || 'Valutazioni per competenze salvate con successo' })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

function openCertificateModal(student) {
  activeCertificateStudent.value = student
  certificateModalOpen.value = true
}

function printCertificate() {
  window.print()
}

function getLevelBadgeColor(level) {
  switch (level) {
    case 'A': return 'positive'
    case 'B': return 'primary'
    case 'C': return 'warning'
    case 'D': return 'negative'
    default: return 'primary'
  }
}

function getDescriptorForLevel(compId, level) {
  switch (level) {
    case 'A':
      return 'L\'alunno/a svolge compiti e risolve problemi complessi, mostrando padronanza nell\'uso delle conoscenze e delle abilità; propone e sostiene le proprie opinioni e assume in modo responsabile decisioni consapevoli.'
    case 'B':
      return 'L\'alunno/a svolge compiti e risolve problemi in situazioni nuove, compie scelte consapevoli, mostrando di saper utilizzare le conoscenze e le abilità acquisite.'
    case 'C':
      return 'L\'alunno/a svolge compiti semplici anche in situazioni nuove, mostrando di possedere conoscenze e abilità fondamentali e di saper applicare basilari regole e procedure apprese.'
    case 'D':
      return 'L\'alunno/a, se opportunamente guidato/a, svolge compiti semplici in situazioni note.'
    default:
      return 'Livello intermedio di padronanza conforme alle indicazioni nazionali.'
  }
}

onMounted(async () => {
  await fetchOptions()
  await fetchCompetenciesData()
})
</script>

<style scoped>
.hover-bg-slate:hover {
  background-color: #f8fafc;
}
.border-top {
  border-top: 1px solid #cbd5e1;
}
@media print {
  .no-print {
    display: none !important;
  }
  .print-certificate-card {
    box-shadow: none !important;
    border: none !important;
    padding: 0 !important;
  }
}
</style>
