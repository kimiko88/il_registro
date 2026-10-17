<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('competenciesDM14.invalsiTitle') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('competenciesDM14.invalsiSubtitle') }}
        </p>
      </div>

      <q-btn
        color="primary"
        icon="upload_file"
        :label="$t('competenciesDM14.importInvalsi')"
        class="rounded-lg shadow-sm"
        @click="showImportDialog = true"
      />
    </div>

    <!-- KPI CARDS -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100">
          <q-card-section>
            <div class="text-caption text-slate-500">Punteggio Medio Italiano</div>
            <div class="text-h4 text-weight-bolder text-primary q-my-xs">212.4</div>
            <q-badge color="positive">+12.4 vs Media Nazionale (200.0)</q-badge>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100">
          <q-card-section>
            <div class="text-caption text-slate-500">Punteggio Medio Matematica</div>
            <div class="text-h4 text-weight-bolder text-indigo-7 q-my-xs">208.1</div>
            <q-badge color="positive">+8.1 vs Media Nazionale (200.0)</q-badge>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100">
          <q-card-section>
            <div class="text-caption text-slate-500">Inglese Reading (B1/B2)</div>
            <div class="text-h4 text-weight-bolder text-teal-7 q-my-xs">218.0</div>
            <q-badge color="positive">+18.0 vs Media Nazionale (200.0)</q-badge>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100">
          <q-card-section>
            <div class="text-caption text-slate-500">Inglese Listening (B1/B2)</div>
            <div class="text-h4 text-weight-bolder text-amber-9 q-my-xs">205.5</div>
            <q-badge color="positive">+5.5 vs Media Nazionale (200.0)</q-badge>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- COMPARATIVE TABLE FOR RAV & PDM -->
    <q-card class="glass-card shadow-soft border-slate-100 overflow-hidden q-mb-xl">
      <q-card-section>
        <div class="text-h6 text-weight-bold q-mb-sm">Correlazione Esiti INVALSI vs Media Voti Consiglio</div>
        <div class="text-caption text-slate-500 q-mb-md">
          Strumento per l'autovalutazione di istituto (RAV), individuazione del cheating e calcolo dell'effetto scuola (PdM).
        </div>

        <q-table
          :rows="results"
          :columns="columns"
          row-key="id"
          flat
          dense
          :pagination="{ rowsPerPage: 10 }"
        >
          <template v-slot:body-cell-level="props">
            <q-td :props="props">
              <q-badge
                :color="props.value >= 4 ? 'positive' : props.value === 3 ? 'primary' : 'warning'"
                class="text-weight-bold"
              >
                Livello {{ props.value }} / 5
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-delta="props">
            <q-td :props="props">
              <span :class="props.value >= 0 ? 'text-positive text-weight-bold' : 'text-negative text-weight-bold'">
                {{ props.value >= 0 ? '+' : '' }}{{ props.value.toFixed(1) }}
              </span>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG IMPORT BATCH -->
    <q-dialog v-model="showImportDialog">
      <q-card style="min-width: 400px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">Importazione Esiti Ufficiali INVALSI</div>
          <div class="text-caption text-slate-500">Carica il tracciato restituito dall'INVALSI per le classi campione</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-file outlined dense label="Seleziona file CSV o Excel INVALSI" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" label="Importa Dati" @click="confirmImport" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const showImportDialog = ref(false)

const columns = [
  { name: 'grade_level', label: 'Grado', field: 'grade_level', align: 'left' },
  { name: 'subject', label: 'Disciplina', field: 'subject', align: 'left' },
  { name: 'score_wle', label: 'Punteggio WLE', field: 'score_wle', align: 'center' },
  { name: 'level', label: 'Livello INVALSI', field: 'level', align: 'center' },
  { name: 'national_mean', label: 'Media Nazionale', field: 'national_mean', align: 'center' },
  { name: 'delta', label: 'Scostamento Delta', field: 'score_delta_vs_national', align: 'center' },
  { name: 'report_grade', label: 'Media Voti Registro', field: 'report_class_grade', align: 'center' }
]

const results = ref([
  { id: '1', grade_level: 'Grado 8 (3ª Media)', subject: 'Italiano', score_wle: 215.0, level: 4, national_mean: 200.0, score_delta_vs_national: 15.0, report_class_grade: 7.8 },
  { id: '2', grade_level: 'Grado 8 (3ª Media)', subject: 'Matematica', score_wle: 208.5, level: 3, national_mean: 200.0, score_delta_vs_national: 8.5, report_class_grade: 7.2 },
  { id: '3', grade_level: 'Grado 8 (3ª Media)', subject: 'Inglese Reading', score_wle: 220.0, level: 4, national_mean: 200.0, score_delta_vs_national: 20.0, report_class_grade: 8.4 },
  { id: '4', grade_level: 'Grado 8 (3ª Media)', subject: 'Inglese Listening', score_wle: 206.0, level: 3, national_mean: 200.0, score_delta_vs_national: 6.0, report_class_grade: 7.5 },
  { id: '5', grade_level: 'Grado 10 (2ª Superiore)', subject: 'Italiano', score_wle: 210.0, level: 4, national_mean: 200.0, score_delta_vs_national: 10.0, report_class_grade: 7.0 },
  { id: '6', grade_level: 'Grado 10 (2ª Superiore)', subject: 'Matematica', score_wle: 204.0, level: 3, national_mean: 200.0, score_delta_vs_national: 4.0, report_class_grade: 6.8 }
])

function confirmImport() {
  $q.notify({ type: 'positive', message: 'Tracciato INVALSI importato e riconciliato con i registri' })
  showImportDialog.value = false
}
</script>
