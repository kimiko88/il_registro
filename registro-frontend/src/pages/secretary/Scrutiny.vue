<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="how_to_vote" color="primary" class="q-mr-sm" />
          Supervisione Scrutini
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Panoramica dello stato degli scrutini finali di tutte le classi
        </p>
      </div>

      <div class="row items-center q-gutter-md">
        <q-select
          v-model="selectedSchoolYear"
          :options="['2025/2026', '2024/2025']"
          label="Anno Scolastico"
          outlined
          dense
          class="bg-white"
          style="min-width: 160px"
        />

        <q-btn
          color="primary"
          unelevated
          icon="download"
          label="Esporta Tutti"
          class="rounded-lg"
          no-caps
          @click="exportAll"
        />
      </div>
    </div>

    <!-- Scrutiny Overview Table Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
        <div class="text-subtitle1 text-weight-bold text-slate-800">Stato Scrutini per Classe</div>
        <q-input v-model="filterText" dense placeholder="Cerca classe..." borderless class="bg-white q-px-sm rounded border">
          <template v-slot:append>
            <q-icon name="search" size="xs" />
          </template>
        </q-input>
      </q-card-section>

      <q-table
        flat
        :rows="filteredOverview"
        :columns="columns"
        row-key="class_id"
        :loading="scrutinyStore.loading"
        class="no-shadow"
      >
        <!-- Status Column -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge
              :color="getStatusBadgeColor(props.row.status)"
              class="text-weight-bold q-px-sm q-py-xs uppercase"
            >
              {{ getStatusLabel(props.row.status) }}
            </q-badge>
          </q-td>
        </template>

        <!-- Completed Subjects Column -->
        <template v-slot:body-cell-completed_subjects="props">
          <q-td :props="props" class="text-weight-bold">
            {{ props.row.completed_subjects }} / {{ props.row.total_subjects }}
          </q-td>
        </template>

        <!-- Actions Column -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-xs">
            <q-btn
              flat
              dense
              color="primary"
              icon="visibility"
              label="Dettaglio"
              no-caps
              @click="openDetail(props.row)"
            />
            <q-btn
              v-if="props.row.status === 'in_progress'"
              color="positive"
              unelevated
              dense
              size="sm"
              icon="lock"
              label="Finalizza"
              no-caps
              @click="confirmFinalize(props.row)"
            />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Class Report Detail Dialog -->
    <q-dialog v-model="detailDialog" maximized>
      <q-card class="bg-slate-50 column">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">
            Scrutinio Classe {{ selectedClass?.class_name }}
          </div>
          <div class="row items-center q-gutter-sm">
            <q-btn color="white" text-color="primary" icon="file_download" label="Esporta Excel" no-caps @click="exportClassExcel" />
            <q-btn icon="close" flat round dense v-close-popup />
          </div>
        </q-card-section>

        <q-card-section class="col scroll q-pa-md">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md">
            <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Tabella Voti Finale e Esiti</div>

            <div class="overflow-x-auto">
              <table class="w-full border-collapse text-left text-body2">
                <thead>
                  <tr class="bg-slate-100 border-b border-slate-300">
                    <th class="p-3 text-weight-bold">Studente</th>
                    <th class="p-3 text-weight-bold">Matematica</th>
                    <th class="p-3 text-weight-bold">Italiano</th>
                    <th class="p-3 text-weight-bold">Inglese</th>
                    <th class="p-3 text-weight-bold">Storia</th>
                    <th class="p-3 text-weight-bold text-center">Esito Finale</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="st in reportData?.students || []" :key="st.student_id" class="border-b border-slate-200 hover:bg-slate-50">
                    <td class="p-3 text-weight-bold text-slate-800">{{ st.name }}</td>
                    <td class="p-3">{{ st.grades?.Matematica || '-' }}</td>
                    <td class="p-3">{{ st.grades?.Italiano || '-' }}</td>
                    <td class="p-3">{{ st.grades?.Inglese || '-' }}</td>
                    <td class="p-3">{{ st.grades?.Storia || '-' }}</td>
                    <td class="p-3 text-center">
                      <q-chip
                        size="sm"
                        :color="st.outcome === 'Ammesso' ? 'positive' : st.outcome === 'Sospeso' ? 'warning' : 'negative'"
                        text-color="white"
                        class="text-weight-bold"
                      >
                        {{ st.outcome }}
                      </q-chip>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Footer Summary Counters -->
            <div class="row q-col-gutter-md border-t border-slate-200 q-pt-md q-mt-md text-center">
              <div class="col-4">
                <div class="bg-green-50 border border-green-200 rounded-lg p-2">
                  <div class="text-h6 text-weight-bold text-positive">{{ reportData?.admitted || 0 }}</div>
                  <div class="text-caption text-slate-600">Ammessi</div>
                </div>
              </div>
              <div class="col-4">
                <div class="bg-amber-50 border border-amber-200 rounded-lg p-2">
                  <div class="text-h6 text-weight-bold text-warning">{{ reportData?.suspended || 0 }}</div>
                  <div class="text-caption text-slate-600">Sospesi</div>
                </div>
              </div>
              <div class="col-4">
                <div class="bg-red-50 border border-red-200 rounded-lg p-2">
                  <div class="text-h6 text-weight-bold text-negative">{{ reportData?.rejected || 0 }}</div>
                  <div class="text-caption text-slate-600">Non Ammessi</div>
                </div>
              </div>
            </div>
          </q-card>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Finalize Confirmation Dialog -->
    <q-dialog v-model="finalizeModal">
      <q-card style="min-width: 400px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-warning text-white row items-center">
          <q-icon name="warning" size="md" class="q-mr-sm" />
          <div class="text-h6 text-weight-bold">Conferma Finalizzazione Scrutinio</div>
        </q-card-section>

        <q-card-section class="q-pa-md">
          <p class="text-body1 text-slate-700 q-mb-xs">
            Sei sicuro di voler finalizzare lo scrutinio della classe <strong>{{ classToFinalize?.class_name }}</strong>?
          </p>
          <p class="text-caption text-negative text-weight-bold q-my-none">
            Attenzione: questa azione è irreversibile. I docenti non potranno più modificare i voti.
          </p>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="positive"
            unelevated
            label="Conferma e Finalizza"
            :loading="finalizing"
            no-caps
            @click="executeFinalize"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useScrutinyStore } from '@/stores/scrutiny'

