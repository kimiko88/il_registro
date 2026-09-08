<template>
  <q-page class="q-pa-md q-pa-lg-xl timecard-page">
    <!-- Hero Header -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="teal-8" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="schedule" size="14px" class="q-mr-xs" />
            {{ t('timecard.badge') || 'GESTIONE ORARIO • CCNL SCUOLA 36H' }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ selectedMonthLabel }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="calendar_month" color="teal-8" class="q-mr-sm" size="36px" />
          {{ t('timecard.title') || 'Cartellino & Piano Ferie' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('timecard.subtitle') || 'Riepilogo ore lavorate, timbrature badge, saldo straordinari e gestione istanze ferie/permessi' }}
        </div>
      </div>

      <!-- Controls & Actions -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          color="teal-8"
          icon="add"
          :label="t('timecard.requestLeaveBtn') || 'Richiedi Ferie / Permesso'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewLeaveDialog"
        />
        <q-btn
          outline
          color="primary"
          icon="download"
          :label="t('timecard.exportCsv') || 'Esporta CSV'"
          no-caps
          rounded
          class="bg-white"
          @click="exportCsv"
        />
        <q-btn
          flat
          round
          dense
          color="primary"
          icon="refresh"
          :loading="loading"
          @click="loadCurrentTab"
        >
          <q-tooltip>{{ t('common.refresh') || 'Aggiorna' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      class="bg-white rounded-xl shadow-1 text-slate-600 q-mb-lg"
      active-color="teal-8"
      indicator-color="teal-8"
      align="left"
      narrow-indicator
      @update:model-value="loadCurrentTab"
    >
      <q-tab name="cartellino" icon="access_time" :label="t('timecard.tabCartellino') || 'Cartellino Mensile'" />
      <q-tab name="ferie" icon="beach_access" :label="t('timecard.tabFerie') || 'Piano Ferie & Permessi'" />
      <q-tab v-if="isDSGAOrAdmin" name="dsga_overview" icon="badge" :label="t('timecard.tabOverview') || 'Riepilogo Personale (DSGA)'" />
    </q-tabs>

    <!-- TAB PANELS -->
    <q-tab-panels v-model="activeTab" animated class="bg-transparent">
      <!-- PANEL 1: CARTELLINO MENSILE -->
      <q-tab-panel name="cartellino" class="q-pa-none">
        <!-- Month & User Selector Bar -->
        <div class="row items-center q-col-gutter-md q-mb-lg">
          <div class="col-12 col-sm-4 col-md-3">
            <q-input
              v-model="selectedMonth"
              type="month"
              outlined
              dense
              bg-color="white"
              :label="t('timecard.selectMonth') || 'Mese di riferimento'"
              @update:model-value="loadTimecard"
            />
          </div>
          <div v-if="isDSGAOrAdmin" class="col-12 col-sm-6 col-md-4">
            <q-select
              v-model="selectedStaffUserId"
              :options="staffOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              bg-color="white"
              :label="t('timecard.selectStaff') || 'Dipendente'"
              @update:model-value="loadTimecard"
            />
          </div>
        </div>

        <!-- KPI Cards Summary -->
        <div class="row q-col-gutter-md q-mb-lg">
          <div class="col-6 col-md-3">
            <q-card class="stat-card bg-teal-50 border-teal-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-teal-100 text-teal-800">
                  <q-icon name="timelapse" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-teal-800 text-weight-bold">Ore Lavorate</div>
                  <div class="text-h5 text-weight-bolder text-teal-900">{{ timecardData.worked_hours || 0 }} h</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6 col-md-3">
            <q-card class="stat-card bg-blue-50 border-blue-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-blue-100 text-blue-800">
                  <q-icon name="schedule" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-blue-800 text-weight-bold">Ore Contrattuali</div>
                  <div class="text-h5 text-weight-bolder text-blue-900">{{ timecardData.contract_hours || 156 }} h</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6 col-md-3">
            <q-card class="stat-card" :class="overtimeBalance >= 0 ? 'bg-emerald-50 border-emerald-200' : 'bg-orange-50 border-orange-200'">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon" :class="overtimeBalance >= 0 ? 'bg-emerald-100 text-emerald-800' : 'bg-orange-100 text-orange-800'">
                  <q-icon :name="overtimeBalance >= 0 ? 'trending_up' : 'trending_down'" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-weight-bold" :class="overtimeBalance >= 0 ? 'text-emerald-800' : 'text-orange-800'">
                    Saldo / Straordinario
                  </div>
                  <div class="text-h5 text-weight-bolder" :class="overtimeBalance >= 0 ? 'text-emerald-900' : 'text-orange-900'">
                    {{ overtimeBalance >= 0 ? '+' : '' }}{{ overtimeBalance.toFixed(1) }} h
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-6 col-md-3">
            <q-card class="stat-card bg-amber-50 border-amber-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-amber-100 text-amber-800">
                  <q-icon name="beach_access" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-amber-800 text-weight-bold">Ferie / Assenze Mese</div>
                  <div class="text-h5 text-weight-bolder text-amber-900">
                    {{ (timecardData.leave_days || 0) + (timecardData.absence_days || 0) }} gg
                  </div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>

        <!-- Dettaglio Giornaliero Cartellino -->
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-card-section class="border-b border-slate-100 row items-center justify-between">
            <div class="text-subtitle1 text-weight-bold text-slate-800">
              Dettaglio Giornaliero Timbrature Badge
            </div>
            <div class="text-caption text-slate-400">
              CCNL Scuola: 36h/settimana • 7.2h medie su 5gg
            </div>
          </q-card-section>

          <q-table
            :rows="dailyEntries"
            :columns="dailyColumns"
            row-key="date"
            :loading="loading"
            class="no-shadow"
            :pagination="{ rowsPerPage: 31 }"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge :color="getEntryStatusColor(props.row.status)" class="q-px-sm q-py-xs">
                  {{ props.row.status || 'Presente' }}
                </q-badge>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 2: PIANO FERIE & PERMESSI -->
      <q-tab-panel name="ferie" class="q-pa-none">
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-card-section class="border-b border-slate-100 row items-center justify-between">
            <div class="text-subtitle1 text-weight-bold text-slate-800">
              Richieste di Ferie, Permessi e Malattia
            </div>
            <q-btn
              color="teal-8"
              icon="add"
              label="Nuova Istanza"
              no-caps
              dense
              rounded
              class="q-px-md"
              @click="openNewLeaveDialog"
            />
          </q-card-section>

          <q-table
            :rows="leavesList"
            :columns="leaveColumns"
            row-key="id"
            :loading="loading"
            no-data-label="Nessuna richiesta ferie registrata"
            class="no-shadow"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-type="props">
              <q-td :props="props">
                <q-badge :color="getLeaveTypeColor(props.row.type)" class="text-weight-bold">
                  {{ getLeaveTypeLabel(props.row.type) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge :color="getLeaveStatusColor(props.row.status)" class="q-px-sm q-py-xs">
                  {{ getLeaveStatusLabel(props.row.status) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" align="right">
                <!-- DSGA Approval buttons -->
                <div v-if="isDSGAOrAdmin && props.row.status === 'pending'" class="row items-center justify-end q-gutter-xs">
                  <q-btn
                    color="positive"
                    flat
                    dense
                    icon="check"
                    label="Approva"
                    size="sm"
                    no-caps
                    @click="approveLeave(props.row)"
                  />
                  <q-btn
                    color="negative"
                    flat
                    dense
                    icon="close"
                    label="Rifiuta"
                    size="sm"
                    no-caps
                    @click="rejectLeave(props.row)"
                  />
                </div>
                <!-- Delete user's own pending request -->
                <q-btn
                  v-else-if="props.row.status === 'pending' && props.row.user_id === user?.id"
                  color="grey-6"
                  flat
                  dense
                  icon="delete"
                  size="sm"
                  @click="deleteLeave(props.row)"
                />
                <span v-else class="text-caption text-slate-400">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 3: RIEPILOGO PERSONALE (DSGA) -->
      <q-tab-panel v-if="isDSGAOrAdmin" name="dsga_overview" class="q-pa-none">
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-card-section class="border-b border-slate-100 row items-center justify-between">
            <div class="text-subtitle1 text-weight-bold text-slate-800">
              Prospetto Personale ATA — Mese: {{ selectedMonth }}
            </div>
            <q-btn
              outline
              color="primary"
              icon="download"
              label="Esporta Tutti (CSV)"
              no-caps
              dense
              rounded
              class="q-px-md"
              @click="exportAllCsv"
            />
          </q-card-section>

          <q-table
            :rows="allTimecards"
            :columns="allTimecardColumns"
            row-key="user_id"
            :loading="loading"
            class="no-shadow"
            :pagination="{ rowsPerPage: 20 }"
          >
            <template v-slot:body-cell-balance="props">
              <q-td :props="props">
                <span :class="props.row.overtime_hours >= 0 ? 'text-positive text-weight-bold' : 'text-negative text-weight-bold'">
                  {{ props.row.overtime_hours >= 0 ? '+' : '' }}{{ Number(props.row.overtime_hours).toFixed(1) }} h
                </span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog: Nuova Richiesta Ferie / Permesso -->
    <q-dialog v-model="leaveDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="beach_access" color="teal-8" class="q-mr-sm" size="24px" />
            {{ t('timecard.dialogLeaveTitle') || 'Richiesta Assenza / Permesso' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Tipologia Istanza *</label>
            <q-select
              v-model="leaveForm.type"
              :options="leaveTypeOptions"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Data Inizio *</label>
              <q-input v-model="leaveForm.start_date" type="date" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Data Fine *</label>
              <q-input v-model="leaveForm.end_date" type="date" outlined dense />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Numero Giorni Lavorativi</label>
              <q-input v-model.number="leaveForm.days" type="number" step="0.5" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Ore (se permesso breve)</label>
              <q-input v-model.number="leaveForm.hours" type="number" step="0.5" outlined dense />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Note / Protocollo</label>
            <q-input v-model="leaveForm.notes" outlined dense type="textarea" rows="2" placeholder="Note per il DSGA o la segreteria..." />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="teal-8"
            :label="t('timecard.submitLeave') || 'Invia Richiesta'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="savingLeave"
            @click="submitLeave"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import staffAttendanceService from '@/services/staffAttendanceService'
import userService from '@/services/userService'

const $q = useQuasar()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)

const activeTab = ref('cartellino')
const loading = ref(false)
const savingLeave = ref(false)

const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)
const selectedStaffUserId = ref('')

const isDSGAOrAdmin = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'dsga' || r === 'admin' || r === 'superadmin'
})

const selectedMonthLabel = computed(() => {
  if (!selectedMonth.value) return ''
  const [y, m] = selectedMonth.value.split('-')
  const d = new Date(Number(y), Number(m) - 1, 1)
  return d.toLocaleDateString(locale.value || 'it-IT', { month: 'long', year: 'numeric' })
})

// Timecard
const timecardData = ref({})
const dailyEntries = ref([])
const overtimeBalance = computed(() => {
  return (timecardData.value.worked_hours || 0) - (timecardData.value.contract_hours || 156)
})

const dailyColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'entry_time', label: 'Entrata', field: 'entry_time', align: 'center' },
  { name: 'exit_time', label: 'Uscita', field: 'exit_time', align: 'center' },
  { name: 'hours', label: 'Ore Effettive', field: row => row.hours ? `${row.hours} h` : '—', align: 'center' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'notes', label: 'Note / Giustificativo', field: 'notes', align: 'left' }
]

// Leaves
const leavesList = ref([])
const leaveDialog = ref(false)
const leaveForm = ref({
  type: 'ferie',
  start_date: '',
  end_date: '',
  days: 1,
  hours: 0,
  notes: ''
})

const leaveTypeOptions = [
  { label: 'Ferie Ordinarie', value: 'ferie' },
  { label: 'Permesso Retribuito (motivi personali)', value: 'permesso' },
  { label: 'Permesso Breve (ore)', value: 'permesso_breve' },
  { label: 'Malattia', value: 'malattia' },
  { label: 'Diritto allo Studio (150 ore)', value: 'permesso_studio' },
  { label: 'Recupero Straordinario', value: 'recupero' }
]

const leaveColumns = [
  { name: 'user_name', label: 'Dipendente', field: 'user_name', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', align: 'center' },
  { name: 'period', label: 'Periodo', field: row => `${row.start_date} → ${row.end_date}`, align: 'left' },
  { name: 'days', label: 'Giorni / Ore', field: row => row.days ? `${row.days} gg` : `${row.hours} h`, align: 'center' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'notes', label: 'Note', field: 'notes', align: 'left' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

// DSGA Overview
const allTimecards = ref([])
const staffOptions = ref([])

const allTimecardColumns = [
  { name: 'user_name', label: 'Dipendente', field: 'user_name', align: 'left', sortable: true },
  { name: 'role', label: 'Ruolo', field: 'role', align: 'center' },
  { name: 'contract_hours', label: 'Ore Contrattuali', field: 'contract_hours', align: 'center' },
  { name: 'worked_hours', label: 'Ore Lavorate', field: 'worked_hours', align: 'center' },
  { name: 'balance', label: 'Saldo Straordinari', align: 'center' },
  { name: 'leave_days', label: 'Ferie Godute', field: 'leave_days', align: 'center' },
  { name: 'sick_days', label: 'Malattia', field: 'sick_days', align: 'center' }
]

async function loadCurrentTab() {
  if (activeTab.value === 'cartellino') {
    await loadTimecard()
  } else if (activeTab.value === 'ferie') {
    await loadLeaves()
  } else if (activeTab.value === 'dsga_overview') {
    await loadAllTimecards()
  }
}

async function loadTimecard() {
  loading.value = true
  try {
    const params = {
      month: selectedMonth.value
    }
    if (selectedStaffUserId.value && isDSGAOrAdmin.value) {
      params.user_id = selectedStaffUserId.value
    }
    const res = await staffAttendanceService.getTimecard(params)
    timecardData.value = res || {}
    dailyEntries.value = res.entries || []
  } catch (err) {
    // Non-blocking fallback
    timecardData.value = {
      worked_hours: 0,
      contract_hours: 156,
      overtime_hours: 0,
      leave_days: 0,
      absence_days: 0
    }
    dailyEntries.value = []
  } finally {
    loading.value = false
  }
}

async function loadLeaves() {
  loading.value = true
  try {
    const res = await staffAttendanceService.listLeaves({})
    leavesList.value = res || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento richieste ferie' })
  } finally {
    loading.value = false
  }
}

async function loadAllTimecards() {
  loading.value = true
  try {
    const res = await staffAttendanceService.getTimecard({ month: selectedMonth.value, all: true })
    allTimecards.value = Array.isArray(res) ? res : [res]
  } catch {
    allTimecards.value = []
  } finally {
    loading.value = false
  }
}

function openNewLeaveDialog() {
  leaveForm.value = {
    type: 'ferie',
    start_date: new Date().toISOString().substring(0, 10),
    end_date: new Date().toISOString().substring(0, 10),
    days: 1,
    hours: 0,
    notes: ''
  }
  leaveDialog.value = true
}

async function submitLeave() {
  if (!leaveForm.value.start_date || !leaveForm.value.end_date) {
    $q.notify({ type: 'warning', message: 'Date obbligatorie' })
    return
  }
  savingLeave.value = true
  try {
    await staffAttendanceService.createLeave(leaveForm.value)
    $q.notify({ type: 'positive', message: 'Richiesta inoltrata con successo al DSGA' })
    leaveDialog.value = false
    await loadLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore invio richiesta' })
  } finally {
    savingLeave.value = false
  }
}

async function approveLeave(row) {
  try {
    await staffAttendanceService.approveLeave(row.id, { notes: 'Approvato da DSGA' })
    $q.notify({ type: 'positive', message: 'Richiesta approvata' })
    await loadLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore approvazione' })
  }
}

async function rejectLeave(row) {
  try {
    await staffAttendanceService.rejectLeave(row.id, { reason: 'Rifiutato per esigenze di servizio' })
    $q.notify({ type: 'positive', message: 'Richiesta respinta' })
    await loadLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore rifiuto' })
  }
}

async function deleteLeave(row) {
  try {
    await staffAttendanceService.deleteLeave(row.id)
    $q.notify({ type: 'positive', message: 'Richiesta eliminata' })
    await loadLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore eliminazione richiesta' })
  }
}

async function exportCsv() {
  try {
    const blob = await staffAttendanceService.exportTimecard({
      month: selectedMonth.value,
      user_id: selectedStaffUserId.value
    })
    const url = window.URL.createObjectURL(new Blob([blob]))
    const a = document.createElement('a')
    a.href = url
    a.download = `cartellino_${selectedMonth.value}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    $q.notify({ type: 'negative', message: 'Errore esportazione CSV' })
  }
}

async function exportAllCsv() {
  try {
    const blob = await staffAttendanceService.exportTimecard({
      month: selectedMonth.value
    })
    const url = window.URL.createObjectURL(new Blob([blob]))
    const a = document.createElement('a')
    a.href = url
    a.download = `cartellino_personale_${selectedMonth.value}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    $q.notify({ type: 'negative', message: 'Errore esportazione CSV' })
  }
}

function getLeaveTypeLabel(t) {
  switch (t) {
    case 'ferie': return 'Ferie'
    case 'permesso': return 'Permesso Retribuito'
    case 'permesso_breve': return 'Permesso Breve'
    case 'malattia': return 'Malattia'
    case 'permesso_studio': return 'Diritto allo Studio'
    case 'recupero': return 'Recupero'
    default: return t
  }
}

function getLeaveTypeColor(t) {
  switch (t) {
    case 'ferie': return 'teal-8'
    case 'permesso': return 'blue-8'
    case 'permesso_breve': return 'indigo-8'
    case 'malattia': return 'negative'
    case 'permesso_studio': return 'purple-8'
    default: return 'grey-7'
  }
}

function getLeaveStatusLabel(s) {
  switch (s) {
    case 'pending': return 'In attesa'
    case 'approved': return 'Approvata'
    case 'rejected': return 'Respinta'
    default: return s
  }
}

function getLeaveStatusColor(s) {
  switch (s) {
    case 'pending': return 'amber-9'
    case 'approved': return 'positive'
    case 'rejected': return 'negative'
    default: return 'grey-7'
  }
}

function getEntryStatusColor(st) {
  switch (st) {
    case 'Presente': return 'positive'
    case 'Assente': return 'negative'
    case 'Ferie': return 'teal-7'
    case 'Permesso': return 'blue-7'
    default: return 'grey-7'
  }
}

onMounted(async () => {
  await loadTimecard()
  if (isDSGAOrAdmin.value) {
    try {
      const res = await userService.getUsers({})
      const raw = res.data?.users || res.data || []
      staffOptions.value = raw
        .filter(u => ['assistente_amministrativo', 'collaboratore_scolastico', 'collaboratore_ds', 'dsga'].includes(u.role))
        .map(u => ({
          id: u.id,
          name: `${u.last_name || ''} ${u.first_name || ''} (${u.role})`
        }))
    } catch {
      // Non blocking
    }
  }
})
</script>

<style scoped>
.timecard-page {
  background: #f8fafc;
  min-height: 100vh;
}

.stat-card {
  border-radius: 16px;
  border-width: 1px;
  border-style: solid;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
