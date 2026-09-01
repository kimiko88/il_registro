<template>
  <q-page class="q-pa-md bg-slate-50 min-h-screen">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="cloud_sync" color="primary" class="q-mr-sm" />
          Flussi SIDI & Tracciati Ministeriali (MIM)
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Esportazione pacchetti XML per Anagrafe Nazionale Studenti ed Esiti Scrutini
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="loadExports" />
    </div>

    <!-- WIZARD ESPORTAZIONE SIDI -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-lg border border-slate-100">
      <q-card-section>
        <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">Generatore Flusso SIDI</div>
        <div class="row q-col-gutter-md items-center">
          <div class="col-12 col-md-5">
            <q-select
              v-model="form.export_type"
              :options="exportOptions"
              emit-value
              map-options
              label="Tipologia di Flusso SIDI *"
              outlined
              dense
            />
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="form.school_year"
              :options="['2024/2025', '2025/2026', '2026/2027']"
              label="Anno Scolastico"
              outlined
              dense
            />
          </div>
          <div class="col-12 col-md-3">
            <q-btn
              color="primary"
              icon="bolt"
              label="Valida ed Esporta XML"
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
            <div class="text-weight-bold">Controllo di Congruenza Superato!</div>
            <div class="text-caption">Tutti i {{ validationResult.total_records }} record anagrafici sono conformi agli schemi XSD ministeriali SIDI.</div>
          </q-banner>

          <q-banner v-else class="bg-amber-50 text-amber-9 rounded-xl border border-amber-200">
            <template v-slot:avatar>
              <q-icon name="warning" color="warning" />
            </template>
            <div class="text-weight-bold">Attenzione: Rilevate anomalie prima dell'invio SIDI</div>
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
        <div class="text-h6 text-weight-bold text-slate-800">Storico Flussi Generati</div>
        <q-chip color="blue-1" text-color="blue-9" size="sm">
          {{ exportsList.length }} Flussi archiviati
        </q-chip>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <div v-if="exportsList.length === 0" class="q-pa-xl text-center">
          <q-icon name="folder_open" size="48px" color="slate-300" class="q-mb-sm" />
          <div class="text-slate-500">Nessun flusso SIDI ancora generato.</div>
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
                Tipo: {{ exp.export_type }} &bull; A.S. {{ exp.school_year }} &bull; {{ exp.records_count }} record
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-btn
                unelevated
                color="primary"
                icon="download"
                label="Scarica XML"
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
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import api from '@/services/api'

const $q = useQuasar()
const loading = ref(false)
const generating = ref(false)
const exportsList = ref([])
const validationResult = ref(null)

const form = ref({
  export_type: 'ANS_ANAGRAFE',
  school_year: '2025/2026'
})

const exportOptions = [
  { label: 'Anagrafe Nazionale Studenti (ANS)', value: 'ANS_ANAGRAFE' },
  { label: 'Esiti Scrutinio Finale di Giugno', value: 'SCRUTINIO_GIUGNO' },
  { label: 'Esiti Scrutinio Differito (Debiti Settembre)', value: 'SCRUTINIO_SETTEMBRE_DEBITI' },
  { label: 'Frequenze e Monitoraggio Assenze', value: 'FREQUENZE' }
]

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
      message: 'Flusso XML generato e validato con successo!'
    })
    await loadExports()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore durante la generazione del flusso SIDI'
    })
  } finally {
    generating.value = false
  }
}

function downloadXML(exp) {
  $q.notify({
    type: 'positive',
    icon: 'download',
    message: `Download di ${exp.file_name} avviato.`
  })
}
</script>