const $q = useQuasar()
const { t } = useI18n()
const scrutinyStore = useScrutinyStore()

const selectedSchoolYear = ref('2025/2026')
const filterText = ref('')

const detailDialog = ref(false)
const selectedClass = ref(null)
const reportData = ref(null)

const finalizeModal = ref(false)
const classToFinalize = ref(null)
const finalizing = ref(false)

const columns = [
  { name: 'class_name', label: 'Classe', align: 'left', field: 'class_name', sortable: true },
  { name: 'completed_subjects', label: 'Docenti Completati', align: 'center', field: 'completed_subjects' },
  { name: 'status', label: 'Stato Scrutinio', align: 'center', field: 'status', sortable: true },
  { name: 'last_updated', label: 'Ultimo Aggiornamento', align: 'center', field: 'last_updated' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const filteredOverview = computed(() => {
  if (!filterText.value) return scrutinyStore.overview
  return scrutinyStore.overview.filter(c =>
    c.class_name.toLowerCase().includes(filterText.value.toLowerCase())
  )
})

onMounted(() => {
  scrutinyStore.fetchOverview()
})

function getStatusBadgeColor(status) {
  switch (status) {
    case 'completed': return 'positive'
    case 'in_progress': return 'warning'
    case 'pending': default: return 'grey-7'
  }
}

function getStatusLabel(status) {
  switch (status) {
    case 'completed': return 'Completato'
    case 'in_progress': return 'In Corso'
    case 'pending': default: return 'In Attesa'
  }
}

async function openDetail(row) {
  selectedClass.value = row
  reportData.value = await scrutinyStore.fetchClassReport(row.class_id)
  detailDialog.value = true
}

function confirmFinalize(row) {
  classToFinalize.value = row
  finalizeModal.value = true
}

async function executeFinalize() {
  if (!classToFinalize.value) return
  finalizing.value = true
  try {
    await scrutinyStore.finalizeScrutiny(classToFinalize.value.class_id)
    $q.notify({ type: 'positive', message: 'Scrutinio finalizzato con successo' })
    finalizeModal.value = false
  } catch {
    $q.notify({ type: 'negative', message: 'Errore durante la finalizzazione' })
  } finally {
    finalizing.value = false
  }
}

function exportAll() {
  scrutinyStore.exportAll()
}

function exportClassExcel() {
  $q.notify({ type: 'info', message: 'Esportazione Excel avviata...' })
}
</script>
