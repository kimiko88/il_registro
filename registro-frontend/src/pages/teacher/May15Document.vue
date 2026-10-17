<template>
  <q-page class="q-pa-md may15-page bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="menu_book" color="primary" />
          {{ t('may15.title') || 'Documento del 15 Maggio' }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ t('may15.subtitle') || 'Relazione del Consiglio di Classe per la Commissione d\'Esame di Stato (Art. 17 D.Lgs. 62/2017)' }}
        </p>
      </div>

      <div class="row items-center q-gutter-sm">
        <q-chip
          :color="statusColor"
          text-color="white"
          class="text-weight-bold q-px-sm"
        >
          <q-icon :name="statusIcon" class="q-mr-xs" />
          {{ statusLabel }}
        </q-chip>

        <q-btn
          outline
          color="primary"
          icon="visibility"
          :label="t('may15.previewPrint') || 'Anteprima Ufficiale'"
          @click="showPrintPreview = true"
        />

        <q-btn
          color="primary"
          icon="save"
          :label="t('common.saveDraft') || 'Salva Bozza'"
          :loading="saving"
          @click="save(false)"
        />

        <q-btn
          color="positive"
          icon="verified"
          :label="t('may15.approveCdC') || 'Approva CdC'"
          :loading="saving"
          @click="approveCdC"
        />

        <q-btn
          color="deep-purple-7"
          icon="publish"
          :label="t('may15.publish') || 'Pubblica per Commissione'"
          :loading="publishing"
          @click="publish"
        />
      </div>
    </div>

    <!-- Filters & Selectors -->
    <q-card flat bordered class="q-mb-md rounded-xl bg-white shadow-soft">
      <q-card-section class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedClassId"
            :options="classOptions"
            emit-value
            map-options
            outlined
            dense
            :label="t('may15.selectClass') || 'Seleziona Classe 5ª'"
            @update:model-value="loadDocument"
          />
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="academicYear"
            :options="['2025/2026', '2024/2025', '2023/2024']"
            outlined
            dense
            :label="t('may15.academicYear') || 'Anno Scolastico'"
            @update:model-value="loadDocument"
          />
        </div>
      </q-card-section>
    </q-card>

    <!-- Main Content Panels -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="presentation" icon="group" :label="t('may15.tabPresentation') || '1. Presentazione Classe'" no-caps />
        <q-tab name="continuity" icon="history_edu" :label="t('may15.tabContinuity') || '2. Continuità Didattica'" no-caps />
        <q-tab name="pcto" icon="business_center" :label="t('may15.tabPCTO') || '3. PCTO & Ed. Civica'" no-caps />
        <q-tab name="simulations" icon="assignment" :label="t('may15.tabSimulations') || '4. Simulazioni Esame'" no-caps />
        <q-tab name="rubrics" icon="fact_check" :label="t('may15.tabRubrics') || '5. Griglie di Valutazione'" no-caps />
        <q-tab name="clil" icon="language" :label="t('may15.tabCLIL') || '6. Moduli CLIL'" no-caps />
      </q-tabs>

      <q-separator />

      <div v-if="loading" class="q-pa-xl text-center">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <q-tab-panels v-else v-model="activeTab" animated class="q-pa-md">
        <!-- 1. Presentazione Classe -->
        <q-tab-panel name="presentation" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.presentationTitle') || '1. Profilo Generale della Classe e Percorso Formativo' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.presentationHelp') || 'Inserire la composizione del gruppo classe, il clima relazionale, la motivazione, l\'evoluzione maturata nel triennio e gli obiettivi formativi globali.' }}
          </div>
          <q-input
            v-model="form.class_presentation"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.presentationPlaceholder') || 'La classe 5ª sez. A è composta da 21 studenti...'"
          />
        </q-tab-panel>

        <!-- 2. Continuità Didattica -->
        <q-tab-panel name="continuity" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.continuityTitle') || '2. Continuità Didattica e Quadro dei Docenti del Triennio' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.continuityHelp') || 'Illustrare l\'avvicendamento dei docenti delle singole discipline nelle classi 3ª, 4ª e 5ª e l\'impatto sulla programmazione didattica.' }}
          </div>
          <q-input
            v-model="form.teaching_continuity"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.continuityPlaceholder') || 'Si segnala continuità didattica triennale per le discipline di indirizzo...'"
          />
        </q-tab-panel>

        <!-- 3. PCTO & Educazione Civica -->
        <q-tab-panel name="pcto" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.pctoTitle') || '3. Percorsi per le Competenze Trasversali e per l\'Orientamento (PCTO) ed Educazione Civica' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.pctoHelp') || 'Riepilogo delle ore svolte in azienda, tirocini, project work, visite tecniche e nuclei tematici svolti nell\'ambito delle 33 ore di Educazione Civica.' }}
          </div>
          <q-input
            v-model="form.pcto_pathways"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.pctoPlaceholder') || 'Tutti gli studenti hanno completato il monte ore minimo prescritto con esperienze presso...'"
          />
        </q-tab-panel>

        <!-- 4. Simulazioni Esame -->
        <q-tab-panel name="simulations" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.simulationsTitle') || '4. Simulazioni delle Prove Nazionali d\'Esame' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.simulationsHelp') || 'Date, tipologie e risultanze delle simulazioni della Prima Prova Scritta (Italiano), della Seconda Prova Scritta e delle simulazioni del Colloquio Orale interdisciplinare.' }}
          </div>
          <q-input
            v-model="form.exam_simulations"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.simulationsPlaceholder') || 'Prima simulazione Prima Prova svolta il... Seconda simulazione svolta il...'"
          />
        </q-tab-panel>

        <!-- 5. Griglie di Valutazione -->
        <q-tab-panel name="rubrics" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.rubricsTitle') || '5. Criteri e Griglie di Valutazione Adottate' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.rubricsHelp') || 'Griglie nazionali e d\'istituto per la valutazione della 1ª prova, 2ª prova, colloquio e criteri di attribuzione della lode.' }}
          </div>
          <q-input
            v-model="form.evaluation_rubrics"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.rubricsPlaceholder') || 'Si applicano le griglie di valutazione allegate all\'O.M. Esami di Stato...'"
          />
        </q-tab-panel>

        <!-- 6. Moduli CLIL -->
        <q-tab-panel name="clil" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-primary">
            {{ t('may15.clilTitle') || '6. Insegnamento di Disciplina Non Linguistica in Lingua Straniera (CLIL)' }}
          </div>
          <div class="text-caption text-grey-7">
            {{ t('may15.clilHelp') || 'Disciplina coinvolta, lingua veicolare (es. Inglese B2), argomenti trattati, ore svolte e modalità di verifica.' }}
          </div>
          <q-input
            v-model="form.clil_modules"
            type="textarea"
            outlined
            rows="10"
            :placeholder="t('may15.clilPlaceholder') || 'Modulo di Scienze applicate in lingua inglese (30 ore)...'"
          />
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog Anteprima Ufficiale di Stampa -->
    <q-dialog v-model="showPrintPreview" max-width="850px" full-width>
      <q-card class="q-pa-lg print-card">
        <q-card-section class="row items-center justify-between no-print">
          <div class="text-h6 text-weight-bold text-primary">{{ t('may15.previewTitle') || 'Anteprima Documento Ufficiale' }}</div>
          <div class="row q-gutter-sm">
            <q-btn color="primary" icon="print" :label="t('common.print') || 'Stampa'" @click="printDoc" />
            <q-btn flat round icon="close" v-close-popup />
          </div>
        </q-card-section>

        <q-separator class="no-print q-mb-md" />

        <!-- Printable Document Body -->
        <div class="official-print-content q-pa-md">
          <div class="text-center q-mb-lg">
            <div class="text-subtitle2 text-weight-bold text-uppercase">Repubblica Italiana - Ministero dell'Istruzione e del Merito</div>
            <div class="text-h5 text-weight-bold text-uppercase q-mt-xs">Esame di Stato Conclusivo del Secondo Ciclo di Istruzione</div>
            <div class="text-subtitle1 text-weight-medium">Anno Scolastico {{ academicYear }}</div>
            <div class="text-h6 text-weight-bolder text-primary q-mt-sm">DOCUMENTO DEL CONSIGLIO DI CLASSE</div>
            <div class="text-caption text-grey-8">(ai sensi dell'art. 17, comma 1, del Decreto Legislativo 13 aprile 2017, n. 62)</div>
            <div class="text-subtitle1 text-weight-bold q-mt-xs">Classe: {{ selectedClassName }}</div>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">1. Profilo Generale della Classe</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.class_presentation || 'Nessuna presentazione registrata' }}</p>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">2. Continuità Didattica nel Triennio</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.teaching_continuity || 'Nessuna nota di continuità inserita' }}</p>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">3. Percorsi PCTO ed Educazione Civica</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.pcto_pathways || 'Nessun percorso PCTO registrato' }}</p>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">4. Simulazioni delle Prove d'Esame</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.exam_simulations || 'Nessuna simulazione d\'esame registrata' }}</p>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">5. Griglie e Criteri di Valutazione</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.evaluation_rubrics || 'Nessuna griglia inserita' }}</p>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold border-bottom q-pb-xs">6. Moduli di Apprendimento CLIL</div>
            <p class="text-body2 q-mt-xs whitespace-pre">{{ form.clil_modules || 'Nessun modulo CLIL registrato' }}</p>
          </div>

          <div class="row justify-between q-mt-xl q-pt-lg border-top">
            <div class="text-center">
              <div>Il Coordinatore del Consiglio di Classe</div>
              <div class="q-mt-md font-signature">___________________________</div>
            </div>
            <div class="text-center">
              <div>Il Dirigente Scolastico / Presidente CdC</div>
              <div class="q-mt-md font-signature">___________________________</div>
            </div>
          </div>
        </div>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import may15Service from '@/services/may15Service'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()

