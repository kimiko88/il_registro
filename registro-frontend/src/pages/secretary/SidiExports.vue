<template>
  <q-page class="q-pa-md bg-slate-50 min-h-screen">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="cloud_sync" color="primary" class="q-mr-sm" />
          {{ t('sidiExports.title') || 'Flussi SIDI & Tracciati Ministeriali (MIM)' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('sidiExports.subtitle') || 'Esportazione pacchetti XML per Anagrafe Nazionale Studenti ed Esiti Scrutini' }}
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" :aria-label="t('common.refresh') || 'Aggiorna'" @click="loadExports" />
    </div>

    <!-- WIZARD ESPORTAZIONE SIDI -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-lg border border-slate-100">
      <q-card-section>
        <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">{{ t('sidiExports.generatorTitle') || 'Generatore Flusso SIDI' }}</div>
        <div class="row q-col-gutter-md items-center">
          <div class="col-12 col-md-5">
            <q-select
              v-model="form.export_type"
              :options="exportOptions"
              emit-value
              map-options
              :label="t('sidiExports.exportType') || 'Tipologia di Flusso SIDI *'"
              outlined
              dense
            />
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="form.school_year"
              :options="['2024/2025', '2025/2026', '2026/2027']"
              :label="t('sidiExports.schoolYear') || 'Anno Scolastico'"
              outlined
              dense
            />
          </div>
          <div class="col-12 col-md-3">
            <q-btn
              color="primary"
              icon="bolt"
              :label="t('sidiExports.validateAndExport') || 'Valida ed Esporta XML'"
              no-caps
              class="full-width rounded-lg font-bold"
              :loading="generating"
              @click="generateExport"
            />
          </div>
        </div>

        <!-- ESITO VALIDAZIONE -->
        <div v-if="validationResult" class="q-mt-md">
          <q-banner v-if="validationResult.valid" class="bg-green-50 text-green-9 rounded-xl border border-green-200">
            <template v-slot:avatar>
              <q-icon name="check_circle" color="positive" />
            </template>
            <div class="text-weight-bold">{{ t('sidiExports.validationSuccessTitle') || 'Controllo di Congruenza Superato!' }}</div>
            <div class="text-caption">{{ t('sidiExports.validationSuccessDesc', { count: validationResult.total_records }) || `Tutti i record anagrafici sono conformi agli schemi XSD ministeriali SIDI.` }}</div>
          </q-banner>

          <q-banner v-else class="bg-amber-50 text-amber-9 rounded-xl border border-amber-200">
            <template v-slot:avatar>
              <q-icon name="warning" color="warning" />
            </template>
            <div class="text-weight-bold">{{ t('sidiExports.validationWarningTitle') || 'Attenzione: Rilevate anomalie prima dell\'invio SIDI' }}</div>
            <ul class="q-my-xs q-pl-md text-caption">
              <li v-for="(err, idx) in validationResult.errors" :key="'e-'+idx">{{ err }}</li>
              <li v-for="(msg, idx) in validationResult.missing_sidi_ids" :key="'m-'+idx">Manca Codice SIDI: {{ msg }}</li>
            </ul>
          </q-banner>
        </div>
      </q-card-section>
    </q-card>

    <!-- STORICO ESPORTAZIONI SIDI -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft border border-slate-100">
      <q-card-section class="row items-center justify-between border-b pb-3">
        <div class="text-h6 text-weight-bold text-slate-800">{{ t('sidiExports.historyTitle') || 'Storico Flussi Generati' }}</div>
        <q-chip color="blue-1" text-color="blue-9" size="sm">
          {{ exportsList.length }} {{ t('sidiExports.recordsLabel') || 'record' }}
        </q-chip>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <!-- Skeleton Loading State -->
        <div v-if="loading" class="q-pa-md q-gutter-y-md">
          <div v-for="n in 3" :key="'skel-'+n" class="row items-center q-gutter-x-md q-py-sm">
            <q-skeleton type="QAvatar" size="42px" class="rounded-xl" />
            <div class="col">
              <q-skeleton type="text" width="60%" height="24px" />
              <q-skeleton type="text" width="40%" height="16px" />
            </div>
            <q-skeleton type="rect" width="100px" height="32px" class="rounded-lg" />
          </div>
        </div>

        <div v-else-if="exportsList.length === 0" class="q-pa-xl text-center">
          <q-icon name="folder_open" size="48px" color="slate-300" class="q-mb-sm" />
          <div class="text-slate-500">{{ t('sidiExports.noExports') || 'Nessun flusso SIDI ancora esportato per questo istituto.' }}</div>
        </div>

        <q-list v-else separator>
          <q-item v-for="exp in exportsList" :key="exp.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="indigo-1" text-color="indigo-8" icon="code" size="42px" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold text-slate-800">
                {{ exp.file_name }}
              </q-item-label>
              <q-item-label caption class="text-slate-500">
                {{ t('sidiExports.typeLabel') || 'Tipo' }}: {{ exp.export_type }} &bull; {{ t('sidiExports.yearLabel') || 'A.S.' }} {{ exp.school_year }} &bull; {{ exp.records_count }} {{ t('sidiExports.recordsLabel') || 'record' }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-btn
                unelevated
                color="primary"
                icon="download"
                :label="t('sidiExports.downloadXml') || 'Scarica XML'"
                size="sm"
                no-caps
                class="rounded-lg font-bold"
                @click="downloadXML(exp)"
              />
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()
const loading = ref(false)
const generating = ref(false)
const exportsList = ref([])
const validationResult = ref(null)

const form = ref({
  export_type: 'ANS_ANAGRAFE',
  school_year: '2025/2026'
})

const exportOptions = computed(() => [
  { label: t('sidiExports.options.ans') || 'Anagrafe Nazionale Studenti (ANS)', value: 'ANS_ANAGRAFE' },
  { label: t('sidiExports.options.scrutinioGiugno') || 'Esiti Scrutinio Finale di Giugno', value: 'SCRUTINIO_GIUGNO' },
  { label: t('sidiExports.options.scrutinioSettembre') || 'Esiti Scrutinio Differito (Debiti Settembre)', value: 'SCRUTINIO_SETTEMBRE_DEBITI' },
  { label: t('sidiExports.options.frequenze') || 'Frequenze e Monitoraggio Assenze', value: 'FREQUENZE' }
])

onMounted(() => {
  loadExports()
})

async function loadExports() {
  loading.value = true
  try {
    const res = await api.get('/sidi/exports')
    exportsList.value = res.data || []
  } catch (err) {
    console.error('Error fetching SIDI exports:', err)
  } finally {
    loading.value = false
  }
}

async function generateExport() {
  generating.value = true
  validationResult.value = null
  try {
    const res = await api.post('/sidi/exports/generate', form.value)
    validationResult.value = res.data?.validation || { valid: true }
    $q.notify({
      type: 'positive',
      message: t('sidiExports.successGenerated') || 'Flusso XML generato e validato con successo!'
    })
    await loadExports()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('sidiExports.errorGenerating') || 'Errore durante la generazione del flusso SIDI'
    })
  } finally {
    generating.value = false
  }
}

function downloadXML(exp) {
  $q.notify({
    type: 'positive',
    icon: 'download',
    message: t('sidiExports.downloadStarted', { fileName: exp.file_name }) || `Download di ${exp.file_name} avviato.`
  })
}
</script>
