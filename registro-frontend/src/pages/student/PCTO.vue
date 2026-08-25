<template>
  <q-page class="q-pa-md bg-slate-50 min-h-screen">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="work_outline" color="primary" class="q-mr-sm" />
          PCTO & Percorsi per le Competenze Trasversali
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Esperienze formative e tirocinio scuola-lavoro
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-chip color="primary" text-color="white" icon="timer" class="text-weight-bold">
          Totale Svolto: {{ totalHours }} / {{ targetHours }} Ore
        </q-chip>
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchPCTO" />
      </div>
    </div>

    <!-- Progress Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-lg border border-slate-100">
      <q-card-section>
        <div class="row items-center justify-between q-mb-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800">
            Avanzamento Monte Ore Obbligatorio
          </div>
          <span class="text-caption text-weight-bold text-primary">
            {{ Math.round(progressValue * 100) }}% Completato
          </span>
        </div>
        <q-linear-progress
          size="20px"
          :value="progressValue"
          color="primary"
          stripe
          rounded
          class="rounded-lg"
        />
        <div class="row justify-between text-caption text-slate-500 q-mt-xs">
          <span>0 ore</span>
          <span>Target percorso: {{ targetHours }} ore</span>
        </div>
      </q-card-section>
    </q-card>

    <div class="row q-col-gutter-lg">
      <!-- Projects List -->
      <div class="col-12 col-md-8">
        <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">I Miei Progetti di PCTO</div>

        <div v-if="loading" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <div v-else-if="projects.length === 0" class="q-pa-xl text-center bg-white rounded-xl border border-slate-100 shadow-soft">
          <q-icon name="folder_off" size="64px" color="slate-300" class="q-mb-md" />
          <div class="text-h6 text-slate-700">Nessun progetto PCTO assegnato</div>
          <div class="text-caption text-slate-500">I progetti assegnati dalla segreteria o dal tuo tutor compariranno qui.</div>
        </div>

        <div v-else class="space-y-4">
          <q-card
            v-for="project in projects"
            :key="project.id"
            flat
            bordered
            class="rounded-xl bg-white shadow-soft border border-slate-100 overflow-hidden q-mb-md"
          >
            <q-card-section class="bg-slate-50 border-b border-slate-100">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-h6 text-weight-bold text-slate-800">{{ project.title }}</div>
                  <div class="text-caption text-slate-500 row items-center q-mt-xs">
                    <q-icon name="apartment" size="xs" class="q-mr-xs text-primary" />
                    <span class="text-weight-medium text-slate-700">{{ project.company_name || 'Ente / Struttura Convenzionata' }}</span>
                    <span class="q-mx-xs">&bull;</span>
                    <span>{{ project.type || 'Convenzione' }}</span>
                  </div>
                </div>
                <div class="text-right">
                  <q-chip
                    :color="project.status === 'Completed' ? 'positive' : 'indigo-1'"
                    :text-color="project.status === 'Completed' ? 'white' : 'indigo-9'"
                    size="sm"
                    class="font-bold"
                  >
                    {{ project.status || 'In Corso' }}
                  </q-chip>
                  <div class="text-caption text-weight-bold text-primary q-mt-xs">
                    {{ project.hours_completed || project.hours_done || 0 }} / {{ project.total_hours }} Ore
                  </div>
                </div>
              </div>
            </q-card-section>

            <q-card-section class="q-pa-md">
              <div class="row q-col-gutter-md text-caption text-slate-600">
                <div class="col-12 col-sm-6">
                  <div class="text-slate-400 text-weight-medium">Tutor / Referente:</div>
                  <div class="text-slate-800 text-weight-bold q-mt-xs">
                    {{ project.company_tutor_name || project.tutor_name || 'Da definire' }}
                  </div>
                </div>
                <div class="col-12 col-sm-6">
                  <div class="text-slate-400 text-weight-medium">Periodo Svolgimento:</div>
                  <div class="text-slate-800 font-medium q-mt-xs">
                    {{ formatDate(project.start_date) }} &ndash; {{ formatDate(project.end_date) }}
                  </div>
                </div>
              </div>

              <div v-if="project.description" class="q-mt-md text-caption text-slate-600">
                <div class="text-slate-400 text-weight-medium">Descrizione Attività:</div>
                <p class="q-mt-xs q-mb-none">{{ project.description }}</p>
              </div>
            </q-card-section>

            <q-separator />

            <q-card-actions align="between" class="bg-slate-50 q-pa-sm">
              <q-btn
                flat
                icon="list_alt"
                label="Registro Attività"
                color="primary"
                no-caps
                @click="openDetails(project)"
              />
              <q-btn
                unelevated
                icon="add_circle"
                label="Registra Ore"
                color="primary"
                no-caps
                class="rounded-lg"
                @click="openLogDialog(project)"
              />
            </q-card-actions>
          </q-card>
        </div>
      </div>

      <!-- Information & Guidelines Sidebar -->
      <div class="col-12 col-md-4">
        <!-- Riferimenti Normativi -->
        <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-md border border-slate-100">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-slate-800">Guida & Normativa PCTO</div>
          </q-card-section>
          <q-list separator>
            <q-item clickable v-ripple href="https://www.istruzione.it/pcto-orientamento/" target="_blank">
              <q-item-section avatar><q-icon name="public" color="primary" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">Portale Ministeriale PCTO</q-item-label>
                <q-item-label caption>Linee guida e modulistica nazionale</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="open_in_new" size="xs" /></q-item-section>
            </q-item>
            <q-item clickable v-ripple href="https://www.inail.it/" target="_blank">
              <q-item-section avatar><q-icon name="security" color="teal" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">Sicurezza sul Lavoro</q-item-label>
                <q-item-label caption>Corso obbligatorio D.Lgs 81/08</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="open_in_new" size="xs" /></q-item-section>
            </q-item>
          </q-list>
        </q-card>

        <!-- Riepilogo Regole -->
        <q-card flat bordered class="rounded-xl bg-blue-50 border border-blue-200 text-blue-900 shadow-soft">
          <q-card-section>
            <div class="row items-center q-mb-xs">
              <q-icon name="info" size="sm" class="q-mr-xs" />
              <div class="text-subtitle2 text-weight-bold">Promemoria per lo Studente</div>
            </div>
            <p class="text-caption q-mb-none leading-relaxed">
              Le ore inserite nel diario di bordo devono essere approvate dal docente tutor PCTO per essere conteggiate validamente nel monte ore finale del diploma.
            </p>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Log Hours Dialog -->
    <q-dialog v-model="showLogDialog">
      <q-card style="min-width: 420px; max-width: 90vw;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">Registra Ore PCTO</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div>
            <div class="text-caption text-slate-500">Progetto:</div>
            <div class="text-weight-bold text-slate-800">{{ selectedProject?.title }}</div>
          </div>

          <q-input
            v-model="logForm.date"
            label="Data Svolgimento (YYYY-MM-DD)"
            type="date"
            outlined
            dense
          />
          <q-input
            v-model.number="logForm.hours"
            label="Numero di Ore Svolte"
            type="number"
            min="0.5"
            step="0.5"
            outlined
            dense
          />
          <q-input
            v-model="logForm.activity"
            label="Descrizione Attività Svolta"
            type="textarea"
            outlined
            dense
            rows="3"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup no-caps />
          <q-btn
            color="primary"
            label="Invia Registrazione"
            :loading="submittingLog"
            no-caps
            @click="submitHours"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Activity Details Dialog -->
    <q-dialog v-model="showDetailsDialog">
      <q-card style="min-width: 500px; max-width: 90vw;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-slate-800 text-white row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">Diario delle Attività</div>
            <div class="text-caption opacity-80">{{ selectedProject?.title }}</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div v-if="loadingDetails" class="text-center q-pa-lg">
            <q-spinner-dots color="primary" size="30px" />
          </div>
          <div v-else-if="projectLogs.length === 0" class="text-center text-slate-500 q-pa-lg">
            Nessuna ora ancora registrata per questo progetto.
          </div>
          <q-list v-else separator>
            <q-item v-for="log in projectLogs" :key="log.id">
              <q-item-section avatar>
                <q-avatar :color="log.verified ? 'green-1' : 'amber-1'" :text-color="log.verified ? 'positive' : 'warning'" icon="schedule" size="36px" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800">{{ log.activity_description }}</q-item-label>
                <q-item-label caption class="text-slate-500">{{ formatDate(log.date) }} &bull; {{ log.hours }} ore</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip :color="log.verified ? 'positive' : 'warning'" text-color="white" size="xs" class="font-bold">
                  {{ log.verified ? 'Verificato' : 'In Attesa' }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Chiudi" v-close-popup no-caps />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { pctoService } from '@/services/pctoService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const projects = ref([])
const loading = ref(false)
const targetHours = 90

const showLogDialog = ref(false)
const submittingLog = ref(false)
const selectedProject = ref(null)
const logForm = ref({
  date: new Date().toISOString().slice(0, 10),
  hours: 4,
  activity: ''
})

const showDetailsDialog = ref(false)
const loadingDetails = ref(false)
const projectLogs = ref([])

const totalHours = computed(() => {
  return projects.value.reduce((acc, p) => acc + (p.hours_completed || p.hours_done || 0), 0)
})

const progressValue = computed(() => {
  if (targetHours === 0) return 0
  return Math.min(totalHours.value / targetHours, 1)
})

const formatDate = (d) => {
  if (!d) return 'N/D'
  return new Date(d).toLocaleDateString('it-IT')
}

onMounted(() => {
  fetchPCTO()
})

async function fetchPCTO() {
  loading.value = true
  try {
    const res = await pctoService.getMyProjects()
    projects.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ message: 'Errore nel caricamento dei progetti PCTO', color: 'negative' })
  } finally {
    loading.value = false
  }
}