const selectedClassId = ref('')
const academicYear = ref('2025/2026')
const activeTab = ref('presentation')
const classOptions = ref([])
const loading = ref(false)
const saving = ref(false)
const publishing = ref(false)
const showPrintPreview = ref(false)

const form = ref({
  status: 'bozza',
  class_presentation: '',
  teaching_continuity: '',
  pcto_pathways: '',
  exam_simulations: '',
  evaluation_rubrics: '',
  clil_modules: ''
})

const selectedClassName = computed(() => {
  const c = classOptions.value.find(item => item.value === selectedClassId.value)
  return c ? c.label : ''
})

const statusLabel = computed(() => {
  switch (form.value.status) {
    case 'pubblicato':
      return t('may15.statusPublished') || 'Pubblicato per la Commissione'
    case 'approvato_cdc':
      return t('may15.statusApproved') || 'Approvato dal CdC'
    default:
      return t('may15.statusDraft') || 'Bozza di Lavoro'
  }
})

const statusColor = computed(() => {
  switch (form.value.status) {
    case 'pubblicato':
      return 'deep-purple-7'
    case 'approvato_cdc':
      return 'positive'
    default:
      return 'warning'
  }
})

const statusIcon = computed(() => {
  switch (form.value.status) {
    case 'pubblicato':
      return 'verified'
    case 'approvato_cdc':
      return 'check_circle'
    default:
      return 'edit_note'
  }
})

