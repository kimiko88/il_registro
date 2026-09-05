<template>
  <q-page class="q-pa-md bg-slate-50 min-h-screen">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="work_outline" color="primary" class="q-mr-sm" />
          {{ t('studentPcto.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('studentPcto.subtitle') }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-chip color="primary" text-color="white" icon="timer" class="text-weight-bold">
          {{ t('studentPcto.totalDone') }}: {{ totalHours }} / {{ targetHours }} {{ t('studentPcto.hours') }}
        </q-chip>
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchPCTO" />
      </div>
    </div>

    <!-- Progress Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-lg border border-slate-100">
      <q-card-section>
        <div class="row items-center justify-between q-mb-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800">
            {{ t('studentPcto.progressTitle') }}
          </div>
          <span class="text-caption text-weight-bold text-primary">
            {{ Math.round(progressValue * 100) }}% {{ t('studentPcto.completed') }}
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
          <span>{{ t('studentPcto.zeroHours') }}</span>
          <span>{{ t('studentPcto.targetPath', { target: targetHours }) }}</span>
        </div>
      </q-card-section>
    </q-card>

    <div class="row q-col-gutter-lg">
      <!-- Projects List -->
      <div class="col-12 col-md-8">
        <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">{{ t('studentPcto.myProjects') }}</div>

        <!-- Skeleton Loader for Projects -->
        <div v-if="loading" class="space-y-4" role="status" :aria-label="t('common.loading') || 'Caricamento in corso'">
          <q-card v-for="n in 2" :key="n" flat bordered class="rounded-xl bg-white shadow-soft q-mb-md border border-slate-100 q-pa-md">
            <q-skeleton type="rect" height="32px" class="q-mb-sm" />
            <q-skeleton type="text" width="60%" class="q-mb-md" />
            <div class="row q-col-gutter-md">
              <div class="col-6"><q-skeleton type="text" /></div>
              <div class="col-6"><q-skeleton type="text" /></div>
            </div>
          </q-card>
        </div>

        <div v-else-if="projects.length === 0" class="q-pa-xl text-center bg-white rounded-xl border border-slate-100 shadow-soft">
          <q-icon name="folder_off" size="64px" color="slate-300" class="q-mb-md" />
          <div class="text-h6 text-slate-700">{{ t('studentPcto.noProjects') }}</div>
          <div class="text-caption text-slate-500">{{ t('studentPcto.noProjectsDesc') }}</div>
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
                    <span class="text-weight-medium text-slate-700">{{ project.company_name || t('studentPcto.companyPlaceholder') }}</span>
                    <span class="q-mx-xs">&bull;</span>
                    <span>{{ project.type || t('studentPcto.convention') }}</span>
                  </div>
                </div>
                <div class="text-right">
                  <q-chip
                    :color="project.status === 'Completed' ? 'positive' : 'indigo-1'"
                    :text-color="project.status === 'Completed' ? 'white' : 'indigo-9'"
                    size="sm"
                    class="font-bold"
                  >
                    {{ project.status === 'Completed' ? (t('common.completed') || 'Completato') : (project.status || t('studentPcto.inProgress')) }}
                  </q-chip>
                  <div class="text-caption text-weight-bold text-primary q-mt-xs">
                    {{ project.hours_completed || project.hours_done || 0 }} / {{ project.total_hours }} {{ t('studentPcto.hours') }}
                  </div>
                </div>
              </div>
            </q-card-section>

            <q-card-section class="q-pa-md">
              <div class="row q-col-gutter-md text-caption text-slate-600">
                <div class="col-12 col-sm-6">
                  <div class="text-slate-400 text-weight-medium">{{ t('studentPcto.tutor') }}</div>
                  <div class="text-slate-800 text-weight-bold q-mt-xs">
                    {{ project.company_tutor_name || project.tutor_name || t('studentPcto.toDefine') }}
                  </div>
                </div>
                <div class="col-12 col-sm-6">
                  <div class="text-slate-400 text-weight-medium">{{ t('studentPcto.period') }}</div>
                  <div class="text-slate-800 font-medium q-mt-xs">
                    {{ formatDate(project.start_date) }} &ndash; {{ formatDate(project.end_date) }}
                  </div>
                </div>
              </div>

              <div v-if="project.description" class="q-mt-md text-caption text-slate-600">
                <div class="text-slate-400 text-weight-medium">{{ t('studentPcto.activityDesc') }}</div>
                <p class="q-mt-xs q-mb-none">{{ project.description }}</p>
              </div>
            </q-card-section>

            <q-separator />

            <q-card-actions align="between" class="bg-slate-50 q-pa-sm">
              <q-btn
                flat
                icon="list_alt"
                :label="t('studentPcto.activityLog')"
                color="primary"
                no-caps
                @click="openDetails(project)"
              />
              <q-btn
                unelevated
                icon="add_circle"
                :label="t('studentPcto.logHoursBtn')"
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
            <div class="text-subtitle1 text-weight-bold text-slate-800">{{ t('studentPcto.guideTitle') }}</div>
          </q-card-section>
          <q-list separator>
            <q-item clickable v-ripple href="https://www.istruzione.it/pcto-orientamento/" target="_blank" rel="noopener noreferrer">
              <q-item-section avatar><q-icon name="public" color="primary" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">{{ t('studentPcto.ministryPortal') }}</q-item-label>
                <q-item-label caption>{{ t('studentPcto.ministryPortalDesc') }}</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="open_in_new" size="xs" /></q-item-section>
            </q-item>
            <q-item clickable v-ripple href="https://www.inail.it/" target="_blank" rel="noopener noreferrer">
              <q-item-section avatar><q-icon name="security" color="teal" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">{{ t('studentPcto.safetyAtWork') }}</q-item-label>
                <q-item-label caption>{{ t('studentPcto.safetyAtWorkDesc') }}</q-item-label>
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
              <div class="text-subtitle2 text-weight-bold">{{ t('studentPcto.studentReminderTitle') }}</div>
            </div>
            <p class="text-caption q-mb-none leading-relaxed">
              {{ t('studentPcto.studentReminderText') }}
            </p>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Log Hours Dialog -->
    <q-dialog v-model="showLogDialog">
      <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">{{ t('studentPcto.logModalTitle') }}</div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div>
            <div class="text-caption text-slate-500">{{ t('studentPcto.projectLabel') }}</div>
            <div class="text-weight-bold text-slate-800">{{ selectedProject?.title }}</div>
          </div>

          <q-input
            v-model="logForm.date"
            :label="t('studentPcto.dateLabel')"
            type="date"
            outlined
            dense
          />
          <q-input
            v-model.number="logForm.hours"
            :label="t('studentPcto.hoursCount')"
            type="number"
            min="0.5"
            step="0.5"
            outlined
            dense
          />
          <q-input
            v-model="logForm.activity"
            :label="t('studentPcto.activityDetailLabel')"
            type="textarea"
            outlined
            dense
            rows="3"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup no-caps />
          <q-btn
            color="primary"
            :label="t('studentPcto.submitLogBtn')"
            :loading="submittingLog"
            no-caps
            @click="submitHours"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Activity Details Dialog -->
    <q-dialog v-model="showDetailsDialog">
      <q-card style="width: min(600px, 95vw); max-width: 95vw;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-slate-800 text-white row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">{{ t('studentPcto.activityDiaryTitle') }}</div>
            <div class="text-caption opacity-80">{{ selectedProject?.title }}</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div v-if="loadingDetails" class="q-pa-lg">
            <q-skeleton type="rect" height="40px" class="q-mb-sm" />
            <q-skeleton type="rect" height="40px" class="q-mb-sm" />
          </div>
          <div v-else-if="projectLogs.length === 0" class="text-center text-slate-500 q-pa-lg">
            {{ t('studentPcto.noHoursYet') }}
          </div>
          <q-list v-else separator>
            <q-item v-for="log in projectLogs" :key="log.id">
              <q-item-section avatar>
                <q-avatar :color="log.verified ? 'green-1' : 'amber-1'" :text-color="log.verified ? 'positive' : 'warning'" icon="schedule" size="36px" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800">{{ log.activity_description }}</q-item-label>
                <q-item-label caption class="text-slate-500">{{ formatDate(log.date) }} &bull; {{ log.hours }} {{ t('studentPcto.hours') }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip :color="log.verified ? 'positive' : 'warning'" text-color="white" size="xs" class="font-bold">
                  {{ log.verified ? t('studentPcto.verified') : t('studentPcto.pending') }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat :label="t('common.close') || 'Chiudi'" v-close-popup no-caps />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { pctoService } from '@/services/pctoService'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'

const $q = useQuasar()
const { t } = useI18n()
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
  return new Date(d).toLocaleDateString()
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
    $q.notify({ message: t('common.error') || 'Errore nel caricamento dei progetti PCTO', color: 'negative' })
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
    $q.notify({ type: 'warning', message: t('studentPcto.fillAllFields') || 'Inserisci un numero di ore valido' })
    return
  }
  if (!logForm.value.activity) {
    $q.notify({ type: 'warning', message: t('studentPcto.fillAllFields') || 'Inserisci la descrizione dell\'attività svolta' })
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
    $q.notify({ type: 'positive', message: t('studentPcto.logSuccess') || 'Ore registrate con successo!' })
    showLogDialog.value = false
    await fetchPCTO()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || t('studentPcto.logError') || 'Errore nella registrazione delle ore' })
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
