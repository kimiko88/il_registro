<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('teacherHours.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('teacherHours.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          color="primary"
          icon="add_circle"
          :label="$t('teacherHours.logActivityTitle')"
          class="rounded-lg shadow-sm"
          @click="showLogDialog = true"
        />
        <q-btn
          outline
          color="secondary"
          icon="timer"
          :label="$t('teacherHours.requestLeaveBtn')"
          class="rounded-lg shadow-sm"
          @click="showLeaveDialog = true"
        />
      </div>
    </div>

    <!-- CARTELLINO MONTE ORE 40+40 (ART. 29 CCNL) -->
    <div class="row q-col-gutter-md q-mb-lg">
      <!-- 1. Collegio Docenti (40h) -->
      <div class="col-12 col-md-6">
        <q-card class="glass-card shadow-soft border-slate-100 full-height">
          <q-card-section>
            <div class="row items-center justify-between q-mb-sm">
              <div class="text-h6 text-weight-bold text-slate-800 row items-center q-gutter-x-sm">
                <q-icon name="groups" color="primary" size="md" />
                <span>{{ $t('teacherHours.collegioCard') }}</span>
              </div>
              <q-badge
                :color="obligation.collegio_done_hours > 40 ? 'warning' : 'positive'"
                class="q-pa-xs text-weight-bold text-caption"
              >
                {{ obligation.collegio_done_hours.toFixed(1) }} / {{ obligation.collegio_planned_hours.toFixed(1) }}h
              </q-badge>
            </div>
            <p class="text-caption text-slate-500 q-mb-md">
              {{ $t('teacherHours.collegioDesc') }}
            </p>

            <q-linear-progress
              :value="Math.min(obligation.collegio_done_hours / 40.0, 1.0)"
              size="12px"
              rounded
              color="primary"
              track-color="indigo-1"
              class="q-mb-md"
            />

            <div class="row text-center q-col-gutter-sm">
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.planned') }}</div>
                <div class="text-h6 text-weight-bold text-slate-700">40h</div>
              </div>
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.done') }}</div>
                <div class="text-h6 text-weight-bold text-primary">{{ obligation.collegio_done_hours.toFixed(1) }}h</div>
              </div>
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.remaining') }}</div>
                <div :class="['text-h6 text-weight-bold', obligation.collegio_remaining_hours < 0 ? 'text-warning' : 'text-positive']">
                  {{ obligation.collegio_remaining_hours.toFixed(1) }}h
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. Consigli di Classe (40h) -->
      <div class="col-12 col-md-6">
        <q-card class="glass-card shadow-soft border-slate-100 full-height">
          <q-card-section>
            <div class="row items-center justify-between q-mb-sm">
              <div class="text-h6 text-weight-bold text-slate-800 row items-center q-gutter-x-sm">
                <q-icon name="meeting_room" color="indigo-7" size="md" />
                <span>{{ $t('teacherHours.consigliCard') }}</span>
              </div>
              <q-badge
                :color="obligation.consigli_done_hours > 40 ? 'warning' : 'positive'"
                class="q-pa-xs text-weight-bold text-caption"
              >
                {{ obligation.consigli_done_hours.toFixed(1) }} / {{ obligation.consigli_planned_hours.toFixed(1) }}h
              </q-badge>
            </div>
            <p class="text-caption text-slate-500 q-mb-md">
              {{ $t('teacherHours.consigliDesc') }}
            </p>

            <q-linear-progress
              :value="Math.min(obligation.consigli_done_hours / 40.0, 1.0)"
              size="12px"
              rounded
              color="indigo-7"
              track-color="indigo-1"
              class="q-mb-md"
            />

            <div class="row text-center q-col-gutter-sm">
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.planned') }}</div>
                <div class="text-h6 text-weight-bold text-slate-700">40h</div>
              </div>
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.done') }}</div>
                <div class="text-h6 text-weight-bold text-indigo-7">{{ obligation.consigli_done_hours.toFixed(1) }}h</div>
              </div>
              <div class="col-4">
                <div class="text-caption text-slate-500">{{ $t('teacherHours.remaining') }}</div>
                <div :class="['text-h6 text-weight-bold', obligation.consigli_remaining_hours < 0 ? 'text-warning' : 'text-positive']">
                  {{ obligation.consigli_remaining_hours.toFixed(1) }}h
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- TABS: STORICO ATTIVITÀ & PERMESSI BREVI -->
    <q-card class="glass-card shadow-soft border-slate-100 overflow-hidden q-mb-xl">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-white"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="logs" icon="history" label="Storico Riunioni & Attività Svolte" />
        <q-tab name="leaves" icon="schedule" :label="$t('teacherHours.shortLeavesTitle')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        <!-- TAB 1: ATTIVITÀ REGISTRATE -->
        <q-tab-panel name="logs" class="q-pa-md">
          <q-table
            :rows="logs"
            :columns="logColumns"
            row-key="id"
            flat
            dense
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-category="props">
              <q-td :props="props">
                <q-chip
                  dense
                  :color="props.value === 'collegio' ? 'primary' : props.value === 'consigli' ? 'indigo-7' : 'teal'"
                  text-color="white"
                  class="text-weight-bold"
                >
                  {{ props.value.toUpperCase() }}
                </q-chip>
              </q-td>
            </template>
            <template v-slot:body-cell-hours="props">
              <q-td :props="props" class="text-weight-bolder">
                {{ props.value }}h
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- TAB 2: PERMESSI BREVI -->
        <q-tab-panel name="leaves" class="q-pa-md">
          <div class="q-mb-md bg-amber-50 border border-amber-200 rounded-borders q-pa-sm text-amber-9 text-caption row items-center">
            <q-icon name="info" size="sm" class="q-mr-sm" />
            <span>{{ $t('teacherHours.shortLeaveNotice') }}</span>
          </div>

          <q-table
            :rows="leaves"
            :columns="leaveColumns"
            row-key="id"
            flat
            dense
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge
                  :color="props.value === 'recovered' ? 'positive' : props.value === 'approved' ? 'primary' : 'warning'"
                  class="q-py-xs q-px-sm"
                >
                  {{ props.value }}
                </q-badge>
                <q-badge v-if="props.row.is_expired" color="negative" class="q-ml-xs">
                  Scaduto (2 Mesi)
                </q-badge>
              </q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn
                  v-if="props.row.status === 'approved' && !props.row.is_expired"
                  flat
                  dense
                  color="positive"
                  icon="check"
                  label="Recupera"
                  @click="openRecoverDialog(props.row)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- DIALOG REGISTRA ATTIVITÀ -->
    <q-dialog v-model="showLogDialog">
      <q-card style="min-width: 400px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">{{ $t('teacherHours.logActivityTitle') }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-select
            v-model="newLog.category"
            :options="['collegio', 'consigli', 'eccedenti', 'disposizione']"
            outlined
            dense
            :label="$t('teacherHours.category')"
          />
          <q-input
            v-model.number="newLog.hours"
            type="number"
            step="0.5"
            min="0.5"
            outlined
            dense
            :label="$t('teacherHours.hours')"
          />
          <q-input
            v-model="newLog.activity_date"
            type="date"
            outlined
            dense
            :label="$t('teacherHours.date')"
          />
          <q-input
            v-model="newLog.description"
            outlined
            dense
            :label="$t('teacherHours.description')"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="$t('common.save')" @click="saveActivity" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG RICHIEDI PERMESSO BREVE -->
    <q-dialog v-model="showLeaveDialog">
      <q-card style="min-width: 400px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">{{ $t('teacherHours.requestLeaveBtn') }}</div>
          <div class="text-caption text-slate-500">Massimo 2 ore giornaliere come previsto dal CCNL</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newLeave.request_date"
            type="date"
            outlined
            dense
            :label="$t('teacherHours.date')"
          />
          <q-input
            v-model.number="newLeave.hours"
            type="number"
            step="0.5"
            min="0.5"
            max="2.0"
            outlined
            dense
            :label="$t('teacherHours.hours')"
          />
          <q-input
            v-model="newLeave.reason"
            outlined
            dense
            label="Motivazione (es. Visita Medica, Motivi Personali)"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="$t('common.save')" @click="saveShortLeave" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG RECUPERO ORE -->
    <q-dialog v-model="showRecoverDialog">
      <q-card style="min-width: 350px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">{{ $t('teacherHours.recordRecovery') }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="recoverData.recovered_date"
            type="date"
            outlined
            dense
            :label="$t('teacherHours.date')"
          />
          <q-input
            v-model.number="recoverData.recovered_hours"
            type="number"
            step="0.5"
            min="0.5"
            outlined
            dense
            :label="$t('teacherHours.recoveredHours')"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="positive" :label="$t('common.save')" @click="confirmRecovery" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const activeTab = ref('logs')
const obligation = ref({
  collegio_planned_hours: 40.0,
  consigli_planned_hours: 40.0,
  collegio_done_hours: 0.0,
  collegio_remaining_hours: 40.0,
  consigli_done_hours: 0.0,
  consigli_remaining_hours: 40.0,
})
const logs = ref([])
const leaves = ref([])

const showLogDialog = ref(false)
const newLog = ref({
  category: 'collegio',
  hours: 2.0,
  activity_date: new Date().toISOString().substring(0, 10),
  description: ''
})

const showLeaveDialog = ref(false)
const newLeave = ref({
  request_date: new Date().toISOString().substring(0, 10),
  hours: 1.0,
  reason: ''
})

const showRecoverDialog = ref(false)
const selectedLeave = ref(null)
const recoverData = ref({
  recovered_date: new Date().toISOString().substring(0, 10),
  recovered_hours: 1.0
})

const logColumns = [
  { name: 'date', label: 'Data', field: 'activity_date', align: 'left', sortable: true },
  { name: 'category', label: 'Categoria', field: 'category', align: 'left' },
  { name: 'hours', label: 'Ore', field: 'hours', align: 'center' },
  { name: 'description', label: 'Descrizione', field: 'description', align: 'left' }
]

const leaveColumns = [
  { name: 'date', label: 'Data Permesso', field: 'request_date', align: 'left' },
  { name: 'hours', label: 'Ore', field: 'hours', align: 'center' },
  { name: 'reason', label: 'Motivo', field: 'reason', align: 'left' },
  { name: 'deadline', label: 'Scadenza Recupero (2 Mesi)', field: 'recovery_deadline', align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

async function loadData() {
  try {
    const res = await api.get('/teacher-hours/dashboard')
    if (res.data) {
      if (res.data.obligation) obligation.value = res.data.obligation
      if (res.data.logs) logs.value = res.data.logs
      if (res.data.leaves) leaves.value = res.data.leaves
    }
  } catch (err) {
    console.error('Error loading teacher hours:', err)
  }
}

async function saveActivity() {
  try {
    newLog.value.obligation_id = obligation.value.id
    await api.post('/teacher-hours/logs', newLog.value)
    $q.notify({ type: 'positive', message: 'Attività registrata con successo' })
    showLogDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nel salvataggio' })
  }
}

async function saveShortLeave() {
  try {
    await api.post('/teacher-hours/leaves', newLeave.value)
    $q.notify({ type: 'positive', message: 'Richiesta di permesso breve inviata' })
    showLeaveDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nella richiesta' })
  }
}

function openRecoverDialog(row) {
  selectedLeave.value = row
  recoverData.value.recovered_hours = row.hours
  showRecoverDialog.value = true
}

async function confirmRecovery() {
  if (!selectedLeave.value) return
  try {
    await api.post(`/teacher-hours/leaves/${selectedLeave.value.id}/recover`, recoverData.value)
    $q.notify({ type: 'positive', message: 'Recupero registrato' })
    showRecoverDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione del recupero' })
  }
}

onMounted(() => {
  loadData()
})
</script>