async function loadClasses() {
  try {
    const res = await api.get('/classes')
    const list = res.data?.classes || res.data || []
    classOptions.value = list.map(c => ({
      label: c.name || `Classe ${c.id}`,
      value: c.id
    }))
    if (classOptions.value.length > 0 && !selectedClassId.value) {
      selectedClassId.value = classOptions.value[0].value
      await loadDocument()
    }
  } catch (err) {
    console.error('Error loading classes:', err)
  }
}

async function loadDocument() {
  if (!selectedClassId.value) return
  loading.value = true
  try {
    const res = await may15Service.getDocument(selectedClassId.value, academicYear.value)
    const doc = res.data || {}
    form.value = {
      status: doc.status || 'bozza',
      class_presentation: doc.class_presentation || '',
      teaching_continuity: doc.teaching_continuity || '',
      pcto_pathways: doc.pcto_pathways || '',
      exam_simulations: doc.exam_simulations || '',
      evaluation_rubrics: doc.evaluation_rubrics || '',
      clil_modules: doc.clil_modules || ''
    }
  } catch (err) {
    console.error('Error loading may15 document:', err)
  } finally {
    loading.value = false
  }
}

async function save(notify = true) {
  if (!selectedClassId.value) return
  saving.value = true
  try {
    await may15Service.saveDocument(selectedClassId.value, {
      academic_year: academicYear.value,
      status: form.value.status,
      ...form.value
    })
    if (notify) {
      $q.notify({ type: 'positive', message: t('common.success') || 'Salvato con successo' })
    }
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

async function approveCdC() {
  form.value.status = 'approvato_cdc'
  await save(false)
  $q.notify({
    type: 'positive',
    icon: 'check_circle',
    message: t('may15.approvedSuccess') || 'Documento approvato dal Consiglio di Classe'
  })
}

async function publish() {
  if (!selectedClassId.value) return
  publishing.value = true
  try {
    await may15Service.publishDocument(selectedClassId.value, academicYear.value)
    form.value.status = 'pubblicato'
    $q.notify({
      type: 'positive',
      icon: 'verified',
      message: t('may15.publishedSuccess') || 'Documento del 15 Maggio pubblicato ufficialmente per la Commissione!'
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore durante la pubblicazione' })
  } finally {
    publishing.value = false
  }
}

function printDoc() {
  window.print()
}

onMounted(() => {
  loadClasses()
})
</script>

<style scoped>
.whitespace-pre {
  white-space: pre-wrap;
}
.border-bottom {
  border-bottom: 2px solid #e2e8f0;
}
.border-top {
  border-top: 1px solid #cbd5e1;
}
@media print {
  .no-print {
    display: none !important;
  }
  .print-card {
    box-shadow: none !important;
    border: none !important;
    padding: 0 !important;
  }
}
</style>