function openLogDialog(project) {
  selectedProject.value = project
  logForm.value = {
    date: new Date().toISOString().slice(0, 10),
    hours: 4,
    activity: ''
  }
  showLogDialog.value = true
}

async function submitHours() {
  if (!logForm.value.hours || logForm.value.hours <= 0) {
    $q.notify({ type: 'warning', message: 'Inserisci un numero di ore valido' })
    return
  }
  if (!logForm.value.activity) {
    $q.notify({ type: 'warning', message: 'Inserisci la descrizione dell\'attività svolta' })
    return
  }

  submittingLog.value = true
  try {
    await pctoService.logHours({
      project_id: selectedProject.value.id,
      date: logForm.value.date,
      hours: parseFloat(logForm.value.hours),
      activity: logForm.value.activity
    })
    $q.notify({ type: 'positive', message: 'Ore registrate con successo!' })
    showLogDialog.value = false
    await fetchPCTO()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore nella registrazione delle ore' })
  } finally {
    submittingLog.value = false
  }
}

async function openDetails(project) {
  selectedProject.value = project
  showDetailsDialog.value = true
  loadingDetails.value = true
  try {
    const res = await pctoService.getProjectDetails(project.id)
    projectLogs.value = res.data?.logs || []
  } catch (e) {
    projectLogs.value = []
  } finally {
    loadingDetails.value = false
  }
}
</script>

